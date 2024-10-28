package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/context"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	v1Routes "github.com/guilhermealegre/pethub-auth-service/api/v1/http"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/request"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/response"
	v1 "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
)

type Controller struct {
	*domain.DefaultController
	model v1.IModel
}

func NewController(app domain.IApp, model v1.IModel) v1.IController {
	return &Controller{
		DefaultController: domain.NewDefaultController(app),
		model:             model,
	}
}

func (c *Controller) Register() {
	engine := c.App().Http().Router()
	v1Routes.GetTokenByEmail.SetRoute(engine, c.GetTokenByEmail)
	v1Routes.GetTokenByPhoneNumber.SetRoute(engine, c.GetTokenByPhoneNumber)
	v1Routes.GetTokenByExternalProviders.SetRoute(engine, c.GetTokenByExternalProviders)
	v1Routes.CallbackByExternalProviders.SetRoute(engine, c.CallbackByExternalProviders)
	v1Routes.EmailSignup.SetRoute(engine, c.EmailSignup)
	v1Routes.EmailSignupConfirmation.SetRoute(engine, c.EmailSignupConfirmation)
	v1Routes.PhoneNumberSignup.SetRoute(engine, c.PhoneNumberSignup)
	v1Routes.CreatePassword.SetRoute(engine, c.CreatePassword)
	v1Routes.Logout.SetRoute(engine, c.Logout)
	v1Routes.Refresh.SetRoute(engine, c.Refresh)
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

	obj, err := c.model.GetToken(ctx, req.Body.Email, v1.Email, req.Body.Password)

	c.Json(ctx, obj.FromDomainToAPI(), err)
}

func (c *Controller) GetTokenByPhoneNumber(gCtx *gin.Context) {

}

func (c *Controller) GetTokenByExternalProviders(gCtx *gin.Context) {

}

func (c *Controller) CallbackByExternalProviders(gCtx *gin.Context) {

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
func (c *Controller) EmailSignup(gCtx *gin.Context) {
	ctx := context.NewContext(gCtx)

	req := request.SignUpEmailRequest{}
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	if err := c.App().Validator().Validate(ctx, req); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	err := c.model.Signup(ctx, req.Body.Email, v1.Email)

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

	err := c.model.Signup(ctx, req.Body.PhoneNumber, v1.PhoneNumber)

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

	tokenPair, err := c.model.CreatePassword(ctx, ctx.GetUser().Email, req.Body.Password)
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

	err := c.model.Logout(ctx, ctx.GetUser().Id)

	c.Json(ctx, response.SuccessResponse{Success: err == nil}, err)
}

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
	var req request.EmailSignupConfirmation
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		c.Json(ctx, nil, err)
		return
	}

	tokenPair, err := c.model.SignupEmailConfirmation(ctx, req.Body.Email, req.Body.Code)

	c.Json(ctx, tokenPair, err)
}
