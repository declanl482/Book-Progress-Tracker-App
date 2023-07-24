package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/declanl482/go-book-tracker-app/types"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleLoginUser(c *gin.Context) {
	var credentials types.Credentials

	// Bind the form data to the credentials.
	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that there exists a record with the given email.
	user, err := s.Storer.GetUserByEmail(credentials.Email)
	if err != nil || user == nil {
		fmt.Println("failed to fetch!")
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid credentials"})
		return
	}

	fmt.Println("database user password:", user.Password)
	fmt.Println("client credentials password:", credentials.Password)

	// Verify the provided password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(strings.TrimSpace(credentials.Password))); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create a JWT access token using the authenticated user ID.
	auth := NewAuth(config.Config.AccessTokenSecretKey) // Replace with your secret key.
	tokenString, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// Send the access token in the response.
	c.JSON(http.StatusOK, gin.H{"access_token": tokenString})

}

func (s *Server) handleCreateUser(c *gin.Context) {
	// Declare a new user variable.
	var newUser types.User

	// Bind the JSON request body to the new user variable.
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the new user's email is already in use.

	emailTaken, err := s.Storer.IsEmailTaken(newUser.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if emailTaken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is taken"})
		return
	}

	// Hash the new user's password.

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUser.Password = string(hashedPassword)

	// Create the new user in the database.
	createdUser, err := s.Storer.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusCreated, createdUser)
}

func (s *Server) handleGetUser(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Fetch the user from the database.
	fetchedUser, err := s.Storer.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// There is no user with the requested id.
	if fetchedUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// Check that the client is authorized to view the fetched user.
	if userID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot view this user"})
		return
	}

	// Return the user data in the response.
	c.IndentedJSON(http.StatusOK, fetchedUser)
}

func (s *Server) handleUpdateUser(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Fetch the user from the database.
	fetchedUser, err := s.Storer.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// There is no user with the requested id.
	if fetchedUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// Check that the client is authorized to update the fetched user.
	if userID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot update this user"})
		return
	}

	// Check if the password field is included in the JSON request.
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := requestBody["password"]; ok {
		// Password is newly updated, rehash it.

		fmt.Println("The password is newly updated.")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody["password"].(string)), 14)

		fmt.Println("Hashed the new password:", hashedPassword)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		requestBody["password"] = string(hashedPassword)

		fmt.Println("Just checking the request body is hashed:", requestBody["password"])
	}

	// Update the updated_at time stamp.
	requestBody["updated_at"] = time.Now()

	// convert the map to the fetchedUser struct of type *types.User
	if err := mapstructure.Decode(requestBody, &fetchedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Update the fetched user in the database.
	updatedUser, err := s.Storer.UpdateUser(fetchedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (s *Server) handleDeleteUser(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	// Fetch the user from the database.
	fetchedUser, err := s.Storer.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// There is no user with the requested id.
	if fetchedUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check that the client is authorized to delete the fetched user.
	if userID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot delete this user"})
		return
	}

	// Delete the user from the database.
	err = s.Storer.DeleteUser(fetchedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusNoContent, nil)
}

func (s *Server) handleCreateBook(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var newBook *types.Book

	// Bind the request body to the new book variable.
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook.OwnerID = currentUser.ID

	// invalid pages count / pages read.

	if newBook.PagesCount <= 0 || newBook.PagesRead < 0 || newBook.PagesRead > newBook.PagesCount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pages data"})
		return
	}

	// Create the book in the database.
	createdBook, err := s.Storer.CreateBook(newBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusCreated, createdBook)
}

func (s *Server) handleGetBooks(c *gin.Context) {
	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	books, err := s.Storer.GetBooks(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch books"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusOK, books)
}

func (s *Server) handleGetBook(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	// Fetch the book from the database
	book, err := s.Storer.GetBook(bookID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch book"})
		return
	}

	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Verify that the user can access the fetched book.
	if book.OwnerID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot view this book"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusOK, book)

}

func (s *Server) handleUpdateBook(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	// Fetch the book from the database.
	fetchedBook, err := s.Storer.GetBook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book"})
		return
	}

	// There is no book with the requested id.
	if fetchedBook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	// Check that the client is authorized to update the fetched book.
	if fetchedBook.OwnerID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot update this book"})
		return
	}

	// Bind the JSON request body to the fetched book variable.
	if err := c.ShouldBindJSON(&fetchedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the updated_at time stamp.
	fetchedBook.UpdatedAt = time.Now()

	// Update the fetched book in the database.
	updatedBook, err := s.Storer.UpdateBook(fetchedBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusOK, updatedBook)

}

func (s *Server) handleDeleteBook(c *gin.Context) {

	// Get the authenticated user from the context.
	currentUser := c.MustGet("currentUser").(*types.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract the id param from the URL request path.
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	// Fetch the book from the database.
	fetchedBook, err := s.Storer.GetBook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book"})
		return
	}

	// There is no book with the requested id.
	if fetchedBook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	// Check that the client is authorized to delete the fetched book.
	if fetchedBook.OwnerID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot delete this book"})
		return
	}

	// Delete the book from the database.
	err = s.Storer.DeleteBook(fetchedBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
		return
	}

	// SUCCESS.
	c.IndentedJSON(http.StatusNoContent, nil)
}
