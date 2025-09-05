package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	Uuid      uuid.UUID `json:"uuid"`
	Path      string    `json:"-"`
	Size      int64     `json:"size"`
	Mimetype  string    `json:"mime"`
	Encoded   bool      `json:"encoded"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

func NewImage() Image {
	return Image{
		Uuid: uuid.New(),
	}
}

type Status string

const (
	Error      Status = "Error"
	Submitted  Status = "Submitted"
	InProgress Status = "In Progress"
	Complete   Status = "Complete"
	Cancelled  Status = "Cancelled"
)

type EmbedJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"`
	StatusMessage string    `json:"status-message"`
	// TODO Separate models per layer, for nested objects
	ImageUUID         uuid.UUID  `json:"image-uuid"`
	Message           string     `json:"message"`
	EmbeddedImageUUID *uuid.UUID `json:"embedded-image-uuid"`
	CreatedAt         time.Time  `json:"created-at"`
	UpdatedAt         time.Time  `json:"updated-at"`
}

func NewEmbedJob(image Image) EmbedJob {
	return EmbedJob{
		Uuid:              uuid.New(),
		Status:            Submitted,
		StatusMessage:     "",
		ImageUUID:         image.Uuid,
		EmbeddedImageUUID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

type ExtractJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"`
	StatusMessage string    `json:"status-message"`
	ImageUUID     uuid.UUID `json:"image-uuid"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"created-at"`
	UpdatedAt     time.Time `json:"updated-at"`
}

func NewExtractJob(image Image) ExtractJob {
	return ExtractJob{
		Uuid:          uuid.New(),
		Status:        Submitted,
		StatusMessage: "",
		ImageUUID:     image.Uuid,
		Message:       "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
