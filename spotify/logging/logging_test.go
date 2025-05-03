package logging

import (
	"context"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	// DefaultLogger is a variable, not a function
	if DefaultLogger == nil {
		t.Error("Expected non-nil DefaultLogger, got nil")
	}
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test-logger")
	
	if logger == nil {
		t.Error("Expected non-nil logger, got nil")
	}
	
	if logger.name != "test-logger" {
		t.Errorf("Expected logger name 'test-logger', got '%s'", logger.name)
	}
}

func TestLoggerWithContext(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger("test-logger")
	
	// Use WithContext method instead
	loggerWithCtx := logger.WithContext(ctx)
	
	if loggerWithCtx == nil {
		t.Error("Expected non-nil logger from context, got nil")
	}
}

func TestLoggerLevels(t *testing.T) {
	// This is a simple test to ensure the log level constants exist
	// LogLevel is an int type, not string
	levels := []LogLevel{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
	}
	
	// Just check that the levels are different
	for i := 1; i < len(levels); i++ {
		if levels[i] <= levels[i-1] {
			t.Errorf("Expected increasing log levels, but level %v is not greater than %v", levels[i], levels[i-1])
		}
	}
}