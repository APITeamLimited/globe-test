package js

import (
	"context"
	"errors"
	// log "github.com/Sirupsen/logrus"
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
	vm := r.Runtime.VM.Copy()
	callable, err := vm.Get(entrypoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &VU***REMOVED***runner: r, vm: vm, callable: callable***REMOVED***, nil
***REMOVED***

type VU struct ***REMOVED***
	ID int64

	runner   *Runner
	vm       *otto.Otto
	callable otto.Value

	ctx context.Context
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
