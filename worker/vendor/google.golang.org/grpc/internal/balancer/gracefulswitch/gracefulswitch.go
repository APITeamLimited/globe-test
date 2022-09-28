/*
 *
 * Copyright 2022 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package gracefulswitch implements a graceful switch load balancer.
package gracefulswitch

import (
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

var errBalancerClosed = errors.New("gracefulSwitchBalancer is closed")
var _ balancer.Balancer = (*Balancer)(nil)

// NewBalancer returns a graceful switch Balancer.
func NewBalancer(cc balancer.ClientConn, opts balancer.BuildOptions) *Balancer ***REMOVED***
	return &Balancer***REMOVED***
		cc:    cc,
		bOpts: opts,
	***REMOVED***
***REMOVED***

// Balancer is a utility to gracefully switch from one balancer to
// a new balancer. It implements the balancer.Balancer interface.
type Balancer struct ***REMOVED***
	bOpts balancer.BuildOptions
	cc    balancer.ClientConn

	// mu protects the following fields and all fields within balancerCurrent
	// and balancerPending. mu does not need to be held when calling into the
	// child balancers, as all calls into these children happen only as a direct
	// result of a call into the gracefulSwitchBalancer, which are also
	// guaranteed to be synchronous. There is one exception: an UpdateState call
	// from a child balancer when current and pending are populated can lead to
	// calling Close() on the current. To prevent that racing with an
	// UpdateSubConnState from the channel, we hold currentMu during Close and
	// UpdateSubConnState calls.
	mu              sync.Mutex
	balancerCurrent *balancerWrapper
	balancerPending *balancerWrapper
	closed          bool // set to true when this balancer is closed

	// currentMu must be locked before mu. This mutex guards against this
	// sequence of events: UpdateSubConnState() called, finds the
	// balancerCurrent, gives up lock, updateState comes in, causes Close() on
	// balancerCurrent before the UpdateSubConnState is called on the
	// balancerCurrent.
	currentMu sync.Mutex
***REMOVED***

// swap swaps out the current lb with the pending lb and updates the ClientConn.
// The caller must hold gsb.mu.
func (gsb *Balancer) swap() ***REMOVED***
	gsb.cc.UpdateState(gsb.balancerPending.lastState)
	cur := gsb.balancerCurrent
	gsb.balancerCurrent = gsb.balancerPending
	gsb.balancerPending = nil
	go func() ***REMOVED***
		gsb.currentMu.Lock()
		defer gsb.currentMu.Unlock()
		cur.Close()
	***REMOVED***()
***REMOVED***

// Helper function that checks if the balancer passed in is current or pending.
// The caller must hold gsb.mu.
func (gsb *Balancer) balancerCurrentOrPending(bw *balancerWrapper) bool ***REMOVED***
	return bw == gsb.balancerCurrent || bw == gsb.balancerPending
***REMOVED***

// SwitchTo initializes the graceful switch process, which completes based on
// connectivity state changes on the current/pending balancer. Thus, the switch
// process is not complete when this method returns. This method must be called
// synchronously alongside the rest of the balancer.Balancer methods this
// Graceful Switch Balancer implements.
func (gsb *Balancer) SwitchTo(builder balancer.Builder) error ***REMOVED***
	gsb.mu.Lock()
	if gsb.closed ***REMOVED***
		gsb.mu.Unlock()
		return errBalancerClosed
	***REMOVED***
	bw := &balancerWrapper***REMOVED***
		gsb: gsb,
		lastState: balancer.State***REMOVED***
			ConnectivityState: connectivity.Connecting,
			Picker:            base.NewErrPicker(balancer.ErrNoSubConnAvailable),
		***REMOVED***,
		subconns: make(map[balancer.SubConn]bool),
	***REMOVED***
	balToClose := gsb.balancerPending // nil if there is no pending balancer
	if gsb.balancerCurrent == nil ***REMOVED***
		gsb.balancerCurrent = bw
	***REMOVED*** else ***REMOVED***
		gsb.balancerPending = bw
	***REMOVED***
	gsb.mu.Unlock()
	balToClose.Close()
	// This function takes a builder instead of a balancer because builder.Build
	// can call back inline, and this utility needs to handle the callbacks.
	newBalancer := builder.Build(bw, gsb.bOpts)
	if newBalancer == nil ***REMOVED***
		// This is illegal and should never happen; we clear the balancerWrapper
		// we were constructing if it happens to avoid a potential panic.
		gsb.mu.Lock()
		if gsb.balancerPending != nil ***REMOVED***
			gsb.balancerPending = nil
		***REMOVED*** else ***REMOVED***
			gsb.balancerCurrent = nil
		***REMOVED***
		gsb.mu.Unlock()
		return balancer.ErrBadResolverState
	***REMOVED***

	// This write doesn't need to take gsb.mu because this field never gets read
	// or written to on any calls from the current or pending. Calls from grpc
	// to this balancer are guaranteed to be called synchronously, so this
	// bw.Balancer field will never be forwarded to until this SwitchTo()
	// function returns.
	bw.Balancer = newBalancer
	return nil
***REMOVED***

// Returns nil if the graceful switch balancer is closed.
func (gsb *Balancer) latestBalancer() *balancerWrapper ***REMOVED***
	gsb.mu.Lock()
	defer gsb.mu.Unlock()
	if gsb.balancerPending != nil ***REMOVED***
		return gsb.balancerPending
	***REMOVED***
	return gsb.balancerCurrent
***REMOVED***

// UpdateClientConnState forwards the update to the latest balancer created.
func (gsb *Balancer) UpdateClientConnState(state balancer.ClientConnState) error ***REMOVED***
	// The resolver data is only relevant to the most recent LB Policy.
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil ***REMOVED***
		return errBalancerClosed
	***REMOVED***
	// Perform this call without gsb.mu to prevent deadlocks if the child calls
	// back into the channel. The latest balancer can never be closed during a
	// call from the channel, even without gsb.mu held.
	return balToUpdate.UpdateClientConnState(state)
***REMOVED***

// ResolverError forwards the error to the latest balancer created.
func (gsb *Balancer) ResolverError(err error) ***REMOVED***
	// The resolver data is only relevant to the most recent LB Policy.
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil ***REMOVED***
		return
	***REMOVED***
	// Perform this call without gsb.mu to prevent deadlocks if the child calls
	// back into the channel. The latest balancer can never be closed during a
	// call from the channel, even without gsb.mu held.
	balToUpdate.ResolverError(err)
***REMOVED***

// ExitIdle forwards the call to the latest balancer created.
//
// If the latest balancer does not support ExitIdle, the subConns are
// re-connected to manually.
func (gsb *Balancer) ExitIdle() ***REMOVED***
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil ***REMOVED***
		return
	***REMOVED***
	// There is no need to protect this read with a mutex, as the write to the
	// Balancer field happens in SwitchTo, which completes before this can be
	// called.
	if ei, ok := balToUpdate.Balancer.(balancer.ExitIdler); ok ***REMOVED***
		ei.ExitIdle()
		return
	***REMOVED***
	gsb.mu.Lock()
	defer gsb.mu.Unlock()
	for sc := range balToUpdate.subconns ***REMOVED***
		sc.Connect()
	***REMOVED***
***REMOVED***

// UpdateSubConnState forwards the update to the appropriate child.
func (gsb *Balancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) ***REMOVED***
	gsb.currentMu.Lock()
	defer gsb.currentMu.Unlock()
	gsb.mu.Lock()
	// Forward update to the appropriate child.  Even if there is a pending
	// balancer, the current balancer should continue to get SubConn updates to
	// maintain the proper state while the pending is still connecting.
	var balToUpdate *balancerWrapper
	if gsb.balancerCurrent != nil && gsb.balancerCurrent.subconns[sc] ***REMOVED***
		balToUpdate = gsb.balancerCurrent
	***REMOVED*** else if gsb.balancerPending != nil && gsb.balancerPending.subconns[sc] ***REMOVED***
		balToUpdate = gsb.balancerPending
	***REMOVED***
	gsb.mu.Unlock()
	if balToUpdate == nil ***REMOVED***
		// SubConn belonged to a stale lb policy that has not yet fully closed,
		// or the balancer was already closed.
		return
	***REMOVED***
	balToUpdate.UpdateSubConnState(sc, state)
***REMOVED***

// Close closes any active child balancers.
func (gsb *Balancer) Close() ***REMOVED***
	gsb.mu.Lock()
	gsb.closed = true
	currentBalancerToClose := gsb.balancerCurrent
	gsb.balancerCurrent = nil
	pendingBalancerToClose := gsb.balancerPending
	gsb.balancerPending = nil
	gsb.mu.Unlock()

	currentBalancerToClose.Close()
	pendingBalancerToClose.Close()
***REMOVED***

// balancerWrapper wraps a balancer.Balancer, and overrides some Balancer
// methods to help cleanup SubConns created by the wrapped balancer.
//
// It implements the balancer.ClientConn interface and is passed down in that
// capacity to the wrapped balancer. It maintains a set of subConns created by
// the wrapped balancer and calls from the latter to create/update/remove
// SubConns update this set before being forwarded to the parent ClientConn.
// State updates from the wrapped balancer can result in invocation of the
// graceful switch logic.
type balancerWrapper struct ***REMOVED***
	balancer.Balancer
	gsb *Balancer

	lastState balancer.State
	subconns  map[balancer.SubConn]bool // subconns created by this balancer
***REMOVED***

func (bw *balancerWrapper) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) ***REMOVED***
	if state.ConnectivityState == connectivity.Shutdown ***REMOVED***
		bw.gsb.mu.Lock()
		delete(bw.subconns, sc)
		bw.gsb.mu.Unlock()
	***REMOVED***
	// There is no need to protect this read with a mutex, as the write to the
	// Balancer field happens in SwitchTo, which completes before this can be
	// called.
	bw.Balancer.UpdateSubConnState(sc, state)
***REMOVED***

// Close closes the underlying LB policy and removes the subconns it created. bw
// must not be referenced via balancerCurrent or balancerPending in gsb when
// called. gsb.mu must not be held.  Does not panic with a nil receiver.
func (bw *balancerWrapper) Close() ***REMOVED***
	// before Close is called.
	if bw == nil ***REMOVED***
		return
	***REMOVED***
	// There is no need to protect this read with a mutex, as Close() is
	// impossible to be called concurrently with the write in SwitchTo(). The
	// callsites of Close() for this balancer in Graceful Switch Balancer will
	// never be called until SwitchTo() returns.
	bw.Balancer.Close()
	bw.gsb.mu.Lock()
	for sc := range bw.subconns ***REMOVED***
		bw.gsb.cc.RemoveSubConn(sc)
	***REMOVED***
	bw.gsb.mu.Unlock()
***REMOVED***

func (bw *balancerWrapper) UpdateState(state balancer.State) ***REMOVED***
	// Hold the mutex for this entire call to ensure it cannot occur
	// concurrently with other updateState() calls. This causes updates to
	// lastState and calls to cc.UpdateState to happen atomically.
	bw.gsb.mu.Lock()
	defer bw.gsb.mu.Unlock()
	bw.lastState = state

	if !bw.gsb.balancerCurrentOrPending(bw) ***REMOVED***
		return
	***REMOVED***

	if bw == bw.gsb.balancerCurrent ***REMOVED***
		// In the case that the current balancer exits READY, and there is a pending
		// balancer, you can forward the pending balancer's cached State up to
		// ClientConn and swap the pending into the current. This is because there
		// is no reason to gracefully switch from and keep using the old policy as
		// the ClientConn is not connected to any backends.
		if state.ConnectivityState != connectivity.Ready && bw.gsb.balancerPending != nil ***REMOVED***
			bw.gsb.swap()
			return
		***REMOVED***
		// Even if there is a pending balancer waiting to be gracefully switched to,
		// continue to forward current balancer updates to the Client Conn. Ignoring
		// state + picker from the current would cause undefined behavior/cause the
		// system to behave incorrectly from the current LB policies perspective.
		// Also, the current LB is still being used by grpc to choose SubConns per
		// RPC, and thus should use the most updated form of the current balancer.
		bw.gsb.cc.UpdateState(state)
		return
	***REMOVED***
	// This method is now dealing with a state update from the pending balancer.
	// If the current balancer is currently in a state other than READY, the new
	// policy can be swapped into place immediately. This is because there is no
	// reason to gracefully switch from and keep using the old policy as the
	// ClientConn is not connected to any backends.
	if state.ConnectivityState != connectivity.Connecting || bw.gsb.balancerCurrent.lastState.ConnectivityState != connectivity.Ready ***REMOVED***
		bw.gsb.swap()
	***REMOVED***
***REMOVED***

func (bw *balancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) ***REMOVED***
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) ***REMOVED***
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	***REMOVED***
	bw.gsb.mu.Unlock()

	sc, err := bw.gsb.cc.NewSubConn(addrs, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) ***REMOVED*** // balancer was closed during this call
		bw.gsb.cc.RemoveSubConn(sc)
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	***REMOVED***
	bw.subconns[sc] = true
	bw.gsb.mu.Unlock()
	return sc, nil
***REMOVED***

func (bw *balancerWrapper) ResolveNow(opts resolver.ResolveNowOptions) ***REMOVED***
	// Ignore ResolveNow requests from anything other than the most recent
	// balancer, because older balancers were already removed from the config.
	if bw != bw.gsb.latestBalancer() ***REMOVED***
		return
	***REMOVED***
	bw.gsb.cc.ResolveNow(opts)
***REMOVED***

func (bw *balancerWrapper) RemoveSubConn(sc balancer.SubConn) ***REMOVED***
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) ***REMOVED***
		bw.gsb.mu.Unlock()
		return
	***REMOVED***
	bw.gsb.mu.Unlock()
	bw.gsb.cc.RemoveSubConn(sc)
***REMOVED***

func (bw *balancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) ***REMOVED***
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) ***REMOVED***
		bw.gsb.mu.Unlock()
		return
	***REMOVED***
	bw.gsb.mu.Unlock()
	bw.gsb.cc.UpdateAddresses(sc, addrs)
***REMOVED***

func (bw *balancerWrapper) Target() string ***REMOVED***
	return bw.gsb.cc.Target()
***REMOVED***
