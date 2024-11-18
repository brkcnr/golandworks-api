package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/service"
	"github.com/brkcnr/golandworks-api/internal/transport/httpserver"
)

func TestNew(t *testing.T) {
	// Initialize dependencies
	todoSvc := &service.TodoService{}
	
	// Create new server
	server := httpserver.New(todoSvc)
	
	// Test cases for different endpoints
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /todo should return OK",
			method:         "GET",
			path:          "/todo",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /todo should return OK",
			method:         "POST",
			path:          "/todo",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /search should return OK",
			method:         "GET",
			path:          "/search",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid path should return NotFound",
			method:         "GET",
			path:          "/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Serve the request
			server.ServeHTTP(w, req)
			
			// Check status code
			if got := w.Code; got != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, got)
			}
		})
	}
} 