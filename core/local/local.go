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
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

type vuHandle struct ***REMOVED***
	VU     lib.VU
	Cancel context.CancelFunc
	Lock   sync.RWMutex

	runLock sync.Mutex
***REMOVED***

type Executor struct ***REMOVED***
	Runner lib.Runner
	Logger *log.Logger

	runLock sync.Mutex
	wg      sync.WaitGroup

	vus       []*vuHandle
	vusLock   sync.RWMutex
	numVUs    int64
	numVUsMax int64

	iterations, endIterations int64
	time, endTime             int64

	pauseLock sync.RWMutex
	pause     chan interface***REMOVED******REMOVED***

	// Lock for: ctx, flow, out
	lock sync.RWMutex

	// Current context, nil if a test isn't running right now.
	ctx context.Context

	// Engineward output channel for samples.
	out chan<- []stats.Sample

	// Flow control for VUs; iterations are run only after reading from this channel.
	flow chan struct***REMOVED******REMOVED***
***REMOVED***

func New(r lib.Runner) *Executor ***REMOVED***
	return &Executor***REMOVED***
		Runner:        r,
		Logger:        log.StandardLogger(),
		endIterations: -1,
		endTime:       -1,
	***REMOVED***
***REMOVED***

func (e *Executor) Run(ctx context.Context, out chan<- []stats.Sample) error ***REMOVED***
	e.runLock.Lock()
	defer e.runLock.Unlock()

	e.lock.Lock()
	e.ctx = ctx
	e.out = out
	e.flow = make(chan struct***REMOVED******REMOVED***)
	e.lock.Unlock()

	e.scale(ctx, lib.Max(0, atomic.LoadInt64(&e.numVUs)))

	defer func() ***REMOVED***
		e.lock.Lock()
		e.ctx = nil
		e.out = nil
		e.flow = nil
		e.lock.Unlock()

		e.wg.Wait()
	***REMOVED***()

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
			leftovers := time.Since(lastTick)
			select ***REMOVED***
			case <-pause:
				lastTick = time.Now().Add(-leftovers)
			case <-ctx.Done():
				return nil
			***REMOVED***
		***REMOVED***

		select ***REMOVED***
		case t := <-ticker.C:
			d := t.Sub(lastTick)
			lastTick = t

			at := atomic.AddInt64(&e.time, int64(d))
			end := atomic.LoadInt64(&e.endTime)
			if end >= 0 && at >= end ***REMOVED***
				return nil
			***REMOVED***
		case e.flow <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			at := atomic.AddInt64(&e.iterations, 1)
			end := atomic.LoadInt64(&e.endIterations)
			if end >= 0 && at >= end ***REMOVED***
				return nil
			***REMOVED***
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Executor) runVU(ctx context.Context, handle *vuHandle) ***REMOVED***
	e.lock.RLock()
	flow := e.flow
	out := e.out
	e.lock.RUnlock()

	for range flow ***REMOVED***
		samples, err := handle.VU.RunOnce(ctx)
		if err != nil ***REMOVED***
			if s, ok := err.(fmt.Stringer); ok ***REMOVED***
				e.Logger.Error(s.String())
			***REMOVED*** else ***REMOVED***
				e.Logger.Error(err.Error())
			***REMOVED***
			continue
		***REMOVED***
		out <- samples
	***REMOVED***
***REMOVED***

func (e *Executor) scale(ctx context.Context, num int64) ***REMOVED***
	e.vusLock.Lock()
	defer e.vusLock.Unlock()

	for i, handle := range e.vus ***REMOVED***
		if i <= int(num) && handle.Cancel == nil ***REMOVED***
			ctx, cancel := context.WithCancel(ctx)
			handle.Cancel = cancel
			go e.runVU(ctx, handle)
		***REMOVED*** else if handle.Cancel != nil ***REMOVED***
			handle.Cancel()
			handle.Cancel = nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Executor) IsRunning() bool ***REMOVED***
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.ctx != nil
***REMOVED***

func (e *Executor) GetIterations() int64 ***REMOVED***
	return atomic.LoadInt64(&e.iterations)
***REMOVED***

func (e *Executor) GetEndIterations() null.Int ***REMOVED***
	v := atomic.LoadInt64(&e.endIterations)
	if v < 0 ***REMOVED***
		return null.Int***REMOVED******REMOVED***
	***REMOVED***
	return null.IntFrom(v)
***REMOVED***

func (e *Executor) SetEndIterations(i null.Int) ***REMOVED***
	if !i.Valid ***REMOVED***
		i.Int64 = -1
	***REMOVED***
	atomic.StoreInt64(&e.endIterations, i.Int64)
***REMOVED***

func (e *Executor) GetTime() time.Duration ***REMOVED***
	return time.Duration(atomic.LoadInt64(&e.time))
***REMOVED***

func (e *Executor) GetEndTime() lib.NullDuration ***REMOVED***
	v := atomic.LoadInt64(&e.endTime)
	if v < 0 ***REMOVED***
		return lib.NullDuration***REMOVED******REMOVED***
	***REMOVED***
	return lib.NullDurationFrom(time.Duration(v))
***REMOVED***

func (e *Executor) SetEndTime(t lib.NullDuration) ***REMOVED***
	if !t.Valid ***REMOVED***
		t.Duration = -1
	***REMOVED***
	atomic.StoreInt64(&e.endTime, int64(t.Duration))
***REMOVED***

func (e *Executor) IsPaused() bool ***REMOVED***
	e.pauseLock.RLock()
	defer e.pauseLock.RUnlock()
	return e.pause != nil
***REMOVED***

func (e *Executor) SetPaused(paused bool) ***REMOVED***
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

	if numVUsMax := atomic.LoadInt64(&e.numVUsMax); num > numVUsMax ***REMOVED***
		return errors.Errorf("can't raise vu count (to %d) above vu cap (%d)", num, numVUsMax)
	***REMOVED***

	e.lock.Lock()
	defer e.lock.Unlock()

	if ctx := e.ctx; ctx != nil ***REMOVED***
		e.scale(ctx, num)
	***REMOVED***

	atomic.StoreInt64(&e.numVUs, num)

	return nil
***REMOVED***

func (e *Executor) GetVUsMax() int64 ***REMOVED***
	return atomic.LoadInt64(&e.numVUsMax)
***REMOVED***

func (e *Executor) SetVUsMax(max int64) error ***REMOVED***
	if max < 0 ***REMOVED***
		return errors.New("vu cap can't be negative")
	***REMOVED***

	if numVUs := atomic.LoadInt64(&e.numVUs); max < numVUs ***REMOVED***
		return errors.Errorf("can't lower vu cap (to %d) below vu count (%d)", max, numVUs)
	***REMOVED***

	numVUsMax := atomic.LoadInt64(&e.numVUsMax)

	if max < numVUsMax ***REMOVED***
		e.vus = e.vus[:max]
		atomic.StoreInt64(&e.numVUsMax, max)
		return nil
	***REMOVED***

	e.vusLock.Lock()
	defer e.vusLock.Unlock()

	vus := e.vus
	for i := numVUsMax; i < max; i++ ***REMOVED***
		var handle vuHandle
		if e.Runner != nil ***REMOVED***
			vu, err := e.Runner.NewVU()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			handle.VU = vu
		***REMOVED***
		vus = append(vus, &handle)
	***REMOVED***
	e.vus = vus

	atomic.StoreInt64(&e.numVUsMax, max)

	return nil
***REMOVED***
