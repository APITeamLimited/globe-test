package ottojs

import (
	"github.com/robertkrimen/otto"
	"time"
)

func (vu *VUContext) Sleep(call otto.FunctionCall) otto.Value ***REMOVED***
	t, err := call.Argument(0).ToFloat()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	time.Sleep(time.Duration(t) * time.Second)
	return otto.UndefinedValue()
***REMOVED***
