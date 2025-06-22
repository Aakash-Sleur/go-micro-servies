package services

import (
	"fmt"
	"time"

	"github.com/Aakash-Sleur/go-micro-order/dto"
	"github.com/Aakash-Sleur/go-micro-order/models"
	"github.com/Aakash-Sleur/go-micro-order/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(orderData dto.OrderDTO, userId string) (*models.Order, error) {
	uid, err := uuid.Parse(userId)
	if err != nil {
		// handle error appropriately, e.g., return or log
		return nil, err
	}
	order := &models.Order{
		ID:          uuid.New(),
		UserID:      uid,
		OrderNumber: uuid.New().String(),
		Status:      orderData.Status,
		TotalAmount: orderData.TotalAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	fmt.Print(orderData.Items[0])
	for _, item := range orderData.Items {
		orderItem := models.OrderItem{
			ID:         uuid.New(),
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.UnitPrice * float64(item.Quantity),
			// Set ProductSnapshot from item.Images
			Images:    item.Images,
			CreatedAt: time.Now(),
		}

		order.Items = append(order.Items, orderItem)
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrders(filter map[string]any) ([]*models.Order, error) {
	return s.repo.GetAll(filter)
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	return s.repo.GetById(id)
}

func (s *OrderService) CancelOrder(id string) (*models.Order, error) {
	return s.repo.Cancel(id)
}
