// Package config provides configuration utilities for the AllMiTools server
package config

import (
	"log"
	"os"
)

// Logger is the default logger for the application
var Logger = log.New(os.Stdout, "[AllMiTools] ", log.LstdFlags)
