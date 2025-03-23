package authservice

import (
	"auth_service/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	store *UserStore
}

func NewUserController(store *UserStore) *UserController {
	return &UserController{
		store: store,
	}
}

func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", c.registerHandler)
	router.POST("/login", c.loginHandler)

}

func (c *UserController) registerHandler(ctx *gin.Context) {
	var payload dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// check if user exists
	_, err := c.store.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		ctx.JSON(400, gin.H{"error": "User already exists"})
		return
	}
	ctxx := context.Background()
	apiKey := GenerateAPIKey()

	// Hash the password
	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// create new user object from the payload with hashed password
	user := NewUser(payload.Username, hashedPassword, payload.Email, apiKey)
	if err := c.store.CreateUser(ctxx, user); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User registered successfully"})
}

func (c *UserController) loginHandler(ctx *gin.Context) {
	// get the payload from the request body
	var payload dto.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctxx := context.Background()
	// get the user from the database
	user, err := c.store.GetUser(ctxx, payload.Username)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// check if the password is correct
	if ComparePassword(user.Password, payload.Password) == false {
		ctx.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	apiKey := user.ApiKey

	// return the api key
	ctx.JSON(200, gin.H{"api_key": apiKey})
}
