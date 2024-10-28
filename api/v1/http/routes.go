package http

import (
	"net/http"

	infra "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/http"
)

var (
	GroupV1  = infra.NewGroup("api").Group("v1")
	GroupV1P = GroupV1.Group("p")
	//auth
	GroupV1Auth  = GroupV1.Group("auth")
	GroupV1PAuth = GroupV1P.Group("auth")

	//documentation
	GroupV1PDocumentation     = GroupV1P.Group("documentation")
	GroupV1PDocumentationUser = GroupV1PDocumentation.Group("auth")

	SwaggerUserDocs    = GroupV1PDocumentationUser.NewEndpoint("/docs", http.MethodGet)
	SwaggerUserSwagger = GroupV1PDocumentationUser.NewEndpoint("/swagger", http.MethodGet)

	Alive           = GroupV1.NewEndpoint("/alive", http.MethodGet)
	PublicAliveUser = GroupV1P.NewEndpoint("/alive/auth", http.MethodGet)

	//auth
	GetTokenByEmail             = GroupV1PAuth.NewEndpoint("/email/login", http.MethodPost)
	GetTokenByPhoneNumber       = GroupV1PAuth.NewEndpoint("/phone-number/login", http.MethodPost)
	GetTokenByExternalProviders = GroupV1PAuth.NewEndpoint("/:provider", http.MethodPost)
	CallbackByExternalProviders = GroupV1PAuth.NewEndpoint("/:provider/callback", http.MethodPost)
	EmailSignup                 = GroupV1PAuth.NewEndpoint("/email/signup", http.MethodPost)
	EmailSignupConfirmation     = GroupV1PAuth.NewEndpoint("/email/signup/confirmation", http.MethodPost)
	PhoneNumberSignup           = GroupV1PAuth.NewEndpoint("/phone-number/signup", http.MethodPost)
	CreatePassword              = GroupV1Auth.NewEndpoint("/signup/create-password", http.MethodPost)
	Logout                      = GroupV1PAuth.NewEndpoint("/logout", http.MethodPost)
	Refresh                     = GroupV1PAuth.NewEndpoint("/refresh-token", http.MethodPost)
)
