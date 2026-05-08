package concepts

import (
	"context"
	"errors"

	"github.com/impez/kora/internal/api"
)

type Handler struct {
	Service *Service
}

func (h *Handler) ArchiveConcept(ctx context.Context, req api.ArchiveConceptRequestObject) (api.ArchiveConceptResponseObject, error) {
	err := h.Service.Archive(ctx, req.Id)
	if errors.Is(err, ErrNotFound) {
		return api.ArchiveConcept404Response{}, nil
	}
	if err != nil {
		return api.ArchiveConcept401Response{}, nil
	}
	return api.ArchiveConcept204Response{}, nil
}

func (h *Handler) RestoreConcept(ctx context.Context, req api.RestoreConceptRequestObject) (api.RestoreConceptResponseObject, error) {
	err := h.Service.Restore(ctx, req.Id)
	if errors.Is(err, ErrNotFound) {
		return api.RestoreConcept404Response{}, nil
	}
	if err != nil {
		return api.RestoreConcept401Response{}, nil
	}
	return api.RestoreConcept204Response{}, nil
}
