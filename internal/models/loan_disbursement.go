package models

import "time"

type LoanDisbursement struct {
	ID                          int       `json:"id" db:"id"`
	LoanID                      int       `json:"loan_id" db:"loan_id"`
	FieldOfficerEmployeeID      string    `json:"field_officer_employee_id" db:"field_officer_employee_id"`
	DisbursementDate            time.Time `json:"disbursement_date" db:"disbursement_date"`
	AgreementLetterSignedUrl    string    `json:"agreement_letter_signed_url" db:"agreement_letter_signed_url"`
	CreatedAt                   time.Time `json:"created_at" db:"created_at"`
}