package service

import (
	"time"

	"github.com/Aakash-Sleur/go-micro-product/models"
	"github.com/Aakash-Sleur/go-micro-product/repository"
	"github.com/google/uuid"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) CreateCartItem(userID, productID uuid.UUID, quantity int) (*models.CartItem, error) {
	cartItem := &models.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(cartItem)
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (s *CartService) GetCartItems(limit, skip int64, userId string, filter map[string]any) ([]*models.CartItem, error) {
	return s.repo.GetCartItems(limit, skip, userId, filter)
}

func (s *CartService) ClearCart(userId string) error {
	return s.repo.ClearCart(userId)
}

func (s *CartService) Remove(id string, userId string) error {
	return s.repo.Remove(id, userId)
}
