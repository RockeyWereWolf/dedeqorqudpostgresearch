package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func main() {
    // Get the database connection parameters from environment variables
    host := os.Getenv("PGHOST")
    port := os.Getenv("PGPORT")
    user := os.Getenv("PGUSER")
    password := os.Getenv("PGPASSWORD")
    dbname := os.Getenv("PGDATABASE")
    sslmode := os.Getenv("PGSSLMODE")

    // Create a connection string using the parameters
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)

    // Open a connection to the database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

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
