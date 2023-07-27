package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DatabaseHostname     string
	DatabasePort         string
	DatabaseName         string
	DatabaseUsername     string
	DatabasePassword     string
	DatabaseTimezone     string
	AccessTokenSecretKey string
}

var Config Configuration

func LoadConfigurationVariables() (*Configuration, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	Config = Configuration{
		DatabaseHostname:     os.Getenv("DATABASE_HOSTNAME"),
		DatabasePort:         os.Getenv("DATABASE_PORT"),
		DatabaseName:         os.Getenv("DATABASE_NAME"),
		DatabaseUsername:     os.Getenv("DATABASE_USERNAME"),
		DatabasePassword:     os.Getenv("DATABASE_PASSWORD"),
		DatabaseTimezone:     os.Getenv("DATABASE_TIMEZONE"),
		AccessTokenSecretKey: os.Getenv("ACCESS_TOKEN_SECRET_KEY"),
	}
	return &Config, nil
}

type TestConfiguration struct {
	TestDatabaseHostname     string
	TestDatabasePort         string
	TestDatabaseName         string
	TestDatabaseUsername     string
	TestDatabasePassword     string
	TestDatabaseTimezone     string
	TestAccessTokenSecretKey string
}

var TestConfig TestConfiguration

func LoadTestConfigurationVariables() (*TestConfiguration, error) {

	err := godotenv.Load("../.env.test")
	if err != nil {
		return nil, err
	}

	TestConfig = TestConfiguration{
		TestDatabaseHostname:     os.Getenv("TEST_DATABASE_HOSTNAME"),
		TestDatabasePort:         os.Getenv("TEST_DATABASE_PORT"),
		TestDatabaseName:         os.Getenv("TEST_DATABASE_NAME"),
		TestDatabaseUsername:     os.Getenv("TEST_DATABASE_USERNAME"),
		TestDatabasePassword:     os.Getenv("TEST_DATABASE_PASSWORD"),
		TestDatabaseTimezone:     os.Getenv("TEST_DATABASE_TIMEZONE"),
		TestAccessTokenSecretKey: os.Getenv("TEST_ACCESS_TOKEN_SECRET_KEY"),
	}
	return &TestConfig, nil
}
