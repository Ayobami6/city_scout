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
	ctxx := context.Background()
	apiKey := GenerateAPIKey()
	// create new user object from the payload
	user := NewUser(payload.Username, payload.Password, payload.Email, apiKey)
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
	user, err := c.store.GetUserByEmail(ctxx, payload.Email)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// check if the password is correct
	if user.Password != payload.Password {
		ctx.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	// generate a jwt token
	token, err := GenerateJWT(user.Email, user.ApiKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// return the token
	ctx.JSON(200, gin.H{"token": token})
}
