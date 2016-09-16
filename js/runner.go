package js

import (
	"context"
	"errors"
	// log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"sync"
	"sync/atomic"
)

var ErrDefaultExport = errors.New("you must export a 'default' function")

const entrypoint = "__$$entrypoint$$__"

type Runner struct ***REMOVED***
	Runtime      *Runtime
	DefaultGroup *lib.Group
	Groups       []*lib.Group

	groupIDCounter int64
	testIDCounter  int64
	groupsMutex    sync.Mutex
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

	r := &Runner***REMOVED***Runtime: runtime***REMOVED***
	r.DefaultGroup = lib.NewGroup("", nil, nil)
	r.Groups = []*lib.Group***REMOVED***r.DefaultGroup***REMOVED***

	return r, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	u := &VU***REMOVED***
		runner: r,
		vm:     r.Runtime.VM.Copy(),
		group:  r.DefaultGroup,
	***REMOVED***

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

func (r *Runner) GetGroups() []*lib.Group ***REMOVED***
	return r.Groups
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
	group, ok := u.group.Group(name, &(u.runner.groupIDCounter))
	if !ok ***REMOVED***
		u.runner.groupsMutex.Lock()
		u.runner.Groups = append(u.runner.Groups, group)
		u.runner.groupsMutex.Unlock()
	***REMOVED***
	u.group = group
	defer func() ***REMOVED*** u.group = group.Parent ***REMOVED***()

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
		for _, name := range obj.Keys() ***REMOVED***
			val, err := obj.Get(name)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			var result bool

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
					result = false
				case val.IsBoolean():
					b, err := val.ToBoolean()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					result = b
				case val.IsNumber():
					f, err := val.ToFloat()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					result = (f != 0)
				case val.IsString():
					s, err := val.ToString()
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					result = (s != "")
				***REMOVED***
				break
			***REMOVED***

			test, _ := u.group.Test(name, &(u.runner.testIDCounter))
			if result ***REMOVED***
				atomic.AddInt64(&(test.Passes), 1)
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&(test.Fails), 1)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return otto.UndefinedValue()
***REMOVED***
