package database_test

import (
	"example/go-book-tracker-app/internal/config"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDatabase(t *testing.T) {

	// (1) Tests the opening of a connection to the specified PostgreSQL testing database.
	t.Run("TestOpenDatabaseConnection", func(t *testing.T) {
		// Load the application-level configuration variables.
		err := config.LoadTestingConfigurationVariables()

		// TEST MAY FAIL HERE.
		if assert.NoError(t, err, "\nFailed to load testing-level configuration variables.\n") {
			// Get the testing database connection string.
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
				config.TestConfig.TestDatabaseHostname,
				config.TestConfig.TestDatabaseUsername,
				config.TestConfig.TestDatabasePassword,
				config.TestConfig.TestDatabaseName,
				config.TestConfig.TestDatabasePort,
				config.TestConfig.TestDatabaseTimezone)

			// Open a connection to the testing database.

			_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				t.Fatalf("Failed to open a connection to the database: %v", err)
			}

			// TEST PASSED.
			t.Logf("Successfully opened a connection to the test database!")
		}

	})

	// // (2) Tests the retrieval of the PostgreSQL testing database instance of type *gorm.DB.
	// t.Run("TestGetInstanceOfDatabase", func(t *testing.T) {

	// 	// Connect to the testing database.
	// 	testingDatabase, err := database.ConnectToTestingDatabase()
	// 	if assert.NoError(t, err, "")

	// 	//testingDatabase := database.GetInstanceOfTestingDatabase()
	// })

	// (2) Ping the database.

	// t.Run("TestPingDatabase", func(t* testing.T)) {

	// }

	// (3) Automigrate 'users' and 'books' table into the database.

}
