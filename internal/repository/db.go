package repository

import (
	"fmt"
	"kliptopia-api/internal/config"
	"strconv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	config := config.LoadConfig()
	logger := log.New()
	password := config.Postgres.Password
	host := config.Postgres.Host
	user := config.Postgres.User
	dbname := config.Postgres.Database

	port, err := strconv.Atoi(config.Postgres.Port)
	if err != nil {
		port = 5432
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password,dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database. %v", err)
		return nil, err
	}
	return db, nil
}