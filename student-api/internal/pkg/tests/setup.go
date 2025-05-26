package tests

import (
	logging "student-api/internal/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupLoggerForTests initializes a zap logger for use in tests
func SetupLoggerForTests() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel) // Only log errors during tests
	logger, _ := config.Build()
	logging.Logger = logger
}
