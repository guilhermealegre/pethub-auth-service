package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

type CustomValidator struct {
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

func (c *CustomValidator) Tag() string {
	return "custom"
}

func (c *CustomValidator) Func(a domain.IApp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		// split values using ` `. eg. notoneof=bob rob job
		match := strings.Split(fl.Param(), " ")
		// convert field value to string
		value := fl.Field().String()
		for _, s := range match {
			// match value with struct filed tag
			if s == value {
				return false
			}
		}
		return true
	}
}
