-- AllMiTools Initial Database Schema
-- Migration: 001_initial_schema.sql
-- Description: Creates the initial tables for the AllMiTools application
-- Date: 2025-05-24

-- Enable UUID extension for generating unique IDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create text_storage table
CREATE TABLE IF NOT EXISTS text_storage (
    -- Unique identifier for the text entry
    id VARCHAR(36) PRIMARY KEY,
    
    -- The text content to store (can contain HTML, CSS, or JavaScript)
    content TEXT NOT NULL,
    
    -- Flag to indicate if this entry should be saved permanently
    save_flag BOOLEAN NOT NULL DEFAULT false,
    
    -- Timestamp when the entry was created
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on id for faster lookups
CREATE INDEX IF NOT EXISTS idx_text_storage_id ON text_storage(id);

-- Create index on save_flag for filtering saved entries
CREATE INDEX IF NOT EXISTS idx_text_storage_save_flag ON text_storage(save_flag);

-- Create index on created_at for time-based queries
CREATE INDEX IF NOT EXISTS idx_text_storage_created_at ON text_storage(created_at);

-- Add comments to table and columns for better documentation
COMMENT ON TABLE text_storage IS 'Stores text content with unique identifiers';
COMMENT ON COLUMN text_storage.id IS 'Unique identifier for the text entry';
COMMENT ON COLUMN text_storage.content IS 'The text content (can contain HTML, CSS, or JavaScript)';
COMMENT ON COLUMN text_storage.save_flag IS 'Flag to indicate if this entry should be saved permanently';
COMMENT ON COLUMN text_storage.created_at IS 'Timestamp when the entry was created';

-- Create a function to clean up old unsaved entries (retention policy)
CREATE OR REPLACE FUNCTION cleanup_unsaved_text_entries()
RETURNS void AS $$
BEGIN
    -- Delete unsaved entries older than 30 days
    DELETE FROM text_storage 
    WHERE save_flag = false 
    AND created_at < NOW() - INTERVAL '30 days';
END;
$$ LANGUAGE plpgsql;

-- Create a comment for the cleanup function
COMMENT ON FUNCTION cleanup_unsaved_text_entries() IS 'Deletes unsaved text entries older than 30 days';
