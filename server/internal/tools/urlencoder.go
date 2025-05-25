// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// URLEncoderParams represents the parameters for the URL encoder
type URLEncoderParams struct {
	Text string `json:"text"` // Text to encode
}

// ValidateURLEncoderParams validates the parameters for the URL encoder
func ValidateURLEncoderParams(params URLEncoderParams) error {
	if params.Text == "" {
		return fmt.Errorf("text parameter is required")
	}
	return nil
}

// ParseURLEncoderParams parses the URL encoder parameters from an HTTP request
// It handles both POST and GET requests
func ParseURLEncoderParams(r *http.Request) (URLEncoderParams, error) {
	// Parse parameters from either POST or GET
	var text string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				return URLEncoderParams{}, fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			text = r.FormValue("text")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params URLEncoderParams
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return URLEncoderParams{}, fmt.Errorf("error parsing JSON data: %v", err)
			}
			defer r.Body.Close()
			
			// Extract text from parsed JSON
			text = params.Text
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return URLEncoderParams{}, fmt.Errorf("error parsing form data: %v", err)
			}
			text = r.FormValue("text")
		}
	} else {
		// Get parameters from query string
		text = r.URL.Query().Get("text")
	}

	// Validate parameters
	if text == "" {
		return URLEncoderParams{}, fmt.Errorf("text parameter is required")
	}

	// Create and return parameters
	return URLEncoderParams{
		Text: text,
	}, nil
}

// ExecuteURLEncoder executes the URL encoder with the given HTTP request
// It parses parameters from the request and returns the URL-encoded string
func ExecuteURLEncoder(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseURLEncoderParams(r)
	if err != nil {
		return "", err
	}

	// Validate parameters
	if err := ValidateURLEncoderParams(params); err != nil {
		return "", err
	}

	// URL encode the text
	encodedText := url.QueryEscape(params.Text)

	// Return the encoded text
	return encodedText, nil
}
