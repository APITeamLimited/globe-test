package js2

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"time"
)

type Runner struct ***REMOVED***
	Filename string
	Source   string
	Client   *fasthttp.Client
***REMOVED***

func New(filename, src string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Filename: filename,
		Source:   src,
		Client: &fasthttp.Client***REMOVED***
			MaxConnsPerHost: math.MaxInt32,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	mDuration := sampler.Stats("request.duration")
	mErrors := sampler.Counter("request.error")

	vm := otto.New()
	script, err := vm.Compile(r.Filename, r.Source)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't compile script")
		return
	***REMOVED***

	vm.Set("sleep", func(call otto.FunctionCall) otto.Value ***REMOVED***
		t, err := call.Argument(0).ToInteger()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		time.Sleep(time.Duration(t) * time.Second)
		return otto.UndefinedValue()
	***REMOVED***)
	vm.Set("get", func(call otto.FunctionCall) otto.Value ***REMOVED***
		url, err := call.Argument(0).ToString()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		startTime := time.Now()
		err = r.Client.Do(req, res)
		duration := time.Since(startTime)

		mDuration.WithFields(sampler.Fields***REMOVED***
			"url":    url,
			"method": "GET",
			"status": res.StatusCode(),
		***REMOVED***).Duration(duration)

		if err != nil ***REMOVED***
			log.WithError(err).Error("Request error")
			mErrors.WithFields(sampler.Fields***REMOVED***
				"url":    url,
				"method": "GET",
				"error":  err,
			***REMOVED***).Int(1)
		***REMOVED***

		return otto.UndefinedValue()
	***REMOVED***)

	for ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			log.WithError(err).Error("Script error")
		***REMOVED***
	***REMOVED***
***REMOVED***
