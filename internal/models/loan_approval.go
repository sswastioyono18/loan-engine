package models

import "time"

type LoanApproval struct {
	ID                       int       `json:"id" db:"id"`
	LoanID                   int       `json:"loan_id" db:"loan_id"`
	FieldValidatorEmployeeID string    `json:"field_validator_employee_id" db:"field_validator_employee_id"`
	ApprovalDate             time.Time `json:"approval_date" db:"approved_at"`
	ProofImageUrl            string    `json:"proof_image_url" db:"proof_image_url"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
}