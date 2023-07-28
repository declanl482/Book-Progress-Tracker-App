package storage

import (
	"errors"
	"fmt"

	"github.com/declanl482/go-book-tracker-app/backend/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

var appDB *gorm.DB
var appDBErr error

func GetDatabaseInstance() (*PostgresStorage, error) {
	if appDB == nil {
		return nil, errors.New("database connection error: database instance is nil")
	}
	return &PostgresStorage{db: appDB}, nil
}
func OpenDatabaseConnection(dsn string) (*gorm.DB, error) {
	// Open a connection to the database.
	appDB, appDBErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return appDB, appDBErr
}

func MigrateTablesToDatabase(db *gorm.DB) error {
	// Migrate desired tables to database using pre-defined types.
	appDBErr = db.AutoMigrate(&types.User{}, &types.Book{})
	return appDBErr
}

func NewPostgresStorage(hostname string, username string, password string, dbname string, port string, timezone string) (*PostgresStorage, error) {
	// Configure the postgres dsn.

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		hostname,
		username,
		password,
		dbname,
		port,
		timezone)

	// Open a connection to the database.
	appDB, appDBErr = OpenDatabaseConnection(dsn)
	if appDBErr != nil {
		fmt.Println("Error connecting to application database:", appDBErr)
		return nil, appDBErr
	}

	// Migrate tables to the database.
	appDBErr = MigrateTablesToDatabase(appDB)
	if appDBErr != nil {
		fmt.Println("Error migrating tables to the application database:", appDBErr)
		return nil, appDBErr
	}

	// The application database connection and table migration was successful.
	// Return an instance of the application database without errors.
	fmt.Println("Connection to the application database was successful.")
	// Return an instance of PostgresStorage.
	return &PostgresStorage{db: appDB}, nil
}

func (s *PostgresStorage) IsEmailTaken(email string) (bool, error) {
	// Check if the new user's email is already in use.
	existingUser, err := s.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return existingUser != nil, nil
}

func (s *PostgresStorage) CreateUser(user *types.User) (*types.User, error) {

	// Create the new user in the database.
	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the newly created user.
	return user, nil
}

func (s *PostgresStorage) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Record not found.
		}
		return nil, result.Error // Return the error for any other error.
	}
	return &user, nil
}

func (s *PostgresStorage) GetUser(id int) (*types.User, error) {
	// Declare a variable for the fetched user.
	var fetchedUser types.User

	// Retrieve the user from the database.
	// Preload the books associated with the user.
	result := s.db.Preload("Books").First(&fetchedUser, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User with the given id was not found.
		}
		return nil, result.Error // Database error when fetching.
	}

	// Return the fetched user.
	return &fetchedUser, nil
}

func (s *PostgresStorage) UpdateUser(user *types.User) (*types.User, error) {
	result := s.db.Model(&user).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *PostgresStorage) DeleteUser(user *types.User) error {
	result := s.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *PostgresStorage) CreateBook(book *types.Book) (*types.Book, error) {
	result := s.db.Create(book)
	if result.Error != nil {
		return nil, result.Error
	}
	return book, nil
}

func (s *PostgresStorage) GetBooks(id int) (*[]types.Book, error) {
	var books []types.Book

	result := s.db.Where("owner_id = ?", id).Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return &books, nil
}

func (s *PostgresStorage) GetBook(id int) (*types.Book, error) {
	var book types.Book

	result := s.db.First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &book, nil
}

func (s *PostgresStorage) UpdateBook(book *types.Book) (*types.Book, error) {
	result := s.db.Model(&book).Updates(book)
	if result.Error != nil {
		return nil, result.Error
	}
	return book, nil
}

func (s *PostgresStorage) DeleteBook(book *types.Book) error {
	result := s.db.Delete(&book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
