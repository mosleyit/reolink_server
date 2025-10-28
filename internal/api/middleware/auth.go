package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// UsernameKey is the context key for username
	UsernameKey contextKey = "username"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Authenticate is a middleware that validates JWT tokens
func Authenticate(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			// Try to extract token from Authorization header first
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Check Bearer prefix
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					utils.RespondError(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must be 'Bearer {token}'", nil)
					return
				}
				tokenString = parts[1]
			} else {
				// For WebSocket connections, try query parameter
				tokenString = r.URL.Query().Get("token")
				if tokenString == "" {
					utils.RespondError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header or token query parameter is required", nil)
					return
				}
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				logger.Debug("JWT validation failed",
					zap.Error(err),
					zap.String("token", tokenString[:min(len(tokenString), 20)]+"..."),
				)
				utils.RespondError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token", nil)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				utils.RespondError(w, http.StatusUnauthorized, "INVALID_CLAIMS", "Invalid token claims", nil)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)

			// Continue with authenticated request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUsername extracts username from context
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
