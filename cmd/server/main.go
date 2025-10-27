package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/mosleyit/reolink_server/internal/api"
	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/config"
	"github.com/mosleyit/reolink_server/internal/events"
	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/internal/storage/db"
	"github.com/mosleyit/reolink_server/internal/storage/repository"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = "1.0.0"
	buildTime  = "unknown"
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Reolink Server",
		zap.String("version", version),
		zap.String("build_time", buildTime),
	)

	// Initialize camera manager
	cameraManager := camera.NewManager(nil)
	logger.Info("Camera manager initialized")

	// Initialize event processor
	eventProcessor := events.NewProcessor(cameraManager, nil)
	logger.Info("Event processor initialized")

	// Initialize event store (Redis)
	var eventStore *events.Store
	if cfg.Redis.Host != "" {
		storeConfig := &events.StoreConfig{
			RedisAddr:     cfg.Redis.GetRedisAddr(),
			RedisPassword: cfg.Redis.Password,
			RedisDB:       cfg.Redis.DB,
			StreamName:    "reolink:events",
		}
		var err error
		eventStore, err = events.NewStore(storeConfig)
		if err != nil {
			logger.Warn("Failed to initialize event store, events will not be persisted",
				zap.Error(err))
		} else {
			eventProcessor.Subscribe(eventStore)
			logger.Info("Event store initialized and subscribed")
		}
	}

	// Start event processor
	ctx := context.Background()
	if err := eventProcessor.Start(ctx); err != nil {
		logger.Fatal("Failed to start event processor", zap.Error(err))
	}
	logger.Info("Event processor started")

	// Initialize database connection
	database, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Run database migrations
	if err := database.RunMigrations(ctx, "migrations"); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Initialize repositories
	cameraRepo := repository.NewCameraRepository(database)
	eventRepo := repository.NewEventRepository(database)
	recordingRepo := repository.NewRecordingRepository(database)
	userRepo := repository.NewUserRepository(database)
	logger.Info("Database repositories initialized",
		zap.String("camera_repo", "ready"),
		zap.String("event_repo", "ready"),
		zap.String("recording_repo", "ready"),
		zap.String("user_repo", "ready"))

	// Load cameras from database into camera manager
	cameras, err := cameraRepo.List(ctx)
	if err != nil {
		logger.Warn("Failed to load cameras from database", zap.Error(err))
	} else {
		loadedCount := 0
		for _, camera := range cameras {
			if err := cameraManager.AddCamera(ctx, camera); err != nil {
				logger.Error("Failed to add camera to manager",
					zap.Error(err),
					zap.String("camera_id", camera.ID),
					zap.String("camera_name", camera.Name))
			} else {
				loadedCount++
			}
		}
		logger.Info("Cameras loaded from database",
			zap.Int("total", len(cameras)),
			zap.Int("loaded", loadedCount),
			zap.Int("failed", len(cameras)-loadedCount))
	}

	// Create event processor adapter for the router
	type eventProcessorAdapter struct {
		processor *events.Processor
	}
	// Implement Subscribe method that adapts the interface
	var adapter service.EventProcessor = service.NewProcessorAdapter(func(sub service.EventSubscriber) {
		eventProcessor.Subscribe(sub)
	})

	// Create HTTP router with dependencies
	router := api.NewRouter(&api.RouterDependencies{
		Config:         cfg,
		CameraManager:  cameraManager,
		EventProcessor: adapter,
		DB:             database.DB,
		CameraRepo:     cameraRepo,
		EventRepo:      eventRepo,
		RecordingRepo:  recordingRepo,
		UserRepo:       userRepo,
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting",
			zap.String("address", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Stop event processor
	if err := eventProcessor.Stop(); err != nil {
		logger.Error("Failed to stop event processor", zap.Error(err))
	}

	// Close event store
	if eventStore != nil {
		if err := eventStore.Close(); err != nil {
			logger.Error("Failed to close event store", zap.Error(err))
		}
	}

	// Stop camera manager
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := cameraManager.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown camera manager", zap.Error(err))
	}

	// Close database connection
	if err := database.Close(); err != nil {
		logger.Error("Failed to close database connection", zap.Error(err))
	} else {
		logger.Info("Database connection closed")
	}

	logger.Info("Server stopped")
}
