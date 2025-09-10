package services

import (
	"context"
	"database/sql"
	"errors"
	"image/png"
	"io"
	"log"
	"os"
	"steg/shared/models"
	"steg/shared/repositories"
	"steg/steganography"

	"github.com/google/uuid"
)

var NotEmbedJob = errors.New("not embed job")

type EmbedJobService struct {
	JobRepository   *repositories.EmbedJobRepository
	ImageRepository *repositories.ImageRepository
}

func NewEmbedJobService(jobRepository *repositories.EmbedJobRepository, imageRepository *repositories.ImageRepository) *EmbedJobService {
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

func (service *EmbedJobService) GetEmbedJob(ctx context.Context, uuid uuid.UUID) (*models.EmbedJob, error) {
	return service.JobRepository.GetByID(ctx, uuid)
}

func (service *EmbedJobService) DeleteJob(ctx context.Context, uuid uuid.UUID) error {
	return service.JobRepository.Delete(ctx, uuid)
}

func (service *EmbedJobService) GetImage(ctx context.Context, job *models.EmbedJob) (*models.Image, error) {
	return service.ImageRepository.GetByID(ctx, job.ImageUuid)
}

func (service *EmbedJobService) Embed(ctx context.Context, job *models.EmbedJob) error {
	image, err := service.GetImage(ctx, job)
	if err != nil {
		return err
	}
	// Get's a lock on the row to check if in progress
	updatedJob, err := service.JobRepository.CheckAndMarkInProgress(ctx, job)
	if err != nil {
		return err
	}
	go func() {
		ctx := context.Background()
		err := service.StartEmbedMessage(ctx, updatedJob, image)
		if err != nil {
			err = service.JobRepository.MarkFailed(ctx, updatedJob, err)
			if err != nil {
				log.Println("Error marking job as failed", err)
			}
		}
	}()
	return nil
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
	// Use no compression and no filtering to preserve LSBs
	encoder := &png.Encoder{
		CompressionLevel: png.NoCompression,
		BufferPool:       nil,
	}
	err = encoder.Encode(toEmbed, embedImage)
	if err != nil {
		return nil, err
	}
	_, err = toEmbed.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	// Get the actual file size
	fileInfo, err := toEmbed.Stat()
	if err != nil {
		return nil, err
	}

	embed, err := service.ImageRepository.Create(ctx, toEmbed.Name(), image.Mimetype, fileInfo.Size(), true)
	if err != nil {
		return nil, err
	}
	return embed, nil
}

// JobService interface implementation

func (service *EmbedJobService) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
	job, err := service.GetEmbedJob(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return job != nil, nil
}

func (service *EmbedJobService) GetType() models.JobType {
	return models.Extract
}
