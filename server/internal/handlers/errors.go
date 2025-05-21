// Package handlers provides HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/templates"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404
	w.WriteHeader(http.StatusNotFound)
	
	// Check if the client accepts JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("404 - Page Not Found: %s", r.URL.Path),
		})
		return
	}
	
	// Default to HTML response using template
	data := map[string]interface{}{
		"Title":       "Page Not Found",
		"CurrentPage": "",
		"Path":        r.URL.Path,
	}
	
	// Render the template
	err := templates.TemplateManager.RenderTemplate(w, "404", data)
	if err != nil {
		// Fallback if template rendering fails
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "404 - Page Not Found: %s", r.URL.Path)
	}
}
