package tracer

import (
	"context"
	"encoding/json"
	"fmt"

	contextInfra "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/context"
	"go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	tracerConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer/config"
)

// Tracer service
type Tracer struct {
	// Name
	name string
	// Configuration
	config *tracerConfig.Config
	// App
	app domain.IApp
	// tracer
	trace.Tracer
	// exporter
	*otlptrace.Exporter
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile tracer configuration file
	configFile = "tracer.yaml"
)

// New creates a new tracer service
func New(app domain.IApp, config *tracerConfig.Config) *Tracer {
	tracer := &Tracer{
		name: "Tracer",
		app:  app,
	}

	if config != nil {
		tracer.config = config
	}

	return tracer
}

// Name gets the service name
func (t *Tracer) Name() string {
	return t.name
}

// Start starts the tracer service
func (t *Tracer) Start() (err error) {
	if t.config == nil {
		t.config = &tracerConfig.Config{}
		t.config.AdditionalConfig = t.additionalConfigType
		if err = config.Load(t.ConfigFile(), t.config); err != nil {
			err = errorCodes.ErrorLoadingConfigFile().Formats(t.ConfigFile(), err)
			message.ErrorMessage(t.Name(), err)
			return err
		}
	}

	if t.config.Enabled {
		exporter, errOtl := otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithEndpoint(t.config.CollectorHostPort),
			otlptracehttp.WithInsecure(),
		)
		if errOtl != nil {
			return errOtl
		}

		// Creating the TracerProvider using a batch span processor to aggregate spans before export
		// and registering it as the global tracer provider
		otel.SetTracerProvider(sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(t.app.Name()),
				attribute.String("service", t.app.Name()),
			)),
			sdkTrace.WithSpanProcessor(sdkTrace.NewBatchSpanProcessor(exporter))),
		)

		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		t.Tracer = otel.Tracer(t.app.Name())
		t.Exporter = exporter
	} else {
		// setting the tracer as a no-op to avoid panics
		t.Tracer = noop.Tracer{}
	}

	t.started = true

	return nil
}

// Started true if started
func (t *Tracer) Started() bool {
	return t.started
}

// Stop stops the tracer service
func (t *Tracer) Stop() error {
	if !t.started {
		return nil
	}

	if t.config.Enabled {
		return t.Exporter.Shutdown(context.Background())
	}

	t.started = false
	return nil
}

// Config gets the configurations
func (t *Tracer) Config() *tracerConfig.Config {
	return t.config
}

// ConfigFile gets the configuration file
func (t *Tracer) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (t *Tracer) WithAdditionalConfigType(obj interface{}) domain.ITracer {
	t.additionalConfigType = obj
	return t
}

// Trace traces data
func (t *Tracer) Trace(ctx context.Context, spanName string, data map[string]any, err error) {
	_, span := t.Tracer.Start(ctx, spanName)
	defer span.End()

	t.setAttributesWithValidation(ctx, span, data, err)
}

// TraceCurrentSpan traces data to a current span
func (t *Tracer) TraceCurrentSpan(ctx context.Context, data map[string]any, err error) {
	t.setAttributesWithValidation(ctx, trace.SpanFromContext(ctx), data, err)
}

// setAttributesWithValidation sets the attributes to a span, with validations
func (t *Tracer) setAttributesWithValidation(ctx context.Context, span trace.Span, data map[string]any, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool(TracerTagError, true))
		span.SetAttributes(attribute.String(TracerTagStatusCode, "ERROR"))
		span.SetAttributes(attribute.String(TracerTagStatusDescription, err.Error()))
	} else {
		span.SetAttributes(attribute.String(TracerTagStatusCode, "OK"))
	}

	for key, value := range data {
		if key == TracerTagParams || key == TracerTagRequestBody || key == TracerTagResponseBody {
			if ctx != nil {
				if ctx.Value(contextInfra.CtxMethod) != nil && ctx.Value(contextInfra.CtxPath) != nil {
					if !t.config.SensitiveUris.Contains(
						fmt.Sprintf("%s", ctx.Value(contextInfra.CtxMethod)),
						fmt.Sprintf("%s", ctx.Value(contextInfra.CtxPath))) {
						span.SetAttributes(getAttribute(key, value))
					}
				}
			}
		} else {
			span.SetAttributes(getAttribute(key, value))
		}
	}
}

// getAttribute gets the specific attribute according to value type
func getAttribute(key string, value any) (attr attribute.KeyValue) {
	if value != nil {
		switch v := value.(type) {
		case string:
			attr = attribute.String(key, v)
		case []string:
			attr = attribute.StringSlice(key, v)
		case bool:
			attr = attribute.Bool(key, v)
		case []bool:
			attr = attribute.BoolSlice(key, v)
		case int:
			attr = attribute.Int(key, v)
		case []int:
			attr = attribute.IntSlice(key, v)
		case int64:
			attr = attribute.Int64(key, v)
		case []int64:
			attr = attribute.Int64Slice(key, v)
		case float64:
			attr = attribute.Float64(key, v)
		case []float64:
			attr = attribute.Float64Slice(key, v)
		default:
			valueB, _ := json.Marshal(v)
			attr = attribute.String(key, string(valueB))
		}
	}

	return attr
}
