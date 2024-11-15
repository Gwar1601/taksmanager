// internal/handlers/auth.go
package handlers

import (
	"net/http"
	"taskmanager/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	// Implementacja rejestracji użytkownika
}

func LoginHandler(c *gin.Context) {
	// Implementacja logowania użytkownika i generowania tokena JWT
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
