package services

import (
	"context"
	"gowes/models"
	"gowes/repositories"
	"log"
	"mime/multipart"
)

type ProductService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error)
	Create(companyID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader, addOnIDList []string) (models.Product, error)
	FindByID(productID string) (models.Product, error)
	DeleteById(productID string) error
	Update(productID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader) (models.Product, error)
}

type productService struct {
	productRepository repositories.ProductRepository
	storageRepository repositories.StorageRepository
}

func NewProductService(productRepository repositories.ProductRepository, storageRepository repositories.StorageRepository) ProductService {
	return &productService{productRepository: productRepository, storageRepository: storageRepository}
}

func (s *productService) FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error) {
	products, total, err := s.productRepository.FindAll(companyID, params)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *productService) Create(companyID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader, addOnIDList []string) (models.Product, error) {
	imageURL, err := s.storageRepository.SaveImage(context.Background(), imageFile, imageHeader)
	if err != nil {
		return models.Product{}, err
	}

	payload.ImageURL = imageURL

	product, err := s.productRepository.Create(companyID, payload)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (s *productService) FindByID(productID string) (models.Product, error) {
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (s *productService) DeleteById(productID string) error {
	// DeleteById mengembalikan image_url via RETURNING — tidak perlu FindByID terpisah
	imageURL, err := s.productRepository.DeleteById(productID)
	if err != nil {
		return err
	}

	if imageURL != "" {
		// Best-effort: log error tapi jangan gagalkan response karena produk sudah terhapus
		if err := s.storageRepository.DeleteImage(context.Background(), imageURL); err != nil {
			log.Printf("warning: failed to delete image %q from storage: %v", imageURL, err)
		}
	}

	return nil
}

func (s *productService) Update(productID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader) (models.Product, error) {
	// Ambil produk lama untuk mendapatkan image_url yang ada
	existing, err := s.productRepository.FindByID(productID)
	if err != nil {
		return models.Product{}, err
	}

	// Gunakan image_url lama sebagai default
	payload.ImageURL = existing.ImageURL

	// Jika ada gambar baru dikirim, upload terlebih dahulu
	if imageFile != nil && imageHeader != nil {
		newImageURL, err := s.storageRepository.SaveImage(context.Background(), imageFile, imageHeader)
		if err != nil {
			return models.Product{}, err
		}

		// Hapus gambar lama dari storage (best-effort)
		if existing.ImageURL != "" {
			if err := s.storageRepository.DeleteImage(context.Background(), existing.ImageURL); err != nil {
				log.Printf("warning: failed to delete old image %q from storage: %v", existing.ImageURL, err)
			}
		}

		payload.ImageURL = newImageURL
	}

	updated, err := s.productRepository.Update(productID, payload)
	if err != nil {
		return models.Product{}, err
	}

	return updated, nil
}
