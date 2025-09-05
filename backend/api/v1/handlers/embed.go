package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	models2 "steg/api/v1/models"
	"steg/api/v1/services"

	"github.com/google/uuid"
)

type EmbedJobCrud struct {
	service *services.EmbedJobService
}

func NewEmbedJobCrud(service *services.EmbedJobService) *EmbedJobCrud {
	return &EmbedJobCrud{service}
}

func (h *EmbedJobCrud) Create(ctx context.Context, reader io.Reader) (*models2.EmbedJob, error) {
	var embedRequest models2.CreateEmbedRequest
	err := json.NewDecoder(reader).Decode(&embedRequest)
	if err != nil {
		return nil, err
	}

	err = embedRequest.Validate()
	if err != nil {
		return nil, err
	}
	return h.service.CreateJob(ctx, embedRequest.ValidatedUUID)
}

func (h *EmbedJobCrud) Read(ctx context.Context, uuid uuid.UUID) (*models2.EmbedJob, bool, error) {
	job, err := h.service.GetJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return job, true, nil
}

func (h *EmbedJobCrud) List(ctx context.Context) ([]*models2.EmbedJob, error) {
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
