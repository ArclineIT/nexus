package handler

import (
	"encoding/json"
	"net/http"

	"git.arcline.it/ArclineIT/nexus/internal/auth"
	"git.arcline.it/ArclineIT/nexus/internal/config"
	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// LoginRequest is the expected body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /auth/login.
// TODO: validate credentials against database.
func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Email == "" || req.Password == "" {
			respondError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		// TODO: verify password, look up user from DB
		// For now, generate a token pair with a placeholder user ID
		userID := uuid.New()

		tokens, err := auth.GenerateTokenPair(h.cfg, userID, req.Email)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to generate tokens")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

// RefreshRequest is the expected body for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.RefreshToken == "" {
			respondError(w, http.StatusBadRequest, "refresh_token is required")
			return
		}

		// Validate the refresh token
		claims, err := auth.ValidateToken(h.cfg, req.RefreshToken)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}

		if claims.TokenType != "refresh" {
			respondError(w, http.StatusUnauthorized, "token is not a refresh token")
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "invalid user ID in token")
			return
		}

		tokens, err := auth.GenerateTokenPair(h.cfg, userID, claims.Email)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to generate tokens")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

// Me handles GET /auth/me — returns the authenticated user's info.
func (h *AuthHandler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The Authenticate middleware has already injected user info into context.
		userID := r.Context().Value("user_id")
		userEmail := r.Context().Value("user_email")

		respondJSON(w, http.StatusOK, map[string]any{
			"id":    userID,
			"email": userEmail,
		})
	}
}

