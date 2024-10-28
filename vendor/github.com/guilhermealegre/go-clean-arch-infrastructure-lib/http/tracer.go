package http

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/context"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer"
)

// traceRequest traces the request information
func (h *Http) traceRequest(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	attrs := make(map[string]any)
	if ctx.GetBody() != nil {
		attrs[tracer.TracerTagRequestBody] = string(ctx.GetBody())
	}
	if ctx.GetParams() != nil {
		attrs[tracer.TracerTagParams] = ctx.GetParams()
	}

	writer := newResponseWriter(gCtx)
	writer.captureResponseWriter(gCtx)

	gCtx.Next()

	attrs[tracer.TracerTagResponseBody] = writer.getBody()

	h.app.Tracer().TraceCurrentSpan(ctx.RequestContext(), attrs, nil)
}
