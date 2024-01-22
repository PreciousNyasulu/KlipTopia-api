package config

import (
	"kliptopia-api/internal/models"
	"os"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadConfig() models.Config {
	err := godotenv.Load("../.env")
	logger :=log.New()

	if err != nil {
		logger.Error("Error loading .env file: %s", err)
		panic(nil)
	}

	return models.Config{
		RabbitMQ: models.RabbitMQConfig{
			Host:     os.Getenv("RABBITMQ_HOST"),
			Port:     os.Getenv("RABBITMQ_PORT"),
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASSWORD"),
			Queue:    os.Getenv("RABBITMQ_QUEUE"),
			Url: 	  os.Getenv("RABBITMQ_URL"),
		},
		Postgres: models.PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
		Authentication: models.AuthConfig{
			TOKEN_SIGNING_SECRET: os.Getenv("TOKEN_SIGNING_SECRET"),
		},
	}
}
