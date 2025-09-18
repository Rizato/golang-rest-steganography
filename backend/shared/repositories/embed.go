package repositories

import (
	"context"
	"database/sql"
	"log"
	"steg/shared/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmbedJobRepository struct {
	dbPool *pgxpool.Pool
}

func NewEmbedJobRepository(dbPool *pgxpool.Pool) *EmbedJobRepository {
	return &EmbedJobRepository{dbPool}
}

func (repository *EmbedJobRepository) Create(ctx context.Context, image_uuid uuid.UUID, message string) (*models.EmbedJob, error) {
	var createdId uuid.UUID
	err := repository.dbPool.QueryRow(ctx, "INSERT INTO embed_jobs (image_uuid, message) VALUES ($1, $2) RETURNING id", image_uuid, message).Scan(&createdId)
	if err != nil {
		return nil, err
	}
	return repository.GetByID(ctx, createdId)
}

func (repository *EmbedJobRepository) List(ctx context.Context) ([]*models.EmbedJob, error) {
	rows, err := repository.dbPool.Query(ctx, "SELECT id, status, status_message, image_uuid, message, embedded_uuid, created_at, updated_at FROM embed_jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*models.EmbedJob
	for rows.Next() {
		var job models.EmbedJob
		err = rows.Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUuid, &job.Message, &job.EmbeddedUuid, &job.CreatedAt, &job.UpdatedAt)
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

func (repository *EmbedJobRepository) GetByID(ctx context.Context, uuid uuid.UUID) (*models.EmbedJob, error) {
	var job models.EmbedJob
	err := repository.dbPool.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, embedded_uuid, created_at, updated_at FROM embed_jobs WHERE id = $1", uuid.String()).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUuid, &job.Message, &job.EmbeddedUuid, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		log.Println("Error getting embed_job:", err)
		return nil, err
	}

	return &job, nil
}

func (repository *EmbedJobRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return repository.dbPool.QueryRow(ctx, "DELETE FROM embed_jobs WHERE id = $1", uuid.String()).Scan()
}

func (repository *EmbedJobRepository) Complete(ctx context.Context, job *models.EmbedJob, embedded *models.Image) error {
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
	jobForUpdate.EmbeddedUuid = &embedded.Uuid
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

func (repository *EmbedJobRepository) MarkFailed(ctx context.Context, job *models.EmbedJob, givenError error) error {
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

func (repository *EmbedJobRepository) CheckAndMarkInProgress(ctx context.Context, job *models.EmbedJob) (*models.EmbedJob, error) {
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction", err)
		return job, err
	}
	defer tx.Rollback(ctx)
	jobForUpdate, err := repository.getForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		log.Println("Error getting job for update", err)
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
		log.Println("Error updating job", err)
		return job, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error commiting job", err)
		return job, err
	}
	return jobForUpdate, nil
}

func (repository *EmbedJobRepository) getForUpdate(ctx context.Context, tx pgx.Tx, uuid uuid.UUID) (*models.EmbedJob, error) {
	var job models.EmbedJob
	err := tx.QueryRow(ctx, "SELECT id, status, status_message, image_uuid, message, created_at, updated_at FROM embed_jobs WHERE id = $1 FOR UPDATE", uuid).Scan(&job.Uuid, &job.Status, &job.StatusMessage, &job.ImageUuid, &job.Message, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (repository *EmbedJobRepository) save(ctx context.Context, tx pgx.Tx, job *models.EmbedJob) error {
	log.Println("Saving embed_job", job)
	result, err := tx.Exec(ctx, "UPDATE embed_jobs SET status = $1, status_message = $2, embedded_uuid = $3, updated_at = $4 WHERE id = $5", job.Status, job.StatusMessage, job.EmbeddedUuid, job.UpdatedAt, job.Uuid)
	if result.RowsAffected() != 1 {
		return sql.ErrNoRows
	}
	return err
}
