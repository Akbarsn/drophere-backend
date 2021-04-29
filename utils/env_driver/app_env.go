package env_driver

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppEnvironment struct {
	Mode                     string
	Port                     int
	StorageRootDirectoryName string
	TemplatePath             string
	PasswordRecovery         PasswordRecoveryEnvironment
}

type PasswordRecoveryEnvironment struct {
	TokenDuration  int
	RecoveryWebURL string
	MailerEmail    string
	MailerName     string
}

func NewAppEnvironmentDriver() (*AppEnvironment, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	tokenDuration, err := strconv.Atoi(os.Getenv("APP_PASSWORD_RECOVERY_TOKEN_DURATION"))
	if err != nil {
		return nil, err
	}

	var appPort int
	if os.Getenv("APP_PORT") == "" {
		appPort = 8080
	} else {
		appPort, err = strconv.Atoi(os.Getenv("APP_PORT"))
		if err != nil {
			return nil, err
		}
	}

	return &AppEnvironment{
		Mode:                     os.Getenv("APP_MODE"),
		Port:                     appPort,
		StorageRootDirectoryName: os.Getenv("APP_STORAGE_ROOT_DIRECTORY_NAME"),
		TemplatePath:             os.Getenv("APP_TEMPLATE_PATH"),
		PasswordRecovery: PasswordRecoveryEnvironment{
			TokenDuration:  tokenDuration,
			RecoveryWebURL: os.Getenv("APP_PASSWORD_RECOVERY_WEB_URL"),
			MailerEmail:    os.Getenv("APP_PASSWORD_MAILER_EMAIL"),
			MailerName:     os.Getenv("APP_PASSWORD_MAILER_NAME"),
		},
	}, nil
}
