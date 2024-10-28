package context

import (
	"context"
	"encoding/json"
	"net/http"

	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/metadata"

	msg "github.com/guilhermealegre/go-clean-arch-core-lib/pagination"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
	meta       interface{}
	pagination *msg.Pagination
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

func (c *Context) Values(key any) any {
	return c.Value(key)
}

func (c *Context) Params() gin.Params {
	return c.Context.Params
}

func (c *Context) Keys() map[string]any {
	if c == nil || c.Context == nil {
		return nil
	}

	return c.Context.Keys
}

func (c *Context) Request() *http.Request {
	return c.Context.Request
}

func (c *Context) Response() gin.ResponseWriter {
	return c.Context.Writer
}

func (c *Context) SetMethod(method string) {
	c.Context.Set(CtxMethod, method)
}

func (c *Context) SetPath(path string) {
	c.Context.Set(CtxPath, path)
}

func (c *Context) SetBody(body []byte) {
	c.Context.Set(CtxBody, body)
}

func (c *Context) SetParams(params map[string]any) {
	c.Context.Set(CtxParams, params)
}

func (c *Context) GetBody() []byte {
	body, exists := c.Context.Get(CtxBody)
	if exists {
		if b, ok := body.([]byte); ok {
			return b
		}
	}
	return nil
}

func (c *Context) GetParams() map[string]any {
	return c.Context.GetStringMap(CtxParams)
}

func (c *Context) GetMethod() string {
	return c.Context.GetString(CtxMethod)
}

func (c *Context) GetPath() string {
	return c.Context.GetString(CtxPath)
}

func (c *Context) AddMeta(meta interface{}) contextDomain.IContext {
	c.meta = meta
	return c
}

func (c *Context) AddPagination(pagination *msg.Pagination) contextDomain.IContext {
	c.pagination = pagination
	return c
}

func (c *Context) GetMeta() any {
	return c.meta
}

func (c *Context) GetPagination() *msg.Pagination {
	return c.pagination
}

func (c *Context) FromGrpc(ctx context.Context) contextDomain.IContext {
	gCtx := &gin.Context{
		Request: &http.Request{},
	}

	// setting the grpc context
	gCtx.Request = gCtx.Request.WithContext(ctx)

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		md, _ = metadata.FromOutgoingContext(ctx)
	}

	if len(md.Get(ContextGrpcKeys)) > 0 {
		_ = json.Unmarshal([]byte(md.Get(ContextGrpcKeys)[0]), &gCtx.Keys)
	}

	return NewContext(gCtx)
}

func (c *Context) ToGrpc() context.Context {
	if c.Keys() == nil {
		return c
	}
	res := make(map[string]string)
	b, _ := json.Marshal(c.Keys())
	res[ContextGrpcKeys] = string(b)
	if c.Request() == nil {
		return metadata.NewOutgoingContext(context.Background(), metadata.New(res))
	}
	return metadata.NewOutgoingContext(c.Request().Context(), metadata.New(res))
}

func (c *Context) RequestContext() context.Context {
	if c.Request() != nil {
		// return a ctx with gin and span data
		return c.Request().WithContext(
			trace.ContextWithSpan(c.Context, trace.SpanFromContext(c.Request().Context()))).
			Context()
	}

	return c
}
