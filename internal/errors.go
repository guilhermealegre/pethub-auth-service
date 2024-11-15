package internal

import (
	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors/config"
)

// Generic error codes
var (
	ErrorGeneric            = config.GetError("1", "%s", errors.Error)
	ErrorInvalidInputFields = config.GetError("2", "The field: %s as invalid value: %v", errors.Info)
	ErrorInvalidLogin       = config.GetError("3", "Invalid login", errors.Info)
	ErrorInvalidEmail       = config.GetError("4", "Invalid User Email", errors.Info)
	ErrorInvalidPhoneNumber = config.GetError("5", "Invalid User Phone Number", errors.Info)
	ErrorTokenConfirmation  = config.GetError("6", "Invalid code confirmation", errors.Info)
	ErrorInvalidProvider    = config.GetError("7", "Invalid provider", errors.Info)

	// notification errors
	ErrorInvalidTemplate = config.GetError("7", "Invalid email template", errors.Info)
)
