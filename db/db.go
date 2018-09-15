package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq" // Blank import to make native sql package work with PostgreSQL
)

// JSONNullString is a wrapper around sql.NullString which implements
// the MarshalJSON and UnmarshalJSON methods
type JSONNullString struct {
	sql.NullString
}

// MarshalJSON marshals the data based on whether it is valid
// or not
func (s JSONNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON takes the raw data and decides whether it is
// valid or not, and sets the data fields accordingly
func (s *JSONNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		s.Valid = true
		s.String = *x
	} else {
		s.Valid = false
	}
	return nil
}

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
