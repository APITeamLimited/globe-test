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
	ThresholdTickRate = 2 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

// Special error used to signal that a VU wants a taint, without logging an error.
var ErrVUWantsTaint = errors.New("test is tainted")

// The Engine is the beating heart of K6.
type Engine struct ***REMOVED***
	Runner  Runner
	Options Options

	Metrics    map[*stats.Metric]stats.Sink
	Thresholds map[string]Thresholds

	atTime time.Duration

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
	e.running = true
	<-ctx.Done()
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
