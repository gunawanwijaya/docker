package z

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func OTel(ctx context.Context, logName, traceName, metricName string) (log *zerolog.Logger, trc trace.Tracer, mtr metric.Meter) {
	log = Context.ZerologLogger(ctx)
	trc = otel.Tracer(traceName)
	mtr = otel.Meter(metricName)
	return
}

var Context interface {
	PutLogLogger(ctx context.Context, logger *log.Logger) context.Context
	LogLogger(ctx context.Context) *log.Logger
	PutZerologLogger(ctx context.Context, logger *zerolog.Logger) context.Context
	ZerologLogger(ctx context.Context) *zerolog.Logger
} = struct {
	ctxKeyLogLogger
	ctxKeyZerologLogger
}{}

// ------------------------------------------------------------------------

type ctxKeyLogLogger struct{}

var theLogLogger *log.Logger

func (key ctxKeyLogLogger) PutLogLogger(ctx context.Context, val *log.Logger) context.Context {
	if val == nil {
		panic("ctxKeyLogLogger")
	} else if theLogLogger == nil {
		theLogLogger = val
	}
	return context.WithValue(ctx, key, val)
}
func (key ctxKeyLogLogger) LogLogger(ctx context.Context) *log.Logger {
	v, ok := ctx.Value(key).(*log.Logger)
	if !ok {
		v = theLogLogger
	}
	return v
}

// ------------------------------------------------------------------------

type ctxKeyZerologLogger struct{}

var theZerologLogger *zerolog.Logger

func (key ctxKeyZerologLogger) PutZerologLogger(ctx context.Context, val *zerolog.Logger) context.Context {
	if val == nil {
		panic("ctxKeyZerologLogger")
	} else if theZerologLogger == nil {
		theZerologLogger = val
	}
	return context.WithValue(ctx, key, val)
}
func (key ctxKeyZerologLogger) ZerologLogger(ctx context.Context) *zerolog.Logger {
	v, ok := ctx.Value(key).(*zerolog.Logger)
	if !ok {
		v = theZerologLogger
	}
	if v == nil {
		zn := zerolog.Nop()
		v = &zn
	}
	return v
}

// ------------------------------------------------------------------------

func ZerologSpanContext(sc trace.SpanContext) zerolog.LogObjectMarshaler {
	return zerologSpanContext{sc}
}

type zerologSpanContext struct{ trace.SpanContext }

func (x zerologSpanContext) MarshalZerologObject(e *zerolog.Event) {
	e = e.
		Stringer("_traceID", x.TraceID())
}
