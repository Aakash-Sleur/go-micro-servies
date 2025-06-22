package routes

import (
	"github.com/Aakash-Sleur/go-micro-order/config"
	"github.com/Aakash-Sleur/go-micro-order/controllers"
	"github.com/Aakash-Sleur/go-micro-order/middleware"
	"github.com/Aakash-Sleur/go-micro-order/repository"
	"github.com/Aakash-Sleur/go-micro-order/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	db := config.DB

	repo := repository.NewOrderRepository(db)
	service := services.NewOrderService(repo)
	controller := controllers.NewOrderController(service)

	api := r.Group("/api")

	orderRoutes := api.Group("/orders/v1")
	{
		orderRoutes.POST("", middleware.JWTMiddleware(), controller.CreateOrder)
		orderRoutes.GET("", controller.GetAllOrders)
		orderRoutes.GET("/:id", controller.GetOrderByID)
		orderRoutes.PUT("/:id/cancel", middleware.JWTMiddleware(), controller.CancelOrder)
	}
}
