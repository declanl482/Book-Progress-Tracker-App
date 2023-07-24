package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"example/go-book-tracker-app/api/routers"
	"example/go-book-tracker-app/config"
	"example/go-book-tracker-app/database"

	_ "github.com/lib/pq"
)

func main() {

	err := config.LoadApplicationConfigurationVariables()
	if err != nil {
		fmt.Println("Failed to load application configuration variables:", err)
		return
	}

	_, err = database.ConnectToApplicationDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	router := gin.Default()
	routers.ConfigureAuthenticationRoutes(router)
	routers.ConfigureUserRoutes(router)
	routers.ConfigureBookRoutes(router)

	router.Run("localhost:8000")
}
