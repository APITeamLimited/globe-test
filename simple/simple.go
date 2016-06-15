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
***REMOVED***

func New() *Runner ***REMOVED***
	return &Runner***REMOVED******REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	mDuration := sampler.Stats("request.duration")
	mErrors := sampler.Counter("request.error")

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	client := fasthttp.Client***REMOVED***MaxConnsPerHost: math.MaxInt32***REMOVED***

	for ***REMOVED***
		req.Reset()
		req.SetRequestURI(t.URL)
		res.Reset()

		startTime := time.Now()
		if err := client.Do(req, res); err != nil ***REMOVED***
			log.WithError(err).Error("Request error")
			mErrors.WithField("url", t.URL).Int(1)
		***REMOVED***
		duration := time.Since(startTime)

		mDuration.WithField("url", t.URL).Duration(duration)

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***
