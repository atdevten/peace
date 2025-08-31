package middleware

import (
	"fmt"
	"net/http"
	"strings"

	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	"github.com/atdevten/peace/internal/domain/value_objects"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides Gin-compatible middleware functions
type AuthMiddleware struct {
	jwtService appjwt.Service
}

func NewAuthMiddleware(jwtService appjwt.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth is a Gin middleware that requires JWT authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate access token
		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := value_objects.NewUserIDFromString(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		// Parse email
		email, err := value_objects.NewEmail(claims.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid email in token",
			})
			c.Abort()
			return
		}

		// Store user info in Gin context (store pointer to match getters)
		c.Set("user_id", userID)
		c.Set("user_email", email)

		// Continue to next handler
		c.Next()
	}
}

// OptionalAuth is a Gin middleware that optionally validates JWT if present
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without user info
			c.Next()
			return
		}

		// Check if it starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, continue without user info
			c.Next()
			return
		}

		tokenString := parts[1]

		// Validate access token
		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			// Invalid token, continue without user info
			c.Next()
			return
		}

		// Parse user ID
		userID, err := value_objects.NewUserIDFromString(claims.UserID)
		if err != nil {
			// Invalid user ID, continue without user info
			c.Next()
			return
		}

		// Parse email
		email, err := value_objects.NewEmail(claims.Email)
		if err != nil {
			// Invalid email, continue without user info
			c.Next()
			return
		}

		// Store user info in Gin context
		c.Set("user_id", userID)
		c.Set("user_email", email)

		// Continue to next handler
		c.Next()
	}
}

// Helper functions to extract user info from Gin context
func GetUserIDFromGinContext(c *gin.Context) (*value_objects.UserID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, false
	}

	fmt.Println("userID", userID)

	userIDValue, ok := userID.(*value_objects.UserID)
	if !ok {
		return nil, false
	}

	return userIDValue, ok
}

func GetUserEmailFromGinContext(c *gin.Context) (*value_objects.Email, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return nil, false
	}

	emailValue, ok := email.(*value_objects.Email)

	return emailValue, ok
}
