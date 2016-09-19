package js

import (
	"github.com/robertkrimen/otto"
	"sync/atomic"
	"time"
)

type JSAPI struct ***REMOVED***
	vu *VU
***REMOVED***

func (a JSAPI) Sleep(secs float64) ***REMOVED***
	time.Sleep(time.Duration(secs * float64(time.Second)))
***REMOVED***

func (a JSAPI) DoGroup(call otto.FunctionCall) otto.Value ***REMOVED***
	name := call.Argument(0).String()
	group, ok := a.vu.group.Group(name, &(a.vu.runner.groupIDCounter))
	if !ok ***REMOVED***
		a.vu.runner.groupsMutex.Lock()
		a.vu.runner.Groups = append(a.vu.runner.Groups, group)
		a.vu.runner.groupsMutex.Unlock()
	***REMOVED***
	a.vu.group = group
	defer func() ***REMOVED*** a.vu.group = group.Parent ***REMOVED***()

	fn := call.Argument(1)
	if !fn.IsFunction() ***REMOVED***
		panic(call.Otto.MakeSyntaxError("fn must be a function"))
	***REMOVED***

	val, err := fn.Call(call.This)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return val
***REMOVED***

func (a JSAPI) DoTest(call otto.FunctionCall) otto.Value ***REMOVED***
	if len(call.ArgumentList) < 2 ***REMOVED***
		return otto.UndefinedValue()
	***REMOVED***

	arg0 := call.Argument(0)
	for _, v := range call.ArgumentList[1:] ***REMOVED***
		obj := v.Object()
		if obj == nil ***REMOVED***
			panic(call.Otto.MakeTypeError("tests must be objects"))
		***REMOVED***
		for _, name := range obj.Keys() ***REMOVED***
			val, err := obj.Get(name)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			result, err := Test(val, arg0)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			test, ok := a.vu.group.Test(name, &(a.vu.runner.testIDCounter))
			if !ok ***REMOVED***
				a.vu.runner.testsMutex.Lock()
				a.vu.runner.Tests = append(a.vu.runner.Tests, test)
				a.vu.runner.testsMutex.Unlock()
			***REMOVED***

			if result ***REMOVED***
				atomic.AddInt64(&(test.Passes), 1)
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&(test.Fails), 1)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return otto.UndefinedValue()
***REMOVED***

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
