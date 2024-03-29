package controller

import (
	"fmt"
	"kliptopia-api/internal/models"
	"net/http"
	"time"

	mr_rabbit "kliptopia-api/internal/rabbitmq_processes"

	"github.com/gin-gonic/gin"
)

func Copy(c *gin.Context) {
	var clipCopyData models.QueueMessage

	clipCopyData.CopiedAt = time.Now()
	if err := c.ShouldBindJSON(&clipCopyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to parse request message, %s", err.Error())})
		return
	}

	if err := validate.Struct(&clipCopyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	_username, _ := c.Get("username")
	username := fmt.Sprintf("%s",_username)
	if err := mr_rabbit.PushMessageToQueue(username, &clipCopyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message published"})
}
