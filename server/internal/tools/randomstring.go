// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RandomStringParams represents the parameters for the random string generator
type RandomStringParams struct {
	Length     int  `json:"length"`     // Length of the random string
	MixedCase  bool `json:"mixedCase"`  // Whether to include both uppercase and lowercase letters
}

// ValidateRandomStringParams validates the parameters for the random string generator
func ValidateRandomStringParams(params RandomStringParams) error {
	if params.Length <= 0 {
		return fmt.Errorf("length must be greater than 0")
	}
	if params.Length > 1000 {
		return fmt.Errorf("length must be less than or equal to 1000")
	}
	return nil
}

// ParseRandomStringParams parses the random string generator parameters from an HTTP request
// It handles both POST and GET requests
func ParseRandomStringParams(r *http.Request) (RandomStringParams, error) {
	// Parse parameters from either POST or GET
	var lengthStr, mixedCaseStr string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return RandomStringParams{}, fmt.Errorf("error parsing form data: %v", err)
		}

		// Get parameters from form data
		lengthStr = r.FormValue("length")
		mixedCaseStr = r.FormValue("mixedCase")
	} else {
		// Get parameters from query string
		lengthStr = r.URL.Query().Get("length")
		mixedCaseStr = r.URL.Query().Get("mixedCase")
	}

	// Set default values
	length := 10
	mixedCase := false

	// Parse length parameter if provided
	if lengthStr != "" {
		parsedLength, err := strconv.Atoi(lengthStr)
		if err != nil {
			return RandomStringParams{}, fmt.Errorf("invalid length parameter: %s", lengthStr)
		}
		length = parsedLength
	}

	// Parse mixedCase parameter if provided
	if mixedCaseStr != "" {
		// Handle HTML form checkbox which sends "on" when checked
		if mixedCaseStr == "on" {
			mixedCase = true
		} else {
			// Try to parse as boolean
			parsedMixedCase, err := strconv.ParseBool(mixedCaseStr)
			if err != nil {
				return RandomStringParams{}, fmt.Errorf("invalid mixedCase parameter: %s (use true, false, or on)", mixedCaseStr)
			}
			mixedCase = parsedMixedCase
		}
	}

	// Create and return parameters
	return RandomStringParams{
		Length:    length,
		MixedCase: mixedCase,
	}, nil
}

// GenerateRandomString generates a random string of the specified length
// If mixedCase is true, it will include both uppercase and lowercase letters
// Otherwise, it will only include lowercase letters
func GenerateRandomString(length int, mixedCase bool) string {
	// Define character sets
	const lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	const uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	
	// Initialize random number generator with current time as seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Create a byte slice to hold the result
	result := make([]byte, length)
	
	// Determine the character set to use
	var charset string
	if mixedCase {
		charset = lowercaseLetters + uppercaseLetters
	} else {
		charset = lowercaseLetters
	}
	
	// Generate the random string
	for i := 0; i < length; i++ {
		result[i] = charset[r.Intn(len(charset))]
	}
	
	return string(result)
}

// ExecuteRandomString executes the random string generator with the given HTTP request
// It parses parameters from the request and returns the generated random string
func ExecuteRandomString(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseRandomStringParams(r)
	if err != nil {
		return "", err
	}

	// Validate parameters
	if err := ValidateRandomStringParams(params); err != nil {
		return "", err
	}

	// Generate the random string
	randomString := GenerateRandomString(params.Length, params.MixedCase)

	// Return the random string
	return randomString, nil
}
