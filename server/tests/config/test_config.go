// Package config provides test configuration for the AllMiTools server
package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// TestConfig contains configuration for tests
type TestConfig struct {
	// TemplatesDir is the path to the templates directory
	TemplatesDir string
	
	// TestDataDir is the path to the test data directory
	TestDataDir string
	
	// ServerPort is the port to use for test server
	ServerPort string
}

// DefaultTestConfig returns the default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		TemplatesDir: getTemplatesDir(),
		TestDataDir:  getTestDataDir(),
		ServerPort:   "8081", // Use a different port for tests
	}
}

// getTemplatesDir returns the path to the templates directory
func getTemplatesDir() string {
	// Get the path to the server directory
	_, filename, _, _ := runtime.Caller(0)
	serverDir := filepath.Join(filepath.Dir(filename), "../..")
	return filepath.Join(serverDir, "templates")
}

// getTestDataDir returns the path to the test data directory
func getTestDataDir() string {
	// Get the path to the server directory
	_, filename, _, _ := runtime.Caller(0)
	serverDir := filepath.Join(filepath.Dir(filename), "../..")
	testDataDir := filepath.Join(serverDir, "tests", "data")
	
	// Create the test data directory if it doesn't exist
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		os.MkdirAll(testDataDir, 0755)
	}
	
	return testDataDir
}

// CreateTestDataFile creates a test data file with the given content
func CreateTestDataFile(filename, content string) (string, error) {
	config := DefaultTestConfig()
	filePath := filepath.Join(config.TestDataDir, filename)
	
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	// Write the content
	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}
	
	return filePath, nil
}
