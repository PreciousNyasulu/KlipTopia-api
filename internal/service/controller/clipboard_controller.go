package controller

import (
	"kliptopia-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Copy(c *gin.Context) {
	var clipCopyData models.QueueMessage

	if err := c.ShouldBindJSON(&clipCopyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&clipCopyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	

}
