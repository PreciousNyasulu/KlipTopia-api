package controller

import (
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/auth"
	"net/http"
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
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

	if auth.CheckUser(user.Email) {
		c.JSON(http.StatusBadRequest,gin.H{"message":"User already exists"})
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

func LoginHandler(c *gin.Context){
	var user models.AuthRequestBody

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch auth.Login(user){
	case "success":
		token,err :=GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"message":"failed to generate token"})
		}
		c.JSON(http.StatusOK,gin.H{"message":"success","token":token})
		return
	case "invalid":
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Wrong password"})
		return
	case "not found":
		c.JSON(http.StatusNotFound,gin.H{"message":"User not found."})
		return
	}

	c.JSON(http.StatusInternalServerError,gin.H{"message":"Something happened while processing your request."})
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

		// Check if the header is in the format "Bearer token"
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
			// TODO: Add your secret key for token validation
			return []byte("yourSecretKey"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Token is valid, continue to the next handler
		c.Set("username", claims.Subject)
		c.Next()
	}
}

// GenerateToken generates a new JWT token for the given user.
func GenerateToken(user models.AuthRequestBody) (string, error) {
	// TODO: Add your secret key for token signing
	secretKey := []byte(user.Email)
	
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}