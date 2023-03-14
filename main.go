package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	//"math/rand"
        "strconv"

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
	
        // Read the SQL file containing the sample data
	data, err := ioutil.ReadFile("sample.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Split the file contents into individual SQL statements
	statements := strings.Split(string(data), ";")

	// Execute each SQL statement to insert data into the documents table
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			log.Fatal(err)
		}
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
			SELECT id, title, main_character, content, ts_headline(content, q, 'StartSel = <mark>, StopSel = </mark>, MaxWords = 100, MinWords = 10, ShortWord = 3, HighlightAll = true') AS snippet
			FROM books, to_tsquery($1) AS q
			WHERE to_tsvector('english', content) @@ q
		`, query)
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			log.Error(err)
			return
		}
		defer rows.Close()

		func getSentenceSuffix(n int) string {
			if n > 3 || n < 1 {
				return "th"
			}
			switch n {
			case 1:
				return "st"
			case 2:
				return "nd"
			case 3:
				return "rd"
			}
			return ""
		}

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
			var title, main_character, content, snippet string
			if err := rows.Scan(&id, &title, &main_character, &content, &snippet); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				log.Error(err)
				return
			}
			sentences := strings.Split(snippet, ".")
			var validSentences []string
			for i, sentence := range sentences {
				if strings.Contains(strings.ToLower(sentence), strings.ToLower(query)) {
					validSentences = append(validSentences, fmt.Sprintf("%d%s", i+1, getSentenceSuffix(i+1)))
				}
			}
			numSentences := strconv.Itoa(len(validSentences))
			snippet = strings.Join(validSentences, ". ")

			fmt.Fprintf(w, `
				<div>
					<h3>Book %dst</h3>
					<p><em>%s</em></p>
					<p>Found in %s sentence(s): %s</p>
					<p>%s</p>
				</div>
			`, id, title, numSentences, strings.Join(validSentences, ", "), snippet)
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
