// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"errors"
	"fmt"
	"net/http"

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
		// Parse form data for POST requests
		if err := r.ParseForm(); err != nil {
			return "", fmt.Errorf("failed to parse form data: %w", err)
		}
		id = r.FormValue("id")
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
