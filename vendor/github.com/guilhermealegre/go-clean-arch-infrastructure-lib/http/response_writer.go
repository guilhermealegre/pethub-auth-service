package http

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// responseWriter is a wrapper of gin.ResponseWriter
type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// newResponseWriter returns a responseWriter
func newResponseWriter(gCtx *gin.Context) *responseWriter {
	return &responseWriter{
		ResponseWriter: gCtx.Writer,
		Body:           new(bytes.Buffer),
	}
}

// Write captures the body
func (w *responseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

// captureResponseWriter captures the response writer
func (w *responseWriter) captureResponseWriter(gCtx *gin.Context) {
	gCtx.Writer = w
}

// getBody returns the response body
func (w *responseWriter) getBody() string {
	return w.Body.String()
}
