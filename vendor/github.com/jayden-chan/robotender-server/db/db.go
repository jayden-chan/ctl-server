package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq" // Blank import to make native sql package work with PostgreSQL
)

// Drink is a struct storing the data for a drink
type Drink struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Glass string          `json:"glass"`
	Tall  bool            `json:"tall"`
	Spec  json.RawMessage `json:"spec"`
	Note  string          `json:"note"`
}

// DrinkArray is an collection of Drink structs
type DrinkArray struct {
	Drinks []Drink `json:"drinks"`
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

// GetPresetType returns the presets of the specified type
func GetPresetType(category string) (drinks DrinkArray, err error) {
	var rows *sql.Rows

	if category == "" {
		query := `
		SELECT id, name, type, glass, tall, spec, note FROM drinks
		WHERE type IN ($1, $2, $3) AND custom = false`

		rows, err = Query(query, "standard", "shots", "complex")
	} else if category == "standard" || category == "shots" || category == "complex" {
		query := `
		SELECT id, name, type, glass, tall, spec, note FROM drinks
		WHERE type = $1 AND custom = false`

		rows, err = Query(query, category)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r Drink

		err = rows.Scan(&r.ID, &r.Name, &r.Type, &r.Glass, &r.Tall, &r.Spec, &r.Note)
		if err != nil {
			log.Println("getDrinkType Scan Error:", err)
		}
		drinks.Drinks = append(drinks.Drinks, r)
	}
	return
}

// GetFavorites returns the favorite drinks for the user
func GetFavorites(userID string) (drinks DrinkArray, err error) {
	query := `
	SELECT id, name, type, glass, tall, spec, note FROM drinks
	INNER JOIN favorites
	ON drinks.id = favorites.drink_id
	WHERE favorites.user_id = $1`

	rows, err := Query(query, userID)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var r Drink

		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Glass, &r.Tall, &r.Spec, &r.Note)
		if err != nil {
			log.Println("getDrinkType Scan Error:", err)
		}
		drinks.Drinks = append(drinks.Drinks, r)
	}
	return
}

// GetCustoms returns the custom drinks for the specified user
func GetCustoms(customerID string) (drinks DrinkArray, err error) {
	query := `
	SELECT id, name, type, glass, tall, spec, note FROM drinks
	WHERE custom = true AND customer_id = $1`

	rows, err := Query(query, customerID)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var r Drink

		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Glass, &r.Tall, &r.Spec, &r.Note)
		if err != nil {
			log.Println("getDrinkType Scan Error:", err)
		}
		drinks.Drinks = append(drinks.Drinks, r)
	}
	return
}
