package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all configuration for the Nexus Control Panel.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	Audience             string
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string
	Format string // "json" or "text"
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            getEnv("NEXUS_HOST", "0.0.0.0"),
			Port:            getEnv("NEXUS_PORT", "8080"),
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv("NEXUS_DB_HOST", "localhost"),
			Port:     getEnv("NEXUS_DB_PORT", "5432"),
			User:     getEnv("NEXUS_DB_USER", "nexus"),
			Password: getEnv("NEXUS_DB_PASSWORD", ""),
			Name:     getEnv("NEXUS_DB_NAME", "nexus"),
			SSLMode:  getEnv("NEXUS_DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:            getEnv("NEXUS_JWT_SECRET", "change-me-in-production"),
			AccessTokenDuration:  15 * time.Minute,
			RefreshTokenDuration: 7 * 24 * time.Hour,
			Issuer:               getEnv("NEXUS_AUTH_ISSUER", "nexus.arcline.it"),
			Audience:             getEnv("NEXUS_AUTH_AUDIENCE", "nexus.arcline.it"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("NEXUS_LOG_LEVEL", "info"),
			Format: getEnv("NEXUS_LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
