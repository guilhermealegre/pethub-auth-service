package domain

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"

	coreErrors "github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	"github.com/guilhermealegre/go-clean-arch-core-lib/response"
)

type DefaultController struct {
	IController
	app IApp
}

func NewDefaultController(app IApp) *DefaultController {
	return &DefaultController{
		app: app,
	}
}

func (c *DefaultController) App() IApp {
	return c.app
}

func (c *DefaultController) Json(ctx contextDomain.IContext, data interface{}, err ...error) {
	ctx.Header("Content-Type", "application/json")
	if len(err) > 0 &&
		err[0] != nil &&
		c.app.Logger() != nil &&
		c.app.Logger().Log() != nil {

		var ve validator.ValidationErrors
		if errors.As(err[0], &ve) {
			err = ReplaceTagErrors(ve)
		}

		statusCode, resp := response.GetResponse(data, ctx.GetPagination(), ctx.GetMeta(), err...)

		body, _ := json.Marshal(resp)
		c.app.Logger().Log().Multi(
			err,
			&LoggerInfo{
				Context: ctx,
				Response: Response{
					StatusCode: statusCode,
					Body:       string(body),
				},
			},
		)
		ctx.JSON(statusCode, resp)
	} else {
		ctx.JSON(response.GetResponse(data, ctx.GetPagination(), ctx.GetMeta(), err...))
	}
}

type LoggerInfo struct {
	Context  contextDomain.IContext `json:"context"`
	Level    coreErrors.Level       `json:"level"`
	Msg      string                 `json:"msg"`
	SubType  SubType                `json:"log"`
	Response Response               `json:"response"`
}

type Backend struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
	Body   string `json:"body"`
}

type Response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

type Frontend struct {
	Name             string  `json:"name"`
	Message          string  `json:"message"`
	Time             string  `json:"time"`
	UserAgent        string  `json:"userAgent"`
	Location         string  `json:"location"`
	ViewportHeight   int     `json:"viewportHeight"`
	ViewportWidth    int     `json:"viewportWidth"`
	Stack            string  `json:"stack"`
	Type             string  `json:"type"`
	History          string  `json:"history"`
	BrowserName      string  `json:"browserName"`
	BrowserVersion   float64 `json:"browserVersion"`
	BrowserSupported bool    `json:"browserSupported"`
	Os               string  `json:"os"`
	Payload          string  `json:"payload"`
}

type SubType string
