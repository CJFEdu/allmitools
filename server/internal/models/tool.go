// Package models contains the data structures for the AllMiTools server
package models

import (
	"errors"
	"fmt"
)

// ToolParameter represents a parameter for a tool
type ToolParameter struct {
	Name        string `json:"name"`        // Name of the parameter
	Type        string `json:"type"`        // Type of the parameter (string, int, bool, etc.)
	Description string `json:"description"` // Description of the parameter
	Required    bool   `json:"required"`    // Whether the parameter is required
	Default     any    `json:"default"`     // Default value for the parameter (if any)
}

// Validate checks if the parameter is valid
func (p *ToolParameter) Validate() error {
	if p.Name == "" {
		return errors.New("parameter name cannot be empty")
	}
	
	if p.Type == "" {
		return errors.New("parameter type cannot be empty")
	}
	
	// Validate parameter type
	validTypes := map[string]bool{
		"string":  true,
		"int":     true,
		"float":   true,
		"bool":    true,
		"array":   true,
		"object":  true,
	}
	
	if !validTypes[p.Type] {
		return fmt.Errorf("invalid parameter type: %s", p.Type)
	}
	
	return nil
}

// ToolInfo represents information about a tool
type ToolInfo struct {
	Name        string          `json:"name"`        // Name of the tool
	Description string          `json:"description"` // Description of the tool
	Version     string          `json:"version"`     // Version of the tool
	Author      string          `json:"author"`      // Author of the tool
	Parameters  []ToolParameter `json:"parameters"`  // Parameters for the tool
}

// Validate checks if the tool info is valid
func (t *ToolInfo) Validate() error {
	if t.Name == "" {
		return errors.New("tool name cannot be empty")
	}
	
	if t.Description == "" {
		return errors.New("tool description cannot be empty")
	}
	
	if t.Version == "" {
		return errors.New("tool version cannot be empty")
	}
	
	// Output type validation removed
	
	// Validate parameters
	for i, param := range t.Parameters {
		if err := param.Validate(); err != nil {
			return fmt.Errorf("parameter %d (%s): %w", i, param.Name, err)
		}
	}
	
	return nil
}

// AvailableTools is a map of available tools
var AvailableTools = map[string]ToolInfo{
	"random-number": {
		Name:        "random-number",
		Description: "Generate a random number within a specified range",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "min",
				Type:        "int",
				Description: "Minimum value (inclusive)",
				Required:    false,
				Default:     1,
			},
			{
				Name:        "max",
				Type:        "int",
				Description: "Maximum value (inclusive)",
				Required:    false,
				Default:     100,
			},
		},
	},
	"text-file": {
		Name:        "text-file",
		Description: "Generate a downloadable text file from provided content",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "content",
				Type:        "string",
				Description: "Content to be saved as a text file",
				Required:    true,
			},
			{
				Name:        "filename",
				Type:        "string",
				Description: "Optional filename for the text file",
				Required:    false,
				Default:     "download.txt",
			},
		},
	},
	"text-formatter": {
		Name:        "text-formatter",
		Description: "Format text with various options",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "text",
				Type:        "string",
				Description: "Text to format",
				Required:    true,
			},
			{
				Name:        "uppercase",
				Type:        "bool",
				Description: "Convert text to uppercase (if false, converts to lowercase)",
				Required:    false,
				Default:     false,
			},
		},
	},
}

// GetToolInfo returns information about a tool
func GetToolInfo(toolName string) (ToolInfo, error) {
	tool, exists := AvailableTools[toolName]
	if !exists {
		return ToolInfo{}, fmt.Errorf("tool not found: %s", toolName)
	}
	return tool, nil
}

// ListTools returns a list of all available tools
func ListTools() []ToolInfo {
	tools := make([]ToolInfo, 0, len(AvailableTools))
	for _, tool := range AvailableTools {
		tools = append(tools, tool)
	}
	return tools
}
