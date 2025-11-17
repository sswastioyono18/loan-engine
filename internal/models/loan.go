package models

import "time"

type Loan struct {
	ID                  int       `json:"id" db:"id"`
	LoanID              string    `json:"loan_id" db:"loan_id"`
	BorrowerID          int       `json:"borrower_id" db:"borrower_id"`
	PrincipalAmount     float64   `json:"principal_amount" db:"principal_amount"`
	Rate                float64   `json:"rate" db:"rate"` // Interest rate percentage
	ROI                 float64   `json:"roi" db:"roi"`   // Return of investment percentage
	AgreementLetterLink string    `json:"agreement_letter_link" db:"agreement_letter_link"`
	CurrentState        string    `json:"current_state" db:"current_state"`
	TotalInvestedAmount float64   `json:"total_invested_amount" db:"total_invested_amount"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}