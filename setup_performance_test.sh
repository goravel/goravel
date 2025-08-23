#!/bin/bash

# Setup script for DB vs ORM performance testing
# This script sets up a PostgreSQL database for testing

set -e

echo "🚀 Setting up PostgreSQL for performance testing..."

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is required but not installed. Please install Docker first."
    exit 1
fi

# Stop and remove existing container if it exists
echo "🔄 Cleaning up existing containers..."
docker stop goravel-postgres 2>/dev/null || true
docker rm goravel-postgres 2>/dev/null || true

# Start PostgreSQL container
echo "🐘 Starting PostgreSQL container..."
docker run --name goravel-postgres \
    -e POSTGRES_DB=goravel \
    -e POSTGRES_USER=root \
    -e POSTGRES_PASSWORD=password \
    -p 5432:5432 \
    -d postgres:15

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Test connection
echo "🔗 Testing database connection..."
for i in {1..30}; do
    if docker exec goravel-postgres pg_isready -U root -d goravel > /dev/null 2>&1; then
        echo "✅ PostgreSQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ PostgreSQL failed to start after 30 seconds"
        exit 1
    fi
    sleep 1
done

# Update .env file
echo "📝 Updating .env file..."
if [ -f .env ]; then
    # Update existing .env
    sed -i.bak 's/^DB_PASSWORD=.*/DB_PASSWORD=password/' .env
    echo "✅ Updated existing .env file"
else
    # Copy from .env.example and update
    cp .env.example .env
    sed -i.bak 's/^DB_PASSWORD=.*/DB_PASSWORD=password/' .env
    echo "✅ Created .env file from .env.example"
fi

# Generate app key if needed
echo "🔑 Checking application key..."
if ! grep -q "APP_KEY=.*" .env || grep -q "APP_KEY=$" .env; then
    echo "🔑 Generating application key..."
    go run . artisan key:generate
fi

# Run migrations
echo "📋 Running database migrations..."
go run . artisan migrate || {
    echo "⚠️  Migration failed. This might be normal if tables don't exist yet."
    echo "    You can manually run migrations later with: go run . artisan migrate"
}

echo ""
echo "🎉 Setup complete! You can now run performance tests:"
echo ""
echo "   # Run basic tests:"
echo "   go test -v ./tests/performance/"
echo ""
echo "   # Run benchmarks:"
echo "   go test -bench=. ./tests/performance/"
echo ""
echo "   # Run with memory profiling:"
echo "   go test -bench=. -benchmem ./tests/performance/"
echo ""
echo "💡 To stop the database when done:"
echo "   docker stop goravel-postgres"
echo ""
echo "🗑️  To remove the database completely:"
echo "   docker stop goravel-postgres && docker rm goravel-postgres"