// Package models contains the data structures for the AllMiTools server
package models

// URLEncoderTool returns the tool info for the URL encoder
func URLEncoderTool() ToolInfo {
	return ToolInfo{
		Name:        "url-encoder",
		Description: "Convert a string to a URL-encoded string",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters: []ToolParameter{
			{
				Name:        "text",
				Type:        "string",
				Description: "Text to URL encode",
				Required:    true,
			},
		},
	}
}
