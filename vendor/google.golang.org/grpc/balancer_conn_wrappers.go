/*
 *
 * Copyright 2017 gRPC authors.
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

package grpc

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/internal/buffer"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/resolver"
)

// scStateUpdate contains the subConn and the new state it changed to.
type scStateUpdate struct ***REMOVED***
	sc    balancer.SubConn
	state connectivity.State
	err   error
***REMOVED***

// exitIdle contains no data and is just a signal sent on the updateCh in
// ccBalancerWrapper to instruct the balancer to exit idle.
type exitIdle struct***REMOVED******REMOVED***

// ccBalancerWrapper is a wrapper on top of cc for balancers.
// It implements balancer.ClientConn interface.
type ccBalancerWrapper struct ***REMOVED***
	cc          *ClientConn
	balancerMu  sync.Mutex // synchronizes calls to the balancer
	balancer    balancer.Balancer
	hasExitIdle bool
	updateCh    *buffer.Unbounded
	closed      *grpcsync.Event
	done        *grpcsync.Event

	mu       sync.Mutex
	subConns map[*acBalancerWrapper]struct***REMOVED******REMOVED***
***REMOVED***

func newCCBalancerWrapper(cc *ClientConn, b balancer.Builder, bopts balancer.BuildOptions) *ccBalancerWrapper ***REMOVED***
	ccb := &ccBalancerWrapper***REMOVED***
		cc:       cc,
		updateCh: buffer.NewUnbounded(),
		closed:   grpcsync.NewEvent(),
		done:     grpcsync.NewEvent(),
		subConns: make(map[*acBalancerWrapper]struct***REMOVED******REMOVED***),
	***REMOVED***
	go ccb.watcher()
	ccb.balancer = b.Build(ccb, bopts)
	_, ccb.hasExitIdle = ccb.balancer.(balancer.ExitIdler)
	return ccb
***REMOVED***

// watcher balancer functions sequentially, so the balancer can be implemented
// lock-free.
func (ccb *ccBalancerWrapper) watcher() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case t := <-ccb.updateCh.Get():
			ccb.updateCh.Load()
			if ccb.closed.HasFired() ***REMOVED***
				break
			***REMOVED***
			switch u := t.(type) ***REMOVED***
			case *scStateUpdate:
				ccb.balancerMu.Lock()
				ccb.balancer.UpdateSubConnState(u.sc, balancer.SubConnState***REMOVED***ConnectivityState: u.state, ConnectionError: u.err***REMOVED***)
				ccb.balancerMu.Unlock()
			case *acBalancerWrapper:
				ccb.mu.Lock()
				if ccb.subConns != nil ***REMOVED***
					delete(ccb.subConns, u)
					ccb.cc.removeAddrConn(u.getAddrConn(), errConnDrain)
				***REMOVED***
				ccb.mu.Unlock()
			case exitIdle:
				if ccb.cc.GetState() == connectivity.Idle ***REMOVED***
					if ei, ok := ccb.balancer.(balancer.ExitIdler); ok ***REMOVED***
						// We already checked that the balancer implements
						// ExitIdle before pushing the event to updateCh, but
						// check conditionally again as defensive programming.
						ccb.balancerMu.Lock()
						ei.ExitIdle()
						ccb.balancerMu.Unlock()
					***REMOVED***
				***REMOVED***
			default:
				logger.Errorf("ccBalancerWrapper.watcher: unknown update %+v, type %T", t, t)
			***REMOVED***
		case <-ccb.closed.Done():
		***REMOVED***

		if ccb.closed.HasFired() ***REMOVED***
			ccb.balancerMu.Lock()
			ccb.balancer.Close()
			ccb.balancerMu.Unlock()
			ccb.mu.Lock()
			scs := ccb.subConns
			ccb.subConns = nil
			ccb.mu.Unlock()
			ccb.UpdateState(balancer.State***REMOVED***ConnectivityState: connectivity.Connecting, Picker: nil***REMOVED***)
			ccb.done.Fire()
			// Fire done before removing the addr conns.  We can safely unblock
			// ccb.close and allow the removeAddrConns to happen
			// asynchronously.
			for acbw := range scs ***REMOVED***
				ccb.cc.removeAddrConn(acbw.getAddrConn(), errConnDrain)
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ccb *ccBalancerWrapper) close() ***REMOVED***
	ccb.closed.Fire()
	<-ccb.done.Done()
***REMOVED***

func (ccb *ccBalancerWrapper) exitIdle() bool ***REMOVED***
	if !ccb.hasExitIdle ***REMOVED***
		return false
	***REMOVED***
	ccb.updateCh.Put(exitIdle***REMOVED******REMOVED***)
	return true
***REMOVED***

func (ccb *ccBalancerWrapper) handleSubConnStateChange(sc balancer.SubConn, s connectivity.State, err error) ***REMOVED***
	// When updating addresses for a SubConn, if the address in use is not in
	// the new addresses, the old ac will be tearDown() and a new ac will be
	// created. tearDown() generates a state change with Shutdown state, we
	// don't want the balancer to receive this state change. So before
	// tearDown() on the old ac, ac.acbw (acWrapper) will be set to nil, and
	// this function will be called with (nil, Shutdown). We don't need to call
	// balancer method in this case.
	if sc == nil ***REMOVED***
		return
	***REMOVED***
	ccb.updateCh.Put(&scStateUpdate***REMOVED***
		sc:    sc,
		state: s,
		err:   err,
	***REMOVED***)
***REMOVED***

func (ccb *ccBalancerWrapper) updateClientConnState(ccs *balancer.ClientConnState) error ***REMOVED***
	ccb.balancerMu.Lock()
	defer ccb.balancerMu.Unlock()
	return ccb.balancer.UpdateClientConnState(*ccs)
***REMOVED***

func (ccb *ccBalancerWrapper) resolverError(err error) ***REMOVED***
	ccb.balancerMu.Lock()
	defer ccb.balancerMu.Unlock()
	ccb.balancer.ResolverError(err)
***REMOVED***

func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) ***REMOVED***
	if len(addrs) <= 0 ***REMOVED***
		return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
	***REMOVED***
	ccb.mu.Lock()
	defer ccb.mu.Unlock()
	if ccb.subConns == nil ***REMOVED***
		return nil, fmt.Errorf("grpc: ClientConn balancer wrapper was closed")
	***REMOVED***
	ac, err := ccb.cc.newAddrConn(addrs, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	acbw := &acBalancerWrapper***REMOVED***ac: ac***REMOVED***
	acbw.ac.mu.Lock()
	ac.acbw = acbw
	acbw.ac.mu.Unlock()
	ccb.subConns[acbw] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return acbw, nil
***REMOVED***

func (ccb *ccBalancerWrapper) RemoveSubConn(sc balancer.SubConn) ***REMOVED***
	// The RemoveSubConn() is handled in the run() goroutine, to avoid deadlock
	// during switchBalancer() if the old balancer calls RemoveSubConn() in its
	// Close().
	ccb.updateCh.Put(sc)
***REMOVED***

func (ccb *ccBalancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) ***REMOVED***
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok ***REMOVED***
		return
	***REMOVED***
	acbw.UpdateAddresses(addrs)
***REMOVED***

func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) ***REMOVED***
	ccb.mu.Lock()
	defer ccb.mu.Unlock()
	if ccb.subConns == nil ***REMOVED***
		return
	***REMOVED***
	// Update picker before updating state.  Even though the ordering here does
	// not matter, it can lead to multiple calls of Pick in the common start-up
	// case where we wait for ready and then perform an RPC.  If the picker is
	// updated later, we could call the "connecting" picker when the state is
	// updated, and then call the "ready" picker after the picker gets updated.
	ccb.cc.blockingpicker.updatePicker(s.Picker)
	ccb.cc.csMgr.updateState(s.ConnectivityState)
***REMOVED***

func (ccb *ccBalancerWrapper) ResolveNow(o resolver.ResolveNowOptions) ***REMOVED***
	ccb.cc.resolveNow(o)
***REMOVED***

func (ccb *ccBalancerWrapper) Target() string ***REMOVED***
	return ccb.cc.target
***REMOVED***

// acBalancerWrapper is a wrapper on top of ac for balancers.
// It implements balancer.SubConn interface.
type acBalancerWrapper struct ***REMOVED***
	mu sync.Mutex
	ac *addrConn
***REMOVED***

func (acbw *acBalancerWrapper) UpdateAddresses(addrs []resolver.Address) ***REMOVED***
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	if len(addrs) <= 0 ***REMOVED***
		acbw.ac.cc.removeAddrConn(acbw.ac, errConnDrain)
		return
	***REMOVED***
	if !acbw.ac.tryUpdateAddrs(addrs) ***REMOVED***
		cc := acbw.ac.cc
		opts := acbw.ac.scopts
		acbw.ac.mu.Lock()
		// Set old ac.acbw to nil so the Shutdown state update will be ignored
		// by balancer.
		//
		// TODO(bar) the state transition could be wrong when tearDown() old ac
		// and creating new ac, fix the transition.
		acbw.ac.acbw = nil
		acbw.ac.mu.Unlock()
		acState := acbw.ac.getState()
		acbw.ac.cc.removeAddrConn(acbw.ac, errConnDrain)

		if acState == connectivity.Shutdown ***REMOVED***
			return
		***REMOVED***

		newAC, err := cc.newAddrConn(addrs, opts)
		if err != nil ***REMOVED***
			channelz.Warningf(logger, acbw.ac.channelzID, "acBalancerWrapper: UpdateAddresses: failed to newAddrConn: %v", err)
			return
		***REMOVED***
		acbw.ac = newAC
		newAC.mu.Lock()
		newAC.acbw = acbw
		newAC.mu.Unlock()
		if acState != connectivity.Idle ***REMOVED***
			go newAC.connect()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (acbw *acBalancerWrapper) Connect() ***REMOVED***
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	go acbw.ac.connect()
***REMOVED***

func (acbw *acBalancerWrapper) getAddrConn() *addrConn ***REMOVED***
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	return acbw.ac
***REMOVED***
