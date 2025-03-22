package logger

import (
	"context"
	"testing"
)

func TestContextWithLogger(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger()
	ctxWithLogger := ContextWithLogger(ctx, logger)

	// Extract logger from new context
	extractedLogger, ok := ctxWithLogger.Value(loggerCtx{}).(*Logger)
	if !ok || extractedLogger != logger {
		t.Errorf("Logger was not properly added to context")
	}
}

func TestLoggerFromContext_WithLogger(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger()
	ctxWithLogger := ContextWithLogger(ctx, logger)

	// Extract logger using loggerFromContext
	extractedLogger := fromContext(ctxWithLogger)
	if extractedLogger != logger {
		t.Errorf("Did not retrieve correct logger from context")
	}
}

func TestLoggerFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	// Extract logger from a context without a logger
	extractedLogger := fromContext(ctx)
	if extractedLogger != Default() {
		t.Errorf("Did not retrieve default logger when context has no logger")
	}
}