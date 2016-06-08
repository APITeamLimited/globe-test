package speedboat

import (
	"golang.org/x/net/context"
)

const (
	AbortTest FlowControl = 0
)

type FlowControl int

func (op FlowControl) Error() string ***REMOVED***
	switch op ***REMOVED***
	case 0:
		return "OP: Abort Test"
	default:
		return "Unknown flow control OP"
	***REMOVED***
***REMOVED***

type Runner interface ***REMOVED***
	RunVU(ctx context.Context, t Test, id int)
***REMOVED***
