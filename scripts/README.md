# Database Initialization Script

This script initializes your Railway PostgreSQL database with the required tables and sample data.

## Quick Setup

1. **Get your Railway PostgreSQL connection details:**
   - Go to your Railway dashboard
   - Click on your PostgreSQL service
   - Copy the connection details from the Variables tab

2. **Set environment variables:**
   ```bash
   export DB_HOST=your-railway-postgres-host
   export DB_PORT=5432
   export DB_USER=your-postgres-username
   export DB_PASSWORD=your-postgres-password
   export DB_NAME=your-database-name
   export DB_SSLMODE=require
   ```

3. **Run the initialization script:**
   ```bash
   cd scripts
   go mod tidy
   go run init-database.go
   ```

## What this script does:

- ✅ Creates database tables (`users`, `policy_files`, `audit_logs`)
- ✅ Creates a default admin user:
  - Username: `admin`
  - Password: `SecureAdmin123!`
- ✅ Populates the database with sample security policies and onboarding documents

## After running this script:

1. Your database will have all the required tables
2. You can log in to your frontend with the admin credentials
3. Your application should work properly

## Alternative Solutions:

### Option 1: Deploy Backend Service (Recommended)
Deploy your Go backend as a separate Railway service so it can automatically initialize the database.

### Option 2: Update Current Deployment
Modify your `railway.json` to use `Dockerfile.monolith` instead of `Dockerfile.frontend` to deploy both frontend and backend together.

## Troubleshooting:

- Make sure your PostgreSQL service is running on Railway
- Verify your connection details are correct
- Check that your database allows connections from external IPs 