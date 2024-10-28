package state_machine

import (
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
)

// IService service interface
type IService interface {
	// Name name of the service
	Name() string
	// Start starts the service
	Start() error
	// Stop stops the service
	Stop() error
	// Started true if service started
	Started() bool
}

// IStateMachineService interface
type IStateMachineService interface {
	IService
	Get(name string) IStateMachine
}

type IStateMachine interface {
	GetName() string
	Load(filePath string) error
	ProcessTransition(ctx contextDomain.IContext, nextState string, obj any) (success bool, err error)
	AddCheckFunction(name string, handler HandlerFunc)
	AddOnErrorFunction(name string, handler HandlerFunc)
	AddOnSuccessFunction(name string, handler HandlerFunc)
	AddExecuteFunction(handler HandlerExecFunction)
	AddCurrentStateFunction(handler CurrentStateFunc)
	AddStateMachineToTrigger(name string, stateMachine IStateMachine) IStateMachine
	AddAdapterFunction(name string, handler HandlerAdapterFunction)
	AddFilterFunction(name string, handler HandlerFilterFunction)
}
