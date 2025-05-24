// Package main is the entry point for the AllMiTools server application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/database"
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

	// Database maintenance routes (protected by auth middleware)
	privateRouter.HandleFunc("/maintenance/cleanup", handlers.DatabaseCleanupHandler).Methods("POST")

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

// scheduleCleanup runs the database cleanup task on a schedule
func scheduleCleanup() {
	cleanupInterval := 24 * time.Hour // Run once per day
	cleanupTicker := time.NewTicker(cleanupInterval)
	defer cleanupTicker.Stop()

	// Run an initial cleanup on startup
	log.Println("Running initial database cleanup...")
	handlers.ScheduledDatabaseCleanup()

	// Then run on the schedule
	for {
		select {
		case <-cleanupTicker.C:
			log.Println("Running scheduled database cleanup...")
			handlers.ScheduledDatabaseCleanup()
		}
	}
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

	// Initialize the database connection
	log.Println("Initializing database connection...")
	if err := database.Initialize(); err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}

	// Start scheduled database cleanup
	go scheduleCleanup()

	// Create and configure the router
	router := newRouter(config)

	// Create a new server with a timeout
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on http://localhost%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-sigChan

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Close database connection
	log.Println("Closing database connection...")
	if err := database.Shutdown(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server gracefully stopped")
}
