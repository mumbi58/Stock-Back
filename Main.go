package main

import (
	"awesomeProject9/database"
	"awesomeProject9/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Echo
	e := echo.New()

	// Connect to the database
	db := database.InitDB()
	defer db.Close()

	// Middleware to enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGINS")},
		AllowHeaders: []string{os.Getenv("CORS_ALLOW_HEADERS")},
	}))
	// Register routes
	routes.RegisterRoutes(e)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	e.Logger.Fatal(e.Start(":" + port))
}
