package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Aakash-Sleur/go-micro-auth/config"
	"github.com/Aakash-Sleur/go-micro-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Read token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		fmt.Print(parts, len(parts), strings.ToLower(parts[0]), "<--parts")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			return
		}

		tokenString := parts[1]
		fmt.Println("JWT Token from header:", tokenString)

		// 2. Parse token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.Load().JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 3. Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["userId"] == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// 4. Extract userId as string
		userId, ok := claims["userId"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
			return
		}

		// 5. Find user from DB
		var user models.User
		if err := config.DB.First(&user, "id = ?", userId).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// 6. Set user in context
		c.Set("currentUser", user)

		c.Next()
	}
}
