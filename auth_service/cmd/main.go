package main

import (
	"log"
	"os"

	authservice "auth_service"
	handler "azure_functions_go"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		protected.GET("/safe_route_function", handler.GetRouteHandler)
	}

	log.Printf("Starting Gin-based Azure Function on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
