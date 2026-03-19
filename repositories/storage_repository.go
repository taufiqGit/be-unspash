package repositories

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type StorageRepository interface {
	SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, fileURL string) error
}

type storageRepository struct {
	client    *s3.Client
	bucket    string
	publicURL string // base URL untuk akses publik object
}

// NewStorageRepository membuat StorageRepository baru menggunakan AWS S3 client.
//
// Env vars yang dibaca:
//   - S3_BUCKET     : nama bucket
//   - S3_PUBLIC_URL : base URL publik, contoh: https://bucket.nos.wjv-1.neo.id
//                     Jika kosong, fallback ke https://<S3_ENDPOINT>/<bucket>
func NewStorageRepository(client *s3.Client) StorageRepository {
	bucket := os.Getenv("S3_BUCKET")

	// Tentukan public URL base
	publicURL := os.Getenv("S3_PUBLIC_URL")
	if publicURL == "" {
		endpoint := os.Getenv("S3_ENDPOINT")
		// Fallback: path-style URL → https://<endpoint>/<bucket>
		publicURL = strings.TrimRight(endpoint, "/") + "/" + bucket
	}
	publicURL = strings.TrimRight(publicURL, "/")

	return &storageRepository{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}
}

// SaveImage mengupload file gambar ke S3 bucket dan mengembalikan public URL-nya.
func (r *storageRepository) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	key := "images/" + uuid.New().String() + ext

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = resolveContentType(ext)
	}

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Return full public URL: https://<publicURL>/images/<uuid>.<ext>
	return r.publicURL + "/" + key, nil
}

// DeleteImage menghapus object dari S3 berdasarkan public URL-nya.
func (r *storageRepository) DeleteImage(ctx context.Context, fileURL string) error {
	// Ekstrak S3 key dari URL: hapus prefix publicURL
	prefix := r.publicURL + "/"
	key := strings.TrimPrefix(fileURL, prefix)
	if key == fileURL {
		return fmt.Errorf("invalid file URL, does not match public URL prefix: %s", fileURL)
	}

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %w", err)
	}

	return nil
}

// resolveContentType mengembalikan MIME type berdasarkan ekstensi file.
func resolveContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
