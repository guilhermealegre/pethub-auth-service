package http

import (
	"net/url"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

	"github.com/gin-gonic/gin"
)

// Endpoint Struct
type Endpoint struct {
	// Group
	group *Group
	// Path
	path string
	// Method
	method string
	// Middleware
	middlewares []gin.HandlerFunc
}

// NewEndpoint creates a new endpoint
func NewEndpoint(group *Group, path string, method string) *Endpoint {
	return &Endpoint{
		group:       group,
		path:        path,
		method:      method,
		middlewares: make([]gin.HandlerFunc, 0),
	}
}

// Path endpoint path
func (e *Endpoint) Path() string {
	return e.path
}

// Method http method
func (e *Endpoint) Method() string {
	return e.method
}

// AddMiddlewares sets middlewares with handle functions
func (e *Endpoint) AddMiddlewares(middlewares domain.IMiddleware) {
	e.middlewares = append(e.middlewares, middlewares.GetHandlers()...)
}

// SetRoute sets a route with handle functions
func (e *Endpoint) SetRoute(engine *gin.Engine, handlerFunc ...gin.HandlerFunc) {
	e.middlewares = append(e.middlewares, handlerFunc...)
	e.group.Init(&engine.RouterGroup).Handle(e.Method(), e.Path(), e.middlewares...)
}

// FullPath full endpoint path
func (e *Endpoint) FullPath() string {
	path, _ := url.JoinPath(e.group.String(), e.path)
	return path
}
