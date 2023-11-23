package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func EncryptPassword(password string) (string,error){
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
}

func GenerateUUID() uint32 {
	uuid := uuid.New()
	return uuid.ID()
}