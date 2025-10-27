# Reolink Server - Technical Specification

## 1. Camera Manager Implementation

### 1.1 Camera Registry
```go
type Camera struct {
    ID           string
    Name         string
    Host         string
    Username     string
    Password     string // Encrypted
    Model        string
    FirmwareVer  string
    Capabilities map[string]bool
    Status       CameraStatus
    LastSeen     time.Time
    Tags         []string
    GroupID      *string
    Client       *reolink.Client // SDK client instance
}

type CameraStatus string
const (
    StatusOnline     CameraStatus = "online"
    StatusOffline    CameraStatus = "offline"
    StatusError      CameraStatus = "error"
    StatusConnecting CameraStatus = "connecting"
)
```

### 1.2 Manager Interface
```go
type CameraManager interface {
    // Lifecycle
    AddCamera(ctx context.Context, config CameraConfig) (*Camera, error)
    RemoveCamera(ctx context.Context, id string) error
    UpdateCamera(ctx context.Context, id string, config CameraConfig) error
    
    // Operations
    GetCamera(ctx context.Context, id string) (*Camera, error)
    ListCameras(ctx context.Context, filters CameraFilters) ([]*Camera, error)
    GetCameraStatus(ctx context.Context, id string) (*CameraStatus, error)
    
    // Control
    ExecuteCommand(ctx context.Context, id string, cmd Command) error
    
    // Health
    StartHealthMonitoring(ctx context.Context)
    StopHealthMonitoring(ctx context.Context)
}
```

### 1.3 Health Monitoring Strategy
- **Interval**: Every 30 seconds
- **Method**: Call `GetDeviceInfo()` from SDK
- **Failure Handling**:
  - 1st failure: Log warning, retry in 10s
  - 2nd failure: Mark as error, retry in 30s
  - 3rd failure: Mark as offline, retry in 60s
  - Success after failure: Mark as online, reset counter
- **Circuit Breaker**: After 5 consecutive failures, stop polling for 5 minutes

### 1.4 Connection Pooling
- Maintain one SDK client per camera
- Reuse client for all API calls
- Auto-reconnect on token expiration
- Lazy initialization on first use

## 2. Event Processing System

### 2.1 Event Types
```go
type EventType string
const (
    EventMotionDetected    EventType = "motion_detected"
    EventAIPerson         EventType = "ai_person"
    EventAIVehicle        EventType = "ai_vehicle"
    EventAIPet            EventType = "ai_pet"
    EventAudioAlarm       EventType = "audio_alarm"
    EventRecordingStart   EventType = "recording_start"
    EventRecordingStop    EventType = "recording_stop"
    EventCameraOnline     EventType = "camera_online"
    EventCameraOffline    EventType = "camera_offline"
)

type Event struct {
    ID          string
    CameraID    string
    CameraName  string
    Type        EventType
    Severity    Severity
    Timestamp   time.Time
    Metadata    map[string]interface{}
    SnapshotURL *string
    VideoURL    *string
    Acknowledged bool
    AcknowledgedAt *time.Time
}

type Severity string
const (
    SeverityInfo     Severity = "info"
    SeverityWarning  Severity = "warning"
    SeverityCritical Severity = "critical"
)
```

### 2.2 Event Polling Strategy
**Motion Detection:**
- Poll `GetMdState()` every 5 seconds per camera
- Compare with previous state
- Generate event on state change

**AI Detection:**
- Poll `GetAiState()` every 5 seconds per camera
- Check for active detections
- Generate event for each detection type

**Optimization:**
- Use goroutine pool (worker pattern)
- Batch polling for multiple cameras
- Cache last state to detect changes
- Configurable poll intervals per event type

### 2.3 Event Flow
```
Camera → Poller → Event Generator → Message Queue → Event Processor → Storage
                                                   ↓
                                              WebSocket Broadcaster
```

### 2.4 Event Storage (TimescaleDB)
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camera_id UUID NOT NULL REFERENCES cameras(id),
    event_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    metadata JSONB,
    snapshot_url TEXT,
    video_url TEXT,
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to hypertable for time-series optimization
SELECT create_hypertable('events', 'timestamp');

-- Create indexes
CREATE INDEX idx_events_camera_id ON events(camera_id);
CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_timestamp ON events(timestamp DESC);
CREATE INDEX idx_events_acknowledged ON events(acknowledged) WHERE NOT acknowledged;
```

## 3. Stream Management

### 3.1 Stream Types
- **RTSP**: Direct camera RTSP URL (primary)
- **RTMP**: Direct camera RTMP URL
- **FLV**: Direct camera FLV URL
- **HLS**: Transcoded HLS (optional, for web browsers)

### 3.2 Stream Proxy Architecture
```
Client → Server Proxy → Camera RTSP Stream
         ↓
         Authentication Check
         ↓
         Session Management
         ↓
         Optional Transcoding (FFmpeg)
```

### 3.3 HLS Transcoding (Optional)
```go
type Transcoder struct {
    ffmpegPath string
    outputDir  string
}

func (t *Transcoder) StartHLSStream(rtspURL string, sessionID string) error {
    // ffmpeg -i rtsp://camera/stream -c:v copy -c:a aac -f hls \
    //        -hls_time 2 -hls_list_size 5 -hls_flags delete_segments \
    //        output.m3u8
}
```

### 3.4 Stream Session Management
```go
type StreamSession struct {
    ID          string
    CameraID    string
    UserID      string
    StreamType  StreamType
    StartedAt   time.Time
    LastAccess  time.Time
    ExpiresAt   time.Time
}

// Auto-cleanup inactive sessions after 5 minutes
```

## 4. API Design

### 4.1 Request/Response Format
**Standard Response:**
```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "timestamp": "2025-10-27T10:30:00Z"
}
```

**Error Response:**
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "CAMERA_NOT_FOUND",
    "message": "Camera with ID 'cam-123' not found",
    "details": {}
  },
  "timestamp": "2025-10-27T10:30:00Z"
}
```

### 4.2 Pagination
```
GET /api/v1/events?page=1&limit=50&sort=-timestamp
```

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 1234,
      "total_pages": 25
    }
  }
}
```

### 4.3 Filtering
```
GET /api/v1/events?camera_id=cam-123&type=motion_detected&start_date=2025-10-20&end_date=2025-10-27
```

### 4.4 Authentication
**JWT-based authentication:**
```
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "password"
}

Response:
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-10-28T10:30:00Z"
  }
}

Subsequent requests:
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

## 5. WebSocket Event Streaming

### 5.1 Connection Protocol
```javascript
// Client connects
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/events?token=JWT_TOKEN');

// Subscribe to specific cameras
ws.send(JSON.stringify({
  action: 'subscribe',
  cameras: ['cam-123', 'cam-456']
}));

// Subscribe to specific event types
ws.send(JSON.stringify({
  action: 'subscribe',
  event_types: ['motion_detected', 'ai_person']
}));

// Receive events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event received:', data);
};
```

### 5.2 Event Message Format
```json
{
  "type": "event",
  "event": {
    "id": "evt-789",
    "camera_id": "cam-123",
    "camera_name": "Front Door",
    "event_type": "motion_detected",
    "severity": "warning",
    "timestamp": "2025-10-27T10:30:00Z",
    "metadata": {
      "channel": 0,
      "sensitivity": 50
    },
    "snapshot_url": "/api/v1/events/evt-789/snapshot"
  }
}
```

### 5.3 Heartbeat
- Server sends ping every 30 seconds
- Client responds with pong
- Disconnect if no pong received within 60 seconds

## 6. Database Schema

### 6.1 Cameras Table
```sql
CREATE TABLE cameras (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_encrypted TEXT NOT NULL,
    model VARCHAR(100),
    firmware_version VARCHAR(50),
    capabilities JSONB DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'offline',
    last_seen TIMESTAMPTZ,
    tags TEXT[] DEFAULT '{}',
    group_id UUID REFERENCES camera_groups(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(host)
);
```

### 6.2 Camera Groups Table
```sql
CREATE TABLE camera_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 6.3 Camera Configs Table
```sql
CREATE TABLE camera_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL, -- 'image', 'alarm', 'network', etc.
    config_data JSONB NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(camera_id, config_type)
);
```

### 6.4 Recordings Table
```sql
CREATE TABLE recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camera_id UUID NOT NULL REFERENCES cameras(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    duration_seconds INT,
    stream_type VARCHAR(20), -- 'main' or 'sub'
    recording_type VARCHAR(50), -- 'MD', 'TIMING', 'AI_PEOPLE', etc.
    storage_path TEXT,
    thumbnail_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recordings_camera_id ON recordings(camera_id);
CREATE INDEX idx_recordings_start_time ON recordings(start_time DESC);
```

### 6.5 Users Table (for API authentication)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user', -- 'admin', 'user', 'viewer'
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## 7. Redis Cache Strategy

### 7.1 Cache Keys
```
camera:{id}:status          → Camera status (TTL: 30s)
camera:{id}:info            → Camera info (TTL: 5m)
camera:{id}:capabilities    → Camera capabilities (TTL: 1h)
stream:{session_id}         → Stream session (TTL: 5m)
user:{id}:token             → JWT token (TTL: 24h)
```

### 7.2 Cache Invalidation
- On camera update: Invalidate `camera:{id}:*`
- On camera delete: Delete all `camera:{id}:*`
- On status change: Update `camera:{id}:status`

## 8. Performance Considerations

### 8.1 Concurrent Camera Operations
- Use worker pool pattern (e.g., 10 workers)
- Queue camera operations
- Prevent overwhelming cameras with requests

### 8.2 Event Processing
- Use buffered channels for event queue
- Batch insert events to database (every 1s or 100 events)
- Use NATS for distributed event processing (future scaling)

### 8.3 Database Optimization
- Use prepared statements
- Connection pooling (max 25 connections)
- TimescaleDB compression for old events
- Automatic data retention policy (90 days)

### 8.4 API Rate Limiting
- Per-user: 100 requests/minute
- Per-IP: 1000 requests/minute
- Use Redis for rate limit counters

## 9. Error Handling

### 9.1 Camera Errors
```go
type CameraError struct {
    CameraID string
    Op       string // Operation that failed
    Err      error
}

func (e *CameraError) Error() string {
    return fmt.Sprintf("camera %s: %s: %v", e.CameraID, e.Op, e.Err)
}
```

### 9.2 Retry Strategy
- Network errors: Exponential backoff (1s, 2s, 4s, 8s, 16s)
- Authentication errors: Re-login once, then fail
- Timeout errors: Retry with increased timeout
- Max retries: 3

## 10. Monitoring & Metrics

### 10.1 Prometheus Metrics
```
# Camera metrics
reolink_cameras_total{status="online|offline|error"}
reolink_camera_requests_total{camera_id, operation, status}
reolink_camera_request_duration_seconds{camera_id, operation}

# Event metrics
reolink_events_total{camera_id, event_type}
reolink_events_processing_duration_seconds

# API metrics
reolink_api_requests_total{method, endpoint, status}
reolink_api_request_duration_seconds{method, endpoint}

# Stream metrics
reolink_active_streams{camera_id, stream_type}
```

### 10.2 Health Check Endpoint
```
GET /health

Response:
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "24h30m15s",
  "components": {
    "database": "healthy",
    "redis": "healthy",
    "cameras": {
      "total": 10,
      "online": 9,
      "offline": 1
    }
  }
}
```

## 11. Deployment

### 11.1 Docker Compose
```yaml
version: '3.8'
services:
  server:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
  
  postgres:
    image: timescale/timescaledb:latest-pg15
    environment:
      - POSTGRES_DB=reolink_server
      - POSTGRES_USER=reolink
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
```

This technical specification provides the detailed implementation guidelines for building the Reolink server. Would you like me to start implementing any specific component?

