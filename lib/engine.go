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
	ShutdownTimeout   = 10 * time.Second
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
	VU     VU
	Buffer []stats.Sample
	Cancel context.CancelFunc
***REMOVED***

type Engine struct ***REMOVED***
	Runner    Runner
	Status    Status
	Stages    []Stage
	Collector stats.Collector
	Pause     sync.WaitGroup
	Metrics   map[*stats.Metric]stats.Sink

	thresholds  map[string][]Threshold
	thresholdVM *otto.Otto

	ctx    context.Context
	vus    []*vuEntry
	nextID int64

	vuMutex   sync.Mutex
	waitGroup sync.WaitGroup
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
		thresholds:  make(map[string][]Threshold),
		thresholdVM: otto.New(),
	***REMOVED***
	e.Pause.Add(1)

	return e, nil
***REMOVED***

func (e *Engine) Apply(opts Options) error ***REMOVED***
	if opts.Run.Valid ***REMOVED***
		e.SetRunning(opts.Run.Bool)
	***REMOVED***
	if opts.VUsMax.Valid ***REMOVED***
		if err := e.SetMaxVUs(opts.VUs.Int64); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if opts.VUs.Valid ***REMOVED***
		if !opts.VUsMax.Valid ***REMOVED***
			if err := e.SetMaxVUs(opts.VUs.Int64); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if err := e.SetVUs(opts.VUs.Int64); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if opts.Duration.Valid ***REMOVED***
		duration, err := time.ParseDuration(opts.Duration.String)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		e.Stages = []Stage***REMOVED***Stage***REMOVED***Duration: null.IntFrom(int64(duration))***REMOVED******REMOVED***
	***REMOVED***

	if opts.Quit.Valid ***REMOVED***
		e.Status.Quit = opts.Quit
	***REMOVED***
	if opts.QuitOnTaint.Valid ***REMOVED***
		e.Status.QuitOnTaint = opts.QuitOnTaint
	***REMOVED***

	if opts.Thresholds != nil ***REMOVED***
		e.thresholds = opts.Thresholds

		// Make sure all scripts are compiled!
		for m, scripts := range e.thresholds ***REMOVED***
			for i, script := range scripts ***REMOVED***
				if script.Script != nil ***REMOVED***
					continue
				***REMOVED***

				s, err := e.thresholdVM.Compile(fmt.Sprintf("threshold$%s:%i", m, i), script.Source)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				script.Script = s
				scripts[i] = script
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (e *Engine) Run(ctx context.Context, opts Options) error ***REMOVED***
	subctx, cancel := context.WithCancel(context.Background())
	e.ctx = subctx
	e.nextID = 1

	e.Apply(opts)

	if e.Collector != nil ***REMOVED***
		e.waitGroup.Add(1)
		go func() ***REMOVED***
			e.Collector.Run(subctx)
			log.Debug("Engine: Collector shut down")
			e.waitGroup.Done()
		***REMOVED***()
	***REMOVED*** else ***REMOVED***
		log.Debug("Engine: No Collector")
	***REMOVED***

	e.waitGroup.Add(1)
	go func() ***REMOVED***
		e.runThresholds(subctx)
		log.Debug("Engine: Thresholds shut down")
		e.waitGroup.Done()
	***REMOVED***()

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
				if vu.Buffer == nil ***REMOVED***
					continue
				***REMOVED***

				buffer := vu.Buffer
				vu.Buffer = nil
				e.consumeBuffer(buffer)
			***REMOVED***

			if !ok ***REMOVED***
				e.SetRunning(false)

				if e.Status.Quit.Bool ***REMOVED***
					break loop
				***REMOVED*** else ***REMOVED***
					log.Info("Test finished, press Ctrl+C to exit")
					<-ctx.Done()
					break loop
				***REMOVED***
			***REMOVED***

			if e.Status.QuitOnTaint.Bool && e.Status.Tainted.Bool ***REMOVED***
				log.Warn("Test tainted, ending early...")
				break loop
			***REMOVED***

			e.consumeEngineStats()
		case <-ctx.Done():
			break loop
		***REMOVED***
	***REMOVED***

	e.SetRunning(false)
	e.SetVUs(0)
	e.SetMaxVUs(0)
	e.consumeEngineStats()

	cancel()

	log.Debug("Engine: Waiting for subsystem shutdown...")

	done := make(chan interface***REMOVED******REMOVED***)
	go func() ***REMOVED***
		e.waitGroup.Wait()
		close(done)
	***REMOVED***()
	timeout := time.After(ShutdownTimeout)
	select ***REMOVED***
	case <-done:
	case <-timeout:
		log.Warn("VUs took too long to finish, shutting down anyways")
	***REMOVED***

	return nil
***REMOVED***

func (e *Engine) IsRunning() bool ***REMOVED***
	return e.ctx != nil
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
		e.waitGroup.Add(1)
		go func() ***REMOVED***
			id := e.nextID
			e.runVU(ctx, id, entry)
			log.WithField("id", id).Debug("Engine: VU terminated")
			e.waitGroup.Done()
		***REMOVED***()
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

		if samples != nil ***REMOVED***
			vu.Buffer = append(vu.Buffer, samples...)
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
				scripts, ok := e.thresholds[m.Name]
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
					v, err := e.thresholdVM.Run(script.Script)
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
	if e.Collector != nil ***REMOVED***
		e.Collector.Collect(buffer)
	***REMOVED***
***REMOVED***

func (e *Engine) consumeEngineStats() ***REMOVED***
	t := time.Now()
	e.consumeBuffer([]stats.Sample***REMOVED***
		stats.Sample***REMOVED***Metric: MetricVUs, Time: t, Value: float64(e.Status.VUs.Int64)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricVUsMax, Time: t, Value: float64(e.Status.VUsMax.Int64)***REMOVED***,
	***REMOVED***)
***REMOVED***
