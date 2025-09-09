package services

import (
	"context"
	"image/png"
	"log"
	"os"
	"steg/shared/models"
	"steg/shared/repositories"
	"steg/steganography"

	"github.com/google/uuid"
)

type ExtractJobService struct {
	JobRepository   *repositories.ExtractJobRepository
	ImageRepository *repositories.ImageRepository
}

func NewExtractJobService(jobRepository *repositories.ExtractJobRepository, imageRepository *repositories.ImageRepository) *ExtractJobService {
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
			if err != nil {
				log.Println("Error marking job as failed", err)
			}
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
	toExtract, err := png.Decode(file)
	if err != nil {
		return "", err
	}
	return steganography.ExtractLsb(toExtract)
}
