package models

import (
	"github.com/google/uuid"
	"time"
)

type ImageProcessor interface {
	ProcessImage()
	GetUUID() uuid.UUID
}

type Status string

const (
	Error      Status = "Error"
	Submitted  Status = "Submitted"
	InProgress Status = "In Progress"
	Complete   Status = "Complete"
)

type EmbedJob struct {
	Uuid              uuid.UUID  `json:"uuid"`
	Status            Status     `json:"status"`
	StatusMessage     string     `json:"status-message"`
	ImageUUID         uuid.UUID  `json:"image-uuid"`
	EmbeddedImageUUID *uuid.UUID `json:"embedded-image-uuid"`
	CreatedAt         time.Time  `json:"created-at"`
	UpdatedAt         time.Time  `json:"updated-at"`
}

func NewEmbedJob(imageUUID uuid.UUID) *EmbedJob {
	return &EmbedJob{
		Uuid:              uuid.New(),
		Status:            Submitted,
		StatusMessage:     "",
		ImageUUID:         imageUUID,
		EmbeddedImageUUID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func (j *EmbedJob) GetUUID() uuid.UUID {
	return j.Uuid
}

func (j *EmbedJob) ProcessImage() {
	// Download file
	j.Status = InProgress
	// Add steg
	// Update state
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

func NewExtractJob(imageUUID uuid.UUID) *ExtractJob {
	return &ExtractJob{
		Uuid:          uuid.New(),
		Status:        Submitted,
		StatusMessage: "",
		ImageUUID:     imageUUID,
		Message:       "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (j *ExtractJob) GetUUID() uuid.UUID {
	return j.Uuid
}

func (j *ExtractJob) ProcessImage() {
	// Download file
	j.Status = InProgress
	// Check for steg
	// Update state
}
