#!/bin/sh

# Database connection settings
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-profile}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}

# Create the database if it doesn't exist
echo "Creating database $DB_NAME if it doesn't exist..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME"

# Apply migrations
echo "Applying migrations..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f ../migrations/001_create_profiles_table.sql

echo "Database initialization complete!" 