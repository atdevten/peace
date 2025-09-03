package value_objects

import (
	"errors"
	"strings"
)

// FirstName represents a user's first name
type FirstName struct {
	value string
}

// NewFirstName creates a new FirstName value object
func NewFirstName(value string) (*FirstName, error) {
	if value == "" {
		return nil, errors.New("first name cannot be empty")
	}

	// Trim whitespace
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("first name cannot be empty")
	}

	// Check length constraints
	if len(trimmed) < 2 {
		return nil, errors.New("first name must be at least 2 characters long")
	}

	if len(trimmed) > 50 {
		return nil, errors.New("first name cannot be longer than 50 characters")
	}

	// Check for invalid characters (only letters, spaces, hyphens, and apostrophes)
	for _, char := range trimmed {
		if !isValidNameCharacter(char) {
			return nil, errors.New("first name contains invalid characters")
		}
	}

	return &FirstName{value: trimmed}, nil
}

// NewOptionalFirstName creates a new FirstName value object from optional string
func NewOptionalFirstName(value *string) (*FirstName, error) {
	if value == nil {
		return nil, nil
	}

	firstName, err := NewFirstName(*value)
	if err != nil {
		return nil, err
	}

	return firstName, nil
}

// String returns the string representation of the first name
func (f *FirstName) String() string {
	return f.value
}

// Value returns the underlying value
func (f *FirstName) Value() string {
	return f.value
}

// IsEmpty checks if the first name is empty
func (f *FirstName) IsEmpty() bool {
	return f.value == ""
}
