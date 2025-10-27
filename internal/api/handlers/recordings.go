package handlers

import (
	"net/http"

	"github.com/mosleyit/reolink_server/pkg/utils"
)

// ListRecordings handles GET /api/v1/recordings
func ListRecordings(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement recording listing
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Recording listing not yet implemented", nil)
}

// GetRecording handles GET /api/v1/recordings/{id}
func GetRecording(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get recording
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Get recording not yet implemented", nil)
}

// DownloadRecording handles GET /api/v1/recordings/{id}/download
func DownloadRecording(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement download recording
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Download recording not yet implemented", nil)
}

// SearchRecordings handles POST /api/v1/recordings/search
func SearchRecordings(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement search recordings
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Search recordings not yet implemented", nil)
}

