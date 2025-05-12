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
psql -h $DB_HOST -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME
if [ $? -ne 0 ]; then
    echo "Creating database $DB_NAME..."
    createdb -h $DB_HOST -U $DB_USER $DB_NAME
else
    echo "Database $DB_NAME already exists."
fi

echo "Restoring data..."
pg_restore --no-owner -h $DB_HOST -d $DB_NAME -U $DB_USER ./dumps/kaspi_pay.custom

echo "Database restored"
