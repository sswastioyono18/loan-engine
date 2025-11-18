package models

import "time"

type LoanStateHistory struct {
	ID               int       `json:"id" db:"id"`
	LoanID           int       `json:"loan_id" db:"loan_id"`
	PreviousState    string    `json:"previous_state" db:"old_state"`
	NewState         string    `json:"new_state" db:"new_state"`
	TransitionReason string    `json:"transition_reason" db:"reason"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}