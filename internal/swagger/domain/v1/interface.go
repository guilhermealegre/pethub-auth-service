package v1

import (
	"github.com/gin-gonic/gin"
)

type IController interface {
	Docs(version int) func(ctx *gin.Context)
	Swagger(version int) func(ctx *gin.Context)
	StaticFile(version int) (relativePath, filePath string)
}
