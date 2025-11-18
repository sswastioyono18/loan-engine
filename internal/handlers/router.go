package handlers

import (
	"net/http"
	"time"

	"github.com/sswastioyono18/loan-engine/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(serviceFactory *services.ServiceFactory) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize handlers
	authHandler := NewAuthHandler(serviceFactory.AuthService())
	borrowerHandler := NewBorrowerHandler(serviceFactory.BorrowerService())
	loanHandler := NewLoanHandler(
		serviceFactory.LoanService(),
		serviceFactory.EmailService,
		serviceFactory.StorageService,
	)
	investorHandler := NewInvestorHandler(serviceFactory.InvestorService())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Authentication routes
		r.Post("/auth/register", authHandler.RegisterUser)
		r.Post("/auth/login", authHandler.LoginUser)
		r.Post("/auth/refresh", authHandler.RefreshToken)

		// Borrower routes
		r.Post("/borrowers", borrowerHandler.CreateBorrower)
		r.Get("/borrowers/{id}", borrowerHandler.GetBorrowerByID)
		r.Put("/borrowers/{id}", borrowerHandler.UpdateBorrower)
		r.Delete("/borrowers/{id}", borrowerHandler.DeleteBorrower)
		r.Get("/borrowers", borrowerHandler.ListBorrowers)

		// Investor routes
		r.Post("/investors", investorHandler.CreateInvestor)
		r.Get("/investors/{id}", investorHandler.GetInvestorByID)
		r.Put("/investors/{id}", investorHandler.UpdateInvestor)
		r.Delete("/investors/{id}", investorHandler.DeleteInvestor)
		r.Get("/investors", investorHandler.ListInvestors)

		// Loan routes
		r.Post("/loans", loanHandler.CreateLoan)
		r.Get("/loans/{id}", loanHandler.GetLoanByID)
		r.Put("/loans/{id}", loanHandler.UpdateLoan)
		r.Delete("/loans/{id}", loanHandler.DeleteLoan)
		r.Get("/loans", loanHandler.ListLoans)
		r.Get("/loans/state/{state}", loanHandler.GetLoansByState)

		// Loan state transition routes
		r.Post("/loans/{id}/approve", loanHandler.ApproveLoan)
		r.Post("/loans/{id}/invest", loanHandler.InvestInLoan)
		r.Post("/loans/{id}/disburse", loanHandler.DisburseLoan)
	})

	return router
}
