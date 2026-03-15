package repositories

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type StorageRepository interface {
	SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, filename string) error
}

type storageRepository struct {
	db *sql.DB
}

func NewStorageRepository(db *sql.DB) StorageRepository {
	return &storageRepository{db: db}
}

func (r *storageRepository) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext

	dst, err := os.Create("static/images/" + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

func (r *storageRepository) DeleteImage(ctx context.Context, filename string) error {
	path := "static/images/" + filename

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("file not found")
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
