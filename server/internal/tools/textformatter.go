// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// TextFormatterParams represents the parameters for the text formatter
type TextFormatterParams struct {
	Text      string `json:"text"`      // Text to format
	Uppercase bool   `json:"uppercase"` // Convert text to uppercase (if false, converts to lowercase)
}

// No result struct needed - we'll just return the formatted text directly

// ValidateTextFormatterParams validates the parameters for the text formatter
func ValidateTextFormatterParams(params TextFormatterParams) error {
	// Text is required
	if params.Text == "" {
		return ErrMissingRequiredParameter("text")
	}

	return nil
}

// ExecuteTextFormatter executes the text formatter tool with the given HTTP request
// It parses parameters from the request and returns the formatted text
func ExecuteTextFormatter(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseTextFormatterParams(r)
	if err != nil {
		return "", err
	}

	// Validate parameters
	if err := ValidateTextFormatterParams(params); err != nil {
		return "", err
	}

	result := params.Text

	// Apply formatting
	if params.Uppercase {
		result = strings.ToUpper(result)
	} else {
		// Default to lowercase when uppercase is false
		result = strings.ToLower(result)
	}

	return result, nil
}

// ParseTextFormatterParams parses the text formatter parameters from an HTTP request
// It handles both POST and GET requests
func ParseTextFormatterParams(r *http.Request) (TextFormatterParams, error) {
	// Parse parameters from either POST or GET
	var text string
	var uppercaseStr string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				return TextFormatterParams{}, fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			text = r.FormValue("text")
			uppercaseStr = r.FormValue("uppercase")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params TextFormatterParams
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return TextFormatterParams{}, fmt.Errorf("error parsing JSON data: %v", err)
			}
			defer r.Body.Close()
			
			// If JSON values are provided, use them directly
			if params.Text != "" {
				return params, nil
			}
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return TextFormatterParams{}, fmt.Errorf("error parsing form data: %v", err)
			}
			text = r.FormValue("text")
			uppercaseStr = r.FormValue("uppercase")
		}
	} else {
		// Get parameters from query string
		text = r.URL.Query().Get("text")
		uppercaseStr = r.URL.Query().Get("uppercase")
	}

	// Parse boolean parameters
	uppercase := false
	if uppercaseStr == "true" || uppercaseStr == "on" || uppercaseStr == "1" {
		uppercase = true
	}

	// Create and return parameters
	return TextFormatterParams{
		Text:      text,
		Uppercase: uppercase,
	}, nil
}
