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
	GetTokenInternalProviders           = GroupV1PAuth.NewEndpoint("/:provider/login", http.MethodPost)
	LoginByExternalProvider             = GroupV1PAuth.NewEndpoint("/:provider/login", http.MethodGet)
	GetTokenByCallBackExternalProviders = GroupV1PAuth.NewEndpoint("/:provider/callback", http.MethodGet)
	SignupInternalProviders             = GroupV1PAuth.NewEndpoint("/:provider/signup", http.MethodPost)
	SignupInternalProvidersConfirmation = GroupV1PAuth.NewEndpoint("/:provider/signup/confirmation", http.MethodPost)
	CreatePassword                      = GroupV1Auth.NewEndpoint("/signup/create-password", http.MethodPost)
	Logout                              = GroupV1PAuth.NewEndpoint("/logout", http.MethodPost)
	Refresh                             = GroupV1PAuth.NewEndpoint("/refresh-token", http.MethodPost)
)
