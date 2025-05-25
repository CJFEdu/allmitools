package logging

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
		body:           &bytes.Buffer{},
	}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// RequestLoggerMiddleware is a middleware that logs HTTP requests
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request logging is enabled
		loggingEnabled, _ := strconv.ParseBool(os.Getenv("REQUEST_LOGGING_ENABLED"))
		if !loggingEnabled {
			// If logging is disabled, just call the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Skip logging for certain endpoints
		if shouldSkipLogging(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Start timer
		startTime := time.Now()

		// Create a copy of the request body
		var requestBody string
		if r.Body != nil && r.Method != http.MethodGet {
			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create a custom response writer to capture the status code
		rw := newResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate response time
		duration := time.Since(startTime)
		responseTimeMs := int(duration.Milliseconds())

		// Create a request log entry
		log := &RequestLog{
			Timestamp:      startTime,
			Endpoint:       r.URL.Path,
			Method:         r.Method,
			ContentType:    r.Header.Get("Content-Type"),
			RequestBody:    sanitizeRequestBody(requestBody, r.URL.Path),
			QueryParams:    r.URL.RawQuery,
			ResponseStatus: rw.statusCode,
			ResponseTimeMs: responseTimeMs,
			UserAgent:      r.Header.Get("User-Agent"),
			IPAddress:      getClientIP(r),
		}

		// Save the log entry asynchronously
		go saveRequestLog(log)
	})
}

// shouldSkipLogging determines if logging should be skipped for certain endpoints
func shouldSkipLogging(path string) bool {
	// Skip logging for health check endpoints
	if strings.HasPrefix(path, "/health") {
		return true
	}

	// Skip logging for static files
	if strings.HasPrefix(path, "/static/") {
		return true
	}

	return false
}

// sanitizeRequestBody removes sensitive information from request bodies
func sanitizeRequestBody(body string, path string) string {
	// Skip logging request bodies for authentication endpoints
	if strings.Contains(path, "/auth/") || strings.Contains(path, "/login") {
		return "[REDACTED - AUTH ENDPOINT]"
	}

	// For all other endpoints, return the body as is
	// In a production environment, you might want to implement more sophisticated
	// sanitization to remove sensitive data like passwords, tokens, etc.
	return body
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (common when behind a proxy)
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// The client IP is the first address in the list
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for X-Real-IP header (used by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// saveRequestLog saves a request log entry to the database
func saveRequestLog(reqLog *RequestLog) {
	// Get the request log DAO
	dao, err := GetRequestLogDAO()
	if err != nil {
		log.Printf("Failed to get request log DAO: %v", err)
		return
	}

	// Insert the log entry
	err = dao.InsertRequestLog(reqLog)
	if err != nil {
		log.Printf("Failed to insert request log: %v", err)
	}
}
