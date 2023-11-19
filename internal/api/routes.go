package api

import (
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(){
	r := gin.Default()

	r.POST("/health-test",func (c *gin.Context)  {
		c.JSON(http.StatusOK,models.Health_check{Healthy: true})
	})

	// auth routes
	r.POST("/auth/register",controller.CreateUserHandler)

	r.Run(":9000")
}