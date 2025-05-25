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
|   |-- /templates          // Template management
|   |   |-- manager.go      // Template manager for loading and rendering templates
|-- /templates              // HTML templates for rendering pages
|   |-- layout.html         // Base layout template
|   |-- home.html           // Homepage template
|   |-- docs_base.html      // Documentation base template
|   |-- docs_tool.html      // Documentation tool template
|   |-- tool.html           // Tool page template
|   |-- 404.html            // Not found page template
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

### Configuration

The server can be configured using environment variables or a `.env` file. An example `.env.example` file is provided in the server directory. Copy this file to `.env` and modify the values as needed.

```bash
# Copy the example .env file
cp server/.env.example server/.env

# Edit the .env file with your preferred settings
```

Available configuration options:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| PORT | Server port | 3000 |
| TEMPLATES_DIR | Templates directory | templates |
| LOG_LEVEL | Logging level | info |
| PRIVATE_USE_PASSWORD | SHA-256 hash of password for private tools | (required for private tools) |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_NAME | PostgreSQL database name | allmitools |
| DB_USER | PostgreSQL username | allmitools_user |
| DB_PASSWORD | PostgreSQL password | (required for database connection) |
| DB_SSL_MODE | PostgreSQL SSL mode | disable |

### Database Setup

#### PostgreSQL Installation

For detailed instructions on installing and configuring PostgreSQL on Ubuntu 24.04, refer to the [PostgresInstall.md](../PostgresInstall.md) guide in the root directory of the repository. This guide includes security best practices and configuration steps.

#### Database Migration

Before running the server for the first time, you need to create the database schema:

```bash
# Navigate to the server directory
cd server

# Run the SQL migration file
psql -U allmitools_user -d allmitools -h localhost -f migrations/001_initial_schema.sql
```

Alternatively, you can run the migration directly from the PostgreSQL command line:

```bash
# Connect to PostgreSQL as the allmitools_user
psql -U allmitools_user -d allmitools -h localhost

# Inside the PostgreSQL prompt, run:
\i migrations/001_initial_schema.sql
```

### Running the server
```bash
cd server
go run main.go
```

The server will start on port 3000 by default (or the port specified in your environment variables).

### Shutting down the server

The server is designed to handle graceful shutdown to ensure that all database connections are properly closed and in-flight requests are completed. To properly shut down the server:

1. **Using keyboard interrupt**: Press `Ctrl+C` in the terminal where the server is running. The server will detect this signal and perform a graceful shutdown.

2. **Using system signals**: You can also send a SIGTERM signal to the server process:

```bash
# Find the process ID (PID) of the server
pgrep -f "go run main.go"

# Send a SIGTERM signal to the process
kill -TERM <PID>
```

During shutdown, the server will:
1. Stop accepting new connections
2. Complete any in-flight requests (with a 15-second timeout)
3. Close all database connections properly
4. Log the shutdown process

This ensures that no data is lost and all resources are properly released.

### Database Maintenance

The server includes an automated database maintenance system that helps keep the database clean and efficient. This system primarily handles the cleanup of temporary data that is no longer needed.

#### Scheduled Cleanup

By default, the server runs a scheduled cleanup task once every 24 hours. This task:

- Removes text entries from the `text_storage` table that are older than 7 days and have `save_flag` set to `false`
- Logs the number of entries removed during each cleanup operation
- Handles any errors that occur during the cleanup process

This automated cleanup ensures that the database doesn't grow unnecessarily large with temporary data that is no longer needed.

#### Manual Cleanup

Administrators can also trigger a manual cleanup operation through the private maintenance endpoint:

```bash
# Using curl (requires authentication cookie)
curl -X POST -b "auth_cookie=your_auth_cookie" http://localhost:3000/private/maintenance/cleanup
```

The endpoint returns a JSON response with information about the cleanup operation:

```json
{
  "success": true,
  "message": "Successfully removed 42 expired text entries",
  "entries_removed": 42,
  "timestamp": "2025-05-24T18:27:01Z"
}
```

#### Customizing Cleanup Behavior

To modify the cleanup schedule or retention period, you'll need to edit the following files:

- `server/main.go`: Change the `cleanupInterval` in the `scheduleCleanup` function to adjust how often the cleanup runs
- `server/internal/handlers/cleanup.go`: Modify the retention period (currently 7 days) in the `DatabaseCleanupHandler` function

## Website Sections

1. **Homepage** (`/`) - Introduction and list of available tools
2. **Documentation Section** (`/docs/...`) - Details about each tool
3. **Tools Section** (`/tools/...`) - Access to the actual tools

### Available Tools

The server currently includes the following tools:

1. **Random Number Generator** (`/tools/random-number`) - Generates a random number within a specified range
   - Parameters: `min` (default: 1), `max` (default: 100)

2. **Date Formatter** (`/tools/date`) - Formats the current date according to specified parameters
   - Parameters: `format` (default: "2006-01-02"), `offset` (default: 0)

3. **Day Tool** (`/tools/day`) - Returns the current day of the month
   - No parameters required

4. **Month Tool** (`/tools/month`) - Returns the current month as a string
   - No parameters required

5. **Year Tool** (`/tools/year`) - Returns the current year
   - No parameters required

6. **Text File Tool** (`/tools/text-file`) - Generates a downloadable text file from provided content
   - Parameters: `content` (required), `filename` (default: "download.txt")
   - Returns a downloadable text file with the specified content

7. **Text Formatter** (`/tools/text-formatter`) - Formats text with various options
   - Parameters: `text` (required), `uppercase` (default: false), `lowercase` (default: false)

### Private Tools

The server also includes private tools that require authentication. These tools are accessible at `/private/tools/{tool_name}` after logging in through the `/login` page.

1. **Text Storage** (`/private/tools/text-storage`) - Stores text content in the database
   - Parameters: `content` (required), `save` (default: false)
   - Returns a unique string ID for the stored text
   - The `save` parameter determines whether the text should be permanently saved

2. **Text Retrieval** (`/private/tools/text-retrieval`) - Retrieves text content from the database
   - Parameters: `id` (required)
   - Returns the text content associated with the provided ID

Private tools require authentication using a password. The password hash is stored in the `.env` file as `PRIVATE_USE_PASSWORD`. You can generate a password hash using the utility in `cmd/hashpassword/main.go`.

### Output Formats

Each tool supports multiple output formats:

- **HTML** - For browser viewing
- **JSON** - For API integration
- **Raw** - Plain text output

The output format is determined by the tool's configuration and can be influenced by the client's Accept header.

## Template Rendering

The server uses Go's html/template package for rendering HTML pages. The template system includes:

1. **Template Manager** - Handles loading and rendering of templates
2. **Content Negotiation** - Automatically detects the client's preferred content type and responds accordingly
3. **Template Structure**:
   - Base layout template with common elements (header, footer, styles)
   - Page-specific templates for homepage, documentation, and tools
   - Error templates (e.g., 404 Not Found)

All handlers support both HTML and JSON responses based on the client's Accept header, making the server suitable for both browser-based usage and API integration.

## Testing

### Running Tests

To run all tests:

```bash
cd server
go test ./tests/unit/...
```

### Test Structure

The project includes comprehensive tests for all components:

1. **Model Tests** - Tests for data structures and validation
2. **Handler Tests** - Tests for HTTP handlers and routing
3. **Template Tests** - Tests for template rendering and content negotiation
4. **Tool Tests** - Tests for tool functionality:
   - Random Number Generator tests with parameter validation
   - Date Formatter tests with mocked time for consistent results
   - Day/Month/Year tools tests with mocked time

All tools are tested for:
- Parameter validation
- Error handling
- Output correctness

Testing .
