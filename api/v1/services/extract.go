package services

import (
	"context"
	"os"
	"steg/api/v1/models"
	"steg/api/v1/repository"

	"github.com/google/uuid"
)

type ExtractJobService struct {
	JobRepository   *repository.ExtractJobRepository
	ImageRepository *repository.ImageRepository
}

func NewExtractJobService(jobRepository *repository.ExtractJobRepository, imageRepository *repository.ImageRepository) *ExtractJobService {
	return &ExtractJobService{
		JobRepository:   jobRepository,
		ImageRepository: imageRepository,
	}
}

func (service *ExtractJobService) CreateJob(ctx context.Context, image_uuid uuid.UUID) (*models.ExtractJob, error) {
	// Get image to check
	_, err := service.ImageRepository.GetByID(ctx, image_uuid)
	if err != nil {
		return nil, err
	}
	return service.JobRepository.Create(ctx, image_uuid)
}

func (service ExtractJobService) ListJobs(ctx context.Context) ([]*models.ExtractJob, error) {
	return service.JobRepository.List(ctx)
}

func (service *ExtractJobService) GetJob(ctx context.Context, uuid uuid.UUID) (*models.ExtractJob, error) {
	return service.JobRepository.GetByID(ctx, uuid)
}

func (service *ExtractJobService) DeleteJob(ctx context.Context, uuid uuid.UUID) error {
	return service.JobRepository.Delete(ctx, uuid)
}

func (service *ExtractJobService) GetImage(ctx context.Context, job *models.ExtractJob) (*models.Image, error) {
	return service.ImageRepository.GetByID(ctx, job.ImageUUID)
}

func (service *ExtractJobService) Start(ctx context.Context, job *models.ExtractJob) (*models.ExtractJob, error) {
	image, err := service.GetImage(ctx, job)
	if err != nil {
		return job, err
	}
	// Get's a lock on the row to check if in progress
	job, err = service.CheckAndMarkInProgress(ctx, job)
	if err != nil {
		return job, err
	}
	go func() {
		ctx := context.Background()
		err := service.StartExtractMessage(ctx, job, image)
		if err != nil {
			// TODO What to do with this? Don't want orphaned in progress images
			service.MarkFailed(ctx, job, err)
		}
	}()
	return job, nil
}

func (service *ExtractJobService) StartExtractMessage(ctx context.Context, job *models.ExtractJob, image *models.Image) error {
	message, err := service.ExtractMessage(image)
	if err != nil {
		return err
	}
	tx, err := service.JobRepository.StartTransaction(ctx)
	if err != nil {
		return err
	}
	jobForUpdate, err := service.JobRepository.GetForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	jobForUpdate.Message = message
	jobForUpdate.Status = models.Complete
	jobForUpdate.StatusMessage = "Completed"
	err = service.JobRepository.Save(ctx, tx, jobForUpdate)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (service *ExtractJobService) ExtractMessage(image *models.Image) (string, error) {
	// Does the steg on the image
	file, err := os.Open(image.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// TODO Decode the steg
	return "Temporary Extract Message", nil
}

func (service *ExtractJobService) MarkFailed(ctx context.Context, job *models.ExtractJob, givenError error) (*models.ExtractJob, error) {
	// Mark it failed
	tx, err := service.JobRepository.StartTransaction(ctx)
	if err != nil {
		return job, err
	}
	jobForUpdate, err := service.JobRepository.GetForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		return job, err
	}
	jobForUpdate.Status = models.Error
	jobForUpdate.StatusMessage = givenError.Error()
	err = service.JobRepository.Save(ctx, tx, jobForUpdate)
	if err != nil {
		tx.Rollback(ctx)
		return job, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return job, err
	}
	return jobForUpdate, nil
}

func (service *ExtractJobService) CheckAndMarkInProgress(ctx context.Context, job *models.ExtractJob) (*models.ExtractJob, error) {
	tx, err := service.JobRepository.StartTransaction(ctx)
	if err != nil {
		return job, err
	}
	jobForUpdate, err := service.JobRepository.GetForUpdate(ctx, tx, job.Uuid)
	if err != nil {
		return job, err
	}
	// Another job already claimed it
	// TODO Handle if it is completed already, but allow retries with cancelled or error
	if jobForUpdate.Status != models.InProgress {
		tx.Rollback(ctx)
		return job, AlreadyInProgress
	}
	jobForUpdate.Status = models.InProgress
	jobForUpdate.StatusMessage = "In Progress"
	err = service.JobRepository.Save(ctx, tx, jobForUpdate)
	if err != nil {
		tx.Rollback(ctx)
		return job, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return job, err
	}
	return jobForUpdate, nil
}
