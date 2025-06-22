package routes

import (
	"github.com/Aakash-Sleur/go-micro-auth/config"
	"github.com/Aakash-Sleur/go-micro-auth/controllers"
	"github.com/Aakash-Sleur/go-micro-auth/middleware"
	"github.com/Aakash-Sleur/go-micro-auth/repository"
	"github.com/Aakash-Sleur/go-micro-auth/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewAuthService(userRepo)
	userController := controllers.NewUserController(userService)

	api := router.Group("/api")
	authRoutes := api.Group("/auth/v1")
	{
		authRoutes.POST("/signup", userController.SignUp)
		authRoutes.POST("/signin", userController.Signin)
		authRoutes.GET("/current-user", middleware.JWTMiddleware(), userController.GetCurrentUser)
	}
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})
}
