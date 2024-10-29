package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatal("No POSTGRES_URL environment variable")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Fprintf(w, "Successfully connected to the database!")
}
