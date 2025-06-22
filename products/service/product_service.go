package service

import (
	"time"

	"github.com/Aakash-Sleur/go-micro-product/dto"
	"github.com/Aakash-Sleur/go-micro-product/models"
	"github.com/Aakash-Sleur/go-micro-product/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product dto.ProductDto, userId string) (*models.Product, error) {
	uid, err := uuid.Parse(userId)
	if err != nil {
		// handle error appropriately, e.g., return or log
		return nil, err
	}
	newProduct := &models.Product{
		ID:            uuid.New(),
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Weight:        product.Weight,
		Images:        product.Images,
		UserID:        uid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(newProduct); err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (s *ProductService) GetAllProducts(limit, offset int, filter map[string]any) ([]*models.Product, error) {
	return s.repo.GetProducts(limit, offset, filter)
}

func (s *ProductService) GetById(id string) (*models.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *ProductService) UpdateProduct(id string, updates map[string]any) error {
	return s.repo.Update(id, updates)
}
