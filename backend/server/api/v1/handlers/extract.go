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

type ExtractJobCrud struct {
	service *services.ExtractJobService
}

func NewExtractJobCrud(service *services.ExtractJobService) *ExtractJobCrud {
	return &ExtractJobCrud{service}
}

func (h *ExtractJobCrud) Create(ctx context.Context, reader io.Reader) (*models.ExtractJob, error) {
	var extractRequest requests.CreateExtractRequest
	err := json.NewDecoder(reader).Decode(&extractRequest)
	if err != nil {
		return nil, err
	}

	err = extractRequest.Validate()
	if err != nil {
		return nil, err
	}
	return h.service.CreateJob(ctx, extractRequest.ValidatedUUID)
}

func (h *ExtractJobCrud) Read(ctx context.Context, uuid uuid.UUID) (*models.ExtractJob, bool, error) {
	job, err := h.service.GetJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return job, true, nil
}

func (h *ExtractJobCrud) List(ctx context.Context) ([]*models.ExtractJob, error) {
	return h.service.ListJobs(ctx)
}

func (h *ExtractJobCrud) Delete(ctx context.Context, uuid uuid.UUID) (bool, error) {
	err := h.service.DeleteJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
