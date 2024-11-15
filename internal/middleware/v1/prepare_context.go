package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

type PrepareContextMiddleware struct {
	app domain.IApp
}

func NewPrepareContextMiddleware(app domain.IApp) domain.IMiddleware {
	return &PrepareContextMiddleware{
		app: app,
	}
}

func (c *PrepareContextMiddleware) RegisterMiddlewares() {
	c.app.Http().Router().RouterGroup.Use(c.GetHandlers()...)
}

func (c *PrepareContextMiddleware) GetHandlers() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//auth.LoadJWTSecret(c.app.Http().Config().JwtSecret),
	}
}
