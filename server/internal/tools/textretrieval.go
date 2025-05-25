// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/database"
)

// ExecuteTextRetrieval executes the text retrieval tool
// This tool retrieves text content from the database using a unique ID
// Parameters:
//   - id: The unique ID of the text to retrieve (required)
func ExecuteTextRetrieval(r *http.Request) (string, error) {
	// Parse parameters
	id := ""

	// Handle both GET and POST requests
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse form data for POST requests
			if err := r.ParseForm(); err != nil {
				return "", fmt.Errorf("failed to parse form data: %w", err)
			}
			id = r.FormValue("id")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params struct {
				ID string `json:"id"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return "", fmt.Errorf("failed to parse JSON data: %w", err)
			}
			defer r.Body.Close()
			
			// Use the ID from JSON
			id = params.ID
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return "", fmt.Errorf("failed to parse form data: %w", err)
			}
			id = r.FormValue("id")
		}
	} else {
		// Parse query parameters for GET requests
		id = r.URL.Query().Get("id")
	}

	// Validate parameters
	if id == "" {
		return "", errors.New("id parameter is required")
	}

	// Get the DAO
	dao, err := database.GetTextStorageDAO()
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	// Retrieve the text
	entry, err := dao.GetTextByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text: %w", err)
	}

	// Return the content
	return entry.Content, nil
}
