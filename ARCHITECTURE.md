# Reolink Server Architecture Plan

## Overview
A Go-based server to ingest data from Reolink cameras and provide control capabilities through a REST API and real-time event streaming. The server will manage multiple cameras, process events, store metadata, and provide a minimal frontend for testing.

## Tech Stack

### Backend
- **Language**: Go 1.24+
- **HTTP Framework**: Chi router (lightweight, idiomatic)
- **Database**: PostgreSQL (structured data) + TimescaleDB extension (time-series events)
- **Cache**: Redis (camera state, session management)
- **Message Queue**: NATS or Redis Streams (event processing) (I think we should use redis to not expand the stack)
- **WebSocket**: github.com/coder/websocket (the successor to nhooyr.io/websocket)
- **Configuration**: Viper (YAML/ENV support)
- **Logging**: zap
- **Metrics**: Prometheus client

### Frontend (Minimal)
- **Framework**: Vanilla JS or Alpine.js (lightweight)
- **UI**: Tailwind CSS
- **Video Player**: Video.js or HLS.js for stream playback
- **Real-time**: WebSocket or Server-Sent Events (SSE)

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx (for stream proxying and SSL termination)
- **Stream Processing**: FFmpeg (optional, for transcoding)

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Web Frontend │  │  Mobile App  │  │  API Clients │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          │ HTTP/WS          │ HTTP/WS          │ HTTP
          │                  │                  │
┌─────────▼──────────────────▼──────────────────▼─────────────────┐
│                      API Gateway / Router                        │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                │
│  │ REST API   │  │ WebSocket  │  │ Stream     │                │
│  │ Endpoints  │  │ Handler    │  │ Proxy      │                │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘                │
└────────┼───────────────┼───────────────┼────────────────────────┘
         │               │               │
┌────────▼───────────────▼───────────────▼────────────────────────┐
│                     Application Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Camera       │  │ Event        │  │ Stream       │          │
│  │ Manager      │  │ Processor    │  │ Manager      │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                  │                  │                  │
│  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐          │
│  │ Health       │  │ Alert        │  │ Recording    │          │
│  │ Monitor      │  │ Handler      │  │ Manager      │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────┬──────────────────┬──────────────────┬────────────────┘
           │                  │                  │
┌──────────▼──────────────────▼──────────────────▼────────────────┐
│                      Data Layer                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ PostgreSQL   │  │ Redis        │  │ NATS/Redis   │          │
│  │ (Metadata)   │  │ (Cache)      │  │ (Events)     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────┬──────────────────────────────────────────────────────┘
           │
┌──────────▼──────────────────────────────────────────────────────┐
│                    Reolink Cameras                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Camera 1     │  │ Camera 2     │  │ Camera N     │          │
│  │ (API/RTSP)   │  │ (API/RTSP)   │  │ (API/RTSP)   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
reolink_server/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── camera.go          # Camera control endpoints
│   │   │   ├── events.go          # Event query endpoints
│   │   │   ├── stream.go          # Stream proxy endpoints
│   │   │   └── websocket.go       # WebSocket handler
│   │   ├── middleware/            # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── cors.go
│   │   └── router.go              # Route definitions
│   ├── camera/
│   │   ├── manager.go             # Camera lifecycle management
│   │   ├── registry.go            # Camera registry
│   │   ├── health.go              # Health monitoring
│   │   └── pool.go                # Connection pooling
│   ├── events/
│   │   ├── processor.go           # Event processing
│   │   ├── subscriber.go          # Event subscription
│   │   ├── types.go               # Event type definitions
│   │   └── store.go               # Event storage
│   ├── stream/
│   │   ├── proxy.go               # Stream proxy
│   │   ├── manager.go             # Stream lifecycle
│   │   └── transcoder.go          # Optional transcoding
│   ├── storage/
│   │   ├── postgres/              # PostgreSQL repositories
│   │   │   ├── camera.go
│   │   │   ├── events.go
│   │   │   └── recordings.go
│   │   ├── redis/                 # Redis cache
│   │   │   └── cache.go
│   │   └── models/                # Data models
│   │       ├── camera.go
│   │       ├── event.go
│   │       └── recording.go
│   ├── config/
│   │   └── config.go              # Configuration management
│   └── logger/
│       └── logger.go              # Logging setup
├── pkg/
│   └── utils/                     # Shared utilities
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── index.html
│   └── templates/
├── migrations/                     # Database migrations
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── configs/
│   ├── config.yaml                # Default configuration
│   └── config.example.yaml
├── scripts/
│   ├── setup.sh                   # Setup script
│   └── docker-entrypoint.sh
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Core Components

### 1. Camera Manager
**Responsibilities:**
- Maintain registry of all cameras (IP, credentials, capabilities)
- Establish and maintain connections using your SDK
- Health monitoring (ping cameras periodically)
- Auto-reconnection on failures
- Capability detection (AI features, PTZ, etc.)

**Key Features:**
- Connection pooling for efficient API calls
- Concurrent camera management with goroutines
- Circuit breaker pattern for failed cameras
- Camera grouping/tagging

### 2. Event Processor
**Responsibilities:**
- Poll cameras for motion detection events
- Subscribe to AI detection events (people, vehicles, pets)
- Process alarm triggers
- Normalize events into common format
- Publish to message queue for downstream processing

**Event Types:**
- Motion Detection
- AI Detection (Person, Vehicle, Pet, Face)
- Audio Alarm
- Recording Start/Stop
- Camera Online/Offline
- Configuration Changes

### 3. Stream Manager
**Responsibilities:**
- Generate authenticated stream URLs (RTSP/RTMP/FLV)
- Proxy streams through server (optional)
- Manage concurrent stream sessions
- Optional: Transcode to HLS/DASH for web playback

### 4. Storage Layer
**Database Schema:**

**cameras table:**
- id, name, host, username, password_encrypted
- model, firmware_version, capabilities (JSONB)
- status, last_seen, created_at, updated_at
- tags (array), group_id

**events table (TimescaleDB hypertable):**
- id, camera_id, event_type, severity
- timestamp, metadata (JSONB)
- snapshot_url, video_clip_url
- acknowledged, acknowledged_at

**recordings table:**
- id, camera_id, file_name, file_size
- start_time, end_time, duration
- stream_type, recording_type
- storage_path, thumbnail_url

**camera_configs table:**
- camera_id, config_type, config_data (JSONB)
- version, updated_at

### 5. REST API Endpoints

**Camera Management:**
```
POST   /api/v1/cameras              # Add camera
GET    /api/v1/cameras              # List cameras
GET    /api/v1/cameras/:id          # Get camera details
PUT    /api/v1/cameras/:id          # Update camera
DELETE /api/v1/cameras/:id          # Remove camera
GET    /api/v1/cameras/:id/status   # Get camera status
POST   /api/v1/cameras/:id/reboot   # Reboot camera
```

**Camera Control:**
```
POST   /api/v1/cameras/:id/ptz/move       # PTZ control
POST   /api/v1/cameras/:id/ptz/preset     # Go to preset
GET    /api/v1/cameras/:id/snapshot       # Take snapshot
POST   /api/v1/cameras/:id/led            # Control LED
POST   /api/v1/cameras/:id/siren          # Trigger siren
```

**Configuration:**
```
GET    /api/v1/cameras/:id/config/:type   # Get config (image, alarm, etc.)
PUT    /api/v1/cameras/:id/config/:type   # Update config
```

**Events:**
```
GET    /api/v1/events                     # List events (with filters)
GET    /api/v1/events/:id                 # Get event details
PUT    /api/v1/events/:id/acknowledge     # Acknowledge event
GET    /api/v1/cameras/:id/events         # Camera-specific events
```

**Streams:**
```
GET    /api/v1/cameras/:id/stream/rtsp    # Get RTSP URL
GET    /api/v1/cameras/:id/stream/flv     # Get FLV URL
GET    /api/v1/cameras/:id/stream/hls     # Get HLS playlist (if transcoding)
```

**Recordings:**
```
GET    /api/v1/recordings                 # List recordings
GET    /api/v1/recordings/:id             # Get recording details
GET    /api/v1/recordings/:id/download    # Download recording
POST   /api/v1/recordings/search          # Search recordings
```

### 6. WebSocket/SSE Events
**Real-time event streaming:**
```
WS     /api/v1/ws/events                  # Subscribe to all events
WS     /api/v1/ws/cameras/:id/events      # Subscribe to camera events
SSE    /api/v1/sse/events                 # SSE alternative
```

**Event Message Format:**
```json
{
  "type": "motion_detected",
  "camera_id": "cam-123",
  "camera_name": "Front Door",
  "timestamp": "2025-10-27T10:30:00Z",
  "data": {
    "channel": 0,
    "sensitivity": 50,
    "snapshot_url": "/api/v1/events/evt-456/snapshot"
  }
}
```

## Minimal Frontend Features

### Dashboard
- Camera grid view with live thumbnails
- Camera status indicators (online/offline)
- Recent events feed
- System health metrics

### Camera View
- Live stream player (HLS/FLV)
- PTZ controls (if supported)
- Snapshot button
- LED/Siren controls
- Quick settings (brightness, motion detection on/off)

### Events Page
- Filterable event list (by camera, type, date)
- Event timeline visualization
- Snapshot previews
- Acknowledge/dismiss events

### Settings Page
- Add/edit/remove cameras
- Configure event notifications
- System settings

## Implementation Phases

### Phase 1: Foundation (Week 1)
- [x] Project structure setup
- [ ] Go module initialization
- [ ] Basic HTTP server with Chi router
- [ ] Configuration management (Viper)
- [ ] Logging setup (zerolog)
- [ ] Database setup (PostgreSQL + migrations)
- [ ] Redis connection

### Phase 2: Camera Management (Week 2)
- [ ] Camera registry implementation
- [ ] SDK integration for camera operations
- [ ] Health monitoring system
- [ ] Connection pooling
- [ ] Basic CRUD API endpoints

### Phase 3: Event Processing (Week 3)
- [ ] Event polling mechanism
- [ ] Event normalization
- [ ] Event storage (TimescaleDB)
- [ ] Message queue integration
- [ ] WebSocket event streaming

### Phase 4: Streaming (Week 4)
- [ ] Stream URL generation
- [ ] Stream proxy implementation
- [ ] Optional: HLS transcoding
- [ ] Stream session management

### Phase 5: Frontend (Week 5)
- [ ] Basic HTML/CSS/JS setup
- [ ] Camera list view
- [ ] Live stream player
- [ ] Event feed
- [ ] Control interface

### Phase 6: Polish & Testing (Week 6)
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Unit tests
- [ ] Integration tests
- [ ] Docker containerization
- [ ] Deployment documentation

## Configuration Example

```yaml
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  name: reolink_server
  user: reolink
  password: ${DB_PASSWORD}
  max_connections: 25

redis:
  host: localhost
  port: 6379
  password: ${REDIS_PASSWORD}
  db: 0

cameras:
  health_check_interval: 30s
  reconnect_interval: 60s
  max_retries: 3
  
events:
  poll_interval: 5s
  retention_days: 90
  
logging:
  level: info
  format: json
```

## Security Considerations

1. **Authentication**: JWT-based API authentication
2. **Camera Credentials**: Encrypt passwords in database
3. **HTTPS**: SSL/TLS for all external communication
4. **CORS**: Proper CORS configuration
5. **Rate Limiting**: Prevent API abuse
6. **Input Validation**: Validate all user inputs
7. **Stream Security**: Token-based stream access

## Next Steps

1. Initialize the project structure
2. Set up the development environment
3. Implement core camera manager
4. Build REST API foundation
5. Add event processing
6. Create minimal frontend

Would you like me to start implementing any specific component?

