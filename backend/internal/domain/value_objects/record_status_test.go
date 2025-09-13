package value_objects

import (
	"testing"
)

func TestNewMentalHealthRecordStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid status open",
			input:       "public",
			wantValue:   "public",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid status closed",
			input:       "private",
			wantValue:   "private",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty status",
			input:       "",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid mental health record status: ",
		},
		{
			name:        "invalid status",
			input:       "invalid",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid mental health record status: invalid",
		},
		{
			name:        "status with spaces",
			input:       " active ",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "invalid mental health record status: active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMentalHealthRecordStatus(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMentalHealthRecordStatus() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewMentalHealthRecordStatus() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewMentalHealthRecordStatus() unexpected error = %v", err)
				return
			}

			// Check value
			if got.String() != tt.wantValue {
				t.Errorf("NewMentalHealthRecordStatus() = %v, want %v", got.String(), tt.wantValue)
			}
		})
	}
}

func TestMentalHealthRecordStatus_String(t *testing.T) {
	tests := []struct {
		name  string
		value MentalHealthRecordStatus
		want  string
	}{
		{
			name:  "public status",
			value: RecordStatusPublic,
			want:  "public",
		},
		{
			name:  "private status",
			value: RecordStatusPrivate,
			want:  "private",
		},
		{
			name:  "empty status",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("MentalHealthRecordStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
