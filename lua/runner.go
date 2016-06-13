package lua

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
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
		Client:   &fasthttp.Client***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	mDuration := sampler.Stats("request.duration")
	mErrors := sampler.Counter("request.error")

	L := lua.NewState()
	defer L.Close()

	fn, err := L.LoadString(r.Source)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't compile script")
		return
	***REMOVED***

	L.SetGlobal("sleep", L.NewFunction(func(L *lua.LState) int ***REMOVED***
		time.Sleep(time.Duration(L.ToInt64(1)) * time.Second)
		return 0
	***REMOVED***))
	L.SetGlobal("get", L.NewFunction(func(L *lua.LState) int ***REMOVED***
		url := L.ToString(1)

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
	***REMOVED***))

	for ***REMOVED***
		L.Push(fn)
		if err := L.PCall(0, 0, nil); err != nil ***REMOVED***
			log.WithError(err).Error("Script error")
		***REMOVED***
	***REMOVED***
***REMOVED***
