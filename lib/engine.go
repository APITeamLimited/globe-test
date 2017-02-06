/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TickRate        = 1 * time.Millisecond
	MetricsRate     = 1 * time.Second
	CollectRate     = 10 * time.Millisecond
	ThresholdsRate  = 2 * time.Second
	ShutdownTimeout = 10 * time.Second

	BackoffAmount = 50 * time.Millisecond
	BackoffMax    = 10 * time.Second
)

type vuEntry struct ***REMOVED***
	VU     VU
	Cancel context.CancelFunc

	Samples    []stats.Sample
	Iterations int64
	lock       sync.Mutex
***REMOVED***

type submetric struct ***REMOVED***
	Name       string
	Conditions map[string]string
	Metric     *stats.Metric
***REMOVED***

func parseSubmetric(name string) (string, map[string]string) ***REMOVED***
	halves := strings.SplitN(strings.TrimSuffix(name, "***REMOVED***"), "***REMOVED***", 2)
	if len(halves) != 2 ***REMOVED***
		return halves[0], nil
	***REMOVED***

	kvs := strings.Split(halves[1], ",")
	conditions := make(map[string]string, len(kvs))
	for _, kv := range kvs ***REMOVED***
		if kv == "" ***REMOVED***
			continue
		***REMOVED***

		parts := strings.SplitN(kv, ":", 2)

		key := strings.TrimSpace(strings.Trim(parts[0], `"'`))
		if len(parts) != 2 ***REMOVED***
			conditions[key] = ""
			continue
		***REMOVED***

		value := strings.TrimSpace(strings.Trim(parts[1], `"'`))
		conditions[key] = value
	***REMOVED***
	return halves[0], conditions
***REMOVED***

// The Engine is the beating heart of K6.
type Engine struct ***REMOVED***
	Runner    Runner
	Options   Options
	Collector Collector
	Logger    *log.Logger

	Stages      []Stage
	Thresholds  map[string]Thresholds
	Metrics     map[*stats.Metric]stats.Sink
	MetricsLock sync.Mutex

	// Submetrics, mapped from parent metric names.
	submetrics map[string][]*submetric

	// Stage tracking.
	atTime          time.Duration
	atStage         int
	atStageSince    time.Duration
	atStageStartVUs int64

	// VU tracking.
	vus       int64
	vusMax    int64
	vuEntries []*vuEntry
	vuStop    chan interface***REMOVED******REMOVED***
	vuPause   chan interface***REMOVED******REMOVED***

	// Atomic counters.
	numIterations int64
	numErrors     int64

	thresholdsTainted bool

	// Subsystem-related.
	lock      sync.Mutex
	subctx    context.Context
	subcancel context.CancelFunc
	subwg     sync.WaitGroup
***REMOVED***

func NewEngine(r Runner, o Options) (*Engine, error) ***REMOVED***
	e := &Engine***REMOVED***
		Runner:  r,
		Options: o,
		Logger:  log.StandardLogger(),

		Metrics:    make(map[*stats.Metric]stats.Sink),
		Thresholds: make(map[string]Thresholds),

		vuStop: make(chan interface***REMOVED******REMOVED***),
	***REMOVED***
	e.clearSubcontext()

	if o.Duration.Valid ***REMOVED***
		d, err := time.ParseDuration(o.Duration.String)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "options.duration")
		***REMOVED***
		e.Stages = []Stage***REMOVED******REMOVED***Duration: d***REMOVED******REMOVED***
	***REMOVED***
	if o.Stages != nil ***REMOVED***
		e.Stages = o.Stages
	***REMOVED***
	if o.VUsMax.Valid ***REMOVED***
		if err := e.SetVUsMax(o.VUsMax.Int64); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if o.VUs.Valid ***REMOVED***
		if err := e.SetVUs(o.VUs.Int64); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if o.Paused.Valid ***REMOVED***
		e.SetPaused(o.Paused.Bool)
	***REMOVED***

	if o.Thresholds != nil ***REMOVED***
		e.Thresholds = o.Thresholds
		e.submetrics = make(map[string][]*submetric)

		for name := range e.Thresholds ***REMOVED***
			if !strings.Contains(name, "***REMOVED***") ***REMOVED***
				continue
			***REMOVED***

			parent, conds := parseSubmetric(name)
			e.submetrics[parent] = append(e.submetrics[parent], &submetric***REMOVED***
				Name:       name,
				Conditions: conds,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return e, nil
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	collectorctx, collectorcancel := context.WithCancel(ctx)
	collectorch := make(chan interface***REMOVED******REMOVED***)
	if e.Collector != nil ***REMOVED***
		go func(ctx context.Context) ***REMOVED***
			e.Collector.Run(ctx)
			close(collectorch)
		***REMOVED***(collectorctx)
	***REMOVED*** else ***REMOVED***
		close(collectorch)
	***REMOVED***

	e.lock.Lock()
	***REMOVED***
		// Run metrics emission.
		e.subwg.Add(1)
		go func(ctx context.Context) ***REMOVED***
			e.runMetricsEmission(ctx)
			e.subwg.Done()
		***REMOVED***(e.subctx)

		// Run metrics collection.
		e.subwg.Add(1)
		go func(ctx context.Context) ***REMOVED***
			e.runCollection(ctx)
			e.subwg.Done()
		***REMOVED***(e.subctx)

		// Run thresholds.
		e.subwg.Add(1)
		go func(ctx context.Context) ***REMOVED***
			e.runThresholds(ctx)
			e.subwg.Done()
		***REMOVED***(e.subctx)
	***REMOVED***
	e.lock.Unlock()

	close(e.vuStop)
	defer func() ***REMOVED***
		e.vuStop = make(chan interface***REMOVED******REMOVED***)
		e.SetPaused(false)

		// Shut down subsystems, wait for graceful termination.
		e.clearSubcontext()
		e.subwg.Wait()

		// Process any leftover samples.
		e.processSamples(e.collect()...)
		collectorcancel()
		<-collectorch

		// Emit final metrics.
		e.emitMetrics()
	***REMOVED***()

	// Set tracking to defaults.
	e.lock.Lock()
	e.atTime = 0
	e.atStage = 0
	e.atStageSince = 0
	e.atStageStartVUs = e.vus
	e.numIterations = 0
	e.numErrors = 0
	e.lock.Unlock()

	var lastTick time.Time
	ticker := time.NewTicker(TickRate)

	maxIterations := e.Options.Iterations.Int64
	for ***REMOVED***
		// Don't do anything while the engine is paused.
		vuPause := e.vuPause
		if vuPause != nil ***REMOVED***
			select ***REMOVED***
			case <-vuPause:
			case <-ctx.Done():
				return nil
			***REMOVED***
		***REMOVED***

		// If we have an iteration cap, exit once we hit it.
		if maxIterations > 0 && e.numIterations == e.vusMax*maxIterations ***REMOVED***
			return nil
		***REMOVED***

		// Calculate the time delta between now and the last tick.
		now := time.Now()
		if lastTick.IsZero() ***REMOVED***
			lastTick = now
		***REMOVED***
		dT := now.Sub(lastTick)
		lastTick = now

		// Update state.
		keepRunning, err := e.processStages(dT)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !keepRunning ***REMOVED***
			return nil
		***REMOVED***

		select ***REMOVED***
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) IsRunning() bool ***REMOVED***
	select ***REMOVED***
	case <-e.vuStop:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (e *Engine) SetPaused(v bool) ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	if v && e.vuPause == nil ***REMOVED***
		e.vuPause = make(chan interface***REMOVED******REMOVED***)
	***REMOVED*** else if !v && e.vuPause != nil ***REMOVED***
		close(e.vuPause)
		e.vuPause = nil
	***REMOVED***
***REMOVED***

func (e *Engine) IsPaused() bool ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.vuPause != nil
***REMOVED***

func (e *Engine) SetVUs(v int64) error ***REMOVED***
	if v < 0 ***REMOVED***
		return errors.New("vus can't be negative")
	***REMOVED***

	e.lock.Lock()
	defer e.lock.Unlock()

	if v > e.vusMax ***REMOVED***
		return errors.New("more vus than allocated requested")
	***REMOVED***

	// Scale up
	for i := e.vus; i < v; i++ ***REMOVED***
		vu := e.vuEntries[i]
		if vu.Cancel != nil ***REMOVED***
			panic(errors.New("fatal miscalculation: attempted to re-schedule active VU"))
		***REMOVED***

		ctx, cancel := context.WithCancel(e.subctx)
		vu.Cancel = cancel

		e.subwg.Add(1)
		go func() ***REMOVED***
			e.runVU(ctx, vu)
			e.subwg.Done()
		***REMOVED***()
	***REMOVED***

	// Scale down
	for i := e.vus - 1; i >= v; i-- ***REMOVED***
		vu := e.vuEntries[i]
		vu.Cancel()
		vu.Cancel = nil
	***REMOVED***

	e.vus = v
	return nil
***REMOVED***

func (e *Engine) GetVUs() int64 ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.vus
***REMOVED***

func (e *Engine) SetVUsMax(v int64) error ***REMOVED***
	if v < 0 ***REMOVED***
		return errors.New("vus-max can't be negative")
	***REMOVED***

	e.lock.Lock()
	defer e.lock.Unlock()

	if v < e.vus ***REMOVED***
		return errors.New("can't reduce vus-max below vus")
	***REMOVED***

	// Scale up
	for len(e.vuEntries) < int(v) ***REMOVED***
		var entry vuEntry
		if e.Runner != nil ***REMOVED***
			vu, err := e.Runner.NewVU()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			entry.VU = vu
		***REMOVED***
		e.vuEntries = append(e.vuEntries, &entry)
	***REMOVED***

	// Scale down
	if len(e.vuEntries) > int(v) ***REMOVED***
		e.vuEntries = e.vuEntries[:int(v)]
	***REMOVED***

	e.vusMax = v
	return nil
***REMOVED***

func (e *Engine) GetVUsMax() int64 ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.vusMax
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	return e.thresholdsTainted
***REMOVED***

func (e *Engine) AtTime() time.Duration ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.atTime
***REMOVED***

func (e *Engine) TotalTime() time.Duration ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	var total time.Duration
	for _, stage := range e.Stages ***REMOVED***
		if stage.Duration <= 0 ***REMOVED***
			return 0
		***REMOVED***
		total += stage.Duration
	***REMOVED***
	return total
***REMOVED***

func (e *Engine) clearSubcontext() ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.subcancel != nil ***REMOVED***
		e.subcancel()
	***REMOVED***
	subctx, subcancel := context.WithCancel(context.Background())
	e.subctx = subctx
	e.subcancel = subcancel
***REMOVED***

func (e *Engine) processStages(dT time.Duration) (bool, error) ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	e.atTime += dT

	// If there are no stages, just keep going indefinitely at a stable VU count.
	if len(e.Stages) == 0 ***REMOVED***
		return true, nil
	***REMOVED***

	stage := e.Stages[e.atStage]
	if stage.Duration > 0 && e.atTime > e.atStageSince+stage.Duration ***REMOVED***
		if e.atStage != len(e.Stages)-1 ***REMOVED***
			e.atStage++
			e.atStageSince = e.atTime
			e.atStageStartVUs = e.vus
			stage = e.Stages[e.atStage]
		***REMOVED*** else ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***
	if stage.Target.Valid ***REMOVED***
		from := e.atStageStartVUs
		to := stage.Target.Int64
		t := 1.0
		if stage.Duration > 0 ***REMOVED***
			t = Clampf(float64(e.atTime)/float64(e.atStageSince+stage.Duration), 0.0, 1.0)
		***REMOVED***
		if err := e.SetVUs(Lerp(from, to, t)); err != nil ***REMOVED***
			return false, errors.Wrapf(err, "stage #%d", e.atStage+1)
		***REMOVED***
	***REMOVED***

	return true, nil
***REMOVED***

func (e *Engine) runVU(ctx context.Context, vu *vuEntry) ***REMOVED***
	maxIterations := e.Options.Iterations.Int64

	// nil runners that produce nil VUs are used for testing.
	if vu.VU == nil ***REMOVED***
		<-ctx.Done()
		return
	***REMOVED***

	// Sleep until the engine starts running.
	<-e.vuStop

	backoffCounter := 0
	backoff := time.Duration(0)
	for ***REMOVED***
		// Exit if the VU has run all its intended iterations.
		if maxIterations > 0 && vu.Iterations >= maxIterations ***REMOVED***
			return
		***REMOVED***

		// If the engine is paused, sleep until it resumes.
		vuPause := e.vuPause
		if vuPause != nil ***REMOVED***
			<-vuPause
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***

		succ := e.runVUOnce(ctx, vu)
		if !succ ***REMOVED***
			backoffCounter++
			backoff += BackoffAmount * time.Duration(backoffCounter)
			if backoff > BackoffMax ***REMOVED***
				backoff = BackoffMax
			***REMOVED***
			select ***REMOVED***
			case <-time.After(backoff):
			case <-ctx.Done():
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			backoff = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) runVUOnce(ctx context.Context, vu *vuEntry) bool ***REMOVED***
	samples, err := vu.VU.RunOnce(ctx)

	// Expired VUs usually have request cancellation errors, and thus skewed metrics and
	// unhelpful "request cancelled" errors. Don't process those.
	select ***REMOVED***
	case <-ctx.Done():
		return true
	default:
	***REMOVED***

	if err != nil ***REMOVED***
		if serr, ok := err.(fmt.Stringer); ok ***REMOVED***
			e.Logger.Error(serr.String())
		***REMOVED*** else ***REMOVED***
			e.Logger.WithError(err).Error("VU Error")
		***REMOVED***
		samples = append(samples,
			stats.Sample***REMOVED***
				Time:   time.Now(),
				Metric: metrics.Errors,
				Tags:   map[string]string***REMOVED***"error": err.Error()***REMOVED***,
				Value:  1,
			***REMOVED***,
		)
	***REMOVED***

	vu.lock.Lock()
	vu.Samples = append(vu.Samples, samples...)
	vu.lock.Unlock()

	atomic.AddInt64(&vu.Iterations, 1)
	atomic.AddInt64(&e.numIterations, 1)
	if err != nil ***REMOVED***
		atomic.AddInt64(&e.numErrors, 1)
		return false
	***REMOVED***
	return true
***REMOVED***

func (e *Engine) runMetricsEmission(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(MetricsRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.emitMetrics()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) emitMetrics() ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	t := time.Now()
	e.processSamples(
		stats.Sample***REMOVED***
			Time:   t,
			Metric: metrics.VUs,
			Value:  float64(e.vus),
		***REMOVED***,
		stats.Sample***REMOVED***
			Time:   t,
			Metric: metrics.VUsMax,
			Value:  float64(e.vusMax),
		***REMOVED***,
		stats.Sample***REMOVED***
			Time:   t,
			Metric: metrics.Iterations,
			Value:  float64(atomic.LoadInt64(&e.numIterations)),
		***REMOVED***,
	)
***REMOVED***

func (e *Engine) runThresholds(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(ThresholdsRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.processThresholds()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) processThresholds() ***REMOVED***
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	e.thresholdsTainted = false
	for m, s := range e.Metrics ***REMOVED***
		ts, ok := e.Thresholds[m.Name]
		if !ok ***REMOVED***
			continue
		***REMOVED***

		e.Logger.WithField("m", m.Name).Debug("running thresholds")
		succ, err := ts.Run(s)
		if err != nil ***REMOVED***
			e.Logger.WithField("m", m.Name).WithError(err).Error("Threshold error")
			continue
		***REMOVED***
		if !succ ***REMOVED***
			e.Logger.WithField("m", m.Name).Debug("Thresholds failed")
			e.thresholdsTainted = true
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) runCollection(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(CollectRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.processSamples(e.collect()...)
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) collect() []stats.Sample ***REMOVED***
	e.lock.Lock()
	entries := e.vuEntries
	e.lock.Unlock()

	samples := []stats.Sample***REMOVED******REMOVED***
	for _, vu := range entries ***REMOVED***
		vu.lock.Lock()
		if len(vu.Samples) > 0 ***REMOVED***
			samples = append(samples, vu.Samples...)
			vu.Samples = nil
		***REMOVED***
		vu.lock.Unlock()
	***REMOVED***
	return samples
***REMOVED***

func (e *Engine) processSamples(samples ...stats.Sample) ***REMOVED***
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	for _, sample := range samples ***REMOVED***
		sink := e.Metrics[sample.Metric]
		if sink == nil ***REMOVED***
			sink = sample.Metric.NewSink()
			e.Metrics[sample.Metric] = sink
		***REMOVED***
		sink.Add(sample)

		for _, sm := range e.submetrics[sample.Metric.Name] ***REMOVED***
			passing := true
			for k, v := range sm.Conditions ***REMOVED***
				if sample.Tags[k] != v ***REMOVED***
					passing = false
					break
				***REMOVED***
			***REMOVED***
			if !passing ***REMOVED***
				continue
			***REMOVED***

			if sm.Metric == nil ***REMOVED***
				sm.Metric = stats.New(sm.Name, sample.Metric.Type, sample.Metric.Contains)
				e.Metrics[sm.Metric] = sm.Metric.NewSink()
			***REMOVED***
			e.Metrics[sm.Metric].Add(sample)
		***REMOVED***
	***REMOVED***

	if e.Collector != nil ***REMOVED***
		e.Collector.Collect(samples)
	***REMOVED***
***REMOVED***
