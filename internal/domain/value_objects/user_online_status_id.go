package value_objects

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// UserOnlineStatusID represents a unique identifier for user online status
type UserOnlineStatusID struct {
	value string
}

// NewUserOnlineStatusID creates a new UserOnlineStatusID
func NewUserOnlineStatusID() *UserOnlineStatusID {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("failed to generate user online status ID: %v", err))
	}
	return &UserOnlineStatusID{
		value: hex.EncodeToString(bytes),
	}
}

// NewUserOnlineStatusIDFromString creates a UserOnlineStatusID from string
func NewUserOnlineStatusIDFromString(value string) (*UserOnlineStatusID, error) {
	if value == "" {
		return nil, fmt.Errorf("user online status ID cannot be empty")
	}

	if len(value) != 32 {
		return nil, fmt.Errorf("user online status ID must be 32 characters long")
	}

	return &UserOnlineStatusID{value: value}, nil
}

// String returns the string representation of the ID
func (u *UserOnlineStatusID) String() string {
	return u.value
}

// Value returns the underlying value
func (u *UserOnlineStatusID) Value() string {
	return u.value
}

// Equals checks if two IDs are equal
func (u *UserOnlineStatusID) Equals(other *UserOnlineStatusID) bool {
	if other == nil {
		return false
	}
	return u.value == other.value
}
