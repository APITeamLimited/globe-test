package lib

import (
	"context"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"sync"
	"time"
)

const (
	TickRate          = 1 * time.Millisecond
	ThresholdTickRate = 2 * time.Second
)

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
	Runner      Runner
	Status      Status
	Stages      []Stage
	Collector   stats.Collector
	Quit        bool
	QuitOnTaint bool
	Pause       sync.WaitGroup

	Metrics    map[*stats.Metric]stats.Sink
	Thresholds map[string][]*otto.Script

	thresholdVM *otto.Otto

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
		Metrics:     make(map[*stats.Metric]stats.Sink),
		Thresholds:  make(map[string][]*otto.Script),
		thresholdVM: otto.New(),
	***REMOVED***

	e.Status.Running = null.BoolFrom(false)
	e.Pause.Add(1)

	e.Status.VUs = null.IntFrom(0)
	e.Status.VUsMax = null.IntFrom(0)

	return e, nil
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	if len(e.Stages) == 0 ***REMOVED***
		return errors.New("Engine has no stages")
	***REMOVED***

	e.ctx = ctx
	e.nextID = 1

	if e.Collector != nil ***REMOVED***
		go e.Collector.Run(ctx)
	***REMOVED*** else ***REMOVED***
		log.Debug("Engine: No Collector")
	***REMOVED***

	go e.runThresholds(ctx)

	e.consumeEngineStats()

	ticker := time.NewTicker(TickRate)
	lastTick := time.Now()

loop:
	for ***REMOVED***
		select ***REMOVED***
		case now := <-ticker.C:
			timeDelta := now.Sub(lastTick)
			e.Status.AtTime.Int64 += int64(timeDelta)
			lastTick = now

			stage, left, ok := StageAt(e.Stages, time.Duration(e.Status.AtTime.Int64))
			if stage.StartVUs.Valid && stage.EndVUs.Valid ***REMOVED***
				progress := (float64(stage.Duration.Int64-int64(left)) / float64(stage.Duration.Int64))
				vus := Lerp(stage.StartVUs.Int64, stage.EndVUs.Int64, progress)
				e.SetVUs(vus)
			***REMOVED***

			for _, vu := range e.vus ***REMOVED***
				e.consumeBuffer(vu.Buffer)
			***REMOVED***

			if !ok ***REMOVED***
				e.SetRunning(false)

				if e.Quit ***REMOVED***
					break loop
				***REMOVED*** else ***REMOVED***
					log.Info("Test finished, press Ctrl+C to exit")
					<-ctx.Done()
					break loop
				***REMOVED***
			***REMOVED***

			if e.QuitOnTaint && e.Status.Tainted.Bool ***REMOVED***
				log.Warn("Test tainted, ending early...")
				break loop
			***REMOVED***

			e.consumeEngineStats()
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
	if e.Status.VUs.Int64 == v ***REMOVED***
		return nil
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***"from": e.Status.VUs.Int64, "to": v***REMOVED***).Debug("Setting VUs")

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
	if e.Status.VUsMax.Int64 == v ***REMOVED***
		return nil
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***"from": e.Status.VUsMax.Int64, "to": v***REMOVED***).Debug("Setting Max VUs")

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

func (e *Engine) AddThreshold(metric, src string) error ***REMOVED***
	script, err := e.thresholdVM.Compile("__threshold__", src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	e.Thresholds[metric] = append(e.Thresholds[metric], script)

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

func (e *Engine) runThresholds(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(ThresholdTickRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			for m, sink := range e.Metrics ***REMOVED***
				scripts, ok := e.Thresholds[m.Name]
				if !ok ***REMOVED***
					continue
				***REMOVED***

				sample := sink.Format()
				for key, value := range sample ***REMOVED***
					if m.Contains == stats.Time ***REMOVED***
						value = value / float64(time.Millisecond)
					***REMOVED***
					// log.WithFields(log.Fields***REMOVED***"k": key, "v": value***REMOVED***).Debug("setting threshold data")
					e.thresholdVM.Set(key, value)
				***REMOVED***

				taint := false
				for _, script := range scripts ***REMOVED***
					v, err := e.thresholdVM.Run(script.String())
					if err != nil ***REMOVED***
						log.WithError(err).WithField("metric", m.Name).Error("Threshold Error")
						taint = true
						continue
					***REMOVED***
					// log.WithFields(log.Fields***REMOVED***"metric": m.Name, "v": v, "s": sample***REMOVED***).Debug("threshold tick")
					bV, err := v.ToBoolean()
					if err != nil ***REMOVED***
						log.WithError(err).WithField("metric", m.Name).Error("Threshold result is invalid")
						taint = true
						continue
					***REMOVED***
					if !bV ***REMOVED***
						taint = true
					***REMOVED***
				***REMOVED***

				for key, _ := range sample ***REMOVED***
					e.thresholdVM.Set(key, otto.UndefinedValue())
				***REMOVED***

				if taint ***REMOVED***
					m.Tainted = true
					e.Status.Tainted.Bool = true
				***REMOVED***
			***REMOVED***
		case <-ctx.Done():
			return
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
		case stats.Rate:
			s = &stats.RateSink***REMOVED******REMOVED***
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
