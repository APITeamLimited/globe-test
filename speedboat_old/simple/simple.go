package simple

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/proto/httpwrap"
	"github.com/loadimpact/speedboat/stats"
	"math"
	"net/http"
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
	Client    http.Client
	Request   http.Request
	Collector *stats.Collector
***REMOVED***

func New(url string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		URL: url,
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	req, err := http.NewRequest("GET", r.URL, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &VU***REMOVED***
		Runner: r,
		Client: http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				MaxIdleConnsPerHost: math.MaxInt32,
			***REMOVED***,
		***REMOVED***,
		Request:   *req,
		Collector: stats.NewCollector(),
	***REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	req := u.Request
	_, _, sample, err := httpwrap.Do(ctx, &u.Client, &req, httpwrap.Params***REMOVED***TakeSample: true***REMOVED***)
	if err != nil ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		default:
			log.WithError(err).Error("Request Error")
			u.Collector.Add(stats.Sample***REMOVED***
				Stat: &mErrors,
				Tags: stats.Tags***REMOVED***
					"method": req.Method,
					"url":    req.URL.String(),
					"error":  err.Error(),
				***REMOVED***,
				Values: stats.Value(1),
			***REMOVED***)
			return err
		***REMOVED***
	***REMOVED***

	sample.Stat = &mRequests
	u.Collector.Add(sample)

	return nil
***REMOVED***
