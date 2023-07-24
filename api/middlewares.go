package api

import (
	"net/http"
	"strings"

	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/gin-gonic/gin"
)

// Middleware to validate the access token, get the current user, and set it in the context.
func (s *Server) RequireValidAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token from the request header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access token required"})
			c.Abort()
			return
		}

		accessToken := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validate the access token and get the user details.
		auth := NewAuth(config.Config.AccessTokenSecretKey)
		userID, err := auth.ValidateAccessToken(accessToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		// Retrieve the user from the database using the userID.
		user, err := s.Storer.GetUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
			c.Abort()
			return
		}

		// Add the user to the context.
		c.Set("currentUser", user)

		// Continue to the next handler.
		c.Next()
	}
}

func (s *Server) DBConnectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.Storer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, X-Custom-Header")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			requestedMethod := c.GetHeader("Access-Control-Request-Method")
			if !IsMethodSupported(requestedMethod) {
				c.AbortWithStatus(http.StatusMethodNotAllowed)
				return
			}
			c.AbortWithStatus(http.StatusNoContent) // Return empty response on success
			return
		}
	}
}

func IsMethodSupported(method string) bool {
	// check if the method is supported
	supportedMethods := []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	for _, supportedMethod := range supportedMethods {
		if method == supportedMethod {
			return true
		}
	}
	return false
}
