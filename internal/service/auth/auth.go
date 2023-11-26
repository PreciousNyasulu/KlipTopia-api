package auth

import (
	// "encoding/json"
	"kliptopia-api/internal/models"
	"kliptopia-api/internal/repository"
	"kliptopia-api/internal/utils"
	"time"

	"github.com/google/uuid"
)

var DB, _ = repository.Connect()

// GetUser retrieves a user by ID and returns it as JSON
func GetUser() (models.User, error) {
	var user models.User
	if err := DB.First(&user, 1).Error; err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(RequestBody models.AuthRequestBody) (bool,error){
	passwordHash,err := utils.EncryptPassword(RequestBody.Password)
	if err != nil {
		return false, err
	}
	
	user := models.User{
		Username: RequestBody.Username,
		Email: RequestBody.Email,
		PasswordHash: passwordHash,
		FirstName: RequestBody.FirstName,
		LastName: RequestBody.LastName,
		ProfilePictureURL: "",
		RegistrationDate: time.Now(),
		LastLoginDate: time.Now(),
		Role: "user",
		AccountStatus: "unverified",		
		VerificationStatus: "unverified",
		PasswordResetToken: uuid.Nil,
		PasswordResetExpiry: time.Now().AddDate(0,3,0), //three months
		TwoFactorEnabled: false,
	}

	if err := DB.Create(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func CheckUser(email string) bool{
	var rowCount int64
	DB.Table("users").Where("email",email).Count(&rowCount)
	return rowCount > 0
}

func Login(user models.AuthRequestBody) string{
	var results []models.User
	DB.Where(models.User{Email: user.Email}).Or(models.User{Username: user.Username}).Find(&results)
	
	for _, _user := range results {
		passwordIsCorrect := utils.VerifyPassword(user.Password,_user.PasswordHash)
		if passwordIsCorrect == nil {
			return "success"
		}
		return "invalid"
	}
	return "not found"
}

