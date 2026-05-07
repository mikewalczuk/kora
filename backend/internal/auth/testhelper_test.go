package auth

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/impez/kora/internal/database"
	"github.com/impez/kora/internal/testutils"
)

const schema = `
	CREATE TABLE users (
		id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username      TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)
`

var (
	sharedPool *pgxpool.Pool
	sharedDB   *database.Queries
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var teardown func()
	sharedPool, teardown = testutils.NewPostgresPool(ctx, schema)
	sharedDB = database.New(sharedPool)

	code := m.Run()
	teardown()
	os.Exit(code)
}

func setupService(t *testing.T) *Service {
	t.Helper()
	t.Cleanup(func() {
		sharedPool.Exec(context.Background(), "TRUNCATE TABLE users")
	})
	return &Service{DB: sharedDB}
}

func seedUser(t *testing.T, username, password string) {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if _, err := sharedPool.Exec(context.Background(),
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		username, string(hash),
	); err != nil {
		t.Fatalf("seed user: %v", err)
	}
}
