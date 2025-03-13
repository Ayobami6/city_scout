package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getRouteHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello " + name,
	})

}

func main() {
	//  create a new router
	router := gin.Default()

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Proctor-Env"},
		AllowCredentials: true,
	}))
	api := router.Group("/api")

	api.GET("/getRouteFunction", getRouteHandler)
	log.Printf("Starting Gin-based Azure Function on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
