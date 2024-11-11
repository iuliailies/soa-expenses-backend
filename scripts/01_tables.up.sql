-- Drop tables if they already exist to avoid conflicts when running the script multiple times
DROP TABLE IF EXISTS EXPENSES;
DROP TABLE IF EXISTS USERS;

-- Create the USERS table
CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,           -- Unique identifier for each user
    name VARCHAR(100) NOT NULL,      -- Name of the user
    email VARCHAR(100) UNIQUE NOT NULL,  -- Email of the user, must be unique
    password_hash VARCHAR(255) NOT NULL,   -- Hashed password for secure storage
    weekly_spending_limit DECIMAL(10, 2)  -- Weekly spending limit set by the user
);

-- Create the EXPENSES table
CREATE TABLE EXPENSES (
    id SERIAL PRIMARY KEY,           -- Unique identifier for each expense
    user_id INT NOT NULL,            -- Foreign key to associate expense with a user
    amount DECIMAL(10, 2) NOT NULL,  -- Amount of the expense
    date DATE NOT NULL,              -- Date of the expense
    category VARCHAR(50),            -- Category of the expense (e.g., Food, Travel)
    
    -- Define foreign key constraint to link expenses to the USERS table
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES USERS(id)
        ON DELETE CASCADE
);

-- Indexes for faster lookup
-- CREATE INDEX idx_user_email ON USERS(email);
-- CREATE INDEX idx_expense_user_id ON EXPENSES(user_id);
