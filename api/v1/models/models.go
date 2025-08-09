package models

import (
	"github.com/google/uuid"
	"time"
)

type Model interface {
	GetSchema() interface{}
}

type ServerFile struct {
	Uuid     uuid.UUID `json:"uuid"`
	Path     string    `json:"-"`
	Size     int64     `json:"size"`
	Mimetype string    `json:"mime"`
}

func NewServerFile() ServerFile {
	return ServerFile{
		Uuid: uuid.New(),
	}
}

func (i ServerFile) GetSchema() interface{} {
	return nil
}

type Job interface {
	Start()
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
	Uuid              uuid.UUID  `json:"uuid"`
	Status            Status     `json:"status"`
	StatusMessage     string     `json:"status-message"`
	ImageUUID         uuid.UUID  `json:"image-uuid"`
	Message           string     `json:"message"`
	EmbeddedImageUUID *uuid.UUID `json:"embedded-image-uuid"`
	CreatedAt         time.Time  `json:"created-at"`
	UpdatedAt         time.Time  `json:"updated-at"`
}

func NewEmbedJob(image ServerFile) EmbedJob {
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

func (j EmbedJob) GetSchema() interface{} {
	return nil
}

func (j EmbedJob) Start() {
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

func NewExtractJob(image ServerFile) ExtractJob {
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

func (j ExtractJob) GetSchema() interface{} {
	return nil
}

func (j ExtractJob) Start() {
	// Download file
	j.Status = InProgress
	// Check for steg
	// Update state
}
