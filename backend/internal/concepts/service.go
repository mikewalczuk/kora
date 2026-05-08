package concepts

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/impez/kora/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrNotFound = errors.New("concept not found")

type Service struct {
	DB *database.Queries
}

func (s *Service) Archive(ctx context.Context, id uuid.UUID) error {
	return s.DB.ArchiveConcept(ctx, database.ArchiveConceptParams{
		ID:     pgtype.UUID{Bytes: id, Valid: true},
		Reason: pgtype.Text{String: "user_archived", Valid: true},
	})
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) error {
	return s.DB.RestoreConcept(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
