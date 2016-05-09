package bridge

import (
	"errors"
	"reflect"
)

type Func struct ***REMOVED***
	Func    reflect.Value
	In, Out []Type
***REMOVED***

// Creates a bridged function.
// Panics if raw is not a function; this is a blatant programming error.
func BridgeFunc(raw interface***REMOVED******REMOVED***) Func ***REMOVED***
	fn := Func***REMOVED***Func: reflect.ValueOf(raw)***REMOVED***
	fnT := fn.Func.Type()

	// We can only bridge functions
	if fn.Func.Kind() != reflect.Func ***REMOVED***
		panic(errors.New("That's not a function >_>"))
	***REMOVED***

	for i := 0; i < fnT.NumIn(); i++ ***REMOVED***
		fn.In = append(fn.In, BridgeType(fnT.In(i)))
	***REMOVED***
	for i := 0; i < fnT.NumOut(); i++ ***REMOVED***
		fn.Out = append(fn.Out, BridgeType(fnT.Out(i)))
	***REMOVED***

	return fn
***REMOVED***
