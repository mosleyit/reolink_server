package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondBadRequest(w, "Invalid request body", nil)
		return
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		utils.RespondBadRequest(w, "Username and password are required", nil)
		return
	}

	// Authenticate user
	loginResp, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password", nil)
		return
	}

	// Return response
	utils.RespondJSON(w, http.StatusOK, loginResp)
}
