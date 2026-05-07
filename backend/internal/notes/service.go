package notes

import (
	"context"
	"errors"

	"github.com/impez/kora/internal/auth"
	"github.com/impez/kora/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrNotFound = errors.New("note not found")

type Service struct {
	DB   *database.Queries
	Auth *auth.Service
}

type ListInput struct {
	Page  int32
	Limit int32
}

type ListResult struct {
	Notes []database.Note
	Total int64
}

type CreateInput struct {
	Title   string
	Content string
}

type UpdateInput struct {
	Title   *string
	Content *string
}

func (s *Service) currentUserID(ctx context.Context) (pgtype.UUID, error) {
	return s.Auth.CurrentUserID(ctx)
}

func (s *Service) List(ctx context.Context, in ListInput) (ListResult, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return ListResult{}, err
	}

	offset := (in.Page - 1) * in.Limit
	notes, err := s.DB.ListNotesByUser(ctx, database.ListNotesByUserParams{
		UserID: userID,
		Limit:  in.Limit,
		Offset: offset,
	})
	if err != nil {
		return ListResult{}, err
	}

	total, err := s.DB.CountNotesByUser(ctx, userID)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{Notes: notes, Total: total}, nil
}

func (s *Service) Create(ctx context.Context, in CreateInput) (database.Note, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return database.Note{}, err
	}

	return s.DB.CreateNote(ctx, database.CreateNoteParams{
		UserID:  userID,
		Title:   in.Title,
		Content: in.Content,
	})
}

func (s *Service) Get(ctx context.Context, id pgtype.UUID) (database.Note, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return database.Note{}, err
	}

	note, err := s.DB.GetNoteByID(ctx, database.GetNoteByIDParams{ID: id, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return database.Note{}, ErrNotFound
	}
	return note, err
}

func (s *Service) Update(ctx context.Context, id pgtype.UUID, in UpdateInput) (database.Note, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return database.Note{}, err
	}

	var title pgtype.Text
	if in.Title != nil {
		title = pgtype.Text{String: *in.Title, Valid: true}
	}
	var content pgtype.Text
	if in.Content != nil {
		content = pgtype.Text{String: *in.Content, Valid: true}
	}

	note, err := s.DB.UpdateNote(ctx, database.UpdateNoteParams{
		ID:      id,
		UserID:  userID,
		Title:   title,
		Content: content,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return database.Note{}, ErrNotFound
	}
	return note, err
}

func (s *Service) Delete(ctx context.Context, id pgtype.UUID) error {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return err
	}

	return s.DB.DeleteNote(ctx, database.DeleteNoteParams{ID: id, UserID: userID})
}
