package practices

import (
	"context"
	"errors"

	"github.com/impez/kora/internal/api"
)

type Handler struct {
	Service *Service
}

func (h *Handler) CreatePractice(ctx context.Context, req api.CreatePracticeRequestObject) (api.CreatePracticeResponseObject, error) {
	id, err := h.Service.Create(ctx, CreateInput{NoteID: req.Body.NoteId})
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
