package models

import "time"

type User struct {
	ID          int       `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	UserType    string    `json:"user_type" db:"user_type"` // staff, investor, admin
	FullName    string    `json:"full_name" db:"name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}