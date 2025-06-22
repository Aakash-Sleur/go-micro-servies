package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Aakash-Sleur/go-micro-product/dto"
	"github.com/Aakash-Sleur/go-micro-product/models"
	"github.com/Aakash-Sleur/go-micro-product/service"
	"github.com/Aakash-Sleur/go-micro-product/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	service *service.ProductService
}

func NewProductService(service *service.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req dto.ProductDto

	userInterface, exists := ctx.Get("currentUser")

	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	user := userInterface.(models.User)

	fmt.Println(userInterface, "userinterface")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	product, err := c.service.CreateProduct(req, user.ID.String())

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Product Created successfully", gin.H{"product": product})
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	// Parse pagination params
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid limit")
		return
	}

	skip, err := strconv.Atoi(ctx.DefaultQuery("skip", "0"))
	if err != nil || skip < 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid skip")
		return
	}

	// Initialize filter map
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

	// Fetch products
	products, err := c.service.GetAllProducts(limit, skip, filter)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Fetched successfully", gin.H{
		"products": products,
		"skip":     skip,
	})
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Product ID is required")
		return
	}

	product, err := c.service.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Product not found")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Product fetched successfully", gin.H{
		"product": product,
	})
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Product ID is required")
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Convert "images" field manually if it exists
	if rawImages, ok := updates["images"]; ok {
		interfaceSlice, ok := rawImages.([]interface{})
		if !ok {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "'images' must be an array of strings")
			return
		}

		stringSlice := make(models.StringSlice, len(interfaceSlice))
		for i, v := range interfaceSlice {
			str, ok := v.(string)
			if !ok {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "all 'images' must be strings")
				return
			}
			stringSlice[i] = str
		}
		updates["images"] = stringSlice
	}

	// Null check for all other fields
	for k, v := range updates {
		if v == nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("Field '%s' cannot be null", k))
			return
		}
	}

	updates["updated_at"] = time.Now()

	if err := c.service.UpdateProduct(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Product not found")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Product updated successfully", nil)
}
