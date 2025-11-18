package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEmailServiceSendInvestmentConfirmation(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	toEmail := "investor@example.com"
	agreementLink := "https://example.com/agreement.pdf"
	loanDetails := "Loan details"

	err := emailService.SendInvestmentConfirmation(ctx, toEmail, agreementLink, loanDetails)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(emailService.SentEmails))
	
	sentEmail := emailService.SentEmails[0]
	assert.Equal(t, toEmail, sentEmail.To)
	assert.Equal(t, "Investment Confirmation", sentEmail.Subject)
	assert.Contains(t, sentEmail.Body, agreementLink)
	assert.Contains(t, sentEmail.Body, loanDetails)
}

func TestMockEmailServiceSendDisbursementNotification(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	toEmail := "borrower@example.com"
	loanDetails := "Loan has been disbursed"

	err := emailService.SendDisbursementNotification(ctx, toEmail, loanDetails)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(emailService.SentEmails))
	
	sentEmail := emailService.SentEmails[0]
	assert.Equal(t, toEmail, sentEmail.To)
	assert.Equal(t, "Loan Disbursement Notification", sentEmail.Subject)
	assert.Contains(t, sentEmail.Body, loanDetails)
}

func TestMockEmailServiceSendApprovalNotification(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	toEmail := "borrower@example.com"
	loanDetails := "Loan has been approved"

	err := emailService.SendApprovalNotification(ctx, toEmail, loanDetails)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(emailService.SentEmails))
	
	sentEmail := emailService.SentEmails[0]
	assert.Equal(t, toEmail, sentEmail.To)
	assert.Equal(t, "Loan Approval Notification", sentEmail.Subject)
	assert.Contains(t, sentEmail.Body, loanDetails)
}

func TestMockEmailServiceMultipleEmails(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	
	// Send investment confirmation
	emailService.SendInvestmentConfirmation(ctx, "investor1@example.com", "link1", "details1")
	
	// Send disbursement notification
	emailService.SendDisbursementNotification(ctx, "borrower@example.com", "details2")
	
	// Send approval notification
	emailService.SendApprovalNotification(ctx, "borrower2@example.com", "details3")

	assert.Equal(t, 3, len(emailService.SentEmails))
	
	// Verify first email
	assert.Equal(t, "investor1@example.com", emailService.SentEmails[0].To)
	assert.Equal(t, "Investment Confirmation", emailService.SentEmails[0].Subject)
	
	// Verify second email
	assert.Equal(t, "borrower@example.com", emailService.SentEmails[1].To)
	assert.Equal(t, "Loan Disbursement Notification", emailService.SentEmails[1].Subject)
	
	// Verify third email
	assert.Equal(t, "borrower2@example.com", emailService.SentEmails[2].To)
	assert.Equal(t, "Loan Approval Notification", emailService.SentEmails[2].Subject)
}

func TestMockEmailServiceGetSentEmails(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	emailService.SendInvestmentConfirmation(ctx, "test@example.com", "link", "details")

	sentEmails := emailService.GetSentEmails()
	assert.Equal(t, 1, len(sentEmails))
	assert.Equal(t, "test@example.com", sentEmails[0].To)
}

func TestMockEmailServiceClearSentEmails(t *testing.T) {
	emailService := NewMockEmailService()

	ctx := context.Background()
	emailService.SendInvestmentConfirmation(ctx, "test@example.com", "link", "details")
	
	assert.Equal(t, 1, len(emailService.SentEmails))
	
	emailService.ClearSentEmails()
	assert.Equal(t, 0, len(emailService.SentEmails))
}