// Package database provides functionality for database operations
package database

import (
	"log"
	"sync"
)

var (
	// Global database manager instance
	dbManager *DBManager
	// Mutex to ensure thread-safe initialization
	initMutex sync.Mutex
	// Flag to track initialization status
	initialized bool
)

// Initialize initializes the database connection
// This should be called once during application startup
func Initialize() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	log.Println("Initializing database connection...")
	manager, err := NewManager()
	if err != nil {
		return err
	}

	dbManager = manager
	initialized = true
	log.Println("Database connection initialized successfully")
	return nil
}

// GetManager returns the global database manager instance
// Initializes the connection if not already initialized
func GetManager() (*DBManager, error) {
	initMutex.Lock()
	defer initMutex.Unlock()

	if !initialized {
		manager, err := NewManager()
		if err != nil {
			return nil, err
		}
		dbManager = manager
		initialized = true
	}

	return dbManager, nil
}

// Shutdown closes the database connection
// This should be called during application shutdown
func Shutdown() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if !initialized || dbManager == nil {
		return nil
	}

	log.Println("Shutting down database connection...")
	err := dbManager.Close()
	if err != nil {
		return err
	}

	initialized = false
	dbManager = nil
	log.Println("Database connection closed successfully")
	return nil
}

// GetTextStorageDAO returns a new TextStorageDAO instance
// using the global database manager
func GetTextStorageDAO() (*TextStorageDAO, error) {
	manager, err := GetManager()
	if err != nil {
		return nil, err
	}
	return NewTextStorageDAO(manager), nil
}
