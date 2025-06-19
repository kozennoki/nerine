package logger_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/logger"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger, err := logger.New()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Verify it's actually a zap.Logger
	if _, ok := interface{}(logger).(*zap.Logger); !ok {
		t.Error("Expected logger to be of type *zap.Logger")
	}

	// Test that we can use the logger without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger panicked when used: %v", r)
		}
	}()

	logger.Info("Test log message")
	logger.Sync() // Flush any buffered log entries
}

func TestNew_ProducesProductionLogger(t *testing.T) {
	t.Parallel()

	logger, err := logger.New()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Test that the logger is configured for production
	// by checking it doesn't panic and can handle structured logging
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Production logger panicked: %v", r)
		}
	}()

	logger.Info("production test",
		zap.String("key", "value"),
		zap.Int("number", 1),
	)
	logger.Sync()
}
