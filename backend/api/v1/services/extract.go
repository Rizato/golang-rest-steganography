package services

import (
	"context"
	"fmt"
	"os"
	"steg/api/v1/models"
	repository2 "steg/api/v1/repository"

	"github.com/google/uuid"
)

type ExtractJobService struct {
	JobRepository   *repository2.ExtractJobRepository
	ImageRepository *repository2.ImageRepository
}

func NewExtractJobService(jobRepository *repository2.ExtractJobRepository, imageRepository *repository2.ImageRepository) *ExtractJobService {
	return &ExtractJobService{
		JobRepository:   jobRepository,
		ImageRepository: imageRepository,
	}
}

func (service *ExtractJobService) CreateJob(ctx context.Context, imageUuid uuid.UUID) (*models.ExtractJob, error) {
	// Get image to check
	_, err := service.ImageRepository.GetByID(ctx, imageUuid)
	if err != nil {
		return nil, err
	}
	return service.JobRepository.Create(ctx, imageUuid)
}

func (service *ExtractJobService) ListJobs(ctx context.Context) ([]*models.ExtractJob, error) {
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
	// Gets a lock on the row to check if in progress
	updatedJob, err := service.JobRepository.CheckAndMarkInProgress(ctx, job)
	if err != nil {
		return updatedJob, err
	}
	go func() {
		ctx := context.Background()
		err := service.StartExtractMessage(ctx, job, image)
		if err != nil {
			err = service.JobRepository.MarkFailed(ctx, job, err)
			fmt.Println("Error marking job as failed", err)
		}
	}()
	return updatedJob, nil
}

func (service *ExtractJobService) StartExtractMessage(ctx context.Context, job *models.ExtractJob, image *models.Image) error {
	message, err := service.ExtractMessage(image)
	if err != nil {
		return err
	}
	return service.JobRepository.Complete(ctx, job, message)
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
