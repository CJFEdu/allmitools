// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"fmt"
	"net/http"
	"net/url"
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
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return URLEncoderParams{}, fmt.Errorf("error parsing form data: %v", err)
		}

		// Get parameters from form data
		text = r.FormValue("text")
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
