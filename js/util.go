package js

import (
	"context"
	"github.com/robertkrimen/otto"
)

func Check(val, arg0 otto.Value) (bool, error) ***REMOVED***
	switch ***REMOVED***
	case val.IsFunction():
		val, err := val.Call(otto.UndefinedValue(), arg0)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return Check(val, arg0)
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

func newSnippetRunner(src string) (*Runner, error) ***REMOVED***
	rt, err := New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	rt.VM.Set("__initapi__", InitAPI***REMOVED***r: rt***REMOVED***)
	defer rt.VM.Set("__initapi__", nil)

	exp, err := rt.load("__snippet__", []byte(src))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewRunner(rt, exp)
***REMOVED***

func runSnippet(src string) error ***REMOVED***
	r, err := newSnippetRunner(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	vu, err := r.NewVU()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = vu.RunOnce(context.Background())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
