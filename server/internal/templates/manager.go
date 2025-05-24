// Package templates provides template management for the AllMiTools server
package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// Manager handles template loading and rendering
type Manager struct {
	templatesDir string
	templates    map[string]*template.Template
	mutex        sync.RWMutex
}

// NewManager creates a new template manager
func NewManager(templatesDir string) *Manager {
	return &Manager{
		templatesDir: templatesDir,
		templates:    make(map[string]*template.Template),
	}
}

// LoadTemplates loads all templates from the templates directory
func (m *Manager) LoadTemplates() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Define the templates to load
	templateNames := []string{
		"home",
		"docs_base",
		"docs_tool",
		"tool",
		"404",
		"login",
		"private_tools_list",
		"private_docs_base",
	}

	// Load each template
	for _, name := range templateNames {
		// Parse the layout template and the specific template
		tmpl, err := template.ParseFiles(
			filepath.Join(m.templatesDir, "layout.html"),
			filepath.Join(m.templatesDir, name+".html"),
		)
		if err != nil {
			return fmt.Errorf("error loading template %s: %w", name, err)
		}

		// Store the template
		m.templates[name] = tmpl
	}

	return nil
}

// RenderTemplate renders a template with the given data
func (m *Manager) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	m.mutex.RLock()
	tmpl, exists := m.templates[name]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("template %s does not exist", name)
	}

	// Set the content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	return tmpl.Execute(w, data)
}

// GetTemplate returns a template by name
func (m *Manager) GetTemplate(name string) (*template.Template, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tmpl, exists := m.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s does not exist", name)
	}

	return tmpl, nil
}
