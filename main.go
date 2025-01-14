package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
    connStr := fmt.Sprintf("postgresql://postgres:%s@host.docker.internal:5432/Expense_Tracker", config.PostgresPassword)
    // connStr := fmt.Sprintf("postgresql://postgres:%s@localhost:5432/Expense_Tracker", config.PostgresPassword)

	store, err := NewPostgresStore(connStr)
    if err != nil {
        log.Fatalf("Unable to initialize store: %v\n", err)
    }
    defer store.conn.Close()

    fmt.Println("Successfully connected to the database!")

    rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

    time.Sleep(10 * time.Second)

	publisher, err := NewRabbitMQPublisher(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()
    log.Println("successfully connected to rabbitmq!")

    // create http handlers and start the server
    h := NewHandler(store, publisher)

    mux := chi.NewRouter()
    RegisterRouters(mux,h)

    fmt.Println("Starting server on :8081")
    http.ListenAndServe("0.0.0.0:8081", mux)
}
