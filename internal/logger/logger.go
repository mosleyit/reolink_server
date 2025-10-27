package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mosleyit/reolink_server/internal/config"
)

var (
	// Log is the global logger instance
	Log *zap.Logger
)

// Init initializes the global logger
func Init(cfg config.LoggingConfig) error {
	var zapConfig zap.Config

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure based on format
	if cfg.Format == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.DisableCaller = !cfg.EnableCaller
	zapConfig.DisableStacktrace = !cfg.EnableStacktrace

	// Set output
	if cfg.Output != "" && cfg.Output != "stdout" {
		zapConfig.OutputPaths = []string{cfg.Output}
	}

	// Build logger
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// WithContext returns a logger with context fields
func WithContext(fields ...zap.Field) *zap.Logger {
	if Log == nil {
		// Fallback to a basic logger if not initialized
		Log, _ = zap.NewProduction()
	}
	return Log.With(fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		// Fallback
		logger, _ := zap.NewProduction()
		logger.Fatal(msg, fields...)
	}
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Panic(msg, fields...)
	} else {
		// Fallback
		logger, _ := zap.NewProduction()
		logger.Panic(msg, fields...)
	}
}

// NewLogger creates a new logger instance with the given name
func NewLogger(name string) *zap.Logger {
	if Log == nil {
		Log, _ = zap.NewProduction()
	}
	return Log.Named(name)
}
