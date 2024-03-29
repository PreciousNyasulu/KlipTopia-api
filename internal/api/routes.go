package api

import (
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(){
	r := gin.Default()

	r.GET("/health-check",func (c *gin.Context)  {
		c.JSON(http.StatusOK,models.Health_check{Healthy: true})
	})

	// auth routes
	r.POST("/api/auth/register",controller.CreateUserHandler)
	r.POST("/api/auth/login",controller.LoginHandler)
	r.POST("/api/auth/logout" ,controller.AuthMiddleware(),controller.LogoutHandler)

	//clipboard event
	r.POST("/api/clipboard/copy",controller.AuthMiddleware(), controller.Copy)

	r.Run(":9000")
}