package domain

import (
	"unicode"

	"github.com/go-playground/validator/v10"

	"regexp"
	"strings"

	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z.]+`)

	tagMappingError = map[string]func(fe validator.FieldError) error{
		"required": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldEmpty, camelCaseField(fe.Field()))
		},
		"gt": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMinValue, camelCaseField(fe.Field()), fe.Param())
		},
		"gte": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMinValue, camelCaseField(fe.Field()), fe.Param())
		},
		"lt": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMaxValue, camelCaseField(fe.Field()), fe.Param())
		},
		"lte": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMaxValue, camelCaseField(fe.Field()), fe.Param())
		},
		"len": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldSize, camelCaseField(fe.Field()), fe.Param())
		},
		"max": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMaxSize, camelCaseField(fe.Field()), fe.Param())
		},
		"min": func(fe validator.FieldError) error {
			return format(fe, errorCodes.ErrorFieldMinSize, camelCaseField(fe.Field()), fe.Param())
		},
	}
)

func format(fe validator.FieldError, errFunc func() errors.ErrorDetails, args ...interface{}) error {
	namespace := nonAlphanumericRegex.ReplaceAllString(fe.StructNamespace(), "")
	path := strings.Split(namespace, ".")
	for i := range path[1:] {
		path[i+1] = camelCaseField(path[i+1])
	}
	return errFunc().Formats(args...).SetField(strings.Join(path[1:], "."))
}

// camelCaseField we assume that all field names start by a lower letter
func camelCaseField(str string) string {
	value := []rune(str)
	value[0] = unicode.ToLower(value[0])
	return string(value)
}

func ReplaceTagErrors(ve validator.ValidationErrors) (errs []error) {
	for _, fe := range ve {
		var newErr error
		if errFunc, ok := tagMappingError[fe.ActualTag()]; ok {
			newErr = errFunc(fe)
		} else {
			newErr = format(fe, errorCodes.ErrorInvalidInputFields, fe.Namespace(), camelCaseField(fe.Field()), fe.Tag())
		}
		errs = append(errs, newErr)
	}

	return errs
}
