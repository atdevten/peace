package logging

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLogger_Close(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	// Test case 1: Logger with file output should close file handle
	t.Run("file logger closes handle", func(t *testing.T) {
		config := &Config{
			Level:      "info",
			Format:     "json",
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			Caller:     false,
			CallerSkip: 2,
			FilePath:   logFile,
		}

		logger := NewLogger(config)
		if logger.File == nil {
			t.Fatal("Expected file handle to be created")
		}

		// Verify file exists and is writable
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Fatal("Log file should exist after logger creation")
		}

		// Test logging works
		ctx := context.Background()
		logger.Info(ctx, "test message")

		// Close logger
		if err := logger.Close(); err != nil {
			t.Fatalf("Expected no error when closing logger, got: %v", err)
		}

		// Verify file handle is closed by trying to write to it
		// This should fail if the file handle is properly closed
		_, err := logger.File.WriteString("test")
		if err == nil {
			t.Error("Expected error when writing to closed file handle")
		}
	})

	// Test case 2: Logger without file output should not panic
	t.Run("console logger handles close gracefully", func(t *testing.T) {
		config := &Config{
			Level:      "info",
			Format:     "pretty",
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			Caller:     false,
			CallerSkip: 2,
			FilePath:   "", // No file path = console output
		}

		logger := NewLogger(config)
		if logger.File != nil {
			t.Fatal("Expected no file handle for console logger")
		}

		// Close should not panic or return error
		if err := logger.Close(); err != nil {
			t.Fatalf("Expected no error when closing console logger, got: %v", err)
		}
	})

	// Test case 3: Multiple close calls should be safe
	t.Run("multiple close calls are safe", func(t *testing.T) {
		config := &Config{
			Level:      "info",
			Format:     "json",
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			Caller:     false,
			CallerSkip: 2,
			FilePath:   filepath.Join(tempDir, "test2.log"),
		}

		logger := NewLogger(config)

		// First close
		if err := logger.Close(); err != nil {
			t.Fatalf("First close failed: %v", err)
		}

		// Second close should not panic
		if err := logger.Close(); err != nil {
			t.Fatalf("Second close failed: %v", err)
		}
	})
}

func TestLogger_FileDescriptorLeak(t *testing.T) {
	// This test verifies that file descriptors are properly released
	tempDir := t.TempDir()

	// Create multiple loggers to test for FD leaks
	for i := 0; i < 10; i++ {
		logFile := filepath.Join(tempDir, "test", "log", "file", "path", "that", "does", "not", "exist", "yet", "test.log")
		config := &Config{
			Level:      "info",
			Format:     "json",
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			Caller:     false,
			CallerSkip: 2,
			FilePath:   logFile,
		}

		logger := NewLogger(config)

		// Use the logger
		ctx := context.Background()
		logger.Info(ctx, "test message")
		logger.Error(ctx, nil, "test error")

		// Close the logger
		if err := logger.Close(); err != nil {
			t.Fatalf("Failed to close logger %d: %v", i, err)
		}
	}

	// If we get here without "too many open files" error, the test passes
	t.Log("Successfully created and closed 10 loggers without FD leak")
}
