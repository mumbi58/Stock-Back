package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	dsn := "pascal:@mesopotamia123@tcp(localhost:3306)/storage"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
