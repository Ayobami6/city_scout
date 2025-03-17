package authservice

import (
	"net/http"
	"os"
	"strings"

	"azure_functions_go/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response(
				401,
				"Unauthorized: Bearer token required",
				nil,
			))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		expectedToken := os.Getenv("AUTH_TOKEN")
		if expectedToken == "" {
			expectedToken = "your-default-secret-token" // Default for development, use env in production
		}

		if tokenString != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response(
				401,
				"Unauthorized: Invalid token",
				nil,
			))
			return
		}
		c.Next()
	}

}
