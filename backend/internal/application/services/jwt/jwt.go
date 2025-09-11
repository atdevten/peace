package jwt

import (
	"github.com/atdevten/peace/internal/domain/value_objects"
)

// Service defines a technology-agnostic token service for auth
type Service interface {
	GenerateAccessToken(userID value_objects.UserID, email value_objects.Email) (string, error)
	GenerateRefreshToken(userID value_objects.UserID, email value_objects.Email) (string, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
}

// Claims are normalized token claims used across the application
type Claims struct {
	UserID string
	Email  string
	Type   string
	// Note: expiration and issued-at are validated inside the service implementation
}
