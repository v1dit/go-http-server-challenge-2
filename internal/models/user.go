package models

import "time"

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // never send password in JSON
	CreatedAt time.Time `json:"created_at"`
}

// Session represents an authenticated user session.
type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}
