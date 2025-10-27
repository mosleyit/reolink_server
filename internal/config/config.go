package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Cameras  CamerasConfig  `mapstructure:"cameras"`
	Events   EventsConfig   `mapstructure:"events"`
	Streams  StreamsConfig  `mapstructure:"streams"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Auth     AuthConfig     `mapstructure:"auth"`
	API      APIConfig      `mapstructure:"api"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Name                  string        `mapstructure:"name"`
	User                  string        `mapstructure:"user"`
	Password              string        `mapstructure:"password"`
	SSLMode               string        `mapstructure:"sslmode"`
	MaxConnections        int           `mapstructure:"max_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	MaxRetries int    `mapstructure:"max_retries"`
	PoolSize   int    `mapstructure:"pool_size"`
}

// CamerasConfig holds camera management configuration
type CamerasConfig struct {
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
	ReconnectInterval   time.Duration `mapstructure:"reconnect_interval"`
	MaxRetries          int           `mapstructure:"max_retries"`
	RequestTimeout      time.Duration `mapstructure:"request_timeout"`
	WorkerPoolSize      int           `mapstructure:"worker_pool_size"`
}

// EventsConfig holds event processing configuration
type EventsConfig struct {
	PollInterval  time.Duration `mapstructure:"poll_interval"`
	RetentionDays int           `mapstructure:"retention_days"`
	BatchSize     int           `mapstructure:"batch_size"`
	BatchInterval time.Duration `mapstructure:"batch_interval"`
	BufferSize    int           `mapstructure:"buffer_size"`
}

// StreamsConfig holds stream management configuration
type StreamsConfig struct {
	SessionTimeout       time.Duration `mapstructure:"session_timeout"`
	CleanupInterval      time.Duration `mapstructure:"cleanup_interval"`
	EnableHLSTranscoding bool          `mapstructure:"enable_hls_transcoding"`
	HLSSegmentDuration   time.Duration `mapstructure:"hls_segment_duration"`
	HLSPlaylistSize      int           `mapstructure:"hls_playlist_size"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level            string `mapstructure:"level"`
	Format           string `mapstructure:"format"`
	Output           string `mapstructure:"output"`
	EnableCaller     bool   `mapstructure:"enable_caller"`
	EnableStacktrace bool   `mapstructure:"enable_stacktrace"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiration time.Duration `mapstructure:"jwt_expiration"`
	BcryptCost    int           `mapstructure:"bcrypt_cost"`
}

// APIConfig holds API configuration
type APIConfig struct {
	RateLimitPerMinute  int      `mapstructure:"rate_limit_per_minute"`
	RateLimitPerIP      int      `mapstructure:"rate_limit_per_ip"`
	EnableCORS          bool     `mapstructure:"enable_cors"`
	CORSAllowedOrigins  []string `mapstructure:"cors_allowed_origins"`
	CORSAllowedMethods  []string `mapstructure:"cors_allowed_methods"`
	CORSAllowedHeaders  []string `mapstructure:"cors_allowed_headers"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	if c.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("jwt secret is required")
	}

	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("jwt secret must be at least 32 characters")
	}

	return nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetServerAddr returns the server address
func (c *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

