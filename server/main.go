// Package main is the entry point for the AllMiTools server application
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/CJFEdu/allmitools/server/internal/handlers"
	"github.com/CJFEdu/allmitools/server/internal/middleware"
	"github.com/CJFEdu/allmitools/server/internal/templates"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// serverConfig holds the configuration for the server
type serverConfig struct {
	Port         int
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

	// Authentication routes
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Private tools routes (protected by auth middleware)
	privateRouter := r.PathPrefix("/private").Subrouter()
	privateRouter.Use(middleware.AuthMiddleware)

	// Private tools listing
	privateRouter.HandleFunc("/tools", handlers.PrivateToolsListHandler).Methods("GET")
	privateRouter.HandleFunc("/tools/", handlers.PrivateToolsListHandler).Methods("GET")

	// Private tool execution
	privateRouter.HandleFunc("/tools/{tool_name}", handlers.PrivateToolsHandler).Methods("GET", "POST")

	// Private documentation
	privateRouter.HandleFunc("/docs", handlers.PrivateDocsBaseHandler).Methods("GET")
	privateRouter.HandleFunc("/docs/", handlers.PrivateDocsBaseHandler).Methods("GET")
	privateRouter.HandleFunc("/docs/{tool_name}", handlers.PrivateDocsToolHandler).Methods("GET")

	// Set custom 404 handler
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	return r
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values or environment variables")
	} else {
		log.Println("Loaded configuration from .env file")
	}
}

// getEnvInt gets an integer environment variable or returns the default value
func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %d\n", key, defaultVal)
		return defaultVal
	}
	
	return val
}

// getEnvString gets a string environment variable or returns the default value
func getEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func main() {
	// Load environment variables
	loadEnv()
	
	// Define the server configuration
	config := serverConfig{
		Port:         getEnvInt("PORT", 3000),
		TemplatesDir: getEnvString("TEMPLATES_DIR", "templates"),
	}
	
	log.Printf("Using configuration: Port=%d, TemplatesDir=%s\n", config.Port, config.TemplatesDir)

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
