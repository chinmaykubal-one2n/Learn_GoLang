package metrics

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
)

var (
	// Global metrics
	APIRequestCount    metric.Int64Counter
	APIRequestDuration metric.Float64Histogram
	meter              metric.Meter
)

// InitializeMetrics sets up OpenTelemetry metrics
func InitializeMetrics(ctx context.Context) (func(context.Context) error, error) {
	var (
		serviceName  = os.Getenv("SERVICE_NAME")
		collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		insecure     = os.Getenv("INSECURE_MODE")
	)

	// Configure metric exporter
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(collectorURL),
	}

	if len(insecure) > 0 {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	} else {
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: %v", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create meter
	meter = otel.Meter("student-api-metrics")

	// Initialize metrics
	if err := initMetrics(); err != nil {
		return nil, err
	}

	return meterProvider.Shutdown, nil
}

// initMetrics creates the metric instruments
func initMetrics() error {
	var err error

	// Counter for API requests
	APIRequestCount, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	// Histogram for request duration
	APIRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

// RecordAPICall records both the count and duration of an API call
func RecordAPICall(ctx context.Context, method, path string, statusCode int, duration float64) {
	// Common attributes
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("http.status_code", statusCode),
	}

	// Record request count
	APIRequestCount.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Record request duration
	APIRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}
