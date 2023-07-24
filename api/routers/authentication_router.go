package routers

import (
	"example/go-book-tracker-app/api/handlers"
	"example/go-book-tracker-app/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigureAuthenticationRoutes(router *gin.Engine) {

	// Register the routes for Authentication

	auth := router.Group("/auth")
	auth.Use(middlewares.CORSMiddleware())
	auth.Use(middlewares.DBConnectionMiddleware())
	{
		auth.POST("/register", handlers.CreateUser)
		auth.POST("/login", handlers.LoginUser)
	}
}
