// Ensure db package handles initialization and retrieval correctly
package db

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
    "os"
)

var (
    db *gorm.DB
)

// Init initializes the database connection
func Init() {
    dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME") + "?parseTime=true"

    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Enable debug mode if needed
    // db.Debug().AutoMigrate(&models.User{}, &models.Role{}) // Example for auto migration
}

// GetDB returns the instance of *gorm.DB
func GetDB() *gorm.DB {
    return db
}
