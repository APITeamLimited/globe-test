package simple

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"time"
)

type Runner struct ***REMOVED***
	Test speedboat.Test

	mDuration *sampler.Metric
	mErrors   *sampler.Metric
***REMOVED***

type VU struct ***REMOVED***
	Runner  *Runner
	Client  fasthttp.Client
	Request fasthttp.Request
***REMOVED***

func New(t speedboat.Test) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Test:      t,
		mDuration: sampler.Stats("request.duration"),
		mErrors:   sampler.Counter("request.error"),
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (speedboat.VU, error) ***REMOVED***
	vu := &VU***REMOVED***
		Runner: r,
		Client: fasthttp.Client***REMOVED***MaxConnsPerHost: math.MaxInt32***REMOVED***,
	***REMOVED***

	vu.Request.SetRequestURI(r.Test.URL)

	return vu, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	startTime := time.Now()
	err := u.Client.Do(&u.Request, res)
	duration := time.Since(startTime)

	u.Runner.mDuration.WithFields(sampler.Fields***REMOVED***
		"url":    u.Runner.Test.URL,
		"method": "GET",
		"status": res.StatusCode(),
	***REMOVED***).Duration(duration)

	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
		u.Runner.mErrors.WithFields(sampler.Fields***REMOVED***
			"url":    u.Runner.Test.URL,
			"method": "GET",
			"status": res.StatusCode(),
		***REMOVED***).Int(1)
		return err
	***REMOVED***

	return nil
***REMOVED***
