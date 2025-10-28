package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/mosleyit/reolink_server/internal/api/handlers"
	apimiddleware "github.com/mosleyit/reolink_server/internal/api/middleware"
	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/config"
	"github.com/mosleyit/reolink_server/internal/storage/repository"
)

// Router holds the HTTP router and dependencies
type Router struct {
	config             *config.Config
	mux                *chi.Mux
	authHandler        *handlers.AuthHandler
	cameraHandler      *handlers.CameraHandler
	eventHandler       *handlers.EventHandler
	recordingHandler   *handlers.RecordingHandler
	eventStreamHandler *handlers.EventStreamHandler
	healthHandler      *handlers.HealthHandler
}

// RouterDependencies holds all dependencies needed by the router
type RouterDependencies struct {
	Config         *config.Config
	CameraManager  *camera.Manager
	EventProcessor service.EventProcessor
	DB             *sql.DB
	CameraRepo     *repository.CameraRepository
	EventRepo      *repository.EventRepository
	RecordingRepo  *repository.RecordingRepository
	UserRepo       *repository.UserRepository
}

// NewRouter creates a new HTTP router
func NewRouter(deps *RouterDependencies) *Router {
	// Create services
	authService := service.NewAuthService(deps.UserRepo, deps.Config.Auth.JWTSecret, deps.Config.Auth.JWTExpiration)
	cameraService := service.NewCameraService(deps.CameraManager, deps.CameraRepo, deps.EventRepo)
	eventService := service.NewEventService(deps.EventRepo)
	recordingService := service.NewRecordingService(deps.RecordingRepo, deps.CameraManager)

	// Create event stream service if processor is provided
	var eventStreamService *service.EventStreamService
	if deps.EventProcessor != nil {
		// Create adapter to convert EventProcessor interface
		adapter := service.NewProcessorAdapter(func(sub service.EventSubscriber) {
			deps.EventProcessor.Subscribe(sub)
		})
		eventStreamService = service.NewEventStreamService(adapter)
	}

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	cameraHandler := handlers.NewCameraHandler(cameraService)
	eventHandler := handlers.NewEventHandler(eventService)
	recordingHandler := handlers.NewRecordingHandler(recordingService)
	var eventStreamHandler *handlers.EventStreamHandler
	if eventStreamService != nil {
		eventStreamHandler = handlers.NewEventStreamHandler(eventStreamService)
	}
	healthHandler := handlers.NewHealthHandler(deps.DB)

	r := &Router{
		config:             deps.Config,
		mux:                chi.NewRouter(),
		authHandler:        authHandler,
		cameraHandler:      cameraHandler,
		eventHandler:       eventHandler,
		recordingHandler:   recordingHandler,
		eventStreamHandler: eventStreamHandler,
		healthHandler:      healthHandler,
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// setupMiddleware configures global middleware
func (r *Router) setupMiddleware() {
	// Request ID
	r.mux.Use(middleware.RequestID)

	// Real IP
	r.mux.Use(middleware.RealIP)

	// Logging
	r.mux.Use(apimiddleware.Logger)

	// Recovery from panics
	r.mux.Use(middleware.Recoverer)

	// Timeout
	r.mux.Use(middleware.Timeout(60 * time.Second))

	// CORS
	if r.config.API.EnableCORS {
		r.mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   r.config.API.CORSAllowedOrigins,
			AllowedMethods:   r.config.API.CORSAllowedMethods,
			AllowedHeaders:   r.config.API.CORSAllowedHeaders,
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// Compress responses
	r.mux.Use(middleware.Compress(5))
}

// setupRoutes configures all API routes
func (r *Router) setupRoutes() {
	// Health check (no auth required)
	r.mux.Get("/health", r.healthHandler.HealthCheck)
	r.mux.Get("/ready", r.healthHandler.ReadinessCheck)

	// API v1 routes
	r.mux.Route("/api/v1", func(rt chi.Router) {
		// Public routes
		rt.Group(func(pub chi.Router) {
			pub.Post("/auth/login", r.authHandler.Login)
		})

		// Protected routes (require authentication)
		rt.Group(func(protected chi.Router) {
			// Apply JWT authentication middleware
			protected.Use(apimiddleware.Authenticate(r.config.Auth.JWTSecret))

			// Camera management
			protected.Route("/cameras", func(cam chi.Router) {
				cam.Get("/", r.cameraHandler.ListCameras)
				cam.Post("/", r.cameraHandler.AddCamera)
				cam.Get("/{id}", r.cameraHandler.GetCamera)
				cam.Put("/{id}", r.cameraHandler.UpdateCamera)
				cam.Delete("/{id}", r.cameraHandler.DeleteCamera)
				cam.Get("/{id}/status", r.cameraHandler.GetCameraStatus)
				cam.Post("/{id}/reboot", r.cameraHandler.RebootCamera)
				cam.Get("/{id}/snapshot", r.cameraHandler.GetSnapshot)

				// PTZ control
				cam.Post("/{id}/ptz/move", r.cameraHandler.PTZMove)
				cam.Post("/{id}/ptz/preset", r.cameraHandler.PTZPreset)

				// LED/Siren control
				cam.Post("/{id}/led", r.cameraHandler.ControlLED)
				cam.Post("/{id}/siren", r.cameraHandler.TriggerSiren)

				// Configuration
				cam.Get("/{id}/config/{type}", r.cameraHandler.GetCameraConfig)
				cam.Put("/{id}/config/{type}", r.cameraHandler.UpdateCameraConfig)

				// Events for specific camera
				cam.Get("/{id}/events", r.cameraHandler.GetCameraEvents)

				// Streams
				cam.Get("/{id}/stream/rtsp", r.cameraHandler.GetRTSPURL)
				cam.Get("/{id}/stream/flv", r.cameraHandler.GetFLVURL)
				cam.Get("/{id}/stream/hls", r.cameraHandler.GetHLSURL)
			})

			// Events
			protected.Route("/events", func(evt chi.Router) {
				evt.Get("/", r.eventHandler.ListEvents)
				evt.Get("/{id}", r.eventHandler.GetEvent)
				evt.Put("/{id}/acknowledge", r.eventHandler.AcknowledgeEvent)
				evt.Get("/{id}/snapshot", r.eventHandler.GetEventSnapshot)
			})

			// Recordings
			protected.Route("/recordings", func(rec chi.Router) {
				rec.Get("/", r.recordingHandler.ListRecordings)
				rec.Get("/{id}", r.recordingHandler.GetRecording)
				rec.Get("/{id}/download", r.recordingHandler.DownloadRecording)
				rec.Post("/search", r.recordingHandler.SearchRecordings)
				rec.Delete("/{id}", r.recordingHandler.DeleteRecording)
			})

			// WebSocket for real-time events
			if r.eventStreamHandler != nil {
				protected.Get("/ws/events", r.eventStreamHandler.WebSocketEvents)
				protected.Get("/ws/cameras/{id}/events", r.eventStreamHandler.WebSocketCameraEvents)

				// SSE alternative
				protected.Get("/sse/events", r.eventStreamHandler.SSEEvents)
			} else {
				// Fallback to legacy handlers if event stream not available
				protected.Get("/ws/events", handlers.WebSocketEvents)
				protected.Get("/ws/cameras/{id}/events", handlers.WebSocketCameraEvents)
				protected.Get("/sse/events", handlers.SSEEvents)
			}
		})
	})

	// Serve static files for frontend
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.mux.Get("/", handlers.ServeIndex)
}
