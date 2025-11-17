package models

import "time"

type Borrower struct {
	ID                int       `json:"id" db:"id"`
	BorrowerIDNumber  string    `json:"borrower_id_number" db:"borrower_id_number"`
	FullName          string    `json:"full_name" db:"full_name"`
	Email             string    `json:"email" db:"email"`
	Phone             string    `json:"phone" db:"phone"`
	Address           string    `json:"address" db:"address"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}