package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sswastioyono18/loan-engine/internal/handlers"
	"github.com/sswastioyono18/loan-engine/internal/repositories"
	"github.com/sswastioyono18/loan-engine/internal/services"
	"github.com/sswastioyono18/loan-engine/pkg/external"
	"github.com/sswastioyono18/loan-engine/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := util.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	borrowerRepo := repositories.NewBorrowerRepository(db)
	loanRepo := repositories.NewLoanRepository(db)
	loanApprovalRepo := repositories.NewLoanApprovalRepository(db)
	loanDisbursementRepo := repositories.NewLoanDisbursementRepository(db)
	investorRepo := repositories.NewInvestorRepository(db)
	loanInvestmentRepo := repositories.NewLoanInvestmentRepository(db)
	loanStateHistoryRepo := repositories.NewLoanStateHistoryRepository(db)
	_ = repositories.NewUserRepository(db) // Initialize for potential future use

	// Initialize external services (mocks for now)
	emailService := external.NewEmailService()
	storageService := external.NewStorageService()

	// Initialize services
	borrowerService := services.NewBorrowerService(borrowerRepo)
	loanService := services.NewLoanService(loanRepo, loanApprovalRepo, loanDisbursementRepo, loanInvestmentRepo, loanStateHistoryRepo, investorRepo, emailService, storageService)
	investorService := services.NewInvestorService(investorRepo)

	// Initialize handlers
	borrowerHandler := handlers.NewBorrowerHandler(borrowerService)
	loanHandler := handlers.NewLoanHandler(loanService, emailService, storageService)
	investorHandler := handlers.NewInvestorHandler(investorService)

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Borrower routes
		r.Post("/borrowers", borrowerHandler.CreateBorrower)
		r.Get("/borrowers/{id}", borrowerHandler.GetBorrowerByID)
		r.Put("/borrowers/{id}", borrowerHandler.UpdateBorrower)
		r.Delete("/borrowers/{id}", borrowerHandler.DeleteBorrower)
		r.Get("/borrowers", borrowerHandler.ListBorrowers)

		// Loan routes
		r.Post("/loans", loanHandler.CreateLoan)
		r.Get("/loans/{id}", loanHandler.GetLoanByID)
		r.Put("/loans/{id}", loanHandler.UpdateLoan)
		r.Delete("/loans/{id}", loanHandler.DeleteLoan)
		r.Get("/loans", loanHandler.ListLoans)
		r.Post("/loans/{id}/approve", loanHandler.ApproveLoan)
		r.Post("/loans/{id}/invest", loanHandler.InvestInLoan)
		r.Post("/loans/{id}/disburse", loanHandler.DisburseLoan)
		r.Get("/loans/state/{state}", loanHandler.GetLoansByState)

		// Investor routes
		r.Post("/investors", investorHandler.CreateInvestor)
		r.Get("/investors/{id}", investorHandler.GetInvestorByID)
		r.Put("/investors/{id}", investorHandler.UpdateInvestor)
		r.Delete("/investors/{id}", investorHandler.DeleteInvestor)
		r.Get("/investors", investorHandler.ListInvestors)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
