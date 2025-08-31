package value_objects

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewMentalHealthRecordID(t *testing.T) {
	// Test that NewMentalHealthRecordID generates a valid UUID
	recordID := NewMentalHealthRecordID()

	// Check that it's not empty
	if recordID.String() == "" {
		t.Error("NewMentalHealthRecordID() generated empty string")
	}

	// Check that it's a valid UUID
	if _, err := uuid.Parse(recordID.String()); err != nil {
		t.Errorf("NewMentalHealthRecordID() generated invalid UUID: %v", err)
	}

	// Check that it's not zero
	if recordID.IsZero() {
		t.Error("NewMentalHealthRecordID() generated zero value")
	}
}

func TestNewMentalHealthRecordIDFromString(t *testing.T) {
	// Generate a valid UUID for testing
	validUUID := uuid.New().String()

	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid UUID",
			input:       validUUID,
			wantValue:   validUUID,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty string",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid UUID length: 0",
		},
		{
			name:        "invalid UUID format",
			input:       "not-a-uuid",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid UUID length: 10",
		},
		{
			name:        "invalid UUID with wrong characters",
			input:       "12345678-1234-1234-1234-123456789abc",
			wantValue:   "12345678-1234-1234-1234-123456789abc",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "UUID with wrong version",
			input:       "00000000-0000-0000-0000-000000000000",
			wantValue:   "00000000-0000-0000-0000-000000000000",
			wantErr:     false,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMentalHealthRecordIDFromString(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMentalHealthRecordIDFromString() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewMentalHealthRecordIDFromString() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewMentalHealthRecordIDFromString() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewMentalHealthRecordIDFromString() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestMentalHealthRecordID_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "valid UUID",
			value: "550e8400-e29b-41d4-a716-446655440000",
			want:  "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
		{
			name:  "random string",
			value: "test-string",
			want:  "test-string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MentalHealthRecordID{value: tt.value}
			if got := r.String(); got != tt.want {
				t.Errorf("MentalHealthRecordID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMentalHealthRecordID_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "zero when empty",
			value: "",
			want:  true,
		},
		{
			name:  "not zero when has value",
			value: "550e8400-e29b-41d4-a716-446655440000",
			want:  false,
		},
		{
			name:  "not zero when has random string",
			value: "test-string",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MentalHealthRecordID{value: tt.value}
			if got := r.IsZero(); got != tt.want {
				t.Errorf("MentalHealthRecordID.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
