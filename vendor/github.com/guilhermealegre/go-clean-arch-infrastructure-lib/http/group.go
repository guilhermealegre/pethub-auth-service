package http

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

// Group group struct
type Group struct {
	name     string
	fullName string
	*gin.RouterGroup
	group       *Group
	middlewares []domain.IMiddleware
	endpoints   []Endpoint
	initalized  bool
}

// NewGroup creates a new group
func NewGroup(name string) *Group {
	n := fmt.Sprintf("/%s", strings.ReplaceAll(name, "/", ""))
	return &Group{
		name:     n,
		fullName: n,
	}
}

// Group appends a group
func (g *Group) Group(name string) *Group {
	n := strings.ReplaceAll(name, "/", "")
	path, _ := url.JoinPath(g.fullName, n)
	return &Group{
		name:     n,
		fullName: path,
		group:    g,
	}
}

// Init initialize a group
func (g *Group) Init(routerGroup *gin.RouterGroup) *Group {
	if g.RouterGroup == nil {
		if g.group != nil {
			g.RouterGroup = g.group.Init(routerGroup).RouterGroup.
				Group(g.name)
		} else {
			g.RouterGroup = routerGroup.Group(g.name)
		}
	}

	if !g.initalized {
		g.RouterGroup.Use(g.GetHandlers()...)
		g.initalized = true
	}

	return g
}

// AddMiddleware adds a new middleware
func (g *Group) AddMiddleware(m domain.IMiddleware) {
	g.middlewares = append(g.middlewares, m)
}

// NewEndpoint created a new endpoint at a group
func (g *Group) NewEndpoint(path string, method string) Endpoint {
	e := Endpoint{
		group:  g,
		path:   path,
		method: method,
	}
	g.endpoints = append(g.endpoints, e)
	return e
}

// String prints the group
func (g *Group) String() string {
	return g.fullName
}

// GetHandlers get handlers
func (g *Group) GetHandlers() (handlers []gin.HandlerFunc) {
	for _, m := range g.middlewares {
		handlers = append(handlers, m.GetHandlers()...)
	}
	return
}
