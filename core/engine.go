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

package core

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

const (
	TickRate        = 1 * time.Millisecond
	MetricsRate     = 1 * time.Second
	CollectRate     = 50 * time.Millisecond
	ThresholdsRate  = 2 * time.Second
	ShutdownTimeout = 10 * time.Second

	BackoffAmount = 50 * time.Millisecond
	BackoffMax    = 10 * time.Second
)

// The Engine is the beating heart of K6.
type Engine struct ***REMOVED***
	runLock sync.Mutex // y tho? TODO: remove?

	//TODO: make most of the stuff here private!
	ExecutionScheduler lib.ExecutionScheduler
	executionState     *lib.ExecutionState

	Options      lib.Options
	Collectors   []lib.Collector
	NoThresholds bool
	NoSummary    bool

	logger *logrus.Logger

	Metrics     map[string]*stats.Metric
	MetricsLock sync.Mutex

	Samples chan stats.SampleContainer

	// Assigned to metrics upon first received sample.
	thresholds map[string]stats.Thresholds
	submetrics map[string][]*stats.Submetric

	// Are thresholds tainted?
	thresholdsTainted bool
***REMOVED***

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(ex lib.ExecutionScheduler, o lib.Options, logger *logrus.Logger) (*Engine, error) ***REMOVED***
	if ex == nil ***REMOVED***
		return nil, errors.New("missing ExecutionScheduler instance")
	***REMOVED***

	e := &Engine***REMOVED***
		ExecutionScheduler: ex,
		executionState:     ex.GetState(),

		Options: o,
		Metrics: make(map[string]*stats.Metric),
		Samples: make(chan stats.SampleContainer, o.MetricSamplesBufferSize.Int64),
		logger:  logger,
	***REMOVED***

	e.thresholds = o.Thresholds
	e.submetrics = make(map[string][]*stats.Submetric)
	for name := range e.thresholds ***REMOVED***
		if !strings.Contains(name, "***REMOVED***") ***REMOVED***
			continue
		***REMOVED***

		parent, sm := stats.NewSubmetric(name)
		e.submetrics[parent] = append(e.submetrics[parent], sm)
	***REMOVED***

	return e, nil
***REMOVED***

// Init is used to initialize the execuction scheduler. That's a costly operation, since it
// initializes all of the planned VUs and could potentially take a long time.
func (e *Engine) Init(ctx context.Context) error ***REMOVED***
	return e.ExecutionScheduler.Init(ctx, e.Samples)
***REMOVED***

func (e *Engine) setRunStatus(status lib.RunStatus) ***REMOVED***
	if len(e.Collectors) == 0 ***REMOVED***
		return
	***REMOVED***

	for _, c := range e.Collectors ***REMOVED***
		c.SetRunStatus(status)
	***REMOVED***
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.runLock.Lock()
	defer e.runLock.Unlock()

	e.logger.Debug("Engine: Starting with parameters...")

	collectorwg := sync.WaitGroup***REMOVED******REMOVED***
	collectorctx, collectorcancel := context.WithCancel(context.Background())
	if len(e.Collectors) > 0 ***REMOVED***
		for _, collector := range e.Collectors ***REMOVED***
			collectorwg.Add(1)
			go func(collector lib.Collector) ***REMOVED***
				collector.Run(collectorctx)
				collectorwg.Done()
			***REMOVED***(collector)
		***REMOVED***
	***REMOVED***

	subctx, subcancel := context.WithCancel(context.Background())
	subwg := sync.WaitGroup***REMOVED******REMOVED***

	// Run metrics emission.
	subwg.Add(1)
	go func() ***REMOVED***
		e.runMetricsEmission(subctx)
		e.logger.Debug("Engine: Emission terminated")
		subwg.Done()
	***REMOVED***()

	// Run thresholds.
	if !e.NoThresholds ***REMOVED***
		subwg.Add(1)
		go func() ***REMOVED***
			e.runThresholds(subctx, subcancel)
			e.logger.Debug("Engine: Thresholds terminated")
			subwg.Done()
		***REMOVED***()
	***REMOVED***

	// Run the execution scheduler.
	errC := make(chan error)
	subwg.Add(1)
	go func() ***REMOVED***
		errC <- e.ExecutionScheduler.Run(subctx, e.Samples)
		e.logger.Debug("Engine: Execution scheduler terminated")
		subwg.Done()
	***REMOVED***()

	sampleContainers := []stats.SampleContainer***REMOVED******REMOVED***
	defer func() ***REMOVED***
		// Shut down subsystems.
		subcancel()

		// Process samples until the subsystems have shut down.
		// Filter out samples produced past the end of a test.
		go func() ***REMOVED***
			if errC != nil ***REMOVED***
				<-errC
				errC = nil
			***REMOVED***
			subwg.Wait()
			close(e.Samples)
		***REMOVED***()

		for sc := range e.Samples ***REMOVED***
			sampleContainers = append(sampleContainers, sc)
		***REMOVED***
		if len(sampleContainers) > 0 ***REMOVED***
			e.processSamples(sampleContainers)
		***REMOVED***

		// Process final thresholds.
		if !e.NoThresholds ***REMOVED***
			e.processThresholds(nil)
		***REMOVED***

		// Finally, shut down collector.
		collectorcancel()
		collectorwg.Wait()
	***REMOVED***()

	ticker := time.NewTicker(CollectRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			if len(sampleContainers) > 0 ***REMOVED***
				e.processSamples(sampleContainers)
				sampleContainers = []stats.SampleContainer***REMOVED******REMOVED***
			***REMOVED***
		case sc := <-e.Samples:
			sampleContainers = append(sampleContainers, sc)
		case err := <-errC:
			errC = nil
			if err != nil ***REMOVED***
				e.logger.WithError(err).Debug("run: execution scheduler returned an error")
				e.setRunStatus(lib.RunStatusAbortedSystem)
				return err
			***REMOVED***
			e.logger.Debug("run: execution scheduler terminated")
			return nil
		case <-ctx.Done():
			e.logger.Debug("run: context expired; exiting...")
			e.setRunStatus(lib.RunStatusAbortedUser)
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	return e.thresholdsTainted
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
	t := time.Now()

	executionState := e.ExecutionScheduler.GetState()
	e.processSamples([]stats.SampleContainer***REMOVED***stats.ConnectedSamples***REMOVED***
		Samples: []stats.Sample***REMOVED***
			***REMOVED***
				Time:   t,
				Metric: metrics.VUs,
				Value:  float64(executionState.GetCurrentlyActiveVUsCount()),
				Tags:   e.Options.RunTags,
			***REMOVED***, ***REMOVED***
				Time:   t,
				Metric: metrics.VUsMax,
				Value:  float64(executionState.GetInitializedVUsCount()),
				Tags:   e.Options.RunTags,
			***REMOVED***,
		***REMOVED***,
		Tags: e.Options.RunTags,
		Time: t,
	***REMOVED******REMOVED***)
***REMOVED***

func (e *Engine) runThresholds(ctx context.Context, abort func()) ***REMOVED***
	ticker := time.NewTicker(ThresholdsRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.processThresholds(abort)
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) processThresholds(abort func()) ***REMOVED***
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	t := e.executionState.GetCurrentTestRunDuration()
	abortOnFail := false

	e.thresholdsTainted = false
	for _, m := range e.Metrics ***REMOVED***
		if len(m.Thresholds.Thresholds) == 0 ***REMOVED***
			continue
		***REMOVED***
		m.Tainted = null.BoolFrom(false)

		e.logger.WithField("m", m.Name).Debug("running thresholds")
		succ, err := m.Thresholds.Run(m.Sink, t)
		if err != nil ***REMOVED***
			e.logger.WithField("m", m.Name).WithError(err).Error("Threshold error")
			continue
		***REMOVED***
		if !succ ***REMOVED***
			e.logger.WithField("m", m.Name).Debug("Thresholds failed")
			m.Tainted = null.BoolFrom(true)
			e.thresholdsTainted = true
			if !abortOnFail && m.Thresholds.Abort ***REMOVED***
				abortOnFail = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if abortOnFail && abort != nil ***REMOVED***
		//TODO: When sending this status we get a 422 Unprocessable Entity
		e.setRunStatus(lib.RunStatusAbortedThreshold)
		abort()
	***REMOVED***
***REMOVED***

func (e *Engine) processSamplesForMetrics(sampleCointainers []stats.SampleContainer) ***REMOVED***
	for _, sampleCointainer := range sampleCointainers ***REMOVED***
		samples := sampleCointainer.GetSamples()

		if len(samples) == 0 ***REMOVED***
			continue
		***REMOVED***

		for _, sample := range samples ***REMOVED***
			m, ok := e.Metrics[sample.Metric.Name]
			if !ok ***REMOVED***
				m = stats.New(sample.Metric.Name, sample.Metric.Type, sample.Metric.Contains)
				m.Thresholds = e.thresholds[m.Name]
				m.Submetrics = e.submetrics[m.Name]
				e.Metrics[m.Name] = m
			***REMOVED***
			m.Sink.Add(sample)

			for _, sm := range m.Submetrics ***REMOVED***
				if !sample.Tags.Contains(sm.Tags) ***REMOVED***
					continue
				***REMOVED***

				if sm.Metric == nil ***REMOVED***
					sm.Metric = stats.New(sm.Name, sample.Metric.Type, sample.Metric.Contains)
					sm.Metric.Sub = *sm
					sm.Metric.Thresholds = e.thresholds[sm.Name]
					e.Metrics[sm.Name] = sm.Metric
				***REMOVED***
				sm.Metric.Sink.Add(sample)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) processSamples(sampleCointainers []stats.SampleContainer) ***REMOVED***
	if len(sampleCointainers) == 0 ***REMOVED***
		return
	***REMOVED***

	// TODO: optimize this...
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	// TODO: run this and the below code in goroutines?
	if !(e.NoSummary && e.NoThresholds) ***REMOVED***
		e.processSamplesForMetrics(sampleCointainers)
	***REMOVED***

	if len(e.Collectors) > 0 ***REMOVED***
		for _, collector := range e.Collectors ***REMOVED***
			collector.Collect(sampleCointainers)
		***REMOVED***
	***REMOVED***
***REMOVED***
