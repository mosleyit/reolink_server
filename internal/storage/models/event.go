package models

import (
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventMotionDetected EventType = "motion_detected"
	EventAIPerson       EventType = "ai_person"
	EventAIVehicle      EventType = "ai_vehicle"
	EventAIPet          EventType = "ai_pet"
	EventAudioAlarm     EventType = "audio_alarm"
	EventRecordingStart EventType = "recording_start"
	EventRecordingStop  EventType = "recording_stop"
	EventCameraOnline   EventType = "camera_online"
	EventCameraOffline  EventType = "camera_offline"
)

// EventSeverity represents the severity level of an event
type EventSeverity string

const (
	SeverityInfo     EventSeverity = "info"
	SeverityWarning  EventSeverity = "warning"
	SeverityCritical EventSeverity = "critical"
)

// Event represents a camera event
type Event struct {
	ID             string        `json:"id" db:"id"`
	CameraID       string        `json:"camera_id" db:"camera_id"`
	CameraName     string        `json:"camera_name" db:"camera_name"`
	Type           EventType     `json:"type" db:"type"`
	Severity       EventSeverity `json:"severity" db:"severity"`
	Timestamp      time.Time     `json:"timestamp" db:"timestamp"`
	Acknowledged   bool          `json:"acknowledged" db:"acknowledged"`
	AcknowledgedAt *time.Time    `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	Metadata       string        `json:"metadata,omitempty" db:"metadata"` // JSON string
	SnapshotPath   string        `json:"snapshot_path,omitempty" db:"snapshot_path"`
	VideoClipURL   string        `json:"video_clip_url,omitempty" db:"video_clip_url"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
}

// EventMetadata represents additional event information
type EventMetadata struct {
	Channel    int                    `json:"channel,omitempty"`
	Confidence float64                `json:"confidence,omitempty"`
	Region     []int                  `json:"region,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}
