package value_objects

import (
	"fmt"
	"strings"
)

type Content struct {
	value string
}

func NewContent(value string) (*Content, error) {
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	if len(value) > 1000 {
		return nil, fmt.Errorf("content cannot exceed 1000 characters")
	}

	return &Content{
		value: strings.TrimSpace(value),
	}, nil
}

func (c *Content) Value() string {
	return c.value
}

func (c *Content) String() string {
	return c.value
}
