# AllMiTools Server

This is the server component of the AllMiTools project. It provides a Go-powered website for no-code automation tools using Gorilla/Mux for routing.

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8.svg)](https://golang.org/)
[![Gorilla Mux](https://img.shields.io/badge/Gorilla_Mux-1.8.0-blue.svg)](https://github.com/gorilla/mux)

## Project Structure

```
/server
|-- /internal
|   |-- /handlers           // HTTP handlers for different routes
|   |   |-- home.go         // Homepage handler
|   |   |-- docs.go         // Documentation handlers
|   |   |-- tools.go        // Tools handler
|   |   |-- errors.go       // For custom error handlers like 404
|   |-- /models             // Data structures for the application
|   |   |-- tool.go         // Tool models and validation
|-- /templates              // HTML templates (future implementation)
|-- /tests                  // Test files
|   |-- /unit               // Unit tests
|   |   |-- main_test.go    // Tests for handlers
|   |   |-- models_test.go  // Tests for models
|   |   |-- router_test.go  // Tests for router configuration
|   |   |-- utils_test.go   // Tests for test utilities
|   |-- /integration        // Integration tests (future implementation)
|   |-- /utils              // Test utilities
|   |   |-- http_test_helpers.go  // HTTP test helpers
|   |   |-- test_utils.go         // General test utilities
|   |-- /mocks              // Mock data for testing
|   |-- /config             // Test configuration
|-- go.mod                  // Go module definition
|-- go.sum                  // Go module checksums
|-- main.go                 // Main application entry point
```

## Getting Started

### Prerequisites
- Go 1.18 or higher
- Dependencies:
  - github.com/gorilla/mux v1.8.0
  - github.com/stretchr/testify v1.8.4 (for testing)

### Running the server
```bash
cd server
go run main.go
```

The server will start on port 8080 by default.

## Website Sections

1. **Homepage** (`/`) - Introduction and list of available tools
2. **Documentation Section** (`/docs/...`) - Details about each tool
3. **Tools Section** (`/tools/...`) - Access to the actual tools

### Available Tools

The server currently includes the following tools:

1. **Random Number Generator** (`/tools/random-number`) - Generates a random number within a specified range
   - Parameters: `min` (default: 1), `max` (default: 100)

2. **Text Formatter** (`/tools/text-formatter`) - Formats text with various options
   - Parameters: `text` (required), `uppercase` (default: false), `lowercase` (default: false)

### Output Formats

Each tool supports multiple output formats:

- **HTML** - For browser viewing
- **JSON** - For API integration
- **Raw** - Plain text output

The output format is determined by the tool's configuration and can be influenced by the client's Accept header.
