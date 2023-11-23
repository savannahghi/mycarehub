package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/cmd"
)

const (
	waitSeconds = 30

	JaegerCollectorEndpoint = "JAEGER_COLLECTOR_ENDPOINT"
)

func startTracing() (*trace.TracerProvider, error) {
	serviceName := fmt.Sprintf("mycarehub-%v", serverutils.GetRunningEnvironment())
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(serverutils.MustGetEnvVar(JaegerCollectorEndpoint)),
			otlptracehttp.WithHeaders(headers),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
		trace.WithSpanProcessor(
			sentryotel.NewSentrySpanProcessor(),
		),
	)

	otel.SetTracerProvider(tracerprovider)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	return tracerprovider, nil
}

func main() {
	//  Run command line arguments
	cmd.Execute()

	ctx := context.Background()
	err := serverutils.Sentry()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	traceProvider, err := startTracing()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			serverutils.LogStartupError(ctx, err)
		}
	}()

	_ = traceProvider.Tracer("mycarehub")

	port, err := strconv.Atoi(serverutils.MustGetEnvVar(serverutils.PortEnvVarName))
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	srv := presentation.PrepareServer(ctx, port, presentation.AllowedOrigins)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverutils.LogStartupError(ctx, err)
		}
	}()

	// Block until we receive a sigint (CTRL+C) signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*waitSeconds)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until timeout
	err = srv.Shutdown(ctx)
	log.Printf("graceful shutdown started; the timeout is %d secs", waitSeconds)
	if err != nil {
		log.Printf("error during clean shutdown: %s", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
