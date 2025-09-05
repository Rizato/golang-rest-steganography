package services

import (
	"context"
	"errors"
)

var (
	AlreadyInProgress = errors.New("already in progress")
)

type JobService[T any] interface {
	Start(ctx context.Context, job *T) (*T, error)
}
