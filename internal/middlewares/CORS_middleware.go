package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {

			if !IsMethodSupported(c.Request.Header.Get("Access-Control-Request-Method")) {
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
