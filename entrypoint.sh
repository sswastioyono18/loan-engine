#!/bin/sh
set -e

echo "Running database migrations..."
./migrate -action up

echo "Starting application..."
exec ./main
