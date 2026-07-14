package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"git.arcline.it/ArclineIT/nexus/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the custom JWT claims for Nexus.
type Claims struct {
	jwt.RegisteredClaims
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
}

// TokenPair holds both access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
	TokenType    string `json:"token_type"` // always "Bearer"
}

// GenerateTokenPair creates a new access + refresh token pair for a user.
func GenerateTokenPair(cfg *config.Config, userID uuid.UUID, email string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Auth.Issuer,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{cfg.Auth.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Auth.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		Email:     email,
		TokenType: "access",
	}

	accessToken, err := signToken(cfg, accessClaims)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Auth.Issuer,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{cfg.Auth.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Auth.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		Email:     email,
		TokenType: "refresh",
	}

	refreshToken, err := signToken(cfg, refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(cfg.Auth.AccessTokenDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken parses and validates a JWT token string.
func ValidateToken(cfg *config.Config, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.Auth.JWTSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// HashToken returns the SHA-256 hash of a token string (for storage).
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GenerateSecureToken creates a cryptographically random token.
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func signToken(cfg *config.Config, claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}
