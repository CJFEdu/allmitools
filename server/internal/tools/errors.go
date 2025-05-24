// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"fmt"
)

// Common error functions for parameter validation
func ErrMissingRequiredParameter(paramName string) error {
	return fmt.Errorf("missing required parameter: %s", paramName)
}

func ErrInvalidParameter(message string) error {
	return fmt.Errorf("invalid parameter: %s", message)
}
