package tracer

import (
	"context"

	tracerConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer/config"
	"github.com/stretchr/testify/mock"
)

func NewTracerMock() *TracerMock {
	return &TracerMock{}
}

type TracerMock struct {
	mock.Mock
}

func (t *TracerMock) Name() string {
	args := t.Called()
	return args.Get(0).(string)
}

func (t *TracerMock) Start() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TracerMock) Stop() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TracerMock) Config() *tracerConfig.Config {
	args := t.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*tracerConfig.Config)
}

func (t *TracerMock) ConfigFile() string {
	args := t.Called()
	return args.Get(0).(string)
}

func (t *TracerMock) Trace(ctx context.Context, spanName string, data map[string]any, err error) {
	t.Called(ctx, spanName, data, err)
}

func (t *TracerMock) TraceCurrentSpan(ctx context.Context, data map[string]any, err error) {
	t.Called(ctx, data, err)
}

// Started true if started
func (t *TracerMock) Started() bool {
	args := t.Called()
	return args.Get(0).(bool)
}
