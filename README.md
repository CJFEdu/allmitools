# AllMiTools Server

This is the server component of the AllMiTools project. It provides a Go-powered website for no-code automation tools using Gorilla/Mux for routing.

## Project Structure

```
/server
|-- /internal
|   |-- /handlers           // HTTP handlers for different routes
|   |   |-- home.go
|   |   |-- docs.go
|   |   |-- tools.go
|   |   |-- errors.go       // For custom error handlers like 404
|   |-- /models             // Data structures (e.g., for tool info)
|-- /templates              // HTML templates
|   |-- homepage.html
|   |-- docs_page.html
|   |-- tool_output.html
|   |-- 404.html            // Custom 404 page template
|-- /tests                  // Test files
|   |-- /unit               // Unit tests
|   |-- /integration        // Integration tests
|-- go.mod
|-- go.sum
|-- main.go                 // Main application entry point
```

## Getting Started

### Prerequisites
- Go 1.17 or higher
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

1. **Homepage** - Introduction and list of available tools
2. **Documentation Section** (`/docs/...`) - Details about each tool
3. **Tools Section** (`/tools/...`) - Access to the actual tools

Each tool supports multiple output formats:
- HTML (for browser viewing)
- JSON (for API consumption)
- Raw (plain text output)
