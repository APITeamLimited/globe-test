package js

import (
	"context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
)

var ErrDefaultExport = errors.New("you must export a 'default' function")

const entrypoint = "__$$entrypoint$$__"

type Runner struct ***REMOVED***
	Runtime *Runtime
***REMOVED***

func NewRunner(runtime *Runtime, exports otto.Value) (*Runner, error) ***REMOVED***
	expObj := exports.Object()
	if expObj == nil ***REMOVED***
		return nil, ErrDefaultExport
	***REMOVED***

	// Values "remember" which VM they belong to, so to get a callable that works across VM copies,
	// we have to stick it in the global scope, then retrieve it again from the new instance.
	callable, err := expObj.Get("default")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !callable.IsFunction() ***REMOVED***
		return nil, ErrDefaultExport
	***REMOVED***
	if err := runtime.VM.Set(entrypoint, callable); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***Runtime: runtime***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	u := &VU***REMOVED***runner: r, vm: r.Runtime.VM.Copy()***REMOVED***

	callable, err := u.vm.Get(entrypoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	u.callable = callable

	if err := u.vm.Set("__vu_impl__", u); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return u, nil
***REMOVED***

type VU struct ***REMOVED***
	ID int64

	runner   *Runner
	vm       *otto.Otto
	callable otto.Value

	ctx   context.Context
	group *lib.Group
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	u.ctx = ctx
	if _, err := u.callable.Call(otto.UndefinedValue()); err != nil ***REMOVED***
		u.ctx = nil
		return nil, err
	***REMOVED***
	u.ctx = nil
	return nil, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	return nil
***REMOVED***

func (u *VU) DoGroup(call otto.FunctionCall) otto.Value ***REMOVED***
	name := call.Argument(0).String()
	fn := call.Argument(1)
	if !fn.IsFunction() ***REMOVED***
		panic(call.Otto.MakeSyntaxError("fn must be a function"))
	***REMOVED***
	log.WithField("name", name).Info("Group")

	val, err := fn.Call(call.This)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return val
***REMOVED***

func (u *VU) DoTest(call otto.FunctionCall) otto.Value ***REMOVED***
	if len(call.ArgumentList) < 2 ***REMOVED***
		return otto.UndefinedValue()
	***REMOVED***

	arg0 := call.Argument(0)
	for _, v := range call.ArgumentList[1:] ***REMOVED***
		obj := v.Object()
		if obj == nil ***REMOVED***
			panic(call.Otto.MakeTypeError("tests must be objects"))
		***REMOVED***
		for _, key := range obj.Keys() ***REMOVED***
			val, err := obj.Get(key)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			var res bool

		typeSwitchLoop:
			for ***REMOVED***
				switch ***REMOVED***
				case val.IsFunction():
					val, err = val.Call(otto.UndefinedValue(), arg0)
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					continue typeSwitchLoop
				case val.IsUndefined() || val.IsNull():
					res = false
				case val.IsBoolean():
					b, err := val.ToBoolean()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					res = b
				case val.IsNumber():
					f, err := val.ToFloat()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					res = (f != 0)
				case val.IsString():
					s, err := val.ToString()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					res = (s != "")
				***REMOVED***
				break
			***REMOVED***

			log.WithFields(log.Fields***REMOVED***
				"arg0": arg0,
				"key":  key,
				"res":  res,
			***REMOVED***).Info("Test")
		***REMOVED***
	***REMOVED***
	return otto.UndefinedValue()
***REMOVED***
