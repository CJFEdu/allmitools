// Package middleware contains HTTP middleware for the AllMiTools server
package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	// CookieName is the name of the authentication cookie
	CookieName = "allmitools_auth"
	// CookieMaxAge is the maximum age of the authentication cookie in seconds (24 hours)
	CookieMaxAge = 86400
)

var (
	// Initialize secure cookie with random keys
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
)

// AuthMiddleware is middleware that checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated via cookie
		if IsAuthenticated(r) {
			// User is authenticated, proceed to the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Check if the user provided a password in the request
		password := ""
		if r.Method == http.MethodPost {
			// Check Content-Type header to determine how to parse the data
			contentType := r.Header.Get("Content-Type")

			// If it's a form submission, parse form data
			if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
			   strings.Contains(contentType, "multipart/form-data") {
				if err := r.ParseForm(); err == nil {
					password = r.FormValue("password")
				}
			} else if strings.Contains(contentType, "application/json") {
				// Parse JSON data
				var loginData struct {
					Password string `json:"password"`
				}

				// Limit request body size to prevent DoS attacks
				body, err := io.ReadAll(io.LimitReader(r.Body, 1024))
				if err == nil {
					// Restore the request body so it can be read by other handlers
					r.Body.Close()
					r.Body = io.NopCloser(bytes.NewBuffer(body))

					// Decode JSON
					if err := json.Unmarshal(body, &loginData); err == nil {
						password = loginData.Password
					}
				}
			} else {
				// Default to form parsing for backward compatibility
				if err := r.ParseForm(); err == nil {
					password = r.FormValue("password")
				}
			}
		} else {
			password = r.URL.Query().Get("password")
		}

		// If password is provided, verify it
		if password != "" && VerifyPassword(password) {
			// Password is correct, set cookie and proceed
			SetAuthCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		// User is not authenticated, redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

// IsAuthenticated checks if the user is authenticated via cookie
func IsAuthenticated(r *http.Request) bool {
	// Get the cookie
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return false
	}

	// Decode the cookie value
	value := make(map[string]string)
	if err = cookieHandler.Decode(CookieName, cookie.Value, &value); err != nil {
		return false
	}

	// Check if the authenticated flag is set
	return value["authenticated"] == "true"
}

// SetAuthCookie sets the authentication cookie
func SetAuthCookie(w http.ResponseWriter) {
	// Create a map to store in the cookie
	value := map[string]string{
		"authenticated": "true",
		"timestamp":     fmt.Sprintf("%d", time.Now().Unix()),
	}

	// Encode the cookie value
	if encoded, err := cookieHandler.Encode(CookieName, value); err == nil {
		// Create a new cookie
		cookie := &http.Cookie{
			Name:     CookieName,
			Value:    encoded,
			Path:     "/",
			MaxAge:   CookieMaxAge,
			HttpOnly: true,
			Secure:   true, // Set to true in production
			SameSite: http.SameSiteStrictMode,
		}

		// Set the cookie
		http.SetCookie(w, cookie)
	}
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(w http.ResponseWriter) {
	// Create a new cookie with negative MaxAge to clear it
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}

	// Set the cookie
	http.SetCookie(w, cookie)
}

// VerifyPassword verifies the provided password against the hash in the .env file
func VerifyPassword(password string) bool {
	// Get the password hash from the environment
	storedHash := os.Getenv("PRIVATE_USE_PASSWORD")
	if storedHash == "" {
		// If no password hash is set, authentication fails
		return false
	}

	// Hash the provided password
	hasher := sha256.New()
	hasher.Write([]byte(password))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	// Compare the hashes
	return computedHash == storedHash
}

// HashPassword generates a SHA-256 hash for a password
// This is a utility function that can be used to generate a hash for the .env file
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
