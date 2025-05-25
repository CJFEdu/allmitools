// Package handlers contains HTTP handlers for the AllMiTools server
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CJFEdu/allmitools/server/internal/models"
	"github.com/CJFEdu/allmitools/server/internal/templates"
)

// PrivateDocsBaseHandler handles requests to the private documentation base page
func PrivateDocsBaseHandler(w http.ResponseWriter, r *http.Request) {
	// Get all private tools
	privateTools := models.GetAllPrivateTools()

	// Prepare data for the template
	data := map[string]interface{}{
		"Title":        "Private Tools Documentation",
		"CurrentPage":  "private-docs",
		"PrivateTools": privateTools,
	}

	// Render the template
	err := templates.TemplateManager.RenderTemplate(w, "private_docs_base", data)
	if err != nil {
		// Fallback if template rendering fails
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Private Tools Documentation</h1>")
		fmt.Fprintf(w, "<p>Select a tool from the list below to view its documentation:</p>")
		fmt.Fprintf(w, "<ul>")
		for _, tool := range privateTools {
			fmt.Fprintf(w, "<li><a href='/private/docs/%s'>%s</a> - %s</li>", tool.Name, tool.Name, tool.Description)
		}
		fmt.Fprintf(w, "</ul>")
		fmt.Fprintf(w, "</body></html>")
	}
}

// PrivateDocsToolHandler handles requests to view documentation for a specific private tool
func PrivateDocsToolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["tool_name"]

	// Get tool info
	toolInfo, err := models.GetPrivateToolInfo(toolName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		
		// Render the 404 template
		data := map[string]interface{}{
			"Title":       "Documentation Not Found",
			"CurrentPage": "private-docs",
		}
		
		err := templates.TemplateManager.RenderTemplate(w, "404", data)
		if err != nil {
			// Fallback if template rendering fails
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>")
			fmt.Fprintf(w, "<h1>Documentation Not Found</h1>")
			fmt.Fprintf(w, "<p>The documentation for private tool '%s' was not found.</p>", toolName)
			fmt.Fprintf(w, "<p><a href='/private/docs'>Back to documentation</a></p>")
			fmt.Fprintf(w, "</body></html>")
		}
		return
	}

	// Prepare data for the template
	data := map[string]interface{}{
		"Title":       fmt.Sprintf("%s Documentation", toolInfo.Name),
		"CurrentPage": "private-docs",
		"Tool":        toolInfo,
		"IsPrivate":   true,
	}

	// Render the template
	err = templates.TemplateManager.RenderTemplate(w, "tool_docs", data)
	if err != nil {
		// Fallback if template rendering fails
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>%s Documentation</h1>", toolInfo.Name)
		fmt.Fprintf(w, "<p>%s</p>", toolInfo.Description)
		fmt.Fprintf(w, "<h2>Parameters</h2>")
		if len(toolInfo.Parameters) > 0 {
			fmt.Fprintf(w, "<ul>")
			for _, param := range toolInfo.Parameters {
				fmt.Fprintf(w, "<li><strong>%s</strong> (%s): %s</li>", param.Name, param.Type, param.Description)
			}
			fmt.Fprintf(w, "</ul>")
		} else {
			fmt.Fprintf(w, "<p>This tool does not have any parameters.</p>")
		}
		fmt.Fprintf(w, "<p><a href='/private/tools/%s'>Use this tool</a></p>", toolInfo.Name)
		fmt.Fprintf(w, "</body></html>")
	}
}
