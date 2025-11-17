package external

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
)

type StorageService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName, contentType string) (string, error)
	DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, fileID string) error
	GetFileURL(ctx context.Context, fileID string) (string, error)
}

type MockStorageService struct {
	storedFiles map[string]*bytes.Buffer
	fileURLs    map[string]string
}

func NewStorageService() *MockStorageService {
	return &MockStorageService{
		storedFiles: make(map[string]*bytes.Buffer),
		fileURLs:    make(map[string]string),
	}
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