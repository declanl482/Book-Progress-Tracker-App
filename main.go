package main

import (
	"fmt"

	"github.com/declanl482/go-book-tracker-app/api"
	"github.com/declanl482/go-book-tracker-app/config"
	"github.com/declanl482/go-book-tracker-app/storage"
)

func main() {

	// listenAddress := flag.String("listenAddress", ":8000", "the server address")
	listenAddress := ":8000"

	configuration, err := config.LoadConfigurationVariables()
	if err != nil {
		fmt.Println("Failed to load configuration variables:", err)
		return
	}

	hostname := configuration.DatabaseHostname
	username := configuration.DatabaseUsername
	password := configuration.DatabasePassword
	name := configuration.DatabaseName
	port := configuration.DatabasePort
	timezone := configuration.DatabaseTimezone

	// Create a new instance of PostgresStorage.
	postgresStorage, err := storage.NewPostgresStorage(hostname, username, password, name, port, timezone)
	if err != nil {
		// Handle the error if any.
		panic(err)
	}

	// Create a new instance of the Server with the UserStorage and BookStorage implementations.
	server := api.NewServer(listenAddress, postgresStorage)

	// Start the server.
	err = server.Start()
	if err != nil {
		// Handle the error if any.
		panic(err)
	}
}
