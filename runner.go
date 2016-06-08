package speedboat

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	AbortTest FlowControl = 0

	LoggerKey ContextKey = 0
)

type FlowControl int

type ContextKey int

func (op FlowControl) Error() string ***REMOVED***
	switch op ***REMOVED***
	case 0:
		return "OP: Abort Test"
	default:
		return "Unknown flow control OP"
	***REMOVED***
***REMOVED***

func WithLogger(ctx context.Context, logger *log.Logger) context.Context ***REMOVED***
	return context.WithValue(ctx, LoggerKey, logger)
***REMOVED***

func GetLogger(ctx context.Context) *log.Logger ***REMOVED***
	return ctx.Value(LoggerKey).(*log.Logger)
***REMOVED***

type Runner interface ***REMOVED***
	RunVU(ctx context.Context, t Test, id int)
***REMOVED***
