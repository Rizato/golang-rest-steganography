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
	MissingMessageError   = errors.New("missing message")
	MessageTooLargeError  = errors.New("message too large")
	MissingImageUUIDError = errors.New("missing image uuid")
	ImageUUIDInvalidError = errors.New("image uuid invald")
	ImageNotFoundError    = errors.New("image not found")
)

type EmbedJobCrud struct {
	ds *models.Datastore
}

func NewEmbedJobCrud(ds *models.Datastore) *EmbedJobCrud {
	return &EmbedJobCrud{ds}
}

func (h *EmbedJobCrud) Create(reader io.Reader) (models.EmbedJob, error) {
	var embedRequest CreateEmbedRequest
	var embedJob models.EmbedJob
	err := json.NewDecoder(reader).Decode(&embedRequest)
	if err != nil {
		return embedJob, err
	}

	if embedRequest.message == "" {
		return embedJob, MissingMessageError
	}

	if len(embedRequest.message) > 1000 {
		return embedJob, MessageTooLargeError
	}

	if embedRequest.imageUUID == "" {
		return embedJob, MissingImageUUIDError
	}

	if len(embedRequest.imageUUID) != 36 {
		return embedJob, ImageUUIDInvalidError
	}

	parsedUUID, err := uuid.Parse(embedRequest.imageUUID)
	if err != nil {
		return embedJob, ImageUUIDInvalidError
	}
	image, found := h.ds.GetImage(parsedUUID)
	if !found {
		return embedJob, ImageNotFoundError
	}
	embedJob = models.NewEmbedJob(image.Uuid)
	h.ds.SaveEmbedJob(embedJob)

	return embedJob, nil
}

func (h *EmbedJobCrud) Read(uuid uuid.UUID) (models.EmbedJob, bool, error) {
	job, found := h.ds.GetEmbedJob(uuid)
	return job, found, nil
}

func (h *EmbedJobCrud) List() ([]models.EmbedJob, error) {
	return h.ds.GetEmbedJobs(), nil
}

func (h *EmbedJobCrud) Delete(uuid uuid.UUID) (bool, error) {
	h.ds.DeleteEmbedJob(uuid)
	return true, nil
}

type ExtractJobCrud struct {
	ds *models.Datastore
}

func NewExtractJobCrud(ds *models.Datastore) *ExtractJobCrud {
	return &ExtractJobCrud{ds}
}

func (h *ExtractJobCrud) Create(reader io.Reader) (models.ExtractJob, error) {
	var extractRequest CreateExtractRequest
	var extractJob models.ExtractJob
	err := json.NewDecoder(reader).Decode(&extractRequest)
	if err != nil {
		return extractJob, err
	}

	if extractRequest.imageUUID == "" {
		return extractJob, MissingImageUUIDError
	}

	if len(extractRequest.imageUUID) != 36 {
		return extractJob, ImageUUIDInvalidError
	}

	parsedUUID, err := uuid.Parse(extractRequest.imageUUID)
	if err != nil {
		return extractJob, ImageUUIDInvalidError
	}
	image, found := h.ds.GetImage(parsedUUID)
	if !found {
		return extractJob, ImageNotFoundError
	}

	extractJob = models.NewExtractJob(image.Uuid)
	h.ds.SaveExtractJob(extractJob)

	return extractJob, nil
}

func (h *ExtractJobCrud) Read(uuid uuid.UUID) (models.ExtractJob, bool, error) {
	job, found := h.ds.GetExtractJob(uuid)
	return job, found, nil
}

func (h *ExtractJobCrud) List() ([]models.ExtractJob, error) {
	return h.ds.GetExtractJobs(), nil
}

func (h *ExtractJobCrud) Delete(uuid uuid.UUID) (bool, error) {
	h.ds.DeleteExtractJob(uuid)
	return true, nil
}

// FileCrud implements RD of crud, because we have separate handlers for file content
type FileCrud struct {
	ds *models.Datastore
}

func NewFileCrud(ds *models.Datastore) *FileCrud {
	return &FileCrud{ds}
}

func (h *FileCrud) Read(uuid uuid.UUID) (models.ServerFile, bool, error) {
	job, found := h.ds.GetImage(uuid)
	return job, found, nil
}

func (h *FileCrud) Delete(uuid uuid.UUID) (bool, error) {
	image, found := h.ds.GetImage(uuid)
	if !found {
		return false, nil
	}
	// Delete from datastore, and filesystem
	err := os.Remove(image.Path)
	if err != nil {
		return false, err
	}
	h.ds.DeleteImage(image.Uuid)
	return true, nil
}
