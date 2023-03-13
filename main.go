package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func main() {
	// Get the database connection parameters from environment variables
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname := os.Getenv("PGDATABASE")

	// Create a connection string using the parameters
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	var err error
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Define the HTTP handlers
	http.HandleFunc("/", homePage)
	http.HandleFunc("/search", searchHandler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Render the home page with a search form
func homePage(w http.ResponseWriter, r *http.Request) {
	// Render the search form
	fmt.Fprintf(w, `
		<h1>Full Text Search</h1>
		<form method="POST" action="/search">
			<input type="text" name="query" placeholder="Search...">
			<input type="submit" value="Search">
		</form>
	`)
}

// Perform a full-text search and render the results
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get the search query from the form data
	query := r.FormValue("query")

	// Execute the search query
	rows, err := db.Query(context.Background(), `
		SELECT text
		FROM documents
		WHERE to_tsvector('english', text) @@ to_tsquery('english', $1)
	`, query)
	if err != nil {
		log.Printf("Error performing search: %v", err)
		http.Error(w, "Error performing search", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Render the search results
	fmt.Fprintf(w, "<h1>Search Results</h1>")
	for rows.Next() {
		var text string
		err = rows.Scan(&text)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "<p>%s</p>", text)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}
}
