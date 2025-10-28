# Deployment Guide

This guide covers deploying the Reolink Server in various environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Production Considerations](#production-considerations)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **CPU**: 2+ cores recommended (4+ for HLS transcoding)
- **RAM**: 2GB minimum, 4GB+ recommended
- **Storage**: 20GB+ for application and logs, additional space for recordings
- **Network**: Stable connection to camera network

### Software Requirements

- **Go**: 1.24+ (for building from source)
- **PostgreSQL**: 15+ with TimescaleDB 2.22.1 extension
- **Redis**: 7+
- **FFmpeg**: Latest version (required for HLS transcoding)
- **Docker**: 20.10+ and Docker Compose 2.0+ (for containerized deployment)

## Environment Variables

Create a `.env` file in the project root:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=30s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=reolink_server
DB_USER=reolink
DB_PASSWORD=your_secure_password
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Authentication
JWT_SECRET=your_very_secure_random_secret_key_here
JWT_EXPIRATION=24h
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

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

# CORS (for production, specify allowed origins)
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token
CORS_ALLOW_CREDENTIALS=true
```

## Docker Deployment

### Using Docker Compose (Recommended)

1. **Create docker-compose.yml**:

```yaml
version: '3.8'

services:
  postgres:
    image: timescale/timescaledb:latest-pg16
    environment:
      POSTGRES_DB: reolink_server
      POSTGRES_USER: reolink
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U reolink"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  reolink_server:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    env_file:
      - .env
    ports:
      - "8080:8080"
    volumes:
      - ./configs:/app/configs
      - hls_output:/tmp/hls
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  hls_output:
```

2. **Create Dockerfile**:

```dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reolink_server ./cmd/server

FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates ffmpeg tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/reolink_server .
COPY --from=builder /build/web ./web

# Create non-root user
RUN addgroup -g 1000 reolink && \
    adduser -D -u 1000 -G reolink reolink && \
    chown -R reolink:reolink /app

USER reolink

EXPOSE 8080

CMD ["./reolink_server"]
```

3. **Deploy**:

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f reolink_server

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Using Docker Only

```bash
# Build image
docker build -t reolink-server:latest .

# Run container
docker run -d \
  --name reolink-server \
  -p 8080:8080 \
  -e DB_HOST=your_db_host \
  -e REDIS_HOST=your_redis_host \
  --env-file .env \
  reolink-server:latest
```

## Manual Deployment

### 1. Install Dependencies

#### Ubuntu/Debian

```bash
# Install PostgreSQL with TimescaleDB
sudo apt-get update
sudo apt-get install -y postgresql-15 postgresql-contrib

# Add TimescaleDB repository
sudo sh -c "echo 'deb https://packagecloud.io/timescale/timescaledb/ubuntu/ $(lsb_release -c -s) main' > /etc/apt/sources.list.d/timescaledb.list"
wget --quiet -O - https://packagecloud.io/timescale/timescaledb/gpgkey | sudo apt-key add -
sudo apt-get update
sudo apt-get install -y timescaledb-2-postgresql-15

# Install Redis
sudo apt-get install -y redis-server

# Install FFmpeg
sudo apt-get install -y ffmpeg

# Install Go
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 2. Setup Database

```bash
# Create database and user
sudo -u postgres psql << EOF
CREATE DATABASE reolink_server;
CREATE USER reolink WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE reolink_server TO reolink;
\c reolink_server
CREATE EXTENSION IF NOT EXISTS timescaledb;
EOF
```

### 3. Build Application

```bash
# Clone repository
git clone https://github.com/mosleyit/reolink_server.git
cd reolink_server

# Install dependencies
go mod download

# Build
go build -o bin/reolink_server ./cmd/server

# Make executable
chmod +x bin/reolink_server
```

### 4. Create Systemd Service

Create `/etc/systemd/system/reolink-server.service`:

```ini
[Unit]
Description=Reolink Server
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=reolink
Group=reolink
WorkingDirectory=/opt/reolink_server
EnvironmentFile=/opt/reolink_server/.env
ExecStart=/opt/reolink_server/bin/reolink_server
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/tmp/hls

[Install]
WantedBy=multi-user.target
```

### 5. Deploy

```bash
# Create user
sudo useradd -r -s /bin/false reolink

# Copy application
sudo mkdir -p /opt/reolink_server
sudo cp -r . /opt/reolink_server/
sudo chown -R reolink:reolink /opt/reolink_server

# Create HLS output directory
sudo mkdir -p /tmp/hls
sudo chown reolink:reolink /tmp/hls

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable reolink-server
sudo systemctl start reolink-server

# Check status
sudo systemctl status reolink-server

# View logs
sudo journalctl -u reolink-server -f
```

## Production Considerations

### Security

1. **Use HTTPS**: Deploy behind a reverse proxy (Nginx, Caddy) with SSL/TLS
2. **Strong Secrets**: Generate strong JWT secrets and passwords
3. **Firewall**: Restrict access to database and Redis ports
4. **Regular Updates**: Keep dependencies and OS updated
5. **Least Privilege**: Run service as non-root user

### Performance

1. **Database Tuning**: Optimize PostgreSQL for your workload
2. **Connection Pooling**: Adjust max connections based on load
3. **Redis Memory**: Configure maxmemory and eviction policies
4. **HLS Cleanup**: Implement periodic cleanup of old HLS segments
5. **Load Balancing**: Use multiple instances behind a load balancer for high availability

### Backup

```bash
# Database backup
pg_dump -U reolink reolink_server > backup_$(date +%Y%m%d).sql

# Automated daily backups
cat > /etc/cron.daily/reolink-backup << 'EOF'
#!/bin/bash
pg_dump -U reolink reolink_server | gzip > /backups/reolink_$(date +%Y%m%d).sql.gz
find /backups -name "reolink_*.sql.gz" -mtime +30 -delete
EOF
chmod +x /etc/cron.daily/reolink-backup
```

## Monitoring

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready
```

### Metrics

Consider integrating with:
- **Prometheus**: For metrics collection
- **Grafana**: For visualization
- **Loki**: For log aggregation

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL is running
   sudo systemctl status postgresql
   
   # Test connection
   psql -h localhost -U reolink -d reolink_server
   ```

2. **Redis Connection Failed**
   ```bash
   # Check Redis is running
   sudo systemctl status redis
   
   # Test connection
   redis-cli ping
   ```

3. **FFmpeg Not Found**
   ```bash
   # Install FFmpeg
   sudo apt-get install ffmpeg
   
   # Verify installation
   which ffmpeg
   ```

4. **Permission Denied**
   ```bash
   # Check file ownership
   ls -la /opt/reolink_server
   
   # Fix permissions
   sudo chown -R reolink:reolink /opt/reolink_server
   ```

### Logs

```bash
# Application logs (systemd)
sudo journalctl -u reolink-server -f

# Application logs (Docker)
docker-compose logs -f reolink_server

# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log

# Redis logs
sudo tail -f /var/log/redis/redis-server.log
```

### Performance Issues

1. **High CPU Usage**: Check HLS transcoding sessions, reduce concurrent streams
2. **High Memory Usage**: Review database connection pool settings
3. **Slow Queries**: Enable PostgreSQL slow query log, add indexes
4. **Network Issues**: Check camera connectivity, firewall rules

## Support

For additional help:
- GitHub Issues: https://github.com/mosleyit/reolink_server/issues
- Documentation: https://github.com/mosleyit/reolink_server/tree/main/docs

