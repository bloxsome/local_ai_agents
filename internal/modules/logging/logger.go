// Package logging provides logging functionality for the Ollama Agents application
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// DEBUG level for detailed debugging information
	DEBUG LogLevel = iota
	// INFO level for general information
	INFO
	// WARNING level for warning messages
	WARNING
	// ERROR level for error messages
	ERROR
	// CRITICAL level for critical error messages
	CRITICAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}[l]
}

// ParseLogLevel converts a string to a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "CRITICAL":
		return CRITICAL
	default:
		return WARNING // Default to WARNING
	}
}

// Logger is a custom logger for the Ollama Agents application
type Logger struct {
	level  LogLevel
	logger *log.Logger
	file   *os.File
}

// NewLogger creates a new Logger instance
func NewLogger(logFile string, level string) (*Logger, error) {
	// Ensure the log file directory exists
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open the log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a multi-writer to write to both file and stderr
	multiWriter := io.MultiWriter(file, os.Stderr)

	// Create the logger
	logger := log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		level:  ParseLogLevel(level),
		logger: logger,
		file:   file,
	}, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log logs a message at the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level >= l.level {
		msg := fmt.Sprintf(format, args...)
		l.logger.Printf("[%s] %s", level.String(), msg)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Critical logs a critical message
func (l *Logger) Critical(format string, args ...interface{}) {
	l.log(CRITICAL, format, args...)
}

// Exception logs an exception with stack trace
func (l *Logger) Exception(format string, err error, args ...interface{}) {
	if l.level <= ERROR {
		msg := fmt.Sprintf(format, args...)
		l.logger.Printf("[ERROR] %s: %+v", msg, err)
	}
}

// Global logger instance
var logger *Logger

// SetupLogging initializes the global logger
func SetupLogging(logFile, logLevel string) error {
	var err error
	logger, err = NewLogger(logFile, logLevel)
	if err != nil {
		return err
	}
	fmt.Printf("Logging setup complete. Log file: %s\n", logFile)
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if logger == nil {
		// If logger is not initialized, create a default logger to stderr
		logger = &Logger{
			level:  WARNING,
			logger: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile),
		}
	}
	return logger
}
