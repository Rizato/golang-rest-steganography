package services

import (
	"context"
)

type JobService[T any] interface {
	Start(ctx context.Context, job *T) (*T, error)
}
