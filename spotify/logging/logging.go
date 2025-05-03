package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// LevelDebug is for detailed troubleshooting information
	LevelDebug LogLevel = iota
	// LevelInfo is for general operational information
	LevelInfo
	// LevelWarn is for potentially harmful situations
	LevelWarn
	// LevelError is for error events that might still allow the application to continue running
	LevelError
)

var (
	// DefaultLogger is the default logger instance
	DefaultLogger *Logger

	// Current log level, defaults to Info
	currentLevel = LevelInfo
)

func init() {
	// Initialize the default logger
	DefaultLogger = NewLogger("spotify-provider")

	// Check for log level environment variable
	logLevelEnv := strings.ToLower(os.Getenv("TF_LOG_SPOTIFY"))
	switch logLevelEnv {
	case "debug":
		currentLevel = LevelDebug
	case "info":
		currentLevel = LevelInfo
	case "warn", "warning":
		currentLevel = LevelWarn
	case "error":
		currentLevel = LevelError
	}
}

// Logger provides structured logging capabilities
type Logger struct {
	name string
}

// NewLogger creates a new logger with the given name
func NewLogger(name string) *Logger {
	return &Logger{name: name}
}

// formatMessage formats a log message with timestamp, level, and context
func (l *Logger) formatMessage(level LogLevel, msg string, args map[string]interface{}) string {
	// Format timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Format level
	var levelStr string
	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
	case LevelInfo:
		levelStr = "INFO"
	case LevelWarn:
		levelStr = "WARN"
	case LevelError:
		levelStr = "ERROR"
	}

	// Format base message
	base := fmt.Sprintf("%s [%s] %s: %s", timestamp, levelStr, l.name, msg)

	// Add context if available
	if len(args) > 0 {
		contextStrings := make([]string, 0, len(args))
		for k, v := range args {
			contextStrings = append(contextStrings, fmt.Sprintf("%s=%v", k, v))
		}
		base = fmt.Sprintf("%s {%s}", base, strings.Join(contextStrings, ", "))
	}

	return base
}

// log logs a message at the specified level
func (l *Logger) log(level LogLevel, msg string, args map[string]interface{}) {
	// Skip logging if the level is below the current level
	if level < currentLevel {
		return
	}

	log.Println(l.formatMessage(level, msg, args))
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, msg, argsToMap(args...))
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, argsToMap(args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, msg, argsToMap(args...))
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, argsToMap(args...))
}

// WithContext returns a new logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// In a more advanced implementation, we could extract values from the context
	// For now, just return the same logger
	return l
}

// Helper functions for the default logger

// Debug logs a debug message using the default logger
func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger
func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

// argsToMap converts variadic arguments to a map
// Arguments should be provided in key-value pairs
func argsToMap(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Return empty map if no args
	if len(args) == 0 {
		return result
	}

	// If a single argument is provided and it's already a map, use it
	if len(args) == 1 {
		if m, ok := args[0].(map[string]interface{}); ok {
			return m
		}
	}

	// Process key-value pairs
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if ok {
				result[key] = args[i+1]
			}
		} else {
			// Handle odd number of arguments
			key, ok := args[i].(string)
			if ok {
				result[key] = "<missing value>"
			}
		}
	}

	return result
}