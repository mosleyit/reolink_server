package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHealthHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := NewHealthHandler(db)
	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "returns healthy status",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "healthy", body["status"])
				assert.Equal(t, "1.0.0", body["version"])
				assert.NotEmpty(t, body["uptime"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			handler := NewHealthHandler(db)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			handler.HealthCheck(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// Parse response
			var response struct {
				Success bool                   `json:"success"`
				Data    map[string]interface{} `json:"data"`
			}
			err = parseJSON(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.True(t, response.Success)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response.Data)
			}
		})
	}
}

func TestHealthHandler_ReadinessCheck(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		dbNil          bool
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "database healthy",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "ready", body["status"])
				components := body["components"].(map[string]interface{})
				assert.Equal(t, "healthy", components["database"])
			},
		},
		{
			name: "database unhealthy",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(errors.New("connection refused"))
			},
			expectedStatus: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "degraded", body["status"])
				components := body["components"].(map[string]interface{})
				assert.Contains(t, components["database"], "unhealthy")
				assert.Contains(t, components["database"], "connection refused")
			},
		},
		{
			name:           "database not configured",
			dbNil:          true,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "ready", body["status"])
				components := body["components"].(map[string]interface{})
				assert.Equal(t, "not configured", components["database"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var handler *HealthHandler
			var mock sqlmock.Sqlmock

			if tt.dbNil {
				handler = NewHealthHandler(nil)
			} else {
				db, m, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
				require.NoError(t, err)
				defer db.Close()
				mock = m

				if tt.setupMock != nil {
					tt.setupMock(mock)
				}

				handler = NewHealthHandler(db)
			}

			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			w := httptest.NewRecorder()

			handler.ReadinessCheck(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// Parse response
			var response struct {
				Success bool                   `json:"success"`
				Data    map[string]interface{} `json:"data"`
			}
			err := parseJSON(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response.Success)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response.Data)
			}

			// Verify all expectations were met
			if mock != nil {
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			}
		})
	}
}

func TestHealthHandler_ReadinessCheck_Timeout(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	// Simulate a slow database that will timeout
	mock.ExpectPing().WillDelayFor(5 * 1000000000) // 5 seconds

	handler := NewHealthHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ReadinessCheck(w, req)

	// Should return degraded status due to timeout
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err = parseJSON(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response.Data
	assert.Equal(t, "degraded", data["status"])
	components := data["components"].(map[string]interface{})
	assert.Contains(t, components["database"], "unhealthy")
}

// Helper function to parse JSON response
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
