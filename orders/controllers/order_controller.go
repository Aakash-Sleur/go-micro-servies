package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Aakash-Sleur/go-micro-order/dto"
	"github.com/Aakash-Sleur/go-micro-order/models"
	"github.com/Aakash-Sleur/go-micro-order/services"
	"github.com/Aakash-Sleur/go-micro-order/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController(service *services.OrderService) *OrderController {
	return &OrderController{service: service}
}

// POST /orders
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req dto.OrderDTO

	userInterface, exists := ctx.Get("currentUser")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authorized")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to parse user from context")
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Creating order for user %s with payload: %+v", user.ID.String(), req)

	order, err := c.service.CreateOrder(req, user.ID.String())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Order created successfully", gin.H{
		"order": order,
	})
}

// GET /orders/:id
func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := c.service.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Order not found")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Order fetched successfully", gin.H{
		"order": order,
	})
}

// GET /orders
func (c *OrderController) GetAllOrders(ctx *gin.Context) {
	filter := make(map[string]any)

	// Extract filter[field]=value from query parameters
	for key, values := range ctx.Request.URL.Query() {
		if len(values) == 0 {
			continue
		}

		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			filterKey := key[7 : len(key)-1]

			// If multiple values per filter key are supported:
			if len(values) > 1 {
				filter[filterKey] = values
			} else {
				filter[filterKey] = values[0]
			}
		}
	}

	// Optional: Parse filter JSON string (if passed as a full JSON in `filter`)
	if filterStr := ctx.Query("filter"); filterStr != "" {
		parsedFilter, err := utils.ParseFilterString(filterStr)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid filter format")
			return
		}

		// Merge JSON-parsed filter into existing filter
		for k, v := range parsedFilter {
			filter[k] = v
		}
	}

	orders, err := c.service.GetOrders(filter)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Orders fetched successfully", gin.H{
		"orders": orders,
	})
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := c.service.CancelOrder(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Order not found")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Order cancelled successfully", gin.H{
		"order": order,
	})
}
