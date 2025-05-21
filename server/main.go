// Package main is the entry point for the AllMiTools server application
package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/CJFEdu/allmitools/server/internal/handlers"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/gorilla/mux"
)

// serverConfig holds the configuration for the server
type serverConfig struct {
	Port        int
	TemplatesDir string
}

// newRouter creates and configures a new router with all the routes
func newRouter(config serverConfig) *mux.Router {
	// Initialize the router
	r := mux.NewRouter()

	// Register routes
	// Homepage route
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Documentation routes
	r.HandleFunc("/docs", handlers.DocsBaseHandler).Methods("GET")
	r.HandleFunc("/docs/", handlers.DocsBaseHandler).Methods("GET")
	r.HandleFunc("/docs/{tool_name}", handlers.DocsToolHandler).Methods("GET")

	// Tools routes
	r.HandleFunc("/tools/{tool_name}", handlers.ToolsHandler).Methods("GET", "POST")

	// Set custom 404 handler
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	return r
}



func main() {
	// Define the server configuration
	config := serverConfig{
		Port:        8080,
		TemplatesDir: filepath.Join("templates"),
	}

	// Initialize the template manager
	log.Println("Initializing template manager...")
	if err := templates.Initialize("."); err != nil {
		log.Fatalf("Error initializing template manager: %v", err)
	}

	// Create and configure the router
	router := newRouter(config)

	// Start the server
	serverAddr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Server starting on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
