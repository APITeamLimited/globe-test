/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package scheduler

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/loadimpact/k6/lib"
)

// This is a helper type used in executors where we have to dynamically control
// the number of VUs that are simultaneously running. For the moment, it is used in the VariableLoopingVUs and
//
// TODO: something simpler?
type vuHandle struct ***REMOVED***
	mutex     *sync.RWMutex
	parentCtx context.Context
	getVU     func() (lib.VU, error)
	returnVU  func(lib.VU)

	canStartIter chan struct***REMOVED******REMOVED***

	ctx    context.Context
	cancel func()
	logger *logrus.Entry
***REMOVED***

func newStoppedVUHandle(
	parentCtx context.Context, getVU func() (lib.VU, error), returnVU func(lib.VU), logger *logrus.Entry,
) *vuHandle ***REMOVED***
	lock := &sync.RWMutex***REMOVED******REMOVED***
	ctx, cancel := context.WithCancel(parentCtx)
	return &vuHandle***REMOVED***
		mutex:     lock,
		parentCtx: parentCtx,
		getVU:     getVU,
		returnVU:  returnVU,

		canStartIter: make(chan struct***REMOVED******REMOVED***),

		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	***REMOVED***
***REMOVED***

func (vh *vuHandle) start() ***REMOVED***
	vh.mutex.Lock()
	vh.logger.Debugf("Start")
	close(vh.canStartIter)
	vh.mutex.Unlock()
***REMOVED***

func (vh *vuHandle) gracefulStop() ***REMOVED***
	vh.mutex.Lock()
	select ***REMOVED***
	case <-vh.canStartIter:
		vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
		vh.logger.Debugf("Graceful stop")
	default:
		// do nothing, the signalling channel was already closed by either
		// hardStop() or gracefulStop()
	***REMOVED***
	vh.mutex.Unlock()
***REMOVED***

func (vh *vuHandle) hardStop() ***REMOVED***
	vh.mutex.Lock()
	vh.logger.Debugf("Hard stop")
	vh.cancel()                                          // cancel the previous context
	vh.ctx, vh.cancel = context.WithCancel(vh.parentCtx) // create a new context
	select ***REMOVED***                                             // if needed,
	case <-vh.canStartIter:
		vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
	default:
		// do nothing, the signalling channel was already closed by either
		// hardStop() or gracefulStop()
	***REMOVED***
	vh.mutex.Unlock()
***REMOVED***

//TODO: simplify this somehow - I feel like there should be a better way to
//implement this logic... maybe with sync.Cond?
func (vh *vuHandle) runLoopsIfPossible(runIter func(context.Context, lib.VU)) ***REMOVED***
	executorDone := vh.parentCtx.Done()

	var vu lib.VU
	defer func() ***REMOVED***
		if vu != nil ***REMOVED***
			vh.returnVU(vu)
			vu = nil
		***REMOVED***
	***REMOVED***()

mainLoop:
	for ***REMOVED***
		vh.mutex.RLock()
		canStartIter, ctx := vh.canStartIter, vh.ctx
		vh.mutex.RUnlock()

		// Wait for either the executor to be done, or for us to be unpaused
		select ***REMOVED***
		case <-canStartIter:
			// Best case, we're currently running, so we do nothing here, we
			// just continue straight ahead.
		case <-executorDone:
			return // The whole executor is done, nothing more to do.
		default:
			// We're not running, but the executor isn't done yet, so we wait
			// for either one of those conditions. But before that, we'll return
			// our VU to the pool, if we have it.
			if vu != nil ***REMOVED***
				vh.returnVU(vu)
				vu = nil
			***REMOVED***
			select ***REMOVED***
			case <-canStartIter:
				// continue on, we were unblocked...
			case <-ctx.Done():
				// hardStop was called, start a fresh iteration to get the new
				// context and signal channel
				continue mainLoop
			case <-executorDone:
				return // The whole executor is done, nothing more to do.
			***REMOVED***
		***REMOVED***

		// Probably not needed, but just in case - if both running and
		// executorDone were actice, check that the executor isn't done.
		select ***REMOVED***
		case <-executorDone:
			return
		default:
		***REMOVED***

		if vu == nil ***REMOVED*** // Ensure we have a VU
			freshVU, err := vh.getVU()
			if err != nil ***REMOVED***
				return
			***REMOVED***
			vu = freshVU
		***REMOVED***

		runIter(ctx, vu)
	***REMOVED***
***REMOVED***
