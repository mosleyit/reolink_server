package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "simple string data",
			data:           "test message",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"success":true,"data":"test message"`,
		},
		{
			name: "struct data",
			data: map[string]string{
				"key": "value",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"success":true,"data":{"key":"value"}`,
		},
		{
			name:           "nil data",
			data:           nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"success":true,"timestamp"`, // nil data is omitted from JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondJSON(w, tt.expectedStatus, tt.data)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			// Verify it's valid JSON
			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.True(t, response.Success)
		})
	}
}

func TestRespondError(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		code           string
		message        string
		details        interface{}
		expectedStatus int
	}{
		{
			name:           "bad request error",
			status:         http.StatusBadRequest,
			code:           "INVALID_INPUT",
			message:        "Invalid input provided",
			details:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error with details",
			status:         http.StatusNotFound,
			code:           "NOT_FOUND",
			message:        "Resource not found",
			details:        map[string]string{"id": "123"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal server error",
			status:         http.StatusInternalServerError,
			code:           "INTERNAL_ERROR",
			message:        "Something went wrong",
			details:        nil,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondError(w, tt.status, tt.code, tt.message, tt.details)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.False(t, response.Success)
			assert.NotNil(t, response.Error)
			assert.Equal(t, tt.code, response.Error.Code)
			assert.Equal(t, tt.message, response.Error.Message)
		})
	}
}

func TestRespondPaginated(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	page := 1
	limit := 10
	total := 3

	w := httptest.NewRecorder()
	RespondPaginated(w, http.StatusOK, data, page, limit, total)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// The response is wrapped in a Response struct
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)

	// Extract the paginated data from response.Data
	dataBytes, err := json.Marshal(response.Data)
	require.NoError(t, err)

	var paginatedResp PaginatedResponse
	err = json.Unmarshal(dataBytes, &paginatedResp)
	require.NoError(t, err)

	assert.NotNil(t, paginatedResp.Pagination)
	assert.Equal(t, page, paginatedResp.Pagination.Page)
	assert.Equal(t, limit, paginatedResp.Pagination.Limit)
	assert.Equal(t, total, paginatedResp.Pagination.Total)
	assert.Equal(t, 1, paginatedResp.Pagination.TotalPages)
}

func TestRespondCreated(t *testing.T) {
	data := map[string]string{"id": "123"}

	w := httptest.NewRecorder()
	RespondCreated(w, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

func TestRespondNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	RespondNoContent(w)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestRespondBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	RespondBadRequest(w, "Invalid input", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "BAD_REQUEST", response.Error.Code)
	assert.Equal(t, "Invalid input", response.Error.Message)
}

func TestRespondUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	RespondUnauthorized(w, "Invalid credentials")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "UNAUTHORIZED", response.Error.Code)
	assert.Equal(t, "Invalid credentials", response.Error.Message)
}

func TestRespondNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	RespondNotFound(w, "Resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "NOT_FOUND", response.Error.Code)
	assert.Equal(t, "Resource not found", response.Error.Message)
}

func TestPaginationInfo_TotalPages(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		limit     int
		wantPages int
	}{
		{
			name:      "exact division",
			total:     100,
			limit:     10,
			wantPages: 10,
		},
		{
			name:      "with remainder",
			total:     105,
			limit:     10,
			wantPages: 11,
		},
		{
			name:      "less than one page",
			total:     5,
			limit:     10,
			wantPages: 1,
		},
		{
			name:      "zero total",
			total:     0,
			limit:     10,
			wantPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := &PaginationInfo{
				Total: tt.total,
				Limit: tt.limit,
			}
			totalPages := (pagination.Total + pagination.Limit - 1) / pagination.Limit
			if pagination.Total == 0 {
				totalPages = 0
			}
			assert.Equal(t, tt.wantPages, totalPages)
		})
	}
}

func TestResponse_Structure(t *testing.T) {
	response := Response{
		Success:   true,
		Data:      "test data",
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.Equal(t, "test data", response.Data)
	assert.Nil(t, response.Error)
	assert.NotZero(t, response.Timestamp)
}

func TestErrorInfo_Structure(t *testing.T) {
	errorInfo := ErrorInfo{
		Code:    "TEST_ERROR",
		Message: "Test error message",
	}

	assert.Equal(t, "TEST_ERROR", errorInfo.Code)
	assert.Equal(t, "Test error message", errorInfo.Message)
}

// Benchmark tests
func BenchmarkRespondJSON(b *testing.B) {
	data := map[string]string{"key": "value"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		RespondJSON(w, http.StatusOK, data)
	}
}

func BenchmarkRespondError(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		RespondError(w, http.StatusBadRequest, "ERROR", "Error message", nil)
	}
}

func BenchmarkRespondPaginated(b *testing.B) {
	data := []string{"item1", "item2", "item3"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		RespondPaginated(w, http.StatusOK, data, 1, 10, 3)
	}
}
