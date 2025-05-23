package logging

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(ctx context.Context, serviceName string, otlpEndpoint string) error {
	// Create an OTLP HTTP exporter to send logs to SigNoz
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithInsecure(), // Disable TLS if SigNoz is running locally
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Create a processor for batching log records
	processor := sdklog.NewBatchProcessor(exporter)

	// Create a LoggerProvider
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
	)

	// Set up zap with the otelzap core bridge
	core := otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(provider))
	zapLogger := zap.New(core)

	// Assign to global variable
	Logger = zapLogger

	// Clean up the provider on application shutdown
	go func() {
		<-ctx.Done()
		if err := provider.Shutdown(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, "Error shutting down logger provider:", err)
		}
	}()

	return nil
}
