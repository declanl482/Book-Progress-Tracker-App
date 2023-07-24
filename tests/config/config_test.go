package config_test

import (
	"example/go-book-tracker-app/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationVariables(t *testing.T) {
	// (1a) Tests the loading of configuration variables for the testing environment.
	t.Run("LoadTestingConfigurationVariables", func(t *testing.T) {
		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n")

		// TEST PASSED.
		t.Logf("Successfully loaded testing-level configuration variables.")
	})

	// (1b) Verifies that the loaded testing-level configuration variables match the expected test-specific values.
	t.Run("ExpectedTestingConfigurationValues", func(t *testing.T) {

		// Load the testing-level configuration variables from the JSON file.
		jsonConfig, err := config.LoadJSONTestingConfigurationVariables("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/testing_config.json")

		// TEST MAY FAIL HERE.
		assert.NoError(t, err, "Failed to load JSON configuration variables for the testing database.")

		// Load the testing-level configuration variables.
		err = config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n")

		// TEST MAY FAIL HERE.
		// Assert that the loaded testing-level configuration variables match the expected test-specific values.
		assert.Equal(t, jsonConfig.TestDatabaseHostname, config.TestConfig.TestDatabaseHostname, "Testing database hostname should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseHostname, jsonConfig.TestDatabaseHostname)
		assert.Equal(t, jsonConfig.TestDatabasePort, config.TestConfig.TestDatabasePort, "Testing database port should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePort, jsonConfig.TestDatabasePort)
		assert.Equal(t, jsonConfig.TestDatabaseName, config.TestConfig.TestDatabaseName, "Testing database name should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseName, jsonConfig.TestDatabaseName)
		assert.Equal(t, jsonConfig.TestDatabaseUsername, config.TestConfig.TestDatabaseUsername, "Testing database username should match expected valuee. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseUsername, jsonConfig.TestDatabaseUsername)
		assert.Equal(t, jsonConfig.TestDatabasePassword, config.TestConfig.TestDatabasePassword, "Testing database password should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePassword, jsonConfig.TestDatabasePassword)
		assert.Equal(t, jsonConfig.TestDatabaseTimezone, config.TestConfig.TestDatabaseTimezone, "Testing database timezone should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseTimezone, jsonConfig.TestDatabaseTimezone)
		assert.Equal(t, jsonConfig.TestAccessTokenSecretKey, config.TestConfig.TestAccessTokenSecretKey, "Testing access-token secret key should match expected value. (Got: %v; Expected: %v)", config.TestConfig.TestAccessTokenSecretKey, jsonConfig.TestAccessTokenSecretKey)

		// TEST PASSED.
		t.Logf("Successfully matched the testing-level configuration values with expected test-specific values.")
	})
}
