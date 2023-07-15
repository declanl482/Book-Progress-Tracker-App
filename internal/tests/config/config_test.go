package config_test

import (
	"example/go-book-tracker-app/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationVariables(t *testing.T) {

	// (1a) Tests the loading of configuration variables for the application environment.
	t.Run("TestLoadApplicationConfigurationVariables", func(t *testing.T) {
		// Load the application-level configuration variables.
		err := config.LoadApplicationConfigurationVariables()

		// TEST FAILED
		if assert.NoError(t, err, "\nFailed to load application-level configuration variables.\n") {
			// TEST PASSED
			t.Logf("Successfully loaded application-level configuration variables.")
		}
	})

	// (1b) Verifies that the loaded application-level configuration variables match the expected test-specific values.
	t.Run("TestMatchingApplicationConfigurationValues", func(t *testing.T) {
		// Set up test-specific environment variables in the application environment, using dummy values.
		os.Setenv("DATABASE_HOSTNAME", "app_database_hostname")
		os.Setenv("DATABASE_PORT", "app_database_port")
		os.Setenv("DATABASE_NAME", "app_database_name")
		os.Setenv("DATABASE_USERNAME", "app_database_username")
		os.Setenv("DATABASE_PASSWORD", "app_database_password")
		os.Setenv("ACCESS_TOKEN_SECRET_KEY", "app_access_token_secret_key")

		// Load the application-level configuration variables.
		err := config.LoadApplicationConfigurationVariables()

		// Assert that there is no error when loading application-level configuration variables.
		// TEST MAY FAIL HERE.
		if assert.NoError(t, err, "\nFailed to load application-level configuration variables.\n") {
			// Assert that the loaded application-level configuration variables match the expected test-specific values.
			assert.Equal(t, "app_database_hostname", config.AppConfig.DatabaseHostname, "Database hostname should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.DatabaseHostname, "app_database_hostname")
			assert.Equal(t, "app_database_port", config.AppConfig.DatabasePort, "Database port should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.DatabasePort, "app_database_port")
			assert.Equal(t, "app_database_name", config.AppConfig.DatabaseName, "Database name should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.DatabaseName, "app_database_name")
			assert.Equal(t, "app_database_username", config.AppConfig.DatabaseUsername, "Database username should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.DatabaseUsername, "app_database_username")
			assert.Equal(t, "app_database_password", config.AppConfig.DatabasePassword, "Database password should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.DatabasePassword, "app_database_password")
			assert.Equal(t, "app_access_token_secret_key", config.AppConfig.AccessTokenSecretKey, "Access-token secret key should match dummy application value. (Got: %v; Expected: %v)", config.AppConfig.AccessTokenSecretKey, "app_access_token_secret_key")

			// TEST PASSED.
			t.Logf("Successfully matched the application-level configuration values with expected test-specific values.")
		}
	})

	// (2a) Tests the loading of configuration variables for the testing environment.
	t.Run("TestLoadTestingConfigurationVariables", func(t *testing.T) {
		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// TEST PASSED.
			t.Logf("Successfully loaded testing-level configuration variables.")
		}
	})

	// (2b) Verifies that the loaded testing-level configuration variables match the expected test-specific values.
	t.Run("TestMatchingTestingConfigurationValues", func(t *testing.T) {
		// Set up test-specific environment variables in the testing environment, using dummy values.
		os.Setenv("TEST_DATABASE_HOSTNAME", "test_database_hostname")
		os.Setenv("TEST_DATABASE_PORT", "test_database_port")
		os.Setenv("TEST_DATABASE_NAME", "test_database_name")
		os.Setenv("TEST_DATABASE_USERNAME", "test_database_username")
		os.Setenv("TEST_DATABASE_PASSWORD", "test_database_password")
		os.Setenv("TEST_ACCESS_TOKEN_SECRET_KEY", "test_access_token_secret_key")

		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// Assert that there is no error when loading testing-level configuration variables.
		// TEST MAY FAIL HERE.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// Assert that the loaded testing-level configuration variables match the expected test-specific values.
			assert.Equal(t, "test_database_hostname", config.TestConfig.TestDatabaseHostname, "Testing database hostname should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseHostname, "test_database_hostname")
			assert.Equal(t, "test_database_port", config.TestConfig.TestDatabasePort, "Testing database port should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePort, "test_database_port")
			assert.Equal(t, "test_database_name", config.TestConfig.TestDatabaseName, "Testing database name should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseName, "test_database_name")
			assert.Equal(t, "test_database_username", config.TestConfig.TestDatabaseUsername, "Testing database username should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabaseUsername, "test_database_username")
			assert.Equal(t, "test_database_password", config.TestConfig.TestDatabasePassword, "Testing database password should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestDatabasePassword, "test_database_password")
			assert.Equal(t, "test_access_token_secret_key", config.TestConfig.TestAccessTokenSecretKey, "Testing access-token secret key should match dummy test value. (Got: %v; Expected: %v)", config.TestConfig.TestAccessTokenSecretKey, "test_access_token_secret_key")

			// TEST PASSED.
			t.Logf("Successfully matched the testing-level configuration values with expected test-specific values.")
		}
	})

}
