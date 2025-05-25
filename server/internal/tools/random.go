// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RandomNumberParams represents the parameters for the random number generator
type RandomNumberParams struct {
	Min int `json:"min"` // Minimum value (inclusive)
	Max int `json:"max"` // Maximum value (inclusive)
}

// ValidateRandomNumberParams validates the parameters for the random number generator
func ValidateRandomNumberParams(params RandomNumberParams) error {
	if params.Min > params.Max {
		return errors.New("minimum value cannot be greater than maximum value")
	}
	return nil
}

// GenerateRandomNumber generates a random number within the specified range
// It returns the generated number and any error that occurred
func GenerateRandomNumber(params RandomNumberParams) (int, error) {
	// Validate parameters
	if err := ValidateRandomNumberParams(params); err != nil {
		return 0, err
	}

	// Initialize random number generator with current time as seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random number within range
	// Add 1 to make the range inclusive of max
	return r.Intn(params.Max-params.Min+1) + params.Min, nil
}

// ParseRandomNumberParams parses the random number generator parameters from an HTTP request
// It handles both POST and GET requests
func ParseRandomNumberParams(r *http.Request) (RandomNumberParams, error) {
	// Parse parameters from either POST or GET
	var minStr, maxStr string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				return RandomNumberParams{}, fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			minStr = r.FormValue("min")
			maxStr = r.FormValue("max")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params RandomNumberParams
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return RandomNumberParams{}, fmt.Errorf("error parsing JSON data: %v", err)
			}
			defer r.Body.Close()
			
			// If JSON values are provided, use them directly
			if params.Min != 0 || params.Max != 0 {
				return params, nil
			}
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return RandomNumberParams{}, fmt.Errorf("error parsing form data: %v", err)
			}
			minStr = r.FormValue("min")
			maxStr = r.FormValue("max")
		}
	} else {
		// Get parameters from query string
		minStr = r.URL.Query().Get("min")
		maxStr = r.URL.Query().Get("max")
	}

	// Set default values
	min := 1
	max := 100

	// Parse min parameter if provided
	if minStr != "" {
		parsedMin, err := strconv.Atoi(minStr)
		if err != nil {
			return RandomNumberParams{}, fmt.Errorf("invalid min parameter: %s", minStr)
		}
		min = parsedMin
	}

	// Parse max parameter if provided
	if maxStr != "" {
		parsedMax, err := strconv.Atoi(maxStr)
		if err != nil {
			return RandomNumberParams{}, fmt.Errorf("invalid max parameter: %s", maxStr)
		}
		max = parsedMax
	}

	// Create and return parameters
	return RandomNumberParams{
		Min: min,
		Max: max,
	}, nil
}

// ExecuteRandomNumber executes the random number generator with the given HTTP request
// It parses parameters from the request and returns the generated random number as a string
func ExecuteRandomNumber(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseRandomNumberParams(r)
	if err != nil {
		return "", err
	}

	// Generate random number
	randNum, err := GenerateRandomNumber(params)
	if err != nil {
		return "", err
	}

	// Return the number as a string
	return fmt.Sprintf("%d", randNum), nil
}
