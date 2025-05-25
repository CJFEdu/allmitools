// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/database"
	"github.com/CJFEdu/allmitools/server/internal/logging"
	"github.com/CJFEdu/allmitools/server/internal/middleware"
)

// CleanupResult represents the result of a cleanup operation
type CleanupResult struct {
	Success            bool      `json:"success"`
	Message            string    `json:"message"`
	TextEntriesRemoved int64     `json:"text_entries_removed"`
	LogEntriesRemoved  int64     `json:"log_entries_removed"`
	Timestamp          time.Time `json:"timestamp"`
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
	textEntriesRemoved, err := dao.DeleteExpiredEntries(7 * 24 * time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to clean up text entries: %v", err),
		})
		return
	}

	// Clean up request logs (older than 7 days)
	var logEntriesRemoved int64 = 0
	logDao, err := logging.GetRequestLogDAO()
	if err != nil {
		log.Printf("Warning: Failed to get request log DAO: %v", err)
	} else {
		logEntriesRemoved, err = logDao.DeleteOldRequestLogs(7)
		if err != nil {
			log.Printf("Warning: Failed to clean up request logs: %v", err)
		}
	}

	// Log the cleanup operation
	log.Printf("Database cleanup completed: %d expired text entries and %d request logs removed", 
		textEntriesRemoved, logEntriesRemoved)

	// Create the result
	result := CleanupResult{
		Success:            true,
		Message:            fmt.Sprintf("Successfully removed %d expired text entries and %d request logs", textEntriesRemoved, logEntriesRemoved),
		TextEntriesRemoved: textEntriesRemoved,
		LogEntriesRemoved:  logEntriesRemoved,
		Timestamp:          time.Now(),
	}

	// Return the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ScheduledDatabaseCleanup performs a scheduled cleanup of the database
// This function can be called periodically by a goroutine
func ScheduledDatabaseCleanup() {
	// Clean up expired text entries
	cleanupTextEntries()

	// Clean up old request logs
	cleanupRequestLogs()
}

// cleanupTextEntries removes expired text entries from the database
func cleanupTextEntries() {
	// Get the text storage DAO
	dao, err := database.GetTextStorageDAO()
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to get database connection: %v", err)
		return
	}

	// Delete expired text entries (older than 7 days with save_flag=false)
	entriesRemoved, err := dao.DeleteExpiredEntries(7 * 24 * time.Hour)
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to clean up text entries: %v", err)
		return
	}

	// Log the cleanup operation
	log.Printf("Scheduled text entries cleanup completed: %d expired entries removed", entriesRemoved)
}

// cleanupRequestLogs removes old request logs from the database
func cleanupRequestLogs() {
	// Get the request log DAO
	logDao, err := logging.GetRequestLogDAO()
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to get request log DAO: %v", err)
		return
	}

	// Delete request logs older than 7 days
	logsRemoved, err := logDao.DeleteOldRequestLogs(7)
	if err != nil {
		log.Printf("Scheduled cleanup error: Failed to clean up request logs: %v", err)
		return
	}

	// Log the cleanup operation
	log.Printf("Scheduled request logs cleanup completed: %d logs removed", logsRemoved)
}
