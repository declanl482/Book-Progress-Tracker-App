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
		// Set up test-specific environment variables in the testing environment, using dummy values.
		os.Setenv("TEST_DATABASE_HOSTNAME", "test_database_hostname")
		os.Setenv("TEST_DATABASE_PORT", "test_database_port")
		os.Setenv("TEST_DATABASE_NAME", "test_database_name")
		os.Setenv("TEST_DATABASE_USERNAME", "test_database_username")
		os.Setenv("TEST_DATABASE_PASSWORD", "test_database_password")
		os.Setenv("TEST_DATABASE_TIMEZONE", "test_database_timezone")
		os.Setenv("TEST_ACCESS_TOKEN_SECRET_KEY", "test_access_token_secret_key")

		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// TEST MAY FAIL HERE.
			// Assert that the loaded testing-level configuration variables match the expected test-specific values.
			assert.Equal(t, "test_database_hostname", config.TestConfig.TestDatabaseHostname, "Testing database hostname should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseHostname, "test_database_hostname")
			assert.Equal(t, "test_database_port", config.TestConfig.TestDatabasePort, "Testing database port should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePort, "test_database_port")
			assert.Equal(t, "test_database_name", config.TestConfig.TestDatabaseName, "Testing database name should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseName, "test_database_name")
			assert.Equal(t, "test_database_username", config.TestConfig.TestDatabaseUsername, "Testing database username should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseUsername, "test_database_username")
			assert.Equal(t, "test_database_password", config.TestConfig.TestDatabasePassword, "Testing database password should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePassword, "test_database_password")
			assert.Equal(t, "test_database_timezone", config.TestConfig.TestDatabaseTimezone, "Testing database timezone should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseTimezone, "test_database_timezone")
			assert.Equal(t, "test_access_token_secret_key", config.TestConfig.TestAccessTokenSecretKey, "Testing access-token secret key should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestAccessTokenSecretKey, "test_access_token_secret_key")

			// TEST PASSED.
			t.Logf("Successfully matched the testing-level configuration values with expected test-specific values.")
		}
	})

}
