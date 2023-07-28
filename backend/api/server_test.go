package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/declanl482/go-book-tracker-app/backend/config"
	"github.com/declanl482/go-book-tracker-app/backend/storage"
	"github.com/declanl482/go-book-tracker-app/backend/types"
	"github.com/stretchr/testify/assert"
)

func TestServerFunctions(t *testing.T) {

	// create a new postgres storage (test db), and use it to create a new server
	testConfig, err := config.LoadTestConfigurationVariables()
	assert.NoError(t, err, "expected no error when loading test config variables, got: %v.", err)

	hostname := testConfig.TestDatabaseHostname
	username := testConfig.TestDatabaseUsername
	password := testConfig.TestDatabasePassword
	name := testConfig.TestDatabaseName
	port := testConfig.TestDatabasePort
	timezone := testConfig.TestDatabaseTimezone

	store, err := storage.NewPostgresStorage(hostname, username, password, name, port, timezone)
	assert.NoError(t, err, "expected no error when creating PostgresStorage for testing database, got: %v.", err)
	listenAddress := ":8080"

	// construct a new server
	server := NewServer(listenAddress, store) // use server to access handler functions

	// Start the server.
	go func() {
		err := server.Start() // Test the starting of the server (register middlewares/handlers + run the router on the listen address)
		assert.NoError(t, err, "expected no error when starting server, routers, and middlewares, got: %v.", err)
	}()

	user := &types.User{}
	userJSON, err := json.Marshal(user)
	assert.NoError(t, err, "expected no error marshalling user data, got: %v.", err)

	book := &types.User{}
	bookJSON, err := json.Marshal(book)
	assert.NoError(t, err, "expected no error marshalling book data, got: %v.", err)

	// Test server routes as an unauthorized client and with empty user inputs for login/registration.

	// (1) Register route (Create User).
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// (2) Login route.

	credentials := url.Values{}
	credentials.Set("username", "")
	credentials.Set("password", "")

	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(credentials.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "expected no error unmarshalling access token from login response, got: %v.", err)

	accessToken := response["access_token"]

	// (3) Get User route.
	req = httptest.NewRequest("GET", "/users/:id", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (4) Update User route.
	req = httptest.NewRequest("PATCH", "/users/:id", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (5) Delete User route.
	req = httptest.NewRequest("DELETE", "/users/:id", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (6) Create Book route.
	req = httptest.NewRequest("POST", "/books/", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (7) Get Books route.
	req = httptest.NewRequest("GET", "/books/", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (8) Get Book route.
	req = httptest.NewRequest("GET", "/books/:id", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (9) Update Book route.
	req = httptest.NewRequest("PATCH", "/books/:id", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// (10) Delete Book route.
	req = httptest.NewRequest("DELETE", "/books/:id", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// stop serving the router
	server.router = nil
}
