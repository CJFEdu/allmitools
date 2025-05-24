// Package models contains the data structures for the AllMiTools server
package models

import "fmt"

// PrivateToolInfo represents information about a private tool
// It extends the regular ToolInfo with additional fields for private tools
type PrivateToolInfo struct {
	ToolInfo
	RequiresAuth bool `json:"requires_auth"` // Whether the tool requires authentication
}

// AvailablePrivateTools is a map of available private tools
var AvailablePrivateTools = map[string]PrivateToolInfo{
	"text-storage": TextStorageTool(),
	"text-retrieval": TextRetrievalTool(),
}

// TextStorageTool returns the tool info for the text storage tool
func TextStorageTool() PrivateToolInfo {
	return PrivateToolInfo{
		ToolInfo: ToolInfo{
			Name:        "text-storage",
			Description: "Stores text content in the database with a unique ID",
			Version:     "1.0.0",
			Author:      "AllMiTools Team",
			Parameters: []ToolParameter{
				{
					Name:        "content",
					Description: "The text content to store",
					Type:        "text",
					Required:    true,
					Default:     "",
				},
				{
					Name:        "save",
					Description: "Whether to save the text permanently",
					Type:        "boolean",
					Required:    false,
					Default:     "false",
				},
			},
		},
		RequiresAuth: true,
	}
}

// TextRetrievalTool returns the tool info for the text retrieval tool
func TextRetrievalTool() PrivateToolInfo {
	return PrivateToolInfo{
		ToolInfo: ToolInfo{
			Name:        "text-retrieval",
			Description: "Retrieves text content from the database using a unique ID",
			Version:     "1.0.0",
			Author:      "AllMiTools Team",
			Parameters: []ToolParameter{
				{
					Name:        "id",
					Description: "The unique ID of the text to retrieve",
					Type:        "string",
					Required:    true,
					Default:     "",
				},
			},
		},
		RequiresAuth: true,
	}
}

// GetPrivateToolInfo returns information about a specific private tool
func GetPrivateToolInfo(toolName string) (PrivateToolInfo, error) {
	tool, exists := AvailablePrivateTools[toolName]
	if !exists {
		return PrivateToolInfo{}, fmt.Errorf("tool not found: %s", toolName)
	}
	return tool, nil
}

// GetAllPrivateTools returns a list of all available private tools
func GetAllPrivateTools() []PrivateToolInfo {
	tools := make([]PrivateToolInfo, 0, len(AvailablePrivateTools))
	for _, tool := range AvailablePrivateTools {
		tools = append(tools, tool)
	}
	return tools
}
