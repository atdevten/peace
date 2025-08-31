package value_objects

import (
	"testing"
)

func TestNewContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid content simple",
			input:       "This is a valid content",
			wantValue:   "This is a valid content",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid content with spaces",
			input:       "  This is a valid content  ",
			wantValue:   "This is a valid content",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid content with special characters",
			input:       "Content with @#$%^&*() characters!",
			wantValue:   "Content with @#$%^&*() characters!",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid content with newlines",
			input:       "Content with newlines and tabs",
			wantValue:   "Content with newlines and tabs",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid content maximum length",
			input:       string(make([]byte, 1000)),
			wantValue:   string(make([]byte, 1000)),
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty content",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "content with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "content with only tabs",
			input:       "\t\t\t",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "content with only newlines",
			input:       "\n\n\n",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "content too long",
			input:       string(make([]byte, 1001)),
			wantValue:   "",
			wantErr:     true,
			expectedErr: "content cannot exceed 1000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewContent(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewContent() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewContent() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewContent() unexpected error = %v", err)
				return
			}

			// Check value
			if got.Value() != tt.wantValue {
				t.Errorf("NewContent() = %v, want %v", got.Value(), tt.wantValue)
			}
		})
	}
}

func TestContent_Value(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple content",
			value: "This is content",
			want:  "This is content",
		},
		{
			name:  "content with special characters",
			value: "Content with @#$%^&*()!",
			want:  "Content with @#$%^&*()!",
		},
		{
			name:  "empty content",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Content{value: tt.value}
			if got := c.Value(); got != tt.want {
				t.Errorf("Content.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple content",
			value: "This is content",
			want:  "This is content",
		},
		{
			name:  "content with special characters",
			value: "Content with @#$%^&*()!",
			want:  "Content with @#$%^&*()!",
		},
		{
			name:  "empty content",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Content{value: tt.value}
			if got := c.String(); got != tt.want {
				t.Errorf("Content.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
