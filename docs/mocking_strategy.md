# Mocking Strategy for External Services

## Overview
This document outlines the mocking strategy for external services in the loan engine system, including email services, file storage, and any third-party APIs that need to be mocked during development and testing.

## External Services to Mock

### 1. Email Service
The loan engine needs to send emails to investors when loans are fully invested, containing links to agreement letters.

### 2. File Storage Service
The system needs to store and retrieve proof images, signed agreement letters, and other documents.

### 3. Payment Gateway (Future)
For processing loan disbursements and repayments (if applicable).

## Mocking Approach

### 1. Interface-Based Mocking

#### Email Service Interface
```go
// pkg/external/email_service.go
package external

import "context"

type EmailService interface {
    SendInvestmentConfirmation(ctx context.Context, toEmail, agreementLink, loanDetails string) error
    SendDisbursementNotification(ctx context.Context, toEmail, loanDetails string) error
    SendApprovalNotification(ctx context.Context, toEmail, loanDetails string) error
}

// pkg/external/email_service_impl.go
package external

import (
    "context"
    "fmt"
    "log"
)

type SMTPEmailService struct {
    host     string
    port     int
    username string
    password string
}

func NewSMTPEmailService(host string, port int, username, password string) *SMTPEmailService {
    return &SMTPEmailService{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

func (s *SMTPEmailService) SendInvestmentConfirmation(ctx context.Context, toEmail, agreementLink, loanDetails string) error {
    // Implementation for sending real emails
    log.Printf("Sending investment confirmation to %s with agreement link: %s", toEmail, agreementLink)
    return nil
}

func (s *SMTPEmailService) SendDisbursementNotification(ctx context.Context, toEmail, loanDetails string) error {
    // Implementation for sending real emails
    log.Printf("Sending disbursement notification to %s", toEmail)
    return nil
}

func (s *SMTPEmailService) SendApprovalNotification(ctx context.Context, toEmail, loanDetails string) error {
    // Implementation for sending real emails
    log.Printf("Sending approval notification to %s", toEmail)
    return nil
}
```

#### Mock Email Service
```go
// pkg/external/email_service_mock.go
package external

import (
    "context"
    "fmt"
    "log"
)

type MockEmailService struct {
    // Track sent emails for testing purposes
    SentEmails []SentEmail
}

type SentEmail struct {
    To      string
    Subject string
    Body    string
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
        Body:    fmt.Sprintf("Loan invested successfully. Agreement link: %s", agreementLink),
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
```

### 2. File Storage Service

#### File Storage Interface
```go
// pkg/external/storage_service.go
package external

import (
    "context"
    "io"
)

type StorageService interface {
    UploadFile(ctx context.Context, file io.Reader, fileName, contentType string) (string, error)
    DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error)
    DeleteFile(ctx context.Context, fileID string) error
    GetFileURL(ctx context.Context, fileID string) (string, error)
}

// pkg/external/storage_service_impl.go
package external

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

type LocalStorageService struct {
    storagePath string
}

func NewLocalStorageService(storagePath string) *LocalStorageService {
    // Create storage directory if it doesn't exist
    os.MkdirAll(storagePath, 0755)
    return &LocalStorageService{
        storagePath: storagePath,
    }
}

func (s *LocalStorageService) UploadFile(ctx context.Context, file io.Reader, fileName, contentType string) (string, error) {
    filePath := filepath.Join(s.storagePath, fileName)
    
    outFile, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer outFile.Close()
    
    _, err = io.Copy(outFile, file)
    if err != nil {
        return "", err
    }
    
    return filePath, nil
}

func (s *LocalStorageService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
    return os.Open(fileID)
}

func (s *LocalStorageService) DeleteFile(ctx context.Context, fileID string) error {
    return os.Remove(fileID)
}

func (s *LocalStorageService) GetFileURL(ctx context.Context, fileID string) (string, error) {
    return fmt.Sprintf("/files/%s", filepath.Base(fileID)), nil
}
```

#### Mock File Storage Service
```go
// pkg/external/storage_service_mock.go
package external

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "log"
)

type MockStorageService struct {
    storedFiles map[string]*bytes.Buffer
    fileURLs    map[string]string
}

func NewMockStorageService() *MockStorageService {
    return &MockStorageService{
        storedFiles: make(map[string]*bytes.Buffer),
        fileURLs:    make(map[string]string),
    }
}

func (m *MockStorageService) UploadFile(ctx context.Context, file io.Reader, fileName, contentType string) (string, error) {
    // Read the file content
    content, err := io.ReadAll(file)
    if err != nil {
        return "", err
    }
    
    // Store in memory
    buffer := bytes.NewBuffer(content)
    m.storedFiles[fileName] = buffer
    
    // Generate a mock URL
    mockURL := fmt.Sprintf("http://mock-storage/%s", fileName)
    m.fileURLs[fileName] = mockURL
    
    log.Printf("[MOCK] Uploaded file: %s with content type: %s", fileName, contentType)
    
    return mockURL, nil
}

func (m *MockStorageService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
    content, exists := m.storedFiles[fileID]
    if !exists {
        return nil, fmt.Errorf("file not found: %s", fileID)
    }
    
    log.Printf("[MOCK] Downloaded file: %s", fileID)
    
    return io.NopCloser(bytes.NewReader(content.Bytes())), nil
}

func (m *MockStorageService) DeleteFile(ctx context.Context, fileID string) error {
    delete(m.storedFiles, fileID)
    delete(m.fileURLs, fileID)
    
    log.Printf("[MOCK] Deleted file: %s", fileID)
    
    return nil
}

func (m *MockStorageService) GetFileURL(ctx context.Context, fileID string) (string, error) {
    url, exists := m.fileURLs[fileID]
    if !exists {
        return "", fmt.Errorf("file URL not found: %s", fileID)
    }
    
    return url, nil
}

// Helper method to get stored files for testing
func (m *MockStorageService) GetStoredFiles() map[string]*bytes.Buffer {
    return m.storedFiles
}

// Helper method to check if file exists
func (m *MockStorageService) FileExists(fileName string) bool {
    _, exists := m.storedFiles[fileName]
    return exists
}
```

### 3. Mock Server for Testing

#### Mock HTTP Server
```go
// pkg/external/mock_server.go
package external

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type MockServer struct {
    emailService   *MockEmailService
    storageService *MockStorageService
    server         *http.Server
}

func NewMockServer(emailService *MockEmailService, storageService *MockStorageService) *MockServer {
    return &MockServer{
        emailService:   emailService,
        storageService: storageService,
    }
}

func (m *MockServer) Start(port int) error {
    mux := http.NewServeMux()
    
    // Email service endpoints
    mux.HandleFunc("/email/investment-confirmation", m.handleInvestmentConfirmation)
    mux.HandleFunc("/email/disbursement-notification", m.handleDisbursementNotification)
    mux.HandleFunc("/email/approval-notification", m.handleApprovalNotification)
    
    // Storage service endpoints
    mux.HandleFunc("/storage/upload", m.handleUpload)
    mux.HandleFunc("/storage/download/", m.handleDownload)
    mux.HandleFunc("/storage/delete/", m.handleDelete)
    
    // Health check
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    m.server = &http.Server{
        Addr:         fmt.Sprintf(":%d", port),
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    log.Printf("Starting mock server on port %d", port)
    return m.server.ListenAndServe()
}

func (m *MockServer) Stop() error {
    return m.server.Close()
}

func (m *MockServer) handleInvestmentConfirmation(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        ToEmail       string `json:"to_email"`
        AgreementLink string `json:"agreement_link"`
        LoanDetails   string `json:"loan_details"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    err := m.emailService.SendInvestmentConfirmation(r.Context(), req.ToEmail, req.AgreementLink, req.LoanDetails)
    if err != nil {
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (m *MockServer) handleDisbursementNotification(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        ToEmail     string `json:"to_email"`
        LoanDetails string `json:"loan_details"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    err := m.emailService.SendDisbursementNotification(r.Context(), req.ToEmail, req.LoanDetails)
    if err != nil {
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (m *MockServer) handleApprovalNotification(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        ToEmail     string `json:"to_email"`
        LoanDetails string `json:"loan_details"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    err := m.emailService.SendApprovalNotification(r.Context(), req.ToEmail, req.LoanDetails)
    if err != nil {
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (m *MockServer) handleUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    fileName := r.FormValue("filename")
    contentType := r.FormValue("content_type")
    
    if fileName == "" {
        http.Error(w, "Filename is required", http.StatusBadRequest)
        return
    }
    
    url, err := m.storageService.UploadFile(r.Context(), file, fileName, contentType)
    if err != nil {
        http.Error(w, "Failed to upload file", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (m *MockServer) handleDownload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    fileID := r.URL.Path[len("/storage/download/"):]
    if fileID == "" {
        http.Error(w, "File ID is required", http.StatusBadRequest)
        return
    }
    
    file, err := m.storageService.DownloadFile(r.Context(), fileID)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }
    defer file.Close()
    
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileID))
    
    io.Copy(w, file)
}

func (m *MockServer) handleDelete(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    fileID := r.URL.Path[len("/storage/delete/"):]
    if fileID == "" {
        http.Error(w, "File ID is required", http.StatusBadRequest)
        return
    }
    
    err := m.storageService.DeleteFile(r.Context(), fileID)
    if err != nil {
        http.Error(w, "Failed to delete file", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
```

## Dependency Injection for Mocking

### Configuration for Different Environments
```go
// internal/config/config.go
package config

import (
    "os"
    "strings"
)

type Config struct {
    Environment string
    EmailService EmailServiceConfig
    StorageService StorageServiceConfig
}

type EmailServiceConfig struct {
    Provider string // "smtp" or "mock"
    SMTPHost string
    SMTPPort int
    SMTPUsername string
    SMTPPassword string
}

type StorageServiceConfig struct {
    Provider string // "local", "s3", or "mock"
    LocalPath string
    S3Bucket string
    S3Region string
}

func LoadConfig() *Config {
    env := getEnv("ENV", "development")
    
    config := &Config{
        Environment: env,
        EmailService: EmailServiceConfig{
            Provider: getEnv("EMAIL_PROVIDER", "mock"),
            SMTPHost: getEnv("SMTP_HOST", ""),
            SMTPPort: getEnvAsInt("SMTP_PORT", 587),
            SMTPUsername: getEnv("SMTP_USERNAME", ""),
            SMTPPassword: getEnv("SMTP_PASSWORD", ""),
        },
        StorageService: StorageServiceConfig{
            Provider: getEnv("STORAGE_PROVIDER", "mock"),
            LocalPath: getEnv("LOCAL_STORAGE_PATH", "./storage"),
            S3Bucket: getEnv("S3_BUCKET", ""),
            S3Region: getEnv("S3_REGION", ""),
        },
    }
    
    return config
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        // In a real implementation, you would parse the string to int
        // For simplicity, returning default value here
        return defaultValue
    }
    return defaultValue
}
```

### Service Factory with Mock Support
```go
// internal/services/factory.go (updated)
package services

import (
    "loan-engine/internal/config"
    "loan-engine/internal/repositories"
    "loan-engine/pkg/external"
)

type ServiceFactory struct {
    RepoFactory  *repositories.RepositoryFactory
    EmailService external.EmailService
    StorageService external.StorageService
}

func NewServiceFactory(
    repoFactory *repositories.RepositoryFactory,
    config *config.Config,
) *ServiceFactory {
    var emailService external.EmailService
    var storageService external.StorageService
    
    // Initialize email service based on config
    if strings.ToLower(config.EmailService.Provider) == "mock" {
        emailService = external.NewMockEmailService()
    } else {
        emailService = external.NewSMTPEmailService(
            config.EmailService.SMTPHost,
            config.EmailService.SMTPPort,
            config.EmailService.SMTPUsername,
            config.EmailService.SMTPPassword,
        )
    }
    
    // Initialize storage service based on config
    if strings.ToLower(config.StorageService.Provider) == "mock" {
        storageService = external.NewMockStorageService()
    } else if strings.ToLower(config.StorageService.Provider) == "local" {
        storageService = external.NewLocalStorageService(config.StorageService.LocalPath)
    } else {
        // Add other storage providers as needed (S3, etc.)
        storageService = external.NewMockStorageService() // fallback
    }
    
    return &ServiceFactory{
        RepoFactory:  repoFactory,
        EmailService: emailService,
        StorageService: storageService,
    }
}

// Updated service methods to use external services
func (f *ServiceFactory) LoanService() LoanService {
    return NewLoanService(f.RepoFactory, f.EmailService, f.StorageService)
}
```

## Testing with Mocks

### Example Test with Mocks
```go
// internal/services/loan_service_test.go
package services

import (
    "context"
    "testing"
    "loan-engine/internal/models"
    "loan-engine/pkg/external"
    "github.com/stretchr/testify/assert"
)

func TestLoanInvestmentSendsEmail(t *testing.T) {
    // Setup
    mockEmailService := external.NewMockEmailService()
    mockStorageService := external.NewMockStorageService()
    
    // Create service factory with mocks
    serviceFactory := &ServiceFactory{
        RepoFactory:  createMockRepoFactory(), // Implementation not shown
        EmailService: mockEmailService,
        StorageService: mockStorageService,
    }
    
    loanService := NewLoanService(
        serviceFactory.RepoFactory,
        serviceFactory.EmailService,
        serviceFactory.StorageService,
    )
    
    // Create test data
    loan := &models.Loan{
        ID: 1,
        LoanID: "LOAN001",
        PrincipalAmount: 10000.00,
        CurrentState: "approved",
    }
    
    investment := &models.LoanInvestment{
        InvestorID: 1,
        InvestmentAmount: 10000.00,
    }
    
    // Execute
    err := loanService.InvestInLoan(context.Background(), loan.ID, investment)
    
    // Verify
    assert.NoError(t, err)
    
    // Check that email was sent
    sentEmails := mockEmailService.GetSentEmails()
    assert.Len(t, sentEmails, 1)
    assert.Equal(t, "Investment Confirmation", sentEmails[0].Subject)
}
```

## Docker Compose for Mock Services

```yaml
# docker-compose.mock.yml
version: '3.8'

services:
  # Mock email service
  mock-email-service:
    build:
      context: ..
      dockerfile: Dockerfile.mock-email
    container_name: mock_email_service
    ports:
      - "3001:3000"
    environment:
      - PORT=3000
    networks:
      - loan_engine_network

  # Mock storage service
  mock-storage-service:
    build:
      context: ..
      dockerfile: Dockerfile.mock-storage
    container_name: mock_storage_service
    ports:
      - "3002:3000"
    environment:
      - PORT=3000
    networks:
      - loan_engine_network

networks:
  loan_engine_network:
    external: true  # Use the network from main compose file
```

This mocking strategy allows for:
1. Complete isolation during unit testing
2. Fast test execution without external dependencies
3. Predictable test results
4. Easy verification of external service interactions
5. Flexibility to switch between real and mock services based on environment