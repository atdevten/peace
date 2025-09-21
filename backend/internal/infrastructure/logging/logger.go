package logging

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Level      string
	Format     string
	TimeFormat string
	Caller     bool
	CallerSkip int
	FilePath   string
}

type Logger struct {
	logger zerolog.Logger
	config *Config
	File   *os.File // Store file handle for cleanup
}

func NewLogger(config *Config) *Logger {
	// Set log level
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = config.TimeFormat

	// Create output writer
	var output io.Writer
	var file *os.File
	if config.FilePath != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(err)
		}

		// Create log file
		file, err = os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		output = file
	} else {
		output = os.Stderr
	}

	// Create logger
	var logger zerolog.Logger
	if config.Format == "pretty" {
		// Pretty format for development
		logger = log.Output(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: config.TimeFormat,
		})
	} else {
		// JSON format for production
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	// Add caller information if enabled
	if config.Caller {
		logger = logger.With().CallerWithSkipFrameCount(config.CallerSkip).Logger()
	}

	return &Logger{
		logger: logger,
		config: config,
		File:   file,
	}
}

// GetLogger returns a logger with trace context
func (l *Logger) GetLogger(ctx context.Context) *zerolog.Logger {
	logger := l.logger

	// Add trace context if available
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		logger = logger.With().
			Str("trace_id", spanCtx.TraceID().String()).
			Str("span_id", spanCtx.SpanID().String()).
			Logger()
	}

	return &logger
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, msg string) {
	l.GetLogger(ctx).Info().Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, err error, msg string) {
	l.GetLogger(ctx).Error().Err(err).Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, msg string) {
	l.GetLogger(ctx).Debug().Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, msg string) {
	l.GetLogger(ctx).Warn().Msg(msg)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(ctx context.Context, fields map[string]interface{}) *zerolog.Logger {
	logger := *l.GetLogger(ctx)
	for key, value := range fields {
		logger = logger.With().Interface(key, value).Logger()
	}
	return &logger
}

// Close closes the file handle if it was opened for file logging
func (l *Logger) Close() error {
	if l.File != nil {
		err := l.File.Close()
		l.File = nil // Set to nil to prevent double close
		return err
	}
	return nil
}
