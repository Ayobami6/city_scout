package authservice

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid API key in the request headers
func AuthMiddleware(store *UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the format is "Bearer {api_key}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {api_key}"})
			c.Abort()
			return
		}

		apiKey := parts[1]
		
		// Verify API key against the database
		user, err := store.GetUserByAPIKey(context.Background(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Store user information in the context for later use
		c.Set("user", user)
		c.Next()
	}
}

// GenerateMiddleware creates and returns the auth middleware function
func GenerateMiddleware() gin.HandlerFunc {
	// This is a placeholder - you need to modify main.go to use AuthMiddleware instead
	return func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Middleware not properly configured"})
		c.Abort()
	}
}
