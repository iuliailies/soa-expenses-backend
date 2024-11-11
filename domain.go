package main

import (
	"time"
)

type Expense struct {
    Id       int       `json:"id"`
    UserId   int       `json:"user_id"`
    Amount   int       `json:"amount"`
    Date     time.Time `json:"date"`
    Category string    `json:"category"`
}

type User struct {
    Id                  int    `json:"id"`
    Name                string `json:"name"`
    Email               string `json:"email"`
    PasswordHash        string `json:"-"`
    WeeklySpendingLimit int    `json:"weekly_spending_limit"`
}