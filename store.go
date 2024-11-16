package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	CreateExpense(e Expense) (Expense, error)
	ListExpenses(userID int) ([]Expense, error)

	SetUserWeeklyLimit(userID int, newLimit int) error

    GetUserByEmail(email string) (User, error)
}

type PostgresStore struct {
	conn *pgx.Conn
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %v", err)
    }

    return &PostgresStore{conn: conn}, nil
}
func (p *PostgresStore) CreateExpense(e Expense) (Expense, error) {
    query := `
        INSERT INTO expenses (user_id, amount, date, category)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `

	// QueryRow is used to execute the SQL statement and retrieve the id into e.Id.
    err := p.conn.QueryRow(context.Background(), query, e.UserId, e.Amount, e.Date, e.Category).Scan(&e.Id)
    if err != nil {
        return Expense{}, fmt.Errorf("failed to create expense: %v", err)
    }

    return e, nil
}

func (p *PostgresStore) ListExpenses(userID int) ([]Expense, error) {
    query := `
        SELECT id, user_id, amount, date, category
        FROM expenses
        WHERE user_id = $1;
    `

    rows, err := p.conn.Query(context.Background(), query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to list expenses for user %d: %v", userID, err)
    }
    defer rows.Close()

    expenses := []Expense{}
    for rows.Next() {
        var e Expense
        err := rows.Scan(&e.Id, &e.UserId, &e.Amount, &e.Date, &e.Category)
        if err != nil {
            return nil, fmt.Errorf("failed to scan expense: %v", err)
        }
        expenses = append(expenses, e)
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return expenses, nil
}


func (p *PostgresStore) SetUserWeeklyLimit(userID int, newLimit int) error {
    query := `
        UPDATE users
        SET weekly_spending_limit = $1
        WHERE id = $2;
    `

    cmdTag, err := p.conn.Exec(context.Background(), query, newLimit, userID)
    if err != nil {
        return fmt.Errorf("failed to set weekly spending limit: %v", err)
    }

    if cmdTag.RowsAffected() == 0 {
        return fmt.Errorf("no user found with ID %d", userID)
    }

    return nil
}

func (p *PostgresStore) GetUserByEmail(email string) (User, error) {
    var user User
    query := `SELECT id, name, email, password_hash, weekly_spending_limit FROM users WHERE email = $1`
    err := p.conn.QueryRow(context.Background(), query, email).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.WeeklySpendingLimit)
    if err != nil {
        return User{}, fmt.Errorf("user not found: %v", err)
    }
    return user, nil
}