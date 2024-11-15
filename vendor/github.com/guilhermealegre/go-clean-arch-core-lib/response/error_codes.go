package response

import (
	"net/http"

	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
)

// Generic error codes
var (
	ErrorGeneric        = errors.NewErrorDetails("COR-1", "%s", errors.Error)
	ErrorUnauthorized   = errors.NewErrorDetails("COR-2", "Unauthorized - %s", errors.Error)
	ErrorInvalidJWT     = errors.NewErrorDetails("COR-3", "Invalid JWT - %s", errors.Error)
	ErrorObjectNotFound = errors.NewErrorDetails("COR-4", "Object not found", errors.Error, errors.Opt{Key: errors.OptStatusCode, Value: http.StatusNotFound})
	ErrorForbidden      = errors.NewErrorDetails("COR-5", "Forbidden - %s", errors.Error)
)
