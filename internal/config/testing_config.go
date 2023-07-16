package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type TestingConfig struct {
	TestDatabaseHostname     string
	TestDatabasePort         string
	TestDatabaseName         string
	TestDatabaseUsername     string
	TestDatabasePassword     string
	TestDatabaseTimezone     string
	TestAccessTokenSecretKey string
}

var TestConfig TestingConfig

func LoadTestingConfigurationVariables() error {
	err := godotenv.Load("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/.env.test")
	if err != nil {
		return fmt.Errorf("error in function LoadTestingConfigurationVariables ; failed to load .env.test file: \n%v", err)
	}

	TestConfig = TestingConfig{
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
