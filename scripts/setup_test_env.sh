#!/bin/bash

# Setup Test Environment for Reolink Server
# This script sets up a local test environment with PostgreSQL and Redis

set -e

echo "🚀 Setting up Reolink Server test environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ PostgreSQL is not installed${NC}"
    echo "Please install PostgreSQL first:"
    echo "  Ubuntu/Debian: sudo apt-get install postgresql"
    echo "  macOS: brew install postgresql"
    exit 1
fi

# Check if Redis is installed
if ! command -v redis-cli &> /dev/null; then
    echo -e "${YELLOW}⚠️  Redis is not installed (optional)${NC}"
    echo "To install Redis:"
    echo "  Ubuntu/Debian: sudo apt-get install redis-server"
    echo "  macOS: brew install redis"
fi

# Database configuration
DB_NAME="${DB_NAME:-reolink_server_dev}"
DB_USER="${DB_USER:-reolink}"
DB_PASSWORD="${DB_PASSWORD:-reolink123}"

echo ""
echo "📊 Database Configuration:"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASSWORD"
echo ""

# Create database user if it doesn't exist
echo "👤 Creating database user..."
sudo -u postgres psql -tc "SELECT 1 FROM pg_user WHERE usename = '$DB_USER'" | grep -q 1 || \
    sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';"

# Create database if it doesn't exist
echo "🗄️  Creating database..."
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    sudo -u postgres psql -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"

# Grant privileges
echo "🔐 Granting privileges..."
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

# Install TimescaleDB extension
echo "⏰ Installing TimescaleDB extension..."
sudo -u postgres psql -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;" || \
    echo -e "${YELLOW}⚠️  TimescaleDB extension not available (optional)${NC}"

# Create .env file
echo "📝 Creating .env file..."
cat > .env << EOF
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=30s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Redis Configuration (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Authentication
JWT_SECRET=$(openssl rand -base64 32)
JWT_EXPIRATION=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=console

# Camera Configuration
CAMERA_HEALTH_CHECK_INTERVAL=60s
CAMERA_RECONNECT_INTERVAL=30s
CAMERA_MAX_RECONNECT_ATTEMPTS=5

# Event Processing
EVENT_PROCESSOR_WORKERS=4
EVENT_BATCH_SIZE=100
EVENT_BATCH_TIMEOUT=5s

# Stream Configuration
STREAM_HLS_OUTPUT_DIR=/tmp/hls
STREAM_FFMPEG_PATH=/usr/bin/ffmpeg
STREAM_SESSION_TIMEOUT=30m
EOF

echo -e "${GREEN}✅ .env file created${NC}"

# Build the application
echo "🔨 Building application..."
go build -o bin/reolink_server ./cmd/server

echo ""
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "📋 Next steps:"
echo "  1. Start the server: ./bin/reolink_server"
echo "  2. The server will run migrations automatically"
echo "  3. Default admin credentials:"
echo "     Username: admin"
echo "     Password: admin"
echo "  4. Access the web UI: http://localhost:8080"
echo "  5. API endpoint: http://localhost:8080/api/v1"
echo ""
echo "🧪 Test the login:"
echo "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"username\":\"admin\",\"password\":\"admin\"}'"
echo ""

