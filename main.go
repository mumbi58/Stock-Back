package main

import (
	"log"
	"os"
	"stock-back/db"
	"stock-back/routes"

	//"github.com/labstack/echo/v4"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	db.Init()
}

func main() {
	// Setup routes
	e := routes.Setup()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
