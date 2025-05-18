// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/models"
)

// HomeResponse represents the response for the homepage
type HomeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// HomeHandler handles requests to the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the list of all available tools
	tools := models.ListTools()
	
	// Check if the client accepts JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HomeResponse{
			Success: true,
			Message: "Welcome to the AllMiTools server!",
			Data: map[string]interface{}{
				"toolCount": len(tools),
				"docsUrl":   "/docs",
			},
		})
		return
	}
	
	// Default to HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>Welcome to the AllMiTools Server</h1>")
	fmt.Fprintf(w, "<p>AllMiTools is a collection of no-code automation tools.</p>")
	
	// Display tool count
	fmt.Fprintf(w, "<p>There are currently <strong>%d</strong> tools available:</p>", len(tools))
	
	// Display a list of tools
	fmt.Fprintf(w, "<ul>")
	for _, tool := range tools {
		fmt.Fprintf(w, "<li><a href='/tools/%s'>%s</a> - %s</li>", 
			tool.Name, tool.Name, tool.Description)
	}
	fmt.Fprintf(w, "</ul>")
	
	// Add a link to the documentation
	fmt.Fprintf(w, "<p><a href='/docs'>View Documentation</a></p>")
	
	fmt.Fprintf(w, "</body></html>")
}
