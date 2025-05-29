package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	// "go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
)

var Logger *otelzap.Logger

func InitLogger(ctx context.Context, serviceName string, otlpEndpoint string) error {
	// Step 1: Create an OTLP HTTP exporter
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Step 2: Create a batching processor
	processor := sdklog.NewBatchProcessor(exporter)

	// Step 3: Create the logger provider
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
	)

	// Step 4: Graceful shutdown
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "Error shutting down logger provider:", err)
		}
	}()

	// Step 5: Create a zap logger
	zapLogger := zap.NewExample()
	otelLogger := otelzap.New(zapLogger,
		otelzap.WithLoggerProvider(provider),
		otelzap.WithMinLevel(zap.DebugLevel),
	)

	// Step 6: Assign the logger
	Logger = otelLogger

	return nil
}
