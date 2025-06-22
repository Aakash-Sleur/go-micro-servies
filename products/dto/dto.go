package dto

type ProductDto struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	Price         float64  `json:"price" binding:"required,gt=0"`
	StockQuantity int      `json:"stock_quantity" binding:"required,gte=0"`
	Weight        float64  `json:"weight" binding:"required,gte=0"`
	Images        []string `json:"images" binding:"required"`
}

type CartDto struct {
	ProductID string
	Quantity  int
}
