package models

import "time"

type LoanInvestment struct {
	ID               int       `json:"id" db:"id"`
	LoanID           int       `json:"loan_id" db:"loan_id"`
	InvestorID       int       `json:"investor_id" db:"investor_id"`
	InvestmentAmount float64   `json:"investment_amount" db:"investment_amount"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}