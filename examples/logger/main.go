package main

import (
	"github.com/rk-the-dev/golib-core/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.InitializeLogger("debug", "app.log", 5, 3, 7)

	// Log messages with structured fields
	logger.Info("Service started", logrus.Fields{"module": "main", "version": "1.0.0"})
	logger.Debug("Debugging mode enabled", logrus.Fields{"env": "development"})
	logger.Warn("This is a warning message", logrus.Fields{"retry_attempts": 3})
	logger.Error("An error occurred", logrus.Fields{"error": "database connection failed"})

	// Simulating a fatal error (this will terminate the application)
	// logger.Fatal("Critical failure", logrus.Fields{"reason": "out of memory"})
}
