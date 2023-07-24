package routers

import (
	"example/go-book-tracker-app/api/handlers"
	"example/go-book-tracker-app/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigureUserRoutes(router *gin.Engine) {

	// Register the routes for Users

	users := router.Group("/users")
	users.Use(middlewares.CORSMiddleware())
	users.Use(middlewares.DBConnectionMiddleware())
	users.Use(middlewares.JWTAccessTokenMiddleware())
	{
		users.GET("/", handlers.GetUser)
		users.PATCH("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}
}
