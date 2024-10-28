package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

type PrintRequestMiddleware struct {
	app domain.IApp
}

func NewPrintRequestMiddleware(app domain.IApp) domain.IMiddleware {
	return &PrintRequestMiddleware{
		app: app,
	}
}

func (c *PrintRequestMiddleware) RegisterMiddlewares() {
	c.app.Http().Router().RouterGroup.Use(c.GetHandlers()...)
}

func (c *PrintRequestMiddleware) GetHandlers() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		c.PrintRequest,
	}
}

func (c *PrintRequestMiddleware) PrintRequest(gCtx *gin.Context) {
	gCtx.Next()
	fmt.Printf("%d | %s | %s\n", gCtx.Writer.Status(), gCtx.Request.Method, gCtx.Request.URL.Path)
}
