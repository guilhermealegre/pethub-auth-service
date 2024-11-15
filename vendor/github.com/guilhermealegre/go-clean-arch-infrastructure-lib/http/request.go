package http

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/context"
)

// loadRequestInfo loads the request information to domain context
func loadRequestInfo(gCtx *gin.Context) {
	if gCtx.Request != nil {
		ctx := context.NewContext(gCtx)
		// method
		ctx.SetMethod(gCtx.Request.Method)
		//path
		if gCtx.Request.URL != nil {
			ctx.SetPath(gCtx.Request.URL.Path)
			// params
			if gCtx.Request.URL.Query() != nil {
				params := make(map[string]any)
				for key, values := range gCtx.Request.URL.Query() {
					for _, value := range values {
						params[key] = value
					}
				}
				if len(params) > 0 {
					ctx.SetParams(params)
				}
			}
		}
		// body
		if gCtx.Request.Body != nil {
			body, _ := io.ReadAll(gCtx.Request.Body)
			// resetting the body buffer to the request
			gCtx.Request.Body = NewReader(bytes.NewBuffer(body), true)
			// setting the body as bytes, so it can be read multiple times
			if len(body) > 0 {
				ctx.SetBody(body)
			}
		}
	}
	gCtx.Next()
}
