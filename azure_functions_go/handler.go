package main

import (
	"azure_functions_go/utils"
	"log"
	"net/http"
	"os"

	authservice "auth_service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getRouteHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	data := map[string]interface{}{
		"Greetings": "Hello " + name,
	}
	c.JSON(http.StatusOK, utils.Response(200, "Please hold while we process your safest route to your destination", data))
}

func main() {
	//  create a new router
	router := gin.Default()

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "4300"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	api := router.Group("/api")
	api.POST("/login", authservice.LoginHandler)

	protected := api.Group("/")
	protected.Use(authservice.AuthMiddleware())
	{
		protected.GET("/safe_route_function", getRouteHandler)
	}

	log.Printf("Starting Gin-based Azure Function on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
