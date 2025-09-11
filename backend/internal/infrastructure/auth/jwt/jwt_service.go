package jwt

import (
	"errors"
	"time"

	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	"github.com/atdevten/peace/internal/domain/value_objects"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewService(secretKey string, accessExpiry time.Duration, refreshExpiry time.Duration) appjwt.Service {
	return &jwtService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func (s *jwtService) generateToken(userID value_objects.UserID, email value_objects.Email, tokenType string, expiry time.Duration) (string, error) {
	claims := &jwtClaims{
		UserID: userID.String(),
		Email:  email.String(),
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GenerateAccessToken(userID value_objects.UserID, email value_objects.Email) (string, error) {
	return s.generateToken(userID, email, "access", s.accessExpiry)
}

func (s *jwtService) GenerateRefreshToken(userID value_objects.UserID, email value_objects.Email) (string, error) {
	return s.generateToken(userID, email, "refresh", s.refreshExpiry)
}

func (s *jwtService) validateAndCheckType(tokenString string, expectedType string) (*appjwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		if claims.Type != expectedType {
			return nil, errors.New("invalid token type")
		}
		return &appjwt.Claims{UserID: claims.UserID, Email: claims.Email, Type: claims.Type}, nil
	}

	return nil, errors.New("invalid token")
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*appjwt.Claims, error) {
	return s.validateAndCheckType(tokenString, "access")
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*appjwt.Claims, error) {
	return s.validateAndCheckType(tokenString, "refresh")
}
