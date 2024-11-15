package v1

import (
	"github.com/gin-gonic/gin"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
)

type IController interface {
	GetPublic(ctx *gin.Context)
}

type IModel interface {
	Get(ctx ctxDomain.IContext) (*Alive, error)
	GetPublic(ctx ctxDomain.IContext) (*PublicAlive, error)
}
