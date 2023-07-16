package middlewares

import (
	"example/go-book-tracker-app/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBConnectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the database connection is valid
		if database.GetInstanceOfApplicationDatabase() == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
			c.Abort()
			return
		}
		// Continue to the next handler
		c.Next()
	}
}
