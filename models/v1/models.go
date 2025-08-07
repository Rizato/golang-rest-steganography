package v1

import (
	"github.com/google/uuid"
	"time"
)

type ImageProcessor interface {
	ProcessImage()
}

type Status int

const (
	Error Status = iota
	Submitted
	InProgress
	Complete
)

func (s Status) String() string {
	switch s {
	case Error:
		return "Error"
	case Submitted:
		return "Submitted"
	case InProgress:
		return "In Progress"
	case Complete:
		return "Complete"
	default:
		return "Unknown"
	}
}

type EncodeJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"` // Need to add custom encoding
	StatusMessage string    `json:"status-message"`
	RawImagePath  string    `json:"-"`
	StegImagePath string    `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEncodeJob() *EncodeJob {
	return &EncodeJob{
		Uuid:          uuid.New(),
		Status:        Submitted,
		StatusMessage: "",
		RawImagePath:  "",
		StegImagePath: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (e *EncodeJob) ProcessImage() {
	// Download file
	e.Status = InProgress
	// Add steg
	// Update state
}

type DecodeJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"`
	StatusMessage string    `json:"status-message"`
	ImagePath     string    `json:"-"` // Do not include
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewDecodeJob() *DecodeJob {
	return &DecodeJob{
		Uuid:          uuid.New(),
		Status:        InProgress,
		StatusMessage: "",
		ImagePath:     "",
		Message:       "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (d *DecodeJob) ProcessImage() {
	// Download file
	d.Status = InProgress
	// Check for steg
	// Update state
}
