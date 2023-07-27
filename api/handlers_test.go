package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/declanl482/go-book-tracker-app/storage"
	"github.com/declanl482/go-book-tracker-app/types"
	"github.com/stretchr/testify/assert"
)

func TestHandlerFunctions(t *testing.T) {
	// create a new postgres storage (test db), and use it to create a new server
	configuration, err := config.LoadConfigurationVariables()
	assert.NoError(t, err, "expected no error when loading config variables, got: %v.", err)

	hostname := configuration.TestDatabaseHostname
	username := configuration.TestDatabaseUsername
	password := configuration.TestDatabasePassword
	name := configuration.TestDatabaseName
	port := configuration.TestDatabasePort
	timezone := configuration.TestDatabaseTimezone

	store, err := storage.NewPostgresStorage(hostname, username, password, name, port, timezone)
	assert.NoError(t, err, "expected no error when creating PostgresStorage for testing database, got: %v.", err)

	listenAddress := ":8080"
	server := NewServer(listenAddress, store) // use server to access handler functions
	// Start the server.
	go func() {
		err := server.Start()
		assert.NoError(t, err, "expected no error when starting server, routers, and middlewares, got: %v.", err)
	}()

	time.Sleep(500 * time.Millisecond)

	missingUsernameUser := &types.User{
		ID:        1,
		Username:  "",
		Email:     "foo@bar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	missingEmailUser := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	missingPasswordUser := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "foo@bar.com",
		Password:  "",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	invalidEmailUser := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "foobar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validUser1 := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "foo@bar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validUser2 := &types.User{
		ID:        2,
		Username:  "fuzz",
		Email:     "fuzz@buzz.com",
		Password:  "fuzz",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("AuthHandlers", func(t *testing.T) {
		// Register Missing username.

		missingUsernameJSON, err := json.Marshal(missingUsernameUser)
		assert.NoError(t, err, "expected no error marshalling missing username user data, got: %v.", err)

		// Convert the JSON string to a byte buffer.
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(missingUsernameJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Register Missing email.

		missingEmailJSON, err := json.Marshal(missingEmailUser)
		assert.NoError(t, err, "expected no error marshalling missing email user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(missingEmailJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Register Missing password.

		missingPasswordJSON, err := json.Marshal(missingPasswordUser)
		assert.NoError(t, err, "expected no error marshalling missing password user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(missingPasswordJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Register Invalid email.

		invalidEmailJSON, err := json.Marshal(invalidEmailUser)
		assert.NoError(t, err, "expected no error marshalling invalid email user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(invalidEmailJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// SUCCESS. Register valid user.
		validUserJSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling valid user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUserJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Attempt to register an email that already exists.

		duplicateUserJSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling duplicate user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(duplicateUserJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Log In user with missing email.

		credentials := url.Values{}
		credentials.Set("username", "")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))

		// Set the Content-Type header to application/x-www-form-urlencoded.
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Log in user with missing password.

		credentials = url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Log in non-existent user, valid email and password syntax.

		credentials = url.Values{}
		credentials.Set("username", "nonexistent@nonexistent.com")
		credentials.Set("password", "nonexistent")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)

		// Log in user 1, valid existing email, incorrect password.

		credentials = url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "incorrectpassword")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)

		// Log in user. SUCCESS.

		credentials = url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken := response["access_token"]

		// Delete the user.

		userID := validUser1.ID
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

		// Attempt to log in same credentials

		credentials = url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)

	})

	t.Run("UserHandlers", func(t *testing.T) {

		// SUCCESS. Register valid user.
		validUser1JSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling valid user data, got: %v.", err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUser1JSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Register another valid user.
		validUser2JSON, err := json.Marshal(validUser2)
		assert.NoError(t, err, "expected no error marshalling valid user 2 data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUser2JSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Log in user 1. SUCCESS.

		credentials := url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken := response["access_token"]

		// Get a user (invalid user id).

		userID := "invalid"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Get a user (user not found)
		userID = "1000"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// Get a user (unauthorized access)
		userID = "2"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// Get a user (success).
		userID = "1"

		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Create the custom JSON request body containing fields to be updated.
		updateData := map[string]interface{}{
			"username": "newusername",
			"email":    "newemail@newemail.com",
			"password": "newpassword",
		}

		// Marshal the updateData into a JSON string
		updateDataJSON, err := json.Marshal(updateData)
		assert.NoError(t, err, "expected no error marshalling user update data, got: %v.", err)

		// Update a user (invalid user ID)
		userID = "invalid"
		req = httptest.NewRequest("PATCH", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Update a user (user not found)

		userID = "1000"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// update a user (unauthorized access)

		userID = "2"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// update a user (all fields) (success)

		userID = "1"
		req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Delete a user (Invalid user ID).
		userID = "invalid"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Delete a user (User not found).
		userID = "1000"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// Delete another user (Unauthorized access).
		userID = "2"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// Delete user 1 (SUCCESS).
		userID = "1"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

		// Log in user 2.

		credentials = url.Values{}
		credentials.Set("username", "fuzz@buzz.com")
		credentials.Set("password", "fuzz")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken = response["access_token"]

		// Delete user 2.
		userID = "2"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)
	})

	missingTitleBook := &types.Book{
		ID:         1,
		Title:      "",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 1000,
		PagesRead:  500,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	missingAuthorBook := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Edition:    1,
		Author:     "",
		PagesCount: 1000,
		PagesRead:  500,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	invalidPagesCountBook := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 0,
		PagesRead:  0,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	invalidPagesReadBook := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 1000,
		PagesRead:  -1,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	invalidPagesBook := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 500,
		PagesRead:  1000,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	validBook1EditionFalse := &types.Book{
		ID:         1,
		Title:      "Book 1",
		Author:     "Tom",
		PagesCount: 1000,
		PagesRead:  500,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	validBook1EditionTrue := &types.Book{
		ID:         2,
		Title:      "Book 2",
		Edition:    1,
		Author:     "Tom",
		PagesCount: 1000,
		PagesRead:  500,
		OwnerID:    1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	validBook2EditionFalse := &types.Book{
		ID:         3,
		Title:      "Book 1",
		Author:     "Scott",
		PagesCount: 300,
		PagesRead:  150,
		OwnerID:    2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	validBook2EditionTrue := &types.Book{
		ID:         4,
		Title:      "Book 2",
		Edition:    1,
		Author:     "Scott",
		PagesCount: 300,
		PagesRead:  150,
		OwnerID:    2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	t.Run("BookHandlers", func(t *testing.T) {
		// Create a valid user 1 (SUCCESS).
		validUser1JSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling valid user 1 data, got: %v.", err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUser1JSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Log in a valid user 1 (SUCCESS).

		credentials := url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken := response["access_token"]

		// Create a book for user 1 (Missing title).

		missingTitleJSON, err := json.Marshal(missingTitleBook)
		assert.NoError(t, err, "expected no error marshalling missing title book data, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(missingTitleJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Create a book for user 1 (Missing author).
		missingAuthorJSON, err := json.Marshal(missingAuthorBook)
		assert.NoError(t, err, "expected no error marshalling missing author book data, got: %v.", err)
		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(missingAuthorJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Create a book for user 1 (invalid pages count).
		invalidPagesCountJSON, err := json.Marshal(invalidPagesCountBook)
		assert.NoError(t, err, "expected no error marshalling invalid pages count book data, got: %v.", err)
		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(invalidPagesCountJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Create a book for user 1 (invalid pages read).
		invalidPagesReadJSON, err := json.Marshal(invalidPagesReadBook)
		assert.NoError(t, err, "expected no error marshalling invalid pages read book data, got: %v.", err)
		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(invalidPagesReadJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Create a book for user 1 (pages read > pages count).
		invalidPagesJSON, err := json.Marshal(invalidPagesBook)
		assert.NoError(t, err, "expected no error marshalling invalid pages book data, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(invalidPagesJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Create a book for user 1 (no edition) (success).
		bookEditionFalse1, err := json.Marshal(validBook1EditionFalse)
		assert.NoError(t, err, "expected no error marshalling valid book data w/out edition, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(bookEditionFalse1))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Create another book for user 1 (yes edition) (success).
		bookEditionTrue1, err := json.Marshal(validBook1EditionTrue)
		assert.NoError(t, err, "expected no error marshalling valid book data w/edition, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(bookEditionTrue1))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Create a valid user 2 (SUCCESS).
		validUser2JSON, err := json.Marshal(validUser2)
		assert.NoError(t, err, "expected no error marshalling valid user 2 data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUser2JSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Log in valid user 2.

		credentials = url.Values{}
		credentials.Set("username", "fuzz@buzz.com")
		credentials.Set("password", "fuzz")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken = response["access_token"]

		// Create a book for valid user 2 (no edition, success).

		bookEditionFalse2, err := json.Marshal(validBook2EditionFalse)
		assert.NoError(t, err, "expected no error marshalling valid book data w/out edition, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(bookEditionFalse2))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Create another book for valid user 2 (yes edition, success).
		bookEditionTrue2, err := json.Marshal(validBook2EditionTrue)
		assert.NoError(t, err, "expected no error marshalling valid book data w/edition, got: %v.", err)

		req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(bookEditionTrue2))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Get all books for valid user 2. (success).

		req = httptest.NewRequest("GET", "/books/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Get a book for valid user 2 with a specified id (invalid request id).
		bookID := "invalid"

		req = httptest.NewRequest("GET", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Get a book for valid user 2 with a specified id (user not found).
		bookID = "1000"

		req = httptest.NewRequest("GET", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// Get a book for valid user 2 with a specified id (unauthorized access).
		bookID = "1"

		req = httptest.NewRequest("GET", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// Get a book for valid user 2 with a specified id (success).
		bookID = "3"

		req = httptest.NewRequest("GET", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Create the custom JSON request body containing fields to be updated.
		updateData := map[string]interface{}{
			"title":       "newtitle",
			"author":      "newauthor",
			"edition":     100,
			"pages_count": 2000,
			"pages_read":  1500,
		}

		// Marshal the updateData into a JSON string
		updateDataJSON, err := json.Marshal(updateData)
		assert.NoError(t, err, "expected no error marshalling book update data, got: %v.", err)

		// Update a book for valid user 2 (invalid book ID).
		bookID = "invalid"
		req = httptest.NewRequest("PATCH", fmt.Sprintf("/books/%s", bookID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Update a book for valid user 2. (book not found)

		bookID = "1000"
		req = httptest.NewRequest("PATCH", fmt.Sprintf("/books/%s", bookID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// Update a book for valid user 2. (unauthorized access)
		bookID = "1"
		req = httptest.NewRequest("PATCH", fmt.Sprintf("/books/%s", bookID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// Update a book for valid user 2. (success)
		bookID = "3"
		req = httptest.NewRequest("PATCH", fmt.Sprintf("/books/%s", bookID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Delete a book for valid user 2. (invalid request id)
		bookID = "invalid"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		// Delete a book for valid user 2. (book not found)
		bookID = "1000"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		// Delete a book for valid user 2. (unauthorized access)
		bookID = "1"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)

		// Delete a book for valid user 2. (success)
		bookID = "3"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/books/%s", bookID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

		// Delete user 2 (books delete on cascade). (success)
		userID := "2"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

		// Log in User 1 (success).
		credentials = url.Values{}
		credentials.Set("username", "foo@bar.com")
		credentials.Set("password", "foo")

		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

		accessToken = response["access_token"]

		// Delete user 1 (success, books delete on cascade).

		userID = "1"
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)
	})

	server.router = nil
	time.Sleep(500 * time.Millisecond)
}
