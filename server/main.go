package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
)

var db *sql.DB

func main() {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    migrateDatabase(psqlInfo)

    router := mux.NewRouter()

    router.HandleFunc("/todos", getTodos).Methods("GET")
    router.HandleFunc("/todos", createTodo).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", router))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
    query := "SELECT id, task, done FROM todos"
    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var todos []map[string]interface{}

    for rows.Next() {
        var id int
        var task string
        var done bool
        err := rows.Scan(&id, &task, &done)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        todo := map[string]interface{}{
            "id":   id,
            "task": task,
            "done": done,
        }
        todos = append(todos, todo)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    // ここに、PostgreSQLに新しいtodoを作成するコードを実装
}

func migrateDatabase(psqlInfo string) {
    migrationsDir := "file://db/migrations"
    m, err := migrate.New(migrationsDir, psqlInfo)
    if err != nil {
        log.Fatalf("migration failed: %v", err)
        return
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("migration failed: %v", err)
        return
    }
}
