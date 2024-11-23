package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/crypto/bcrypt"
)

// Handler struct to encapsulate HTTP handling logic
type Handler struct {
    store Store
}

func NewHandler(store Store) *Handler {
    return &Handler{store: store}
}

func RegisterRouters(mux *chi.Mux, handler *Handler) {
    mux.Use(middleware.Logger) // Add logging middleware

    mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true, 
		MaxAge:           300,  
	}))

    mux.Route("/api", func(api chi.Router) {
        api.Post("/expenses", handler.CreateExpense)
        api.Get("/expenses", handler.ListExpenses)
        api.Put("/users/limit", handler.SetUserWeeklyLimit)
    })

    mux.Post("/auth/login", handler.AuthenticateUser)
}

// Handler function to create a new expense
func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from the request header
    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "Unauthorized: missing or invalid user_id", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user_id format", http.StatusBadRequest)
        return
    }

    // Decode the expense from the request body
    var expense Expense
    if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Assign the extracted user ID to the expense
    expense.UserId = userID

    // Create the expense in the database
    createdExpense, err := h.store.CreateExpense(expense)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create expense: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with the created expense
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createdExpense)
}

// Handler function to list all expenses
func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "Unauthorized: missing or invalid user_id", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user_id format", http.StatusBadRequest)
        return
    }

	// Fetch expenses for the user
	expenses, err := h.store.ListExpenses(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list expenses: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// Handler function to update a user's weekly spending limit
func (h *Handler) SetUserWeeklyLimit(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "Unauthorized: missing or invalid user_id", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user_id format", http.StatusBadRequest)
        return
    }

    var payload struct {
        NewLimit int `json:"new_limit"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err = h.store.SetUserWeeklyLimit(userID, payload.NewLimit)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to set weekly spending limit: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content if successful
}

// AuthenticateUser handler to validate email and password, then return user ID
func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    user, err := h.store.GetUserByEmail(credentials.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Return user ID if authentication is successful
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"user_id": user.Id})
}
