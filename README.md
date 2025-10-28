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

All API endpoints (except `/health` and `/ready`) require JWT authentication.

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2024-01-01T12:00:00Z"
}

# Use token in subsequent requests
curl http://localhost:8080/api/v1/cameras \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Camera Management

```bash
# List all cameras
GET /api/v1/cameras
Response: { "cameras": [...], "total": 5 }

# Add camera
POST /api/v1/cameras
{
  "name": "Front Door",
  "host": "192.168.1.100",
  "port": 80,
  "username": "admin",
  "password": "password",
  "enabled": true
}

# Get camera details
GET /api/v1/cameras/{id}
Response: { "id": "...", "name": "Front Door", "host": "...", ... }

# Update camera
PUT /api/v1/cameras/{id}
{
  "name": "Updated Name",
  "enabled": false
}

# Delete camera
DELETE /api/v1/cameras/{id}

# Get camera status
GET /api/v1/cameras/{id}/status
Response: { "online": true, "recording": true, "last_seen": "..." }

# Reboot camera
POST /api/v1/cameras/{id}/reboot
```

### Camera Configuration

```bash
# Get camera configuration
GET /api/v1/cameras/{id}/config?type=encoding
# Supported types: encoding, network, alarm, led, ftp, email, push,
#                  recording, osd, image, audio, ptz, zoom_focus,
#                  isp, ir_lights, status_led, power_led, auto_focus,
#                  day_night, white_balance, auto_reply, battery

# Update camera configuration
PUT /api/v1/cameras/{id}/config
{
  "type": "led",
  "config": {
    "state": "on"
  }
}
# Supported update types: led, ptz, zoom_focus
```

### Camera Control

```bash
# Take snapshot
GET /api/v1/cameras/{id}/snapshot
Returns: JPEG image

# PTZ control
POST /api/v1/cameras/{id}/ptz/move
{
  "operation": "up",      # up, down, left, right, left_up, left_down, right_up, right_down
  "speed": 32,            # 1-64
  "channel": 0            # optional, default 0
}

# PTZ preset
POST /api/v1/cameras/{id}/ptz/preset
{
  "operation": "goto",    # goto, set, remove
  "preset_id": 1
}

# Control LED
POST /api/v1/cameras/{id}/led
{
  "state": "on",          # on, off, auto
  "channel": 0
}

# Trigger siren
POST /api/v1/cameras/{id}/siren
{
  "duration": 5,          # seconds
  "channel": 0
}

# Start/Stop recording
POST /api/v1/cameras/{id}/recording
{
  "action": "start",      # start, stop
  "channel": 0
}
```

### Events

```bash
# List events with filtering
GET /api/v1/events?page=1&limit=50&camera_id=cam-123&type=motion_detected&start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z

# Response
{
  "events": [...],
  "total": 150,
  "page": 1,
  "limit": 50
}

# Get event details
GET /api/v1/events/{id}

# Acknowledge event
PUT /api/v1/events/{id}/acknowledge

# Get event snapshot (if available)
GET /api/v1/events/{id}/snapshot
Returns: JPEG image
```

### Recordings

```bash
# List recordings
GET /api/v1/recordings?camera_id=cam-123&start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z

# Get recording details
GET /api/v1/recordings/{id}

# Download recording
GET /api/v1/recordings/{id}/download
Response: { "url": "...", "method": "GET", "notes": "..." }
```

### Video Streaming

```bash
# Get stream URLs
GET /api/v1/cameras/{id}/stream/rtsp?stream_type=main&channel=0
Response: { "url": "rtsp://..." }

GET /api/v1/cameras/{id}/stream/rtmp?stream_type=sub&channel=0
Response: { "url": "rtmp://..." }

GET /api/v1/cameras/{id}/stream/flv?stream_type=main&channel=0
Response: { "url": "http://..." }

# Proxy FLV stream (direct streaming through server)
GET /api/v1/cameras/{id}/stream/flv/proxy?stream_type=main&channel=0
Returns: FLV video stream

# Start HLS transcoding session
POST /api/v1/cameras/{id}/stream/hls/start
{
  "stream_type": "main",  # main, sub, ext
  "channel": 0
}
Response: {
  "session_id": "uuid",
  "playlist_url": "/api/v1/stream/hls/{session_id}/playlist.m3u8",
  "expires_at": "..."
}

# Get HLS playlist
GET /api/v1/stream/hls/{session_id}/playlist.m3u8

# Get HLS segment
GET /api/v1/stream/hls/{session_id}/segment_001.ts

# Stop HLS session
DELETE /api/v1/stream/hls/{session_id}
```

### Real-time Event Streaming

#### WebSocket

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/events?token=JWT_TOKEN');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
  // { "id": "...", "camera_id": "...", "type": "motion_detected", ... }
};

// Camera-specific events
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/cameras/{id}/events?token=JWT_TOKEN');
```

#### Server-Sent Events (SSE)

```javascript
const eventSource = new EventSource('/api/v1/sse/events?token=JWT_TOKEN');

eventSource.onmessage = (event) => {
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

### Completed ✅

- [x] Project initialization and structure
- [x] Basic HTTP server with routing (Chi)
- [x] Configuration management (Viper)
- [x] Structured logging (Zap)
- [x] Database integration (PostgreSQL + TimescaleDB)
- [x] Redis integration (caching and streams)
- [x] Camera manager implementation (97% SDK coverage - 146/150 methods)
- [x] Event processing system (Redis Streams)
- [x] WebSocket and SSE event streaming
- [x] Stream proxy (FLV proxy + HLS transcoding)
- [x] Authentication and authorization (JWT)
- [x] Frontend development (minimal web interface)
- [x] Comprehensive testing (250+ unit tests)
- [x] API documentation (this README)

### In Progress 🚧

- [ ] Integration tests
- [ ] Deployment guides (Docker, Kubernetes)
- [ ] OpenAPI/Swagger specification

### Future Enhancements 🔮

- [ ] Multi-user support with RBAC
- [ ] Event-based automation rules
- [ ] Mobile app (React Native)
- [ ] Cloud storage integration (S3, Azure Blob)
- [ ] Advanced analytics and reporting
- [ ] Email/SMS notifications
- [ ] ONVIF protocol support
- [ ] Multi-language support

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

