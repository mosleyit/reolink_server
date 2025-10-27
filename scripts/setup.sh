#!/bin/bash

set -e

echo "🚀 Reolink Server Setup Script"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p bin logs

# Copy configuration if it doesn't exist
if [ ! -f configs/config.yaml ]; then
    echo "📝 Creating configuration file..."
    cp configs/config.example.yaml configs/config.yaml
    echo "⚠️  Please edit configs/config.yaml with your settings"
fi

# Copy .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "⚠️  Please edit .env with your settings"
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Build the application
echo "🔨 Building application..."
go build -o bin/reolink-server cmd/server/main.go

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit configs/config.yaml with your database and Redis settings"
echo "2. Set up PostgreSQL and Redis (or use docker-compose up -d)"
echo "3. Run the server: ./bin/reolink-server"
echo ""
echo "Or use Docker Compose:"
echo "  docker-compose up -d"
echo ""

