package routes

import (
	"github.com/Aakash-Sleur/go-micro-product/config"
	"github.com/Aakash-Sleur/go-micro-product/controllers"
	"github.com/Aakash-Sleur/go-micro-product/middleware"
	"github.com/Aakash-Sleur/go-micro-product/repository"
	"github.com/Aakash-Sleur/go-micro-product/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Product Routes
	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productController := controllers.NewProductService(productService)

	api := router.Group("/api")
	productRoutes := api.Group("/products/v1")
	{
		// Handle OPTIONS requests for CORS preflight - FIXED: No trailing slashes
		productRoutes.OPTIONS("", func(c *gin.Context) {
			c.Status(200)
		})
		productRoutes.OPTIONS("/:id", func(c *gin.Context) {
			c.Status(200)
		})

		// FIXED: Remove trailing slashes from all routes
		productRoutes.POST("", middleware.JWTMiddleware(), productController.CreateProduct)
		productRoutes.GET("", productController.GetAllProducts)
		productRoutes.GET("/:id", productController.GetProductById)
		productRoutes.PUT("/:id", middleware.JWTMiddleware(), productController.UpdateProduct)
	}

	// Cart Routes
	cartRepo := repository.NewCartRepository(config.DB)
	cartService := service.NewCartService(cartRepo)
	cartController := controllers.NewCartController(cartService)
	cartRoutes := api.Group("/cart/v1")
	{
		// Handle OPTIONS requests for CORS preflight - FIXED: No trailing slashes
		cartRoutes.OPTIONS("", func(c *gin.Context) {
			c.Status(200)
		})
		cartRoutes.OPTIONS("/:id", func(c *gin.Context) {
			c.Status(200)
		})
		// FIXED: Remove trailing slashes from all routes
		cartRoutes.POST("", middleware.JWTMiddleware(), cartController.CreateCartItem)
		cartRoutes.GET("", middleware.JWTMiddleware(), cartController.GetCartItems)
		cartRoutes.DELETE("", middleware.JWTMiddleware(), cartController.ClearCart)
		cartRoutes.DELETE("/:id", middleware.JWTMiddleware(), cartController.RemoveCartItem)
	}
}
