package repository

import (
	"context"
	"steg/api/v1/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExtractJobRepository struct {
	dbPool *pgxpool.Pool
}

func NewExtractJobRepository(dbPool *pgxpool.Pool) *ExtractJobRepository {
	return &ExtractJobRepository{dbPool}
}

func (receiver *ExtractJobRepository) Create(ctx context.Context, image_uuid uuid.UUID) (*models.ExtractJob, error) {
	var createdId uuid.UUID
	err := receiver.dbPool.QueryRow(ctx, "INSERT INTO extract_jobs (image_uuid) VALUES ($1) RETURNING id", image_uuid).Scan(&createdId)
	if err != nil {
		return nil, err
	}
	return receiver.GetByID(ctx, createdId)
}

func (receiver *ExtractJobRepository) List(ctx context.Context) ([]*models.ExtractJob, error) {
	rows, err := receiver.dbPool.Query(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*models.ExtractJob
	for rows.Next() {
		var job models.ExtractJob
		err = rows.Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUUID, &job.Message, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (receiver *ExtractJobRepository) GetByID(ctx context.Context, uuid uuid.UUID) (*models.ExtractJob, error) {
	var job *models.ExtractJob
	err := receiver.dbPool.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs WHERE id = $1", uuid.String()).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUUID, &job.Message, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (receiver ExtractJobRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return receiver.dbPool.QueryRow(ctx, "DELETE FROM extract_jobs WHERE id = $1", uuid.String()).Scan()
}

func (receiver ExtractJobRepository) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := receiver.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (receiver *ExtractJobRepository) GetForUpdate(ctx context.Context, tx pgx.Tx, uuid uuid.UUID) (*models.ExtractJob, error) {
	var job *models.ExtractJob
	err := tx.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs WHERE id = $1 FOR UPDATE", uuid.String()).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUUID, &job.Message, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (receiver ExtractJobRepository) Save(ctx context.Context, tx pgx.Tx, job *models.ExtractJob) error {
	job.UpdatedAt = time.Now()
	return tx.QueryRow(ctx, "UPDATE extract_jobs SET status=$1 status_message=$2 message=$3 updated_at=$4 WHERE id=$4", job.Status, job.StatusMessage, job.Message, job.UpdatedAt).Scan()
}
