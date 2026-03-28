package tracing

import (
	"context"

	kratos "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "intelligent-guidance-system"
)

func TracingMiddleware() kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			tracer := otel.Tracer(tracerName)
			operation := tr.Operation()
			if operation == "" {
				operation = "unknown"
			}

			propagator := otel.GetTextMapPropagator()
			ctx = propagator.Extract(ctx, &headerCarrier{header: tr.RequestHeader()})

			ctx, span := tracer.Start(ctx, operation,
				trace.WithAttributes(
					attribute.String("transport.kind", tr.Kind().String()),
				),
			)
			defer span.End()

			reply, err := handler(ctx, req)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "success")
			}

			return reply, err
		}
	}
}

func TracingMiddlewareWithName(name string) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			tracer := otel.Tracer(name)
			operation := tr.Operation()
			if operation == "" {
				operation = "unknown"
			}

			propagator := otel.GetTextMapPropagator()
			ctx = propagator.Extract(ctx, &headerCarrier{header: tr.RequestHeader()})

			ctx, span := tracer.Start(ctx, operation,
				trace.WithAttributes(
					attribute.String("service", name),
					attribute.String("transport.kind", tr.Kind().String()),
				),
			)
			defer span.End()

			reply, err := handler(ctx, req)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "success")
			}

			return reply, err
		}
	}
}

type headerCarrier struct {
	header transport.RequestHeader
}

func (h *headerCarrier) Get(key string) string {
	return h.header.Get(key)
}

func (h *headerCarrier) Set(key, value string) {
	h.header.Set(key, value)
}

func (h *headerCarrier) Keys() []string {
	return nil
}

func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

func InjectTraceIntoHeader(ctx context.Context, header transport.RequestHeader) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, &headerCarrier{header: header})
}