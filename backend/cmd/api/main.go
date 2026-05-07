package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/impez/kora/internal/api"
	"github.com/impez/kora/internal/auth"
	"github.com/impez/kora/internal/config"
	"github.com/impez/kora/internal/database"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg.LogFormat)

	pool, err := pgxpool.New(context.Background(), cfg.DB)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	swagger, err := api.GetSpec()
	if err != nil {
		slog.Error("failed to load swagger spec", "error", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(oapimiddleware.OapiRequestValidator(swagger))
	r.Use(injectRequest)

	authHandler := &auth.Handler{
		Service: &auth.Service{DB: database.New(pool), JWTSecret: cfg.JWTSecret},
	}
	api.HandlerFromMux(api.NewStrictHandler(authHandler, nil), r)

	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("server starting", "addr", addr, "log_format", cfg.LogFormat)

	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
	}
}

func injectRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), auth.RequestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setupLogger(format string) {
    var lvl slog.Level
    opts := &slog.HandlerOptions{Level: lvl, AddSource: format == "text"}
    var h slog.Handler
    if format == "json" {
        h = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        h = slog.NewTextHandler(os.Stdout, opts)
    }
    slog.SetDefault(slog.New(h))
}
