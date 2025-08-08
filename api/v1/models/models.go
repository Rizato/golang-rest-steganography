package models

import (
	"github.com/google/uuid"
	"time"
)

type ImageProcessor interface {
	ProcessImage()
	GetUUID() uuid.UUID
	SetImagePath(path string)
}

type Status string

const (
	Error      Status = "Error"
	Submitted  Status = "Submitted"
	InProgress Status = "In Progress"
	Complete   Status = "Complete"
)

type EmbedJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"`
	StatusMessage string    `json:"status-message"`
	ImagePath     string    `json:"-"`
	StegImagePath string    `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEmbedJob() *EmbedJob {
	return &EmbedJob{
		Uuid:          uuid.New(),
		Status:        Submitted,
		StatusMessage: "",
		ImagePath:     "",
		StegImagePath: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (j *EmbedJob) GetUUID() uuid.UUID {
	return j.Uuid
}

func (j *EmbedJob) SetImagePath(path string) {
	j.ImagePath = path
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
	ImagePath     string    `json:"-"` // Do not include
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewExtractJob() *ExtractJob {
	return &ExtractJob{
		Uuid:          uuid.New(),
		Status:        InProgress,
		StatusMessage: "",
		ImagePath:     "",
		Message:       "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (j *ExtractJob) GetUUID() uuid.UUID {
	return j.Uuid
}

func (j *ExtractJob) SetImagePath(path string) {
	j.ImagePath = path
}

func (j *ExtractJob) ProcessImage() {
	// Download file
	j.Status = InProgress
	// Check for steg
	// Update state
}
