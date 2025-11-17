package external

import (
	"context"
	"fmt"
	"log"
)

type EmailService interface {
	SendInvestmentConfirmation(ctx context.Context, toEmail, agreementLink, loanDetails string) error
	SendDisbursementNotification(ctx context.Context, toEmail, loanDetails string) error
	SendApprovalNotification(ctx context.Context, toEmail, loanDetails string) error
}

type MockEmailService struct {
	// Track sent emails for testing purposes
	SentEmails []SentEmail
}

type SentEmail struct {
	To      string
	Subject string
	Body    string
}

func NewEmailService() *MockEmailService {
	return &MockEmailService{
		SentEmails: make([]SentEmail, 0),
	}
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		SentEmails: make([]SentEmail, 0),
	}
}

func (m *MockEmailService) SendInvestmentConfirmation(ctx context.Context, toEmail, agreementLink, loanDetails string) error {
	email := SentEmail{
		To:      toEmail,
		Subject: "Investment Confirmation",
		Body:    fmt.Sprintf("Loan invested successfully. Agreement link: %s. Details: %s", agreementLink, loanDetails),
	}
	
	m.SentEmails = append(m.SentEmails, email)
	log.Printf("[MOCK] Sent investment confirmation to %s with agreement link: %s", toEmail, agreementLink)
	
	return nil
}

func (m *MockEmailService) SendDisbursementNotification(ctx context.Context, toEmail, loanDetails string) error {
	email := SentEmail{
		To:      toEmail,
		Subject: "Loan Disbursement Notification",
		Body:    fmt.Sprintf("Loan has been disbursed. Details: %s", loanDetails),
	}
	
	m.SentEmails = append(m.SentEmails, email)
	log.Printf("[MOCK] Sent disbursement notification to %s", toEmail)
	
	return nil
}

func (m *MockEmailService) SendApprovalNotification(ctx context.Context, toEmail, loanDetails string) error {
	email := SentEmail{
		To:      toEmail,
		Subject: "Loan Approval Notification",
		Body:    fmt.Sprintf("Loan has been approved. Details: %s", loanDetails),
	}
	
	m.SentEmails = append(m.SentEmails, email)
	log.Printf("[MOCK] Sent approval notification to %s", toEmail)
	
	return nil
}

// Helper method to get sent emails for testing
func (m *MockEmailService) GetSentEmails() []SentEmail {
	return m.SentEmails
}

// Helper method to clear sent emails
func (m *MockEmailService) ClearSentEmails() {
	m.SentEmails = make([]SentEmail, 0)
}