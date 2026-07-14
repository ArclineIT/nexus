package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an authenticated user session.
type Session struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	TokenHash    string     `json:"-"`
	RefreshToken string     `json:"-"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ConnectedApp represents an external tool integrated via SSO.
// These are the "spokes" of the Nexus hub.
type ConnectedApp struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description,omitempty"`
	HomepageURL      string    `json:"homepage_url"`
	SSOCallbackURL   string    `json:"sso_callback_url"`
	ClientID         string    `json:"client_id"`
	ClientSecretHash string    `json:"-"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AuditLog records a security-relevant event for compliance and debugging.
type AuditLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Metadata  string    `json:"metadata,omitempty"` // JSON-encoded extra context
	CreatedAt time.Time `json:"created_at"`
}

