package value_objects

import (
	"testing"
)

func TestNewFirstName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid first name simple",
			input:       "John",
			wantValue:   "John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid first name with spaces",
			input:       "  John  ",
			wantValue:   "John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid first name with hyphen",
			input:       "Jean-Pierre",
			wantValue:   "Jean-Pierre",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid first name with apostrophe",
			input:       "O'Connor",
			wantValue:   "O'Connor",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid first name with dot",
			input:       "St. John",
			wantValue:   "St. John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid first name with multiple spaces",
			input:       "Mary Jane",
			wantValue:   "Mary Jane",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty first name",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name cannot be empty",
		},
		{
			name:        "first name with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name cannot be empty",
		},
		{
			name:        "first name too short",
			input:       "A",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name must be at least 2 characters long",
		},
		{
			name:        "first name too long",
			input:       "VeryLongFirstNameThatExceedsTheMaximumAllowedLength",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name cannot be longer than 50 characters",
		},
		{
			name:        "first name with invalid characters",
			input:       "John123",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name contains invalid characters",
		},
		{
			name:        "first name with special characters",
			input:       "John@Doe",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFirstName(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewFirstName() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewFirstName() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewFirstName() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewFirstName() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestNewOptionalFirstName(t *testing.T) {
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
			name:        "valid first name",
			input:       stringPtr("John"),
			wantValue:   "John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "invalid first name",
			input:       stringPtr(""),
			wantValue:   "",
			wantErr:     true,
			expectedErr: "first name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOptionalFirstName(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOptionalFirstName() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewOptionalFirstName() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOptionalFirstName() unexpected error = %v", err)
				return
			}

			// Check value
			if tt.input == nil {
				if got != nil {
					t.Errorf("NewOptionalFirstName() = %v, want nil", got)
				}
			} else {
				if got.String() != tt.wantValue {
					t.Errorf("NewOptionalFirstName() = %v, want %v", got.String(), tt.wantValue)
				}
			}
		})
	}
}

func TestFirstName_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple name",
			value: "John",
			want:  "John",
		},
		{
			name:  "name with spaces",
			value: "Mary Jane",
			want:  "Mary Jane",
		},
		{
			name:  "empty name",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := &FirstName{value: tt.value}
			if got := fn.String(); got != tt.want {
				t.Errorf("FirstName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstName_Value(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple name",
			value: "John",
			want:  "John",
		},
		{
			name:  "name with spaces",
			value: "Mary Jane",
			want:  "Mary Jane",
		},
		{
			name:  "empty name",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := &FirstName{value: tt.value}
			if got := fn.Value(); got != tt.want {
				t.Errorf("FirstName.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstName_IsEmpty(t *testing.T) {
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
			value: "John",
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
			fn := &FirstName{value: tt.value}
			if got := fn.IsEmpty(); got != tt.want {
				t.Errorf("FirstName.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
