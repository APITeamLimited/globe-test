package executor

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/lib"
)

type stateType int32

// states
const (
	stopped stateType = iota
	starting
	running
	toGracefulStop
	toHardStop
)

/*
the below is a state transition table (https://en.wikipedia.org/wiki/State-transition_table)
short names for input:
- start is the method start
- loop is a loop of runLoopsIfPossible
- grace is the method gracefulStop
- hard is the method hardStop
+-------+-------------------------------------+---------------------------------------------------+
| input | current         | next state        | notes                                             |
+-------+-------------------------------------+---------------------------------------------------+
| start | stopped         | starting          | normal                                            |
| start | starting        | starting          | nothing                                           |
| start | running         | running           | nothing                                           |
| start | toGracefulStop  | running           | we raced with the loop stopping, just continue    |
| start | toHardStop      | starting          | same as stopped really                            |
| loop  | stopped         | stopped           | we actually are blocked on canStartIter           |
| loop  | starting        | running           | get new VU and context                            |
| loop  | running         | running           | usually fast path                                 |
| loop  | toGracefulStop  | stopped           | cancel the context and make new one               |
| loop  | toHardStop      | stopped           | cancel the context and make new one               |
| grace | stopped         | stopped           | nothing                                           |
| grace | starting        | stopped           | cancel the context to return the VU               |
| grace | running         | toGracefulStop    | normal one, the actual work is in the loop        |
| grace | toGracefulStop  | toGracefulStop    | nothing                                           |
| grace | toHardSTop      | toHardStop        | nothing                                           |
| hard  | stopped         | stopped           | nothing                                           |
| hard  | starting        | stopped           | short circuit as in the grace case, not necessary |
| hard  | running         | toHardStop        | normal, cancel context and reinitialize it        |
| hard  | toGracefulStop  | toHardStop        | normal, cancel context and reinitialize it        |
| hard  | toHardStop      | toHardStop        | nothing                                           |
+-------+-----------------+-------------------+----------------------------------------------------+
*/

// This is a helper type used in executors where we have to dynamically control
// the number of VUs that are simultaneously running. For the moment, it is used
// in the RampingVUs and the ExternallyControlled executors.
// Notes on the implementation requirements:
// - it needs to be able to start and stop VUs in thread safe fashion
// - for each call to getVU there must be 1 (and only 1) call to returnVU
// - gracefulStop must let an iteration which has started to finish. For reasons of ease of
// implementation and lack of good evidence it's not required to let a not started iteration to
// finish in other words if you call start and then gracefulStop, there is no requirement for
// 1 iteration to have been started.
// - hardStop must stop an iteration in process
// - it's not required but preferable, if where possible to not reactivate VUs and to reuse context
// as this speed ups the execution
type vuHandle struct ***REMOVED***
	mutex                 *sync.Mutex
	parentCtx             context.Context
	getVU                 func() (lib.InitializedVU, error)
	returnVU              func(lib.InitializedVU)
	nextIterationCounters func() (uint64, uint64)
	config                *BaseConfig

	initVU       lib.InitializedVU
	activeVU     lib.ActiveVU
	canStartIter chan struct***REMOVED******REMOVED***

	state stateType // see the table above for meanings
	// stateH []int32 // helper for debugging

	ctx    context.Context
	cancel func()
	logger *logrus.Entry
***REMOVED***

func newStoppedVUHandle(
	parentCtx context.Context, getVU func() (lib.InitializedVU, error),
	returnVU func(lib.InitializedVU),
	nextIterationCounters func() (uint64, uint64),
	config *BaseConfig, logger *logrus.Entry,
) *vuHandle ***REMOVED***
	ctx, cancel := context.WithCancel(parentCtx)

	return &vuHandle***REMOVED***
		mutex:                 &sync.Mutex***REMOVED******REMOVED***,
		parentCtx:             parentCtx,
		getVU:                 getVU,
		nextIterationCounters: nextIterationCounters,
		config:                config,

		canStartIter: make(chan struct***REMOVED******REMOVED***),
		state:        stopped,

		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
		returnVU: returnVU,
	***REMOVED***
***REMOVED***

func (vh *vuHandle) start() (err error) ***REMOVED***
	vh.mutex.Lock()
	defer vh.mutex.Unlock()

	switch vh.state ***REMOVED***
	case starting, running:
		return nil // nothing to do
	case toGracefulStop: // we raced with the loop, lets not return the vu just to get it back
		vh.logger.Debug("Start")
		close(vh.canStartIter)
		vh.changeState(running)
	case stopped, toHardStop: // we need to reactivate the VU and remake the context for it
		vh.logger.Debug("Start")
		vh.initVU, err = vh.getVU()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		vh.activeVU = vh.initVU.Activate(getVUActivationParams(
			vh.ctx, *vh.config, vh.returnVU, vh.nextIterationCounters))
		close(vh.canStartIter)
		vh.changeState(starting)
	***REMOVED***
	return nil
***REMOVED***

// just a helper function for debugging
func (vh *vuHandle) changeState(newState stateType) ***REMOVED***
	// vh.stateH = append(vh.stateH, newState)
	atomic.StoreInt32((*int32)(&vh.state), int32(newState))
***REMOVED***

func (vh *vuHandle) gracefulStop() ***REMOVED***
	vh.mutex.Lock()
	defer vh.mutex.Unlock()
	switch vh.state ***REMOVED***
	case toGracefulStop, toHardStop, stopped:
		return // nothing to do
	case starting: // we raced with the loop and apparently it won't do a single iteration
		vh.cancel()
		vh.ctx, vh.cancel = context.WithCancel(vh.parentCtx)
		vh.changeState(stopped)
	case running:
		vh.changeState(toGracefulStop)
	***REMOVED***

	vh.logger.Debug("Graceful stop")
	vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

func (vh *vuHandle) hardStop() ***REMOVED***
	vh.mutex.Lock()
	defer vh.mutex.Unlock()

	switch vh.state ***REMOVED***
	case toHardStop, stopped:
		return // nothing to do
	case starting: // we raced with the loop and apparently it won't do a single iteration
		vh.changeState(stopped)
	case running, toGracefulStop:
		vh.changeState(toHardStop)
	***REMOVED***
	vh.logger.Debug("Hard stop")
	vh.cancel()
	vh.ctx, vh.cancel = context.WithCancel(vh.parentCtx)
	vh.canStartIter = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

// runLoopsIfPossible is where all the fun is :D. Unfortunately somewhere we need to check most
// of the cases and this is where this happens.
func (vh *vuHandle) runLoopsIfPossible(runIter func(context.Context, lib.ActiveVU) bool) ***REMOVED***
	// We can probably initialize here, but it's also easier to just use the slow path in the second
	// part of the for loop
	defer func() ***REMOVED***
		// not sure if this is needed, because here the parentCtx is canceled and I would argue it doesn't matter
		// if we set the correct state
		vh.mutex.Lock()
		vh.changeState(stopped)
		vh.mutex.Unlock()
	***REMOVED***()

	var (
		executorDone = vh.parentCtx.Done()
		ctx          context.Context
		cancel       func()
		vu           lib.ActiveVU
	)

	for ***REMOVED***
		state := stateType(atomic.LoadInt32((*int32)(&vh.state)))
		if state == running && runIter(ctx, vu) ***REMOVED*** // fast path
			continue
		***REMOVED***

		// slow path - something has changed - get what and wait until we can do more iterations
		vh.mutex.Lock()
		select ***REMOVED***
		case <-executorDone:
			// The whole executor is done, nothing more to do.
			vh.mutex.Unlock()
			return
		default:
		***REMOVED***

		switch vh.state ***REMOVED***
		case running: // start raced us toGracefulStop
			vh.mutex.Unlock()
			continue
		case toGracefulStop:
			if cancel != nil ***REMOVED***
				// we need to cancel the context, to return the vu
				// and because *we* did, lets reinitialize it
				cancel()
				vh.ctx, vh.cancel = context.WithCancel(vh.parentCtx)
			***REMOVED***
			fallthrough // to set the state
		case toHardStop:
			// we have *now* stopped
			vh.changeState(stopped)
		case stopped, starting:
			// there is nothing to do
		***REMOVED***

		canStartIter := vh.canStartIter
		ctx = vh.ctx
		vh.mutex.Unlock()

		// We're are stopped, but the executor isn't done yet, so we wait
		// for either one of those conditions.
		select ***REMOVED***
		case <-canStartIter: // we can start again
			vh.mutex.Lock()
			select ***REMOVED***
			case <-vh.canStartIter: // we check again in case of race
				// reinitialize
				vu, ctx, cancel = vh.activeVU, vh.ctx, vh.cancel
				vh.changeState(running)
			default:
				// well we got raced to here by something ... loop again ...
			***REMOVED***
			vh.mutex.Unlock()
		case <-ctx.Done():
			// hardStop was called, start a fresh iteration to get the new
			// context and signal channel
		case <-executorDone:
			// The whole executor is done, nothing more to do.
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
