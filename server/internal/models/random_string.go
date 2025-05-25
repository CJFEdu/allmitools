// Package models contains the data structures for the AllMiTools server
package models

// RandomStringTool returns the tool info for the random string generator
func RandomStringTool() ToolInfo {
	return ToolInfo{
		Name:        "random-string",
		Description: "Generate a random string of specified length",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "length",
				Type:        "int",
				Description: "Length of the random string",
				Required:    false,
				Default:     10,
			},
			{
				Name:        "mixedCase",
				Type:        "bool",
				Description: "Include both uppercase and lowercase letters (default: false, lowercase only)",
				Required:    false,
				Default:     false,
			},
		},
	}
}
