// Package logging contains functionality for logging HTTP requests
package logging

import (
	"fmt"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/database"
)

// RequestLog represents a log entry for an HTTP request
type RequestLog struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	ContentType    string    `json:"content_type,omitempty"`
	RequestBody    string    `json:"request_body,omitempty"`
	QueryParams    string    `json:"query_params,omitempty"`
	ResponseStatus int       `json:"response_status"`
	ResponseTimeMs int       `json:"response_time_ms"`
	UserAgent      string    `json:"user_agent,omitempty"`
	IPAddress      string    `json:"ip_address,omitempty"`
}

// RequestLogDAO provides database operations for request logs
type RequestLogDAO struct {
	dbManager database.DBManagerInterface
}

// NewRequestLogDAO creates a new RequestLogDAO with the given database manager
func NewRequestLogDAO(dbManager database.DBManagerInterface) (*RequestLogDAO, error) {
	if dbManager == nil {
		return nil, fmt.Errorf("database manager cannot be nil")
	}
	return &RequestLogDAO{dbManager: dbManager}, nil
}

// GetRequestLogDAO returns a RequestLogDAO using the default database manager
func GetRequestLogDAO() (*RequestLogDAO, error) {
	dbManager, err := database.GetManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get database manager: %w", err)
	}
	return NewRequestLogDAO(dbManager)
}

// InsertRequestLog inserts a new request log entry into the database
func (dao *RequestLogDAO) InsertRequestLog(log *RequestLog) error {
	// Use the query params as is, since it's already a string
	queryParamsJSON := log.QueryParams

	// Prepare the SQL statement
	query := `
		INSERT INTO request_logs (
			timestamp, endpoint, method, content_type, request_body, 
			query_params, response_status, response_time_ms, user_agent, ip_address
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id
	`

	// Execute the query
	var id string
	err := dao.dbManager.QueryRowWithRetry(
		query,
		log.Timestamp,
		log.Endpoint,
		log.Method,
		log.ContentType,
		log.RequestBody,
		queryParamsJSON,
		log.ResponseStatus,
		log.ResponseTimeMs,
		log.UserAgent,
		log.IPAddress,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to insert request log: %w", err)
	}

	// Set the ID in the log object
	log.ID = id
	return nil
}

// GetRequestLogs retrieves request logs with optional filtering
func (dao *RequestLogDAO) GetRequestLogs(limit int, offset int) ([]RequestLog, error) {
	// Prepare the SQL statement
	query := `
		SELECT 
			id, timestamp, endpoint, method, content_type, 
			request_body, query_params, response_status, response_time_ms, 
			user_agent, ip_address
		FROM request_logs
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`

	// Execute the query
	rows, err := dao.dbManager.QueryWithRetry(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query request logs: %w", err)
	}
	defer rows.Close()

	// Parse the results
	logs := []RequestLog{}
	for rows.Next() {
		var log RequestLog
		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Endpoint,
			&log.Method,
			&log.ContentType,
			&log.RequestBody,
			&log.QueryParams,
			&log.ResponseStatus,
			&log.ResponseTimeMs,
			&log.UserAgent,
			&log.IPAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request log row: %w", err)
		}
		logs = append(logs, log)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over request log rows: %w", err)
	}

	return logs, nil
}

// DeleteOldRequestLogs deletes request logs older than the specified number of days
func (dao *RequestLogDAO) DeleteOldRequestLogs(days int) (int64, error) {
	// Prepare the SQL statement
	query := `
		DELETE FROM request_logs
		WHERE timestamp < NOW() - INTERVAL '$1 days'
	`

	// Execute the query
	result, err := dao.dbManager.ExecWithRetry(query, days)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old request logs: %w", err)
	}

	// Get the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows count: %w", err)
	}

	return rowsAffected, nil
}

// CountRequestLogs returns the total number of request logs
func (dao *RequestLogDAO) CountRequestLogs() (int, error) {
	// Prepare the SQL statement
	query := `SELECT COUNT(*) FROM request_logs`

	// Execute the query
	var count int
	err := dao.dbManager.QueryRowWithRetry(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count request logs: %w", err)
	}

	return count, nil
}
