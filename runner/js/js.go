package js

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/util"
	"github.com/robertkrimen/otto"
	"net/http"
	"time"
)

type JSRunner struct ***REMOVED***
	BaseVM *otto.Otto
	Script *otto.Script

	httpClient *http.Client
***REMOVED***

func New() (r *JSRunner, err error) ***REMOVED***
	r = &JSRunner***REMOVED******REMOVED***

	// Create a base VM
	r.BaseVM = otto.New()

	// Bridge basic functions
	r.BaseVM.Set("sleep", jsSleepFactory(time.Sleep))
	r.BaseVM.Set("log", jsLogFactory(func(text string) ***REMOVED***
		// out <- runner.NewLogEntry(text)
		log.WithField("text", text).Info("Test Log")
	***REMOVED***))

	// Use a single HTTP client for this
	r.httpClient = &http.Client***REMOVED***
		Transport: &http.Transport***REMOVED***
			DisableKeepAlives: true,
		***REMOVED***,
	***REMOVED***
	r.BaseVM.Set("get", jsHTTPGetFactory(func(url string) (*http.Response, error) ***REMOVED***
		return r.httpClient.Get(url)
	***REMOVED***))

	return r, nil
***REMOVED***

func (r *JSRunner) Load(filename, src string) (err error) ***REMOVED***
	r.Script, err = r.BaseVM.Compile(filename, src)
	return err
***REMOVED***

func (r *JSRunner) RunVU(stop <-chan interface***REMOVED******REMOVED***) <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	out := make(chan interface***REMOVED******REMOVED***)

	go func() ***REMOVED***
		defer close(out)

	runLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-stop:
				break runLoop
			default:
				for res := range r.RunIteration() ***REMOVED***
					out <- res
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

func (r *JSRunner) RunIteration() <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	out := make(chan interface***REMOVED******REMOVED***)

	go func() ***REMOVED***
		defer close(out)
		defer func() ***REMOVED***
			if err := recover(); err != nil ***REMOVED***
				out <- runner.NewError(err.(JSError))
			***REMOVED***
		***REMOVED***()

		// Make a copy of the base VM
		vm := r.BaseVM //.Copy()

		// Log has to be bridged here, as it needs a reference to the channel
		// vm.Set("log", jsLogFactory(func(text string) ***REMOVED***
		// 	out <- runner.NewLogEntry(text)
		// ***REMOVED***))

		startTime := time.Now()
		var err error
		duration := util.Time(func() ***REMOVED***
			_, err = vm.Run(r.Script)
		***REMOVED***)

		if err != nil ***REMOVED***
			out <- runner.NewError(err)
		***REMOVED***

		out <- runner.NewMetric(startTime, duration)
	***REMOVED***()

	return out
***REMOVED***
