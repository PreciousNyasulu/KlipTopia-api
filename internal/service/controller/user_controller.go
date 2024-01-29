package controller

import (
	"errors"
	"fmt"
	conf "kliptopia-api/internal/config"
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/service/auth"
	"net/http"
	"strings"
	"time"

	mr_rabbit "kliptopia-api/internal/rabbitmq_processes"

	"kliptopia-api/internal/repository"

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
			if err.Error() == "user already logged in" {
				c.JSON(http.StatusOK,gin.H{"message": "user already authenticated"})
				return
			}			
		}

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
	DB, err := repository.Connect()
	if err != nil {
		logger.Error("Failed to connect to the data repository")
	}
	defer repository.CloseDB()

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

	var (
		user_id int 
		result_count int64
	)

	user_id = getUserId(user.Username)
	
	DB.Table("auth_tokens").Where(models.AuthToken{User_Id: user_id}).Where("expired_at is NULL").Count(&result_count)
	if result_count > 0 {
		logger.Info(fmt.Sprintf("User %s already authenticated",user.Username))
		return "", errors.New("user already logged in")
	}

	token_repository := models.AuthToken{
		User_Id: user_id,
		Token: tokenString,
		Created_At: time.Now(),
	}

	if err := DB.Create(&token_repository).Error; err != nil{
		logger.Warn(fmt.Sprintf("Token generated but failed to save to the data repository, %v", err))
	}

	return tokenString, nil
}

func LogoutHandler(c *gin.Context){
	DB, err := repository.Connect()
	if err != nil {
		logger.Error("Failed to connect to the data repository")
	}
	defer repository.CloseDB()

	username, _ := c.Get("username")
	err = DB.Table("auth_tokens").Where(models.AuthToken{User_Id: getUserId(fmt.Sprintf("%s",username))}).Where("expired_at is NULL").Update("expired_at",time.Now()).Error
	if err != nil {
		logger.Error("Failed to invalidate the token")
		c.JSON(http.StatusInternalServerError,gin.H{"message":"Failed to invalidate the token"})
		return 
	}
	c.JSON(http.StatusOK, gin.H{"message":"Logout"})
}

func getUserId(username string) int {
	DB, err := repository.Connect()
	if err != nil {
		logger.Error("Failed to connect to the data repository")
	}
	defer repository.CloseDB()

	var user_id int
	//retrieve user_id
	DB.Table("users").Where(models.User{Username: username}).Or(models.User{Email: username}).Select("user_id").Find(&user_id)
	return user_id
}