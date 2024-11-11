package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Config struct {
    PostgresPassword string `json:"postgres_password"`
}

func loadConfig() (*Config, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }
    return &config, nil
}

func main() {
    config, err := loadConfig()
    if err != nil {
        log.Fatalf("Could not load config: %v\n", err)
    }

    // connect to database
    connStr := fmt.Sprintf("postgresql://postgres:%s@localhost:5432/Expense_Tracker", config.PostgresPassword)

	store, err := NewPostgresStore(connStr)
    if err != nil {
        log.Fatalf("Unable to initialize store: %v\n", err)
    }
    defer store.conn.Close(context.Background())

    fmt.Println("Successfully connected to the database!")

    // create http handlers and start the server
    h := NewHandler(store)

    mux := chi.NewRouter()
    RegisterRouters(mux,h)

    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", mux)
}