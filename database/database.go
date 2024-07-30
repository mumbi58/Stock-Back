package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // Assuming you're using MySQL
)

func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	return db
}
