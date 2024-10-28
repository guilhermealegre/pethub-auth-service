package tracer

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer"
)

// tracerMiddleware
type tracerMiddleware struct {
	app domain.IApp
}

// NewTracerMiddleware creates a new tracer hook
func NewTracerMiddleware(app domain.IApp) redis.Hook {
	return &tracerMiddleware{
		app: app,
	}
}

// DialHook dial hook
func (t *tracerMiddleware) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook process hook
func (t *tracerMiddleware) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)

		attr := make(map[string]any)
		attr[tracer.TracerTagRedisCmd] = cmd.String()
		t.app.Tracer().Trace(ctx, t.app.Redis().Name(), attr, err)

		return err
	}
}

// ProcessPipelineHook process pipeline hook
func (t *tracerMiddleware) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)

		attrs := make(map[string]any)
		for i, cmd := range cmds {
			attrs[fmt.Sprintf("%s.%d", tracer.TracerTagRedisCmd, i+1)] = cmd.String()

			if cmd.Err() != nil {
				attrs[fmt.Sprintf("%s.%d.err", tracer.TracerTagRedisCmd, i+1)] = cmd.Err().Error()
			}
		}

		t.app.Tracer().Trace(ctx, t.app.Redis().Name(), attrs, err)

		return err
	}
}
