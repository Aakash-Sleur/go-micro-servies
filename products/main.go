package main

import (
	"log"
	"net/http"

	"github.com/Aakash-Sleur/go-micro-product/config"
	"github.com/Aakash-Sleur/go-micro-product/routes"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add debug logging
		log.Printf("CORS Middleware executed for: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Origin: %s", c.Request.Header.Get("Origin"))

		// Set CORS headers for ALL requests
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // For debugging - change back to specific origin later
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		log.Printf("CORS headers set")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS request")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	config.ConnectToDB()
	config.SyncDatabase()

	r := gin.Default()

	// CRITICAL: Disable redirects that cause 301 responses
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// Add CORS middleware FIRST - this is critical
	r.Use(CORSMiddleware())

	// Add a test route to verify CORS is working
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "CORS test successful"})
	})
	// Then add your routes
	routes.SetupRoutes(r)

	// Debug: Print all registered routes
	log.Println("Registered routes:")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	port := config.Load().PORT
	log.Printf("Service started on port %s", port)
	r.Run(":" + port)
}
