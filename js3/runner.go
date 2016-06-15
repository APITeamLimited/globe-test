package js3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
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

	vm := duktape.New()

	vm.PushGlobalGoFunction("sleep", func(vm *duktape.Context) int ***REMOVED***
		t := vm.ToInt(0)
		time.Sleep(time.Duration(t) * time.Second)
		return 0
	***REMOVED***)
	vm.PushGlobalGoFunction("get", func(vm *duktape.Context) int ***REMOVED***
		url := vm.ToString(0)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		startTime := time.Now()
		err := r.Client.Do(req, res)
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

		return 0
	***REMOVED***)

	vm.PushString(r.Filename)
	if err := vm.PcompileStringFilename(0, r.Source); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't compile script")
		return
	***REMOVED***

	for ***REMOVED***
		vm.DupTop()
		if vm.Pcall(0) != duktape.ErrNone ***REMOVED***
			err := vm.SafeToString(-1)
			log.WithField("error", err).Error("Script error")
		***REMOVED***
		vm.Pop()
	***REMOVED***
***REMOVED***
