-- AllMiTools Request Logging Schema
-- Migration: 002_request_logging.sql
-- Description: Creates the request_logs table for tracking HTTP requests
-- Date: 2025-05-25

-- Create request_logs table
CREATE TABLE IF NOT EXISTS request_logs (
    -- Unique identifier for the log entry
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Timestamp when the request was received
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- The endpoint that was requested (e.g., /tools/random-number)
    endpoint TEXT NOT NULL,
    
    -- The HTTP method used (GET, POST, PUT, DELETE, etc.)
    method TEXT NOT NULL,
    
    -- The content type of the request (e.g., application/json)
    content_type TEXT,
    
    -- The request body (may contain form data or JSON)
    request_body TEXT,
    
    -- The query parameters as a JSON string
    query_params TEXT,
    
    -- The HTTP status code of the response
    response_status INTEGER,
    
    -- The time taken to process the request in milliseconds
    response_time_ms INTEGER,
    
    -- The user agent string from the request
    user_agent TEXT,
    
    -- The IP address of the client (anonymized if needed)
    ip_address TEXT
);

-- Create index on timestamp for efficient cleanup queries
CREATE INDEX IF NOT EXISTS idx_request_logs_timestamp ON request_logs(timestamp);

-- Create index on endpoint and method for performance analysis
CREATE INDEX IF NOT EXISTS idx_request_logs_endpoint_method ON request_logs(endpoint, method);

-- Add comments to table and columns for better documentation
COMMENT ON TABLE request_logs IS 'Stores HTTP request logs for debugging and analysis';
COMMENT ON COLUMN request_logs.id IS 'Unique identifier for the log entry';
COMMENT ON COLUMN request_logs.timestamp IS 'Timestamp when the request was received';
COMMENT ON COLUMN request_logs.endpoint IS 'The endpoint that was requested';
COMMENT ON COLUMN request_logs.method IS 'The HTTP method used';
COMMENT ON COLUMN request_logs.content_type IS 'The content type of the request';
COMMENT ON COLUMN request_logs.request_body IS 'The request body (may contain form data or JSON)';
COMMENT ON COLUMN request_logs.query_params IS 'The query parameters as a JSON string';
COMMENT ON COLUMN request_logs.response_status IS 'The HTTP status code of the response';
COMMENT ON COLUMN request_logs.response_time_ms IS 'The time taken to process the request in milliseconds';
COMMENT ON COLUMN request_logs.user_agent IS 'The user agent string from the request';
COMMENT ON COLUMN request_logs.ip_address IS 'The IP address of the client';

-- Create a function to clean up old request logs (retention policy)
CREATE OR REPLACE FUNCTION cleanup_old_request_logs()
RETURNS void AS $$
BEGIN
    -- Delete request logs older than 7 days
    DELETE FROM request_logs 
    WHERE timestamp < NOW() - INTERVAL '7 days';
END;
$$ LANGUAGE plpgsql;

-- Create a comment for the cleanup function
COMMENT ON FUNCTION cleanup_old_request_logs() IS 'Deletes request logs older than 7 days';
