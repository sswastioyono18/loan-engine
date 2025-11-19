package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sswastioyono18/loan-engine/internal/handlers"
	"github.com/sswastioyono18/loan-engine/internal/repositories"
	"github.com/sswastioyono18/loan-engine/internal/services"
	"github.com/sswastioyono18/loan-engine/pkg/external"
	"github.com/sswastioyono18/loan-engine/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLoanE2EScenario(t *testing.T) {
	ctx := context.Background()

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "loan_engine_db",
				"POSTGRES_USER":     "loan_engine_user",
				"POSTGRES_PASSWORD": "loan_engine_password",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer postgresC.Terminate(ctx)

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "loan_engine_user")
	os.Setenv("DB_PASSWORD", "loan_engine_password")
	os.Setenv("DB_NAME", "loan_engine_db")
	os.Setenv("DB_SSL_MODE", "disable")

	time.Sleep(2 * time.Second)

	db, err := util.InitDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, goose.SetDialect("postgres"))
	require.NoError(t, goose.Up(db.DB, "./migrations"))

	router := setupE2ERouter(db)

	// Step 1: Create Borrower
	borrowerResp := postJSON(t, router, "/api/v1/borrowers", map[string]interface{}{
		"id_number": "B001",
		"name":      "John Doe",
		"email":     "john.doe@example.com",
		"phone":     "+621234567890",
		"address":   "Jalan Tedeng Aling Aling",
	})
	borrowerID := int(borrowerResp["data"].(map[string]interface{})["id"].(float64))
	fmt.Printf("✅ Step 1: Borrower created (ID: %d)\n", borrowerID)

	// Step 2: Create Loan (State: proposed)
	loanResp := postJSON(t, router, "/api/v1/loans", map[string]interface{}{
		"borrower_id":           borrowerID,
		"principal_amount":      1000000.00,
		"rate":                  0.05,
		"roi":                   0.08,
		"agreement_letter_link": "https://example.com/agreement.pdf",
	})
	loanData := loanResp["data"].(map[string]interface{})
	loanID := int(loanData["id"].(float64))
	assert.Equal(t, "proposed", loanData["current_state"])
	fmt.Printf("✅ Step 2: Loan created (ID: %d, State: %s)\n", loanID, loanData["current_state"])

	// Step 3: Approve Loan (State: proposed → approved)
	approveResp := postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/approve", loanID), map[string]interface{}{
		"field_validator_employee_id": "emp001",
		"proof_image_url":             "https://example.com/proof.jpg",
	})
	assert.True(t, approveResp["success"].(bool))
	
	loanResp = getJSON(t, router, fmt.Sprintf("/api/v1/loans/%d", loanID))
	loanData = loanResp["data"].(map[string]interface{})
	assert.Equal(t, "approved", loanData["current_state"])
	fmt.Printf("✅ Step 3: Loan approved (State: %s)\n", loanData["current_state"])

	// Step 4: Create Investor
	investorResp := postJSON(t, router, "/api/v1/investors", map[string]interface{}{
		"investor_id": "INV001",
		"name":        "Jane Smith",
		"email":       "jane.smith@example.com",
		"phone":       "+0987654321",
	})
	investorID := int(investorResp["data"].(map[string]interface{})["id"].(float64))
	fmt.Printf("✅ Step 4: Investor created (ID: %d)\n", investorID)

	// Step 5: Invest in Loan (State: approved → invested)
	investResp := postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/invest", loanID), map[string]interface{}{
		"investor_id":       investorID,
		"investment_amount": 1000000.00,
	})
	assert.True(t, investResp["success"].(bool))
	
	loanResp = getJSON(t, router, fmt.Sprintf("/api/v1/loans/%d", loanID))
	loanData = loanResp["data"].(map[string]interface{})
	assert.Equal(t, "invested", loanData["current_state"])
	assert.Equal(t, 1000000.00, loanData["total_invested_amount"])
	fmt.Printf("✅ Step 5: Loan invested (State: %s, Amount: %.2f)\n", loanData["current_state"], loanData["total_invested_amount"])

	// Step 6: Disburse Loan (State: invested → disbursed)
	disburseResp := postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/disburse", loanID), map[string]interface{}{
		"field_officer_employee_id":   "emp002",
		"agreement_letter_signed_url": "https://example.com/signed-agreement.pdf",
	})
	assert.True(t, disburseResp["success"].(bool))
	
	loanResp = getJSON(t, router, fmt.Sprintf("/api/v1/loans/%d", loanID))
	loanData = loanResp["data"].(map[string]interface{})
	assert.Equal(t, "disbursed", loanData["current_state"])
	fmt.Printf("✅ Step 6: Loan disbursed (State: %s)\n", loanData["current_state"])

	fmt.Println("\n🎉 E2E Test Complete: Loan lifecycle from proposed → approved → invested → disbursed")
}

func TestLoanPartialInvestmentScenario(t *testing.T) {
	ctx := context.Background()

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "loan_engine_db",
				"POSTGRES_USER":     "loan_engine_user",
				"POSTGRES_PASSWORD": "loan_engine_password",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer postgresC.Terminate(ctx)

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "loan_engine_user")
	os.Setenv("DB_PASSWORD", "loan_engine_password")
	os.Setenv("DB_NAME", "loan_engine_db")
	os.Setenv("DB_SSL_MODE", "disable")

	time.Sleep(2 * time.Second)

	db, err := util.InitDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, goose.SetDialect("postgres"))
	require.NoError(t, goose.Up(db.DB, "./migrations"))

	router := setupE2ERouter(db)

	// Create Borrower
	borrowerResp := postJSON(t, router, "/api/v1/borrowers", map[string]interface{}{
		"id_number": "B002",
		"name":      "Jane Doe",
		"email":     "jane.doe@example.com",
		"phone":     "+621234567891",
		"address":   "Jalan Sudirman",
	})
	borrowerID := int(borrowerResp["data"].(map[string]interface{})["id"].(float64))
	fmt.Printf("✅ Borrower created (ID: %d)\n", borrowerID)

	// Create Loan
	loanResp := postJSON(t, router, "/api/v1/loans", map[string]interface{}{
		"borrower_id":           borrowerID,
		"principal_amount":      5000000.00,
		"rate":                  0.05,
		"roi":                   0.08,
		"agreement_letter_link": "https://example.com/agreement.pdf",
	})
	loanData := loanResp["data"].(map[string]interface{})
	loanID := int(loanData["id"].(float64))
	assert.Equal(t, "proposed", loanData["current_state"])
	fmt.Printf("✅ Loan created (ID: %d, Principal: %.2f, State: %s)\n", loanID, loanData["principal_amount"], loanData["current_state"])

	// Approve Loan
	postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/approve", loanID), map[string]interface{}{
		"field_validator_employee_id": "emp001",
		"proof_image_url":             "https://example.com/proof.jpg",
	})
	fmt.Printf("✅ Loan approved\n")

	// Create Investor 1
	investor1Resp := postJSON(t, router, "/api/v1/investors", map[string]interface{}{
		"investor_id": "INV001",
		"name":        "Investor One",
		"email":       "investor1@example.com",
		"phone":       "+0987654321",
	})
	investor1ID := int(investor1Resp["data"].(map[string]interface{})["id"].(float64))

	// Create Investor 2
	investor2Resp := postJSON(t, router, "/api/v1/investors", map[string]interface{}{
		"investor_id": "INV002",
		"name":        "Investor Two",
		"email":       "investor2@example.com",
		"phone":       "+0987654322",
	})
	investor2ID := int(investor2Resp["data"].(map[string]interface{})["id"].(float64))
	fmt.Printf("✅ Investors created (ID: %d, %d)\n", investor1ID, investor2ID)

	// Partial Investment 1 (2M out of 5M)
	postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/invest", loanID), map[string]interface{}{
		"investor_id":       investor1ID,
		"investment_amount": 2000000.00,
	})
	
	loanResp = getJSON(t, router, fmt.Sprintf("/api/v1/loans/%d", loanID))
	loanData = loanResp["data"].(map[string]interface{})
	assert.Equal(t, "approved", loanData["current_state"])
	assert.Equal(t, 2000000.00, loanData["total_invested_amount"])
	fmt.Printf("✅ Partial investment 1: %.2f (State: %s, Total: %.2f/%.2f)\n", 
		2000000.00, loanData["current_state"], loanData["total_invested_amount"], loanData["principal_amount"])

	// Partial Investment 2 (3M out of 5M - completes the loan)
	postJSON(t, router, fmt.Sprintf("/api/v1/loans/%d/invest", loanID), map[string]interface{}{
		"investor_id":       investor2ID,
		"investment_amount": 3000000.00,
	})
	
	loanResp = getJSON(t, router, fmt.Sprintf("/api/v1/loans/%d", loanID))
	loanData = loanResp["data"].(map[string]interface{})
	assert.Equal(t, "invested", loanData["current_state"])
	assert.Equal(t, 5000000.00, loanData["total_invested_amount"])
	fmt.Printf("✅ Partial investment 2: %.2f (State: %s, Total: %.2f/%.2f)\n", 
		3000000.00, loanData["current_state"], loanData["total_invested_amount"], loanData["principal_amount"])

	fmt.Println("\n🎉 Partial Investment Test Complete: Loan fully funded by multiple investors")
}

func setupE2ERouter(db *util.DB) *chi.Mux {
	borrowerRepo := repositories.NewBorrowerRepository(db)
	loanRepo := repositories.NewLoanRepository(db)
	loanApprovalRepo := repositories.NewLoanApprovalRepository(db)
	loanDisbursementRepo := repositories.NewLoanDisbursementRepository(db)
	investorRepo := repositories.NewInvestorRepository(db)
	loanInvestmentRepo := repositories.NewLoanInvestmentRepository(db)
	loanStateHistoryRepo := repositories.NewLoanStateHistoryRepository(db)

	emailService := external.NewEmailService()
	storageService := external.NewStorageService()

	borrowerService := services.NewBorrowerService(borrowerRepo)
	loanService := services.NewLoanService(loanRepo, loanApprovalRepo, loanDisbursementRepo, loanInvestmentRepo, loanStateHistoryRepo, investorRepo, emailService, storageService)
	investorService := services.NewInvestorService(investorRepo)

	borrowerHandler := handlers.NewBorrowerHandler(borrowerService)
	loanHandler := handlers.NewLoanHandler(loanService, emailService, storageService)
	investorHandler := handlers.NewInvestorHandler(investorService)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/borrowers", borrowerHandler.CreateBorrower)
		r.Get("/borrowers/{id}", borrowerHandler.GetBorrowerByID)
		
		r.Post("/loans", loanHandler.CreateLoan)
		r.Get("/loans/{id}", loanHandler.GetLoanByID)
		r.Post("/loans/{id}/approve", loanHandler.ApproveLoan)
		r.Post("/loans/{id}/invest", loanHandler.InvestInLoan)
		r.Post("/loans/{id}/disburse", loanHandler.DisburseLoan)
		
		r.Post("/investors", investorHandler.CreateInvestor)
	})

	return r
}

func postJSON(t *testing.T, router *chi.Mux, path string, payload map[string]interface{}) map[string]interface{} {
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	require.Equal(t, http.StatusOK, w.Code, "Response: %s", w.Body.String())
	
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	return response
}

func getJSON(t *testing.T, router *chi.Mux, path string) map[string]interface{} {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	require.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	return response
}
