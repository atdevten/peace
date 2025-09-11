package value_objects

import (
	"testing"
)

func TestNewTagID(t *testing.T) {
	// Test that NewTagID creates a tag ID with value 0
	tagID := NewTagID()

	// Check that it's not nil
	if tagID == nil {
		t.Error("NewTagID() returned nil")
	}

	// Check that it has value 0
	if tagID.IntValue() != 0 {
		t.Errorf("NewTagID() = %v, want 0", tagID.IntValue())
	}

	// Check string representation
	if tagID.String() != "0" {
		t.Errorf("NewTagID().String() = %v, want '0'", tagID.String())
	}
}

func TestNewTagIDFromInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "positive integer",
			input: 123,
			want:  123,
		},
		{
			name:  "zero",
			input: 0,
			want:  0,
		},
		{
			name:  "negative integer",
			input: -1,
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTagIDFromInt(tt.input)
			if got.IntValue() != tt.want {
				t.Errorf("NewTagIDFromInt() = %v, want %v", got.IntValue(), tt.want)
			}
		})
	}
}

func TestTagID_IntValue(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "positive integer",
			value: 123,
			want:  123,
		},
		{
			name:  "zero",
			value: 0,
			want:  0,
		},
		{
			name:  "negative integer",
			value: -1,
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagID := &TagID{value: tt.value}
			if got := tagID.IntValue(); got != tt.want {
				t.Errorf("TagID.IntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagID_String(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{
			name:  "positive integer",
			value: 123,
			want:  "123",
		},
		{
			name:  "zero",
			value: 0,
			want:  "0",
		},
		{
			name:  "negative integer",
			value: -1,
			want:  "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagID := &TagID{value: tt.value}
			if got := tagID.String(); got != tt.want {
				t.Errorf("TagID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
