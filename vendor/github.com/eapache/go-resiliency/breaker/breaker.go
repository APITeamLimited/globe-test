// Package breaker implements the circuit-breaker resiliency pattern for Go.
package breaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrBreakerOpen is the error returned from Run() when the function is not executed
// because the breaker is currently open.
var ErrBreakerOpen = errors.New("circuit breaker is open")

const (
	closed uint32 = iota
	open
	halfOpen
)

// Breaker implements the circuit-breaker resiliency pattern
type Breaker struct ***REMOVED***
	errorThreshold, successThreshold int
	timeout                          time.Duration

	lock              sync.Mutex
	state             uint32
	errors, successes int
	lastError         time.Time
***REMOVED***

// New constructs a new circuit-breaker that starts closed.
// From closed, the breaker opens if "errorThreshold" errors are seen
// without an error-free period of at least "timeout". From open, the
// breaker half-closes after "timeout". From half-open, the breaker closes
// after "successThreshold" consecutive successes, or opens on a single error.
func New(errorThreshold, successThreshold int, timeout time.Duration) *Breaker ***REMOVED***
	return &Breaker***REMOVED***
		errorThreshold:   errorThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	***REMOVED***
***REMOVED***

// Run will either return ErrBreakerOpen immediately if the circuit-breaker is
// already open, or it will run the given function and pass along its return
// value. It is safe to call Run concurrently on the same Breaker.
func (b *Breaker) Run(work func() error) error ***REMOVED***
	state := atomic.LoadUint32(&b.state)

	if state == open ***REMOVED***
		return ErrBreakerOpen
	***REMOVED***

	return b.doWork(state, work)
***REMOVED***

// Go will either return ErrBreakerOpen immediately if the circuit-breaker is
// already open, or it will run the given function in a separate goroutine.
// If the function is run, Go will return nil immediately, and will *not* return
// the return value of the function. It is safe to call Go concurrently on the
// same Breaker.
func (b *Breaker) Go(work func() error) error ***REMOVED***
	state := atomic.LoadUint32(&b.state)

	if state == open ***REMOVED***
		return ErrBreakerOpen
	***REMOVED***

	// errcheck complains about ignoring the error return value, but
	// that's on purpose; if you want an error from a goroutine you have to
	// get it over a channel or something
	go b.doWork(state, work)

	return nil
***REMOVED***

func (b *Breaker) doWork(state uint32, work func() error) error ***REMOVED***
	var panicValue interface***REMOVED******REMOVED***

	result := func() error ***REMOVED***
		defer func() ***REMOVED***
			panicValue = recover()
		***REMOVED***()
		return work()
	***REMOVED***()

	if result == nil && panicValue == nil && state == closed ***REMOVED***
		// short-circuit the normal, success path without contending
		// on the lock
		return nil
	***REMOVED***

	// oh well, I guess we have to contend on the lock
	b.processResult(result, panicValue)

	if panicValue != nil ***REMOVED***
		// as close as Go lets us come to a "rethrow" although unfortunately
		// we lose the original panicing location
		panic(panicValue)
	***REMOVED***

	return result
***REMOVED***

func (b *Breaker) processResult(result error, panicValue interface***REMOVED******REMOVED***) ***REMOVED***
	b.lock.Lock()
	defer b.lock.Unlock()

	if result == nil && panicValue == nil ***REMOVED***
		if b.state == halfOpen ***REMOVED***
			b.successes++
			if b.successes == b.successThreshold ***REMOVED***
				b.closeBreaker()
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b.errors > 0 ***REMOVED***
			expiry := b.lastError.Add(b.timeout)
			if time.Now().After(expiry) ***REMOVED***
				b.errors = 0
			***REMOVED***
		***REMOVED***

		switch b.state ***REMOVED***
		case closed:
			b.errors++
			if b.errors == b.errorThreshold ***REMOVED***
				b.openBreaker()
			***REMOVED*** else ***REMOVED***
				b.lastError = time.Now()
			***REMOVED***
		case halfOpen:
			b.openBreaker()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *Breaker) openBreaker() ***REMOVED***
	b.changeState(open)
	go b.timer()
***REMOVED***

func (b *Breaker) closeBreaker() ***REMOVED***
	b.changeState(closed)
***REMOVED***

func (b *Breaker) timer() ***REMOVED***
	time.Sleep(b.timeout)

	b.lock.Lock()
	defer b.lock.Unlock()

	b.changeState(halfOpen)
***REMOVED***

func (b *Breaker) changeState(newState uint32) ***REMOVED***
	b.errors = 0
	b.successes = 0
	atomic.StoreUint32(&b.state, newState)
***REMOVED***
