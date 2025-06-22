package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StringSlice is a custom type for handling []string in GORM
type StringSlice []string

// Value implements the driver.Valuer interface for database storage
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	jsonData, err := json.Marshal(s)
	return string(jsonData), err
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), s)
	case []byte:
		return json.Unmarshal(v, s)
	default:
		return errors.New("cannot scan into StringSlice")
	}
}

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
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID   `json:"user_id" gorm:"not null"`
	Name          string      `json:"name" gorm:"not null"`
	Description   string      `json:"description" gorm:"not null"`
	Price         float64     `json:"price" gorm:"not null"`
	StockQuantity int         `json:"stock_quantity" gorm:"default:0"`
	Weight        float64     `json:"weight"`
	Images        StringSlice `json:"images" gorm:"type:jsonb"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
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
