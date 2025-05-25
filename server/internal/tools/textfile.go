// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// TextFileParams represents the parameters for the text file tool
type TextFileParams struct {
	Content string `json:"content"` // Content to be saved as a text file
	Filename string `json:"filename"` // Optional filename for the text file (default: "download.txt")
}

// ValidateTextFileParams validates the parameters for the text file tool
func ValidateTextFileParams(params TextFileParams) error {
	if params.Content == "" {
		return errors.New("content cannot be empty")
	}
	
	return nil
}

// GenerateTextFile prepares the content to be downloaded as a text file
// It returns the content and filename for the text file
func GenerateTextFile(params TextFileParams) (string, string, error) {
	// Validate parameters
	if err := ValidateTextFileParams(params); err != nil {
		return "", "", err
	}

	// Use default filename if not provided
	filename := params.Filename
	if filename == "" {
		filename = "download.txt"
	}

	return params.Content, filename, nil
}

// ExecuteTextFile handles the text file tool which generates a downloadable text file
// from the provided content. It parses the request, generates the text file, and writes
// the response with appropriate headers for file download.
//
// Parameters are extracted from the request based on the Content-Type:
// - For form submissions: content and filename form fields
// - For JSON submissions: Content and Filename JSON properties
// - For GET requests: content and filename query parameters
//
// Returns the file content, filename, and response headers for the caller to use.
func ExecuteTextFile(r *http.Request) (string, string, error) {
	// Parse parameters from either POST or GET
	var content, filename string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				return "", "", fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			content = r.FormValue("content")
			filename = r.FormValue("filename")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params TextFileParams
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return "", "", fmt.Errorf("error parsing JSON data: %v", err)
			}
			defer r.Body.Close()
			
			// Use the values from JSON
			content = params.Content
			filename = params.Filename
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return "", "", fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			content = r.FormValue("content")
			filename = r.FormValue("filename")
		}
	} else {
		// Get parameters from query string
		content = r.URL.Query().Get("content")
		filename = r.URL.Query().Get("filename")
	}

	// Validate parameters
	if content == "" {
		return "", "", fmt.Errorf("content parameter is required")
	}

	// Create parameters for the text file tool
	params := TextFileParams{
		Content:  content,
		Filename: filename,
	}

	// Generate text file
	fileContent, fileName, err := GenerateTextFile(params)
	if err != nil {
		return "", "", err
	}

	return fileContent, fileName, nil
}
