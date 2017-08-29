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
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
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

// The Engine is the beating heart of K6.
type Engine struct ***REMOVED***
	runLock sync.Mutex

	Executor  lib.Executor
	Options   lib.Options
	Collector lib.Collector

	logger *log.Logger

	Stages      []lib.Stage
	Metrics     map[string]*stats.Metric
	MetricsLock sync.RWMutex

	// Assigned to metrics upon first received sample.
	thresholds map[string]stats.Thresholds
	submetrics map[string][]*stats.Submetric

	// Are thresholds tainted?
	thresholdsTainted bool
***REMOVED***

func NewEngine(ex lib.Executor, o lib.Options) (*Engine, error) ***REMOVED***
	if ex == nil ***REMOVED***
		ex = local.New(nil)
	***REMOVED***

	e := &Engine***REMOVED***
		Executor: ex,
		Options:  o,
		Metrics:  make(map[string]*stats.Metric),
	***REMOVED***
	e.SetLogger(log.StandardLogger())

	if err := ex.SetVUsMax(o.VUsMax.Int64); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := ex.SetVUs(o.VUs.Int64); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ex.SetPaused(o.Paused.Bool)

	// Use Stages if available, if not, construct a stage to fill the specified duration.
	// Special case: A valid duration of 0 = an infinite (invalid duration) stage.
	if o.Stages != nil ***REMOVED***
		e.Stages = o.Stages
	***REMOVED*** else if o.Duration.Valid && o.Duration.Duration > 0 ***REMOVED***
		e.Stages = []lib.Stage***REMOVED******REMOVED***Duration: o.Duration***REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		e.Stages = []lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	ex.SetEndTime(SumStages(e.Stages))
	ex.SetEndIterations(o.Iterations)

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

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.runLock.Lock()
	defer e.runLock.Unlock()

	e.logger.Debug("Engine: Starting with parameters...")
	for i, st := range e.Stages ***REMOVED***
		fields := make(log.Fields)
		if st.Target.Valid ***REMOVED***
			fields["tgt"] = st.Target.Int64
		***REMOVED***
		if st.Duration.Valid ***REMOVED***
			fields["d"] = st.Duration.Duration
		***REMOVED***
		e.logger.WithFields(fields).Debugf(" - stage #%d", i)
	***REMOVED***

	fields := make(log.Fields)
	if endTime := e.Executor.GetEndTime(); endTime.Valid ***REMOVED***
		fields["time"] = endTime.Duration
	***REMOVED***
	if endIter := e.Executor.GetEndIterations(); endIter.Valid ***REMOVED***
		fields["iter"] = endIter.Int64
	***REMOVED***
	e.logger.WithFields(fields).Debug(" - end conditions (if any)")

	collectorwg := sync.WaitGroup***REMOVED******REMOVED***
	collectorctx, collectorcancel := context.WithCancel(context.Background())
	if e.Collector != nil ***REMOVED***
		collectorwg.Add(1)
		go func() ***REMOVED***
			e.Collector.Run(collectorctx)
			collectorwg.Done()
		***REMOVED***()
		for !e.Collector.IsReady() ***REMOVED***
			runtime.Gosched()
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
	subwg.Add(1)
	go func() ***REMOVED***
		e.runThresholds(subctx)
		e.logger.Debug("Engine: Thresholds terminated")
		subwg.Done()
	***REMOVED***()

	// Run the executor.
	out := make(chan []stats.Sample)
	errC := make(chan error)
	subwg.Add(1)
	go func() ***REMOVED***
		errC <- e.Executor.Run(subctx, out)
		e.logger.Debug("Engine: Executor terminated")
		subwg.Done()
	***REMOVED***()

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
			close(out)
		***REMOVED***()
		for samples := range out ***REMOVED***
			e.processSamples(samples...)
		***REMOVED***

		// Emit final metrics.
		e.emitMetrics()

		// Process final thresholds.
		e.processThresholds()

		// Finally, shut down collector.
		collectorcancel()
		collectorwg.Wait()
	***REMOVED***()

	ticker := time.NewTicker(TickRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			vus, keepRunning := ProcessStages(e.Stages, e.Executor.GetTime())
			if !keepRunning ***REMOVED***
				e.logger.Debug("run: ProcessStages() returned false; exiting...")
				return nil
			***REMOVED***
			if vus.Valid ***REMOVED***
				if err := e.Executor.SetVUs(vus.Int64); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		case samples := <-out:
			e.processSamples(samples...)
		case err := <-errC:
			errC = nil
			if err != nil ***REMOVED***
				e.logger.WithError(err).Debug("run: executor returned an error")
				return err
			***REMOVED***
			e.logger.Debug("run: executor terminated")
			return nil
		case <-ctx.Done():
			e.logger.Debug("run: context expired; exiting...")
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	return e.thresholdsTainted
***REMOVED***

func (e *Engine) SetLogger(l *log.Logger) ***REMOVED***
	e.logger = l
	e.Executor.SetLogger(l)
***REMOVED***

func (e *Engine) GetLogger() *log.Logger ***REMOVED***
	return e.logger
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
	e.processSamples(
		stats.Sample***REMOVED***
			Time:   t,
			Metric: metrics.VUs,
			Value:  float64(e.Executor.GetVUs()),
		***REMOVED***,
		stats.Sample***REMOVED***
			Time:   t,
			Metric: metrics.VUsMax,
			Value:  float64(e.Executor.GetVUsMax()),
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
	for _, m := range e.Metrics ***REMOVED***
		if len(m.Thresholds.Thresholds) == 0 ***REMOVED***
			continue
		***REMOVED***
		m.Tainted = null.BoolFrom(false)

		e.logger.WithField("m", m.Name).Debug("running thresholds")
		succ, err := m.Thresholds.Run(m.Sink)
		if err != nil ***REMOVED***
			e.logger.WithField("m", m.Name).WithError(err).Error("Threshold error")
			continue
		***REMOVED***
		if !succ ***REMOVED***
			e.logger.WithField("m", m.Name).Debug("Thresholds failed")
			m.Tainted = null.BoolFrom(true)
			e.thresholdsTainted = true
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) processSamples(samples ...stats.Sample) ***REMOVED***
	if len(samples) == 0 ***REMOVED***
		return
	***REMOVED***

	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	for _, sample := range samples ***REMOVED***
		m, ok := e.Metrics[sample.Metric.Name]
		if !ok ***REMOVED***
			m = sample.Metric
			m.Thresholds = e.thresholds[m.Name]
			m.Submetrics = e.submetrics[m.Name]
			e.Metrics[m.Name] = m
		***REMOVED***
		m.Sink.Add(sample)

		for _, sm := range m.Submetrics ***REMOVED***
			passing := true
			for k, v := range sm.Tags ***REMOVED***
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
				sm.Metric.Sub = *sm
				sm.Metric.Thresholds = e.thresholds[sm.Name]
				e.Metrics[sm.Name] = sm.Metric
			***REMOVED***
			sm.Metric.Sink.Add(sample)
		***REMOVED***
	***REMOVED***

	if e.Collector != nil ***REMOVED***
		e.Collector.Collect(samples)
	***REMOVED***
***REMOVED***
