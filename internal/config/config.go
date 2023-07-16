package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	DatabaseHostname     string
	DatabasePort         string
	DatabaseName         string
	DatabaseUsername     string
	DatabasePassword     string
	DatabaseTimezone     string
	AccessTokenSecretKey string
}

var AppConfig ApplicationConfig

func LoadApplicationConfigurationVariables() error {
	err := godotenv.Load("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/.env")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error in function LoadApplicationConfigurationVariables() ; failed to load .env file: %v", err)
	}

	AppConfig = ApplicationConfig{
		DatabaseHostname:     os.Getenv("DATABASE_HOSTNAME"),
		DatabasePort:         os.Getenv("DATABASE_PORT"),
		DatabaseName:         os.Getenv("DATABASE_NAME"),
		DatabaseUsername:     os.Getenv("DATABASE_USERNAME"),
		DatabasePassword:     os.Getenv("DATABASE_PASSWORD"),
		DatabaseTimezone:     os.Getenv("DATABASE_TIMEZONE"),
		AccessTokenSecretKey: os.Getenv("ACCESS_TOKEN_SECRET_KEY"),
	}
	return nil
}
