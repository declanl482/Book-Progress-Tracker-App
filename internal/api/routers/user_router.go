package routers

import (
	"example/go-book-tracker-app/internal/api/handlers"
	"example/go-book-tracker-app/internal/middlewares"
	"example/go-book-tracker-app/internal/middlewares/oauth2"

	"github.com/gin-gonic/gin"
)

func ConfigureUserRoutes(router *gin.Engine) {

	// Register the routes for Users

	users := router.Group("/users")
	users.Use(middlewares.CORSMiddleware())
	users.Use(middlewares.DBConnectionMiddleware())
	users.Use(oauth2.JWTAccessTokenMiddleware())
	{
		users.GET("/", handlers.GetUser)
		users.PATCH("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}
}
