# Event Processing System

This package provides a comprehensive event processing system for Reolink cameras, including motion detection, AI detection (people, vehicles, pets), and camera status events.

## Architecture

The event system consists of three main components:

### 1. Event Processor (`processor.go`)

The processor polls cameras for events and dispatches them to subscribers.

**Features:**
- Configurable polling intervals for motion and AI detection
- Concurrent polling of multiple cameras
- Event buffering and dispatching
- Subscriber pattern for event distribution
- Graceful start/stop with context support

**Configuration:**
```go
config := &events.Config{
    PollInterval:      5 * time.Second,  // General polling interval
    MotionCheckPeriod: 5 * time.Second,  // Motion detection check frequency
    AICheckPeriod:     10 * time.Second, // AI detection check frequency
    EventBufferSize:   1000,             // Event channel buffer size
    MaxWorkers:        10,               // Max concurrent workers
}
```

**Usage:**
```go
// Create processor
processor := events.NewProcessor(cameraManager, config)

// Subscribe to events
processor.Subscribe(eventStore)

// Start processing
ctx := context.Background()
processor.Start(ctx)

// Stop processing
processor.Stop()
```

### 2. Event Store (`store.go`)

The store persists events to Redis Streams for durability and real-time streaming.

**Features:**
- Redis Streams for event persistence
- Real-time event streaming with consumer groups
- Query events by camera ID, type, or time range
- Automatic stream trimming to manage storage
- Implements Subscriber interface for automatic event persistence

**Configuration:**
```go
config := &events.StoreConfig{
    RedisAddr:     "localhost:6379",
    RedisPassword: "",
    RedisDB:       0,
    StreamName:    "reolink:events",
}
```

**Usage:**
```go
// Create store
store, err := events.NewStore(config)
if err != nil {
    log.Fatal(err)
}
defer store.Close()

// Subscribe to processor
processor.Subscribe(store)

// Query events
events, err := store.GetEvents(ctx, "", "", 100)

// Get events for specific camera
events, err := store.GetEventsByCameraID(ctx, "cam-123", 50)

// Stream events in real-time
eventCh, err := store.StreamEvents(ctx, "$")
for event := range eventCh {
    fmt.Printf("Event: %+v\n", event)
}

// Trim old events
store.TrimStream(ctx, 10000)
```

### 3. Event Models (`internal/storage/models/event.go`)

Defines event types and data structures.

**Event Types:**
- `EventMotionDetected` - Motion detected by camera
- `EventAIPerson` - Person detected by AI
- `EventAIVehicle` - Vehicle detected by AI
- `EventAIPet` - Pet (dog/cat) detected by AI
- `EventAudioAlarm` - Audio alarm triggered
- `EventRecordingStart` - Recording started
- `EventRecordingStop` - Recording stopped
- `EventCameraOnline` - Camera came online
- `EventCameraOffline` - Camera went offline

**Event Structure:**
```go
type Event struct {
    ID           string    // Unique event ID
    CameraID     string    // Camera identifier
    CameraName   string    // Camera name
    Type         EventType // Event type
    Timestamp    time.Time // When event occurred
    Acknowledged bool      // Whether event was acknowledged
    Metadata     string    // JSON metadata
    SnapshotPath string    // Path to snapshot image
    CreatedAt    time.Time // When event was created
}
```

## Event Flow

1. **Polling**: Processor polls each camera at configured intervals
2. **Detection**: Checks motion state and AI detection state
3. **Event Creation**: Creates event objects with metadata
4. **Publishing**: Publishes events to internal channel
5. **Dispatching**: Dispatcher reads from channel and notifies subscribers
6. **Persistence**: Store subscriber saves events to Redis Streams
7. **Streaming**: Clients can stream events in real-time from Redis

## Integration Example

```go
// Initialize camera manager
cameraManager := camera.NewManager(nil)

// Add cameras
camera1 := &models.Camera{
    ID:       "cam-001",
    Name:     "Front Door",
    Host:     "192.168.1.100",
    Port:     80,
    Username: "admin",
    Password: "password",
}
cameraManager.AddCamera(camera1)

// Initialize event processor
processor := events.NewProcessor(cameraManager, nil)

// Initialize event store
storeConfig := &events.StoreConfig{
    RedisAddr:  "localhost:6379",
    StreamName: "reolink:events",
}
store, _ := events.NewStore(storeConfig)

// Subscribe store to processor
processor.Subscribe(store)

// Start processing
ctx := context.Background()
processor.Start(ctx)

// Events are now automatically detected, processed, and stored
```

## Testing

The package includes comprehensive unit tests:

```bash
go test -v ./internal/events/
```

**Test Coverage:**
- Processor initialization and configuration
- Event subscription and publishing
- Event dispatching to multiple subscribers
- Start/stop lifecycle
- Event channel buffering and overflow handling
- Benchmark tests for performance

## Performance Considerations

1. **Polling Intervals**: Adjust based on your needs
   - Lower intervals = more responsive but higher CPU/network usage
   - Higher intervals = less load but slower detection

2. **Event Buffer Size**: Size the buffer based on event volume
   - Too small = events may be dropped
   - Too large = more memory usage

3. **Redis Streams**: Trim regularly to manage storage
   - Use `TrimStream()` to limit stream length
   - Consider TTL-based trimming for time-based retention

4. **Concurrent Cameras**: Each camera runs in its own goroutine
   - Scales well with multiple cameras
   - Monitor goroutine count with many cameras

## Future Enhancements

- [ ] Webhook notifications for events
- [ ] Event filtering and rules engine
- [ ] Snapshot capture on motion/AI events
- [ ] Event aggregation and deduplication
- [ ] Consumer groups for distributed processing
- [ ] Metrics and monitoring integration

