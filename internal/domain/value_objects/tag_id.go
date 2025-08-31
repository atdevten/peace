package value_objects

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type TagID struct {
	value int
}

func NewTagID() *TagID {
	return &TagID{value: 0}
}

func NewTagIDFromInt(value int) *TagID {
	return &TagID{value: value}
}

func (t *TagID) IntValue() int {
	return t.value
}

func (t *TagID) String() string {
	return strconv.Itoa(t.value)
}

// Value implements driver.Valuer interface for database serialization
func (t *TagID) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return int64(t.value), nil
}

// Scan implements sql.Scanner interface for database deserialization
func (t *TagID) Scan(value interface{}) error {
	if value == nil {
		t.value = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		t.value = int(v)
	case int32:
		t.value = int(v)
	case int:
		t.value = v
	case string:
		if v == "" {
			t.value = 0
			return nil
		}
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("cannot scan string %q as TagID: %w", v, err)
		}
		t.value = parsed
	default:
		return fmt.Errorf("cannot scan %T as TagID", value)
	}
	return nil
}
