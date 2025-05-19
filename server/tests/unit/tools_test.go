// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"testing"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/tools"
	"github.com/stretchr/testify/assert"
)

// TestRandomNumberValidation tests the validation of random number parameters
func TestRandomNumberValidation(t *testing.T) {
	// Test cases for parameter validation
	testCases := []struct {
		name          string
		params        tools.RandomNumberParams
		expectedError string
	}{
		{
			name: "Valid parameters",
			params: tools.RandomNumberParams{
				Min: 1,
				Max: 100,
			},
			expectedError: "",
		},
		{
			name: "Min greater than Max",
			params: tools.RandomNumberParams{
				Min: 100,
				Max: 1,
			},
			expectedError: "minimum value cannot be greater than maximum value",
		},
		{
			name: "Equal Min and Max",
			params: tools.RandomNumberParams{
				Min: 50,
				Max: 50,
			},
			expectedError: "",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tools.ValidateRandomNumberParams(tc.params)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

// TestGenerateRandomNumber tests the random number generation
func TestGenerateRandomNumber(t *testing.T) {
	// Test cases for random number generation
	testCases := []struct {
		name          string
		params        tools.RandomNumberParams
		expectError   bool
		validateRange bool
	}{
		{
			name: "Generate number in range 1-100",
			params: tools.RandomNumberParams{
				Min: 1,
				Max: 100,
			},
			expectError:   false,
			validateRange: true,
		},
		{
			name: "Generate number with equal min and max",
			params: tools.RandomNumberParams{
				Min: 42,
				Max: 42,
			},
			expectError:   false,
			validateRange: true,
		},
		{
			name: "Error with invalid range",
			params: tools.RandomNumberParams{
				Min: 100,
				Max: 1,
			},
			expectError:   true,
			validateRange: false,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			num, err := tools.GenerateRandomNumber(tc.params)
			
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				if tc.validateRange {
					assert.GreaterOrEqual(t, num, tc.params.Min)
					assert.LessOrEqual(t, num, tc.params.Max)
				}
			}
		})
	}
}

// TestDateFormatterValidation tests the validation of date formatter parameters
func TestDateFormatterValidation(t *testing.T) {
	// Test cases for parameter validation
	testCases := []struct {
		name          string
		params        tools.DateFormatterParams
		expectedError string
	}{
		{
			name: "Valid parameters with format",
			params: tools.DateFormatterParams{
				Format: "2006-01-02",
				Offset: 0,
			},
			expectedError: "",
		},
		{
			name: "Valid parameters with offset",
			params: tools.DateFormatterParams{
				Format: "2006-01-02",
				Offset: 7,
			},
			expectedError: "",
		},
		{
			name: "Empty format",
			params: tools.DateFormatterParams{
				Format: "",
				Offset: 0,
			},
			expectedError: "",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tools.ValidateDateFormatterParams(tc.params)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

// TestFormatDate tests the date formatting functionality
func TestFormatDate(t *testing.T) {
	// Create a fixed time for testing
	fixedTime := time.Date(2025, 5, 18, 10, 0, 0, 0, time.UTC)
	
	// Test cases for date formatting
	testCases := []struct {
		name           string
		params         tools.DateFormatterParams
		expectedOutput string
	}{
		{
			name: "Default format",
			params: tools.DateFormatterParams{
				Format: "",
				Offset: 0,
			},
			expectedOutput: "2025-05-18",
		},
		{
			name: "Custom format",
			params: tools.DateFormatterParams{
				Format: "Jan 2, 2006",
				Offset: 0,
			},
			expectedOutput: "May 18, 2025",
		},
		{
			name: "With positive offset",
			params: tools.DateFormatterParams{
				Format: "2006-01-02",
				Offset: 7,
			},
			expectedOutput: "2025-05-25",
		},
		{
			name: "With negative offset",
			params: tools.DateFormatterParams{
				Format: "2006-01-02",
				Offset: -7,
			},
			expectedOutput: "2025-05-11",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tools.FormatDate(tc.params, fixedTime)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}

// TestGetDateComponent tests the date component functionality
func TestGetDateComponent(t *testing.T) {
	// Create a fixed time for testing
	fixedTime := time.Date(2025, 5, 18, 10, 0, 0, 0, time.UTC)
	
	// Test cases for date components
	testCases := []struct {
		name           string
		component      string
		expectedOutput interface{}
		expectError    bool
	}{
		{
			name:           "Get day",
			component:      "day",
			expectedOutput: 18,
			expectError:    false,
		},
		{
			name:           "Get month",
			component:      "month",
			expectedOutput: "May",
			expectError:    false,
		},
		{
			name:           "Get year",
			component:      "year",
			expectedOutput: 2025,
			expectError:    false,
		},
		{
			name:        "Invalid component",
			component:   "invalid",
			expectError: true,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tools.GetDateComponent(tc.component, fixedTime)
			
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, result)
			}
		})
	}
}

// TestTextFileValidation tests the validation of text file parameters
func TestTextFileValidation(t *testing.T) {
	// Test cases for parameter validation
	testCases := []struct {
		name          string
		params        tools.TextFileParams
		expectedError string
	}{
		{
			name: "Valid parameters with content",
			params: tools.TextFileParams{
				Content:  "This is test content",
				Filename: "test.txt",
			},
			expectedError: "",
		},
		{
			name: "Valid parameters without filename",
			params: tools.TextFileParams{
				Content:  "This is test content",
				Filename: "",
			},
			expectedError: "",
		},
		{
			name: "Empty content",
			params: tools.TextFileParams{
				Content:  "",
				Filename: "test.txt",
			},
			expectedError: "content cannot be empty",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tools.ValidateTextFileParams(tc.params)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

// TestGenerateTextFile tests the text file generation functionality
func TestGenerateTextFile(t *testing.T) {
	// Test cases for text file generation
	testCases := []struct {
		name            string
		params          tools.TextFileParams
		expectedContent string
		expectedName    string
		expectError     bool
	}{
		{
			name: "Generate file with content and filename",
			params: tools.TextFileParams{
				Content:  "This is test content",
				Filename: "test.txt",
			},
			expectedContent: "This is test content",
			expectedName:    "test.txt",
			expectError:     false,
		},
		{
			name: "Generate file with default filename",
			params: tools.TextFileParams{
				Content:  "This is test content",
				Filename: "",
			},
			expectedContent: "This is test content",
			expectedName:    "download.txt",
			expectError:     false,
		},
		{
			name: "Error with empty content",
			params: tools.TextFileParams{
				Content:  "",
				Filename: "test.txt",
			},
			expectError: true,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, filename, err := tools.GenerateTextFile(tc.params)
			
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedContent, content)
				assert.Equal(t, tc.expectedName, filename)
			}
		})
	}
}
