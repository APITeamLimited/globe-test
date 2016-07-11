package simple

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***
)

type Runner struct ***REMOVED***
	URL string
***REMOVED***

type VU struct ***REMOVED***
	Runner    *Runner
	Client    fasthttp.Client
	Request   fasthttp.Request
	Collector *stats.Collector
***REMOVED***

func New(url string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		URL: url,
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	vu := &VU***REMOVED***
		Runner:    r,
		Client:    fasthttp.Client***REMOVED***MaxConnsPerHost: math.MaxInt32***REMOVED***,
		Collector: stats.NewCollector(),
	***REMOVED***

	vu.Request.SetRequestURI(r.URL)

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

	tags := stats.Tags***REMOVED***
		"url":    u.Runner.URL,
		"method": "GET",
		"status": res.StatusCode(),
	***REMOVED***
	u.Collector.Add(stats.Sample***REMOVED***
		Stat:   &mRequests,
		Tags:   tags,
		Values: stats.Values***REMOVED***"duration": float64(duration)***REMOVED***,
	***REMOVED***)

	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
		u.Collector.Add(stats.Sample***REMOVED***
			Stat:   &mErrors,
			Tags:   tags,
			Values: stats.Value(1),
		***REMOVED***)
		return err
	***REMOVED***

	return nil
***REMOVED***
