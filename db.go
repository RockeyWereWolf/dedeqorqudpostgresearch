import (
    "database/sql"
    "io/ioutil"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func initDB() *sql.DB {
    db, err := sql.Open("postgres", "user=postgres password=example_password host=database dbname=postgres sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    // Read the schema file
    schema, err := ioutil.ReadFile("kitabe-dede-qorqud.sql")
    if err != nil {
        log.Fatal(err)
    }

    // Execute the schema file
    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal(err)
    }

    return db
}
