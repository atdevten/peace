#!/bin/bash

# Script to import data only once on startup
# Checks if quotes already exist before importing

CONFIG_FILE="./configs/docker.env"
DATA_FILE="./quotes.csv.zip"
WORKERS=5
BATCH_SIZE=50

echo "🚀 Starting data import check..."

# Check if data file exists
if [ ! -f "$DATA_FILE" ]; then
    echo "⚠️  Data file $DATA_FILE not found. Skipping import."
    exit 0
fi

# Build data-importer
echo "📦 Building data-importer..."
go build -o /tmp/data-importer ./cmd/data-importer

# Check if quotes already exist in database
echo "🔍 Checking if quotes already exist..."
QUOTE_COUNT=$(echo "SELECT COUNT(*) FROM quotes;" | PGPASSWORD=local_password psql -h postgres -U postgres -d peace_local -t -A 2>/dev/null | tr -d ' ')

if [ "$QUOTE_COUNT" -gt "0" ]; then
    echo "✅ Found $QUOTE_COUNT quotes in database. Skipping import."
    exit 0
fi

echo "📝 No quotes found. Starting import..."
echo "📊 Using $WORKERS workers with batch size $BATCH_SIZE"

# Run the importer
/tmp/data-importer -config "$CONFIG_FILE" -file "$DATA_FILE" -workers "$WORKERS" -batch-size "$BATCH_SIZE"

if [ $? -eq 0 ]; then
    echo "✅ Data import completed successfully!"
else
    echo "❌ Data import failed!"
    exit 1
fi
