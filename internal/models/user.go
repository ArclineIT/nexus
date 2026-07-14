package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an identity in the Nexus Control Panel.
// It is the single source of truth for all connected tools.
type User struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	DisplayName    string     `json:"display_name"`
	PasswordHash   string     `json:"-"` // never serialized
	MFAEnabled     bool       `json:"mfa_enabled"`
	MFASecret      string     `json:"-"` // never serialized
	EmailVerified  bool       `json:"email_verified"`
	Active         bool       `json:"active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Role represents a named set of permissions.
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole associates a user with a role.
type UserRole struct {
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Permission defines a granular action that can be allowed or denied.
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"`
}

// RolePermission associates a role with a permission.
type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

