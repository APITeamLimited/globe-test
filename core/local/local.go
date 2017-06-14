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
	"github.com/pkg/errors"
	null "gopkg.in/guregu/null.v3"
)

type vuHandle struct ***REMOVED***
	VU     lib.VU
	Cancel context.CancelFunc
	Lock   sync.RWMutex

	runLock sync.Mutex
***REMOVED***

func (h *vuHandle) Run(ctx context.Context) error ***REMOVED***
	h.Lock.Lock()
	_, cancel := context.WithCancel(ctx)
	h.Cancel = cancel
	h.Lock.Unlock()

	return nil
***REMOVED***

type Executor struct ***REMOVED***
	Runner lib.Runner

	vus       []*vuHandle
	vusLock   sync.RWMutex
	numVUs    int64
	numVUsMax int64

	iterations, endIterations int64
	time, endTime             int64
	paused                    lib.AtomicBool

	runLock sync.Mutex
	ctx     context.Context
***REMOVED***

func New(r lib.Runner) *Executor ***REMOVED***
	return &Executor***REMOVED***Runner: r***REMOVED***
***REMOVED***

func (e *Executor) Run(ctx context.Context, out <-chan []stats.Sample) error ***REMOVED***
	e.runLock.Lock()
	defer e.runLock.Unlock()

	e.ctx = ctx
	<-ctx.Done()
	e.ctx = nil

	return nil
***REMOVED***

func (e *Executor) IsRunning() bool ***REMOVED***
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
	return e.paused.Get()
***REMOVED***

func (e *Executor) SetPaused(paused bool) ***REMOVED***
	e.paused.Set(paused)
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

	e.vusLock.Lock()
	defer e.vusLock.Unlock()

	for i, handle := range e.vus ***REMOVED***
		if i <= int(num) ***REMOVED***
			_, cancel := context.WithCancel(e.ctx)
			handle.Cancel = cancel
		***REMOVED*** else if handle.Cancel != nil ***REMOVED***
			handle.Cancel()
			handle.Cancel = nil
		***REMOVED***
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
