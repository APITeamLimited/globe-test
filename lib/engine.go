package lib

import (
	"context"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"sync"
	"time"
)

const TickRate = 500 * time.Millisecond

var (
	MetricVUs    = &stats.Metric***REMOVED***Name: "vus", Type: stats.Gauge***REMOVED***
	MetricVUsMax = &stats.Metric***REMOVED***Name: "vus_max", Type: stats.Gauge***REMOVED***
	MetricErrors = &stats.Metric***REMOVED***Name: "errors", Type: stats.Counter***REMOVED***

	ErrTooManyVUs = errors.New("More VUs than the maximum requested")
	ErrMaxTooLow  = errors.New("Can't lower max below current VU count")

	// Special error used to taint a test, without printing an error.
	ErrVUWantsTaint = errors.New("[ErrVUWantsTaint is never logged]")
)

type vuEntry struct ***REMOVED***
	VU        VU
	Buffer    []stats.Sample
	ExtBuffer stats.Buffer
	Cancel    context.CancelFunc
***REMOVED***

type Engine struct ***REMOVED***
	Runner    Runner
	Status    Status
	Stages    []Stage
	Collector stats.Collector
	Quit      bool
	Pause     sync.WaitGroup

	Metrics map[*stats.Metric]stats.Sink

	ctx    context.Context
	vus    []*vuEntry
	nextID int64

	vuMutex sync.Mutex
***REMOVED***

func NewEngine(r Runner) (*Engine, error) ***REMOVED***
	e := &Engine***REMOVED***
		Runner: r,
		Status: Status***REMOVED***
			Running: null.BoolFrom(false),
			Tainted: null.BoolFrom(false),
			VUs:     null.IntFrom(0),
			VUsMax:  null.IntFrom(0),
			AtTime:  null.IntFrom(0),
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

	e.consumeEngineStats()
	ticker := time.NewTicker(TickRate)

	if e.Collector != nil ***REMOVED***
		go e.Collector.Run(ctx)
	***REMOVED*** else ***REMOVED***
		log.Debug("Engine: No Collector")
	***REMOVED***

loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.consumeEngineStats()

			for _, vu := range e.vus ***REMOVED***
				e.consumeBuffer(vu.Buffer)
				vu.Buffer = vu.Buffer[:0]
			***REMOVED***

			if e.Status.Running.Bool ***REMOVED***
				e.Status.AtTime.Int64 += int64(TickRate)

				stage, stageLeft, ok := StageAt(e.Stages, time.Duration(e.Status.AtTime.Int64))
				if stage.VUTarget.Valid ***REMOVED***
					t := e.Status.AtTime.Int64
					tx := t - int64(TickRate)
					ty := t + int64(stageLeft)
					x := e.Status.VUs.Int64
					y := stage.VUTarget.Int64
					vus := Ease(t, tx, ty, x, y)

					if vus != e.Status.VUs.Int64 ***REMOVED***
						log.WithField("vus", vus).Debug("Engine: Interpolating VUs...")
						if err := e.SetVUs(vus); err != nil ***REMOVED***
							log.WithError(err).WithField("vus", vus).Error("Engine: VU interpolation failed")
						***REMOVED***
					***REMOVED***
				***REMOVED***

				if !ok ***REMOVED***
					e.SetRunning(false)

					if !e.Quit ***REMOVED***
						log.Info("Test expired, execution paused, pass --quit to exit here")
					***REMOVED*** else ***REMOVED***
						log.Info("Test ended, bye!")
						break loop
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case <-ctx.Done():
			break loop
		***REMOVED***
	***REMOVED***

	e.vus = nil

	e.Status.Running = null.BoolFrom(false)
	e.Status.VUs = null.IntFrom(0)
	e.Status.VUsMax = null.IntFrom(0)
	e.consumeEngineStats()

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
		go e.runVU(ctx, e.nextID, entry)
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
			entry := &vuEntry***REMOVED***VU: vu***REMOVED***
			if e.Collector != nil ***REMOVED***
				entry.ExtBuffer = e.Collector.Buffer()
			***REMOVED***
			vus = append(vus, entry)
		***REMOVED***
		e.vus = vus
	***REMOVED*** else if v < current ***REMOVED***
		e.vus = e.vus[:v]
	***REMOVED***

	e.Status.VUsMax.Int64 = v
	return nil
***REMOVED***

func (e *Engine) runVU(ctx context.Context, id int64, vu *vuEntry) ***REMOVED***
	idString := strconv.FormatInt(id, 10)

waitForPause:
	e.Pause.Wait()

	for ***REMOVED***
		samples, err := vu.VU.RunOnce(ctx)

		// If the context is cancelled, the iteration is likely to produce erroneous output
		// due to cancelled HTTP requests and whatnot. Discard output from such runs.
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***

		if err != nil ***REMOVED***
			e.Status.Tainted.Bool = true

			if err != ErrVUWantsTaint ***REMOVED***
				if s, ok := err.(fmt.Stringer); ok ***REMOVED***
					log.Error(s.String())
				***REMOVED*** else ***REMOVED***
					log.WithError(err).Error("Runtime Error")
				***REMOVED***

				samples = append(samples, stats.Sample***REMOVED***
					Metric: MetricErrors,
					Time:   time.Now(),
					Tags:   map[string]string***REMOVED***"vu": idString, "error": err.Error()***REMOVED***,
					Value:  float64(1),
				***REMOVED***)
			***REMOVED***
		***REMOVED***

		vu.Buffer = append(vu.Buffer, samples...)
		if vu.ExtBuffer != nil ***REMOVED***
			vu.ExtBuffer.Add(samples...)
		***REMOVED***

		if !e.Status.Running.Bool ***REMOVED***
			goto waitForPause
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

func (e *Engine) consumeBuffer(buffer []stats.Sample) ***REMOVED***
	for _, sample := range buffer ***REMOVED***
		e.getSink(sample.Metric).Add(sample)
	***REMOVED***
***REMOVED***

func (e *Engine) consumeEngineStats() ***REMOVED***
	t := time.Now()
	e.consumeBuffer([]stats.Sample***REMOVED***
		stats.Sample***REMOVED***Metric: MetricVUs, Time: t, Value: float64(e.Status.VUs.Int64)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricVUsMax, Time: t, Value: float64(e.Status.VUsMax.Int64)***REMOVED***,
	***REMOVED***)
***REMOVED***
