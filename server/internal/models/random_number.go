// Package models contains the data structures for the AllMiTools server
package models

// RandomNumberTool returns the tool info for the random number generator
func RandomNumberTool() ToolInfo {
	return ToolInfo{
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
	}
}
