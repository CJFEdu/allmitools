// Package models contains the data structures for the AllMiTools server
package models

// DateFormatterTool returns the tool info for the date formatter
func DateFormatterTool() ToolInfo {
	return ToolInfo{
		Name:        "date",
		Description: "Format the current date with optional offset",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "format",
				Type:        "string",
				Description: "Date format string (e.g., '2006-01-02' for YYYY-MM-DD)",
				Required:    false,
				Default:     "2006-01-02",
			},
			{
				Name:        "offset",
				Type:        "int",
				Description: "Offset in days (can be negative)",
				Required:    false,
				Default:     0,
			},
		},
	}
}
