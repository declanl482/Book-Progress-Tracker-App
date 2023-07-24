package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example/go-book-tracker-app/api/models"
	"example/go-book-tracker-app/config"
	"fmt"

	_ "github.com/lib/pq"
)

var applicationDatabase *gorm.DB
var applicationDBError error

func ConfigureApplicationDatabaseDSN() string {
	// Configure the DSN (Data Source Name) using the application configuration variables.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.AppConfig.DatabaseHostname,
		config.AppConfig.DatabaseUsername,
		config.AppConfig.DatabasePassword,
		config.AppConfig.DatabaseName,
		config.AppConfig.DatabasePort,
		config.AppConfig.DatabaseTimezone)
	// Return the application DSN string.
	return dsn
}

func OpenConnectionToApplicationDatabase(dsn string) (*gorm.DB, error) {
	// Open a connection to the PostgreSQL application database using the DSN string.
	applicationDatabase, applicationDBError = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Return an instance of the application database retrieved using GORM.
	// If the connection is unsuccessful, applicationDatabase == nil and err != nil.
	// If the connection is successful, applicationDatabase != nil and err == nil.
	return applicationDatabase, applicationDBError
}

func MigrateTablesToApplicationDatabase() error {
	// Migrate desired tables to the application database using pre-defined models from package models.
	applicationDBError = applicationDatabase.AutoMigrate(&models.User{}, &models.Book{})
	return applicationDBError
}

func ConnectToApplicationDatabase() (*gorm.DB, error) {
	// Configure the DSN (Data Source Name) for the PostgreSQL application database.
	dsn := ConfigureApplicationDatabaseDSN()

	// Open a connection to the PostgreSQL application database.
	applicationDatabase, applicationDBError = OpenConnectionToApplicationDatabase(dsn)

	// Check if the connection to the application database was successful.
	if applicationDBError != nil {
		fmt.Println("Error connecting to application database:", applicationDBError)
		return nil, applicationDBError
	}

	// Automigrate desired tables to the application database using pre-defined models from the models package.
	applicationDBError = MigrateTablesToApplicationDatabase()

	// Check if the migration of tables to the application database was successful.
	if applicationDBError != nil {
		fmt.Println("Error migrating tables to the application database:", applicationDBError)
		return nil, applicationDBError
	}

	// The application database connection and table migration was successful.
	// Return an instance of the application database without errors.
	fmt.Println("Connection to the application database was successful.")
	return applicationDatabase, nil
}

func GetInstanceOfApplicationDatabase() *gorm.DB {
	// Return an instance of the application database of type *gorm.DB.
	return applicationDatabase
}
