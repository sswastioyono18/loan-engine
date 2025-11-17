package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kitabisa/loan-engine/internal/handlers"
	"github.com/kitabisa/loan-engine/internal/repositories"
	"github.com/kitabisa/loan-engine/internal/services"
	"github.com/kitabisa/loan-engine/pkg/external"
)

func main() {
	// Load configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "loan_engine_user")
	dbPassword := getEnv("DB_PASSWORD", "loan_engine_password")
	dbName := getEnv("DB_NAME", "loan_engine_db")
	dbSslMode := getEnv("DB_SSL_MODE", "disable")
	jwtSecret := getEnv("JWT_SECRET", "your_jwt_secret_key_here")
	
	// Build connection string
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)

	// Initialize database connection
	db, err := repositories.NewPostgreSQLDriver(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize external services (mocks for now)
	emailService := external.NewMockEmailService()
	storageService := external.NewMockStorageService()

	// Initialize service factory
	serviceFactory := services.NewServiceFactory(
		repositories.NewRepositoryFactory(db),
		emailService,
		storageService,
		jwtSecret,
	)

	// Create router
	router := handlers.NewRouter(serviceFactory)

	// Get port from environment or use default
	port := getEnv("PORT", "8080")
	
	// Start server
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}