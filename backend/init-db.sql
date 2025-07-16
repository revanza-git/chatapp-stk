-- Initialize database for Security Chatbot
-- This script runs automatically when PostgreSQL container starts

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- For trigram matching and better text search

-- Set default timezone
SET timezone = 'UTC';

-- Create indexes for better performance (GORM will create tables automatically)
-- These will be applied after GORM creates the tables

-- Note: GORM will auto-migrate and create the actual tables
-- This script just sets up extensions and any initial configuration 