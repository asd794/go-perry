package telemetry

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func Init(serviceName string) (func(context.Context) error, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := newStdoutExporter()
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func newStdoutExporter() (sdktrace.SpanExporter, error) {
	return &stdoutExporter{}, nil
}

type stdoutExporter struct{}

func (e *stdoutExporter) ExportSpans(
	ctx context.Context,
	spans []sdktrace.ReadOnlySpan,
) error {
	for _, span := range spans {
		sc := span.SpanContext()
		serviceName := ""
		for _, attr := range span.Resource().Attributes() {
			if attr.Key == "service.name" {
				serviceName = attr.Value.AsString()
				break
			}
		}

		fmt.Fprintf(
			os.Stdout,
			"[TRACE] service=%s span=%s trace_id=%s span_id=%s parent=%s\n",
			serviceName,
			span.Name(),
			sc.TraceID(),
			sc.SpanID(),
			span.Parent().SpanID(),
		)
	}

	return nil
}

func (e *stdoutExporter) Shutdown(ctx context.Context) error {
	return nil
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
