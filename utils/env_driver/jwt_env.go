package env_driver

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type JWTEnvironment struct {
	Secret        string
	Duration      int
	SignAlgorithm string
}

func NewJWTEnvironmentDriver() (*JWTEnvironment, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	duration, err := strconv.Atoi(os.Getenv("JWT_DURATION"))
	if err != nil {
		return nil, err
	}

	return &JWTEnvironment{
		Secret:        os.Getenv("JWT_SECRET"),
		Duration:      duration,
		SignAlgorithm: os.Getenv("JWT_SECRET"),
	}, nil
}
