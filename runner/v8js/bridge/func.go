package bridge

import (
	"errors"
	"fmt"
	"reflect"
)

type Func struct ***REMOVED***
	Func      reflect.Value
	In, Out   []Type
	IsVaradic bool
	VarArg    Type
***REMOVED***

func (f *Func) Call(args []interface***REMOVED******REMOVED***) error ***REMOVED***
	rArgs := make([]reflect.Value, 0, len(args))
	for i, v := range args ***REMOVED***
		t := Type***REMOVED******REMOVED***
		if i >= len(f.In) ***REMOVED***
			if f.IsVaradic ***REMOVED***
				t = f.VarArg
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			t = f.In[i]
		***REMOVED***

		if err := t.Cast(&v); err != nil ***REMOVED***
			return err
		***REMOVED***
		rArgs = append(rArgs, reflect.ValueOf(v))
	***REMOVED***
	f.Func.Call(rArgs)
	return nil
***REMOVED***

func (f *Func) JS(mod, name string) string ***REMOVED***
	return fmt.Sprintf(`function() ***REMOVED*** __internal__._invoke('%s', '%s', Array.prototype.slice.call(arguments)); ***REMOVED***`, mod, name)
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
		if !fnT.IsVariadic() || i != fnT.NumIn()-1 ***REMOVED***
			fn.In = append(fn.In, BridgeType(fnT.In(i)))
		***REMOVED*** else ***REMOVED***
			fn.IsVaradic = true
			fn.VarArg = BridgeType(fnT.In(i).Elem())
		***REMOVED***
	***REMOVED***
	for i := 0; i < fnT.NumOut(); i++ ***REMOVED***
		fn.Out = append(fn.Out, BridgeType(fnT.Out(i)))
	***REMOVED***

	return fn
***REMOVED***
