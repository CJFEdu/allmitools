// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/database"
	"github.com/CJFEdu/allmitools/server/internal/middleware"
)

// CleanupResult represents the result of a cleanup operation
type CleanupResult struct {
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	EntriesRemoved int64    `json:"entries_removed"`
	Timestamp     time.Time `json:"timestamp"`
}

// DatabaseCleanupHandler handles requests to clean up the database
// This handler is protected by the auth middleware
func DatabaseCleanupHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed. Use POST to trigger cleanup.",
		})
		return
	}

	// Check if user is authenticated (should be handled by middleware, but double-check)
	if !middleware.IsAuthenticated(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Authentication required for database cleanup.",
		})
		return
	}

	// Get the text storage DAO
	dao, err := database.GetTextStorageDAO()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get database connection: %v", err),
		})
		return
	}

	// Delete expired text entries (older than 7 days with save_flag=false)
	entriesRemoved, err := dao.DeleteExpiredEntries(7 * 24 * time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to clean up database: %v", err),
		})
		return
	}

	// Log the cleanup operation
	log.Printf("Database cleanup completed: %d expired text entries removed", entriesRemoved)

	// Create the result
	result := CleanupResult{
		Success:       true,
		Message:       fmt.Sprintf("Successfully removed %d expired text entries", entriesRemoved),
		EntriesRemoved: entriesRemoved,
		Timestamp:     time.Now(),
	}

	// Return the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ScheduledDatabaseCleanup performs a scheduled cleanup of the database
// This function can be called periodically by a goroutine
func ScheduledDatabaseCleanup() {
	// Get the text storage DAO
	dao, err := database.GetTextStorageDAO()
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to get database connection: %v", err)
		return
	}

	// Delete expired text entries (older than 7 days with save_flag=false)
	entriesRemoved, err := dao.DeleteExpiredEntries(7 * 24 * time.Hour)
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to clean up database: %v", err)
		return
	}

	// Log the cleanup operation
	log.Printf("Scheduled database cleanup completed: %d expired text entries removed", entriesRemoved)
}
