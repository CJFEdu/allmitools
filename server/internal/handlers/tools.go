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

		// Check if the client accepts JSON
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ToolResponse{
				Success: false,
				Error:   fmt.Sprintf("Tool not found: %s", toolName),
			})
			return
		}

		// Default to HTML response using 404 template
		data := map[string]interface{}{
			"Title":       "Tool Not Found",
			"CurrentPage": "tools",
		}

		// Render the 404 template
		err := templates.TemplateManager.RenderTemplate(w, "404", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>Tool Not Found</h1>")
			fmt.Fprintf(w, "<p>The tool '%s' was not found.</p>", toolName)
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
				Message: fmt.Sprintf("Tool information for: %s", toolName),
				Data:    toolInfo,
			})
			return
		}

		// Default to HTML response using template
		data := map[string]interface{}{
			"Title":       toolInfo.Name,
			"CurrentPage": "tools",
			"Tool":        toolInfo,
		}

		// Render the template
		err := templates.TemplateManager.RenderTemplate(w, "tool", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>%s</h1>", toolInfo.Name)
			fmt.Fprintf(w, "<p>%s</p>", toolInfo.Description)
			fmt.Fprintf(w, "<p><a href='/docs/%s'>View Documentation</a></p>", toolInfo.Name)
			fmt.Fprintf(w, "</body></html>")
		}
		return
	}

	// Execute the tool (for POST or GET with parameters)
	// Parse query parameters and execute the appropriate tool
	var result string
	var toolErr error

	switch toolName {
	case "random-number":
		result, toolErr = tools.ExecuteRandomNumber(r)
	case "random-string":
		result, toolErr = tools.ExecuteRandomString(r)
	case "date":
		result, toolErr = tools.ExecuteDateFormatter(r)
	case "day":
		result, toolErr = tools.ExecuteDateComponent("day")
	case "month":
		result, toolErr = tools.ExecuteDateComponent("month")
	case "year":
		result, toolErr = tools.ExecuteDateComponent("year")
	case "text-formatter":
		result, toolErr = tools.ExecuteTextFormatter(r)
	case "url-encoder":
		result, toolErr = tools.ExecuteURLEncoder(r)
	case "sha256-hasher":
		result, toolErr = tools.ExecuteSHA256Hasher(r)
	case "text-file":
		// Special handling for text file tool (returns a file download)
		fileContent, fileName, err := tools.ExecuteTextFile(r)
		if err != nil {
			// Return error as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ToolResponse{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			// Set headers for file download
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

			// Write the file content to the response
			w.Write([]byte(fileContent))
		}
		return // Early return as we've already written the response
	default:
		// For unknown tools, return an error
		toolErr = fmt.Errorf("unknown tool: %s", toolName)
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
				Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
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
		Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
		Data:    result,
	})
}

// generateHTMLResponse generates an HTML response based on the tool result
func generateHTMLResponse(w http.ResponseWriter, result string) {
	fmt.Fprintf(w, "<html><body>")

	// Just display the result in a paragraph
	fmt.Fprintf(w, "<p>%s</p>", result)

	fmt.Fprintf(w, "</body></html>")
}

// generateRawResponse generates a raw text response based on the tool result
func generateRawResponse(w http.ResponseWriter, result string) {
	// Just write the result string directly
	fmt.Fprintf(w, "%s", result)
}


