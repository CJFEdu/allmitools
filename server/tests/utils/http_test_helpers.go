// Package utils provides testing utilities for the AllMiTools server
package utils

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TemplateTestCase represents a test case for template rendering
type TemplateTestCase struct {
	Name           string
	TemplatePath   string
	TemplateData   interface{}
	ExpectedOutput string
}

// MockResponseWriter is a mock http.ResponseWriter for testing
type MockResponseWriter struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// NewMockResponseWriter creates a new MockResponseWriter
func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		Headers: make(http.Header),
	}
}

// Header returns the header map for the mock response writer
func (m *MockResponseWriter) Header() http.Header {
	return m.Headers
}

// Write writes the data to the mock response writer
func (m *MockResponseWriter) Write(data []byte) (int, error) {
	m.Body = append(m.Body, data...)
	return len(data), nil
}

// WriteHeader sets the status code for the mock response writer
func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

// GetBodyString returns the body as a string
func (m *MockResponseWriter) GetBodyString() string {
	return string(m.Body)
}

// TestTemplateRendering tests template rendering
func TestTemplateRendering(t *testing.T, templatePath string, data interface{}, expectedContains []string) {
	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	assert.NoError(t, err, "Error parsing template: %v", err)

	// Create a buffer to store the rendered template
	w := httptest.NewRecorder()

	// Execute the template
	err = tmpl.Execute(w, data)
	assert.NoError(t, err, "Error executing template: %v", err)

	// Get the result
	result := w.Body.String()

	// Check that the result contains the expected strings
	for _, expected := range expectedContains {
		assert.Contains(t, result, expected, "Template output should contain '%s'", expected)
	}
}

// LoadTestFile loads a file for testing
func LoadTestFile(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Error reading test file: %v", err)
	return data
}

// CreateTempFile creates a temporary file for testing
func CreateTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	assert.NoError(t, err, "Error creating temporary file: %v", err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err, "Error writing to temporary file: %v", err)

	err = tmpFile.Close()
	assert.NoError(t, err, "Error closing temporary file: %v", err)

	return tmpFile.Name()
}

// CleanupTempFile removes a temporary file after testing
func CleanupTempFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err, "Error removing temporary file: %v", err)
}

// GetTestFilePath returns the absolute path to a test file
func GetTestFilePath(t *testing.T, relativePath string) string {
	absPath, err := filepath.Abs(relativePath)
	assert.NoError(t, err, "Error getting absolute path: %v", err)
	return absPath
}

// CompareResponseBody compares the response body with the expected content
func CompareResponseBody(t *testing.T, resp *http.Response, expected string) {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body: %v", err)
	assert.Equal(t, expected, string(body), "Response body does not match expected content")
}

// TestHTMLResponse tests an HTML response
func TestHTMLResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedContains []string) {
	assert.Equal(t, expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"), "Expected Content-Type to be text/html")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body: %v", err)

	for _, expected := range expectedContains {
		assert.Contains(t, string(body), expected, "Response body should contain '%s'", expected)
	}
}

// TestJSONResponse tests a JSON response
func TestJSONResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedJSON interface{}) {
	assert.Equal(t, expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Expected Content-Type to be application/json")

	var result interface{}
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body: %v", err)

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err, "Error unmarshaling JSON response: %v", err)

	expectedData, err := json.Marshal(expectedJSON)
	assert.NoError(t, err, "Error marshaling expected JSON: %v", err)

	var expected interface{}
	err = json.Unmarshal(expectedData, &expected)
	assert.NoError(t, err, "Error unmarshaling expected JSON: %v", err)

	assert.Equal(t, expected, result, "JSON response does not match expected JSON")
}

// ReadResponseBody reads the body of an HTTP response and returns it as a string
func ReadResponseBody(t *testing.T, resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body: %v", err)
	return string(body)
}
