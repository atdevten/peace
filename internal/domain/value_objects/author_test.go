package value_objects

import (
	"testing"
)

func TestNewAuthor(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid author simple",
			input:       "John Doe",
			wantValue:   "John Doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid author with spaces",
			input:       "  John Doe  ",
			wantValue:   "John Doe",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid author with special characters",
			input:       "Jean-Pierre O'Connor",
			wantValue:   "Jean-Pierre O'Connor",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid author with dots",
			input:       "St. John",
			wantValue:   "St. John",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty author",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "author cannot be empty",
		},
		{
			name:        "author with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "author cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthor(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewAuthor() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewAuthor() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewAuthor() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewAuthor() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestAuthor_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple author",
			value: "John Doe",
			want:  "John Doe",
		},
		{
			name:  "author with special characters",
			value: "Jean-Pierre O'Connor",
			want:  "Jean-Pierre O'Connor",
		},
		{
			name:  "empty author",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Author{value: tt.value}
			if got := a.String(); got != tt.want {
				t.Errorf("Author.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
