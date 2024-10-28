package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

type IController interface {
	GetPublic(ctx *gin.Context)
}

type IModel interface {
	Get(ctx domain.IContext) (*Alive, error)
	GetPublic(ctx domain.IContext) (*PublicAlive, error)
}
