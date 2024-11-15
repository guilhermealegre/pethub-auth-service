package validator

import (
	"github.com/go-playground/validator/v10"
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

// Validator
type Validator struct {
	// Name
	name string
	// App
	app domain.IApp
	// Field Validators
	fieldValidators []domain.IFieldValidator
	// Struct Validators
	structValidators []domain.IStructValidator
	// Validator Client
	validator *validator.Validate
	// Started
	started bool
}

// New creates a new validator
func New(app domain.IApp) *Validator {
	return &Validator{
		name:      "Validator",
		app:       app,
		validator: validator.New(),
	}
}

// Name of the service
func (v *Validator) Name() string {
	return v.name
}

// Start starts the service
func (v *Validator) Start() error {
	// default custom field validators
	for tag, handlerFunc := range mapCustomDefaultValidators {
		if err := v.validator.RegisterValidationCtx(tag, handlerFunc); err != nil {
			return err
		}
	}

	// registering field validations
	for _, validatorList := range v.fieldValidators {
		if err := v.validator.RegisterValidationCtx(validatorList.Tag(), validatorList.Func(v.app)); err != nil {
			return err
		}

	}

	// registering struct validations
	for _, validatorList := range v.structValidators {
		v.validator.RegisterStructValidationCtx(validatorList.Func(v.app), validatorList.Struct())
	}

	v.started = true

	return nil
}

// Stop stops the service
func (v *Validator) Stop() error {
	if !v.started {
		return nil
	}
	v.started = false
	return nil
}

// AddFieldValidators adds a custom field validator
func (v *Validator) AddFieldValidators(fv ...domain.IFieldValidator) domain.IValidator {
	v.fieldValidators = append(v.fieldValidators, fv...)
	return v
}

// AddStructValidators adds a custom struct validator
func (v *Validator) AddStructValidators(sv ...domain.IStructValidator) domain.IValidator {
	v.structValidators = append(v.structValidators, sv...)
	return v
}

// Validate validates a struct
func (v *Validator) Validate(ctx contextDomain.IContext, obj any) error {
	return v.validator.StructCtx(ctx.RequestContext(), obj)
}

// Started true if started
func (v *Validator) Started() bool {
	return v.started
}
