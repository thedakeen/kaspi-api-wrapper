#!/bin/bash
set -e

echo "Restoring DB from dump..."

DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-kaspi_pay}
DB_HOST=${DB_HOST:-db}

export PGPASSWORD=$DB_PASSWORD

if [ ! -f ./dumps/kaspi_pay.custom ]; then
    echo "Dump file not found at ./dumps/kaspi_pay.custom"
    exit 1
fi

echo "Checking if DB exists..."
if psql -h $DB_HOST -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo "Database $DB_NAME exists. Dropping and recreating..."
    psql -h $DB_HOST -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME"
    createdb -h $DB_HOST -U $DB_USER $DB_NAME
else
    echo "Creating database $DB_NAME..."
    createdb -h $DB_HOST -U $DB_USER $DB_NAME
fi

echo "Restoring data..."
pg_restore --no-owner -h $DB_HOST -d $DB_NAME -U $DB_USER ./dumps/kaspi_pay.custom

echo "Database restored"
