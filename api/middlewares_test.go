package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/declanl482/go-book-tracker-app/storage"
	"github.com/declanl482/go-book-tracker-app/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareFunctions(t *testing.T) {

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

	// construct a new server
	server := NewServer(listenAddress, store) // use server to access handler functions

	// Start the server.
	go func() {
		err := server.Start() // Test the starting of the server (register middlewares/handlers + run the router on the listen address)
		assert.NoError(t, err, "expected no error when starting server, routers, and middlewares, got: %v.", err)
	}()

	invalidAuth := NewAuth("invalid-access-key")

	validUser1 := &types.User{
		ID:        1,
		Username:  "foo",
		Email:     "foo@bar.com",
		Password:  "foo",
		Books:     []types.Book{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("TestRequireValidAccessToken", func(t *testing.T) {

		// Generate an invalid mock access token.
		invalidAccessToken, err := invalidAuth.GenerateAccessToken(-1)
		assert.NoError(t, err, "expected no error generating invalid access token, got: %v.", err)

		// Create a new HTTP request with the required headers.
		req := httptest.NewRequest("GET", "/books/", nil)
		req.Header.Set("Authorization", "Bearer "+invalidAccessToken)

		// Create a new mock Gin context.
		w := httptest.NewRecorder()
		// server.router.ServeHTTP(w, req)
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// Call the RequireValidAccessToken middleware function passing the mock context.
		server.RequireValidAccessToken()(c)
		// successfully got books, meaning access token was successfully required.
		assert.Equal(t, 401, w.Code)

		// Create a valid user
		validUserJSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling valid user data, got: %v.", err)

		req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUserJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Log the user in, collect the valid access token:
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

		validAccessToken := response["access_token"]

		// Create a new HTTP request with the required headers.
		req = httptest.NewRequest("GET", "/books/", nil)
		req.Header.Set("Authorization", "Bearer "+validAccessToken)

		// Create a new mock Gin context.
		w = httptest.NewRecorder()
		// server.router.ServeHTTP(w, req)
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		// Call the RequireValidAccessToken middleware function passing the mock context.
		server.RequireValidAccessToken()(c)

		// Check that the user was successfully set in context.
		currentUser := c.MustGet("currentUser").(*types.User)

		assert.NotEmpty(t, currentUser)

		// successfully got books, meaning access token was successfully required.
		assert.Equal(t, 200, w.Code)

		// Delete the user
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", validUser1.ID), nil)
		req.Header.Set("Authorization", "Bearer "+validAccessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

	})

	t.Run("TestDBConnectionMiddleware", func(t *testing.T) {

		// Create a valid user
		validUserJSON, err := json.Marshal(validUser1)
		assert.NoError(t, err, "expected no error marshalling valid user data, got: %v.", err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(validUserJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)

		// Log the user in, collect the valid access token:
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

		validAccessToken := response["access_token"]

		req = httptest.NewRequest("GET", "/books/", nil)
		req.Header.Set("Authorization", "Bearer "+validAccessToken)

		// Create a new mock Gin context.
		w = httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the database connection middleware function passing the mock context.
		server.DBConnectionMiddleware()(c)

		// successfully got books, meaning db connection was successful.
		assert.Equal(t, 200, w.Code)

		// delete the user

		req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", validUser1.ID), nil)
		req.Header.Set("Authorization", "Bearer "+validAccessToken)

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)

		dummyListenAddress := ":9000"

		// make a server instance with a null storer.
		serverInstance := NewServer(dummyListenAddress, nil)

		req = httptest.NewRequest("GET", "/books/", nil)
		// Call the database connection middleware function passing the mock context.
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		serverInstance.DBConnectionMiddleware()(c)

		// internal server error, meaning database connection unsuccessful.
		assert.Equal(t, 500, w.Code)
	})

	t.Run("TestCORSMiddleware", func(t *testing.T) {
		// Create a new test context.
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/books/", nil)

		// Call the CORSMiddleware function passing the mock context.
		server.CORSMiddleware()(c)

		// Check the response headers.
		assert.Equal(t, "*", c.Writer.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", c.Writer.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "GET, POST, PATCH, DELETE, OPTIONS", c.Writer.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Origin, Authorization, Content-Type, X-Custom-Header", c.Writer.Header().Get("Access-Control-Allow-Headers"))

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/books/", nil)

		// Test OPTIONS GET method.
		c.Request.Method = "OPTIONS"
		c.Request.Header.Set("Access-Control-Request-Method", "GET")
		server.CORSMiddleware()(c)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Create a new mock Gin context.
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/books/", nil)

		// Test OPTIONS POST method.
		c.Request.Method = "OPTIONS"
		c.Request.Header.Set("Access-Control-Request-Method", "POST")
		server.CORSMiddleware()(c)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// new mock context.
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/books/:id", nil)

		// Test OPTIONS PATCH method.
		c.Request.Method = "OPTIONS"
		c.Request.Header.Set("Access-Control-Request-Method", "PATCH")
		server.CORSMiddleware()(c)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// new mock context
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/books/:id", nil)

		// Test OPTIONS PUT method.
		c.Request.Method = "OPTIONS"
		c.Request.Header.Set("Access-Control-Request-Method", "PUT")
		server.CORSMiddleware()(c)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

		// Test OPTIONS DELETE method.
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/books/:id", nil)

		// Test OPTIONS PUT method.
		c.Request.Method = "OPTIONS"
		c.Request.Header.Set("Access-Control-Request-Method", "DELETE")
		server.CORSMiddleware()(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
