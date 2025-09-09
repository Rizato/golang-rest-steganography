package services

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"steg/api/v1/models"
	"steg/api/v1/repository"
	"steg/steganography"

	"github.com/google/uuid"
)

type EmbedJobService struct {
	JobRepository   *repository.EmbedJobRepository
	ImageRepository *repository.ImageRepository
}

func NewEmbedJobService(jobRepository *repository.EmbedJobRepository, imageRepository *repository.ImageRepository) *EmbedJobService {
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
	return service.ImageRepository.GetByID(ctx, job.ImageUuid)
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
			log.Println("Error marking job as failed", err)
		}
	}()
	return updatedJob, nil
}

func (service *EmbedJobService) StartEmbedMessage(ctx context.Context, job *models.EmbedJob, image *models.Image) error {
	embedded, err := service.EmbedMessage(ctx, job.Message, image)
	if err != nil {
		return err
	}
	return service.JobRepository.Complete(ctx, job, embedded)
}

func (service *EmbedJobService) EmbedMessage(ctx context.Context, message string, image *models.Image) (*models.Image, error) {
	// Does the steg on the image
	file, err := os.Open(image.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	originalImage, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	embedImage, err := steganography.EmbedLsb(originalImage, message)
	if err != nil {
		return nil, err
	}
	toEmbed, err := os.CreateTemp("", "*_embed")
	if err != nil {
		return nil, err
	}
	defer toEmbed.Close()
	err = png.Encode(toEmbed, embedImage)
	if err != nil {
		return nil, err
	}

	embed, err := service.ImageRepository.Create(ctx, toEmbed.Name(), image.Mimetype, image.Size, true)
	if err != nil {
		return nil, err
	}
	return embed, nil
}
