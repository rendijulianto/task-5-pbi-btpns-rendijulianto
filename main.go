package main

import (
	"final-project-rakamin/database"
	"final-project-rakamin/router"
)

func main() {
	// Connect to database
	database.Connect()
	// Run migrations
	database.Migrate()

	// Setup the router
	r := router.SetupRouter()

	// Run the server
	r.Run(":8080")
}

