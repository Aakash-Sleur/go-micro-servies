package dto

import "github.com/google/uuid"

type OrderItemDTO struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"price"` // maps to "price" in JSON
	Images    []string  `json:"images"`
}

type OrderDTO struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	OrderNumber string         `json:"order_number"`
	Status      string         `json:"status"`
	TotalAmount float64        `json:"total_amount"`
	Items       []OrderItemDTO `json:"items"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}
