package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

func ParseFilterString(filterStr string) (map[string]interface{}, error) {
	// This function should parse the filter string into a map.
	// For simplicity, let's assume the filterStr is in JSON format.
	var filter map[string]interface{}
	err := json.Unmarshal([]byte(filterStr), &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}
