package models

import (
	// "database/sql"
	"time"
	"github.com/google/uuid"
)

type Health_check struct{
	Healthy bool `json:"Healthy"`
}

// Config struct to hold the configuration values
type Config struct {
	RabbitMQ RabbitMQConfig
	Postgres PostgresConfig
}

// RabbitMQConfig struct to hold RabbitMQ configuration
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Queue    string
	Url		 string
}

// PostgresConfig struct to hold PostgreSQL configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type User struct {
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"password_hash"`
	FirstName          string 	  `json:"first_name"`
	LastName           string `json:"last_name"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	RegistrationDate   time.Time  `json:"registration_date"`
	LastLoginDate      time.Time  `json:"last_login_date"`
	Role               string     `json:"role"`
	AccountStatus      string     `json:"account_status"`
	VerificationStatus string     `json:"verification_status"`
	PasswordResetToken uuid.UUID  `json:"password_reset_token"`
	PasswordResetExpiry time.Time `json:"password_reset_expiry"`
	TwoFactorEnabled   bool       `json:"two_factor_enabled"`
	// OAuthProvider      string `db:"oauth_provider"`
	// OAuthID            string `json:"oauth_id"`
}

type AuthRequestBody struct{
	Username 	string `json:"username" validate:"required,min=3,max=20"`
	FirstName	string `json:"firstname" validate:"required,min=3,max=20"`
	LastName	string `json:"lastname" validate:"required,min=3,max=20"`
	Email 		string `json:"email" validate:"required,email"`
	Password 	string `json:"password" validate:"required,min=8,max=16"`
}

type SessionData struct{
	UserID		string
	Email		string
	Username	string
}