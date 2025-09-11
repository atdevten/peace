package value_objects

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(email string) (*Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if len(email) > 255 {
		return nil, errors.New("email too long")
	}

	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	return &Email{value: email}, nil
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) IsZero() bool {
	return e.value == ""
}
