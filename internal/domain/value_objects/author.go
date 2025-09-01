package value_objects

import (
	"fmt"
	"strings"
)

type Author struct {
	value string
}

func NewAuthor(value string) (*Author, error) {
	trimmedValue := strings.TrimSpace(value)

	if trimmedValue == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	// Removed 255 character limit since we're using TEXT type
	// Allow reasonable maximum to prevent abuse
	if len(trimmedValue) > 2000 {
		return nil, fmt.Errorf("author cannot exceed 2000 characters")
	}

	return &Author{
		value: trimmedValue,
	}, nil
}

func (a *Author) Value() string {
	return a.value
}

func (a *Author) String() string {
	return a.value
}
