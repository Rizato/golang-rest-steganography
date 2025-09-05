package services

import (
	"context"

	"github.com/google/uuid"
)

type JobService[T any] interface {
	GetJob(ctx context.Context, uuid uuid.UUID) (*T, error)
	Start(ctx context.Context, job *T) (*T, error)
}
