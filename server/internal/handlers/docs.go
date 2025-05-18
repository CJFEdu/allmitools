// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/gorilla/mux"
)

// DocsResponse represents the response for documentation requests
type DocsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DocsBaseHandler handles requests to the documentation base path
func DocsBaseHandler(w http.ResponseWriter, r *http.Request) {
	// Get the list of all available tools
	tools := models.ListTools()
	
	// Check if the client accepts JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DocsResponse{
			Success: true,
			Message: "List of available tools",
			Data:    tools,
		})
		return
	}
	
	// Default to HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>AllMiTools Documentation</h1>")
	fmt.Fprintf(w, "<h2>Available Tools</h2>")
	fmt.Fprintf(w, "<ul>")
	
	for _, tool := range tools {
		fmt.Fprintf(w, "<li><a href='/docs/%s'>%s</a> - %s</li>", 
			tool.Name, tool.Name, tool.Description)
	}
	
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</body></html>")
}

// DocsToolHandler handles requests for documentation about a specific tool
func DocsToolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["tool_name"]
	
	// Get tool info
	toolInfo, err := models.GetToolInfo(toolName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		
		// Check if the client accepts JSON
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(DocsResponse{
				Success: false,
				Error:   fmt.Sprintf("Tool not found: %s", toolName),
			})
			return
		}
		
		// Default to HTML response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Tool Not Found</h1>")
		fmt.Fprintf(w, "<p>The tool '%s' was not found.</p>", toolName)
		fmt.Fprintf(w, "<p><a href='/docs'>Back to documentation index</a></p>")
		fmt.Fprintf(w, "</body></html>")
		return
	}
	
	// Check if the client accepts JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DocsResponse{
			Success: true,
			Message: fmt.Sprintf("Documentation for tool: %s", toolName),
			Data:    toolInfo,
		})
		return
	}
	
	// Default to HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>%s</h1>", toolInfo.Name)
	fmt.Fprintf(w, "<p><strong>Description:</strong> %s</p>", toolInfo.Description)
	fmt.Fprintf(w, "<p><strong>Version:</strong> %s</p>", toolInfo.Version)
	fmt.Fprintf(w, "<p><strong>Author:</strong> %s</p>", toolInfo.Author)
	fmt.Fprintf(w, "<p><strong>Output Type:</strong> %s</p>", toolInfo.OutputType)
	
	// Display parameters
	if len(toolInfo.Parameters) > 0 {
		fmt.Fprintf(w, "<h2>Parameters</h2>")
		fmt.Fprintf(w, "<table border='1'>")
		fmt.Fprintf(w, "<tr><th>Name</th><th>Type</th><th>Description</th><th>Required</th><th>Default</th></tr>")
		
		for _, param := range toolInfo.Parameters {
			defaultValue := ""
			if param.Default != nil {
				defaultValue = fmt.Sprintf("%v", param.Default)
			}
			
			fmt.Fprintf(w, "<tr>")
			fmt.Fprintf(w, "<td>%s</td>", param.Name)
			fmt.Fprintf(w, "<td>%s</td>", param.Type)
			fmt.Fprintf(w, "<td>%s</td>", param.Description)
			fmt.Fprintf(w, "<td>%v</td>", param.Required)
			fmt.Fprintf(w, "<td>%s</td>", defaultValue)
			fmt.Fprintf(w, "</tr>")
		}
		
		fmt.Fprintf(w, "</table>")
	}
	
	// Add a link back to the documentation index
	fmt.Fprintf(w, "<p><a href='/docs'>Back to documentation index</a></p>")
	
	// Add a link to use the tool
	fmt.Fprintf(w, "<p><a href='/tools/%s'>Use this tool</a></p>", toolInfo.Name)
	
	fmt.Fprintf(w, "</body></html>")
}
