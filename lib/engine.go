package lib

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"strconv"
	"sync"
	"time"
)

var (
	MetricActiveVUs   = &stats.Metric***REMOVED***Name: "vus_active", Type: stats.Gauge***REMOVED***
	MetricInactiveVUs = &stats.Metric***REMOVED***Name: "vus_inactive", Type: stats.Gauge***REMOVED***
	MetricErrors      = &stats.Metric***REMOVED***Name: "errors", Type: stats.Counter***REMOVED***
)

type Engine struct ***REMOVED***
	Runner  Runner
	Status  Status
	Metrics map[*stats.Metric]stats.Sink

	ctx       context.Context
	cancelers []context.CancelFunc
	pool      []VU

	vuMutex sync.Mutex
	mMutex  sync.Mutex
***REMOVED***

func NewEngine(r Runner, prepared int64) (*Engine, error) ***REMOVED***
	pool := make([]VU, prepared)
	for i := int64(0); i < prepared; i++ ***REMOVED***
		vu, err := r.NewVU()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pool[i] = vu
	***REMOVED***

	return &Engine***REMOVED***
		Runner:  r,
		Metrics: make(map[*stats.Metric]stats.Sink),
		pool:    pool,
	***REMOVED***, nil
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.ctx = ctx

	e.Status.ID = "default"
	e.Status.Running = true
	e.Status.ActiveVUs = int64(len(e.cancelers))
	e.Status.InactiveVUs = int64(len(e.pool))

	e.reportInternalStats()
	ticker := time.NewTicker(1 * time.Second)

loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.reportInternalStats()
		case <-ctx.Done():
			break loop
		***REMOVED***
	***REMOVED***

	e.cancelers = nil
	e.pool = nil

	e.Status.Running = false
	e.Status.ActiveVUs = 0
	e.Status.InactiveVUs = 0
	e.reportInternalStats()

	return nil
***REMOVED***

func (e *Engine) Scale(vus int64) error ***REMOVED***
	e.vuMutex.Lock()
	defer e.vuMutex.Unlock()

	l := int64(len(e.cancelers))
	switch ***REMOVED***
	case l < vus:
		for i := int64(len(e.cancelers)); i < vus; i++ ***REMOVED***
			vu, err := e.getVU()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			id := i + 1
			if err := vu.Reconfigure(id); err != nil ***REMOVED***
				return err
			***REMOVED***

			ctx, cancel := context.WithCancel(e.ctx)
			e.cancelers = append(e.cancelers, cancel)
			go func() ***REMOVED***
				e.runVU(ctx, id, vu)

				e.vuMutex.Lock()
				e.pool = append(e.pool, vu)
				e.vuMutex.Unlock()
			***REMOVED***()
		***REMOVED***
	case l > vus:
		for _, cancel := range e.cancelers[vus+1:] ***REMOVED***
			cancel()
		***REMOVED***
		e.cancelers = e.cancelers[:vus]
	***REMOVED***

	e.Status.ActiveVUs = int64(len(e.cancelers))
	e.Status.InactiveVUs = int64(len(e.pool))

	return nil
***REMOVED***

func (e *Engine) reportInternalStats() ***REMOVED***
	e.mMutex.Lock()
	t := time.Now()
	e.getSink(MetricActiveVUs).Add(stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(len(e.cancelers))***REMOVED***)
	e.getSink(MetricInactiveVUs).Add(stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(len(e.pool))***REMOVED***)
	e.mMutex.Unlock()
***REMOVED***

func (e *Engine) runVU(ctx context.Context, id int64, vu VU) ***REMOVED***
	idString := strconv.FormatInt(id, 10)
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
			samples, err := vu.RunOnce(ctx)
			e.mMutex.Lock()
			if err != nil ***REMOVED***
				log.WithField("vu", id).WithError(err).Error("Runtime Error")
				e.getSink(MetricErrors).Add(stats.Sample***REMOVED***
					Time:  time.Now(),
					Tags:  map[string]string***REMOVED***"vu": idString, "error": err.Error()***REMOVED***,
					Value: float64(1),
				***REMOVED***)
			***REMOVED***
			for _, s := range samples ***REMOVED***
				e.getSink(s.Metric).Add(s)
			***REMOVED***
			e.mMutex.Unlock()
		***REMOVED***
	***REMOVED***
***REMOVED***

// Returns a pooled VU if available, otherwise make a new one.
func (e *Engine) getVU() (VU, error) ***REMOVED***
	l := len(e.pool)
	if l > 0 ***REMOVED***
		vu := e.pool[l-1]
		e.pool = e.pool[:l-1]
		return vu, nil
	***REMOVED***

	log.Warn("More VUs requested than what was prepared; instantiation during tests is costly and may skew results!")
	return e.Runner.NewVU()
***REMOVED***

// Returns a value sink for a metric, created from the type if unavailable.
func (e *Engine) getSink(m *stats.Metric) stats.Sink ***REMOVED***
	s, ok := e.Metrics[m]
	if !ok ***REMOVED***
		switch m.Type ***REMOVED***
		case stats.Counter:
			s = &stats.CounterSink***REMOVED******REMOVED***
		case stats.Gauge:
			s = &stats.GaugeSink***REMOVED******REMOVED***
		case stats.Trend:
			s = &stats.TrendSink***REMOVED******REMOVED***
		***REMOVED***
		e.Metrics[m] = s
	***REMOVED***
	return s
***REMOVED***
