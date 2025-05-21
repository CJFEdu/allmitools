// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"testing"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestToolParameterValidation tests the validation of tool parameters
func TestToolParameterValidation(t *testing.T) {
	// Test cases for parameter validation
	testCases := []struct {
		name          string
		parameter     models.ToolParameter
		expectedError string
	}{
		{
			name: "Valid parameter",
			parameter: models.ToolParameter{
				Name:        "test",
				Type:        "string",
				Description: "A test parameter",
				Required:    true,
			},
			expectedError: "",
		},
		{
			name: "Empty name",
			parameter: models.ToolParameter{
				Name:        "",
				Type:        "string",
				Description: "A test parameter",
				Required:    true,
			},
			expectedError: "parameter name cannot be empty",
		},
		{
			name: "Empty type",
			parameter: models.ToolParameter{
				Name:        "test",
				Type:        "",
				Description: "A test parameter",
				Required:    true,
			},
			expectedError: "parameter type cannot be empty",
		},
		{
			name: "Invalid type",
			parameter: models.ToolParameter{
				Name:        "test",
				Type:        "invalid",
				Description: "A test parameter",
				Required:    true,
			},
			expectedError: "invalid parameter type: invalid",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.parameter.Validate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

// TestToolInfoValidation tests the validation of tool info
func TestToolInfoValidation(t *testing.T) {
	// Test cases for tool info validation
	testCases := []struct {
		name          string
		toolInfo      models.ToolInfo
		expectedError string
	}{
		{
			name: "Valid tool info",
			toolInfo: models.ToolInfo{
				Name:        "test-tool",
				Description: "A test tool",
				Version:     "1.0.0",
				Author:      "Test Author",
				Parameters: []models.ToolParameter{
					{
						Name:        "param1",
						Type:        "string",
						Description: "A test parameter",
						Required:    true,
					},
				},
				OutputType: "json",
			},
			expectedError: "",
		},
		{
			name: "Empty name",
			toolInfo: models.ToolInfo{
				Name:        "",
				Description: "A test tool",
				Version:     "1.0.0",
				Author:      "Test Author",
				OutputType:  "json",
			},
			expectedError: "tool name cannot be empty",
		},
		{
			name: "Empty description",
			toolInfo: models.ToolInfo{
				Name:        "test-tool",
				Description: "",
				Version:     "1.0.0",
				Author:      "Test Author",
				OutputType:  "json",
			},
			expectedError: "tool description cannot be empty",
		},
		{
			name: "Empty version",
			toolInfo: models.ToolInfo{
				Name:        "test-tool",
				Description: "A test tool",
				Version:     "",
				Author:      "Test Author",
				OutputType:  "json",
			},
			expectedError: "tool version cannot be empty",
		},
		{
			name: "Invalid output type",
			toolInfo: models.ToolInfo{
				Name:        "test-tool",
				Description: "A test tool",
				Version:     "1.0.0",
				Author:      "Test Author",
				OutputType:  "invalid",
			},
			expectedError: "invalid output type: invalid",
		},
		{
			name: "Invalid parameter",
			toolInfo: models.ToolInfo{
				Name:        "test-tool",
				Description: "A test tool",
				Version:     "1.0.0",
				Author:      "Test Author",
				Parameters: []models.ToolParameter{
					{
						Name:        "",
						Type:        "string",
						Description: "A test parameter",
						Required:    true,
					},
				},
				OutputType: "json",
			},
			expectedError: "parameter 0 (): parameter name cannot be empty",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.toolInfo.Validate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

// TestGetToolInfo tests the GetToolInfo function
func TestGetToolInfo(t *testing.T) {
	// Test cases for GetToolInfo
	testCases := []struct {
		name           string
		toolName       string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "Existing tool",
			toolName:    "random-number",
			expectError: false,
		},
		{
			name:           "Non-existent tool",
			toolName:       "non-existent",
			expectError:    true,
			expectedErrMsg: "tool not found: non-existent",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool, err := models.GetToolInfo(tc.toolName)
			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.toolName, tool.Name)
			}
		})
	}
}

// TestListTools tests the ListTools function
func TestListTools(t *testing.T) {
	// Get the list of tools
	tools := models.ListTools()
	
	// Check that we have the expected number of tools
	assert.Equal(t, len(models.AvailableTools), len(tools))
	
	// Check that all tools are present
	toolMap := make(map[string]bool)
	for _, tool := range tools {
		toolMap[tool.Name] = true
	}
	
	for name := range models.AvailableTools {
		assert.True(t, toolMap[name], "Tool %s should be in the list", name)
	}
}
