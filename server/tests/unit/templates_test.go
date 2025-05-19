// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CJFEdu/allmitools/server/internal/handlers"
	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestTemplateInitialization tests the initialization of the template manager
func TestTemplateInitialization(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Verify that the template manager is initialized
	assert.NotNil(t, templates.TemplateManager, "Template manager should not be nil")
}

// TestHomeHandlerTemplateRendering tests the template rendering in the HomeHandler
func TestHomeHandlerTemplateRendering(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err, "Error creating request: %v", err)
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Call the handler
	handlers.HomeHandler(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"), "Handler returned wrong content type")
	
	// Check that the response contains expected HTML elements
	body := rr.Body.String()
	assert.Contains(t, body, "Home - AllMiTools", "Response should contain the title")
	assert.Contains(t, body, "Welcome to AllMiTools", "Response should contain the welcome message")
	assert.Contains(t, body, "Available Tools", "Response should contain the tools section")
}

// TestDocsBaseHandlerTemplateRendering tests the template rendering in the DocsBaseHandler
func TestDocsBaseHandlerTemplateRendering(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/docs", nil)
	assert.NoError(t, err, "Error creating request: %v", err)
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Call the handler
	handlers.DocsBaseHandler(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"), "Handler returned wrong content type")
	
	// Check that the response contains expected HTML elements
	body := rr.Body.String()
	assert.Contains(t, body, "Documentation - AllMiTools", "Response should contain the title")
	assert.Contains(t, body, "Available Tools", "Response should contain the tools section")
}

// TestDocsToolHandlerTemplateRendering tests the template rendering in the DocsToolHandler
func TestDocsToolHandlerTemplateRendering(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Get a valid tool name for testing
	tools := models.ListTools()
	if len(tools) == 0 {
		t.Skip("No tools available for testing")
	}
	
	toolName := tools[0].Name
	
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/docs/"+toolName, nil)
	assert.NoError(t, err, "Error creating request: %v", err)
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Set up the router context with the tool name
	router := mux.NewRouter()
	router.HandleFunc("/docs/{tool_name}", handlers.DocsToolHandler)
	
	// Serve the request
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"), "Handler returned wrong content type")
	
	// Check that the response contains expected HTML elements
	body := rr.Body.String()
	assert.Contains(t, body, toolName, "Response should contain the tool name")
	assert.Contains(t, body, "Documentation", "Response should contain the documentation title")
}

// TestToolsHandlerTemplateRendering tests the template rendering in the ToolsHandler
func TestToolsHandlerTemplateRendering(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Get a valid tool name for testing
	tools := models.ListTools()
	if len(tools) == 0 {
		t.Skip("No tools available for testing")
	}
	
	toolName := tools[0].Name
	
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/tools/"+toolName, nil)
	assert.NoError(t, err, "Error creating request: %v", err)
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Set up the router context with the tool name
	router := mux.NewRouter()
	router.HandleFunc("/tools/{tool_name}", handlers.ToolsHandler)
	
	// Serve the request
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"), "Handler returned wrong content type")
	
	// Check that the response contains expected HTML elements
	body := rr.Body.String()
	assert.Contains(t, body, toolName, "Response should contain the tool name")
	assert.Contains(t, body, "form", "Response should contain a form for the tool")
}

// TestNotFoundHandlerTemplateRendering tests the template rendering in the NotFoundHandler
func TestNotFoundHandlerTemplateRendering(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	assert.NoError(t, err, "Error creating request: %v", err)
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Call the handler
	handlers.NotFoundHandler(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code, "Handler returned wrong status code")
	
	// Check the content type
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"), "Handler returned wrong content type")
	
	// Check that the response contains expected HTML elements
	body := rr.Body.String()
	assert.Contains(t, body, "Page Not Found", "Response should contain the not found message")
	assert.Contains(t, body, "404 - Page Not Found", "Response should contain the not found message")
}

// TestContentTypeNegotiation tests content type negotiation in handlers
func TestContentTypeNegotiation(t *testing.T) {
	// Initialize the template manager
	err := templates.Initialize("../../")
	assert.NoError(t, err, "Template initialization should not return an error")
	
	// Test cases for different handlers and accept headers
	testCases := []struct {
		name         string
		path         string
		acceptHeader string
		contentType  string
	}{
		{"HomeHandlerHTML", "/", "text/html", "text/html; charset=utf-8"},
		{"HomeHandlerJSON", "/", "application/json", "application/json"},
		{"DocsBaseHandlerHTML", "/docs", "text/html", "text/html; charset=utf-8"},
		{"DocsBaseHandlerJSON", "/docs", "application/json", "application/json"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified accept header
			req, err := http.NewRequest("GET", tc.path, nil)
			assert.NoError(t, err, "Error creating request: %v", err)
			req.Header.Set("Accept", tc.acceptHeader)
			
			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			
			// Create a router and register the handlers
			router := mux.NewRouter()
			router.HandleFunc("/", handlers.HomeHandler)
			router.HandleFunc("/docs", handlers.DocsBaseHandler)
			router.HandleFunc("/docs/", handlers.DocsBaseHandler)
			
			// Serve the request
			router.ServeHTTP(rr, req)
			
			// Check the status code
			assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
			
			// Check the content type
			assert.Equal(t, tc.contentType, rr.Header().Get("Content-Type"), "Handler returned wrong content type")
		})
	}
}
