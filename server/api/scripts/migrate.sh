#!/bin/bash
# Database migration script

set -e

cd "$(dirname "$0")/.."

DB_PATH=${DB_PATH:-"data/pawstream.db"}
MIGRATIONS_PATH="migrations"

case "$1" in
    up)
        echo "Running migrations up..."
        go run cmd/api/main.go migrate up
        ;;
    down)
        echo "Running migrations down..."
        go run cmd/api/main.go migrate down
        ;;
    version)
        echo "Checking migration version..."
        go run cmd/api/main.go migrate version
        ;;
    *)
        echo "Usage: $0 {up|down|version}"
        exit 1
        ;;
esac
