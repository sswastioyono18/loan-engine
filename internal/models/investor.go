package models

import "time"

type Investor struct {
	ID        int       `json:"id" db:"id"`
	InvestorID string   `json:"investor_id" db:"investor_id"`
	FullName  string    `json:"full_name" db:"full_name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}