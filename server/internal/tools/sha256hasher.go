// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

// SHA256HasherParams represents the parameters for the SHA-256 hasher
type SHA256HasherParams struct {
	Text string `json:"text"` // Text to hash
}

// ValidateSHA256HasherParams validates the parameters for the SHA-256 hasher
func ValidateSHA256HasherParams(params SHA256HasherParams) error {
	if params.Text == "" {
		return fmt.Errorf("text parameter is required")
	}
	return nil
}

// ParseSHA256HasherParams parses the SHA-256 hasher parameters from an HTTP request
// It handles both POST and GET requests
func ParseSHA256HasherParams(r *http.Request) (SHA256HasherParams, error) {
	// Parse parameters from either POST or GET
	var text string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return SHA256HasherParams{}, fmt.Errorf("error parsing form data: %v", err)
		}

		// Get parameters from form data
		text = r.FormValue("text")
	} else {
		// Get parameters from query string
		text = r.URL.Query().Get("text")
	}

	// Validate parameters
	if text == "" {
		return SHA256HasherParams{}, fmt.Errorf("text parameter is required")
	}

	// Create and return parameters
	return SHA256HasherParams{
		Text: text,
	}, nil
}

// ExecuteSHA256Hasher executes the SHA-256 hasher with the given HTTP request
// It parses parameters from the request and returns the SHA-256 hash as a string
func ExecuteSHA256Hasher(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseSHA256HasherParams(r)
	if err != nil {
		return "", err
	}

	// Validate parameters
	if err := ValidateSHA256HasherParams(params); err != nil {
		return "", err
	}

	// Create a new SHA-256 hash
	hasher := sha256.New()
	
	// Write the text to the hasher
	hasher.Write([]byte(params.Text))
	
	// Get the hash sum as bytes
	hashBytes := hasher.Sum(nil)
	
	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	// Return the hash string
	return hashString, nil
}
