package simple

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/rcrowley/go-metrics"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"time"
)

type Runner struct ***REMOVED***
	Client *fasthttp.Client
***REMOVED***

func New() *Runner ***REMOVED***
	return &Runner***REMOVED***
		Client: &fasthttp.Client***REMOVED***
			MaxIdleConnDuration: time.Duration(0),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	mDuration := metrics.NewRegisteredHistogram("duration", speedboat.Registry, metrics.NewUniformSample(1024))
	for ***REMOVED***
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(t.URL)

		startTime := time.Now()
		if err := r.Client.Do(req, res); err != nil ***REMOVED***
			log.WithError(err).Error("Request error")
		***REMOVED***
		duration := time.Since(startTime)

		mDuration.Update(duration.Nanoseconds())

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***
