// Package database provides functionality for database operations
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TextStorageDAO handles database operations for text storage
type TextStorageDAO struct {
	dbManager DBManagerInterface
}

// TextEntry represents a text entry in the database
type TextEntry struct {
	ID        string    // Unique identifier
	Content   string    // Text content
	SaveFlag  bool      // Whether to save permanently
	CreatedAt time.Time // Creation timestamp
}

// NewTextStorageDAO creates a new TextStorageDAO
func NewTextStorageDAO(dbManager DBManagerInterface) *TextStorageDAO {
	return &TextStorageDAO{
		dbManager: dbManager,
	}
}

// StoreText stores text content in the database
// Returns the ID of the stored text
func (dao *TextStorageDAO) StoreText(content string, saveFlag bool) (string, error) {
	// Validate input
	if content == "" {
		return "", errors.New("content cannot be empty")
	}

	// Generate a unique ID
	id := uuid.New().String()

	// Prepare the SQL statement
	query := `
		INSERT INTO text_storage (id, content, save_flag, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	// Execute the query with retry logic
	var returnedID string
	err := dao.dbManager.QueryRowWithRetry(query, id, content, saveFlag).Scan(&returnedID)
	if err != nil {
		return "", fmt.Errorf("failed to store text: %w", err)
	}

	return returnedID, nil
}

// GetTextByID retrieves text content by ID
func (dao *TextStorageDAO) GetTextByID(id string) (*TextEntry, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	// Prepare the SQL statement
	query := `
		SELECT id, content, save_flag, created_at
		FROM text_storage
		WHERE id = $1
	`

	// Execute the query with retry logic
	var entry TextEntry
	err := dao.dbManager.QueryRowWithRetry(query, id).Scan(
		&entry.ID,
		&entry.Content,
		&entry.SaveFlag,
		&entry.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("text entry with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve text: %w", err)
	}

	return &entry, nil
}

// DeleteExpiredEntries deletes unsaved entries older than the specified duration
func (dao *TextStorageDAO) DeleteExpiredEntries(age time.Duration) (int64, error) {
	// Prepare the SQL statement
	query := `
		DELETE FROM text_storage
		WHERE save_flag = false
		AND created_at < $1
	`

	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-age)

	// Execute the query with retry logic
	result, err := dao.dbManager.ExecWithRetry(query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired entries: %w", err)
	}

	// Get the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rowsAffected, nil
}

// DeleteTextByID deletes a text entry by ID
func (dao *TextStorageDAO) DeleteTextByID(id string) error {
	// Validate input
	if id == "" {
		return errors.New("id cannot be empty")
	}

	// Prepare the SQL statement
	query := `
		DELETE FROM text_storage
		WHERE id = $1
	`

	// Execute the query with retry logic
	result, err := dao.dbManager.ExecWithRetry(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete text: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("text entry with ID %s not found", id)
	}

	return nil
}

// UpdateTextSaveFlag updates the save flag for a text entry
func (dao *TextStorageDAO) UpdateTextSaveFlag(id string, saveFlag bool) error {
	// Validate input
	if id == "" {
		return errors.New("id cannot be empty")
	}

	// Prepare the SQL statement
	query := `
		UPDATE text_storage
		SET save_flag = $2
		WHERE id = $1
	`

	// Execute the query with retry logic
	result, err := dao.dbManager.ExecWithRetry(query, id, saveFlag)
	if err != nil {
		return fmt.Errorf("failed to update text save flag: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("text entry with ID %s not found", id)
	}

	return nil
}

// GetAllSavedEntries retrieves all saved text entries
func (dao *TextStorageDAO) GetAllSavedEntries() ([]*TextEntry, error) {
	// Prepare the SQL statement
	query := `
		SELECT id, content, save_flag, created_at
		FROM text_storage
		WHERE save_flag = true
		ORDER BY created_at DESC
	`

	// Execute the query with retry logic
	rows, err := dao.dbManager.QueryWithRetry(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve saved entries: %w", err)
	}
	defer rows.Close()

	// Process the results
	var entries []*TextEntry
	for rows.Next() {
		var entry TextEntry
		err := rows.Scan(
			&entry.ID,
			&entry.Content,
			&entry.SaveFlag,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		entries = append(entries, &entry)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return entries, nil
}
