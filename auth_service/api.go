package authservice

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"azure_functions_go/utils"

	"github.com/gin-gonic/gin"
)

var tokens []string

// Generate random token
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func LoginHandler(c *gin.Context) {
	// Basic auth middleware would go here in production
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth || username != "admin" || password != "secret" {
		c.JSON(http.StatusUnauthorized, utils.Response(401, "Invalid credentials", nil))
		return
	}

	token, err := randomHex(20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response(500, "Error generating token", nil))
		return
	}

	tokens = append(tokens, token)
	c.JSON(http.StatusOK, utils.Response(200, "Login successful", map[string]interface{}{
		"token": token,
	}))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response(401, "Unauthorized", nil))
			return
		}

		reqToken := strings.Split(bearerToken, " ")[1]
		for _, token := range tokens {
			if token == reqToken {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response(401, "Invalid token", nil))
	}
}
