package models

import (
	"time"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleUser   UserRole = "user"
	RoleViewer UserRole = "viewer"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose in JSON
	Email        string    `json:"email,omitempty" db:"email"`
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string   `json:"username" validate:"required"`
	Password string   `json:"password" validate:"required,min=8"`
	Email    string   `json:"email,omitempty" validate:"omitempty,email"`
	Role     UserRole `json:"role" validate:"required,oneof=admin user viewer"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email    *string   `json:"email,omitempty" validate:"omitempty,email"`
	Password *string   `json:"password,omitempty" validate:"omitempty,min=8"`
	Role     *UserRole `json:"role,omitempty" validate:"omitempty,oneof=admin user viewer"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response with JWT token
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      *User     `json:"user"`
}

