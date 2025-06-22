package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Aakash-Sleur/go-micro-product/models"
	"github.com/Aakash-Sleur/go-micro-product/service"
	"github.com/Aakash-Sleur/go-micro-product/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartController struct {
	service *service.CartService
}

func NewCartController(serice *service.CartService) *CartController {
	return &CartController{service: serice}
}

func (c *CartController) CreateCartItem(ctx *gin.Context) {
	var req struct {
		ProductId string
		quantity  int
	}

	userInterface, exists := ctx.Get("currentUser")

	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	user := userInterface.(models.User)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid product ID format")
		return
	}
	cartItem, err := c.service.CreateCartItem(user.ID, productUUID, req.quantity)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Product Created successfully", gin.H{"cartItem": cartItem})
}

func (c *CartController) GetCartItems(ctx *gin.Context) {
	var req struct {
		Limit int64 `form:"limit"`
		Skip  int64 `form:"skip"`
	}

	userInterface, exists := ctx.Get("currentUser")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	user := userInterface.(models.User)
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filter := make(map[string]interface{})

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

	fmt.Println("Final Filter:", filter)

	cartItems, err := c.service.GetCartItems(req.Limit, req.Skip, user.ID.String(), filter)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Cart items retrieved successfully", gin.H{"cartItems": cartItems})
}

func (c *CartController) RemoveCartItem(ctx *gin.Context) {
	itemId := ctx.Param("id")
	if itemId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "item ID is required")
		return
	}

	userInterface, exists := ctx.Get("currentUser")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	user := userInterface.(models.User)

	err := c.service.Remove(itemId, user.ID.String())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Cart item removed successfully", nil)
}

func (c *CartController) ClearCart(ctx *gin.Context) {
	userInterface, exists := ctx.Get("currentUser")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	user := userInterface.(models.User)

	err := c.service.ClearCart(user.ID.String())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Cart item removed successfully", nil)
}
