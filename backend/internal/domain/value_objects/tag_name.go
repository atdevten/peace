package value_objects

import (
	"errors"
	"strings"
)

type TagName struct {
	value string
}

var (
	ErrTagNameEmpty   = errors.New("tag name cannot be empty")
	ErrTagNameTooLong = errors.New("tag name cannot exceed 100 characters")
	ErrTagNameInvalid = errors.New("tag name contains invalid characters")
)

func NewTagName(value string) (*TagName, error) {
	if err := validateTagName(value); err != nil {
		return nil, err
	}

	return &TagName{value: strings.TrimSpace(value)}, nil
}

func validateTagName(value string) error {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return ErrTagNameEmpty
	}

	if len(trimmed) > 100 {
		return ErrTagNameTooLong
	}

	// Check for invalid characters (only allow alphanumeric, spaces, hyphens, underscores)
	for _, char := range trimmed {
		if !isValidTagNameChar(char) {
			return ErrTagNameInvalid
		}
	}

	return nil
}

func isValidTagNameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == ' ' ||
		char == '-' ||
		char == '_'
}

func (t *TagName) Value() string {
	return t.value
}

func (t *TagName) String() string {
	return t.value
}
