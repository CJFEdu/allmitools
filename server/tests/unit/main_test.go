// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CJFEdu/allmitools/server/internal/handlers"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/CJFEdu/allmitools/server/tests/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestMain is the entry point for running tests
func TestMain(m *testing.M) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	if err != nil {
		// If template initialization fails, print an error and exit
		fmt.Printf("Error initializing templates: %v\n", err)
		os.Exit(1)
	}
	
	// Run the tests
	os.Exit(m.Run())
}

// TestHomeHandler tests the home handler function
func TestHomeHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Create a handler function from the HomeHandler function
	handler := http.HandlerFunc(handlers.HomeHandler)

	// Serve the HTTP request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	
	// Check the response body
	assert.Contains(t, rr.Body.String(), "Welcome to AllMiTools", "handler returned unexpected body")
}

// TestDocsBaseHandler tests the docs base handler function
func TestDocsBaseHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/docs", nil)
	assert.NoError(t, err)
	
	// Set the Accept header to application/json to get a JSON response
	req.Header.Set("Accept", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Create a handler function from the DocsBaseHandler function
	handler := http.HandlerFunc(handlers.DocsBaseHandler)

	// Serve the HTTP request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "handler returned wrong content type")
	
	// Check the response body
	assert.Contains(t, rr.Body.String(), "\"success\":true", "handler returned unexpected body")
}

// TestDocsToolHandler tests the docs tool handler function
func TestDocsToolHandler(t *testing.T) {
	// Create a new router to use the gorilla/mux vars
	r := mux.NewRouter()
	r.HandleFunc("/docs/{tool_name}", handlers.DocsToolHandler)
	
	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()
	
	// Create a client
	client := &http.Client{}
	
	// Create a request with JSON Accept header
	req, err := http.NewRequest("GET", ts.URL + "/docs/random-number", nil)
	assert.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	
	// Make the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "handler returned wrong content type")
	
	// Check the response body
	body := utils.ReadResponseBody(t, resp)
	assert.Contains(t, body, "\"success\":true", "handler returned unexpected body")
}

// TestToolsHandler tests the tools handler function
func TestToolsHandler(t *testing.T) {
	// Create a new router to use the gorilla/mux vars
	r := mux.NewRouter()
	r.HandleFunc("/tools/{tool_name}", handlers.ToolsHandler)
	
	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()
	
	// Create a client
	client := &http.Client{}
	
	// Create a request with JSON Accept header
	req, err := http.NewRequest("GET", ts.URL + "/tools/random-number", nil)
	assert.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	
	// Make the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "handler returned wrong content type")
	
	// Check the response body
	body := utils.ReadResponseBody(t, resp)
	assert.Contains(t, body, "\"success\":true", "handler returned unexpected body")
}

// TestNotFoundHandler tests the not found handler function
func TestNotFoundHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Create a handler function from the NotFoundHandler function
	handler := http.HandlerFunc(handlers.NotFoundHandler)

	// Serve the HTTP request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
	
	// Check the response body
	assert.Contains(t, rr.Body.String(), "404 - Page Not Found", "handler returned unexpected body")
}

// TestNewRouter tests the router configuration
func TestNewRouter(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()
	
	// Configure routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/docs", handlers.DocsBaseHandler).Methods("GET")
	router.HandleFunc("/docs/", handlers.DocsBaseHandler).Methods("GET")
	router.HandleFunc("/docs/{tool_name}", handlers.DocsToolHandler).Methods("GET")
	router.HandleFunc("/tools/{tool_name}", handlers.ToolsHandler).Methods("GET", "POST")
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	
	// Test cases for the router
	testCases := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{"Homepage GET", "GET", "/", http.StatusOK, "Welcome to AllMiTools"},
		{"Docs base path", "GET", "/docs", http.StatusOK, "Documentation - AllMiTools"},
		{"Docs base path with trailing slash", "GET", "/docs/", http.StatusOK, "Documentation - AllMiTools"},
		{"Docs for specific tool", "GET", "/docs/random-number", http.StatusOK, "random-number"},
		{"Tool GET request", "GET", "/tools/random-number", http.StatusOK, "random-number"},
		{"Non-existent path", "GET", "/nonexistent", http.StatusNotFound, "404 - Page Not Found"},
	}
	
	// Create a test server using the router
	ts := httptest.NewServer(router)
	defer ts.Close()
	
	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the request
			var resp *http.Response
			var err error
			
			switch tc.method {
			case "GET":
				resp, err = http.Get(ts.URL + tc.url)
			case "POST":
				resp, err = http.Post(ts.URL + tc.url, "application/json", nil)
			default:
				t.Fatalf("Unsupported method: %s", tc.method)
			}
			
			assert.NoError(t, err)
			defer resp.Body.Close()
			
			// Check the status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected status code %d for %s %s, got %d", 
				tc.expectedStatus, tc.method, tc.url, resp.StatusCode)
			
			// Check the response body
			body := utils.ReadResponseBody(t, resp)
			assert.Contains(t, body, tc.expectedBody, "Expected body to contain '%s' for %s %s, got '%s'", 
				tc.expectedBody, tc.method, tc.url, body)
		})
	}
}
