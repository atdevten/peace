package value_objects

import (
	"errors"
	"regexp"
	"strings"
)

type Username struct {
	value string
}

func NewUsername(username string) (*Username, error) {
	username = strings.TrimSpace(username)

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	if len(username) > 50 {
		return nil, errors.New("username too long")
	}

	// Username should contain only alphanumeric characters and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return nil, errors.New("username can only contain letters, numbers, and underscores")
	}

	return &Username{value: username}, nil
}

func (u *Username) String() string {
	return u.value
}

func (u *Username) IsZero() bool {
	return u.value == ""
}
