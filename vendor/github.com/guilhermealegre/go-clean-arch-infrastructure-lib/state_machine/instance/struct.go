package state_machine

import (
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
)

// StateMachine ...
type StateMachine struct {
	Name                      string `json:"name"`
	execute                   HandlerExecFunction
	stateMachinesToTriggerMap map[string]IStateMachine
	currentState              CurrentStateFunc
	States                    []StateInput                      `json:"states"`
	MapStates                 map[string]map[string]Handlers    `json:"map_states"`
	OnSuccessHandlers         map[string]HandlerFunc            `json:"on_success_handlers"`
	OnErrorHandlers           map[string]HandlerFunc            `json:"on_error_handlers"`
	CheckHandlers             map[string]HandlerFunc            `json:"check_handlers"`
	FilterHandlers            map[string]HandlerFilterFunction  `json:"filter_handlers"`
	AdapterHandlers           map[string]HandlerAdapterFunction `json:"adapter_handlers"`
}

type StateInput struct {
	Name        string            `json:"name"`
	Transitions []TransitionInput `json:"transitions"`
}

type TransitionInput struct {
	// Name
	Name string `json:"name"`
	// Check
	Check []CheckInputStruct `json:"check" mapstructure:"check"`
	// On Success
	OnSuccess []OnSuccessInputStruct `json:"on_success"  mapstructure:"on_success"`
	// On Error
	OnError []OnErrorInputStruct `json:"on_error" mapstructure:"on_error"`
}

type CheckInputStruct struct {
	Func string `json:"func"`
}

type OnSuccessInputStruct struct {
	Check           []CheckInputStruct `json:"check"`
	Func            string             `json:"func"`
	FuncArg         []string           `json:"func_arg"`
	Adapter         string             `json:"adapter"`
	Filter          string             `json:"filter"`
	IsStateMachine  bool               `json:"is_state_machine" mapstructure:"is_state_machine"`
	IgnoreError     bool               `json:"ignore_error,omitempty" mapstructure:"ignore_error"`
	IgnoreNoSuccess bool               `json:"ignore_no_success,omitempty" mapstructure:"ignore_no_success"`
}

type OnErrorInputStruct struct {
	Check []CheckInputStruct `json:"check"`
	Func  string             `json:"func"`
}

type Handlers struct {
	// Check
	Check []CheckStruct `json:"check"`
	// On Success
	OnSuccess []OnSuccessStruct `json:"on_success"`
	// On Error
	OnError []OnErrorStruct `json:"on_error"`
}

type CheckStruct struct {
	Func            string   `json:"func"`
	FuncArg         []string `json:"func_arg"`
	IgnoreError     bool     `json:"ignore_error,omitempty" mapstructure:"ignore_error"`
	IgnoreNoSuccess bool     `json:"ignore_no_success,omitempty" mapstructure:"ignore_no_success"`
}

type OnSuccessStruct struct {
	Check           []CheckStruct `json:"check"`
	Func            string        `json:"func"`
	FuncArg         []string      `json:"func_arg"`
	Adapter         string        `json:"adapter"`
	Filter          string        `json:"filter"`
	IsStateMachine  bool          `json:"is_state_machine" mapstructure:"is_state_machine"`
	IgnoreError     bool          `json:"ignore_error,omitempty" mapstructure:"ignore_error"`
	IgnoreNoSuccess bool          `json:"ignore_no_success,omitempty" mapstructure:"ignore_no_success"`
}

type OnErrorStruct struct {
	Check           []CheckStruct `json:"check"`
	Func            string        `json:"func"`
	FuncArg         []string      `json:"func_arg"`
	IgnoreError     bool          `json:"ignore_error,omitempty" mapstructure:"ignore_error"`
	IgnoreNoSuccess bool          `json:"ignore_no_success,omitempty" mapstructure:"ignore_no_success"`
}

type HandlerAdapterFunction func(ctx contextDomain.IContext, nextState string, obj any) ([]any, error)
type HandlerFilterFunction func(ctx contextDomain.IContext, nextState string, objs []any) ([]any, error)
type HandlerExecFunction func(ctx contextDomain.IContext, nextState string, obj any) (err error)
type HandlerFunc func(ctx contextDomain.IContext, nextState string, arg any, optArg ...string) (success bool, err error)
type CurrentStateFunc func(ctx contextDomain.IContext, obj any) (string, error)
