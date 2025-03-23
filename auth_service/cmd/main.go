package main

import (
	authservice "auth_service"
	handler "azure_functions_go"
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Set up environment variables
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	addr := ":" + port
	dbUrl := authservice.GetEnv("DB_URL", "mongodb://localhost:27017")

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	dbClient, err := authservice.ConnectDB(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Create a new Gin router
	router := gin.Default()

	// Set up CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Set up auth service routes
	mongoDb := dbClient.Database("authservice")
	userStore := authservice.NewUserStore(mongoDb)
	userController := authservice.NewUserController(userStore)

	// Register auth routes
	authGroup := router.Group("/api/v1")
	userController.RegisterRoutes(authGroup)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authservice.AuthMiddleware(userStore))
	{
		protected.GET("/safe_route", handler.GetRouteHandler)
		protected.GET("/search_route", handler.SearchRouteHandler)
		protected.GET("/fastest_route", handler.SearchFastestRouteHandler)
	}

	// Start the server
	log.Printf("Starting City Scout server on port %s...\n", port)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
