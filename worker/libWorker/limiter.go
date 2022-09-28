package libWorker

import (
	"sync"
)

// SlotLimiter can restrict the concurrent execution of tasks to the given `slots` limit
type SlotLimiter chan struct***REMOVED******REMOVED***

// NewSlotLimiter initializes and returns a new SlotLimiter with the given slot count
func NewSlotLimiter(slots int) SlotLimiter ***REMOVED***
	if slots <= 0 ***REMOVED***
		return nil
	***REMOVED***

	ch := make(chan struct***REMOVED******REMOVED***, slots)
	for i := 0; i < slots; i++ ***REMOVED***
		ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return ch
***REMOVED***

// Begin uses up a slot to denote the start of a task exeuction. It's a noop if the number
// of slots is 0, and if no slots are available, it blocks and waits.
func (sl SlotLimiter) Begin() ***REMOVED***
	if sl != nil ***REMOVED***
		<-sl
	***REMOVED***
***REMOVED***

// End restores a slot and should be called at the end of a taks execution, preferably
// from a defer statement right after Begin()
func (sl SlotLimiter) End() ***REMOVED***
	if sl != nil ***REMOVED***
		sl <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// MultiSlotLimiter can restrict the concurrent execution of different groups of tasks
// to the given `slots` limit. Each group is represented with a string ID.
type MultiSlotLimiter struct ***REMOVED***
	m     map[string]SlotLimiter
	slots int
	mutex sync.Mutex
***REMOVED***

// NewMultiSlotLimiter initializes and returns a new MultiSlotLimiter with the given slot count
// TODO: move to lib and use something better than a mutex? sync.Map perhaps?
func NewMultiSlotLimiter(slots int) *MultiSlotLimiter ***REMOVED***
	return &MultiSlotLimiter***REMOVED***make(map[string]SlotLimiter), slots, sync.Mutex***REMOVED******REMOVED******REMOVED***
***REMOVED***

// Slot is used to retrieve the corresponding slot to the given string ID. If no slot with that ID exists,
// it creates it and saves it for future use. It is safe to call this method concurrently.
func (l *MultiSlotLimiter) Slot(s string) SlotLimiter ***REMOVED***
	if l.slots == 0 ***REMOVED***
		return nil
	***REMOVED***

	l.mutex.Lock()
	defer l.mutex.Unlock()

	ll, ok := l.m[s]
	if !ok ***REMOVED***
		tmp := NewSlotLimiter(l.slots)
		ll = tmp
		l.m[s] = ll
	***REMOVED***
	return ll
***REMOVED***
