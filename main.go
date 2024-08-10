package main

import (
	"awesomeProject9/database"
	"awesomeProject9/routes"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	// Initialize the database
	database.Init() // Changed from InitDB to Init

	// Create a new Echo instance
	e := echo.New()

	// Set up routes
	routes.RegisterRoutes(e) // Only pass the Echo instance

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	log.Fatal(e.Start(":" + port))
}
