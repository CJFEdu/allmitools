// Package unit contains unit tests for the AllMiTools server
package unit

import (
	"net/http"
	"testing"

	"github.com/CJFEdu/allmitools/server/tests/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestExecuteRequest tests the ExecuteRequest function
func TestExecuteRequest(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()
	
	// Add a test route
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	}).Methods("GET")
	
	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	
	// Execute the request
	rr := utils.ExecuteRequest(req, router)
	
	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Test response", rr.Body.String())
}

// TestCreateTestRequest tests the CreateTestRequest function
func TestCreateTestRequest(t *testing.T) {
	// Test with no body
	req, err := utils.CreateTestRequest("GET", "/test", nil)
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/test", req.URL.String())
	assert.Equal(t, "", req.Header.Get("Content-Type"))
	
	// Test with JSON body
	body := map[string]interface{}{
		"key": "value",
	}
	req, err = utils.CreateTestRequest("POST", "/test", body)
	assert.NoError(t, err)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "/test", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

// TestCheckResponseCode tests the CheckResponseCode function
func TestCheckResponseCode(t *testing.T) {
	// This should not panic
	utils.CheckResponseCode(t, http.StatusOK, http.StatusOK)
	
	// We can't directly test the failure case in a unit test
	// as it would cause the test to fail
}

// TestRunAPITests tests the RunAPITests function
func TestRunAPITests(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()
	
	// Add a test route
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	}).Methods("GET")
	
	// Create test cases
	testCases := []utils.APITestCase{
		{
			Name:           "Test GET /test",
			Method:         "GET",
			URL:            "/test",
			Body:           nil,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   "Test response",
		},
	}
	
	// Run the tests
	utils.RunAPITests(t, router, testCases)
}
