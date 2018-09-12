package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Blank import to make native sql package work with PostgreSQL
)

// Query queries the database with the connection string from env vars
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	return db.Query(query, args...)
}

// Exec executes database code with the connection string from env vars
func Exec(query string, args ...interface{}) (sql.Result, error) {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	return db.Exec(query, args...)
}
