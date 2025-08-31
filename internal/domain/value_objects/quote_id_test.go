package value_objects

import (
	"testing"
)

func TestNewQuoteID(t *testing.T) {
	// Test that NewQuoteID creates a quote ID with value 0
	quoteID := NewQuoteID()

	// Check that it's not nil
	if quoteID == nil {
		t.Error("NewQuoteID() returned nil")
	}

	// Check that it has value 0
	if quoteID.Value() != 0 {
		t.Errorf("NewQuoteID() = %v, want 0", quoteID.Value())
	}

	// Check string representation
	if quoteID.String() != "0" {
		t.Errorf("NewQuoteID().String() = %v, want '0'", quoteID.String())
	}
}

func TestNewQuoteIDFromInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "positive integer",
			input: 123,
			want:  123,
		},
		{
			name:  "zero",
			input: 0,
			want:  0,
		},
		{
			name:  "negative integer",
			input: -1,
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQuoteIDFromInt(tt.input)
			if got.Value() != tt.want {
				t.Errorf("NewQuoteIDFromInt() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func TestNewQuoteIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   int
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid integer string",
			input:       "123",
			wantValue:   123,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "zero string",
			input:       "0",
			wantValue:   0,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "negative integer string",
			input:       "-1",
			wantValue:   -1,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "empty string",
			input:       "",
			wantValue:   0,
			wantErr:     true,
			expectedErr: "invalid quote id: strconv.Atoi: parsing \"\": invalid syntax",
		},
		{
			name:        "invalid string",
			input:       "not-a-number",
			wantValue:   0,
			wantErr:     true,
			expectedErr: "invalid quote id: strconv.Atoi: parsing \"not-a-number\": invalid syntax",
		},
		{
			name:        "float string",
			input:       "123.45",
			wantValue:   0,
			wantErr:     true,
			expectedErr: "invalid quote id: strconv.Atoi: parsing \"123.45\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuoteIDFromString(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewQuoteIDFromString() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewQuoteIDFromString() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewQuoteIDFromString() unexpected error = %v", err)
				return
			}

			// Check value
			if got.Value() != tt.wantValue {
				t.Errorf("NewQuoteIDFromString() = %v, want %v", got.Value(), tt.wantValue)
			}
		})
	}
}

func TestQuoteID_Value(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "positive integer",
			value: 123,
			want:  123,
		},
		{
			name:  "zero",
			value: 0,
			want:  0,
		},
		{
			name:  "negative integer",
			value: -1,
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QuoteID{value: tt.value}
			if got := q.Value(); got != tt.want {
				t.Errorf("QuoteID.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteID_String(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{
			name:  "positive integer",
			value: 123,
			want:  "123",
		},
		{
			name:  "zero",
			value: 0,
			want:  "0",
		},
		{
			name:  "negative integer",
			value: -1,
			want:  "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QuoteID{value: tt.value}
			if got := q.String(); got != tt.want {
				t.Errorf("QuoteID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
