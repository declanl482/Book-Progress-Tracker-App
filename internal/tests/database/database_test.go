package database_test

import (
	"example/go-book-tracker-app/internal/config"
	"example/go-book-tracker-app/internal/database"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	// (1) Tests the configuration of the testing database DSN.
	t.Run("ConfigureTestingDatabaseDSN", func(t *testing.T) {
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
			// Configure the DSN for the testing database.
			actualDSN := database.ConfigureTestingDatabaseDSN()

			// Hard code the expected DSN for the testing database.
			expectedDSN := "host=test_database_hostname user=test_database_username password=test_database_password dbname=test_database_name port=test_database_port sslmode=disable TimeZone=test_database_timezone"

			// Assert that the actual DSN matches the expected DSN for the testing database.
			assert.Equal(t, expectedDSN, actualDSN, "Configured DSN should match expected DSN. (Got: %v; Expected: %v)", actualDSN, expectedDSN)

			// TEST PASSED.
			t.Logf("Successfully configured the testing database DSN.")
		}
	})

	// (2) Tests the opening of a connection to the specified PostgreSQL testing database.
	t.Run("OpenConnectionToTestingDatabase", func(t *testing.T) {
		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// Configure the DSN for the testing database.
			dsn := database.ConfigureTestingDatabaseDSN()

			// Open a connection to the testing database.
			_, err = database.OpenConnectionToTestingDatabase(dsn)

			// TEST MAY FAIL HERE.
			// Assert that there is no error when opening a connection to the testing database.
			if assert.NoError(t, err, "\nFailed to open a connection to the testing database.\n") {
				// TEST PASSED.
				t.Logf("Successfully opened a connection to the testing database.")
			}
		}
	})

	// (3) Tests the automatic migration of tables to the testing database.
	t.Run("MigrateTablesToTestingDatabase", func(t *testing.T) {

		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// Configure the DSN for the testing database.
			dsn := database.ConfigureTestingDatabaseDSN()

			// Open a connection to the testing database.
			_, err = database.OpenConnectionToTestingDatabase(dsn)

			// TEST MAY FAIL HERE.
			// Assert that there is no error when opening a connection to the testing database.
			if assert.NoError(t, err, "\nFailed to open a connection to the testing database.\n") {
				// Run the table migration
				err := database.MigrateTablesToTestingDatabase()

				// TEST MAY FAIL HERE.
				// Assert that there is no error when migrating tables to the testing database.
				if assert.NoError(t, err, "\nFailed to migrate tables to the testing database.\n") {
					// TEST PASSED.
					t.Logf("Successfully migrated tables to the testing database.")
				}
			}
		}
	})

	// (4) Tests the DSN configuration, database connection, and table migration functions for the testing database,
	//	   encapsulated into a single function.
	t.Run("ConnectToTestingDatabase", func(t *testing.T) {
		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			_, err = database.ConnectToTestingDatabase()

			// TEST MAY FAIL HERE.
			// Assert that there is no error when initializing the testing database.
			if assert.NoError(t, err, "\nFailed to initialize the testing database.\n") {
				// TEST PASSED.
				t.Logf("Successfully initialized the testing database.")
			}
		}
	})

	// (5) Tests the retrieval of an instance of the testing database.
	t.Run("GetInstanceOfTestingDatabase", func(t *testing.T) {
		// Load the testing-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		// Assert that there is no error when loading testing-level configuration variables.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			testingDatabase, err := database.ConnectToTestingDatabase()

			// TEST MAY FAIL HERE.
			// Assert that there is no error when initializing the testing database.
			if assert.NoError(t, err, "\nFailed to initialize the testing database.\n") {
				// TEST PASSED.
				testingDatabaseInstance := database.GetInstanceOfTestingDatabase()

				// TEST MAY FAIL HERE.
				// Assert that the retrieved instance of the testing database properly matches the testing database
				// which the program is connected to.
				assert.Equal(t, testingDatabase, testingDatabaseInstance, "Testing database instance should match the testing database which the program is connected to.")

				// TEST PASSED.
				t.Logf("Successfully retrieved an instance of the testing database.")
			}
		}
	})
}
