package database_test

// import (
// 	"example/go-book-tracker-app/internal/config"
// 	"fmt"
// 	"testing"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func TestDatabase(t *testing.T) {

// 	// (1) Open a connection to the database.
// 	t.Run("OpenDatabaseConnection", func(t *testing.T) {

// 		// Load the configuration
// 		config.Load()

// 		// Get the database connection string
// 		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
// 			config.TestConfig.TestDatabaseHostname,
// 			config.TestConfig.TestDatabaseUsername,
// 			config.TestConfig.TestDatabasePassword,
// 			config.TestConfig.TestDatabaseName,
// 			config.TestConfig.TestDatabasePort,
// 			config.TestConfig.TestDatabaseTimezone)

// 		// Open a connection to the database

// 		_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 		if err != nil {
// 			t.Fatalf("Failed to open a connection to the database: %v", err)
// 		}

// 		// Test passed
// 		t.Logf("Successfully opened a connection to the test database!")

// 	})

// 	// (2) Ping the database.

// 	// (3) Automigrate 'users' and 'books' table into the database.

// 	// (4) Get an instance of the database of type *gorm.DB.

// }
