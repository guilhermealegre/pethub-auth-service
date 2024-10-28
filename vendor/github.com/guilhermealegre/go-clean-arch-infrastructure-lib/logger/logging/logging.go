package logging

import (
	"encoding/json"
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"runtime/debug"
	"strings"

	"github.com/guilhermealegre/go-clean-arch-core-lib/response"

	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/writer"
	"github.com/rs/zerolog"
)

// Init initializes the logging
func (l *Logging) Init(cfg config.Config) domain.ILogging {
	l.config = cfg
	for _, o := range l.config.Output {
		switch o {
		case OutputConsole:
			l.Writers = append(l.Writers, writer.NewConsole())
		case OutputFile:
			l.Writers = append(l.Writers, writer.NewFile(l.config.Path, l.config.File))
		case OutputRabbit:
			var queue, routingKey string
			if l.config.Rabbitmq != nil {
				queue = l.config.Rabbitmq.Queue
				routingKey = l.config.Rabbitmq.RoutingKey
			}

			l.Writers = append(l.Writers, writer.NewRabbit(l.RabbitMq, l.config.Path, queue, routingKey))
		case OutputSQS:
			var routingKey string
			if l.config.SQS != nil {
				routingKey = l.config.SQS.RoutingKey
			}

			l.Writers = append(l.Writers, writer.NewSQS(l.SQS, l.config.Path, l.config.SQS.Connection, l.config.SQS.Queue, routingKey))
		}
	}

	l.log = zerolog.
		New(zerolog.MultiLevelWriter(l.Writers...)).
		With().
		Timestamp().
		Logger()

	return l
}

// Do log the error
func (l *Logging) Do(err error, info ...*domain.LoggerInfo) {
	l.Multi([]error{err}, info...)
}

// Multi logs the multi errors
func (l *Logging) Multi(err []error, info ...*domain.LoggerInfo) {
	if len(l.Writers) == 0 {
		//If no output to write, return
		return
	}

	logSubType := domain.DefaultSubType
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var level errors.Level
	var msg string
	var ctx contextDomain.IContext
	var responseBody string
	var responseStatusCode int

	if len(info) > 0 {
		if info[0].Context != nil {
			ctx = info[0].Context
		}

		if int(info[0].Level) > 0 {
			level = info[0].Level
		}

		if info[0].SubType != "" {
			logSubType = info[0].SubType
		}

		if info[0].Msg != "" {
			msg = info[0].Msg
		}

		if info[0].Response.Body != "" {
			responseBody = info[0].Response.Body
		}

		if info[0].Response.StatusCode > 0 {
			responseStatusCode = info[0].Response.StatusCode
		}
	}

	ev := getErrorEvent(l.log, level, err[0])

	if l.config.StackTrace {
		ev.Strs("stack", getStack())
	}

	var errStr string
	switch err[0].(type) {
	case errors.ErrorDetails, errors.ErrorDetailsList:
		var bError []byte
		if len(err) > 1 {
			bError, _ = json.Marshal(err)
		} else {
			bError, _ = json.Marshal(err[0])
		}
		errStr = string(bError)

	default:
		var bError []byte

		if len(err) > 1 {
			newErrs := errors.ErrorDetailsList{}
			for _, e := range err {
				_ = newErrs.Add(response.ErrorGeneric.Formats(e.Error()))
			}

			bError, _ = json.Marshal(newErrs)
		} else {
			newErr := response.ErrorGeneric.Formats(err[0].Error())
			bError, _ = json.Marshal(newErr)
		}

		errStr = string(bError)
	}

	log := &Log{
		Type:        BackendType,
		SubType:     logSubType,
		App:         l.Default.AppName,
		Environment: l.Default.Environment,
		HostName:    l.Default.Hostname,
		ClientIp:    l.Default.ClientIP,
		Error:       errStr,
		Backend:     &domain.Backend{},
	}

	if ctx != nil {
		if ctx.Request() != nil {
			log.Backend.Request.Method = ctx.Request().Method
			log.Backend.Request.Uri = ctx.FullPath()

			if l.config.Body &&
				!l.config.BodyExcludeUris.Contains(log.Backend.Request.Method, log.Backend.Request.Uri) {
				body := ctx.GetBody()
				if body != nil {
					log.Backend.Request.Body = string(body)
				}
			}
		}

		if ctx.Response() != nil {
			if responseStatusCode == 0 {
				responseStatusCode = ctx.Response().Status()
			}
			log.Backend.Response.StatusCode = responseStatusCode

			if l.config.Body &&
				!l.config.BodyExcludeUris.Contains(log.Backend.Request.Method, log.Backend.Request.Uri) {
				log.Backend.Response.Body = responseBody
			}
		}
	}

	ev.Interface("log", log)
	ev.Msg(msg)
}

// Frontend logs frontend errors
func (l *Logging) Frontend(error string, level errors.Level, fe *domain.Frontend) {
	if len(l.Writers) == 0 {
		//If no output to write, return
		return
	}

	logSubType := domain.DefaultSubType
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var msg string

	log := &Log{
		Type:        FrontendType,
		SubType:     logSubType,
		App:         l.Default.AppName,
		Environment: l.Default.Environment,
		HostName:    l.Default.Hostname,
		ClientIp:    l.Default.ClientIP,
		Error:       error,
		Frontend:    fe,
	}

	ev := getErrorByLevel(l.log, level)
	ev.Interface("log", log)
	ev.Msg(msg)
}

// getErrorEvent gets the error event
func getErrorEvent(log zerolog.Logger, level errors.Level, err interface{}) *zerolog.Event {
	// If level defined in method
	if int(level) > 0 {
		return getErrorByLevel(log, level)
	} else {
		//Get the level by our error codes
		data, ok := err.(errors.ErrorDetails)
		if ok {
			return getErrorByLevel(log, errors.Level(data.Level))
		}
	}

	//default
	return log.Error()
}

// getErrorByLevel gets the error by level
func getErrorByLevel(log zerolog.Logger, level errors.Level) *zerolog.Event {
	switch level {
	case errors.Fatal:
		return log.Fatal()
	case errors.Warning:
		return log.Warn()
	case errors.Info:
		return log.Info()
	case errors.Debug:
		return log.Debug()
	default:
		return log.Error()
	}
}

// getStack gets the debug stack
func getStack() []string {
	stackList := strings.Split(string(debug.Stack()), "\n")
	var k int
	var s string
	for k, s = range stackList {
		if strings.HasPrefix(s, "github.com/gin-gonic/gin.(*Context).Next") {
			break
		}
	}

	return stackList[0:k]
}
