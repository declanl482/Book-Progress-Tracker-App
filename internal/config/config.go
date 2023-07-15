package config

import (
	"fmt"
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

type TConfig struct {
	TestDatabaseHostname     string
	TestDatabasePort         string
	TestDatabaseName         string
	TestDatabaseUsername     string
	TestDatabasePassword     string
	TestDatabaseTimezone     string
	TestAccessTokenSecretKey string
}

var AppConfig Config
var TestConfig TConfig

func LoadApplicationConfigurationVariables() error {
	err := godotenv.Load("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/.env")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error in function LoadApplicationConfigurationVariables() ; failed to load .env file: %v", err)
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
	return nil
}

func LoadTestingConfigurationVariables() error {
	err := godotenv.Load("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/.env.test")
	if err != nil {
		return fmt.Errorf("error in function LoadTestingConfigurationVariables ; failed to load .env file: \n%v", err)
	}

	TestConfig = TConfig{
		TestDatabaseHostname:     os.Getenv("TEST_DATABASE_HOSTNAME"),
		TestDatabasePort:         os.Getenv("TEST_DATABASE_PORT"),
		TestDatabaseName:         os.Getenv("TEST_DATABASE_NAME"),
		TestDatabaseUsername:     os.Getenv("TEST_DATABASE_USERNAME"),
		TestDatabasePassword:     os.Getenv("TEST_DATABASE_PASSWORD"),
		TestDatabaseTimezone:     os.Getenv("TEST_DATABASE_TIMEZONE"),
		TestAccessTokenSecretKey: os.Getenv("TEST_ACCESS_TOKEN_SECRET_KEY"),
	}
	return nil
}
