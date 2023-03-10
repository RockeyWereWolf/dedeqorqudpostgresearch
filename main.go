package main

import (
    "database/sql"
    "time"
    "net/http"
    "fmt"
    "io/ioutil"
    //"log"
    "os"

    _ "github.com/lib/pq"
    log "github.com/sirupsen/logrus"  
)

var db *sql.DB

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
    
    for {
        time.Sleep(time.Second * 10)
    }

    // Read the SQL schema file
    schema, err := ioutil.ReadFile("sql/scheme.sql")
    if err != nil {
        log.Fatal(err)
    }

    // Execute the SQL schema file
    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal(err) 
    } 
    // Define the HTTP handlers
	http.HandleFunc("/", homePage)
	http.HandleFunc("/search", searchHandler)

    // Start the HTTP server
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
	panic(err)
    }
}
//Read and execute sql files

//Initial html page layout
func homePage(w http.ResponseWriter, r *http.Request) {
    // Add the HTML form to get the user's search term
    fmt.Fprintf(w, "<form method='POST' action='/search'>")
    fmt.Fprintf(w, "<input type='text' name='searchTerm'>")
    fmt.Fprintf(w, "<input type='submit' value='Search'>")
    fmt.Fprintf(w, "</form>")
}
//Executing search results
func searchHandler(w http.ResponseWriter, r *http.Request) {
    // Get the search term from the form data
    searchTerm := r.FormValue("searchTerm")

    // Perform the full text search
    results, err := performFullTextSearch(searchTerm)
    if err != nil {
        http.Error(w, "Error performing full text search", http.StatusInternalServerError)
        return
    }

    // Display the results on the web page
    fmt.Fprintf(w, "<h1>Search Results</h1>")
    for _, result := range results {
        fmt.Fprintf(w, "<p>%s</p>", result)
    }
}
//Performing PostgreSQL full text search on the response
func performFullTextSearch(searchTerm string) ([]string, error) {
    // Define the search query
    query := fmt.Sprintf("SELECT text FROM documents WHERE to_tsvector('english', text) @@ to_tsquery('english', '%s')", searchTerm)

    // Execute the search query
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Collect the search results
    var results []string
    for rows.Next() {
        var text string
        err = rows.Scan(&text)
        if err != nil {
            return nil, err
        }
        results = append(results, text)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }

    return results, nil
}
