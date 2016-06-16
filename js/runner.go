package js

import (
	"fmt"
	"github.com/loadimpact/speedboat"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"os"
)

type Runner struct ***REMOVED***
	Test speedboat.Test

	filename string
	source   string
***REMOVED***

type VU struct ***REMOVED***
	Runner *Runner
	VM     *otto.Otto
	Script *otto.Script

	Client fasthttp.Client

	ID int64
***REMOVED***

func New(t speedboat.Test, filename, source string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Test:     t,
		filename: filename,
		source:   source,
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (speedboat.VU, error) ***REMOVED***
	vm := otto.New()

	script, err := vm.Compile(r.filename, r.source)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vm.Set("print", func(call otto.FunctionCall) otto.Value ***REMOVED***
		fmt.Fprintln(os.Stderr, call.Argument(0))
		return otto.UndefinedValue()
	***REMOVED***)

	return &VU***REMOVED***
		Runner: r,
		VM:     vm,
		Script: script,
	***REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	if _, err := u.VM.Run(u.Script); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
