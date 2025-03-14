package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the global log instance
var Logger *logrus.Logger

// InitializeLogger initializes the logger with configurations
func InitializeLogger(logLevel string, logFile string, maxSizeMB int, maxBackups int, maxAgeDays int) {
	Logger = logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // Default to INFO
	}
	Logger.SetLevel(level)

	// Set log format as JSON with timestamp
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Configure log output
	if logFile != "" {
		Logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    maxSizeMB,  // Max size in MB before rotation
			MaxBackups: maxBackups, // Max number of old log files to retain
			MaxAge:     maxAgeDays, // Max age in days before deletion
			Compress:   true,       // Compress old logs
		})
	} else {
		Logger.SetOutput(os.Stdout)
	}
}

// Debug logs a debug message
func Debug(msg string, fields logrus.Fields) {
	Logger.WithFields(fields).Debug(msg)
}

// Info logs an info message
func Info(msg string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(msg)
}

// Warn logs a warning message
func Warn(msg string, fields logrus.Fields) {
	Logger.WithFields(fields).Warn(msg)
}

// Error logs an error message
func Error(msg string, fields logrus.Fields) {
	Logger.WithFields(fields).Error(msg)
}

// Fatal logs a fatal error and exits
func Fatal(msg string, fields logrus.Fields) {
	Logger.WithFields(fields).Fatal(msg)
}
