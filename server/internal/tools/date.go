// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// ParseDateFormatterParams parses the date formatter parameters from an HTTP request
// It handles both POST and GET requests
func ParseDateFormatterParams(r *http.Request) (DateFormatterParams, error) {
	// Parse parameters from either POST or GET
	var format, offsetStr string

	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				return DateFormatterParams{}, fmt.Errorf("error parsing form data: %v", err)
			}

			// Get parameters from form data
			format = r.FormValue("format")
			offsetStr = r.FormValue("offset")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var params DateFormatterParams
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return DateFormatterParams{}, fmt.Errorf("error parsing JSON data: %v", err)
			}
			defer r.Body.Close()
			
			// Return the parsed parameters directly
			return params, nil
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				return DateFormatterParams{}, fmt.Errorf("error parsing form data: %v", err)
			}
			format = r.FormValue("format")
			offsetStr = r.FormValue("offset")
		}
	} else {
		// Get parameters from query string
		format = r.URL.Query().Get("format")
		offsetStr = r.URL.Query().Get("offset")
	}

	// Set default values
	offset := 0

	// Parse offset parameter if provided
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return DateFormatterParams{}, fmt.Errorf("invalid offset parameter: %s", offsetStr)
		}
		offset = parsedOffset
	}

	// Create and return parameters
	return DateFormatterParams{
		Format: format,
		Offset: offset,
	}, nil
}

// ExecuteDateFormatter executes the date formatter with the given HTTP request
// It parses parameters from the request and returns the formatted date as a string
func ExecuteDateFormatter(r *http.Request) (string, error) {
	// Parse parameters from the request
	params, err := ParseDateFormatterParams(r)
	if err != nil {
		return "", err
	}

	// Format date
	formattedDate, err := FormatDate(params, time.Now())
	if err != nil {
		return "", err
	}

	// Return the formatted date as a string
	return formattedDate, nil
}

// ExecuteDateComponent executes a date component request (day, month, or year)
// It returns the component value as a string and any error that occurred
func ExecuteDateComponent(component string) (string, error) {
	// Get the specified date component
	result, err := GetDateComponent(component, time.Now())
	if err != nil {
		return "", err
	}

	// Convert the result to a string and return it
	return fmt.Sprintf("%v", result), nil
}
