package value_objects

import (
	"testing"
)

func TestNewLastName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid last name simple",
			input:       "Doe",
			wantValue:   "Doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name with spaces",
			input:       "  Doe  ",
			wantValue:   "Doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name with hyphen",
			input:       "Smith-Jones",
			wantValue:   "Smith-Jones",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name with apostrophe",
			input:       "O'Connor",
			wantValue:   "O'Connor",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name with dot",
			input:       "St. John",
			wantValue:   "St. John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name with multiple spaces",
			input:       "Van der Berg",
			wantValue:   "Van der Berg",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name minimum length",
			input:       "Li",
			wantValue:   "Li",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty last name",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name cannot be empty",
		},
		{
			name:        "last name with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name cannot be empty",
		},
		{
			name:        "last name too short",
			input:       "A",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name must be at least 2 characters long",
		},
		{
			name:        "last name too long",
			input:       "VeryLongLastNameThatExceedsTheMaximumAllowedLengthOfFiftyCharacters",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name cannot be longer than 50 characters",
		},
		{
			name:        "last name with invalid characters",
			input:       "Doe123",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name contains invalid characters",
		},
		{
			name:        "last name with special characters",
			input:       "Doe@Smith",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLastName(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLastName() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewLastName() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewLastName() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewLastName() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestNewOptionalLastName(t *testing.T) {
	tests := []struct {
		name        string
		input       *string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "nil input",
			input:       nil,
			wantValue:   "",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid last name",
			input:       stringPtr("Doe"),
			wantValue:   "Doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "invalid last name",
			input:       stringPtr(""),
			wantValue:   "",
			wantErr:     true,
			expectedErr: "last name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOptionalLastName(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOptionalLastName() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewOptionalLastName() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOptionalLastName() unexpected error = %v", err)
				return
			}

			// Check value
			if tt.input == nil {
				if got != nil {
					t.Errorf("NewOptionalLastName() = %v, want nil", got)
				}
			} else {
				if got.String() != tt.wantValue {
					t.Errorf("NewOptionalLastName() = %v, want %v", got.String(), tt.wantValue)
				}
			}
		})
	}
}

func TestLastName_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple name",
			value: "Doe",
			want:  "Doe",
		},
		{
			name:  "name with spaces",
			value: "Van der Berg",
			want:  "Van der Berg",
		},
		{
			name:  "empty name",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln := &LastName{value: tt.value}
			if got := ln.String(); got != tt.want {
				t.Errorf("LastName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastName_Value(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple name",
			value: "Doe",
			want:  "Doe",
		},
		{
			name:  "name with spaces",
			value: "Van der Berg",
			want:  "Van der Berg",
		},
		{
			name:  "empty name",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln := &LastName{value: tt.value}
			if got := ln.Value(); got != tt.want {
				t.Errorf("LastName.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastName_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "empty when empty",
			value: "",
			want:  true,
		},
		{
			name:  "not empty when has value",
			value: "Doe",
			want:  false,
		},
		{
			name:  "not empty when has spaces",
			value: "  ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln := &LastName{value: tt.value}
			if got := ln.IsEmpty(); got != tt.want {
				t.Errorf("LastName.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
