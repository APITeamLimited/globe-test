package js

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/robertkrimen/otto"
	"time"
)

type JSRunner struct ***REMOVED***
	BaseVM *otto.Otto
	Script *otto.Script
***REMOVED***

func New() (r *JSRunner, err error) ***REMOVED***
	r = &JSRunner***REMOVED******REMOVED***

	// Create a base VM
	r.BaseVM = otto.New()

	// Bridge basic functions
	r.BaseVM.Set("sleep", jsSleepFactory(time.Sleep))

	return r, nil
***REMOVED***

func (r *JSRunner) Load(filename, src string) (err error) ***REMOVED***
	r.Script, err = r.BaseVM.Compile(filename, src)
	return err
***REMOVED***

func (r *JSRunner) RunVU() <-chan runner.Result ***REMOVED***
	out := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(out)

		vm := r.BaseVM.Copy()
		for res := range r.RunIteration(vm) ***REMOVED***
			out <- res
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

func (r *JSRunner) RunIteration(vm *otto.Otto) <-chan runner.Result ***REMOVED***
	out := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(out)

		startTime := time.Now()
		vm.Run(r.Script)
		duration := time.Since(startTime)

		out <- runner.Result***REMOVED***
			Type: "metric",
			Metric: runner.Metric***REMOVED***
				Time:     time.Now(),
				Duration: duration,
			***REMOVED***,
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***
