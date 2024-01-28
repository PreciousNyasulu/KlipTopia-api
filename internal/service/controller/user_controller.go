package controller

import (
	"fmt"
	conf "kliptopia-api/internal/config"
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/auth"
	"net/http"
	"strings"
	"time"

	mr_rabbit "kliptopia-api/internal/rabbitmq_processes"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var validate = validator.New()
var config = conf.LoadConfig()
var logger = *log.New()

func GetAllUsersHandler(c *gin.Context) {
	jsondata, _ := auth.GetUser()
	c.JSON(http.StatusOK, jsondata)
}

func CreateUserHandler(c *gin.Context) {
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

	if auth.CheckUser(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists"})
		return
	}

	// Insert user into the database
	_, err := auth.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginHandler(c *gin.Context) {
	var user models.AuthRequestBody

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch auth.Login(user) {
	case "success":
		token, err := GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
			return
		}

		_, err = mr_rabbit.CreateSessionQueue(user.Username)
		if err == nil {
			logger.Info("created queue for: ", user.Username)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
		return

	case "invalid":
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong password"})
		return
	case "not found":
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"message": "Something happened while processing your request."})
}

// AuthMiddleware is a middleware that checks for authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if the header is missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the header is in the format "Bearer [token]"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the token
		tokenString := parts[1]
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Authentication.TOKEN_SIGNING_SECRET), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"message": fmt.Sprintf("Unauthorized, %v",err)})
			c.Abort()
			return
		}

		// Token is valid, continue to the next handler
		c.Set(contextUsernameKey, claims.Subject)
		c.Next()
	}
}

const (
	contextUsernameKey = "username"
)

// GenerateToken generates a new JWT token for the given user.
func GenerateToken(user models.AuthRequestBody) (string, error) {

	//dont know if this is the correct way to do it (sigh) :(
	secretKey := []byte(config.Authentication.TOKEN_SIGNING_SECRET)

	// generate a jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Username,                         //not sure about setting the email as the subject though
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours, sounds like a good idea for now
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
