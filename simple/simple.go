package simple

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math/rand"
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
	mDuration := sampler.Stats("request.duration")
	mErrors := sampler.Counter("request.error")

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	for ***REMOVED***
		req.Reset()
		req.SetRequestURI(t.URL)
		res.Reset()

		startTime := time.Now()
		if err := r.Client.Do(req, res); err != nil ***REMOVED***
			log.WithError(err).Error("Request error")
			mErrors.WithField("url", t.URL).Int(1)
		***REMOVED***
		duration := time.Since(startTime)

		mDuration.WithField("url", t.URL).Duration(duration)

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
		***REMOVED***
	***REMOVED***
***REMOVED***
