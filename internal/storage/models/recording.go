package models

import (
	"time"
)

// RecordingType represents the type of recording
type RecordingType string

const (
	RecordingTiming     RecordingType = "timing"
	RecordingMotion     RecordingType = "motion"
	RecordingAIPeople   RecordingType = "ai_people"
	RecordingAIVehicle  RecordingType = "ai_vehicle"
	RecordingAIPet      RecordingType = "ai_pet"
	RecordingManual     RecordingType = "manual"
)

// StreamType represents the stream type
type StreamType string

const (
	StreamMain StreamType = "main"
	StreamSub  StreamType = "sub"
	StreamExt  StreamType = "ext"
)

// Recording represents a camera recording
type Recording struct {
	ID            string        `json:"id" db:"id"`
	CameraID      string        `json:"camera_id" db:"camera_id"`
	FileName      string        `json:"file_name" db:"file_name"`
	FileSize      int64         `json:"file_size" db:"file_size"`
	StartTime     time.Time     `json:"start_time" db:"start_time"`
	EndTime       time.Time     `json:"end_time" db:"end_time"`
	Duration      int           `json:"duration" db:"duration"` // seconds
	StreamType    StreamType    `json:"stream_type" db:"stream_type"`
	RecordingType RecordingType `json:"recording_type" db:"recording_type"`
	StoragePath   string        `json:"storage_path" db:"storage_path"`
	ThumbnailURL  string        `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
}

// RecordingSearchRequest represents a request to search recordings
type RecordingSearchRequest struct {
	CameraID      *string        `json:"camera_id,omitempty"`
	StartTime     *time.Time     `json:"start_time,omitempty"`
	EndTime       *time.Time     `json:"end_time,omitempty"`
	RecordingType *RecordingType `json:"recording_type,omitempty"`
	StreamType    *StreamType    `json:"stream_type,omitempty"`
	Limit         int            `json:"limit,omitempty"`
	Offset        int            `json:"offset,omitempty"`
}

