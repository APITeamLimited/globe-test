// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rate provides a rate limiter.
package rate

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// Limit defines the maximum frequency of some events.
// Limit is represented as number of events per second.
// A zero Limit allows no events.
type Limit float64

// Inf is the infinite rate limit; it allows all events (even if burst is zero).
const Inf = Limit(math.MaxFloat64)

// Every converts a minimum time interval between events to a Limit.
func Every(interval time.Duration) Limit ***REMOVED***
	if interval <= 0 ***REMOVED***
		return Inf
	***REMOVED***
	return 1 / Limit(interval.Seconds())
***REMOVED***

// A Limiter controls how frequently events are allowed to happen.
// It implements a "token bucket" of size b, initially full and refilled
// at rate r tokens per second.
// Informally, in any large enough time interval, the Limiter limits the
// rate to r tokens per second, with a maximum burst size of b events.
// As a special case, if r == Inf (the infinite rate), b is ignored.
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
//
// The zero value is a valid Limiter, but it will reject all events.
// Use NewLimiter to create non-zero Limiters.
//
// Limiter has three main methods, Allow, Reserve, and Wait.
// Most callers should use Wait.
//
// Each of the three methods consumes a single token.
// They differ in their behavior when no token is available.
// If no token is available, Allow returns false.
// If no token is available, Reserve returns a reservation for a future token
// and the amount of time the caller must wait before using it.
// If no token is available, Wait blocks until one can be obtained
// or its associated context.Context is canceled.
//
// The methods AllowN, ReserveN, and WaitN consume n tokens.
type Limiter struct ***REMOVED***
	mu     sync.Mutex
	limit  Limit
	burst  int
	tokens float64
	// last is the last time the limiter's tokens field was updated
	last time.Time
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time
***REMOVED***

// Limit returns the maximum overall event rate.
func (lim *Limiter) Limit() Limit ***REMOVED***
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
***REMOVED***

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (lim *Limiter) Burst() int ***REMOVED***
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
***REMOVED***

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewLimiter(r Limit, b int) *Limiter ***REMOVED***
	return &Limiter***REMOVED***
		limit: r,
		burst: b,
	***REMOVED***
***REMOVED***

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *Limiter) Allow() bool ***REMOVED***
	return lim.AllowN(time.Now(), 1)
***REMOVED***

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (lim *Limiter) AllowN(now time.Time, n int) bool ***REMOVED***
	return lim.reserveN(now, n, 0).ok
***REMOVED***

// A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
// A Reservation may be canceled, which may enable the Limiter to permit additional events.
type Reservation struct ***REMOVED***
	ok        bool
	lim       *Limiter
	tokens    int
	timeToAct time.Time
	// This is the Limit at reservation time, it can change later.
	limit Limit
***REMOVED***

// OK returns whether the limiter can provide the requested number of tokens
// within the maximum wait time.  If OK is false, Delay returns InfDuration, and
// Cancel does nothing.
func (r *Reservation) OK() bool ***REMOVED***
	return r.ok
***REMOVED***

// Delay is shorthand for DelayFrom(time.Now()).
func (r *Reservation) Delay() time.Duration ***REMOVED***
	return r.DelayFrom(time.Now())
***REMOVED***

// InfDuration is the duration returned by Delay when a Reservation is not OK.
const InfDuration = time.Duration(1<<63 - 1)

// DelayFrom returns the duration for which the reservation holder must wait
// before taking the reserved action.  Zero duration means act immediately.
// InfDuration means the limiter cannot grant the tokens requested in this
// Reservation within the maximum wait time.
func (r *Reservation) DelayFrom(now time.Time) time.Duration ***REMOVED***
	if !r.ok ***REMOVED***
		return InfDuration
	***REMOVED***
	delay := r.timeToAct.Sub(now)
	if delay < 0 ***REMOVED***
		return 0
	***REMOVED***
	return delay
***REMOVED***

// Cancel is shorthand for CancelAt(time.Now()).
func (r *Reservation) Cancel() ***REMOVED***
	r.CancelAt(time.Now())
***REMOVED***

// CancelAt indicates that the reservation holder will not perform the reserved action
// and reverses the effects of this Reservation on the rate limit as much as possible,
// considering that other reservations may have already been made.
func (r *Reservation) CancelAt(now time.Time) ***REMOVED***
	if !r.ok ***REMOVED***
		return
	***REMOVED***

	r.lim.mu.Lock()
	defer r.lim.mu.Unlock()

	if r.lim.limit == Inf || r.tokens == 0 || r.timeToAct.Before(now) ***REMOVED***
		return
	***REMOVED***

	// calculate tokens to restore
	// The duration between lim.lastEvent and r.timeToAct tells us how many tokens were reserved
	// after r was obtained. These tokens should not be restored.
	restoreTokens := float64(r.tokens) - r.limit.tokensFromDuration(r.lim.lastEvent.Sub(r.timeToAct))
	if restoreTokens <= 0 ***REMOVED***
		return
	***REMOVED***
	// advance time to now
	now, _, tokens := r.lim.advance(now)
	// calculate new number of tokens
	tokens += restoreTokens
	if burst := float64(r.lim.burst); tokens > burst ***REMOVED***
		tokens = burst
	***REMOVED***
	// update state
	r.lim.last = now
	r.lim.tokens = tokens
	if r.timeToAct == r.lim.lastEvent ***REMOVED***
		prevEvent := r.timeToAct.Add(r.limit.durationFromTokens(float64(-r.tokens)))
		if !prevEvent.Before(now) ***REMOVED***
			r.lim.lastEvent = prevEvent
		***REMOVED***
	***REMOVED***
***REMOVED***

// Reserve is shorthand for ReserveN(time.Now(), 1).
func (lim *Limiter) Reserve() *Reservation ***REMOVED***
	return lim.ReserveN(time.Now(), 1)
***REMOVED***

// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
// The Limiter takes this Reservation into account when allowing future events.
// The returned Reservation’s OK() method returns false if n exceeds the Limiter's burst size.
// Usage example:
//   r := lim.ReserveN(time.Now(), 1)
//   if !r.OK() ***REMOVED***
//     // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
//     return
//   ***REMOVED***
//   time.Sleep(r.Delay())
//   Act()
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation ***REMOVED***
	r := lim.reserveN(now, n, InfDuration)
	return &r
***REMOVED***

// Wait is shorthand for WaitN(ctx, 1).
func (lim *Limiter) Wait(ctx context.Context) (err error) ***REMOVED***
	return lim.WaitN(ctx, 1)
***REMOVED***

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error) ***REMOVED***
	lim.mu.Lock()
	burst := lim.burst
	limit := lim.limit
	lim.mu.Unlock()

	if n > burst && limit != Inf ***REMOVED***
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, burst)
	***REMOVED***
	// Check if ctx is already cancelled
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	default:
	***REMOVED***
	// Determine wait limit
	now := time.Now()
	waitLimit := InfDuration
	if deadline, ok := ctx.Deadline(); ok ***REMOVED***
		waitLimit = deadline.Sub(now)
	***REMOVED***
	// Reserve
	r := lim.reserveN(now, n, waitLimit)
	if !r.ok ***REMOVED***
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", n)
	***REMOVED***
	// Wait if necessary
	delay := r.DelayFrom(now)
	if delay == 0 ***REMOVED***
		return nil
	***REMOVED***
	t := time.NewTimer(delay)
	defer t.Stop()
	select ***REMOVED***
	case <-t.C:
		// We can proceed.
		return nil
	case <-ctx.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.Cancel()
		return ctx.Err()
	***REMOVED***
***REMOVED***

// SetLimit is shorthand for SetLimitAt(time.Now(), newLimit).
func (lim *Limiter) SetLimit(newLimit Limit) ***REMOVED***
	lim.SetLimitAt(time.Now(), newLimit)
***REMOVED***

// SetLimitAt sets a new Limit for the limiter. The new Limit, and Burst, may be violated
// or underutilized by those which reserved (using Reserve or Wait) but did not yet act
// before SetLimitAt was called.
func (lim *Limiter) SetLimitAt(now time.Time, newLimit Limit) ***REMOVED***
	lim.mu.Lock()
	defer lim.mu.Unlock()

	now, _, tokens := lim.advance(now)

	lim.last = now
	lim.tokens = tokens
	lim.limit = newLimit
***REMOVED***

// SetBurst is shorthand for SetBurstAt(time.Now(), newBurst).
func (lim *Limiter) SetBurst(newBurst int) ***REMOVED***
	lim.SetBurstAt(time.Now(), newBurst)
***REMOVED***

// SetBurstAt sets a new burst size for the limiter.
func (lim *Limiter) SetBurstAt(now time.Time, newBurst int) ***REMOVED***
	lim.mu.Lock()
	defer lim.mu.Unlock()

	now, _, tokens := lim.advance(now)

	lim.last = now
	lim.tokens = tokens
	lim.burst = newBurst
***REMOVED***

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *Limiter) reserveN(now time.Time, n int, maxFutureReserve time.Duration) Reservation ***REMOVED***
	lim.mu.Lock()

	if lim.limit == Inf ***REMOVED***
		lim.mu.Unlock()
		return Reservation***REMOVED***
			ok:        true,
			lim:       lim,
			tokens:    n,
			timeToAct: now,
		***REMOVED***
	***REMOVED***

	now, last, tokens := lim.advance(now)

	// Calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n)

	// Calculate the wait duration
	var waitDuration time.Duration
	if tokens < 0 ***REMOVED***
		waitDuration = lim.limit.durationFromTokens(-tokens)
	***REMOVED***

	// Decide result
	ok := n <= lim.burst && waitDuration <= maxFutureReserve

	// Prepare reservation
	r := Reservation***REMOVED***
		ok:    ok,
		lim:   lim,
		limit: lim.limit,
	***REMOVED***
	if ok ***REMOVED***
		r.tokens = n
		r.timeToAct = now.Add(waitDuration)
	***REMOVED***

	// Update state
	if ok ***REMOVED***
		lim.last = now
		lim.tokens = tokens
		lim.lastEvent = r.timeToAct
	***REMOVED*** else ***REMOVED***
		lim.last = last
	***REMOVED***

	lim.mu.Unlock()
	return r
***REMOVED***

// advance calculates and returns an updated state for lim resulting from the passage of time.
// lim is not changed.
// advance requires that lim.mu is held.
func (lim *Limiter) advance(now time.Time) (newNow time.Time, newLast time.Time, newTokens float64) ***REMOVED***
	last := lim.last
	if now.Before(last) ***REMOVED***
		last = now
	***REMOVED***

	// Calculate the new number of tokens, due to time that passed.
	elapsed := now.Sub(last)
	delta := lim.limit.tokensFromDuration(elapsed)
	tokens := lim.tokens + delta
	if burst := float64(lim.burst); tokens > burst ***REMOVED***
		tokens = burst
	***REMOVED***
	return now, last, tokens
***REMOVED***

// durationFromTokens is a unit conversion function from the number of tokens to the duration
// of time it takes to accumulate them at a rate of limit tokens per second.
func (limit Limit) durationFromTokens(tokens float64) time.Duration ***REMOVED***
	seconds := tokens / float64(limit)
	return time.Duration(float64(time.Second) * seconds)
***REMOVED***

// tokensFromDuration is a unit conversion function from a time duration to the number of tokens
// which could be accumulated during that duration at a rate of limit tokens per second.
func (limit Limit) tokensFromDuration(d time.Duration) float64 ***REMOVED***
	return d.Seconds() * float64(limit)
***REMOVED***
