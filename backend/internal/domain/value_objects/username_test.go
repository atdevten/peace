package value_objects

import (
	"testing"
)

func TestNewUsername(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid username simple",
			input:       "john_doe",
			wantValue:   "john_doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid username with numbers",
			input:       "user123",
			wantValue:   "user123",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid username with uppercase",
			input:       "User123",
			wantValue:   "User123",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid username with spaces",
			input:       "  john_doe  ",
			wantValue:   "john_doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid username minimum length",
			input:       "abc",
			wantValue:   "abc",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid username maximum length",
			input:       "verylongusernamewithfiftycharactersexactly",
			wantValue:   "verylongusernamewithfiftycharactersexactly",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty username",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username cannot be empty",
		},
		{
			name:        "username with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username cannot be empty",
		},
		{
			name:        "username too short",
			input:       "ab",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username must be at least 3 characters",
		},
		{
			name:        "username too long",
			input:       "verylongusernamewithmorethanfiftycharactersexactly",
			wantValue:   "verylongusernamewithmorethanfiftycharactersexactly",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "username with hyphens",
			input:       "john-doe",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username can only contain letters, numbers, and underscores",
		},
		{
			name:        "username with dots",
			input:       "john.doe",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username can only contain letters, numbers, and underscores",
		},
		{
			name:        "username with special characters",
			input:       "john@doe",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username can only contain letters, numbers, and underscores",
		},
		{
			name:        "username with spaces in middle",
			input:       "john doe",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "username can only contain letters, numbers, and underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUsername(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUsername() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewUsername() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUsername() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewUsername() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestUsername_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple username",
			value: "john_doe",
			want:  "john_doe",
		},
		{
			name:  "username with numbers",
			value: "user123",
			want:  "user123",
		},
		{
			name:  "empty username",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Username{value: tt.value}
			if got := u.String(); got != tt.want {
				t.Errorf("Username.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsername_IsZero(t *testing.T) {
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
			value: "john_doe",
			want:  false,
		},
		{
			name:  "not zero when has spaces",
			value: "  ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Username{value: tt.value}
			if got := u.IsZero(); got != tt.want {
				t.Errorf("Username.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
