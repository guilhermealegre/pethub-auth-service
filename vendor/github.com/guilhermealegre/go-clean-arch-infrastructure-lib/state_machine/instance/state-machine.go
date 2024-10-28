package state_machine

import (
	"os"
	"regexp"
	"strings"

	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	"github.com/spf13/viper"
)

func NewStateMachine() IStateMachine {
	return &StateMachine{
		MapStates:                 make(map[string]map[string]Handlers),
		stateMachinesToTriggerMap: make(map[string]IStateMachine),
		CheckHandlers:             make(map[string]HandlerFunc),
		OnSuccessHandlers:         make(map[string]HandlerFunc),
		OnErrorHandlers:           make(map[string]HandlerFunc),
		AdapterHandlers:           make(map[string]HandlerAdapterFunction),
		FilterHandlers:            make(map[string]HandlerFilterFunction),
	}
}

func (sm *StateMachine) Load(filePath string) error {
	// Read the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	//var test any
	viper.SetConfigFile(filePath)

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	if err = viper.Unmarshal(&sm); err != nil {
		return err
	}

	// initialize the state machine
	for _, state := range sm.States {

		if sm.MapStates[state.Name] == nil {
			sm.MapStates[state.Name] = make(map[string]Handlers)
		}

		for _, transition := range state.Transitions {
			var handlers Handlers
			// add check handlers
			for _, check := range transition.Check {
				funcName, args := splitFunctionAndArguments(check.Func)
				handlers.Check = append(handlers.Check, CheckStruct{
					Func:    funcName,
					FuncArg: args,
				})
			}
			// add on_success handlers
			for _, onSuccess := range transition.OnSuccess {
				funcName, args := splitFunctionAndArguments(onSuccess.Func)
				addOnSuccess := OnSuccessStruct{
					Func:            funcName,
					FuncArg:         args,
					Adapter:         onSuccess.Adapter,
					Filter:          onSuccess.Filter,
					IsStateMachine:  onSuccess.IsStateMachine,
					IgnoreError:     onSuccess.IgnoreError,
					IgnoreNoSuccess: onSuccess.IgnoreNoSuccess,
				}

				for _, check := range onSuccess.Check {
					funcName, args := splitFunctionAndArguments(check.Func)
					addOnSuccess.Check = append(addOnSuccess.Check, CheckStruct{
						Func:    funcName,
						FuncArg: args,
					})
				}

				handlers.OnSuccess = append(handlers.OnSuccess, addOnSuccess)
			}
			// add on_error handlers
			for _, onError := range transition.OnError {
				funcName, args := splitFunctionAndArguments(onError.Func)
				addOnError := OnErrorStruct{
					Func:    funcName,
					FuncArg: args,
				}

				for _, check := range onError.Check {
					funcName, args := splitFunctionAndArguments(check.Func)
					addOnError.Check = append(addOnError.Check, CheckStruct{
						Func:    funcName,
						FuncArg: args,
					})
				}

				handlers.OnError = append(handlers.OnError, addOnError)
			}

			sm.MapStates[state.Name][transition.Name] = handlers
		}
	}

	return nil
}

func (sm *StateMachine) GetName() string {
	return sm.Name
}

func (sm *StateMachine) AddCheckFunction(name string, handler HandlerFunc) {
	sm.CheckHandlers[name] = handler
}

func (sm *StateMachine) AddOnErrorFunction(name string, handler HandlerFunc) {
	sm.OnErrorHandlers[name] = handler
}

func (sm *StateMachine) AddOnSuccessFunction(name string, handler HandlerFunc) {
	sm.OnSuccessHandlers[name] = handler
}

func (sm *StateMachine) AddExecuteFunction(handler HandlerExecFunction) {
	sm.execute = handler
}

func (sm *StateMachine) AddStateMachineToTrigger(name string, stateMachine IStateMachine) IStateMachine {
	sm.stateMachinesToTriggerMap[name] = stateMachine
	return sm
}

func (sm *StateMachine) AddAdapterFunction(name string, handler HandlerAdapterFunction) {
	sm.AdapterHandlers[name] = handler
}

func (sm *StateMachine) AddFilterFunction(name string, handler HandlerFilterFunction) {
	sm.FilterHandlers[name] = handler
}

func (sm *StateMachine) AddCurrentStateFunction(handler CurrentStateFunc) {
	sm.currentState = handler
}

func (sm *StateMachine) ProcessTransition(ctx contextDomain.IContext, nextState string, obj any) (success bool, err error) {
	// Get handlers
	currentState, err := sm.currentState(ctx, obj)
	if err != nil {
		return false, err
	}

	handlers, exitTransition := sm.MapStates[currentState][nextState]
	if !exitTransition {
		return false, errors.ErrorInStateMachineTransition().Formats(currentState, nextState, sm.Name)
	}

	success, err = sm.runCheckFunction(ctx, nextState, handlers.Check, obj)
	if err != nil {
		if success, err := sm.runOnErrorFunction(ctx, nextState, handlers.OnError, obj); err != nil {
			return success, err
		}
		return success, err
	}

	if !success {
		return false, nil
	}

	err = sm.execute(ctx, nextState, obj)
	if err != nil {
		if success, err := sm.runOnErrorFunction(ctx, nextState, handlers.OnError, obj); err != nil {
			return success, err
		}
		return success, err
	}

	success, err = sm.runOnSuccessFunction(ctx, nextState, handlers.OnSuccess, obj)
	if err != nil {
		if success, err := sm.runOnErrorFunction(ctx, nextState, handlers.OnError, obj); err != nil {
			return success, err
		}
		return success, err
	}

	return success, nil
}

func (sm *StateMachine) getCheckFunction(name string) HandlerFunc {
	// it's an internal function
	return sm.CheckHandlers[name]
}
func (sm *StateMachine) getOnErrorFunction(name string) HandlerFunc {
	return sm.OnErrorHandlers[name]
}

func (sm *StateMachine) getOnSuccessFunction(name string) HandlerFunc {
	return sm.OnSuccessHandlers[name]
}

func (sm *StateMachine) getAdapterFunction(name string) HandlerAdapterFunction {
	return sm.AdapterHandlers[name]
}

func (sm *StateMachine) getFilterFunction(name string) HandlerFilterFunction {
	return sm.FilterHandlers[name]
}

func (sm *StateMachine) getStateMachineToTrigger(name string) IStateMachine {
	return sm.stateMachinesToTriggerMap[name]
}

func (sm *StateMachine) runCheckFunction(ctx contextDomain.IContext, nextState string, handlers []CheckStruct, obj any) (success bool, err error) {
	success = true
	for _, handler := range handlers {
		handlerFunc := sm.getCheckFunction(handler.Func)

		success, err = handlerFunc(ctx, nextState, obj, handler.FuncArg...)
		if err != nil && !handler.IgnoreError {
			return false, err
		}

		if !success && !handler.IgnoreNoSuccess {
			return false, nil
		}
	}

	return success, nil
}

func (sm *StateMachine) runOnErrorFunction(ctx contextDomain.IContext, nextState string, handlers []OnErrorStruct, obj any) (bool, error) {
	for _, handler := range handlers {
		success, err := sm.runCheckFunction(ctx, nextState, handler.Check, obj)
		if err != nil && !handler.IgnoreError {
			return false, err
		}

		if !success && !handler.IgnoreNoSuccess {
			return false, nil
		}

		handlerFunc := sm.getOnErrorFunction(handler.Func)
		success, err = handlerFunc(ctx, nextState, obj)
		if err != nil && !handler.IgnoreError {
			return false, err
		}

		if !success && !handler.IgnoreNoSuccess {
			return false, nil
		}
	}

	return true, nil
}

func (sm *StateMachine) runOnSuccessFunction(ctx contextDomain.IContext, nextState string, handlers []OnSuccessStruct, obj any) (bool, error) {
	for _, handler := range handlers {
		objs := []any{obj}

		adapter := sm.getAdapterFunction(handler.Adapter)
		if adapter != nil {
			newObjs, err := adapter(ctx, nextState, obj)
			if err != nil {
				return false, err
			}
			objs = newObjs
		}

		filter := sm.getFilterFunction(handler.Filter)
		if filter != nil {
			newObjs, err := filter(ctx, nextState, objs)
			if err != nil {
				return false, err
			}
			objs = newObjs
		}

		for _, obj := range objs {
			success, err := sm.runCheckFunction(ctx, nextState, handler.Check, obj)
			if err != nil {
				return false, err
			}

			if !success {
				continue
			}

			if handler.IsStateMachine {
				smTrigger := sm.getStateMachineToTrigger(handler.Func)
				if smTrigger != nil {
					success, err := smTrigger.ProcessTransition(ctx, handler.FuncArg[1], obj)
					if err != nil && !handler.IgnoreError {
						return false, err
					}

					if !success && !handler.IgnoreNoSuccess {
						return false, nil
					}
				}
			} else {
				handlerFunc := sm.getOnSuccessFunction(handler.Func)
				success, err := handlerFunc(ctx, nextState, obj, handler.FuncArg...)
				if err != nil && !handler.IgnoreError {
					return false, err
				}

				if !success && !handler.IgnoreNoSuccess {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func splitFunctionAndArguments(input string) (function string, arguments []string) {
	pattern := regexp.MustCompile(`^[\w-]+(?:-[\w-]+)?\(((?:[\w-]+(?:-\w+)?)*(?:,\s*[\w-]+(?:-\w+)*)*)?\)$`)
	if !pattern.MatchString(input) {
		return input, nil
	}

	parts := strings.Split(input, "(")
	functionName := parts[0]
	argumentsString := parts[1]
	argumentsString = strings.TrimSuffix(argumentsString, ")")
	arguments = strings.Split(argumentsString, ",")
	for i, arg := range arguments {
		arguments[i] = strings.TrimSpace(arg)
	}

	return functionName, arguments
}
