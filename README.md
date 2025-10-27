# Reolink Server

A Go-based server for managing Reolink cameras, ingesting events, and providing real-time monitoring capabilities.

## Features

- 🎥 **Multi-Camera Management**: Manage multiple Reolink cameras from a single interface
- 📊 **Event Ingestion**: Capture motion detection, AI detection (people, vehicles, pets), and alarm events
- 🔄 **Real-time Updates**: WebSocket and SSE support for live event streaming
- 🎬 **Stream Management**: RTSP, RTMP, FLV, and optional HLS transcoding
- 💾 **Time-Series Storage**: PostgreSQL with TimescaleDB for efficient event storage
- 🔐 **Secure API**: JWT-based authentication and authorization
- 🌐 **Web Interface**: Minimal frontend for testing and monitoring
- 📈 **Health Monitoring**: Automatic camera health checks and reconnection

## Architecture

Built with:
- **Go 1.24+**: High-performance backend
- **Chi Router**: Lightweight HTTP routing
- **PostgreSQL + TimescaleDB**: Time-series event storage
- **Redis**: Caching and session management
- **Zap**: Structured logging
- **Reolink API Wrapper**: Camera SDK integration

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 15+ (with TimescaleDB extension)
- Redis 7+
- Docker & Docker Compose (optional)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/mosleyit/reolink_server.git
cd reolink_server
```

2. Create environment file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start services:
```bash
docker-compose up -d
```

4. Access the application:
- Web UI: http://localhost:8080
- API: http://localhost:8080/api/v1
- Health: http://localhost:8080/health

### Manual Setup

1. Install dependencies:
```bash
make deps
```

2. Create configuration:
```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your settings
```

3. Set up database:
```bash
# Create PostgreSQL database
createdb reolink_server

# Run migrations (TODO: implement)
make migrate-up
```

4. Build and run:
```bash
make build
./bin/reolink-server
```

Or run directly:
```bash
make run
```

## Configuration

Configuration is managed through `configs/config.yaml` and environment variables. Environment variables take precedence.

Key configuration sections:
- `server`: HTTP server settings
- `database`: PostgreSQL connection
- `redis`: Redis connection
- `cameras`: Camera management settings
- `events`: Event processing configuration
- `streams`: Stream management settings
- `auth`: JWT and authentication
- `api`: API and CORS settings

See `configs/config.example.yaml` for all available options.

## API Documentation

### Authentication

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Use token in subsequent requests
curl http://localhost:8080/api/v1/cameras \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Camera Management

```bash
# List cameras
GET /api/v1/cameras

# Add camera
POST /api/v1/cameras
{
  "name": "Front Door",
  "host": "192.168.1.100",
  "username": "admin",
  "password": "password"
}

# Get camera details
GET /api/v1/cameras/{id}

# Update camera
PUT /api/v1/cameras/{id}

# Delete camera
DELETE /api/v1/cameras/{id}

# Get camera status
GET /api/v1/cameras/{id}/status

# Reboot camera
POST /api/v1/cameras/{id}/reboot
```

### Camera Control

```bash
# Take snapshot
GET /api/v1/cameras/{id}/snapshot

# PTZ control
POST /api/v1/cameras/{id}/ptz/move
{
  "operation": "up",
  "speed": 32
}

# Control LED
POST /api/v1/cameras/{id}/led
{
  "state": "on"
}

# Trigger siren
POST /api/v1/cameras/{id}/siren
{
  "duration": 5
}
```

### Events

```bash
# List events
GET /api/v1/events?page=1&limit=50&camera_id=cam-123&type=motion_detected

# Get event details
GET /api/v1/events/{id}

# Acknowledge event
PUT /api/v1/events/{id}/acknowledge
```

### Streams

```bash
# Get RTSP URL
GET /api/v1/cameras/{id}/stream/rtsp

# Get FLV URL
GET /api/v1/cameras/{id}/stream/flv

# Get HLS URL (if transcoding enabled)
GET /api/v1/cameras/{id}/stream/hls
```

### WebSocket Events

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/events?token=JWT_TOKEN');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
};
```

## Development

### Project Structure

```
reolink_server/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── api/            # HTTP handlers and routing
│   ├── camera/         # Camera management
│   ├── events/         # Event processing
│   ├── stream/         # Stream management
│   ├── storage/        # Database and cache
│   ├── config/         # Configuration
│   └── logger/         # Logging
├── pkg/                # Public libraries
├── web/                # Frontend files
├── migrations/         # Database migrations
└── configs/            # Configuration files
```

### Make Commands

```bash
make help           # Show all available commands
make deps           # Install dependencies
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make fmt            # Format code
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-up      # Start Docker Compose
make docker-down    # Stop Docker Compose
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/camera/...
```

## Roadmap

- [x] Project initialization and structure
- [x] Basic HTTP server with routing
- [x] Configuration management
- [x] Logging setup
- [ ] Database integration (PostgreSQL + TimescaleDB)
- [ ] Redis integration
- [ ] Camera manager implementation
- [ ] Event processing system
- [ ] WebSocket event streaming
- [ ] Stream proxy
- [ ] Authentication and authorization
- [ ] Frontend development
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Comprehensive testing
- [ ] Deployment guides

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

## Acknowledgments

- [Reolink API Wrapper](https://github.com/mosleyit/reolink_api_wrapper) - Go SDK for Reolink cameras
- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [Zap](https://github.com/uber-go/zap) - Blazing fast structured logging

## Support

For issues and questions, please open an issue on GitHub.

