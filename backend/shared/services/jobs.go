package services

import (
	"context"

	"github.com/google/uuid"
)

type JobService interface {
	GetJob(ctx context.Context, uuid uuid.UUID) (any, error)
	Execute(ctx context.Context, job any) error
}
