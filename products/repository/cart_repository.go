package repository

import (
	"github.com/Aakash-Sleur/go-micro-product/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Create(cart *models.CartItem) error {
	return r.db.Create(cart).Error
}

func (r *CartRepository) GetCartItems(limit, skip int64, userId string, filter map[string]any) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem

	query := r.db.Model(&models.CartItem{}).Preload("Product")
	if userId != "" {
		query = query.Where("user_id = ?", userId)
	}

	if len(filter) > 0 {
		query = query.Where(filter)
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	if skip > 0 {
		query = query.Offset(int(skip))
	}

	err := query.Find(&cartItems).Error
	return cartItems, err
}

func (r *CartRepository) Remove(itemId string, userId string) error {
	return r.db.Delete(&models.CartItem{}, itemId).Error
}

func (r *CartRepository) ClearCart(userId string) error {
	return r.db.Where("user_id = ?", userId).Delete(&models.CartItem{}).Error
}
