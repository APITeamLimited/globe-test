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
	"errors"
	"github.com/loadimpact/k6/stats"
	"sync"
	"time"
)

const (
	TickRate          = 1 * time.Millisecond
	CollectRate       = 10 * time.Millisecond
	ThresholdTickRate = 2 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

// Special error used to signal that a VU wants a taint, without logging an error.
var ErrVUWantsTaint = errors.New("test is tainted")

type vuEntry struct ***REMOVED***
	VU     VU
	Cancel context.CancelFunc

	Samples []stats.Sample
	lock    sync.Mutex
***REMOVED***

// The Engine is the beating heart of K6.
type Engine struct ***REMOVED***
	Runner  Runner
	Options Options

	Thresholds  map[string]Thresholds
	Metrics     map[*stats.Metric]stats.Sink
	MetricsLock sync.Mutex

	atTime    time.Duration
	vuEntries []*vuEntry
	vuMutex   sync.Mutex

	// Stubbing these out to pass tests.
	running bool
	paused  bool
	vus     int64
	vusMax  int64

	// Subsystem-related.
	subctx    context.Context
	subcancel context.CancelFunc
	submutex  sync.Mutex
	subwg     sync.WaitGroup
***REMOVED***

func NewEngine(r Runner, o Options) (*Engine, error) ***REMOVED***
	e := &Engine***REMOVED***
		Runner:  r,
		Options: o,

		Metrics:    make(map[*stats.Metric]stats.Sink),
		Thresholds: make(map[string]Thresholds),
	***REMOVED***
	e.clearSubcontext()

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

	return e, nil
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	go e.runCollection(ctx)

	lastTick := time.Time***REMOVED******REMOVED***
	ticker := time.NewTicker(TickRate)

	e.running = true
loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			break loop
		case <-ticker.C:
		***REMOVED***

		// Calculate the time delta between now and the last tick.
		now := time.Now()
		if lastTick.IsZero() ***REMOVED***
			lastTick = now
		***REMOVED***
		dT := now.Sub(lastTick)
		lastTick = now

		// Update the time counter appropriately.
		e.atTime += dT
	***REMOVED***
	e.running = false

	e.clearSubcontext()
	e.subwg.Wait()

	return nil
***REMOVED***

func (e *Engine) IsRunning() bool ***REMOVED***
	return e.running
***REMOVED***

func (e *Engine) SetPaused(v bool) ***REMOVED***
	e.paused = v
***REMOVED***

func (e *Engine) IsPaused() bool ***REMOVED***
	return e.paused
***REMOVED***

func (e *Engine) SetVUs(v int64) error ***REMOVED***
	if v < 0 ***REMOVED***
		return errors.New("vus can't be negative")
	***REMOVED***
	if v > e.vusMax ***REMOVED***
		return errors.New("more vus than allocated requested")
	***REMOVED***

	e.vuMutex.Lock()
	defer e.vuMutex.Unlock()

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
			e.subwg.Done()
			e.runVU(ctx, vu)
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
	return e.vus
***REMOVED***

func (e *Engine) SetVUsMax(v int64) error ***REMOVED***
	if v < 0 ***REMOVED***
		return errors.New("vus-max can't be negative")
	***REMOVED***
	if v < e.vus ***REMOVED***
		return errors.New("can't reduce vus-max below vus")
	***REMOVED***

	e.vuMutex.Lock()
	defer e.vuMutex.Unlock()

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
	return e.vusMax
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	return false
***REMOVED***

func (e *Engine) AtTime() time.Duration ***REMOVED***
	return e.atTime
***REMOVED***

func (e *Engine) TotalTime() (time.Duration, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (e *Engine) clearSubcontext() ***REMOVED***
	e.submutex.Lock()
	defer e.submutex.Unlock()

	if e.subcancel != nil ***REMOVED***
		e.subcancel()
	***REMOVED***
	subctx, subcancel := context.WithCancel(context.Background())
	e.subctx = subctx
	e.subcancel = subcancel
***REMOVED***

func (e *Engine) runVU(ctx context.Context, vu *vuEntry) ***REMOVED***
	// nil runners that produce nil VUs are used for testing.
	if vu.VU == nil ***REMOVED***
		<-ctx.Done()
		return
	***REMOVED***

	for ***REMOVED***
		samples, _ := vu.VU.RunOnce(ctx)

		vu.lock.Lock()
		vu.Samples = append(vu.Samples, samples...)
		vu.lock.Unlock()

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
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
	samples := []stats.Sample***REMOVED******REMOVED***
	for _, vu := range e.vuEntries ***REMOVED***
		if vu.Samples == nil ***REMOVED***
			continue
		***REMOVED***

		vu.lock.Lock()
		samples = append(samples, vu.Samples...)
		vu.Samples = nil
		vu.lock.Unlock()
	***REMOVED***
	return samples
***REMOVED***

func (e *Engine) processSamples(samples ...stats.Sample) ***REMOVED***
	e.MetricsLock.Lock()
	for _, sample := range samples ***REMOVED***
		sink := e.Metrics[sample.Metric]
		if sink == nil ***REMOVED***
			sink = sample.Metric.NewSink()
			e.Metrics[sample.Metric] = sink
		***REMOVED***
		sink.Add(sample)
	***REMOVED***
	e.MetricsLock.Unlock()
***REMOVED***
