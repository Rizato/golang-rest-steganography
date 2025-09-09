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
	// TODO Separate requests per layer, for nested objects
	ImageUuid    uuid.UUID  `json:"image-uuid"`
	Message      string     `json:"message"`
	EmbeddedUuid *uuid.UUID `json:"embedded-uuid"`
	CreatedAt    time.Time  `json:"created-at"`
	UpdatedAt    time.Time  `json:"updated-at"`
}

type ExtractJob struct {
	Uuid          uuid.UUID `json:"uuid"`
	Status        Status    `json:"status"`
	StatusMessage string    `json:"status-message"`
	ImageUUID     uuid.UUID `json:"image-uuid"`
	Message       *string   `json:"message"`
	CreatedAt     time.Time `json:"created-at"`
	UpdatedAt     time.Time `json:"updated-at"`
}
