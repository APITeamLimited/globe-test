package js

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

// eventLoop implements an event with
// handling of unhandled rejected promises.
//
// A specific thing about this event loop is that it will wait to return
// not only until the queue is empty but until nothing is registered that it will run in the future.
// This is in contrast with more common behaviours where it only returns on
// a specific event/action or when the loop is empty.
// This is required as in k6 iterations (for which event loop will be primary used)
// are supposed to be independent and any work started in them needs to finish,
// but also they need to end when all the instructions are done.
// Additionally because of this on any error while the event loop will exit it's
// required to wait on the event loop to be empty before the execution can continue.
type eventLoop struct ***REMOVED***
	lock                sync.Mutex
	queue               []func() error
	wakeupCh            chan struct***REMOVED******REMOVED*** // TODO: maybe use sync.Cond ?
	registeredCallbacks int
	vu                  modules.VU

	// pendingPromiseRejections are rejected promises with no handler,
	// if there is something in this map at an end of an event loop then it will exit with an error.
	// It's similar to what Deno and Node do.
	pendingPromiseRejections map[*goja.Promise]struct***REMOVED******REMOVED***
***REMOVED***

// newEventLoop returns a new event loop with a few helpers attached to it:
// - reporting (and aborting on) unhandled promise rejections
func newEventLoop(vu modules.VU) *eventLoop ***REMOVED***
	e := &eventLoop***REMOVED***
		wakeupCh:                 make(chan struct***REMOVED******REMOVED***, 1),
		pendingPromiseRejections: make(map[*goja.Promise]struct***REMOVED******REMOVED***),
		vu:                       vu,
	***REMOVED***
	vu.Runtime().SetPromiseRejectionTracker(e.promiseRejectionTracker)

	return e
***REMOVED***

func (e *eventLoop) wakeup() ***REMOVED***
	select ***REMOVED***
	case e.wakeupCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***
***REMOVED***

// registerCallback register that a callback will be invoked on the loop, preventing it from returning/finishing.
// The returned function, upon invocation, will queue its argument and wakeup the loop if needed.
// If the eventLoop has since stopped, it will not be executed.
// This function *must* be called from within running on the event loop, but its result can be called from anywhere.
func (e *eventLoop) registerCallback() func(func() error) ***REMOVED***
	e.lock.Lock()
	e.registeredCallbacks++
	e.lock.Unlock()

	return func(f func() error) ***REMOVED***
		e.lock.Lock()
		e.queue = append(e.queue, f)
		e.registeredCallbacks--
		e.lock.Unlock()
		e.wakeup()
	***REMOVED***
***REMOVED***

func (e *eventLoop) promiseRejectionTracker(p *goja.Promise, op goja.PromiseRejectionOperation) ***REMOVED***
	// No locking necessary here as the goja runtime will call this synchronously
	// Read Notes on https://tc39.es/ecma262/#sec-host-promise-rejection-tracker
	if op == goja.PromiseRejectionReject ***REMOVED***
		e.pendingPromiseRejections[p] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED*** // goja.PromiseRejectionHandle so a promise that was previously rejected without handler now got one
		delete(e.pendingPromiseRejections, p)
	***REMOVED***
***REMOVED***

func (e *eventLoop) popAll() (queue []func() error, awaiting bool) ***REMOVED***
	e.lock.Lock()
	queue = e.queue
	e.queue = make([]func() error, 0, len(queue))
	awaiting = e.registeredCallbacks != 0
	e.lock.Unlock()
	return
***REMOVED***

// start will run the event loop until it's empty and there are no uninvoked registered callbacks
// or a queued function returns an error. The provided firstCallback will be the first thing executed.
// After start returns the event loop can be reused as long as waitOnRegistered is called.
func (e *eventLoop) start(firstCallback func() error) error ***REMOVED***
	e.queue = []func() error***REMOVED***firstCallback***REMOVED***
	for ***REMOVED***
		queue, awaiting := e.popAll()

		if len(queue) == 0 ***REMOVED***
			if !awaiting ***REMOVED***
				return nil
			***REMOVED***
			<-e.wakeupCh
			continue
		***REMOVED***

		for _, f := range queue ***REMOVED***
			if err := f(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// This will get a random unhandled rejection instead of the first one, for example.
		// But that seems to be the case in other tools as well so it seems to not be that big of a problem.
		for promise := range e.pendingPromiseRejections ***REMOVED***
			// TODO maybe throw the whole promise up and get make a better message outside of the event loop
			value := promise.Result()
			if o := value.ToObject(e.vu.Runtime()); o != nil ***REMOVED***
				stack := o.Get("stack")
				if stack != nil ***REMOVED***
					value = stack
				***REMOVED***
			***REMOVED***
			// this is the de facto wording in both firefox and deno at least
			return fmt.Errorf("Uncaught (in promise) %s", value) //nolint:stylecheck
		***REMOVED***
	***REMOVED***
***REMOVED***

// Wait on all registered callbacks so we know nothing is still doing work.
func (e *eventLoop) waitOnRegistered() ***REMOVED***
	for ***REMOVED***
		_, awaiting := e.popAll()
		if !awaiting ***REMOVED***
			return
		***REMOVED***
		<-e.wakeupCh
	***REMOVED***
***REMOVED***
