package services

import (
	"context"
	"fmt"
	"os"
	"steg/api/v1/models"
	repository2 "steg/api/v1/repository"

	"github.com/google/uuid"
)

type EmbedJobService struct {
	JobRepository   *repository2.EmbedJobRepository
	ImageRepository *repository2.ImageRepository
}

func NewEmbedJobService(jobRepository *repository2.EmbedJobRepository, imageRepository *repository2.ImageRepository) *EmbedJobService {
	return &EmbedJobService{
		JobRepository:   jobRepository,
		ImageRepository: imageRepository,
	}
}

func (service *EmbedJobService) CreateJob(ctx context.Context, image_uuid uuid.UUID, message string) (*models.EmbedJob, error) {
	// Get image to check
	_, err := service.ImageRepository.GetByID(ctx, image_uuid)
	if err != nil {
		return nil, err
	}
	return service.JobRepository.Create(ctx, image_uuid, message)
}

func (service *EmbedJobService) ListJobs(ctx context.Context) ([]*models.EmbedJob, error) {
	return service.JobRepository.List(ctx)
}

func (service *EmbedJobService) GetJob(ctx context.Context, uuid uuid.UUID) (*models.EmbedJob, error) {
	return service.JobRepository.GetByID(ctx, uuid)
}

func (service *EmbedJobService) DeleteJob(ctx context.Context, uuid uuid.UUID) error {
	return service.JobRepository.Delete(ctx, uuid)
}

func (service *EmbedJobService) GetImage(ctx context.Context, job *models.EmbedJob) (*models.Image, error) {
	return service.ImageRepository.GetByID(ctx, job.ImageUUID)
}

func (service *EmbedJobService) Start(ctx context.Context, job *models.EmbedJob) (*models.EmbedJob, error) {
	image, err := service.GetImage(ctx, job)
	if err != nil {
		return job, err
	}
	// Get's a lock on the row to check if in progress
	updatedJob, err := service.JobRepository.CheckAndMarkInProgress(ctx, job)
	if err != nil {
		return updatedJob, err
	}
	go func() {
		ctx := context.Background()
		err := service.StartEmbedMessage(ctx, job, image)
		if err != nil {
			err = service.JobRepository.MarkFailed(ctx, job, err)
			fmt.Println("Error marking job as failed", err)
		}
	}()
	return updatedJob, nil
}

func (service *EmbedJobService) StartEmbedMessage(ctx context.Context, job *models.EmbedJob, image *models.Image) error {
	embedded, err := service.EmbedMessage(job.Message, image)
	if err != nil {
		return err
	}
	return service.JobRepository.Complete(ctx, job, embedded)
}

func (service *EmbedJobService) EmbedMessage(message string, image *models.Image) (*models.Image, error) {
	// Does the steg on the image
	file, err := os.Open(image.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO Embded the message
	return nil, nil
}
