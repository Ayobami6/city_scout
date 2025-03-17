package handler

import (
	"azure_functions_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRouteHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	data := map[string]interface{}{
		"Greetings": "Hello " + name,
	}
	c.JSON(http.StatusOK, utils.Response(200, "Please hold while we process your safest route to your destination", data))
}
