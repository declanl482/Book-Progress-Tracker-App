package middlewares

import (
	"example/go-book-tracker-app/internal/api/models"
	"example/go-book-tracker-app/internal/config"
	"example/go-book-tracker-app/internal/database"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateJWTAccessToken(userID string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	secretKey := []byte(config.AppConfig.AccessTokenSecretKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWTAccessToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.AccessTokenSecretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Extract the user id claim from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user id claim")
	}
	return userID, nil
}

func GetCurrentUserFromAccessToken(c *gin.Context) *models.User {

	// get access token from authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}
	accessToken := strings.Replace(authHeader, "Bearer ", "", 1)

	// verify the access token
	userID, err := VerifyJWTAccessToken(accessToken)
	if err != nil {
		return nil
	}

	// retrieve the user from the database using their id
	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {

		return nil
	}
	return &user
}

func JWTAccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUserFromAccessToken(c)
		if user == nil {
			// the access token is not verified

			c.JSON(http.StatusUnauthorized, gin.H{"error": "could not verify access token"})
			c.Abort()
			return
		}
		// access token is verified
		// Set the current user in the context
		c.Set("currentUser", user)

		// Continue to the next handler
		c.Next() // move onto the refresh token middleware
	}
}
