package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"svc/internal/feature/ioops"
	"svc/internal/repository/sqlcore"
	"svc/internal/service/core"
	"svc/internal/z"
	"svc/internal/z/zsql"
	"time"

	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = telemetry(ctx)
	// -----------------------------------------------------------------------------------------------------------------
	zsqlConn := zsql.RoundRobin(
		must(zsql.Open("")),
		must(zsql.Open("")),
		must(zsql.Open("")),
	)
	// -----------------------------------------------------------------------------------------------------------------
	sqlcoreRepository := must(sqlcore.New(ctx,
		sqlcore.Configuration{},
		sqlcore.Dependency{
			Conn: zsqlConn,
		},
	))
	// -----------------------------------------------------------------------------------------------------------------
	ioopsFeature := must(ioops.New(ctx,
		ioops.Configuration{},
		ioops.Dependency{
			SQL_Core: sqlcoreRepository,
		},
	))
	// -----------------------------------------------------------------------------------------------------------------
	coreService := must(core.New(ctx,
		core.Configuration{
			ServingHTTP: true,
		},
		core.Dependency{
			IO_Ops: ioopsFeature,
		},
	))
	// -----------------------------------------------------------------------------------------------------------------
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		coreService.Stop(ctx)
	}()
	os.Exit(must(0, coreService.Start(ctx)))
}

func must[T any](v T, err error) T {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err.Error())
		os.Exit(1)
	} else {
		fmt.Fprintln(os.Stdout, "ok")
	}
	return v
}

func telemetry(ctx context.Context) context.Context {
	ic := z.Context
	logFile := must(os.Create("./var/log/service-go.log"))

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	res := must(resource.New(ctx, resource.WithAttributes(
		semconv.ServiceInstanceID("PID/"+strconv.Itoa(os.Getpid())),
		semconv.ServiceNamespace("service"),
		semconv.ServiceName("core"),
		semconv.ServiceVersion("v0.0.0"),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersion(otel.Version()),
	)))

	// -----------------------------------------------------------------------------------------------------------------
	// *log.Logger
	mw := io.MultiWriter(os.Stdout, logFile)
	log := log.New(mw, "", log.LstdFlags)
	ctx = ic.PutLogLogger(ctx, log)

	// -----------------------------------------------------------------------------------------------------------------
	// *zerolog.Logger
	mw = zerolog.MultiLevelWriter(zerolog.NewConsoleWriter(), logFile)
	zlog := zerolog.New(mw).With().Timestamp().Logger()
	ctx = ic.PutZerologLogger(ctx, &zlog)
	otel.SetLogger(zerologr.New(&zlog))

	// -----------------------------------------------------------------------------------------------------------------
	// Otel Meter
	otel.SetMeterProvider(metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(
			must(otlpmetricgrpc.New(ctx,
				otlpmetricgrpc.WithInsecure(),
				otlpmetricgrpc.WithEndpoint("otelcol:4317"),
			)),
			metric.WithInterval(60*time.Second), // default 60s
			metric.WithTimeout(30*time.Second),  // default 30s
		)),
	))

	// -----------------------------------------------------------------------------------------------------------------
	// Otel Tracer
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(
			must(otlptracegrpc.New(ctx,
				otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint("otelcol:4317"),
			)),
			trace.WithBatchTimeout(5*time.Second),
			trace.WithExportTimeout(30*time.Second),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithMaxQueueSize(trace.DefaultMaxQueueSize),
		),
	))

	return ctx
}
