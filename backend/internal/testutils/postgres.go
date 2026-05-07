package testutils

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type noopLogger struct{}

func (noopLogger) Printf(_ string, _ ...any) {}

func init() {
	tclog.SetDefault(noopLogger{})
}

// NewPostgresPool starts a Postgres container, applies schema, and returns the pool
// and a teardown function. Designed for use in TestMain where *testing.T is unavailable.
func NewPostgresPool(ctx context.Context, schema string) (*pgxpool.Pool, func()) {
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
		testcontainers.WithLogger(noopLogger{}),
	)
	if err != nil {
		panic("start postgres: " + err.Error())
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("connection string: " + err.Error())
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic("connect: " + err.Error())
	}

	if _, err := pool.Exec(ctx, schema); err != nil {
		panic("apply schema: " + err.Error())
	}

	teardown := func() {
		pool.Close()
		container.Terminate(ctx)
	}

	return pool, teardown
}

// StartPostgres starts a Postgres container for a single test, registering cleanup via t.Cleanup.
func StartPostgres(ctx context.Context, t *testing.T, schema string) *pgxpool.Pool {
	t.Helper()
	pool, teardown := NewPostgresPool(ctx, schema)
	t.Cleanup(teardown)
	return pool
}
