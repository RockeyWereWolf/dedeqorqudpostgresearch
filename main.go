package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

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
	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the search query from the request parameters
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" {
			// Display an empty search form if no query was provided
			fmt.Fprint(w, `
				<html>
					<head><title>Search</title></head>
					<body>
						<h1>Search</h1>
						<form method="GET">
							<label for="q">Search query:</label>
							<input type="text" name="q" id="q">
							<button type="submit">Search</button>
						</form>
					</body>
				</html>
			`)
			return
		}

		// Execute the full-text search query
		rows, err := db.Query(`
			SELECT id, body, ts_headline(body, q) AS snippet
			FROM documents, to_tsquery($1) AS q
			WHERE to_tsvector('english', body) @@ q
		`, query)
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			log.Error(err)
			return
		}
		defer rows.Close()

		// Display the search results in an HTML page
		fmt.Fprintf(w, `
			<html>
				<head><title>Search Results</title></head>
				<body>
					<h1>Search Results</h1>
					<p>Showing results for: %q</p>
		`, query)
		for rows.Next() {
			var id int
			var body, snippet string
			if err := rows.Scan(&id, &body, &snippet); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				log.Error(err)
				return
			}
			fmt.Fprintf(w, `
				<div>
					<h3>Document #%d</h3>
					<p>%s</p>
					<p><em>%s</em></p>
				</div>
			`, id, snippet, body)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to iterate over results", http.StatusInternalServerError)
			log.Error(err)
			return
		}
		fmt.Fprintf(w, `
				</body>
			</html>
		`)
	})
	// Start the HTTP server and listen for incoming requests
	addr := ":8080"
	log.Infof("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

