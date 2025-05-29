package logging

import (
	"context"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	// "go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
)

var Logger *otelzap.Logger

func InitLogger(ctx context.Context, serviceName string, otlpEndpoint string) (func(context.Context) error, error) {
	// Step 1: Create an OTLP HTTP exporter
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Step 2: Create a batching processor
	processor := sdklog.NewBatchProcessor(exporter)

	// Step 3: Create the logger provider
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
	)

	// Step 4: Create a zap logger
	zapLogger := zap.NewExample()
	otelLogger := otelzap.New(zapLogger,
		otelzap.WithLoggerProvider(provider),
		otelzap.WithMinLevel(zap.DebugLevel),
	)

	// Step 5: Assign the logger
	Logger = otelLogger

	return provider.Shutdown, nil
}
