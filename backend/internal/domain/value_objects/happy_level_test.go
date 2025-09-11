package value_objects

import (
	"testing"
)

func TestNewHappyLevel(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		wantValue   int
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid happy level minimum",
			input:       HappyLevelMin,
			wantValue:   HappyLevelMin,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid happy level maximum",
			input:       HappyLevelMax,
			wantValue:   HappyLevelMax,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid happy level middle",
			input:       5,
			wantValue:   5,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "invalid happy level below minimum",
			input:       0,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "happy level must be between 1 and 10",
		},
		{
			name:        "invalid happy level above maximum",
			input:       11,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "happy level must be between 1 and 10",
		},
		{
			name:        "invalid happy level negative",
			input:       -1,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "happy level must be between 1 and 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHappyLevel(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewHappyLevel() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewHappyLevel() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewHappyLevel() unexpected error = %v", err)
				return
			}

			// Check value
			if got.Value() != tt.wantValue {
				t.Errorf("NewHappyLevel() = %v, want %v", got.Value(), tt.wantValue)
			}
		})
	}
}

func TestHappyLevel_Value(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "value 1",
			value: 1,
			want:  1,
		},
		{
			name:  "value 5",
			value: 5,
			want:  5,
		},
		{
			name:  "value 10",
			value: 10,
			want:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &HappyLevel{value: tt.value}
			if got := hl.Value(); got != tt.want {
				t.Errorf("HappyLevel.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHappyLevel_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{
			name:  "empty when value is 0",
			value: 0,
			want:  true,
		},
		{
			name:  "not empty when value is 1",
			value: 1,
			want:  false,
		},
		{
			name:  "not empty when value is 5",
			value: 5,
			want:  false,
		},
		{
			name:  "not empty when value is 10",
			value: 10,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &HappyLevel{value: tt.value}
			if got := hl.IsEmpty(); got != tt.want {
				t.Errorf("HappyLevel.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
