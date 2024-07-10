package db

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "stock-back/models"
    "stock-back/utils"
)

var DB *gorm.DB

func Init() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // MySQL connection string
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    // Connect to MySQL database
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto migration of models
    err = database.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Assign database connection to utils.DB
    utils.DB = database

    fmt.Println("Database connected and migrated")
}

func GetDB() *gorm.DB {
    return utils.DB
}
