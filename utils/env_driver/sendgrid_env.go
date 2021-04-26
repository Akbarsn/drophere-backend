package env_driver

import (
	"os"

	"github.com/joho/godotenv"
)

type SendgridMailer struct {
	APIKey string
}

func NewSendgridMailerDriver() (*SendgridMailer, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &SendgridMailer{
		APIKey: os.Getenv("MAILER_SENDGRID_API_KEY"),
	}, nil
}
