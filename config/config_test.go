package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TestConfiguration struct {
	TestHostname        string
	TestPort            string
	TestName            string
	TestUsername        string
	TestPassword        string
	TestTimezone        string
	TestAccessSecretKey string
}

var TestConfig TestConfiguration

func LoadTestConfigurationVariables() error {
	err := godotenv.Load("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/.env.test")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error in function LoadConfigurationVariables() ; failed to load .env file: %v", err)
	}

	TestConfig = TestConfiguration{
		TestHostname:        os.Getenv("TEST_HOSTNAME"),
		TestPort:            os.Getenv("TEST_PORT"),
		TestName:            os.Getenv("TEST_NAME"),
		TestUsername:        os.Getenv("TEST_USERNAME"),
		TestPassword:        os.Getenv("TEST_PASSWORD"),
		TestTimezone:        os.Getenv("TEST_TIMEZONE"),
		TestAccessSecretKey: os.Getenv("TEST_ACCESS_SECRET_KEY"),
	}
	return nil
}

func TestConfigurationVariables(t *testing.T) {
	err := LoadTestConfigurationVariables()
	assert.NoError(t, err, "expected no error loading test configuration variables, got: %v.", err)

	err = LoadConfigurationVariables()
	assert.NoError(t, err, "expected no error loading configuration variables, got: %v.", err)

	assert.Equal(t, TestConfig.TestHostname, Config.TestDatabaseHostname)
	assert.Equal(t, TestConfig.TestPort, Config.TestDatabasePort)
	assert.Equal(t, TestConfig.TestName, Config.TestDatabaseName)
	assert.Equal(t, TestConfig.TestUsername, Config.TestDatabaseUsername)
	assert.Equal(t, TestConfig.TestPassword, Config.TestDatabasePassword)
	assert.Equal(t, TestConfig.TestTimezone, Config.TestDatabaseTimezone)
	assert.Equal(t, TestConfig.TestAccessSecretKey, Config.TestAccessTokenSecretKey)

	t.Logf("Successfully loaded and verified configuration variables for testing database.")
}
