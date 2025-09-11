package value_objects

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid email simple",
			input:       "test@example.com",
			wantValue:   "test@example.com",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid email with uppercase",
			input:       "TEST@EXAMPLE.COM",
			wantValue:   "test@example.com",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid email with spaces",
			input:       "  test@example.com  ",
			wantValue:   "test@example.com",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid email with dots",
			input:       "test.user@example.co.uk",
			wantValue:   "test.user@example.co.uk",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid email with plus",
			input:       "test+tag@example.com",
			wantValue:   "test+tag@example.com",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid email with underscore",
			input:       "test_user@example.com",
			wantValue:   "test_user@example.com",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty email",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "email cannot be empty",
		},
		{
			name:        "email with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "email cannot be empty",
		},
		{
			name:        "email too long",
			input:       "verylongemailaddress" + string(make([]byte, 240)) + "@example.com",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "email too long",
		},
		{
			name:        "invalid email missing @",
			input:       "testexample.com",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name:        "invalid email missing domain",
			input:       "test@",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name:        "invalid email missing local part",
			input:       "@example.com",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name:        "invalid email with invalid characters",
			input:       "test!@example.com",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name:        "invalid email with uppercase in domain",
			input:       "test@EXAMPLE.COM",
			wantValue:   "test@example.com",
			wantErr:     false,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewEmail() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewEmail() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewEmail() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewEmail() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple email",
			value: "test@example.com",
			want:  "test@example.com",
		},
		{
			name:  "email with dots",
			value: "test.user@example.co.uk",
			want:  "test.user@example.co.uk",
		},
		{
			name:  "empty email",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{value: tt.value}
			if got := e.String(); got != tt.want {
				t.Errorf("Email.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_IsZero(t *testing.T) {
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
			value: "test@example.com",
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
			e := &Email{value: tt.value}
			if got := e.IsZero(); got != tt.want {
				t.Errorf("Email.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
