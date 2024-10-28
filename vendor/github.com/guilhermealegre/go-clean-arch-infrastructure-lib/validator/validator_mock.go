package validator

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/stretchr/testify/mock"
)

func NewValidatorMock() *ValidatorMock {
	return &ValidatorMock{}
}

type ValidatorMock struct {
	mock.Mock
}

func (v *ValidatorMock) Name() string {
	args := v.Called()
	return args.Get(0).(string)
}

func (v *ValidatorMock) Start() error {
	args := v.Called()
	return args.Error(0)
}

func (v *ValidatorMock) Stop() error {
	args := v.Called()
	return args.Error(0)
}

func (v *ValidatorMock) AddFieldValidators(val ...domain.IFieldValidator) domain.IValidator {
	args := v.Called(val)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IValidator)
}

func (v *ValidatorMock) AddStructValidators(val ...domain.IStructValidator) domain.IValidator {
	args := v.Called(val)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IValidator)
}

func (v *ValidatorMock) Validate(ctx contextDomain.IContext, val any) error {
	args := v.Called(ctx, val)
	return args.Error(0)
}

// WithAdditionalConfigType sets an additional config type
func (v *ValidatorMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := v.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (v *ValidatorMock) Started() bool {
	args := v.Called()
	return args.Get(0).(bool)
}
