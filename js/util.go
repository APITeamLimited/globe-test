package js

import (
	"github.com/robertkrimen/otto"
)

func Test(val, arg0 otto.Value) (bool, error) ***REMOVED***
	switch ***REMOVED***
	case val.IsFunction():
		val, err := val.Call(otto.UndefinedValue(), arg0)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return Test(val, arg0)
	case val.IsBoolean():
		b, err := val.ToBoolean()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return b, nil
	case val.IsNumber():
		f, err := val.ToFloat()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return f != 0, nil
	case val.IsString():
		s, err := val.ToString()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return s != "", nil
	default:
		return false, nil
	***REMOVED***
***REMOVED***

func throw(vm *otto.Otto, v interface***REMOVED******REMOVED***) ***REMOVED***
	if err, ok := v.(error); ok ***REMOVED***
		panic(vm.MakeCustomError("Error", err.Error()))
	***REMOVED***
	panic(v)
***REMOVED***
