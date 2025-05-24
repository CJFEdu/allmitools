// Package database provides functionality for database operations
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// Import PostgreSQL driver
	_ "github.com/lib/pq"
)

// DBManagerInterface defines the interface for database operations
type DBManagerInterface interface {
	ExecWithRetry(query string, args ...interface{}) (sql.Result, error)
	QueryWithRetry(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowWithRetry(query string, args ...interface{}) *sql.Row
	BeginTx() (*sql.Tx, error)
	Ping() error
	Close() error
}

// DBManager manages database connections and operations
type DBManager struct {
	DB           *sql.DB
	MaxRetries   int
	RetryBackoff time.Duration
}

// Config holds database configuration parameters
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewManager creates a new database manager with connection pooling
func NewManager() (*DBManager, error) {
	// Load configuration from environment variables
	config := loadConfigFromEnv()

	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Create manager
	manager := &DBManager{
		DB:           db,
		MaxRetries:   3,
		RetryBackoff: 100 * time.Millisecond,
	}

	// Test connection
	if err := manager.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return manager, nil
}

// loadConfigFromEnv loads database configuration from environment variables
func loadConfigFromEnv() Config {
	port, err := strconv.Atoi(getEnvWithDefault("DB_PORT", "5432"))
	if err != nil {
		port = 5432
	}

	return Config{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnvWithDefault("DB_USER", "allmitools_user"),
		Password: getEnvWithDefault("DB_PASSWORD", ""),
		DBName:   getEnvWithDefault("DB_NAME", "allmitools"),
		SSLMode:  getEnvWithDefault("DB_SSL_MODE", "disable"),
	}
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Ping tests the database connection
func (m *DBManager) Ping() error {
	var err error
	for i := 0; i < m.MaxRetries; i++ {
		err = m.DB.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Database ping attempt %d failed: %v. Retrying in %v...", i+1, err, m.RetryBackoff)
		time.Sleep(m.RetryBackoff)
		// Exponential backoff
		m.RetryBackoff *= 2
	}
	return fmt.Errorf("failed to ping database after %d attempts: %w", m.MaxRetries, err)
}

// Close closes the database connection
func (m *DBManager) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return nil
}

// ExecWithRetry executes a query with retry logic
func (m *DBManager) ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	var (
		result sql.Result
		err    error
		retry  = 0
		backoff = m.RetryBackoff
	)

	for retry < m.MaxRetries {
		result, err = m.DB.Exec(query, args...)
		if err == nil {
			return result, nil
		}

		log.Printf("Database exec attempt %d failed: %v. Retrying in %v...", retry+1, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
		retry++
	}

	return nil, fmt.Errorf("failed to execute query after %d attempts: %w", m.MaxRetries, err)
}

// QueryWithRetry executes a query with retry logic
func (m *DBManager) QueryWithRetry(query string, args ...interface{}) (*sql.Rows, error) {
	var (
		rows    *sql.Rows
		err     error
		retry   = 0
		backoff = m.RetryBackoff
	)

	for retry < m.MaxRetries {
		rows, err = m.DB.Query(query, args...)
		if err == nil {
			return rows, nil
		}

		log.Printf("Database query attempt %d failed: %v. Retrying in %v...", retry+1, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
		retry++
	}

	return nil, fmt.Errorf("failed to execute query after %d attempts: %w", m.MaxRetries, err)
}

// QueryRowWithRetry executes a query that returns a single row with retry logic
func (m *DBManager) QueryRowWithRetry(query string, args ...interface{}) *sql.Row {
	// QueryRow doesn't return an error, so we can't retry on error
	// We'll just return the result directly
	return m.DB.QueryRow(query, args...)
}

// BeginTx starts a new transaction
func (m *DBManager) BeginTx() (*sql.Tx, error) {
	var (
		tx      *sql.Tx
		err     error
		retry   = 0
		backoff = m.RetryBackoff
	)

	for retry < m.MaxRetries {
		tx, err = m.DB.Begin()
		if err == nil {
			return tx, nil
		}

		log.Printf("Database begin transaction attempt %d failed: %v. Retrying in %v...", retry+1, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
		retry++
	}

	return nil, fmt.Errorf("failed to begin transaction after %d attempts: %w", m.MaxRetries, err)
}
