// Package events implements setInterval, setTimeout and co. Not to be used, mostly for testing purposes
package events

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

// RootModule is the global module instance that will create module
// instances for each VU.
type RootModule struct***REMOVED******REMOVED***

// Events represents an instance of the events module.
type Events struct ***REMOVED***
	vu modules.VU

	timerStopCounter uint32
	timerStopsLock   sync.Mutex
	timerStops       map[uint32]chan struct***REMOVED******REMOVED***
***REMOVED***

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &Events***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &Events***REMOVED***
		vu:         vu,
		timerStops: make(map[uint32]chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Exports returns the exports of the k6 module.
func (e *Events) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"setTimeout":    e.setTimeout,
			"clearTimeout":  e.clearTimeout,
			"setInterval":   e.setInterval,
			"clearInterval": e.clearInterval,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func noop() error ***REMOVED*** return nil ***REMOVED***

func (e *Events) getTimerStopCh() (uint32, chan struct***REMOVED******REMOVED***) ***REMOVED***
	id := atomic.AddUint32(&e.timerStopCounter, 1)
	ch := make(chan struct***REMOVED******REMOVED***)
	e.timerStopsLock.Lock()
	e.timerStops[id] = ch
	e.timerStopsLock.Unlock()
	return id, ch
***REMOVED***

func (e *Events) stopTimerCh(id uint32) bool ***REMOVED*** //nolint:unparam
	e.timerStopsLock.Lock()
	defer e.timerStopsLock.Unlock()
	ch, ok := e.timerStops[id]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	delete(e.timerStops, id)
	close(ch)
	return true
***REMOVED***

func (e *Events) call(callback goja.Callable, args []goja.Value) error ***REMOVED***
	// TODO: investigate, not sure GlobalObject() is always the correct value for `this`?
	_, err := callback(e.vu.Runtime().GlobalObject(), args...)
	return err
***REMOVED***

func (e *Events) setTimeout(callback goja.Callable, delay float64, args ...goja.Value) uint32 ***REMOVED***
	runOnLoop := e.vu.RegisterCallback()
	id, stopCh := e.getTimerStopCh()

	if delay < 0 ***REMOVED***
		delay = 0
	***REMOVED***

	go func() ***REMOVED***
		timer := time.NewTimer(time.Duration(delay * float64(time.Millisecond)))
		defer func() ***REMOVED***
			e.stopTimerCh(id)
			if !timer.Stop() ***REMOVED***
				<-timer.C
			***REMOVED***
		***REMOVED***()

		select ***REMOVED***
		case <-timer.C:
			runOnLoop(func() error ***REMOVED***
				return e.call(callback, args)
			***REMOVED***)
		case <-stopCh:
			runOnLoop(noop)
		case <-e.vu.Context().Done():
			e.vu.State().Logger.Warnf("setTimeout %d was stopped because the VU iteration was interrupted", id)
			runOnLoop(noop)
		***REMOVED***
	***REMOVED***()

	return id
***REMOVED***

func (e *Events) clearTimeout(id uint32) ***REMOVED***
	e.stopTimerCh(id)
***REMOVED***

func (e *Events) setInterval(callback goja.Callable, delay float64, args ...goja.Value) uint32 ***REMOVED***
	runOnLoop := e.vu.RegisterCallback()
	id, stopCh := e.getTimerStopCh()

	go func() ***REMOVED***
		ticker := time.NewTicker(time.Duration(delay * float64(time.Millisecond)))
		defer func() ***REMOVED***
			e.stopTimerCh(id)
			ticker.Stop()
		***REMOVED***()

		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				runOnLoop(func() error ***REMOVED***
					runOnLoop = e.vu.RegisterCallback()
					return e.call(callback, args)
				***REMOVED***)
			case <-stopCh:
				runOnLoop(noop)
				return
			case <-e.vu.Context().Done():
				e.vu.State().Logger.Warnf("setInterval %d was stopped because the VU iteration was interrupted", id)
				runOnLoop(noop)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return id
***REMOVED***

func (e *Events) clearInterval(id uint32) ***REMOVED***
	e.stopTimerCh(id)
***REMOVED***
