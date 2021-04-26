package env_driver

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseEnvironment struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func NewDatabaseEnvironmentDriver() (*DatabaseEnvironment, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &DatabaseEnvironment{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}, nil
}
