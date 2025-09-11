package value_objects

import (
	"errors"
	"strings"
)

// LastName represents a user's last name
type LastName struct {
	value string
}

// NewLastName creates a new LastName value object
func NewLastName(value string) (*LastName, error) {
	if value == "" {
		return nil, errors.New("last name cannot be empty")
	}

	// Trim whitespace
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("last name cannot be empty")
	}

	// Check length constraints
	if len(trimmed) < 2 {
		return nil, errors.New("last name must be at least 2 characters long")
	}

	if len(trimmed) > 50 {
		return nil, errors.New("last name cannot be longer than 50 characters")
	}

	// Check for invalid characters (only letters, spaces, hyphens, and apostrophes)
	for _, char := range trimmed {
		if !isValidNameCharacter(char) {
			return nil, errors.New("last name contains invalid characters")
		}
	}

	return &LastName{value: trimmed}, nil
}

// NewOptionalLastName creates a new LastName value object from optional string
func NewOptionalLastName(value *string) (*LastName, error) {
	if value == nil {
		return nil, nil
	}

	lastName, err := NewLastName(*value)
	if err != nil {
		return nil, err
	}

	return lastName, nil
}

// String returns the string representation of the last name
func (l *LastName) String() string {
	return l.value
}

// Value returns the underlying value
func (l *LastName) Value() string {
	return l.value
}

// IsEmpty checks if the last name is empty
func (l *LastName) IsEmpty() bool {
	return l.value == ""
}
