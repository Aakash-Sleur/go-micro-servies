package config

import (
	"log"

	"github.com/Aakash-Sleur/go-micro-order/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := Load().DB_URI
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB")
	}

	log.Println("Connected to DB")

	DB = db
}

func SyncDatabase() {
	// DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.Product{})
	// DB.AutoMigrate(&models.CartItem{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItem{})
}
