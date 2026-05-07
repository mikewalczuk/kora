package notes

import (
	"context"
	"errors"

	"github.com/impez/kora/internal/api"
	"github.com/impez/kora/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	Service *Service
}

func (h *Handler) ListNotes(ctx context.Context, req api.ListNotesRequestObject) (api.ListNotesResponseObject, error) {
	page := int32(1)
	if req.Params.Page != nil {
		page = int32(*req.Params.Page)
	}
	limit := int32(20)
	if req.Params.Limit != nil {
		limit = int32(*req.Params.Limit)
	}

	result, err := h.Service.List(ctx, ListInput{Page: page, Limit: limit})
	if err != nil {
		return api.ListNotes401Response{}, nil
	}

	items := make([]api.Note, len(result.Notes))
	for i, n := range result.Notes {
		items[i] = dbNoteToAPI(n)
	}

	return api.ListNotes200JSONResponse{
		Items: items,
		Total: int(result.Total),
		Page:  int(page),
		Limit: int(limit),
	}, nil
}

func (h *Handler) CreateNote(ctx context.Context, req api.CreateNoteRequestObject) (api.CreateNoteResponseObject, error) {
	note, err := h.Service.Create(ctx, CreateInput{
		Title:   req.Body.Title,
		Content: req.Body.Content,
	})
	if err != nil {
		return api.CreateNote401Response{}, nil
	}

	return api.CreateNote201JSONResponse(dbNoteToAPI(note)), nil
}

func (h *Handler) GetNote(ctx context.Context, req api.GetNoteRequestObject) (api.GetNoteResponseObject, error) {
	id := pgtype.UUID{Bytes: req.Id, Valid: true}
	note, err := h.Service.Get(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return api.GetNote404Response{}, nil
	}
	if err != nil {
		return api.GetNote401Response{}, nil
	}

	return api.GetNote200JSONResponse(dbNoteToAPI(note)), nil
}

func (h *Handler) UpdateNote(ctx context.Context, req api.UpdateNoteRequestObject) (api.UpdateNoteResponseObject, error) {
	id := pgtype.UUID{Bytes: req.Id, Valid: true}
	note, err := h.Service.Update(ctx, id, UpdateInput{
		Title:   req.Body.Title,
		Content: req.Body.Content,
	})
	if errors.Is(err, ErrNotFound) {
		return api.UpdateNote404Response{}, nil
	}
	if err != nil {
		return api.UpdateNote401Response{}, nil
	}

	return api.UpdateNote200JSONResponse(dbNoteToAPI(note)), nil
}

func (h *Handler) DeleteNote(ctx context.Context, req api.DeleteNoteRequestObject) (api.DeleteNoteResponseObject, error) {
	id := pgtype.UUID{Bytes: req.Id, Valid: true}
	err := h.Service.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return api.DeleteNote404Response{}, nil
	}
	if err != nil {
		return api.DeleteNote401Response{}, nil
	}

	return api.DeleteNote204Response{}, nil
}

func dbNoteToAPI(n database.Note) api.Note {
	return api.Note{
		Id:        n.ID.Bytes,
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt.Time.UTC(),
		UpdatedAt: n.UpdatedAt.Time.UTC(),
	}
}
