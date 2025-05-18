// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CJFEdu/allmitools/server/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestRouterConfiguration tests that the router is configured correctly
func TestRouterConfiguration(t *testing.T) {
	// Create a new router for testing
	r := mux.NewRouter()
	
	// Add routes to match the ones in main.go
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/docs", handlers.DocsBaseHandler).Methods("GET")
	r.HandleFunc("/docs/", handlers.DocsBaseHandler).Methods("GET")
	r.HandleFunc("/docs/{tool_name}", handlers.DocsToolHandler).Methods("GET")
	r.HandleFunc("/tools/{tool_name}", handlers.ToolsHandler).Methods("GET", "POST")
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	
	// Test cases for route matching
	testCases := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		shouldMatch    bool
	}{
		{"Homepage GET", "GET", "/", http.StatusOK, true},
		{"Homepage POST (not allowed)", "POST", "/", http.StatusMethodNotAllowed, false},
		{"Docs base path", "GET", "/docs", http.StatusOK, true},
		{"Docs base path with trailing slash", "GET", "/docs/", http.StatusOK, true},
		{"Docs for specific tool", "GET", "/docs/random-number", http.StatusOK, true},
		{"Tool GET request", "GET", "/tools/random-number", http.StatusOK, true},
		{"Tool POST request", "POST", "/tools/random-number", http.StatusOK, true},
		// The NotFoundHandler will handle non-existent paths, so all paths will match the router
		{"Non-existent path", "GET", "/nonexistent", http.StatusNotFound, true},
	}
	
	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			assert.NoError(t, err)
			
			// Use a recorder to capture the response
			rr := httptest.NewRecorder()
			
			// For the special case of non-existent paths, we're testing the NotFoundHandler
			if tc.name == "Non-existent path" {
				// Just serve the request and check the status code
				r.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code %d for %s %s, got %d", 
					tc.expectedStatus, tc.method, tc.url, rr.Code)
				return
			}
			
			// For other cases, check route matching first
			var match mux.RouteMatch
			matched := r.Match(req, &match)
			
			// Assert that the route matched as expected
			assert.Equal(t, tc.shouldMatch, matched, "Route matching should be %v for %s %s", tc.shouldMatch, tc.method, tc.url)
			
			// If the route should match, test the handler
			if tc.shouldMatch && matched {
				r.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code %d for %s %s, got %d", 
					tc.expectedStatus, tc.method, tc.url, rr.Code)
			}
		})
	}
}


