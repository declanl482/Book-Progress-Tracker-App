package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/config"
	"fmt"

	_ "github.com/lib/pq"
)

var testingDatabase *gorm.DB
var testDBError error

func ConfigureTestingDatabaseDSN() string {
	// Configure the DSN (Data Source Name) using the testing configuration variables.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.TestConfig.TestDatabaseHostname,
		config.TestConfig.TestDatabaseUsername,
		config.TestConfig.TestDatabasePassword,
		config.TestConfig.TestDatabaseName,
		config.TestConfig.TestDatabasePort,
		config.TestConfig.TestDatabaseTimezone)
	// Return the testing DSN string.
	return dsn
}

func OpenConnectionToTestingDatabase(dsn string) (*gorm.DB, error) {
	// Open a connection to the PostgreSQL testing database using the DSN string.
	testingDatabase, testDBError = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Return an instance of the testing database retrieved using GORM.
	// If the connection is unsuccessful, testingDatabase == nil and err != nil.
	// If the connection is successful, testingDatabase != nil and err == nil.
	return testingDatabase, testDBError
}

func MigrateTablesToTestingDatabase() error {
	// Migrate desired tables to the testing database using pre-defined models from package models.
	testDBError = testingDatabase.AutoMigrate(&models.User{}, &models.Book{})
	return testDBError
}

func ConnectToTestingDatabase() (*gorm.DB, error) {
	// Configure the DSN (Data Source Name) for the PostgreSQL testing database.
	dsn := ConfigureTestingDatabaseDSN()

	// Open a connection to the PostgreSQL testing database.
	testingDatabase, testDBError = OpenConnectionToTestingDatabase(dsn)

	// Check if the connection to the testing database was successful.
	if testDBError != nil {
		fmt.Println("Error connecting to testing database:", testDBError)
		return nil, testDBError
	}

	// Automigrate desired tables to the testing database using pre-defined models from the models package.
	testDBError = MigrateTablesToTestingDatabase()

	// Check if the migration of tables to the testing database was successful.
	if testDBError != nil {
		fmt.Println("Error migrating tables to the testing database:", testDBError)
		return nil, testDBError
	}

	// The testing database connection and table migration was successful.
	// Return an instance of the testing database without errors.
	fmt.Println("Connection to the testing database was successful.")
	return testingDatabase, nil
}

func GetInstanceOfTestingDatabase() *gorm.DB {
	// Return an instance of the application database of type *gorm.DB.
	return testingDatabase
}
