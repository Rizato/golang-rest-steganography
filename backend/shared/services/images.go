package services

import (
	"context"
	"os"
	"steg/shared/models"
	"steg/shared/repositories"

	"github.com/google/uuid"
)

type ImageService struct {
	ImageRepository *repositories.ImageRepository
}

func NewImageService(imageRepository *repositories.ImageRepository) *ImageService {
	return &ImageService{ImageRepository: imageRepository}
}

func (service *ImageService) AddImage(context context.Context, path string, mimetype string, size int64, encoded bool) (*models.Image, error) {
	image, err := service.ImageRepository.Create(context, path, mimetype, size, encoded)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (service *ImageService) GetImage(context context.Context, uuid uuid.UUID) (*models.Image, error) {
	image, err := service.ImageRepository.GetByID(context, uuid)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (service *ImageService) DeleteImage(context context.Context, uuid uuid.UUID) error {
	// Get so we can get 404
	image, err := service.GetImage(context, uuid)
	if err != nil {
		return err
	}
	err = os.Remove(image.Path)
	if err != nil {
		return err
	}
	return service.ImageRepository.Delete(context, uuid)
}
