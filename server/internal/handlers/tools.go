// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/gorilla/mux"
)

// ToolResponse represents the response from a tool
type ToolResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ToolsHandler handles requests to use specific tools
func ToolsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["tool_name"]
	
	// Get tool info
	toolInfo, err := models.GetToolInfo(toolName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ToolResponse{
			Success: false,
			Error:   fmt.Sprintf("Tool not found: %s", toolName),
		})
		return
	}
	
	// Set the appropriate content type based on the tool's output type
	switch toolInfo.OutputType {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		// For now, just return the tool info as JSON
		json.NewEncoder(w).Encode(ToolResponse{
			Success: true,
			Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
			Data:    toolInfo,
		})
	case "html":
		w.Header().Set("Content-Type", "text/html")
		// For now, just return a simple HTML response
		fmt.Fprintf(w, "<html><body><h1>Tool: %s</h1><p>%s</p></body></html>", 
			toolInfo.Name, toolInfo.Description)
	case "raw":
		w.Header().Set("Content-Type", "text/plain")
		// For now, just return a simple text response
		fmt.Fprintf(w, "Tool: %s\nDescription: %s", toolInfo.Name, toolInfo.Description)
	default:
		// Default to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToolResponse{
			Success: true,
			Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
			Data:    toolInfo,
		})
	}
}
