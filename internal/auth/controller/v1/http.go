package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/context"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	v1Routes "github.com/guilhermealegre/pethub-auth-service/api/v1/http"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/request"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/response"
	"github.com/guilhermealegre/pethub-auth-service/internal"
	v1Domain "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
)

type Controller struct {
	*domain.DefaultController
	model v1Domain.IModel
}

func NewController(app domain.IApp, model v1Domain.IModel) v1Domain.IController {
	return &Controller{
		DefaultController: domain.NewDefaultController(app),
		model:             model,
	}
}

func (c *Controller) Register() {
	engine := c.App().Http().Router()
	v1Routes.GetTokenInternalProviders.SetRoute(engine, c.GetTokenInternalProviders)
	v1Routes.GetTokenByCallBackExternalProviders.SetRoute(engine, c.GetTokenCallBackByExternalProviders)
	v1Routes.LoginByExternalProvider.SetRoute(engine, c.LoginExternalProviders)
	v1Routes.SignupInternalProviders.SetRoute(engine, c.SignupInternalProviders)
	v1Routes.SignupInternalProvidersConfirmation.SetRoute(engine, c.EmailSignupConfirmation)
	v1Routes.CreatePassword.SetRoute(engine, c.CreatePassword)
	v1Routes.Logout.SetRoute(engine, c.Logout)
	v1Routes.Refresh.SetRoute(engine, c.Refresh)
}

func (c *Controller) GetTokenInternalProviders(gCtx *gin.Context) {
	provider := gCtx.Param("provider")
	switch provider {
	case v1Domain.Email:
		c.GetTokenByEmail(gCtx)
	case v1Domain.PhoneNumber:
		c.GetTokenByPhoneNumber(gCtx)
	default:
		ctx := context.NewContext(gCtx)
		c.Json(ctx, nil, internal.ErrorInvalidProvider())
	}
}

/*
	 swagger:route POST /p/auth/{provide} auth GetTokenByEmailRequest

	 Login authenticate user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) GetTokenByEmail(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	var req request.GetTokenByEmailRequest

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	if err := c.App().Validator().Validate(ctx, req); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	obj, err := c.model.GetToken(ctx, req.Body.Email, v1Domain.Email, req.Body.Password)

	c.Json(ctx, obj.FromDomainToAPI(), err)
}

/*
	 swagger:route POST /p/auth/{provide} auth GetTokenByPhoneNumberRequest

	 Login authenticate user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) GetTokenByPhoneNumber(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	var req request.GetTokenByPhoneNumberRequest

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	if err := c.App().Validator().Validate(ctx, req); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	obj, err := c.model.GetToken(ctx, req.Body.PhoneNumber, v1Domain.PhoneNumber, req.Body.Password)

	c.Json(ctx, obj.FromDomainToAPI(), err)
}

/*
	 swagger:route POST /p/auth/email/signup auth SignUpEmailRequest

	 Signup authenticated user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) LoginExternalProviders(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	provider := ctx.Param("provider")
	q := ctx.Request().URL.Query()
	q.Set("provider", provider)
	ctx.Request().URL.RawQuery = q.Encode()

	_, err := c.model.GetTokenByExternalProviders(ctx)
	if err != nil {
		return
	}

	return
}

/*
	 swagger:route POST /p/auth/email/signup auth SignUpEmailRequest

	 Signup authenticated user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) GetTokenCallBackByExternalProviders(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	provider := ctx.Param("provider")
	q := ctx.Request().URL.Query()
	q.Set("provider", provider)
	ctx.Request().URL.RawQuery = q.Encode()

	obj, err := c.model.GetTokenByExternalProviders(ctx)

	c.Json(ctx, obj.FromDomainToAPI(), err)

}

/*
	 swagger:route POST /p/auth/email/signup auth SignUpEmailRequest

	 Signup authenticated user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) SignupInternalProviders(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	var loginIdentifier string
	provider := ctx.Param("provider")

	switch provider {
	case v1Domain.Email:
		req := request.SignUpEmailRequest{}
		if err := ctx.ShouldBindJSON(&req.Body); err != nil {
			c.Json(ctx, nil, err)
			return
		}

		if err := c.App().Validator().Validate(ctx, req); err != nil {
			c.Json(ctx, nil, err)
			return
		}
		loginIdentifier = req.Body.Email

	case v1Domain.PhoneNumber:
		req := request.SignUpPhoneNumberRequest{}
		if err := ctx.ShouldBindJSON(&req.Body); err != nil {
			c.Json(ctx, nil, err)
			return
		}

		if err := c.App().Validator().Validate(ctx, req); err != nil {
			c.Json(ctx, nil, err)
			return
		}
		loginIdentifier = req.Body.CodePhoneNumber + req.Body.PhoneNumber
	}

	err := c.model.Signup(ctx, loginIdentifier, provider)

	c.Json(ctx, response.SuccessResponse{Success: err == nil}, err)
}

/*
	 swagger:route POST /p/auth/phone-number/signup auth SignUpEmailRequest

	 Signup authenticated user

		Produces:
		- application/json

		Responses:
		  200: SwaggerScannerLoginResponse
		  400: ErrorResponse
*/
func (c *Controller) PhoneNumberSignup(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)

	req := request.SignUpPhoneNumberRequest{}
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	if err := c.App().Validator().Validate(ctx, req); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	err := c.model.Signup(ctx, req.Body.PhoneNumber, v1Domain.PhoneNumber)
	c.Json(ctx, response.SuccessResponse{Success: err == nil}, err)
}

func (c *Controller) CreatePassword(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)

	req := request.CreatePasswordRequest{}
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	if err := c.App().Validator().Validate(ctx, req); err != nil {
		c.Json(ctx, nil, err)
		return
	}
	//ctx.GetUser().Email,
	email := "sdfa@asf.com"

	tokenPair, err := c.model.CreatePassword(ctx, email, req.Body.Password)
	c.Json(ctx, tokenPair.FromDomainToAPI(), err)
}

/*
	 swagger:route POST /p/auth/logout auth AuthRequest

	 Logout authenticated user

		Produces:
		- application/json

		Responses:
		  200: SuccessResponse
		  400: ErrorResponse
*/
func (c *Controller) Logout(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	//ctx.GetUser().Id

	err := c.model.Logout(ctx, 1)

	c.Json(ctx, response.SuccessResponse{Success: err == nil}, err)
}

/*
	 swagger:route POST /p/auth/logout auth AuthRequest

	 Logout authenticated user

		Produces:
		- application/json

		Responses:
		  200: SuccessResponse
		  400: ErrorResponse
*/
func (c *Controller) Refresh(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)

	obj, err := c.model.RefreshToken(ctx, "")

	c.Json(ctx, obj.FromDomainToAPI(), err)
}

/*
	 swagger:route POST /p/auth/email/signup/confirm auth EmailSignupConfirmRequest

	 Confirm email signup

		Produces:
		- application/json

		Responses:
		  200: SuccessResponse
		  400: ErrorResponse
*/

func (c *Controller) EmailSignupConfirmation(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)
	var req request.EmailSignupConfirmationRequest
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	tokenPair, err := c.model.SignupEmailConfirmation(ctx, req.Body.Email, req.Body.Code)

	c.Json(ctx, tokenPair, err)
}
