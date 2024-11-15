package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http"
)

type AuthenticateMiddleware struct {
	app domain.IApp
}

func NewAuthenticateMiddleware(app domain.IApp) domain.IMiddleware {
	return &AuthenticateMiddleware{
		app: app,
	}
}

func (c *AuthenticateMiddleware) RegisterMiddlewares() {
	http.GroupV1Auth.AddMiddleware(c)
}

func (c *AuthenticateMiddleware) GetHandlers() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//auth.LoadJWTSecret(c.app.Http().Config().JwtSecret),
		//auth.ValidateToken,
	}
}
