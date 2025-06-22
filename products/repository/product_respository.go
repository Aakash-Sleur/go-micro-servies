package repository

import (
	"github.com/Aakash-Sleur/go-micro-product/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetProductById(id string) (*models.Product, error) {
	var product *models.Product

	err := r.db.Where("id = ?", id).First(&product).Error
	return product, err
}

func (r *ProductRepository) GetProducts(limit, offset int, filter map[string]interface{}) ([]*models.Product, error) {
	var products []*models.Product

	query := r.db.Model(&models.Product{})

	if len(filter) > 0 {
		query = query.Where(filter)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(id string, updates map[string]any) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}
