// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
"encoding/json"
"fmt"
"net/http"
"strconv"
"time"

"github.com/CJFEdu/allmitools/server/internal/models"
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
json.NewEncoder(w).Encode(ToolResponse{
Success: false,
Error:   fmt.Sprintf("Tool not found: %s", toolName),
})
return
}

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
// Parse query parameters
minStr := r.URL.Query().Get("min")
maxStr := r.URL.Query().Get("max")

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
randomNum, err := tools.GenerateRandomNumber(params)
if err != nil {
return nil, err
}

// Return result
return map[string]interface{}{
"number": randomNum,
"min":    min,
"max":    max,
}, nil
}

// executeDateFormatterTool executes the date formatter tool
func executeDateFormatterTool(r *http.Request) (interface{}, error) {
// Parse query parameters
format := r.URL.Query().Get("format")
offsetStr := r.URL.Query().Get("offset")

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
