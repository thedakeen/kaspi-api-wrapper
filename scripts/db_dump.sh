#!/bin/bash
set -e

echo "Creating DB dump..."

DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-kaspi_pay}
DB_HOST=${DB_HOST:-kaspi-wrapper-db}

mkdir -p ./dumps

export PGPASSWORD=$DB_PASSWORD

pg_dump -h $DB_HOST -p 5432 --no-owner -Fc -U $DB_USER $DB_NAME -f ./dumps/kaspi_pay.custom

echo "DB dump created"