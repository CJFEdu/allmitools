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
	"private-demo": PrivateDemoTool(),
}

// PrivateDemoTool returns the tool info for the private demo tool
func PrivateDemoTool() PrivateToolInfo {
	return PrivateToolInfo{
		ToolInfo: ToolInfo{
			Name:        "private-demo",
			Description: "A demo private tool that requires authentication",
			Version:     "1.0.0",
			Author:      "AllMiTools Team",
			Parameters:  []ToolParameter{}, // No parameters needed for this demo
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
