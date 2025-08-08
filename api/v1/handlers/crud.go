package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"os"
	"steg/api/v1/models"
)

var (
	MethodNotAllowedError = errors.New("method not allowed")
	MissingMessageError   = errors.New("missing message")
	MessageTooLargeError  = errors.New("message too large")
	MissingImageUUIDError = errors.New("missing image uuid")
	ImageUUIDInvalidError = errors.New("image uuid invald")
	ImageNotFoundError    = errors.New("image not found")
)

type EmbedJobCrud struct {
	*models.Datastore
}

func NewEmbedJobCrud(ds *models.Datastore) *EmbedJobCrud {
	return &EmbedJobCrud{ds}
}

func (h *EmbedJobCrud) Create(reader io.Reader) (*models.EmbedJob, error) {
	var embedRequest CreateEmbedRequest
	err := json.NewDecoder(reader).Decode(&embedRequest)
	if err != nil {
		return nil, err
	}

	if embedRequest.message == "" {
		return nil, MissingMessageError
	}

	if len(embedRequest.message) > 1000 {
		return nil, MessageTooLargeError
	}

	if embedRequest.imageUUID == "" {
		return nil, MissingImageUUIDError
	}

	if len(embedRequest.imageUUID) != 36 {
		return nil, ImageUUIDInvalidError
	}

	parsedUUID, err := uuid.Parse(embedRequest.imageUUID)
	if err != nil {
		return nil, ImageUUIDInvalidError
	}
	image, found := h.Images[parsedUUID]
	if !found {
		return nil, ImageNotFoundError
	}

	return models.NewEmbedJob(image.GetUUID()), nil
}

func (h *EmbedJobCrud) Read(uuid uuid.UUID) (*models.EmbedJob, bool, error) {
	job, found := h.EmbedJobs[uuid]
	return job, found, nil
}

func (h *EmbedJobCrud) Update(uuid uuid.UUID) (*models.EmbedJob, error) {
	return nil, MethodNotAllowedError
}

func (h *EmbedJobCrud) Delete(uuid uuid.UUID) (bool, error) {
	job, found := h.EmbedJobs[uuid]
	if !found {
		return false, nil
	}
	delete(h.EmbedJobs, job.Uuid)
	return true, nil
}

type ExtractJobCrud struct {
	*models.Datastore
}

func NewExtractJobCrud(ds *models.Datastore) *ExtractJobCrud {
	return &ExtractJobCrud{ds}
}

func (h *ExtractJobCrud) Create(reader io.Reader) (*models.ExtractJob, error) {
	var extractRequest CreateExtractRequest
	err := json.NewDecoder(reader).Decode(&extractRequest)
	if err != nil {
		return nil, err
	}

	if extractRequest.imageUUID == "" {
		return nil, MissingImageUUIDError
	}

	if len(extractRequest.imageUUID) != 36 {
		return nil, ImageUUIDInvalidError
	}

	parsedUUID, err := uuid.Parse(extractRequest.imageUUID)
	if err != nil {
		return nil, ImageUUIDInvalidError
	}
	image, found := h.Images[parsedUUID]
	if !found {
		return nil, ImageNotFoundError
	}

	return models.NewExtractJob(image.GetUUID()), nil
}

func (h *ExtractJobCrud) Read(uuid uuid.UUID) (*models.ExtractJob, bool, error) {
	job, found := h.ExtractJobs[uuid]
	return job, found, nil
}

func (h *ExtractJobCrud) Update(uuid uuid.UUID) (*models.ExtractJob, error) {
	return nil, MethodNotAllowedError
}

func (h *ExtractJobCrud) Delete(uuid uuid.UUID) (bool, error) {
	job, found := h.ExtractJobs[uuid]
	if !found {
		return false, nil
	}
	delete(h.EmbedJobs, job.Uuid)
	return true, nil
}

// FileCrud implements RD of crud, because we have separate handlers for file content
type FileCrud struct {
	*models.Datastore
}

func NewFileCrud(ds *models.Datastore) *FileCrud {
	return &FileCrud{ds}
}

func (h *FileCrud) Create(reader io.Reader) (*models.ServerFile, error) {
	return nil, MethodNotAllowedError
}

func (h *FileCrud) Read(uuid uuid.UUID) (*models.ServerFile, bool, error) {
	job, found := h.Images[uuid]
	return job, found, nil
}

func (h *FileCrud) Update(uuid uuid.UUID) (*models.ServerFile, error) {
	return nil, MethodNotAllowedError
}

func (h *FileCrud) Delete(uuid uuid.UUID) (bool, error) {
	image, found := h.Images[uuid]
	if !found {
		return false, nil
	}
	// Delete from datastore, and filesystem
	err := os.Remove(image.Path)
	if err != nil {
		return false, err
	}
	delete(h.EmbedJobs, image.Uuid)
	return true, nil
}
