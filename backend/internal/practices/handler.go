package practices

import (
	"context"
	"errors"

	"github.com/impez/kora/internal/api"
)

type Handler struct {
	Service *Service
}

func (h *Handler) ListPractices(ctx context.Context, req api.ListPracticesRequestObject) (api.ListPracticesResponseObject, error) {
	in := ListInput{
		NoteID: req.Params.NoteId,
		Limit:  20,
		Page:   1,
	}
	if req.Params.Status != nil {
		in.Status = *req.Params.Status
	}
	if req.Params.Limit != nil {
		in.Limit = *req.Params.Limit
	}
	if req.Params.Page != nil {
		in.Page = *req.Params.Page
	}

	result, err := h.Service.List(ctx, in)
	if err != nil {
		return api.ListPractices401Response{}, nil
	}
	return api.ListPractices200JSONResponse(result), nil
}

func (h *Handler) CreatePractice(ctx context.Context, req api.CreatePracticeRequestObject) (api.CreatePracticeResponseObject, error) {
	id, err := h.Service.Create(ctx, CreateInput{NoteID: req.Body.NoteId})
	if errors.Is(err, ErrConflict) {
		return api.CreatePractice409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Id: id}}, nil
	}
	if err != nil {
		return api.CreatePractice401Response{}, nil
	}
	return api.CreatePractice201JSONResponse(api.CreatePracticeResponse{Id: id}), nil
}

func (h *Handler) GetPractice(ctx context.Context, req api.GetPracticeRequestObject) (api.GetPracticeResponseObject, error) {
	practice, err := h.Service.Get(ctx, req.Id)
	if errors.Is(err, ErrNotFound) {
		return api.GetPractice404Response{}, nil
	}
	if err != nil {
		return api.GetPractice401Response{}, nil
	}
	return api.GetPractice200JSONResponse(practice), nil
}
