package main

import (
    "database/sql"
    "fmt"
    //"io/ioutil"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func main() {
    // Get the database connection parameters from environment variables
    host := "dedeqorqudpostgresearch-db-1"
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

    /* Read the SQL schema file
    schema, err := ioutil.ReadFile("kitabe-dede-qorqud.sql")
    if err != nil {
        log.Fatal(err)
    }

    // Execute the SQL schema file
    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal(err) 
    }*/
}
