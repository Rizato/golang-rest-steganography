package repository

import (
	"context"
	"database/sql"
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

func (repository *ExtractJobRepository) Create(ctx context.Context, image_uuid uuid.UUID) (*models.ExtractJob, error) {
	var createdId uuid.UUID
	err := repository.dbPool.QueryRow(ctx, "INSERT INTO extract_jobs (image_uuid) VALUES ($1) RETURNING id", image_uuid).Scan(&createdId)
	if err != nil {
		return nil, err
	}
	return repository.GetByID(ctx, createdId)
}

func (repository *ExtractJobRepository) List(ctx context.Context) ([]*models.ExtractJob, error) {
	rows, err := repository.dbPool.Query(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs")
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

func (repository *ExtractJobRepository) GetByID(ctx context.Context, uuid uuid.UUID) (*models.ExtractJob, error) {
	var job models.ExtractJob
	err := repository.dbPool.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs WHERE id = $1", uuid.String()).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUUID, &job.Message, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (repository *ExtractJobRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return repository.dbPool.QueryRow(ctx, "DELETE FROM extract_jobs WHERE id = $1", uuid.String()).Scan()
}

func (repository *ExtractJobRepository) Complete(ctx context.Context, job *models.ExtractJob, message string) error {
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	jobForUpdate, err := repository.getForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		return err
	}
	jobForUpdate.Status = models.Complete
	jobForUpdate.StatusMessage = "Completed"
	jobForUpdate.Message = &message
	jobForUpdate.UpdatedAt = time.Now()
	err = repository.save(ctx, tx, jobForUpdate)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ExtractJobRepository) MarkFailed(ctx context.Context, job *models.ExtractJob, givenError error) error {
	// Mark it failed
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	jobForUpdate, err := repository.getForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		return err
	}
	jobForUpdate.Status = models.Error
	jobForUpdate.StatusMessage = givenError.Error()
	jobForUpdate.UpdatedAt = time.Now()
	err = repository.save(ctx, tx, jobForUpdate)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ExtractJobRepository) CheckAndMarkInProgress(ctx context.Context, job *models.ExtractJob) (*models.ExtractJob, error) {
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		return job, err
	}
	defer tx.Rollback(ctx)
	jobForUpdate, err := repository.getForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		return job, err
	}
	// Another job already claimed it
	// TODO Handle if it is completed already, but allow retries with cancelled or error
	if jobForUpdate.Status == models.InProgress {
		return job, AlreadyInProgress
	}
	jobForUpdate.Status = models.InProgress
	jobForUpdate.StatusMessage = "In Progress"
	jobForUpdate.UpdatedAt = time.Now()
	err = repository.save(ctx, tx, jobForUpdate)
	if err != nil {
		return job, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return job, err
	}
	return jobForUpdate, nil
}

func (repository *ExtractJobRepository) getForUpdate(ctx context.Context, tx pgx.Tx, uuid uuid.UUID) (*models.ExtractJob, error) {
	var job models.ExtractJob
	err := tx.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM extract_jobs WHERE id = $1 FOR UPDATE", uuid).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUUID, &job.Message, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (repository *ExtractJobRepository) save(ctx context.Context, tx pgx.Tx, job *models.ExtractJob) error {
	result, err := tx.Exec(ctx, "UPDATE extract_jobs SET status = $1, status_message = $2, message = $3, updated_at = $4 WHERE id = $5", job.Status, job.StatusMessage, job.Message, job.UpdatedAt, job.Uuid)
	if result.RowsAffected() != 1 {
		return sql.ErrNoRows
	}
	return err
}
