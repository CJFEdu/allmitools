// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"fmt"
	"time"
)

// DateFormatterParams represents the parameters for the date formatter
type DateFormatterParams struct {
	Format string `json:"format"` // Format string (e.g., "2006-01-02")
	Offset int    `json:"offset"` // Offset in days (can be negative)
}

// DefaultDateFormat is the default format for the date formatter
const DefaultDateFormat = "2006-01-02"

// ValidateDateFormatterParams validates the parameters for the date formatter
func ValidateDateFormatterParams(params DateFormatterParams) error {
	// If format is empty, we'll use the default format, so no error
	return nil
}

// FormatDate formats the current date according to the specified parameters
// It returns the formatted date and any error that occurred
func FormatDate(params DateFormatterParams, now time.Time) (string, error) {
	// Validate parameters
	if err := ValidateDateFormatterParams(params); err != nil {
		return "", err
	}

	// Use default format if not specified
	format := params.Format
	if format == "" {
		format = DefaultDateFormat
	}

	// Apply offset
	if params.Offset != 0 {
		now = now.AddDate(0, 0, params.Offset)
	}

	// Format date
	return now.Format(format), nil
}

// GetCurrentDay returns the current day of the month
func GetCurrentDay(now time.Time) (int, error) {
	return now.Day(), nil
}

// GetCurrentMonth returns the current month as a string
func GetCurrentMonth(now time.Time) (string, error) {
	return now.Month().String(), nil
}

// GetCurrentYear returns the current year
func GetCurrentYear(now time.Time) (int, error) {
	return now.Year(), nil
}

// GetDateComponent returns a specific component of the current date
// Valid components are "day", "month", and "year"
func GetDateComponent(component string, now time.Time) (interface{}, error) {
	switch component {
	case "day":
		return GetCurrentDay(now)
	case "month":
		return GetCurrentMonth(now)
	case "year":
		return GetCurrentYear(now)
	default:
		return nil, fmt.Errorf("invalid date component: %s", component)
	}
}
