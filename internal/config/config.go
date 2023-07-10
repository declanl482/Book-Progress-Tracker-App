package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHostname     string
	DatabasePort         string
	DatabaseName         string
	DatabaseUsername     string
	DatabasePassword     string
	DatabaseTimezone     string
	AccessTokenSecretKey string
}

var AppConfig Config

func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	AppConfig = Config{
		DatabaseHostname:     os.Getenv("DATABASE_HOSTNAME"),
		DatabasePort:         os.Getenv("DATABASE_PORT"),
		DatabaseName:         os.Getenv("DATABASE_NAME"),
		DatabaseUsername:     os.Getenv("DATABASE_USERNAME"),
		DatabasePassword:     os.Getenv("DATABASE_PASSWORD"),
		DatabaseTimezone:     os.Getenv("DATABASE_TIMEZONE"),
		AccessTokenSecretKey: os.Getenv("ACCESS_TOKEN_SECRET_KEY"),
	}
}
