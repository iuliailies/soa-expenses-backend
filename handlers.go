package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler struct to encapsulate HTTP handling logic
type Handler struct {
    store Store
}

func NewHandler(store Store) *Handler {
    return &Handler{store: store}
}

func RegisterRouters(mux *chi.Mux, handler *Handler)  {

    mux.Use(middleware.Logger) // Add logging middleware

    mux.Post("/expenses", handler.CreateExpense)
    mux.Get("/expenses", handler.ListExpenses)
    mux.Put("/users/{userID}/limit", handler.SetUserWeeklyLimit)

}

// Handler function to create a new expense
func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
    var expense Expense
    if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    createdExpense, err := h.store.CreateExpense(expense)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create expense: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createdExpense)
}

// Handler function to list all expenses
func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
    expenses, err := h.store.ListExpenses()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to list expenses: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(expenses)
}


// Handler function to update a user's weekly spending limit
func (h *Handler) SetUserWeeklyLimit(w http.ResponseWriter, r *http.Request) {
    userIDStr := chi.URLParam(r, "userID")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
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
