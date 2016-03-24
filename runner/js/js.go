package js

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/robertkrimen/otto"
)

type JSRunner struct ***REMOVED***
	BaseVM *otto.Otto
	Script *otto.Script
***REMOVED***

func New() (r *JSRunner, err error) ***REMOVED***
	r = &JSRunner***REMOVED******REMOVED***

	// Create a base VM
	r.BaseVM = otto.New()

	// TODO: Bridge functions here

	return r, nil
***REMOVED***

func (r *JSRunner) Load(filename, src string) (err error) ***REMOVED***
	r.Script, err = r.BaseVM.Compile(filename, src)
	return err
***REMOVED***

func (r *JSRunner) RunVU() <-chan runner.Result ***REMOVED***
	out := make(chan runner.Result)

	go func() ***REMOVED***
		// out <- runner.Result***REMOVED***
		// 	Type: "log",
		// 	LogEntry: runner.LogEntry***REMOVED***
		// 		Time: time.Now(),
		// 		Text: "AaaaaaA",
		// 	***REMOVED***,
		// ***REMOVED***
		vm := r.BaseVM.Copy()
		vm.Run(r.Script)
		close(out)
	***REMOVED***()

	return out
***REMOVED***
