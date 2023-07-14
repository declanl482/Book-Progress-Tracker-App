package handlers

import (
	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/database"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		fmt.Println("error occurs at binding")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("original password:", newUser.Password)

	db := database.GetDB()

	// check if the email is actively in use

	var existingUser models.User
	result := db.Where("email = ?", newUser.Email).First(&existingUser)

	if result.Error == nil {
		// fmt.Println("found someone")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is taken"})
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	fmt.Println("hashed password:", hashedPassword)
	newUser.Password = string(hashedPassword)
	fmt.Println("updated new user password:", newUser.Password)

	// add the new user to the database

	result = db.Create(&newUser)
	if err := result.Error; err != nil {
		fmt.Println("Error occurs when we try to add user to databse using GORM.Create")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	fmt.Println("the user was created")
	c.JSON(http.StatusCreated, newUser)

}

func GetUser(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(*models.User)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	// find the user with the current user's id, preload their books
	result := database.GetDB().Preload("Books").First(&user, currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.ID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot view user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// func GetUser(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}

// 	for i := range users {
// 		if uint64(users[i].ID) == id {
// 			c.IndentedJSON(http.StatusOK, users[i])
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
// }

func UpdateUser(c *gin.Context) {

	// Get the current user sending the update request
	currentUser := c.MustGet("currentUser").(*models.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse the requested id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch the user from the database

	user := models.User{}
	result := database.GetDB().First(&user, id)

	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// check if the requested user belongs to the current user
	if user.ID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot update this user"})
		return
	}

	// Bind the updated JSON data from the request context to the extracted user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UpdatedAt = time.Now()

	// Perform the update operation on the user
	result = database.GetDB().Model(&user).Updates(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {

	// Get the current user sending the delete request
	currentUser := c.MustGet("currentUser").(*models.User)

	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse the requested id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch user from database
	user := models.User{}
	result := database.GetDB().First(&user, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if requesting user id == queried user id
	if user.ID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot delete this user"})
		return
	}

	// Delete the user from the DB

	result = database.GetDB().Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}
