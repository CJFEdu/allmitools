// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/CJFEdu/allmitools/server/internal/middleware"
	"github.com/CJFEdu/allmitools/server/internal/templates"
)

// LoginHandler handles requests to the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if this is a POST request (login attempt)
	if r.Method == http.MethodPost {
		// Get the password from the request
		var password string
		
		// Check Content-Type header to determine how to parse the data
		contentType := r.Header.Get("Content-Type")

		// If it's a form submission, parse form data
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || 
		   strings.Contains(contentType, "multipart/form-data") {
			// Parse the form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// Get the password from the form
			password = r.FormValue("password")
		} else if strings.Contains(contentType, "application/json") {
			// Parse JSON data
			var loginData struct {
				Password string `json:"password"`
			}

			// Limit request body size to prevent DoS attacks
			body, err := io.ReadAll(io.LimitReader(r.Body, 1024))
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			// Restore the request body so it can be read by other handlers
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Decode JSON
			if err := json.Unmarshal(body, &loginData); err != nil {
				http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
				return
			}

			// Get the password from JSON
			password = loginData.Password
		} else {
			// Default to form parsing for backward compatibility
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// Get the password from the form
			password = r.FormValue("password")
		}
		
		// Verify the password
		if password != "" && verifyPassword(password) {
			// Password is correct, set cookie and redirect to home page
			middleware.SetAuthCookie(w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Password is incorrect, render login page with error
		data := map[string]interface{}{
			"Title":       "Login",
			"CurrentPage": "login",
			"Error":       "Invalid password. Please try again.",
		}

		// Render the login template
		err := templates.TemplateManager.RenderTemplate(w, "login", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>Login</h1>")
			fmt.Fprintf(w, "<p style='color: red;'>Invalid password. Please try again.</p>")
			fmt.Fprintf(w, "<form method='post' action='/login'>")
			fmt.Fprintf(w, "<label for='password'>Password:</label><br>")
			fmt.Fprintf(w, "<input type='password' id='password' name='password'><br><br>")
			fmt.Fprintf(w, "<input type='submit' value='Login'>")
			fmt.Fprintf(w, "</form>")
			fmt.Fprintf(w, "</body></html>")
		}
		return
	}

	// This is a GET request, render the login page
	data := map[string]interface{}{
		"Title":       "Login",
		"CurrentPage": "login",
	}

	// Render the login template
	err := templates.TemplateManager.RenderTemplate(w, "login", data)
	if err != nil {
		// Fallback if template rendering fails
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Login</h1>")
		fmt.Fprintf(w, "<form method='post' action='/login'>")
		fmt.Fprintf(w, "<label for='password'>Password:</label><br>")
		fmt.Fprintf(w, "<input type='password' id='password' name='password'><br><br>")
		fmt.Fprintf(w, "<input type='submit' value='Login'>")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "</body></html>")
	}
}

// LogoutHandler handles logout requests
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	middleware.ClearAuthCookie(w)
	
	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// verifyPassword verifies the provided password
func verifyPassword(password string) bool {
	// Use the middleware function to verify the password
	return middleware.VerifyPassword(password)
}
