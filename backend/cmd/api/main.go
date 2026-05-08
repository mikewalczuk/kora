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
	"github.com/impez/kora/internal/notes"
	"github.com/impez/kora/internal/practices"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

// server satisfies api.StrictServerInterface by delegating to feature handlers.
type server struct {
	auth      *auth.Handler
	notes     *notes.Handler
	practices *practices.Handler
}

var _ api.StrictServerInterface = (*server)(nil)

func (s *server) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	return s.auth.Login(ctx, req)
}
func (s *server) Logout(ctx context.Context, req api.LogoutRequestObject) (api.LogoutResponseObject, error) {
	return s.auth.Logout(ctx, req)
}
func (s *server) GetMe(ctx context.Context, req api.GetMeRequestObject) (api.GetMeResponseObject, error) {
	return s.auth.GetMe(ctx, req)
}
func (s *server) ListNotes(ctx context.Context, req api.ListNotesRequestObject) (api.ListNotesResponseObject, error) {
	return s.notes.ListNotes(ctx, req)
}
func (s *server) CreateNote(ctx context.Context, req api.CreateNoteRequestObject) (api.CreateNoteResponseObject, error) {
	return s.notes.CreateNote(ctx, req)
}
func (s *server) GetNote(ctx context.Context, req api.GetNoteRequestObject) (api.GetNoteResponseObject, error) {
	return s.notes.GetNote(ctx, req)
}
func (s *server) UpdateNote(ctx context.Context, req api.UpdateNoteRequestObject) (api.UpdateNoteResponseObject, error) {
	return s.notes.UpdateNote(ctx, req)
}
func (s *server) DeleteNote(ctx context.Context, req api.DeleteNoteRequestObject) (api.DeleteNoteResponseObject, error) {
	return s.notes.DeleteNote(ctx, req)
}
func (s *server) CreatePractice(ctx context.Context, req api.CreatePracticeRequestObject) (api.CreatePracticeResponseObject, error) {
	return s.practices.CreatePractice(ctx, req)
}
func (s *server) GetPractice(ctx context.Context, req api.GetPracticeRequestObject) (api.GetPracticeResponseObject, error) {
	return s.practices.GetPractice(ctx, req)
}

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
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(oapimiddleware.OapiRequestValidator(swagger))
	r.Use(injectRequest)

	db := database.New(pool)
	authSvc := &auth.Service{DB: db, JWTSecret: cfg.JWTSecret}

	srv := &server{
		auth:      &auth.Handler{Service: authSvc},
		notes:     &notes.Handler{Service: &notes.Service{DB: db, Auth: authSvc}},
		practices: &practices.Handler{Service: &practices.Service{}},
	}
	api.HandlerFromMux(api.NewStrictHandler(srv, nil), r)

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
