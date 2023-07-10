package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/config"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *gorm.DB

// DATABASE_HOSTNAME=localhost
// DATABASE_PORT=5432
// DATABASE_NAME=go-book-tracker-db
// DATABASE_USERNAME=postgres
// DATABASE_PASSWORD=Q2xlYW5GcmllbmRseUNvZGVyMTIz
// DATABASE_TIMEZONE=America/New_York

func ConnectToDB() (*gorm.DB, error) {
	// dsn := "host=localhost user=postgres password=Q2xlYW5GcmllbmRseUNvZGVyMTIz dbname=go-book-tracker-db port=5432 sslmode=disable TimeZone=America/New_York"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.AppConfig.DatabaseHostname,
		config.AppConfig.DatabaseUsername,
		config.AppConfig.DatabasePassword,
		config.AppConfig.DatabaseName,
		config.AppConfig.DatabasePort,
		config.AppConfig.DatabaseTimezone)
	var err error
	for {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Database connection was not successful. Retrying...")
			time.Sleep(3 * time.Second)
			continue
		}
		// AutoMigrate models to create missing tables

		err = db.AutoMigrate(&models.User{}, &models.Book{}) // Books []Book does not migrate for Users and OwnerID does not migrate for Books
		if err != nil {
			fmt.Println("Error migrating models:", err)
			return nil, err
		}

		// var users []models.User
		// db.Preload("Books").Find(&users)

		// ALTER TABLE `books` ADD CONSTANT `books_users_fkey` FOREIGN KEY (`owner_id`) REFERENCES `users`(`id`)

		fmt.Println("Database connection was successful!")
		break
	}
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}
