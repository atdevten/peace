package value_objects

import (
	"testing"
)

func TestNewTagName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid tag name simple",
			input:       "motivation",
			wantValue:   "motivation",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name with spaces",
			input:       "  motivation  ",
			wantValue:   "motivation",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name with hyphens",
			input:       "self-improvement",
			wantValue:   "self-improvement",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name with underscores",
			input:       "mental_health",
			wantValue:   "mental_health",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name with numbers",
			input:       "motivation2024",
			wantValue:   "motivation2024",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name minimum length",
			input:       "ab",
			wantValue:   "ab",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid tag name maximum length",
			input:       "verylongtagnamewiththirtycharacters",
			wantValue:   "verylongtagnamewiththirtycharacters",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty tag name",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "tag name cannot be empty",
		},
		{
			name:        "tag name with only spaces",
			input:       "   ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "tag name cannot be empty",
		},
		{
			name:        "tag name too short",
			input:       "a",
			wantValue:   "a",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "tag name too long",
			input:       "verylongtagnamewithmorethanthirtycharactersexactly",
			wantValue:   "verylongtagnamewithmorethanthirtycharactersexactly",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "tag name with special characters",
			input:       "motivation@2024",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "tag name contains invalid characters",
		},
		{
			name:        "tag name with spaces in middle",
			input:       "mental health",
			wantValue:   "mental health",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "tag name with dots",
			input:       "motivation.2024",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "tag name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTagName(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTagName() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewTagName() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTagName() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewTagName() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestTagName_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple tag name",
			value: "motivation",
			want:  "motivation",
		},
		{
			name:  "tag name with hyphens",
			value: "self-improvement",
			want:  "self-improvement",
		},
		{
			name:  "tag name with underscores",
			value: "mental_health",
			want:  "mental_health",
		},
		{
			name:  "empty tag name",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn := &TagName{value: tt.value}
			if got := tn.String(); got != tt.want {
				t.Errorf("TagName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
