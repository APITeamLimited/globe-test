package js

import (
	log "github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"
	"strconv"
	"sync/atomic"
	"time"
)

type JSAPI struct ***REMOVED***
	vu *VU
***REMOVED***

func (a JSAPI) Sleep(secs float64) ***REMOVED***
	time.Sleep(time.Duration(secs * float64(time.Second)))
***REMOVED***

func (a JSAPI) Log(level int, msg string, args []otto.Value) ***REMOVED***
	fields := make(log.Fields, len(args))
	for i, arg := range args ***REMOVED***
		if arg.IsObject() ***REMOVED***
			obj := arg.Object()
			for _, key := range obj.Keys() ***REMOVED***
				v, err := obj.Get(key)
				if err != nil ***REMOVED***
					throw(a.vu.vm, err)
				***REMOVED***
				fields[key] = v.String()
			***REMOVED***
			continue
		***REMOVED***
		fields["arg"+strconv.FormatInt(int64(i), 10)] = arg.String()
	***REMOVED***

	entry := log.WithFields(fields)
	switch level ***REMOVED***
	case 0:
		entry.Debug(msg)
	case 1:
		entry.Info(msg)
	case 2:
		entry.Warn(msg)
	case 3:
		entry.Error(msg)
	***REMOVED***
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
		throw(call.Otto, err)
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
				throw(call.Otto, err)
			***REMOVED***

			result, err := Test(val, arg0)
			if err != nil ***REMOVED***
				throw(call.Otto, err)
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
