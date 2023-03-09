package main

import (
    "database/sql"
    "time"
    "net/http"
    "fmt"
    //"io/ioutil"
    //"log"
    "os"

    _ "github.com/lib/pq"
    log "github.com/sirupsen/logrus"  
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {
    // Get the database connection parameters from environment variables
    host := os.Getenv("PGHOST")
    port := os.Getenv("PGPORT")
    user := os.Getenv("PGUSER")
    password := os.Getenv("PGPASSWORD")
    dbname := os.Getenv("PGDATABASE")
    
    /*Web app sample testing
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+"8080", mux) */
    
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
    schema, err := ioutil.ReadFile("kitabe-dede-qorqud.sql")
    if err != nil {
        log.Fatal(err)
    }

    // Execute the SQL schema file
    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal(err) 
    }
}
