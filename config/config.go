package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "go-book-tracker-app"

func loadEnv() error {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()

	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}
	return nil
}

type Configuration struct {
	DatabaseHostname         string
	DatabasePort             string
	DatabaseName             string
	DatabaseUsername         string
	DatabasePassword         string
	DatabaseTimezone         string
	AccessTokenSecretKey     string
	TestDatabaseHostname     string
	TestDatabasePort         string
	TestDatabaseName         string
	TestDatabaseUsername     string
	TestDatabasePassword     string
	TestDatabaseTimezone     string
	TestAccessTokenSecretKey string
}

var Config Configuration

func LoadConfigurationVariables() (*Configuration, error) {
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	Config = Configuration{
		DatabaseHostname:         os.Getenv("DATABASE_HOSTNAME"),
		DatabasePort:             os.Getenv("DATABASE_PORT"),
		DatabaseName:             os.Getenv("DATABASE_NAME"),
		DatabaseUsername:         os.Getenv("DATABASE_USERNAME"),
		DatabasePassword:         os.Getenv("DATABASE_PASSWORD"),
		DatabaseTimezone:         os.Getenv("DATABASE_TIMEZONE"),
		AccessTokenSecretKey:     os.Getenv("ACCESS_TOKEN_SECRET_KEY"),
		TestDatabaseHostname:     os.Getenv("TEST_DATABASE_HOSTNAME"),
		TestDatabasePort:         os.Getenv("TEST_DATABASE_PORT"),
		TestDatabaseName:         os.Getenv("TEST_DATABASE_NAME"),
		TestDatabaseUsername:     os.Getenv("TEST_DATABASE_USERNAME"),
		TestDatabasePassword:     os.Getenv("TEST_DATABASE_PASSWORD"),
		TestDatabaseTimezone:     os.Getenv("TEST_DATABASE_TIMEZONE"),
		TestAccessTokenSecretKey: os.Getenv("TEST_ACCESS_TOKEN_SECRET_KEY"),
	}
	return &Config, nil
}
