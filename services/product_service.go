package services

import (
	"context"
	"gowes/models"
	"gowes/repositories"
	"mime/multipart"
)

type ProductService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error)
	Create(companyID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader) (models.Product, error)
	FindByID(productID string) (models.Product, error)
	DeleteById(productID string) error
	// Update(productID string, payload models.ProductInput) (models.Product, error)
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

func (s *productService) Create(companyID string, payload models.ProductInput, imageFile multipart.File, imageHeader *multipart.FileHeader) (models.Product, error) {
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
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}

	if product.ImageURL != "" {
		err := s.storageRepository.DeleteImage(context.Background(), product.ImageURL)
		if err != nil {
			return err
		}
	}

	err = s.productRepository.DeleteById(productID)
	if err != nil {
		return err
	}
	return nil
}

// func (s *productService) Update(productID string, payload models.ProductInput) (models.Product, error) {
// 	product, err := s.productRepository.Update(productID, payload)
// 	if err != nil {
// 		return models.Product{}, err
// 	}
// 	return product, nil
// }
