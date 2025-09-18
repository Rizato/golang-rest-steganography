package repositories

import (
	"context"
	"steg/shared/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository struct {
	dbPool *pgxpool.Pool
}

func NewImageRepository(dbPool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{dbPool: dbPool}
}

func (receiver *ImageRepository) List(context context.Context) ([]*models.Image, error) {
	rows, err := receiver.dbPool.Query(context, "SELECT id, path, mimetype, size, encoded, created_at, updated_at size FROM images")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []*models.Image
	for rows.Next() {
		var image models.Image
		err = rows.Scan(&image.Uuid, &image.Path, &image.Mimetype, &image.Size, &image.Encoded, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, &image)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (receiver *ImageRepository) GetByID(context context.Context, uuid uuid.UUID) (*models.Image, error) {
	var image models.Image
	err := receiver.dbPool.QueryRow(context, "SELECT id, path, mimetype, size, encoded, created_at, updated_at FROM images WHERE id = $1", uuid.String()).Scan(&image.Uuid, &image.Path, &image.Mimetype, &image.Size, &image.Encoded, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (receiver *ImageRepository) Create(context context.Context, path string, mimetype string, size int64, encoded bool) (*models.Image, error) {
	var id uuid.UUID
	err := receiver.dbPool.QueryRow(context, "INSERT INTO images (path, mimetype, size, encoded) VALUES ($1, $2, $3, $4) RETURNING id", path, mimetype, size, encoded).Scan(&id)
	if err != nil {
		return nil, err
	}
	return receiver.GetByID(context, id)
}

func (receiver *ImageRepository) Delete(context context.Context, uuid uuid.UUID) error {
	return receiver.dbPool.QueryRow(context, "DELETE FROM images WHERE id = $1", uuid.String()).Scan()
}
