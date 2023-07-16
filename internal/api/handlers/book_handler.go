package handlers

import (
	"net/http"
	"strconv"

	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/database"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(*models.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var newBook models.Book

	// Call ShouldBindJSON to bind received JSON to newBook
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the owner id of the book
	newBook.OwnerID = currentUser.ID

	// save the new book to the database
	result := database.GetInstanceOfApplicationDatabase().Create(&newBook)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	// append the new book to the current user's book collection
	// currentUser.Books = append(currentUser.Books, newBook)

	c.JSON(http.StatusCreated, newBook)
}

// getBooks responds with a list of all books as JSON

func GetBooks(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(*models.User)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var books []models.Book

	// Retrieve all books from the database where book[i].ownerID == currentUser.ID

	result := database.GetInstanceOfApplicationDatabase().Where("owner_id = ?", currentUser.ID).Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve books"})
		return
	}

	c.IndentedJSON(http.StatusOK, books)
}

// getBook locates the book whose ID value matches the id passed in by the client,
// then returns that book as a JSON response.

func GetBook(c *gin.Context) {

	// check that the user is authenticated
	currentUser := c.MustGet("currentUser").(*models.User)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// parse the request url for the id of the requested book
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch the requested book from the database

	book := models.Book{}
	result := database.GetInstanceOfApplicationDatabase().First(&book, id)

	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// check if the requested book belongs to the current user

	if book.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot access this book"})
		return
	}

	// Return the book
	c.JSON(http.StatusOK, book)
}

// updateBook locates the book whose id matches the input id, then updates the book to match
// the request data provided by the client, then returns the updated book as a JSON response.

func UpdateBook(c *gin.Context) {

	// check that the user is authenticated
	currentUser := c.MustGet("currentUser").(*models.User)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// parse the request url for the id of the requested book
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch book from database
	book := models.Book{}

	result := database.GetInstanceOfApplicationDatabase().First(&book, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Check if current user owns the book

	if book.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot update this book"})
		return
	}

	// Bind the updated JSON data from the request context to the extracted book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update timestamp
	book.UpdatedAt = time.Now()

	// Perform the update operation on the user
	result = database.GetInstanceOfApplicationDatabase().Model(&book).Updates(book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// delete book locates a book by its id, returns 204 no content response if successful

func DeleteBook(c *gin.Context) {

	// Get the current user sending the delete request
	currentUser := c.MustGet("currentUser").(*models.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse the requested id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch book from database
	book := models.Book{}
	result := database.GetInstanceOfApplicationDatabase().First(&book, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Check if requesting user id == queried user id
	if book.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot delete this book"})
		return
	}

	// Delete the user from the DB

	result = database.GetInstanceOfApplicationDatabase().Delete(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}
