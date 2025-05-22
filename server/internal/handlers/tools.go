// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/CJFEdu/allmitools/server/internal/tools"
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
	var result interface{}
	var toolErr error
	
	switch toolName {
	case "random-number":
		result, toolErr = executeRandomNumberTool(r)
	case "date":
		result, toolErr = executeDateFormatterTool(r)
	case "day":
		result, toolErr = executeDateComponentTool("day")
	case "month":
		result, toolErr = executeDateComponentTool("month")
	case "year":
		result, toolErr = executeDateComponentTool("year")
	case "text-file":
		// For the text file tool, we handle it differently as it needs to set special headers
		if err := executeTextFileTool(w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ToolResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
		return // Early return as we've already written the response
	default:
		// For unknown tools, just return the tool info
		result = toolInfo
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
		// Generate HTML based on the result
		generateHTMLResponse(w, toolName, result)
		return
	case "raw":
		w.Header().Set("Content-Type", "text/plain")
		// Generate raw text based on the result
		generateRawResponse(w, toolName, result)
		return
	}
}

// If no form parameter or invalid format, use Accept header or tool's default
// Set the appropriate content type based on the tool's output type
switch toolInfo.OutputType {
case "json":
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToolResponse{
		Success: true,
		Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
		Data:    result,
	})
case "html":
	w.Header().Set("Content-Type", "text/html")
	// Generate HTML based on the result
	generateHTMLResponse(w, toolName, result)
case "raw":
	w.Header().Set("Content-Type", "text/plain")
	// Generate raw text based on the result
	generateRawResponse(w, toolName, result)
default:
	// Default to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToolResponse{
	Success: true,
	Message: fmt.Sprintf("Tool '%s' executed successfully", toolName),
	Data:    result,
})
}
}

// executeRandomNumberTool executes the random number generator tool
func executeRandomNumberTool(r *http.Request) (interface{}, error) {
	// Parse parameters from either POST or GET
	var minStr, maxStr string
	
	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("error parsing form data: %v", err)
		}
		
		// Get parameters from form data
		minStr = r.FormValue("min")
		maxStr = r.FormValue("max")
	} else {
		// Get parameters from query string
		minStr = r.URL.Query().Get("min")
		maxStr = r.URL.Query().Get("max")
	}
	
	// Set default values
	min := 1
	max := 100
	
	// Parse min parameter if provided
	if minStr != "" {
		parsedMin, err := strconv.Atoi(minStr)
		if err != nil {
			return nil, fmt.Errorf("invalid min parameter: %s", minStr)
		}
		min = parsedMin
	}

	// Parse max parameter if provided
	if maxStr != "" {
		parsedMax, err := strconv.Atoi(maxStr)
		if err != nil {
			return nil, fmt.Errorf("invalid max parameter: %s", maxStr)
		}
		max = parsedMax
	}

	// Create parameters for the random number generator
	params := tools.RandomNumberParams{
		Min: min,
		Max: max,
	}
	
	// Generate random number
	randNum, err := tools.GenerateRandomNumber(params)
	if err != nil {
		return nil, err
	}
	
	// Return result
	return map[string]interface{}{
		"number": randNum,
		"min":    min,
		"max":    max,
	}, nil
}

// executeDateFormatterTool executes the date formatter tool
func executeDateFormatterTool(r *http.Request) (interface{}, error) {
	// Parse parameters from either POST or GET
	var format, offsetStr string
	
	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("error parsing form data: %v", err)
		}
		
		// Get parameters from form data
		format = r.FormValue("format")
		offsetStr = r.FormValue("offset")
	} else {
		// Get parameters from query string
		format = r.URL.Query().Get("format")
		offsetStr = r.URL.Query().Get("offset")
	}

	// Set default values
	offset := 0

// Parse offset parameter if provided
if offsetStr != "" {
parsedOffset, err := strconv.Atoi(offsetStr)
if err != nil {
return nil, fmt.Errorf("invalid offset parameter: %s", offsetStr)
}
offset = parsedOffset
}

// Create parameters for the date formatter
params := tools.DateFormatterParams{
Format: format,
Offset: offset,
}

// Format date
formattedDate, err := tools.FormatDate(params, time.Now())
if err != nil {
return nil, err
}

// Return result
return map[string]interface{}{
"date":   formattedDate,
"format": format,
"offset": offset,
}, nil
}

// executeDateComponentTool executes the day, month, or year tool
func executeDateComponentTool(component string) (interface{}, error) {
// Get the specified date component
result, err := tools.GetDateComponent(component, time.Now())
if err != nil {
return nil, err
}

// Return result
return map[string]interface{}{
component: result,
}, nil
}

// generateHTMLResponse generates an HTML response based on the tool result
func generateHTMLResponse(w http.ResponseWriter, toolName string, result interface{}) {
fmt.Fprintf(w, "<html><body>")

switch toolName {
case "random-number":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "<p>%v</p>", resultMap["number"])
}
case "date":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "<p>%v</p>", resultMap["date"])
}
case "day":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "<p>%v</p>", resultMap["day"])
}
case "month":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "<p>%v</p>", resultMap["month"])
}
case "year":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "<p>%v</p>", resultMap["year"])
}
default:
// For unknown tools, just display the result as JSON
resultJSON, _ := json.MarshalIndent(result, "", "  ")
fmt.Fprintf(w, "<p>%s</p>", resultJSON)
}

fmt.Fprintf(w, "</body></html>")
}

// generateRawResponse generates a raw text response based on the tool result
func generateRawResponse(w http.ResponseWriter, toolName string, result interface{}) {
switch toolName {
case "random-number":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "%v", resultMap["number"])
}
case "date":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "%v", resultMap["date"])
}
case "day":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "%v", resultMap["day"])
}
case "month":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "%v", resultMap["month"])
}
case "year":
if resultMap, ok := result.(map[string]interface{}); ok {
fmt.Fprintf(w, "%v", resultMap["year"])
}
default:
// For unknown tools, just display the result as JSON
resultJSON, _ := json.MarshalIndent(result, "", "  ")
fmt.Fprintf(w, "%s", resultJSON)
}
}

// executeTextFileTool handles the text file tool which generates a downloadable text file
// from the provided content
func executeTextFileTool(w http.ResponseWriter, r *http.Request) error {
	// Parse parameters from either POST or GET
	var content, filename string
	
	// Check if this is a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("error parsing form data: %v", err)
		}
		
		// Get parameters from form data
		content = r.FormValue("content")
		filename = r.FormValue("filename")
	} else {
		// Get parameters from query string
		content = r.URL.Query().Get("content")
		filename = r.URL.Query().Get("filename")
	}
	
	// Validate parameters
	if content == "" {
		return fmt.Errorf("content parameter is required")
	}

// Create parameters for the text file tool
params := tools.TextFileParams{
Content:  content,
Filename: filename,
}

// Generate text file
fileContent, fileName, err := tools.GenerateTextFile(params)
if err != nil {
return err
}

// Set headers for file download
w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
w.Header().Set("Content-Type", "text/plain")
w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

// Write the file content to the response
_, err = w.Write([]byte(fileContent))
return err
}
