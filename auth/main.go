package main

import (
	"log"
	"time"

	"github.com/Aakash-Sleur/go-micro-auth/config"
	"github.com/Aakash-Sleur/go-micro-auth/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectToDB()
	config.SyncDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	port := config.Load().PORT
	log.Printf("Service started on port %s", port)

	r.Run(":" + port)
}
