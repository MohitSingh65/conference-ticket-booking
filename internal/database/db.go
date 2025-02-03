package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// Connect to database
func Connect() {
	dbURL := "postgres://postgres:@localhost/conference_db?sslmode=disable"

	var err error
	DB, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Connected to the database")

	// Run migrations
	Migrate()
}

func Migrate() {
	schema := `
  CREATE TABLE IF NOT EXISTS tickets (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	fmt.Println("Database migration completed successfully")
}
