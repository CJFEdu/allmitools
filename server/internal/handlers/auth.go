// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"fmt"
	"net/http"

	"github.com/CJFEdu/allmitools/server/internal/middleware"
	"github.com/CJFEdu/allmitools/server/internal/templates"
)

// LoginHandler handles requests to the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if this is a POST request (login attempt)
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Get the password from the form
		password := r.FormValue("password")
		
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
