package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
)

func main() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT to_tsvector('english', 'The quick brown fox jumped over the lazy dog')")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        var result string
        for rows.Next() {
            err := rows.Scan(&result)
            if err != nil {
                log.Fatal(err)
            }
        }

        fmt.Fprintln(w, result)
    })

   
