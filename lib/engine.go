package lib

import (
	"context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"sync"
	"time"
)

var (
	MetricVUs    = &stats.Metric***REMOVED***Name: "vus", Type: stats.Gauge***REMOVED***
	MetricVUsMax = &stats.Metric***REMOVED***Name: "vus_max", Type: stats.Gauge***REMOVED***
	MetricErrors = &stats.Metric***REMOVED***Name: "errors", Type: stats.Counter***REMOVED***

	ErrTooManyVUs = errors.New("More VUs than the maximum requested")
	ErrMaxTooLow  = errors.New("Can't lower max below current VU count")
)

type vuEntry struct ***REMOVED***
	VU     VU
	Cancel context.CancelFunc
***REMOVED***

type Engine struct ***REMOVED***
	Runner  Runner
	Status  Status
	Metrics map[*stats.Metric]stats.Sink
	Pause   sync.WaitGroup

	ctx    context.Context
	vus    []*vuEntry
	nextID int64

	vuMutex sync.Mutex
	mMutex  sync.Mutex
***REMOVED***

func NewEngine(r Runner) (*Engine, error) ***REMOVED***
	e := &Engine***REMOVED***
		Runner: r,
		Status: Status***REMOVED***
			Running: null.BoolFrom(false),
			VUs:     null.IntFrom(0),
			VUsMax:  null.IntFrom(0),
		***REMOVED***,
		Metrics: make(map[*stats.Metric]stats.Sink),
	***REMOVED***

	e.Status.Running = null.BoolFrom(false)
	e.Pause.Add(1)

	e.Status.VUs = null.IntFrom(0)
	e.Status.VUsMax = null.IntFrom(0)

	return e, nil
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.ctx = ctx
	e.nextID = 1

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

	e.vus = nil

	e.Status.Running = null.BoolFrom(false)
	e.Status.VUs = null.IntFrom(0)
	e.Status.VUsMax = null.IntFrom(0)
	e.reportInternalStats()

	return nil
***REMOVED***

func (e *Engine) SetRunning(running bool) ***REMOVED***
	if running && !e.Status.Running.Bool ***REMOVED***
		e.Pause.Done()
		log.Debug("Engine Unpaused")
	***REMOVED*** else if !running && e.Status.Running.Bool ***REMOVED***
		e.Pause.Add(1)
		log.Debug("Engine Paused")
	***REMOVED***
	e.Status.Running.Bool = running
***REMOVED***

func (e *Engine) SetVUs(v int64) error ***REMOVED***
	e.vuMutex.Lock()
	defer e.vuMutex.Unlock()

	if v > e.Status.VUsMax.Int64 ***REMOVED***
		return ErrTooManyVUs
	***REMOVED***

	current := e.Status.VUs.Int64
	for i := current; i < v; i++ ***REMOVED***
		entry := e.vus[i]
		if entry.Cancel != nil ***REMOVED***
			panic(errors.New("ATTEMPTED TO RESCHEDULE RUNNING VU"))
		***REMOVED***

		ctx, cancel := context.WithCancel(e.ctx)
		entry.Cancel = cancel

		if err := entry.VU.Reconfigure(e.nextID); err != nil ***REMOVED***
			return err
		***REMOVED***
		go e.runVU(ctx, e.nextID, entry.VU)
		e.nextID++
	***REMOVED***
	for i := current - 1; i >= v; i-- ***REMOVED***
		entry := e.vus[i]
		entry.Cancel()
		entry.Cancel = nil
	***REMOVED***

	e.Status.VUs.Int64 = v
	return nil
***REMOVED***

func (e *Engine) SetMaxVUs(v int64) error ***REMOVED***
	e.vuMutex.Lock()
	defer e.vuMutex.Unlock()

	if v < e.Status.VUs.Int64 ***REMOVED***
		return ErrMaxTooLow
	***REMOVED***

	current := e.Status.VUsMax.Int64
	if v > current ***REMOVED***
		vus := e.vus
		for i := current; i < v; i++ ***REMOVED***
			vu, err := e.Runner.NewVU()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			vus = append(vus, &vuEntry***REMOVED***VU: vu***REMOVED***)
		***REMOVED***
		e.vus = vus
	***REMOVED*** else if v < current ***REMOVED***
		e.vus = e.vus[:v]
	***REMOVED***

	e.Status.VUsMax.Int64 = v
	return nil
***REMOVED***

func (e *Engine) reportInternalStats() ***REMOVED***
	e.mMutex.Lock()
	t := time.Now()
	e.getSink(MetricVUs).Add(stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(e.Status.VUs.Int64)***REMOVED***)
	e.getSink(MetricVUsMax).Add(stats.Sample***REMOVED***Time: t, Tags: nil, Value: float64(e.Status.VUsMax.Int64)***REMOVED***)
	e.mMutex.Unlock()
***REMOVED***

func (e *Engine) runVU(ctx context.Context, id int64, vu VU) ***REMOVED***
	idString := strconv.FormatInt(id, 10)

waitForPause:
	e.Pause.Wait()

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

			if !e.Status.Running.Bool ***REMOVED***
				goto waitForPause
			***REMOVED***
		***REMOVED***
	***REMOVED***
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
