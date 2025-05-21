// Package templates provides template management for the AllMiTools server
package templates

import (
	"log"
	"path/filepath"
)

// Global template manager instance
var TemplateManager *Manager

// Initialize initializes the template manager
func Initialize(serverRoot string) error {
	// Create the template manager
	templatesDir := filepath.Join(serverRoot, "templates")
	TemplateManager = NewManager(templatesDir)

	// Load templates
	if err := TemplateManager.LoadTemplates(); err != nil {
		log.Printf("Error loading templates: %v", err)
		return err
	}

	log.Println("Templates loaded successfully")
	return nil
}
