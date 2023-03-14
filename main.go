package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	//"math/rand"
        //"strconv"

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
						<h1>Kitabe Dede Qorqud Search</h1>
						<p>This project allows you to search any keywords(eg "Bayindir") in Kitabe Dede Qorqud</p>
						<form method="GET">
							<label for="q">Type your keyword:</label>
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

        // Count the number of sentences in the content
        sentences := strings.Split(content, ".")
        numSentences := len(sentences)

        // Find the sentences containing the search query
        var matches []int
        for i, sentence := range sentences {
	        if strings.Contains(sentence, query) {
		    matches = append(matches, i)
	        }
        }

        // Display the search result for the current book
        fmt.Fprintf(w, `
	        <div>
		        <h3>Book %d</h3>
		            <p><em>%s</em></p>
		            <p>Found in %d sentence(s)</p>
                `, id, title, len(matches))
            for _, match := range matches {
	            sentence := sentences[match]
	            // Highlight the search query in the sentence
	            sentence = strings.ReplaceAll(sentence, query, "<mark>"+query+"</mark>")
	            fmt.Fprintf(w, `
		            <p>%d. %s</p>
	            `, match+1, sentence)
            }
            fmt.Fprintf(w, `
	            <p>Out of %d sentence(s)</p>
	            <p>%s</p>
            </div>
        `, numSentences, snippet)
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
