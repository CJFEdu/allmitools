// Package models contains the data structures for the AllMiTools server
package models

// TextFileTool returns the tool info for the text file generator
func TextFileTool() ToolInfo {
	return ToolInfo{
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
	}
}
