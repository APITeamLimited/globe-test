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
	MetricVUs       = &stats.Metric***REMOVED***Name: "vus", Type: stats.Gauge***REMOVED***
	MetricVUsPooled = &stats.Metric***REMOVED***Name: "vus_pooled", Type: stats.Gauge***REMOVED***
	MetricErrors    = &stats.Metric***REMOVED***Name: "errors", Type: stats.Counter***REMOVED***
)

type Engine struct ***REMOVED***
	Runner  Runner
	Status  Status
	Metrics map[*stats.Metric][]stats.Sample

	ctx       context.Context
	cancelers []context.CancelFunc
	pool      []VU

	vuMutex sync.Mutex
	mMutex  sync.Mutex
***REMOVED***

func NewEngine(r Runner) *Engine ***REMOVED***
	return &Engine***REMOVED***
		Runner:  r,
		Metrics: make(map[*stats.Metric][]stats.Sample),
	***REMOVED***
***REMOVED***

func (e *Engine) Run(ctx context.Context, prepared int64) error ***REMOVED***
	e.ctx = ctx

	e.pool = make([]VU, prepared)
	for i := int64(0); i < prepared; i++ ***REMOVED***
		vu, err := e.Runner.NewVU()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		e.pool[i] = vu
	***REMOVED***

	e.Status.StartTime = time.Now()
	e.Status.Running = true
	e.Status.VUs = 0
	e.Status.Pooled = prepared

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
	e.Status.VUs = 0
	e.Status.Pooled = 0

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

	e.Status.VUs = int64(len(e.cancelers))
	e.Status.Pooled = int64(len(e.pool))

	return nil
***REMOVED***

func (e *Engine) reportInternalStats() ***REMOVED***
	e.mMutex.Lock()
	t := time.Now()
	e.Metrics[MetricVUs] = append(
		e.Metrics[MetricVUs],
		stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(len(e.cancelers))***REMOVED***,
	)
	e.Metrics[MetricVUsPooled] = append(
		e.Metrics[MetricVUsPooled],
		stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(len(e.pool))***REMOVED***,
	)
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
				e.Metrics[MetricErrors] = append(e.Metrics[MetricErrors], stats.Sample***REMOVED***
					Time:  time.Now(),
					Tags:  map[string]string***REMOVED***"vu": idString, "error": err.Error()***REMOVED***,
					Value: float64(1),
				***REMOVED***)
			***REMOVED***
			for _, s := range samples ***REMOVED***
				e.Metrics[s.Metric] = append(e.Metrics[s.Metric], s.Sample)
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
