package config_test

import (
	"example/go-book-tracker-app/internal/config"
	"os"
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
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// TEST PASSED.
			t.Logf("Successfully loaded testing-level configuration variables.")
		}
	})

	// (1b) Verifies that the loaded testing-level configuration variables match the expected test-specific values.
	t.Run("ExpectedTestingConfigurationValues", func(t *testing.T) {

		// Load the testing-level configuration variables from the JSON file.
		jsonConfig, err := config.LoadJSONTestingConfigurationVariables("C:/Users/13dli/go/src/github.com/declanl482/go-book-tracker-app/testing_config.json")
		assert.NoError(t, err, "Failed to load JSON configuration variables for the testing database.")

		// Set up the test-specific environment variables based on the JSON configuration
		os.Setenv("TEST_DATABASE_HOSTNAME", jsonConfig.TestDatabaseHostname)
		os.Setenv("TEST_DATABASE_PORT", jsonConfig.TestDatabasePort)
		os.Setenv("TEST_DATABASE_NAME", jsonConfig.TestDatabaseName)
		os.Setenv("TEST_DATABASE_USERNAME", jsonConfig.TestDatabaseUsername)
		os.Setenv("TEST_DATABASE_PASSWORD", jsonConfig.TestDatabasePassword)
		os.Setenv("TEST_DATABASE_TIMEZONE", jsonConfig.TestDatabaseTimezone)
		os.Setenv("TEST_ACCESS_TOKEN_SECRET_KEY", jsonConfig.TestAccessTokenSecretKey)

		// Load the testing-level configuration variables.
		err = config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// TEST MAY FAIL HERE.
			// Assert that the loaded testing-level configuration variables match the expected test-specific values.
			assert.Equal(t, jsonConfig.TestDatabaseHostname, config.TestConfig.TestDatabaseHostname, "Testing database hostname should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseHostname, jsonConfig.TestDatabaseHostname)
			assert.Equal(t, jsonConfig.TestDatabasePort, config.TestConfig.TestDatabasePort, "Testing database port should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePort, jsonConfig.TestDatabasePort)
			assert.Equal(t, jsonConfig.TestDatabaseName, config.TestConfig.TestDatabaseName, "Testing database name should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseName, jsonConfig.TestDatabaseName)
			assert.Equal(t, jsonConfig.TestDatabaseUsername, config.TestConfig.TestDatabaseUsername, "Testing database username should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseUsername, jsonConfig.TestDatabaseUsername)
			assert.Equal(t, jsonConfig.TestDatabasePassword, config.TestConfig.TestDatabasePassword, "Testing database password should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePassword, jsonConfig.TestDatabasePassword)
			assert.Equal(t, jsonConfig.TestDatabaseTimezone, config.TestConfig.TestDatabaseTimezone, "Testing database timezone should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseTimezone, jsonConfig.TestDatabaseTimezone)
			assert.Equal(t, jsonConfig.TestAccessTokenSecretKey, config.TestConfig.TestAccessTokenSecretKey, "Testing access-token secret key should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestAccessTokenSecretKey, jsonConfig.TestAccessTokenSecretKey)

			// TEST PASSED.
			t.Logf("Successfully matched the testing-level configuration values with expected test-specific values.")
		}
	})

}
