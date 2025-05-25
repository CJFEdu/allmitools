// Package models contains the data structures for the AllMiTools server
package models

// TextFormatterTool returns the tool info for the text formatter
func TextFormatterTool() ToolInfo {
	return ToolInfo{
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
	}
}
