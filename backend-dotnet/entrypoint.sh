#!/bin/bash
set -e

echo "=========================================="
echo "Backend .NET Entrypoint Script"
echo "=========================================="

# Validate required environment variables
if [ -z "$SA_PASSWORD" ]; then
    echo "ERROR: SA_PASSWORD environment variable is not set"
    exit 1
fi

# Wait for SQL Server to be ready
echo "Waiting for SQL Server to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0

until /opt/mssql-tools18/bin/sqlcmd -S sqlserver -U sa -P "$SA_PASSWORD" -Q "SELECT 1" -C -b > /dev/null 2>&1; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo "ERROR: SQL Server did not become ready in time"
    exit 1
  fi
  echo "SQL Server is unavailable - attempt $RETRY_COUNT/$MAX_RETRIES"
  sleep 2
done

echo "✓ SQL Server is up and running!"

# Wait for database to be created (the app will create it via migrations)
echo "Waiting for database to be created..."
sleep 5

# Check if database exists, if not, wait for it
DB_CHECK=0
for i in {1..10}; do
  DB_EXISTS=$(/opt/mssql-tools18/bin/sqlcmd -S sqlserver -U sa -P "$SA_PASSWORD" -Q "SET NOCOUNT ON; SELECT COUNT(*) FROM sys.databases WHERE name = 'WarrantyDb'" -h -1 -W -C 2>/dev/null | tr -d '[:space:]')
  
  if [ "$DB_EXISTS" = "1" ]; then
    echo "✓ Database 'WarrantyDb' found"
    DB_CHECK=1
    break
  fi
  
  echo "Waiting for database creation... attempt $i/10"
  sleep 3
done

if [ $DB_CHECK -eq 0 ]; then
  echo "WARNING: Database 'WarrantyDb' not found yet. Migrations will create it."
fi

# Give migrations some time to complete
sleep 5

# Check if seed data has already been applied
echo "Checking for existing seed data..."
SEED_CHECK=$(/opt/mssql-tools18/bin/sqlcmd -S sqlserver -U sa -P "$SA_PASSWORD" -d WarrantyDb -Q "SET NOCOUNT ON; IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'customers') SELECT COUNT(*) FROM customers ELSE SELECT 0" -h -1 -W -C 2>/dev/null | tr -d '[:space:]')

if [ "$SEED_CHECK" = "0" ] || [ -z "$SEED_CHECK" ]; then
    echo "No seed data found. Running seed data script..."
    if /opt/mssql-tools18/bin/sqlcmd -S sqlserver -U sa -P "$SA_PASSWORD" -d WarrantyDb -i /app/seed_data_up.sql -C; then
        echo "✓ Seed data applied successfully!"
    else
        echo "WARNING: Seed data script encountered errors. Check logs above."
        # Don't exit - let the app start anyway
    fi
else
    echo "✓ Seed data already exists (found $SEED_CHECK customers). Skipping seed data execution."
fi

echo "=========================================="
echo "Starting .NET application..."
echo "=========================================="
exec dotnet Backend.Dotnet.API.dll
