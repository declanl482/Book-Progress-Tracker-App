package handlers

import (
	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/database"
	"example/go-book-tracker-app/internal/middlewares/oauth2"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c *gin.Context) {
	var credentials models.UserCredentials

	// Bind the request body to the UserCredentials struct
	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verify the user . Fetch them from the database via email
	var user models.User
	result := database.GetDB().Where("email = ?", credentials.Email).First(&user)
	if result.Error != nil {
		fmt.Println("failed to fetch!")
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid credentials"})
		return
	}

	fmt.Println(user.Password, credentials.Password)

	fmt.Println(bcrypt.CompareHashAndPassword([]byte("$2a$14$GLfQAHOBXMMHreYUzfL7VOXrqTRmBfVepdkFyO502x9K9FAUv5obO"), []byte{}))

	// verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		fmt.Println("failed to verify password!")
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create an access token
	fmt.Println(strconv.Itoa(int(user.ID)))
	accessToken, err := oauth2.CreateJWTAccessToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	// return the access token to the client
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
