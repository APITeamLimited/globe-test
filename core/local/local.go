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

package local

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

// TODO: totally rewrite this!
// This is an overcomplicated and probably buggy piece of code that is a major PITA to refactor...
// It does a ton of stuff in a very convoluted way, has a and uses a very incomprehensible mix
// of all possible Go synchronization mechanisms (channels, mutexes, rwmutexes, atomics,
// and waitgroups) and has a bunch of contexts and tickers on top...

var _ lib.Executor = &Executor***REMOVED******REMOVED***

type vuHandle struct ***REMOVED***
	sync.RWMutex
	vu     lib.VU
	ctx    context.Context
	cancel context.CancelFunc
***REMOVED***

func (h *vuHandle) run(logger *log.Logger, flow <-chan int64, iterDone chan<- struct***REMOVED******REMOVED***) ***REMOVED***
	h.RLock()
	ctx := h.ctx
	h.RUnlock()

	for ***REMOVED***
		select ***REMOVED***
		case _, ok := <-flow:
			if !ok ***REMOVED***
				return
			***REMOVED***
		case <-ctx.Done():
			return
		***REMOVED***

		if h.vu != nil ***REMOVED***
			err := h.vu.RunOnce(ctx)
			select ***REMOVED***
			case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			default:
				if err != nil ***REMOVED***
					if s, ok := err.(fmt.Stringer); ok ***REMOVED***
						logger.Error(s.String())
					***REMOVED*** else ***REMOVED***
						logger.Error(err.Error())
					***REMOVED***
				***REMOVED***
				iterDone <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			iterDone <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type Executor struct ***REMOVED***
	Runner lib.Runner
	Logger *log.Logger

	runLock sync.Mutex
	wg      sync.WaitGroup

	runSetup    bool
	runTeardown bool

	vus       []*vuHandle
	vusLock   sync.RWMutex
	numVUs    int64
	numVUsMax int64
	nextVUID  int64

	iters     int64 // Completed iterations
	partIters int64 // Partial, incomplete iterations
	endIters  int64 // End test at this many iterations

	time    int64 // Current time
	endTime int64 // End test at this timestamp

	pauseLock sync.RWMutex
	pause     chan interface***REMOVED******REMOVED***

	stages []lib.Stage

	// Lock for: ctx, flow, out
	lock sync.RWMutex

	// Current context, nil if a test isn't running right now.
	ctx context.Context

	// Output channel to which VUs send samples.
	vuOut chan stats.SampleContainer

	// Channel on which VUs sigal that iterations are completed
	iterDone chan struct***REMOVED******REMOVED***

	// Flow control for VUs; iterations are run only after reading from this channel.
	flow chan int64
***REMOVED***

func New(r lib.Runner) *Executor ***REMOVED***
	var bufferSize int64
	if r != nil ***REMOVED***
		bufferSize = r.GetOptions().MetricSamplesBufferSize.Int64
	***REMOVED***

	return &Executor***REMOVED***
		Runner:      r,
		Logger:      log.StandardLogger(),
		runSetup:    true,
		runTeardown: true,
		endIters:    -1,
		endTime:     -1,
		vuOut:       make(chan stats.SampleContainer, bufferSize),
		iterDone:    make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

func (e *Executor) Run(parent context.Context, engineOut chan<- stats.SampleContainer) (reterr error) ***REMOVED***
	e.runLock.Lock()
	defer e.runLock.Unlock()

	if e.Runner != nil && e.runSetup ***REMOVED***
		if err := e.Runner.Setup(parent, engineOut); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	ctx, cancel := context.WithCancel(parent)
	vuFlow := make(chan int64)
	e.lock.Lock()
	vuOut := e.vuOut
	iterDone := e.iterDone
	e.ctx = ctx
	e.flow = vuFlow
	e.lock.Unlock()

	var cutoff time.Time
	defer func() ***REMOVED***
		if e.Runner != nil && e.runTeardown ***REMOVED***
			err := e.Runner.Teardown(parent, engineOut)
			if reterr == nil ***REMOVED***
				reterr = err
			***REMOVED*** else if err != nil ***REMOVED***
				reterr = fmt.Errorf("Teardown error %#v\nPrevious error: %#v", err, reterr)
			***REMOVED***
		***REMOVED***

		close(vuFlow)
		cancel()

		e.lock.Lock()
		e.ctx = nil
		e.vuOut = nil
		e.flow = nil
		e.lock.Unlock()

		wait := make(chan interface***REMOVED******REMOVED***)
		go func() ***REMOVED***
			e.wg.Wait()
			close(wait)
		***REMOVED***()

		for ***REMOVED***
			select ***REMOVED***
			case <-iterDone:
				// Spool through all remaining iterations, do not emit stats since the Run() is over
			case newSampleContainer := <-vuOut:
				if cutoff.IsZero() ***REMOVED***
					engineOut <- newSampleContainer
				***REMOVED*** else if csc, ok := newSampleContainer.(stats.ConnectedSampleContainer); ok && csc.GetTime().Before(cutoff) ***REMOVED***
					engineOut <- newSampleContainer
				***REMOVED*** else ***REMOVED***
					for _, s := range newSampleContainer.GetSamples() ***REMOVED***
						if s.Time.Before(cutoff) ***REMOVED***
							engineOut <- s
						***REMOVED***
					***REMOVED***
				***REMOVED***
			case <-wait:
			***REMOVED***
			select ***REMOVED***
			case <-wait:
				close(vuOut)
				return
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	startVUs := atomic.LoadInt64(&e.numVUs)
	if err := e.scale(ctx, lib.Max(0, startVUs)); err != nil ***REMOVED***
		return err
	***REMOVED***

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	lastTick := time.Now()
	for ***REMOVED***
		// If the test is paused, sleep until either the pause or the test ends.
		// Also shift the last tick to omit time spent paused, but not partial ticks.
		e.pauseLock.RLock()
		pause := e.pause
		e.pauseLock.RUnlock()
		if pause != nil ***REMOVED***
			e.Logger.Debug("Local: Pausing!")
			leftovers := time.Since(lastTick)
			select ***REMOVED***
			case <-pause:
				e.Logger.Debug("Local: No longer paused")
				lastTick = time.Now().Add(-leftovers)
			case <-ctx.Done():
				e.Logger.Debug("Local: Terminated while in paused state")
				return nil
			***REMOVED***
		***REMOVED***

		// Dumb hack: we don't wanna start any more iterations than the max, but we can't
		// conditionally select on a channel either...so, we cheat: swap out the flow channel for a
		// nil channel (writing to nil always blocks) if we don't wanna write an iteration.
		flow := vuFlow
		end := atomic.LoadInt64(&e.endIters)
		partials := atomic.LoadInt64(&e.partIters)
		if end >= 0 && partials >= end ***REMOVED***
			flow = nil
		***REMOVED***

		select ***REMOVED***
		case flow <- partials:
			// Start an iteration if there's a VU waiting. See also: the big comment block above.
			atomic.AddInt64(&e.partIters, 1)
		case t := <-ticker.C:
			// Every tick, increment the clock, see if we passed the end point, and process stages.
			// If the test ends this way, set a cutoff point; any samples collected past the cutoff
			// point are excluded.
			d := t.Sub(lastTick)
			lastTick = t

			end := time.Duration(atomic.LoadInt64(&e.endTime))
			at := time.Duration(atomic.AddInt64(&e.time, int64(d)))
			if end >= 0 && at >= end ***REMOVED***
				e.Logger.WithFields(log.Fields***REMOVED***"at": at, "end": end***REMOVED***).Debug("Local: Hit time limit")
				cutoff = time.Now()
				return nil
			***REMOVED***

			stages := e.stages
			if len(stages) > 0 ***REMOVED***
				vus, keepRunning := ProcessStages(startVUs, stages, at)
				if !keepRunning ***REMOVED***
					e.Logger.WithField("at", at).Debug("Local: Ran out of stages")
					cutoff = time.Now()
					return nil
				***REMOVED***
				if vus.Valid ***REMOVED***
					if err := e.SetVUs(vus.Int64); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case sampleContainer := <-vuOut:
			engineOut <- sampleContainer
		case <-iterDone:
			// Every iteration ends with a write to iterDone. Check if we've hit the end point.
			// If not, make sure to include an Iterations bump in the list!
			var tags *stats.SampleTags
			if e.Runner != nil ***REMOVED***
				tags = e.Runner.GetOptions().RunTags
			***REMOVED***
			engineOut <- stats.Sample***REMOVED***
				Time:   time.Now(),
				Metric: metrics.Iterations,
				Value:  1,
				Tags:   tags,
			***REMOVED***

			end := atomic.LoadInt64(&e.endIters)
			at := atomic.AddInt64(&e.iters, 1)
			if end >= 0 && at >= end ***REMOVED***
				e.Logger.WithFields(log.Fields***REMOVED***"at": at, "end": end***REMOVED***).Debug("Local: Hit iteration limit")
				return nil
			***REMOVED***
		case <-ctx.Done():
			// If the test is cancelled, just set the cutoff point to now and proceed down the same
			// logic as if the time limit was hit.
			e.Logger.Debug("Local: Exiting with context")
			cutoff = time.Now()
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Executor) scale(ctx context.Context, num int64) error ***REMOVED***
	e.Logger.WithField("num", num).Debug("Local: Scaling...")

	e.vusLock.Lock()
	defer e.vusLock.Unlock()

	e.lock.RLock()
	flow := e.flow
	iterDone := e.iterDone
	e.lock.RUnlock()

	for i, handle := range e.vus ***REMOVED***
		handle := handle
		handle.RLock()
		cancel := handle.cancel
		handle.RUnlock()

		if i < int(num) ***REMOVED***
			if cancel == nil ***REMOVED***
				vuctx, cancel := context.WithCancel(ctx)
				handle.Lock()
				handle.ctx = vuctx
				handle.cancel = cancel
				handle.Unlock()

				if handle.vu != nil ***REMOVED***
					if err := handle.vu.Reconfigure(atomic.AddInt64(&e.nextVUID, 1)); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***

				e.wg.Add(1)
				go func() ***REMOVED***
					handle.run(e.Logger, flow, iterDone)
					e.wg.Done()
				***REMOVED***()
			***REMOVED***
		***REMOVED*** else if cancel != nil ***REMOVED***
			handle.Lock()
			handle.cancel()
			handle.cancel = nil
			handle.Unlock()
		***REMOVED***
	***REMOVED***

	atomic.StoreInt64(&e.numVUs, num)
	return nil
***REMOVED***

func (e *Executor) IsRunning() bool ***REMOVED***
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.ctx != nil
***REMOVED***

func (e *Executor) GetRunner() lib.Runner ***REMOVED***
	return e.Runner
***REMOVED***

func (e *Executor) SetLogger(l *log.Logger) ***REMOVED***
	e.Logger = l
***REMOVED***

func (e *Executor) GetLogger() *log.Logger ***REMOVED***
	return e.Logger
***REMOVED***

func (e *Executor) GetStages() []lib.Stage ***REMOVED***
	return e.stages
***REMOVED***

func (e *Executor) SetStages(s []lib.Stage) ***REMOVED***
	e.stages = s
***REMOVED***

func (e *Executor) GetIterations() int64 ***REMOVED***
	return atomic.LoadInt64(&e.iters)
***REMOVED***

func (e *Executor) GetEndIterations() null.Int ***REMOVED***
	v := atomic.LoadInt64(&e.endIters)
	if v < 0 ***REMOVED***
		return null.Int***REMOVED******REMOVED***
	***REMOVED***
	return null.IntFrom(v)
***REMOVED***

func (e *Executor) SetEndIterations(i null.Int) ***REMOVED***
	if !i.Valid ***REMOVED***
		i.Int64 = -1
	***REMOVED***
	e.Logger.WithField("i", i.Int64).Debug("Local: Setting end iterations")
	atomic.StoreInt64(&e.endIters, i.Int64)
***REMOVED***

func (e *Executor) GetTime() time.Duration ***REMOVED***
	return time.Duration(atomic.LoadInt64(&e.time))
***REMOVED***

func (e *Executor) GetEndTime() types.NullDuration ***REMOVED***
	v := atomic.LoadInt64(&e.endTime)
	if v < 0 ***REMOVED***
		return types.NullDuration***REMOVED******REMOVED***
	***REMOVED***
	return types.NullDurationFrom(time.Duration(v))
***REMOVED***

func (e *Executor) SetEndTime(t types.NullDuration) ***REMOVED***
	if !t.Valid ***REMOVED***
		t.Duration = -1
	***REMOVED***
	e.Logger.WithField("d", t.Duration).Debug("Local: Setting end time")
	atomic.StoreInt64(&e.endTime, int64(t.Duration))
***REMOVED***

func (e *Executor) IsPaused() bool ***REMOVED***
	e.pauseLock.RLock()
	defer e.pauseLock.RUnlock()
	return e.pause != nil
***REMOVED***

func (e *Executor) SetPaused(paused bool) ***REMOVED***
	e.Logger.WithField("paused", paused).Debug("Local: Setting paused")
	e.pauseLock.Lock()
	defer e.pauseLock.Unlock()

	if paused && e.pause == nil ***REMOVED***
		e.pause = make(chan interface***REMOVED******REMOVED***)
	***REMOVED*** else if !paused && e.pause != nil ***REMOVED***
		close(e.pause)
		e.pause = nil
	***REMOVED***
***REMOVED***

func (e *Executor) GetVUs() int64 ***REMOVED***
	return atomic.LoadInt64(&e.numVUs)
***REMOVED***

func (e *Executor) SetVUs(num int64) error ***REMOVED***
	if num < 0 ***REMOVED***
		return errors.New("vu count can't be negative")
	***REMOVED***

	if atomic.LoadInt64(&e.numVUs) == num ***REMOVED***
		return nil
	***REMOVED***

	e.Logger.WithField("vus", num).Debug("Local: Setting VUs")

	if numVUsMax := atomic.LoadInt64(&e.numVUsMax); num > numVUsMax ***REMOVED***
		return errors.Errorf("can't raise vu count (to %d) above vu cap (%d)", num, numVUsMax)
	***REMOVED***

	if ctx := e.ctx; ctx != nil ***REMOVED***
		if err := e.scale(ctx, num); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		atomic.StoreInt64(&e.numVUs, num)
	***REMOVED***

	return nil
***REMOVED***

func (e *Executor) GetVUsMax() int64 ***REMOVED***
	return atomic.LoadInt64(&e.numVUsMax)
***REMOVED***

func (e *Executor) SetVUsMax(max int64) error ***REMOVED***
	e.Logger.WithField("max", max).Debug("Local: Setting max VUs")
	if max < 0 ***REMOVED***
		return errors.New("vu cap can't be negative")
	***REMOVED***

	numVUsMax := atomic.LoadInt64(&e.numVUsMax)

	if numVUsMax == max ***REMOVED***
		return nil
	***REMOVED***

	if numVUs := atomic.LoadInt64(&e.numVUs); max < numVUs ***REMOVED***
		return errors.Errorf("can't lower vu cap (to %d) below vu count (%d)", max, numVUs)
	***REMOVED***

	if max < numVUsMax ***REMOVED***
		e.vus = e.vus[:max]
		atomic.StoreInt64(&e.numVUsMax, max)
		return nil
	***REMOVED***

	e.lock.RLock()
	vuOut := e.vuOut
	e.lock.RUnlock()

	e.vusLock.Lock()
	defer e.vusLock.Unlock()

	vus := e.vus
	for i := numVUsMax; i < max; i++ ***REMOVED***
		var handle vuHandle
		if e.Runner != nil ***REMOVED***
			vu, err := e.Runner.NewVU(vuOut)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			handle.vu = vu
		***REMOVED***
		vus = append(vus, &handle)
	***REMOVED***
	e.vus = vus

	atomic.StoreInt64(&e.numVUsMax, max)

	return nil
***REMOVED***

func (e *Executor) SetRunSetup(r bool) ***REMOVED***
	e.runSetup = r
***REMOVED***

func (e *Executor) SetRunTeardown(r bool) ***REMOVED***
	e.runTeardown = r
***REMOVED***
