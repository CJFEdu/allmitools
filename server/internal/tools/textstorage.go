// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/database"
)

// ExecuteTextStorage executes the text storage tool
// This tool stores text content in the database and returns a unique ID
// Parameters:
//   - content: The text content to store (required)
//   - save: Whether to save the text permanently (optional, default: false)
func ExecuteTextStorage(r *http.Request) (string, error) {
	// Parse parameters
	content := ""
	saveFlag := false
	var err error

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
			content = r.FormValue("content")
			saveFlagStr := r.FormValue("save")
			if saveFlagStr != "" {
				saveFlag, err = strconv.ParseBool(saveFlagStr)
				if err != nil {
					return "", fmt.Errorf("invalid value for 'save' parameter: %w", err)
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params struct {
				Content string `json:"content"`
				Save    bool   `json:"save"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return "", fmt.Errorf("failed to parse JSON data: %w", err)
			}
			defer r.Body.Close()
			
			// Use the values from JSON
			content = params.Content
			saveFlag = params.Save
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return "", fmt.Errorf("failed to parse form data: %w", err)
			}
			content = r.FormValue("content")
			saveFlagStr := r.FormValue("save")
			if saveFlagStr != "" {
				saveFlag, err = strconv.ParseBool(saveFlagStr)
				if err != nil {
					return "", fmt.Errorf("invalid value for 'save' parameter: %w", err)
				}
			}
		}
	} else {
		// Parse query parameters for GET requests
		content = r.URL.Query().Get("content")
		saveFlagStr := r.URL.Query().Get("save")
		if saveFlagStr != "" {
			saveFlag, err = strconv.ParseBool(saveFlagStr)
			if err != nil {
				return "", fmt.Errorf("invalid value for 'save' parameter: %w", err)
			}
		}
	}

	// Validate parameters
	if content == "" {
		return "", errors.New("content parameter is required")
	}

	// Get the DAO
	dao, err := database.GetTextStorageDAO()
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	// Store the text
	id, err := dao.StoreText(content, saveFlag)
	if err != nil {
		return "", fmt.Errorf("failed to store text: %w", err)
	}

	// Return the ID
	return fmt.Sprintf("Text stored successfully with ID: %s", id), nil
}
