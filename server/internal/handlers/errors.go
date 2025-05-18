// Package handlers provides HTTP handlers for the AllMiTools server
package handlers

import (
	"fmt"
	"net/http"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404
	w.WriteHeader(http.StatusNotFound)
	
	// For now, just return a simple message
	// This will be replaced with template rendering in the future
	fmt.Fprintf(w, "404 - Page Not Found: %s", r.URL.Path)
}
