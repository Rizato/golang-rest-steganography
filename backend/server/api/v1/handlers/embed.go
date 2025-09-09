package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"steg/server/api/v1/requests"
	"steg/shared/models"
	"steg/shared/services"

	"github.com/google/uuid"
)

type EmbedJobCrud struct {
	service *services.EmbedJobService
}

func NewEmbedJobCrud(service *services.EmbedJobService) *EmbedJobCrud {
	return &EmbedJobCrud{service}
}

func (h *EmbedJobCrud) Create(ctx context.Context, reader io.Reader) (*models.EmbedJob, error) {
	var embedRequest requests.CreateEmbedRequest
	err := json.NewDecoder(reader).Decode(&embedRequest)
	if err != nil {
		return nil, err
	}

	err = embedRequest.Validate()
	if err != nil {
		return nil, err
	}
	return h.service.CreateJob(ctx, embedRequest.ValidatedUUID, embedRequest.Message)
}

func (h *EmbedJobCrud) Read(ctx context.Context, uuid uuid.UUID) (*models.EmbedJob, bool, error) {
	job, err := h.service.GetJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return job, true, nil
}

func (h *EmbedJobCrud) List(ctx context.Context) ([]*models.EmbedJob, error) {
	return h.service.ListJobs(ctx)
}

func (h *EmbedJobCrud) Delete(ctx context.Context, uuid uuid.UUID) (bool, error) {
	err := h.service.DeleteJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
