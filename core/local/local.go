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
	"sync"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	null "gopkg.in/guregu/null.v3"
)

type VUHandle struct ***REMOVED***
	VU      lib.VU
	Cancel  context.CancelFunc
	Samples []stats.Sample
***REMOVED***

type Executor struct ***REMOVED***
	Runner lib.Runner
	VUs    []VUHandle

	runLock   sync.Mutex
	isRunning bool

	iterations, endIterations int64
	time, endTime             int64
	paused                    lib.AtomicBool
	vus, vusMax               int64
***REMOVED***

func New(r lib.Runner) *Executor ***REMOVED***
	return &Executor***REMOVED***Runner: r***REMOVED***
***REMOVED***

func (e *Executor) Run(ctx context.Context, out <-chan []stats.Sample) error ***REMOVED***
	e.runLock.Lock()
	e.isRunning = true
	defer func() ***REMOVED***
		e.isRunning = false
		e.runLock.Unlock()
	***REMOVED***()

	<-ctx.Done()
	return nil
***REMOVED***

func (e *Executor) IsRunning() bool ***REMOVED***
	return e.isRunning
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
	return e.paused.Get()
***REMOVED***

func (e *Executor) SetPaused(paused bool) ***REMOVED***
	e.paused.Set(paused)
***REMOVED***

func (e *Executor) GetVUs() int64 ***REMOVED***
	return atomic.LoadInt64(&e.vus)
***REMOVED***

func (e *Executor) SetVUs(vus int64) error ***REMOVED***
	atomic.StoreInt64(&e.vus, vus)
	return nil
***REMOVED***

func (e *Executor) GetVUsMax() int64 ***REMOVED***
	return atomic.LoadInt64(&e.vusMax)
***REMOVED***

func (e *Executor) SetVUsMax(max int64) error ***REMOVED***
	atomic.StoreInt64(&e.vusMax, max)
	return nil
***REMOVED***
