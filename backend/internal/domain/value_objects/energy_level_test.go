package value_objects

import (
	"testing"
)

func TestNewEnergyLevel(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		wantValue   int
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid energy level minimum",
			input:       EnergyLevelMin,
			wantValue:   EnergyLevelMin,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid energy level maximum",
			input:       EnergyLevelMax,
			wantValue:   EnergyLevelMax,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid energy level middle",
			input:       5,
			wantValue:   5,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "invalid energy level below minimum",
			input:       0,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "energy level must be between 1 and 10",
		},
		{
			name:        "invalid energy level above maximum",
			input:       11,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "energy level must be between 1 and 10",
		},
		{
			name:        "invalid energy level negative",
			input:       -1,
			wantValue:   0,
			wantErr:     true,
			expectedErr: "energy level must be between 1 and 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnergyLevel(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewEnergyLevel() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewEnergyLevel() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewEnergyLevel() unexpected error = %v", err)
				return
			}

			// Check value
			if got.Value() != tt.wantValue {
				t.Errorf("NewEnergyLevel() = %v, want %v", got.Value(), tt.wantValue)
			}
		})
	}
}

func TestEnergyLevel_Value(t *testing.T) {
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
			el := &EnergyLevel{value: tt.value}
			if got := el.Value(); got != tt.want {
				t.Errorf("EnergyLevel.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnergyLevel_IsEmpty(t *testing.T) {
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
			el := &EnergyLevel{value: tt.value}
			if got := el.IsEmpty(); got != tt.want {
				t.Errorf("EnergyLevel.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
