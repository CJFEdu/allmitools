// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/CJFEdu/allmitools/server/internal/tools"
)

// PrivateToolsHandler handles requests to use specific private tools
// This handler is protected by the auth middleware
func PrivateToolsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["tool_name"]

	// Get tool info
	toolInfo, err := models.GetPrivateToolInfo(toolName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		// Check if the client accepts JSON
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ToolResponse{
				Success: false,
				Error:   fmt.Sprintf("Private tool not found: %s", toolName),
			})
			return
		}

		// Default to HTML response using 404 template
		data := map[string]interface{}{
			"Title":       "Private Tool Not Found",
			"CurrentPage": "private-tools",
		}

		// Render the 404 template
		err := templates.TemplateManager.RenderTemplate(w, "404", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>Private Tool Not Found</h1>")
			fmt.Fprintf(w, "<p>The private tool '%s' was not found.</p>", toolName)
			fmt.Fprintf(w, "<p><a href='/'>Back to homepage</a></p>")
			fmt.Fprintf(w, "</body></html>")
		}
		return
	}

	// Determine if we should execute the tool or show the form
	// Execute if: POST request OR GET request with parameters
	// Show form if: GET request without parameters
	shouldExecute := false

	// Always execute for POST requests
	if r.Method == http.MethodPost {
		shouldExecute = true
	}

	// For GET requests, check if there are any query parameters for the tool
	if r.Method == http.MethodGet {
		// Check if any tool-specific parameters are present
		hasParameters := false
		for _, param := range toolInfo.Parameters {
			if r.URL.Query().Get(param.Name) != "" {
				hasParameters = true
				break
			}
		}

		// Also check for output_format parameter
		if r.URL.Query().Get("output_format") != "" {
			hasParameters = true
		}

		// If there are parameters, execute the tool
		if hasParameters {
			shouldExecute = true
		}
	}

	// If we shouldn't execute (GET without parameters), show the form
	if !shouldExecute {
		// Check if the client accepts JSON
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ToolResponse{
				Success: true,
				Message: fmt.Sprintf("Private tool information for: %s", toolName),
				Data:    toolInfo,
			})
			return
		}

		// Default to HTML response using template
		data := map[string]interface{}{
			"Title":       toolInfo.Name,
			"CurrentPage": "private-tools",
			"Tool":        toolInfo,
			"IsPrivate":   true,
		}

		// Render the template
		err := templates.TemplateManager.RenderTemplate(w, "tool", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>%s (Private)</h1>", toolInfo.Name)
			fmt.Fprintf(w, "<p>%s</p>", toolInfo.Description)
			fmt.Fprintf(w, "<p><a href='/private/docs/%s'>View Documentation</a></p>", toolInfo.Name)
			fmt.Fprintf(w, "</body></html>")
		}
		return
	}

	// Execute the tool (for POST or GET with parameters)
	// Parse query parameters and execute the appropriate tool
	var result string
	var toolErr error

	switch toolName {
	case "private-demo":
		result, toolErr = tools.ExecutePrivateDemo(r)
	default:
		// For unknown tools, return an error
		toolErr = fmt.Errorf("unknown private tool: %s", toolName)
		result = ""
	}

	// Handle tool execution error
	if toolErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ToolResponse{
			Success: false,
			Error:   toolErr.Error(),
		})
		return
	}

	// Check if the user specified an output format in the form, query parameter, or Accept header
	outputFormat := ""

	// 1. Check form parameter (highest priority)
	if r.Method == http.MethodPost {
		outputFormat = r.FormValue("output_format")
	}

	// 2. Check query parameter if no form parameter
	if outputFormat == "" {
		outputFormat = r.URL.Query().Get("output_format")
	}

	// 3. Check Accept header if no form or query parameter
	if outputFormat == "" && strings.Contains(r.Header.Get("Accept"), "application/json") {
		outputFormat = "json"
	} else if outputFormat == "" && strings.Contains(r.Header.Get("Accept"), "text/plain") {
		outputFormat = "raw"
	}

	// Determine the appropriate output format
	// Priority: 1. Form parameter, 2. Query parameter, 3. Accept header, 4. Tool's default output type
	if outputFormat != "" {
		// Use the format specified by the user (form, query parameter, or Accept header)
		switch outputFormat {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ToolResponse{
				Success: true,
				Message: fmt.Sprintf("Private tool '%s' executed successfully", toolName),
				Data:    result,
			})
			return
		case "html":
			w.Header().Set("Content-Type", "text/html")
			generateHTMLResponse(w, result)
			return
		case "raw":
			w.Header().Set("Content-Type", "text/plain")
			generateRawResponse(w, result)
			return
		}
	}

	// If no form parameter or invalid format, default to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToolResponse{
		Success: true,
		Message: fmt.Sprintf("Private tool '%s' executed successfully", toolName),
		Data:    result,
	})
}

// PrivateToolsListHandler handles requests to list all available private tools
func PrivateToolsListHandler(w http.ResponseWriter, r *http.Request) {
	// Get all private tools
	privateTools := models.GetAllPrivateTools()

	// Check if the client accepts JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToolResponse{
			Success: true,
			Message: "List of available private tools",
			Data:    privateTools,
		})
		return
	}

	// Default to HTML response using template
	data := map[string]interface{}{
		"Title":        "Private Tools",
		"CurrentPage":  "private-tools",
		"PrivateTools": privateTools,
	}

	// Render the template
	err := templates.TemplateManager.RenderTemplate(w, "private_tools_list", data)
	if err != nil {
		// Fallback if template rendering fails
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Private Tools</h1>")
		fmt.Fprintf(w, "<ul>")
		for _, tool := range privateTools {
			fmt.Fprintf(w, "<li><a href='/private/tools/%s'>%s</a> - %s</li>", tool.Name, tool.Name, tool.Description)
		}
		fmt.Fprintf(w, "</ul>")
		fmt.Fprintf(w, "</body></html>")
	}
}
