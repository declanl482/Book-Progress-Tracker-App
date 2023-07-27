package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/declanl482/go-book-tracker-app/types"
	"github.com/stretchr/testify/assert"
)

func TestStorageFunctions(t *testing.T) {
	// Configure the testing database variables.
	testConfig, err := config.LoadTestConfigurationVariables()
	assert.NoError(t, err, "expected no error when loading test config variables, got: %v.", err)

	testHostname := testConfig.TestDatabaseHostname
	testUsername := testConfig.TestDatabaseUsername
	testPassword := testConfig.TestDatabasePassword
	testDBName := testConfig.TestDatabaseName
	testPort := testConfig.TestDatabasePort
	testTimezone := testConfig.TestDatabaseTimezone

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		testHostname,
		testUsername,
		testPassword,
		testDBName,
		testPort,
		testTimezone)

	db, err := OpenDatabaseConnection(dsn)
	assert.NoError(t, err, "\nFailed to open a connection to the testing database.\n")

	store := PostgresStorage{db: db}

	postgresInstance, err := GetDatabaseInstance()
	assert.Equal(t, store, *postgresInstance)

	// ***Note: At the time of calling our interface functions in the handlers,
	// requests are already validated against the types using c.ShouldBindJSON().
	// We can expect that inputs of type *types.User are valid (no missing username/email/pass or invalid email).

	userInitial := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "foo@bar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedUsername := &types.User{
		ID:        1,
		Username:  "bob",
		Email:     "foo@bar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedEmail := &types.User{
		ID:        1,
		Username:  "bob",
		Email:     "fuzz@buzz.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedPassword := &types.User{
		ID:        1,
		Username:  "bob",
		Email:     "fuzz@buzz.com",
		Password:  "password",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userFinal := &types.User{
		ID:        1,
		Username:  "tom",
		Email:     "beep@boop.com",
		Password:  "password123",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("TestPostgresFunctions", func(t *testing.T) {
		// (1) Tests the configuration of the DSN string.
		t.Run("TestOpenDatabaseConnection", func(t *testing.T) {
			_, err := OpenDatabaseConnection(dsn)

			assert.NoError(t, err, "\nFailed to open a connection to the testing database.\n")

			// TEST PASSED.
			t.Logf("Successfully opened a connection to the testing database.")
		})

		// (2) Tests the migration of tables to the database.
		t.Run("TestMigrateTablesToDatabase", func(t *testing.T) {
			err = MigrateTablesToDatabase(db)
			assert.NoError(t, err, "\nFailed to automigrate tables to testing database")

			// TEST PASSED.
			t.Logf("Successfully migrated tables to the testing database.")
		})

		// (3) Tests the construction of a new PostgresStorage.
		t.Run("TestNewPostgresStorage", func(t *testing.T) {
			_, err := NewPostgresStorage(testHostname, testUsername, testPassword, testDBName, testPort, testTimezone)
			assert.NoError(t, err, "\nFailed to create a new PostgresStorage.\n")

			// TEST PASSED.
			t.Logf("Successfully created a new PostgresStorage.")
		})
	})

	t.Run("UserFunctions", func(t *testing.T) {

		// Create a user (VALID).
		createdUser, err := store.CreateUser(userInitial)

		assert.NoError(t, err, "expected no error creating user, got: %v.", err)
		assert.Equal(t, userInitial, createdUser)

		// Get a user (User not found.)
		fetchedUser, err := store.GetUser(0)

		// If fetchedUser == nil, then nil error for function GetUser means user is not found.
		// For any other error, function returns result.Error.
		assert.NoError(t, err, "expected user not found error, got: user not found error.")
		assert.Equal(t, (*types.User)(nil), fetchedUser) // the user is nil

		// Get a user (success)
		fetchedUser, err = store.GetUser(1)
		assert.NoError(t, err, "expected no error fetching valid user, got: %v.", err)
		assert.Equal(t, userInitial, fetchedUser)

		// Update a user (Username)
		newUsername, err := store.UpdateUser(updatedUsername)
		assert.NoError(t, err, "expected no error updating user's username, got: %v.", err)
		assert.Equal(t, updatedUsername, newUsername)

		// Update a user (Email)
		newEmail, err := store.UpdateUser(updatedEmail)
		assert.NoError(t, err, "expected no error updating user's email, got: %v.", err)
		assert.Equal(t, updatedEmail, newEmail)

		// Update a user (Password) // TODO: Implement logic for rehashing password.
		newPassword, err := store.UpdateUser(updatedPassword)
		assert.NoError(t, err, "expected no error updating user's password, got: %v.", err)
		assert.Equal(t, updatedPassword, newPassword)

		// Update a user (All fields)
		finalUpdate, err := store.UpdateUser(userFinal)
		assert.NoError(t, err, "expected no error updating the user, got: %v.", err)
		assert.Equal(t, userFinal, finalUpdate)

		// Delete a user (Valid).
		err = store.DeleteUser(userFinal)
		assert.NoError(t, err, "expected no error deleting valid user, got: %v.", err)

	})

	// ***Note: At the time of calling our interface functions in the handlers,
	// requests are already validated against the types using c.ShouldBindJSON().
	// We can expect that inputs of type *types.Book are valid (no missing required fields, or invalid inputs).

	bookOwner1 := &types.User{
		ID:        1,
		Username:  "tom",
		Email:     "beep@boop.com",
		Password:  "password123",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bookOwner2 := &types.User{
		ID:        2,
		Username:  "bob",
		Email:     "fuzz@buzz.com",
		Password:  "password",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bookInitial1 := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 100,
		PagesRead:  50,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	bookInitial2 := &types.Book{
		ID:         2,
		Title:      "Book 2",
		Author:     "Bob",
		PagesCount: 300,
		PagesRead:  150,
		OwnerID:    2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	bookUpdated1 := &types.Book{
		ID:         1,
		Title:      "Book 1 New Title",
		Edition:    2,
		Author:     "Tom 2.0",
		PagesCount: 150,
		PagesRead:  75,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	bookUpdated2 := &types.Book{
		ID:         2,
		Title:      "Book 2 New Title",
		Edition:    4,
		Author:     "Bob 2.0",
		PagesCount: 600,
		PagesRead:  300,
		OwnerID:    2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	t.Run("BookFunctions", func(t *testing.T) {

		// Create book owner 1.
		createdOwner1, err := store.CreateUser(bookOwner1)

		assert.NoError(t, err, "expected no error creating user 1, got: %v.", err)
		assert.Equal(t, bookOwner1, createdOwner1)

		// Create book owner 2.
		createdOwner2, err := store.CreateUser(bookOwner2)

		assert.NoError(t, err, "expected no error creating user 2, got: %v.", err)
		assert.Equal(t, bookOwner2, createdOwner2)

		// Create a book (valid, with edition).
		createdBook1, err := store.CreateBook(bookInitial1)
		assert.NoError(t, err, "expected no error creating valid book 1, got: %v.", err)
		assert.Equal(t, bookInitial1, createdBook1)

		// Create another book (valid, without edition).
		createdBook2, err := store.CreateBook(bookInitial2)
		assert.NoError(t, err, "expected no error creating valid book 2, got: %v.", err)
		assert.Equal(t, bookInitial2, createdBook2)

		// Get all books for owner 1 (success).
		books1, err := store.GetBooks(1)
		assert.NoError(t, err, "expected no error fetching user 1's books, go: %v.", err)
		assert.Equal(t, []types.Book{*bookInitial1}, *books1)

		// Get all books for owner 2 (success).
		books2, err := store.GetBooks(2)
		assert.NoError(t, err, "expected no error fetching user 1's books, go: %v.", err)
		assert.Equal(t, []types.Book{*bookInitial2}, *books2)

		// Get a book (not found).
		bookNotFound, err := store.GetBook(0)
		assert.Empty(t, bookNotFound)
		assert.NoError(t, err, "expected book not found error, got: book not found error")
		//assert.Equal(t, (*types.Book)(nil), bookNotFound)

		// Get book 1 (success).
		fetchedBook1, err := store.GetBook(1)
		assert.NoError(t, err, "expected no error fetching valid book 1, got: %v.", err)
		assert.Equal(t, bookInitial1, fetchedBook1)

		// Get book 2 (success).
		fetchedBook2, err := store.GetBook(2)
		assert.NoError(t, err, "expected no error fetching valid book 2, got: %v.", err)
		assert.Equal(t, bookInitial2, fetchedBook2)

		// Update book 1 (all fields: title, edition, author, pagescount, pagesread)
		newBook1, err := store.UpdateBook(bookUpdated1)
		assert.NoError(t, err, "expected no error updating book 1, got: %v.", err)
		assert.Equal(t, bookUpdated1, newBook1)

		// Update book 2 (all fields).
		newBook2, err := store.UpdateBook(bookUpdated2)
		assert.NoError(t, err, "expected no error updating book 2, got: %v.", err)
		assert.Equal(t, bookUpdated2, newBook2)

		// Delete book 1 (success).
		err = store.DeleteBook(bookInitial1)
		assert.NoError(t, err, "expected no error deleting book 1, got: %v.", err)

		// Delete book 2 (success).
		err = store.DeleteBook(bookInitial2)
		assert.NoError(t, err, "expected no error deleting book 2, got: %v.", err)

		// Delete book owner 1.
		err = store.DeleteUser(bookOwner1)
		assert.NoError(t, err, "expected no error deleting valid user 1, got: %v.", err)

		// Delete book owner 2.
		err = store.DeleteUser(bookOwner2)
		assert.NoError(t, err, "expected no error deleting valid user 2, got: %v.", err)
	})
}
