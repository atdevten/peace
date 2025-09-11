package timeutil

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	// Test with a known time
	testTime := time.Date(2023, 12, 1, 10, 30, 45, 0, time.UTC)
	expected := "2023-12-01T10:30:45Z"

	result := FormatTime(testTime)
	if result != expected {
		t.Errorf("FormatTime() = %v, want %v", result, expected)
	}
}

func TestFormatTimePointer(t *testing.T) {
	// Test with nil pointer
	result := FormatTimePointer(nil)
	if result != nil {
		t.Errorf("FormatTimePointer(nil) = %v, want nil", result)
	}

	// Test with valid time pointer
	testTime := time.Date(2023, 12, 1, 10, 30, 45, 0, time.UTC)
	expected := "2023-12-01T10:30:45Z"

	result = FormatTimePointer(&testTime)
	if result == nil {
		t.Error("FormatTimePointer() returned nil for valid time")
	} else if *result != expected {
		t.Errorf("FormatTimePointer() = %v, want %v", *result, expected)
	}
}

func TestParseTime(t *testing.T) {
	// Test valid time string
	timeStr := "2023-12-01T10:30:45Z"
	expected := time.Date(2023, 12, 1, 10, 30, 45, 0, time.UTC)

	result, err := ParseTime(timeStr)
	if err != nil {
		t.Errorf("ParseTime() error = %v", err)
	}
	if !result.Equal(expected) {
		t.Errorf("ParseTime() = %v, want %v", result, expected)
	}

	// Test empty string
	_, err = ParseTime("")
	if err == nil {
		t.Error("ParseTime(\"\") should return error")
	}

	// Test invalid format
	_, err = ParseTime("invalid-time")
	if err == nil {
		t.Error("ParseTime(\"invalid-time\") should return error")
	}
}

func TestParseTimePointer(t *testing.T) {
	// Test nil pointer
	result, err := ParseTimePointer(nil)
	if err != nil {
		t.Errorf("ParseTimePointer(nil) error = %v", err)
	}
	if result != nil {
		t.Errorf("ParseTimePointer(nil) = %v, want nil", result)
	}

	// Test empty string pointer
	emptyStr := ""
	result, err = ParseTimePointer(&emptyStr)
	if err != nil {
		t.Errorf("ParseTimePointer(&\"\") error = %v", err)
	}
	if result != nil {
		t.Errorf("ParseTimePointer(&\"\") = %v, want nil", result)
	}

	// Test valid time string pointer
	timeStr := "2023-12-01T10:30:45Z"
	expected := time.Date(2023, 12, 1, 10, 30, 45, 0, time.UTC)

	result, err = ParseTimePointer(&timeStr)
	if err != nil {
		t.Errorf("ParseTimePointer() error = %v", err)
	}
	if result == nil {
		t.Error("ParseTimePointer() returned nil for valid time")
	} else if !result.Equal(expected) {
		t.Errorf("ParseTimePointer() = %v, want %v", *result, expected)
	}
}

func TestIsValidTimeFormat(t *testing.T) {
	// Test valid format
	if !IsValidTimeFormat("2023-12-01T10:30:45Z") {
		t.Error("IsValidTimeFormat() should return true for valid format")
	}

	// Test invalid format
	if IsValidTimeFormat("invalid-time") {
		t.Error("IsValidTimeFormat() should return false for invalid format")
	}

	// Test empty string
	if IsValidTimeFormat("") {
		t.Error("IsValidTimeFormat() should return false for empty string")
	}
}

func TestNow(t *testing.T) {
	// Test that Now() returns a valid RFC3339 formatted string
	result := Now()
	if !IsValidTimeFormat(result) {
		t.Errorf("Now() returned invalid format: %v", result)
	}
}
