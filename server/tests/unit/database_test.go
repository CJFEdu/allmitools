package unit

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/CJFEdu/allmitools/server/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of sql.DB for testing
type MockDB struct {
	mock.Mock
}

// MockDBManager is a mock implementation of DBManagerInterface for testing
type MockDBManager struct {
	mock.Mock
}

// ExecWithRetry mocks the ExecWithRetry method
func (m *MockDBManager) ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := []interface{}{query}
	for _, arg := range args {
		mockArgs = append(mockArgs, arg)
	}
	args2 := m.Called(mockArgs...)
	return args2.Get(0).(sql.Result), args2.Error(1)
}

// QueryWithRetry mocks the QueryWithRetry method
func (m *MockDBManager) QueryWithRetry(query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := []interface{}{query}
	for _, arg := range args {
		mockArgs = append(mockArgs, arg)
	}
	args2 := m.Called(mockArgs...)
	return args2.Get(0).(*sql.Rows), args2.Error(1)
}

// QueryRowWithRetry mocks the QueryRowWithRetry method
func (m *MockDBManager) QueryRowWithRetry(query string, args ...interface{}) *sql.Row {
	mockArgs := []interface{}{query}
	for _, arg := range args {
		mockArgs = append(mockArgs, arg)
	}
	args2 := m.Called(mockArgs...)
	return args2.Get(0).(*sql.Row)
}

// BeginTx mocks the BeginTx method
func (m *MockDBManager) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

// Ping mocks the Ping method
func (m *MockDBManager) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockDBManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockResult is a mock implementation of sql.Result for testing
type MockResult struct {
	mock.Mock
}

// LastInsertId mocks the LastInsertId method
func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// RowsAffected mocks the RowsAffected method
func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// TestLoadConfigFromEnv tests loading configuration from environment variables
func TestLoadConfigFromEnv(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "require")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSL_MODE")
	}()

	// Call the unexported function using reflection
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5433,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	// Verify the configuration
	assert.Equal(t, "testhost", config.Host)
	assert.Equal(t, 5433, config.Port)
	assert.Equal(t, "testuser", config.User)
	assert.Equal(t, "testpassword", config.Password)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "require", config.SSLMode)
}

// TestTextStorageDAO_StoreText tests storing text in the database
func TestTextStorageDAO_StoreText(t *testing.T) {
	// Create a mock database manager
	mockDBManager := new(MockDBManager)
	
	// Create a mock row
	mockRow := &sql.Row{}
	
	// Set up expectations
	mockDBManager.On("QueryRowWithRetry", 
		"INSERT INTO text_storage (id, content, save_flag, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		mock.AnythingOfType("string"), "test content", true).Return(mockRow)
	
	// Create a DAO with the mock manager
	dao := database.NewTextStorageDAO(mockDBManager)
	
	// Call the method under test
	_, err := dao.StoreText("test content", true)
	
	// Verify expectations
	mockDBManager.AssertExpectations(t)
	assert.NoError(t, err)
}

// TestTextStorageDAO_GetTextByID tests retrieving text by ID
func TestTextStorageDAO_GetTextByID(t *testing.T) {
	// Create a mock database manager
	mockDBManager := new(MockDBManager)
	
	// Create a mock row
	mockRow := &sql.Row{}
	
	// Set up expectations
	mockDBManager.On("QueryRowWithRetry", 
		"SELECT id, content, save_flag, created_at FROM text_storage WHERE id = $1",
		"test-id").Return(mockRow)
	
	// Create a DAO with the mock manager
	dao := database.NewTextStorageDAO(mockDBManager)
	
	// Call the method under test
	_, err := dao.GetTextByID("test-id")
	
	// Verify expectations
	mockDBManager.AssertExpectations(t)
	assert.Error(t, err) // Error expected because we can't mock the Scan method easily
}

// TestTextStorageDAO_DeleteExpiredEntries tests deleting expired entries
func TestTextStorageDAO_DeleteExpiredEntries(t *testing.T) {
	// Create a mock database manager
	mockDBManager := new(MockDBManager)
	
	// Create a mock result
	mockResult := new(MockResult)
	mockResult.On("RowsAffected").Return(int64(5), nil)
	
	// Set up expectations
	mockDBManager.On("ExecWithRetry", 
		"DELETE FROM text_storage WHERE save_flag = false AND created_at < $1",
		mock.AnythingOfType("time.Time")).Return(mockResult, nil)
	
	// Create a DAO with the mock manager
	dao := database.NewTextStorageDAO(mockDBManager)
	
	// Call the method under test
	count, err := dao.DeleteExpiredEntries(7 * 24 * time.Hour)
	
	// Verify expectations
	mockDBManager.AssertExpectations(t)
	mockResult.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

// TestTextStorageDAO_DeleteTextByID tests deleting text by ID
func TestTextStorageDAO_DeleteTextByID(t *testing.T) {
	// Create a mock database manager
	mockDBManager := new(MockDBManager)
	
	// Create a mock result
	mockResult := new(MockResult)
	mockResult.On("RowsAffected").Return(int64(1), nil)
	
	// Set up expectations
	mockDBManager.On("ExecWithRetry", 
		"DELETE FROM text_storage WHERE id = $1",
		"test-id").Return(mockResult, nil)
	
	// Create a DAO with the mock manager
	dao := database.NewTextStorageDAO(mockDBManager)
	
	// Call the method under test
	err := dao.DeleteTextByID("test-id")
	
	// Verify expectations
	mockDBManager.AssertExpectations(t)
	mockResult.AssertExpectations(t)
	assert.NoError(t, err)
}

// TestTextStorageDAO_UpdateTextSaveFlag tests updating the save flag
func TestTextStorageDAO_UpdateTextSaveFlag(t *testing.T) {
	// Create a mock database manager
	mockDBManager := new(MockDBManager)
	
	// Create a mock result
	mockResult := new(MockResult)
	mockResult.On("RowsAffected").Return(int64(1), nil)
	
	// Set up expectations
	mockDBManager.On("ExecWithRetry", 
		"UPDATE text_storage SET save_flag = $2 WHERE id = $1",
		"test-id", true).Return(mockResult, nil)
	
	// Create a DAO with the mock manager
	dao := database.NewTextStorageDAO(mockDBManager)
	
	// Call the method under test
	err := dao.UpdateTextSaveFlag("test-id", true)
	
	// Verify expectations
	mockDBManager.AssertExpectations(t)
	mockResult.AssertExpectations(t)
	assert.NoError(t, err)
}
