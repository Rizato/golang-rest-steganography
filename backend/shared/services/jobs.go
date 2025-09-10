package services

import (
	"context"
	"steg/shared/models"

	"github.com/google/uuid"
)

type JobService interface {
	GetType() models.JobType
	Exists(context.Context, uuid.UUID) (bool, error)
}
