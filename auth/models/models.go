package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Product struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `json:"user_id" gorm:"not null"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null"`
	Price         float64        `json:"price" gorm:"not null"`
	StockQuantity int            `json:"stock_quantity" gorm:"default:0"`
	Weight        float64        `json:"weight"`
	Images        pq.StringArray `json:"images" gorm:"type:text[]"` // for text[]
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID   `json:"user_id" gorm:"not null"`
	OrderNumber string      `json:"order_number" gorm:"uniqueIndex;not null"`
	Status      string      `json:"status" gorm:"default:'pending'"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderItem model
type OrderItem struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID    uuid.UUID `json:"order_id" gorm:"not null"`
	ProductID  uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity   int       `json:"quantity" gorm:"not null"`
	UnitPrice  float64   `json:"unit_price" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	Images     any       `json:"product_snapshot" gorm:"type:jsonb"`
	Product    Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt  time.Time `json:"created_at"`
}

type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
