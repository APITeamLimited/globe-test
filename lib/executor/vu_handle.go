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

package executor

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/loadimpact/k6/lib"
)

// This is a helper type used in executors where we have to dynamically control
// the number of VUs that are simultaneously running. For the moment, it is used
// in the VariableLoopingVUs and the ExternallyControlled executors.
//
// TODO: something simpler?
type vuHandle struct ***REMOVED***
	mutex     *sync.RWMutex
	parentCtx context.Context
	getVU     func() (lib.InitializedVU, error)
	returnVU  func(lib.InitializedVU)
	exec      string
	env       map[string]string

	canStartIter chan struct***REMOVED******REMOVED***

	ctx    context.Context
	cancel func()
	logger *logrus.Entry
***REMOVED***

func newStoppedVUHandle(
	parentCtx context.Context, getVU func() (lib.InitializedVU, error),
	returnVU func(lib.InitializedVU), exec string, env map[string]string,
	logger *logrus.Entry,
) *vuHandle ***REMOVED***
	lock := &sync.RWMutex***REMOVED******REMOVED***
	ctx, cancel := context.WithCancel(parentCtx)
	return &vuHandle***REMOVED***
		mutex:     lock,
		parentCtx: parentCtx,
		getVU:     getVU,
		returnVU:  returnVU,
		exec:      exec,
		env:       env,

		canStartIter: make(chan struct***REMOVED******REMOVED***),

		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	***REMOVED***
***REMOVED***

func (vh *vuHandle) start() ***REMOVED***
	vh.mutex.Lock()
	vh.logger.Debug("Start")
	close(vh.canStartIter)
	vh.mutex.Unlock()
***REMOVED***

func (vh *vuHandle) gracefulStop() ***REMOVED***
	vh.mutex.Lock()
	select ***REMOVED***
	case <-vh.canStartIter:
		vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
		vh.logger.Debug("Graceful stop")
	default:
		// do nothing, the signalling channel was already initialized by hardStop()
	***REMOVED***
	vh.mutex.Unlock()
***REMOVED***

func (vh *vuHandle) hardStop() ***REMOVED***
	vh.mutex.Lock()
	vh.logger.Debug("Hard stop")
	vh.cancel()                                          // cancel the previous context
	vh.ctx, vh.cancel = context.WithCancel(vh.parentCtx) // create a new context
	select ***REMOVED***
	case <-vh.canStartIter:
		vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
	default:
		// do nothing, the signalling channel was already initialized by gracefulStop()
	***REMOVED***
	vh.mutex.Unlock()
***REMOVED***

//TODO: simplify this somehow - I feel like there should be a better way to
//implement this logic... maybe with sync.Cond?
func (vh *vuHandle) runLoopsIfPossible(runIter func(context.Context, lib.ActiveVU)) ***REMOVED***
	executorDone := vh.parentCtx.Done()

	var vu lib.ActiveVU
	var deactivateVU func()

mainLoop:
	for ***REMOVED***
		vh.mutex.RLock()
		canStartIter, ctx := vh.canStartIter, vh.ctx
		vh.mutex.RUnlock()

		// Wait for either the executor to be done, or for us to be un-paused
		select ***REMOVED***
		case <-canStartIter:
			// Best case, we're currently running, so we do nothing here, we
			// just continue straight ahead.
		case <-executorDone:
			// The whole executor is done, nothing more to do.
			return
		default:
			// We're not running, but the executor isn't done yet, so we wait
			// for either one of those conditions. But before that, clear
			// the VU reference to ensure we get a fresh one below.
			vu = nil
			select ***REMOVED***
			case <-canStartIter:
				// continue on, we were unblocked...
			case <-ctx.Done():
				// hardStop was called, start a fresh iteration to get the new
				// context and signal channel
				continue mainLoop
			case <-executorDone:
				// The whole executor is done, nothing more to do.
				return
			***REMOVED***
		***REMOVED***

		// Probably not needed, but just in case - if both running and
		// executorDone were active, check that the executor isn't done.
		select ***REMOVED***
		case <-executorDone:
			return
		default:
		***REMOVED***

		// Ensure we have an active VU
		if vu == nil ***REMOVED***
			initVU, err := vh.getVU()
			if err != nil ***REMOVED***
				return
			***REMOVED***
			deactivateVU = func() ***REMOVED***
				vh.returnVU(initVU)
			***REMOVED***
			vu = initVU.Activate(&lib.VUActivationParams***REMOVED***
				Exec:               vh.exec,
				RunContext:         ctx,
				Env:                vh.env,
				DeactivateCallback: deactivateVU,
			***REMOVED***)
		***REMOVED***

		runIter(ctx, vu)
	***REMOVED***
***REMOVED***
