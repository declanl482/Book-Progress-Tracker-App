package config

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func loadTestEnv() error {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()

	fmt.Println(currentWorkDirectory)

	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env.test`)

	if err != nil {
		return fmt.Errorf("failed to load .env.test file: %v", err)
	}
	return nil
}

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

func LoadTestConfigurationVariables() (*TestConfiguration, error) {

	err := loadTestEnv()
	if err != nil {
		return nil, err
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
	return &TestConfig, nil
}

func TestConfigurationVariables(t *testing.T) {
	testConfig, err := LoadTestConfigurationVariables()
	assert.NoError(t, err, "expected no error loading test configuration variables, got: %v.", err)

	var config *Configuration
	config, err = LoadConfigurationVariables()
	assert.NoError(t, err, "expected no error loading configuration variables, got: %v.", err)

	assert.Equal(t, testConfig.TestHostname, config.TestDatabaseHostname)
	assert.Equal(t, testConfig.TestPort, config.TestDatabasePort)
	assert.Equal(t, testConfig.TestName, config.TestDatabaseName)
	assert.Equal(t, testConfig.TestUsername, config.TestDatabaseUsername)
	assert.Equal(t, testConfig.TestPassword, config.TestDatabasePassword)
	assert.Equal(t, testConfig.TestTimezone, config.TestDatabaseTimezone)
	assert.Equal(t, testConfig.TestAccessSecretKey, config.TestAccessTokenSecretKey)

	t.Logf("Successfully loaded and verified configuration variables for testing database.")
}
