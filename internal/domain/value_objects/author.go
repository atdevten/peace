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

	if len(trimmedValue) > 255 {
		return nil, fmt.Errorf("author cannot exceed 255 characters")
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
