package tracer

import (
	"context"

	"github.com/gocraft/dbr/v2"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer"
)

//const tracerTagTraver = tracer.TracerTagTracer

// tracerMiddleware is an EventReceiver that traces queries
type tracerMiddleware struct {
	app domain.IApp
	dbr.NullEventReceiver
	dbr.TracingEventReceiver
}

type TraceValue struct {
	attrs map[string]any
	err   error
}

// NewTracerMiddleware creates a new tracerMiddleware
func NewTracerMiddleware(app domain.IApp) dbr.EventReceiver {
	return &tracerMiddleware{
		app: app,
	}
}

// SpanStart starts the span
func (t *tracerMiddleware) SpanStart(ctx context.Context, eventName string, query string) context.Context {
	attrs := &TraceValue{
		attrs: map[string]any{
			tracer.TracerTagEventName: eventName,
			tracer.TracerTagQuery:     query,
		},
	}
	return context.WithValue(ctx, tracer.TracerTagTracer, attrs)
}

// SpanError collects the error
func (t *tracerMiddleware) SpanError(ctx context.Context, err error) {
	value, _ := ctx.Value(tracer.TracerTagTracer).(*TraceValue)
	if value != nil {
		value.err = err
	}
}

// SpanFinish closes the span
func (t *tracerMiddleware) SpanFinish(ctx context.Context) {
	value := ctx.Value(tracer.TracerTagTracer).(*TraceValue)
	if value != nil {
		t.app.Tracer().Trace(ctx, t.app.Database().Name(), value.attrs, value.err)
	}
}
