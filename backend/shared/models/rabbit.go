package models

import "github.com/google/uuid"

/**
Models for passing messages to rabbitmq
*/

type JobType string

const (
	Embed   JobType = "EMBED"
	Extract JobType = "EXTRACT"
)

// JobMessage is the message passed to rabbit so it knows what job to start
type JobMessage struct {
	Job  JobType
	Uuid uuid.UUID
}

func NewJobMessage(job JobType, uuid uuid.UUID) JobMessage {
	return JobMessage{job, uuid}
}
