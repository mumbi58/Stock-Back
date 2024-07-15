package main

import (
    "log"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/joho/godotenv"
    //"gorm.io/gorm"
    "github.com/go-playground/validator/v10"
    "stock-back/db"
    "stock-back/routes"
)

func main() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Initialize database connection
    db.Init()

    // Create new Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Initialize validator
    e.Validator = &CustomValidator{Validator: validator.New()}

    // Setup routes
    routes.SetupRouter(e, db.GetDB()) // Pass the *gorm.DB instance

    // Start server
    port := ":8080"
    log.Printf("Server started at %s", port)
    e.Logger.Fatal(e.Start(port))
}

// CustomValidator struct to implement custom validator
type CustomValidator struct {
    Validator *validator.Validate
}

// Validate performs struct validation
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.Validator.Struct(i)
}
