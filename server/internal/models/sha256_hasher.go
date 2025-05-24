// Package models contains the data structures for the AllMiTools server
package models

// SHA256HasherTool returns the tool info for the SHA-256 hasher
func SHA256HasherTool() ToolInfo {
	return ToolInfo{
		Name:        "sha256-hasher",
		Description: "Convert a text string into a SHA-256 hash",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "text",
				Type:        "string",
				Description: "Text to hash",
				Required:    true,
			},
		},
	}
}
