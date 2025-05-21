// Package tools contains the implementation of various tools for the AllMiTools server
package tools

import (
	"errors"
)

// TextFileParams represents the parameters for the text file tool
type TextFileParams struct {
	Content string `json:"content"` // Content to be saved as a text file
	Filename string `json:"filename"` // Optional filename for the text file (default: "download.txt")
}

// ValidateTextFileParams validates the parameters for the text file tool
func ValidateTextFileParams(params TextFileParams) error {
	if params.Content == "" {
		return errors.New("content cannot be empty")
	}
	
	return nil
}

// GenerateTextFile prepares the content to be downloaded as a text file
// It returns the content and filename for the text file
func GenerateTextFile(params TextFileParams) (string, string, error) {
	// Validate parameters
	if err := ValidateTextFileParams(params); err != nil {
		return "", "", err
	}

	// Use default filename if not provided
	filename := params.Filename
	if filename == "" {
		filename = "download.txt"
	}

	return params.Content, filename, nil
}
