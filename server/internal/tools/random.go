// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"errors"
	"math/rand"
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
