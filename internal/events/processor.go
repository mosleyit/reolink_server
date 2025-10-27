package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"go.uber.org/zap"
)

// Processor handles event processing from cameras
type Processor struct {
	cameraManager *camera.Manager
	config        *Config
	subscribers   []Subscriber
	mu            sync.RWMutex
	stopCh        chan struct{}
	wg            sync.WaitGroup
	eventCh       chan *models.Event
}

// Config holds processor configuration
type Config struct {
	PollInterval      time.Duration
	MotionCheckPeriod time.Duration
	AICheckPeriod     time.Duration
	EventBufferSize   int
	MaxWorkers        int
}

// DefaultConfig returns default processor configuration
func DefaultConfig() *Config {
	return &Config{
		PollInterval:      5 * time.Second,
		MotionCheckPeriod: 5 * time.Second,
		AICheckPeriod:     10 * time.Second,
		EventBufferSize:   1000,
		MaxWorkers:        10,
	}
}

// Subscriber interface for event subscribers
type Subscriber interface {
	OnEvent(event *models.Event) error
}

// NewProcessor creates a new event processor
func NewProcessor(cameraManager *camera.Manager, config *Config) *Processor {
	if config == nil {
		config = DefaultConfig()
	}

	return &Processor{
		cameraManager: cameraManager,
		config:        config,
		subscribers:   make([]Subscriber, 0),
		stopCh:        make(chan struct{}),
		eventCh:       make(chan *models.Event, config.EventBufferSize),
	}
}

// Subscribe adds a subscriber to receive events
func (p *Processor) Subscribe(subscriber Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers = append(p.subscribers, subscriber)
}

// Start begins event processing
func (p *Processor) Start(ctx context.Context) error {
	logger.Info("Starting event processor")

	// Start event dispatcher
	p.wg.Add(1)
	go p.dispatchEvents(ctx)

	// Start camera pollers
	cameras := p.cameraManager.ListCameras()
	for _, cam := range cameras {
		// Get the camera client
		client, err := p.cameraManager.GetCamera(cam.ID)
		if err != nil {
			logger.Warn("Failed to get camera client",
				zap.String("camera_id", cam.ID),
				zap.Error(err))
			continue
		}

		p.wg.Add(1)
		go p.pollCamera(ctx, client)
	}

	logger.Info("Event processor started", zap.Int("cameras", len(cameras)))
	return nil
}

// Stop stops event processing
func (p *Processor) Stop() error {
	logger.Info("Stopping event processor")
	close(p.stopCh)
	p.wg.Wait()
	close(p.eventCh)
	logger.Info("Event processor stopped")
	return nil
}

// pollCamera polls a single camera for events
func (p *Processor) pollCamera(ctx context.Context, cameraClient *camera.CameraClient) {
	defer p.wg.Done()

	cameraID := cameraClient.Camera.ID
	cameraName := cameraClient.Camera.Name

	logger.Info("Starting camera poller",
		zap.String("camera_id", cameraID),
		zap.String("camera_name", cameraName))

	motionTicker := time.NewTicker(p.config.MotionCheckPeriod)
	defer motionTicker.Stop()

	aiTicker := time.NewTicker(p.config.AICheckPeriod)
	defer aiTicker.Stop()

	for {
		select {
		case <-p.stopCh:
			logger.Info("Camera poller stopped", zap.String("camera_id", cameraID))
			return
		case <-ctx.Done():
			logger.Info("Camera poller context cancelled", zap.String("camera_id", cameraID))
			return
		case <-motionTicker.C:
			p.checkMotionDetection(ctx, cameraClient)
		case <-aiTicker.C:
			p.checkAIDetection(ctx, cameraClient)
		}
	}
}

// checkMotionDetection checks for motion detection events
func (p *Processor) checkMotionDetection(ctx context.Context, cameraClient *camera.CameraClient) {
	// Check motion state for channel 0 (most cameras have single channel)
	state, err := cameraClient.GetMotionState(ctx, 0)
	if err != nil {
		logger.Debug("Failed to get motion state",
			zap.String("camera_id", cameraClient.Camera.ID),
			zap.Error(err))
		return
	}

	// State 1 typically means motion detected
	if state == 1 {
		event := &models.Event{
			ID:         uuid.New().String(),
			CameraID:   cameraClient.Camera.ID,
			CameraName: cameraClient.Camera.Name,
			Type:       models.EventMotionDetected,
			Timestamp:  time.Now(),
			CreatedAt:  time.Now(),
		}

		metadata := models.EventMetadata{
			Channel: 0,
			Extra: map[string]interface{}{
				"state": state,
			},
		}

		if metadataJSON, err := json.Marshal(metadata); err == nil {
			event.Metadata = string(metadataJSON)
		}

		p.publishEvent(event)
	}
}

// checkAIDetection checks for AI detection events
func (p *Processor) checkAIDetection(ctx context.Context, cameraClient *camera.CameraClient) {
	// Check AI state for channel 0
	aiState, err := cameraClient.GetAIState(ctx, 0)
	if err != nil {
		logger.Debug("Failed to get AI state",
			zap.String("camera_id", cameraClient.Camera.ID),
			zap.Error(err))
		return
	}

	if aiState == nil {
		return
	}

	// Check for pet detection (dog/cat)
	if aiState.DogCat.Support == 1 && aiState.DogCat.AlarmState == 1 {
		p.publishAIEvent(cameraClient, models.EventAIPet, aiState)
	}

	// Check for people detection
	if aiState.People.Support == 1 && aiState.People.AlarmState == 1 {
		p.publishAIEvent(cameraClient, models.EventAIPerson, aiState)
	}

	// Check for vehicle detection
	if aiState.Vehicle.Support == 1 && aiState.Vehicle.AlarmState == 1 {
		p.publishAIEvent(cameraClient, models.EventAIVehicle, aiState)
	}
}

// publishAIEvent publishes an AI detection event
func (p *Processor) publishAIEvent(cameraClient *camera.CameraClient, eventType models.EventType, aiState interface{}) {
	event := &models.Event{
		ID:         uuid.New().String(),
		CameraID:   cameraClient.Camera.ID,
		CameraName: cameraClient.Camera.Name,
		Type:       eventType,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	metadata := models.EventMetadata{
		Channel: 0,
		Extra: map[string]interface{}{
			"ai_state": aiState,
		},
	}

	if metadataJSON, err := json.Marshal(metadata); err == nil {
		event.Metadata = string(metadataJSON)
	}

	p.publishEvent(event)
}

// publishEvent sends an event to the event channel
func (p *Processor) publishEvent(event *models.Event) {
	select {
	case p.eventCh <- event:
		logger.Debug("Event published",
			zap.String("event_id", event.ID),
			zap.String("camera_id", event.CameraID),
			zap.String("type", string(event.Type)))
	default:
		logger.Warn("Event channel full, dropping event",
			zap.String("event_id", event.ID),
			zap.String("camera_id", event.CameraID))
	}
}

// dispatchEvents dispatches events to subscribers
func (p *Processor) dispatchEvents(ctx context.Context) {
	defer p.wg.Done()

	logger.Info("Event dispatcher started")

	for {
		select {
		case <-p.stopCh:
			logger.Info("Event dispatcher stopped")
			return
		case <-ctx.Done():
			logger.Info("Event dispatcher context cancelled")
			return
		case event, ok := <-p.eventCh:
			if !ok {
				logger.Info("Event channel closed")
				return
			}

			p.notifySubscribers(event)
		}
	}
}

// notifySubscribers notifies all subscribers of an event
func (p *Processor) notifySubscribers(event *models.Event) {
	p.mu.RLock()
	subscribers := make([]Subscriber, len(p.subscribers))
	copy(subscribers, p.subscribers)
	p.mu.RUnlock()

	for _, subscriber := range subscribers {
		if err := subscriber.OnEvent(event); err != nil {
			logger.Error("Subscriber error",
				zap.String("event_id", event.ID),
				zap.Error(err))
		}
	}
}

// GetEventChannel returns the event channel for direct access
func (p *Processor) GetEventChannel() <-chan *models.Event {
	return p.eventCh
}

// PublishCameraEvent publishes a camera status event (online/offline)
func (p *Processor) PublishCameraEvent(cameraID, cameraName string, eventType models.EventType) {
	event := &models.Event{
		ID:         uuid.New().String(),
		CameraID:   cameraID,
		CameraName: cameraName,
		Type:       eventType,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	p.publishEvent(event)
}

// AddCamera starts polling a new camera
func (p *Processor) AddCamera(ctx context.Context, cameraClient *camera.CameraClient) {
	p.wg.Add(1)
	go p.pollCamera(ctx, cameraClient)

	logger.Info("Added camera to event processor",
		zap.String("camera_id", cameraClient.Camera.ID))
}
