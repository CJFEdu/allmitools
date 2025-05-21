// Package utils provides testing utilities for the AllMiTools server
package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// APITestCase represents a test case for API endpoints
type APITestCase struct {
	Name           string
	Method         string
	URL            string
	Body           interface{}
	ExpectedStatus int
	ExpectedBody   string
	ExpectedJSON   interface{}
}

// ExecuteRequest executes an HTTP request and returns the response
func ExecuteRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// CheckResponseCode checks if the response code matches the expected code
func CheckResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected response code %d. Got %d", expected, actual)
}

// CreateTestRequest creates a new HTTP request for testing
func CreateTestRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// RunAPITests runs a series of API test cases
func RunAPITests(t *testing.T, router *mux.Router, testCases []APITestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var req *http.Request
			var err error

			req, err = CreateTestRequest(tc.Method, tc.URL, tc.Body)
			assert.NoError(t, err)

			rr := ExecuteRequest(req, router)
			CheckResponseCode(t, tc.ExpectedStatus, rr.Code)

			if tc.ExpectedBody != "" {
				assert.Equal(t, tc.ExpectedBody, rr.Body.String())
			}

			if tc.ExpectedJSON != nil {
				var expected, actual interface{}
				
				// Convert expected JSON to a comparable format
				expectedJSON, err := json.Marshal(tc.ExpectedJSON)
				assert.NoError(t, err)
				err = json.Unmarshal(expectedJSON, &expected)
				assert.NoError(t, err)
				
				// Parse actual response body
				err = json.Unmarshal(rr.Body.Bytes(), &actual)
				assert.NoError(t, err)
				
				// Compare expected and actual JSON
				assert.Equal(t, expected, actual)
			}
		})
	}
}
