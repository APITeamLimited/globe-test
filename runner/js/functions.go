package js

import (
	"github.com/robertkrimen/otto"
	"time"
)

func jsSleepFactory(impl func(time.Duration)) func(otto.FunctionCall) otto.Value ***REMOVED***
	return func(call otto.FunctionCall) otto.Value ***REMOVED***
		seconds, err := call.Argument(0).ToFloat()
		if err != nil ***REMOVED***
			seconds = 0.0
		***REMOVED***
		impl(time.Duration(seconds * float64(time.Second)))
		return otto.UndefinedValue()
	***REMOVED***
***REMOVED***
