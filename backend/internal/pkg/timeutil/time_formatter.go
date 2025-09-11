package timeutil

import (
	"errors"
	"time"
)

const (
	// RFC3339Format is the standard format used across the application
	RFC3339Format = time.RFC3339
	// DateFormat is the date-only format (YYYY-MM-DD)
	DateFormat = "2006-01-02"
)

// FormatTime formats a time.Time to RFC3339 string
func FormatTime(t time.Time) string {
	return t.Format(RFC3339Format)
}

// FormatTimePointer formats a *time.Time to *string in RFC3339 format
// Returns nil if the input is nil
func FormatTimePointer(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(RFC3339Format)
	return &formatted
}

// ParseTime parses an RFC3339 formatted string to time.Time
func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, errors.New("time string cannot be empty")
	}

	parsedTime, err := time.Parse(RFC3339Format, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid time format. Use RFC3339 format (e.g., 2023-12-01T10:00:00Z)")
	}

	return parsedTime, nil
}

// ParseTimePointer parses an RFC3339 formatted string to *time.Time
// Returns nil if the input string is empty or nil
func ParseTimePointer(timeStr *string) (*time.Time, error) {
	if timeStr == nil || *timeStr == "" {
		return nil, nil
	}

	parsedTime, err := ParseTime(*timeStr)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

// Now returns current time formatted as RFC3339 string
func Now() string {
	return FormatTime(time.Now())
}

// IsValidTimeFormat checks if a string is in valid RFC3339 format
func IsValidTimeFormat(timeStr string) bool {
	_, err := time.Parse(RFC3339Format, timeStr)
	return err == nil
}

// ParseDate parses a date string in YYYY-MM-DD format to time.Time
// The time component is set to 00:00:00 UTC
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date string cannot be empty")
	}

	parsedTime, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format. Use YYYY-MM-DD format (e.g., 2023-12-01)")
	}

	// Convert to UTC to be consistent with other time operations
	return parsedTime.UTC(), nil
}

// ParseDatePointer parses a date string in YYYY-MM-DD format to *time.Time
// Returns nil if the input string is empty or nil
func ParseDatePointer(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	parsedTime, err := ParseDate(*dateStr)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

// FormatDate formats a time.Time to YYYY-MM-DD string
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// IsValidDateFormat checks if a string is in valid YYYY-MM-DD format
func IsValidDateFormat(dateStr string) bool {
	_, err := time.Parse(DateFormat, dateStr)
	return err == nil
}
