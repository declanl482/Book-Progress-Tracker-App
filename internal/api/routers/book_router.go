package routers

import (
	"example/go-book-tracker-app/internal/api/handlers"
	"example/go-book-tracker-app/internal/middlewares"
	"example/go-book-tracker-app/internal/middlewares/oauth2"

	"github.com/gin-gonic/gin"
)

func ConfigureBookRoutes(router *gin.Engine) {

	// Register the routes for Books

	books := router.Group("/books")
	books.Use(middlewares.CORSMiddleware())
	books.Use(middlewares.DBConnectionMiddleware())
	books.Use(oauth2.JWTAccessTokenMiddleware())
	{
		books.POST("/", handlers.CreateBook)
		books.GET("/", handlers.GetBooks)
		books.GET("/:id", handlers.GetBook)
		books.PATCH("/:id", handlers.UpdateBook)
		books.DELETE("/:id", handlers.DeleteBook)
	}
}
