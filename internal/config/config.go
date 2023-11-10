package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/joho/godotenv"
)

type config struct{
	rabbitmq_host string
	rabbitmq_port string
	rabbitmq_username string
	rabbitmq_password string
	rabbitmq_queue string
	rabbitmq_url string
}

func LoadEnv() config{
	err := godotenv.Load("production.env")
	if err != nil {
		log.Info("Some error occured. Err: ", err)
	}

	return config{
		rabbitmq_host: os.Getenv("RABBITMQ_HOST"),
		rabbitmq_port: os.Getenv("RABBITMQ_PORT"),
		rabbitmq_username: os.Getenv("RABBITMQ_USERNAME"),
		rabbitmq_password: os.Getenv("RABBITMQ_PASSWORD"),
		rabbitmq_queue: os.Getenv("RABBITMQ_QUEUE"),
		rabbitmq_url: os.Getenv("RABBITMQ_URL"),
	}
}