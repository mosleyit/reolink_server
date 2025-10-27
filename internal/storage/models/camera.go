package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Camera represents a Reolink camera in the system
type Camera struct {
	ID           string             `json:"id" db:"id"`
	Name         string             `json:"name" db:"name"`
	Host         string             `json:"host" db:"host"`
	Port         int                `json:"port" db:"port"`
	Username     string             `json:"username" db:"username"`
	Password     string             `json:"-" db:"password"` // Never expose in JSON
	UseHTTPS     bool               `json:"use_https" db:"use_https"`
	SkipVerify   bool               `json:"skip_verify" db:"skip_verify"`
	Status       string             `json:"status" db:"status"` // online, offline, error
	Model        string             `json:"model" db:"model"`
	FirmwareVer  string             `json:"firmware_version" db:"firmware_version"`
	HardwareVer  string             `json:"hardware_version" db:"hardware_version"`
	Capabilities CameraCapabilities `json:"capabilities" db:"capabilities"`
	Tags         pq.StringArray     `json:"tags" db:"tags"`
	GroupID      *string            `json:"group_id,omitempty" db:"group_id"`
	LastSeen     time.Time          `json:"last_seen" db:"last_seen"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
}

// CameraCapabilities represents camera capabilities stored as JSONB
type CameraCapabilities map[string]bool

// Value implements the driver.Valuer interface for database storage
func (c CameraCapabilities) Value() (driver.Value, error) {
	if c == nil {
		return json.Marshal(map[string]bool{})
	}
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface for database retrieval
func (c *CameraCapabilities) Scan(value interface{}) error {
	if value == nil {
		*c = make(CameraCapabilities)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan CameraCapabilities: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, c)
}

// CameraStatus represents the current status of a camera
type CameraStatus struct {
	CameraID    string    `json:"camera_id"`
	Status      string    `json:"status"` // online, offline, error
	Model       string    `json:"model"`
	FirmwareVer string    `json:"firmware_version"`
	Uptime      int64     `json:"uptime"` // seconds
	LastSeen    time.Time `json:"last_seen"`
	Error       string    `json:"error,omitempty"`
}

// CreateCameraRequest represents a request to add a new camera
type CreateCameraRequest struct {
	Name       string `json:"name" validate:"required"`
	Host       string `json:"host" validate:"required"`
	Port       int    `json:"port"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	UseHTTPS   bool   `json:"use_https"`
	SkipVerify bool   `json:"skip_verify"`
}

// UpdateCameraRequest represents a request to update camera settings
type UpdateCameraRequest struct {
	Name       *string `json:"name,omitempty"`
	Host       *string `json:"host,omitempty"`
	Port       *int    `json:"port,omitempty"`
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	UseHTTPS   *bool   `json:"use_https,omitempty"`
	SkipVerify *bool   `json:"skip_verify,omitempty"`
}
