package js

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/robertkrimen/otto"
	"time"
)

type JSRunner struct ***REMOVED***
	BaseVM *otto.Otto
***REMOVED***

func New() (r *JSRunner, err error) ***REMOVED***
	r = &JSRunner***REMOVED******REMOVED***

	// Create a base VM
	r.BaseVM = otto.New()

	// TODO: Bridge functions here

	return r, nil
***REMOVED***

func (r *JSRunner) Run(filename, src string) <-chan runner.Result ***REMOVED***
	out := make(chan runner.Result)

	go func() ***REMOVED***
		out <- runner.Result***REMOVED***
			Type: "log",
			LogEntry: runner.LogEntry***REMOVED***
				Time: time.Now(),
				Text: "AaaaaaA",
			***REMOVED***,
		***REMOVED***
		close(out)
	***REMOVED***()

	return out
***REMOVED***
