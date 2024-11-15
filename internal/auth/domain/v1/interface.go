package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
)

type IController interface {
	domain.IController
	GetTokenInternalProviders(gCtx *gin.Context)
}

type IModel interface {
	GetToken(ctx ctxDomain.IContext, loginIdentifier, identifierType, password string) (*TokenPair, error)
	Signup(ctx ctxDomain.IContext, loginIdentifier, typeIdentifier string) (err error)
	SignupEmailConfirmation(ctx ctxDomain.IContext, email, code string) (*TokenPair, error)
	Logout(ctx ctxDomain.IContext, idUser int) error
	RefreshToken(ctx ctxDomain.IContext, refreshToken string) (*TokenPair, error)
	CreatePassword(ctx ctxDomain.IContext, email string, password string) (*TokenPair, error)
	GetTokenByExternalProviders(ctx ctxDomain.IContext) (*TokenPair, error)
	LoginExternalProviders(ctx ctxDomain.IContext)
}

type IRepository interface {
	GetAuthDetails(ctx ctxDomain.IContext, tx dbr.SessionRunner, loginIdentifier, identifierType string) (*UserAuthDetails, error)
	CreateAuth(ctx ctxDomain.IContext, signup *UserAuthDetails) error
	SaveConfirmationCode(ctx ctxDomain.IContext, loginIdentifier, code string) error
	GetConfirmationCode(ctx ctxDomain.IContext, loginIdentifier string) (string, error)
	CreatePassword(ctx ctxDomain.IContext, tx session.ITx, email string, hashedPassword []byte) (uuid.UUID, error)
}

type IStreaming interface {
	CreateUser(ctx ctxDomain.IContext, uuid2 uuid.UUID) (int, error)
	GetUserDetails(ctx ctxDomain.IContext, idUser int) (*UserDetails, error)
	SendEmailSignupConfirmationCode(ctx ctxDomain.IContext, email, confirmationCode string) error
	SendSMSSignupConfirmationCode(ctx ctxDomain.IContext, number, confirmationCode string) error
}
