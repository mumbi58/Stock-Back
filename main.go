package main

import (
    "stock-back/db"
    "stock-back/routes"
    "github.com/labstack/echo/v4"
)

func main() {
    // Initialize the database
    db.Init()

    // Create a new Echo instance
    e := echo.New()

    // Set up routes
    routes.SetupRoutes(e, db.GetDB())

    // Start the server
    e.Logger.Fatal(e.Start(":8080"))
}
