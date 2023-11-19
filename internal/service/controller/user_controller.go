package controller

import (
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetAllUsersHandler(c *gin.Context){
	jsondata,_ :=auth.GetUser()
	c.JSON(http.StatusOK,jsondata)
}

func CreateUserHandler(c *gin.Context){
	var user models.AuthRequestBody

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the user struct
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert user into the database
	_,err := auth.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"message":"Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}