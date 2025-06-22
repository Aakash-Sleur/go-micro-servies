package repository

import (
	"github.com/Aakash-Sleur/go-micro-order/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetAll(filter map[string]any) ([]*models.Order, error) {
	var orders []*models.Order
	query := r.db.Model(&models.Order{})
	if len(filter) > 0 {
		query.Where(filter)
	}

	query.Order("created_at DESC")

	err := query.Preload("Items").Find(&orders).Error

	return orders, err
}

func (r *OrderRepository) GetById(id string) (*models.Order, error) {
	var order *models.Order

	err := r.db.Preload("Items").Where("id = ?", id).First(&order).Error
	return order, err
}
func (r *OrderRepository) Cancel(orderId string) (*models.Order, error) {
	var order *models.Order
	err := r.db.First(&order, "id = ?", orderId).Error
	if err != nil {
		return nil, err
	}

	order.Status = "cancelled"
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}
