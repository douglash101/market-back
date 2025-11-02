package middleware

import (
	"context"
	"crypto/subtle"
	"market/internal/domain/user"
	"market/pkg/security"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// SecureMiddleware provides authentication and security headers for API endpoints
func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		setSecurityHeaders(w)

		// Handle CORS preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Validate authentication
		if !isAuthenticated(r) {
			http.Error(w, `{"error":"Unauthorized","message":"Invalid or missing authentication token"}`, http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware provides a simple function-based middleware for direct use
func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		setSecurityHeaders(w)

		// Handle CORS preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		tokenStr := r.Header.Get("Authorization")
		w.Header().Add("Content-Type", "application/json")
		if tokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token é necessário"})
			return
		}

		// Remove "Bearer " do token se presente
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		claims := &user.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return []byte("sua_chave_secreta_aqui"), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token inválido"})
			return
		}

		// Pegando o "sub" do JWT
		userID := claims.UserID
		if userID.String() == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Usuário não identificado"})
			return
		}

		user := security.UserAuth{
			UserID:    userID,
			CompanyID: claims.CompanyID,
		}

		ctx := context.WithValue(r.Context(), security.USER_KEY, user)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

// setSecurityHeaders adds security headers to the response
func setSecurityHeaders(w http.ResponseWriter) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "3600")

	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

	// API response headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// isAuthenticated validates the Bearer token against the SECRET_KEY
func isAuthenticated(r *http.Request) bool {

	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	// Check Bearer token format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(token), []byte("")) == 1
}

// HealthCheckMiddleware provides a simple middleware for health check endpoints (no auth required)
func HealthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set minimal security headers for health checks
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")

		// Add timestamp for monitoring
		w.Header().Set("X-Timestamp", time.Now().UTC().Format(time.RFC3339))

		next.ServeHTTP(w, r)
	})
}
