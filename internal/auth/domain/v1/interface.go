package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

type IController interface {
	domain.IController
	EmailSignup(ctx *gin.Context)
	PhoneNumberSignup(ctx *gin.Context)
	GetTokenByEmail(ctx *gin.Context)
	GetTokenByPhoneNumber(ctx *gin.Context)
	GetTokenByExternalProviders(ctx *gin.Context)
	CallbackByExternalProviders(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type IModel interface {
	GetToken(ctx domain.IContext, loginIdentifier, identifierType, password string) (*TokenPair, error)
	Signup(ctx domain.IContext, loginIdentifier, typeIdentifier string) (err error)
	SignupEmailConfirmation(ctx domain.IContext, email, code string) (*TokenPair, error)
	Logout(ctx domain.IContext, idUser int) error
	RefreshToken(ctx domain.IContext, refreshToken string) (*TokenPair, error)
	CreatePassword(ctx domain.IContext, email string, password string) (*TokenPair, error)
}

type IRepository interface {
	GetAuthDetails(ctx domain.IContext, tx dbr.SessionRunner, loginIdentifier, identifierType string) (*UserAuthDetails, error)
	GetUserDetails(ctx domain.IContext, idUser int) (*UserDetails, error)
	CreateAuth(ctx domain.IContext, signup *UserAuthDetails) error
	SaveConfirmationCode(ctx domain.IContext, loginIdentifier, code string) error
	GetConfirmationCode(ctx domain.IContext, loginIdentifier string) (string, error)
	CreatePassword(ctx domain.IContext, tx *dbr.Tx, idUser int, email string, hashedPassword []byte) error
}

type IStreaming interface {
	CreateUser(ctx domain.IContext, tx *dbr.Tx) (int, error)
}
