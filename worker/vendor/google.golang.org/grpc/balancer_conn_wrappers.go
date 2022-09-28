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
	"strings"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/internal/balancer/gracefulswitch"
	"google.golang.org/grpc/internal/buffer"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/resolver"
)

// ccBalancerWrapper sits between the ClientConn and the Balancer.
//
// ccBalancerWrapper implements methods corresponding to the ones on the
// balancer.Balancer interface. The ClientConn is free to call these methods
// concurrently and the ccBalancerWrapper ensures that calls from the ClientConn
// to the Balancer happen synchronously and in order.
//
// ccBalancerWrapper also implements the balancer.ClientConn interface and is
// passed to the Balancer implementations. It invokes unexported methods on the
// ClientConn to handle these calls from the Balancer.
//
// It uses the gracefulswitch.Balancer internally to ensure that balancer
// switches happen in a graceful manner.
type ccBalancerWrapper struct ***REMOVED***
	cc *ClientConn

	// Since these fields are accessed only from handleXxx() methods which are
	// synchronized by the watcher goroutine, we do not need a mutex to protect
	// these fields.
	balancer        *gracefulswitch.Balancer
	curBalancerName string

	updateCh *buffer.Unbounded // Updates written on this channel are processed by watcher().
	resultCh *buffer.Unbounded // Results of calls to UpdateClientConnState() are pushed here.
	closed   *grpcsync.Event   // Indicates if close has been called.
	done     *grpcsync.Event   // Indicates if close has completed its work.
***REMOVED***

// newCCBalancerWrapper creates a new balancer wrapper. The underlying balancer
// is not created until the switchTo() method is invoked.
func newCCBalancerWrapper(cc *ClientConn, bopts balancer.BuildOptions) *ccBalancerWrapper ***REMOVED***
	ccb := &ccBalancerWrapper***REMOVED***
		cc:       cc,
		updateCh: buffer.NewUnbounded(),
		resultCh: buffer.NewUnbounded(),
		closed:   grpcsync.NewEvent(),
		done:     grpcsync.NewEvent(),
	***REMOVED***
	go ccb.watcher()
	ccb.balancer = gracefulswitch.NewBalancer(ccb, bopts)
	return ccb
***REMOVED***

// The following xxxUpdate structs wrap the arguments received as part of the
// corresponding update. The watcher goroutine uses the 'type' of the update to
// invoke the appropriate handler routine to handle the update.

type ccStateUpdate struct ***REMOVED***
	ccs *balancer.ClientConnState
***REMOVED***

type scStateUpdate struct ***REMOVED***
	sc    balancer.SubConn
	state connectivity.State
	err   error
***REMOVED***

type exitIdleUpdate struct***REMOVED******REMOVED***

type resolverErrorUpdate struct ***REMOVED***
	err error
***REMOVED***

type switchToUpdate struct ***REMOVED***
	name string
***REMOVED***

type subConnUpdate struct ***REMOVED***
	acbw *acBalancerWrapper
***REMOVED***

// watcher is a long-running goroutine which reads updates from a channel and
// invokes corresponding methods on the underlying balancer. It ensures that
// these methods are invoked in a synchronous fashion. It also ensures that
// these methods are invoked in the order in which the updates were received.
func (ccb *ccBalancerWrapper) watcher() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case u := <-ccb.updateCh.Get():
			ccb.updateCh.Load()
			if ccb.closed.HasFired() ***REMOVED***
				break
			***REMOVED***
			switch update := u.(type) ***REMOVED***
			case *ccStateUpdate:
				ccb.handleClientConnStateChange(update.ccs)
			case *scStateUpdate:
				ccb.handleSubConnStateChange(update)
			case *exitIdleUpdate:
				ccb.handleExitIdle()
			case *resolverErrorUpdate:
				ccb.handleResolverError(update.err)
			case *switchToUpdate:
				ccb.handleSwitchTo(update.name)
			case *subConnUpdate:
				ccb.handleRemoveSubConn(update.acbw)
			default:
				logger.Errorf("ccBalancerWrapper.watcher: unknown update %+v, type %T", update, update)
			***REMOVED***
		case <-ccb.closed.Done():
		***REMOVED***

		if ccb.closed.HasFired() ***REMOVED***
			ccb.handleClose()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// updateClientConnState is invoked by grpc to push a ClientConnState update to
// the underlying balancer.
//
// Unlike other methods invoked by grpc to push updates to the underlying
// balancer, this method cannot simply push the update onto the update channel
// and return. It needs to return the error returned by the underlying balancer
// back to grpc which propagates that to the resolver.
func (ccb *ccBalancerWrapper) updateClientConnState(ccs *balancer.ClientConnState) error ***REMOVED***
	ccb.updateCh.Put(&ccStateUpdate***REMOVED***ccs: ccs***REMOVED***)

	var res interface***REMOVED******REMOVED***
	select ***REMOVED***
	case res = <-ccb.resultCh.Get():
		ccb.resultCh.Load()
	case <-ccb.closed.Done():
		// Return early if the balancer wrapper is closed while we are waiting for
		// the underlying balancer to process a ClientConnState update.
		return nil
	***REMOVED***
	// If the returned error is nil, attempting to type assert to error leads to
	// panic. So, this needs to handled separately.
	if res == nil ***REMOVED***
		return nil
	***REMOVED***
	return res.(error)
***REMOVED***

// handleClientConnStateChange handles a ClientConnState update from the update
// channel and invokes the appropriate method on the underlying balancer.
//
// If the addresses specified in the update contain addresses of type "grpclb"
// and the selected LB policy is not "grpclb", these addresses will be filtered
// out and ccs will be modified with the updated address list.
func (ccb *ccBalancerWrapper) handleClientConnStateChange(ccs *balancer.ClientConnState) ***REMOVED***
	if ccb.curBalancerName != grpclbName ***REMOVED***
		// Filter any grpclb addresses since we don't have the grpclb balancer.
		var addrs []resolver.Address
		for _, addr := range ccs.ResolverState.Addresses ***REMOVED***
			if addr.Type == resolver.GRPCLB ***REMOVED***
				continue
			***REMOVED***
			addrs = append(addrs, addr)
		***REMOVED***
		ccs.ResolverState.Addresses = addrs
	***REMOVED***
	ccb.resultCh.Put(ccb.balancer.UpdateClientConnState(*ccs))
***REMOVED***

// updateSubConnState is invoked by grpc to push a subConn state update to the
// underlying balancer.
func (ccb *ccBalancerWrapper) updateSubConnState(sc balancer.SubConn, s connectivity.State, err error) ***REMOVED***
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

// handleSubConnStateChange handles a SubConnState update from the update
// channel and invokes the appropriate method on the underlying balancer.
func (ccb *ccBalancerWrapper) handleSubConnStateChange(update *scStateUpdate) ***REMOVED***
	ccb.balancer.UpdateSubConnState(update.sc, balancer.SubConnState***REMOVED***ConnectivityState: update.state, ConnectionError: update.err***REMOVED***)
***REMOVED***

func (ccb *ccBalancerWrapper) exitIdle() ***REMOVED***
	ccb.updateCh.Put(&exitIdleUpdate***REMOVED******REMOVED***)
***REMOVED***

func (ccb *ccBalancerWrapper) handleExitIdle() ***REMOVED***
	if ccb.cc.GetState() != connectivity.Idle ***REMOVED***
		return
	***REMOVED***
	ccb.balancer.ExitIdle()
***REMOVED***

func (ccb *ccBalancerWrapper) resolverError(err error) ***REMOVED***
	ccb.updateCh.Put(&resolverErrorUpdate***REMOVED***err: err***REMOVED***)
***REMOVED***

func (ccb *ccBalancerWrapper) handleResolverError(err error) ***REMOVED***
	ccb.balancer.ResolverError(err)
***REMOVED***

// switchTo is invoked by grpc to instruct the balancer wrapper to switch to the
// LB policy identified by name.
//
// ClientConn calls newCCBalancerWrapper() at creation time. Upon receipt of the
// first good update from the name resolver, it determines the LB policy to use
// and invokes the switchTo() method. Upon receipt of every subsequent update
// from the name resolver, it invokes this method.
//
// the ccBalancerWrapper keeps track of the current LB policy name, and skips
// the graceful balancer switching process if the name does not change.
func (ccb *ccBalancerWrapper) switchTo(name string) ***REMOVED***
	ccb.updateCh.Put(&switchToUpdate***REMOVED***name: name***REMOVED***)
***REMOVED***

// handleSwitchTo handles a balancer switch update from the update channel. It
// calls the SwitchTo() method on the gracefulswitch.Balancer with a
// balancer.Builder corresponding to name. If no balancer.Builder is registered
// for the given name, it uses the default LB policy which is "pick_first".
func (ccb *ccBalancerWrapper) handleSwitchTo(name string) ***REMOVED***
	// TODO: Other languages use case-insensitive balancer registries. We should
	// switch as well. See: https://github.com/grpc/grpc-go/issues/5288.
	if strings.EqualFold(ccb.curBalancerName, name) ***REMOVED***
		return
	***REMOVED***

	// TODO: Ensure that name is a registered LB policy when we get here.
	// We currently only validate the `loadBalancingConfig` field. We need to do
	// the same for the `loadBalancingPolicy` field and reject the service config
	// if the specified policy is not registered.
	builder := balancer.Get(name)
	if builder == nil ***REMOVED***
		channelz.Warningf(logger, ccb.cc.channelzID, "Channel switches to new LB policy %q, since the specified LB policy %q was not registered", PickFirstBalancerName, name)
		builder = newPickfirstBuilder()
	***REMOVED*** else ***REMOVED***
		channelz.Infof(logger, ccb.cc.channelzID, "Channel switches to new LB policy %q", name)
	***REMOVED***

	if err := ccb.balancer.SwitchTo(builder); err != nil ***REMOVED***
		channelz.Errorf(logger, ccb.cc.channelzID, "Channel failed to build new LB policy %q: %v", name, err)
		return
	***REMOVED***
	ccb.curBalancerName = builder.Name()
***REMOVED***

// handleRemoveSucConn handles a request from the underlying balancer to remove
// a subConn.
//
// See comments in RemoveSubConn() for more details.
func (ccb *ccBalancerWrapper) handleRemoveSubConn(acbw *acBalancerWrapper) ***REMOVED***
	ccb.cc.removeAddrConn(acbw.getAddrConn(), errConnDrain)
***REMOVED***

func (ccb *ccBalancerWrapper) close() ***REMOVED***
	ccb.closed.Fire()
	<-ccb.done.Done()
***REMOVED***

func (ccb *ccBalancerWrapper) handleClose() ***REMOVED***
	ccb.balancer.Close()
	ccb.done.Fire()
***REMOVED***

func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) ***REMOVED***
	if len(addrs) <= 0 ***REMOVED***
		return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
	***REMOVED***
	ac, err := ccb.cc.newAddrConn(addrs, opts)
	if err != nil ***REMOVED***
		channelz.Warningf(logger, ccb.cc.channelzID, "acBalancerWrapper: NewSubConn: failed to newAddrConn: %v", err)
		return nil, err
	***REMOVED***
	acbw := &acBalancerWrapper***REMOVED***ac: ac***REMOVED***
	acbw.ac.mu.Lock()
	ac.acbw = acbw
	acbw.ac.mu.Unlock()
	return acbw, nil
***REMOVED***

func (ccb *ccBalancerWrapper) RemoveSubConn(sc balancer.SubConn) ***REMOVED***
	// Before we switched the ccBalancerWrapper to use gracefulswitch.Balancer, it
	// was required to handle the RemoveSubConn() method asynchronously by pushing
	// the update onto the update channel. This was done to avoid a deadlock as
	// switchBalancer() was holding cc.mu when calling Close() on the old
	// balancer, which would in turn call RemoveSubConn().
	//
	// With the use of gracefulswitch.Balancer in ccBalancerWrapper, handling this
	// asynchronously is probably not required anymore since the switchTo() method
	// handles the balancer switch by pushing the update onto the channel.
	// TODO(easwars): Handle this inline.
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok ***REMOVED***
		return
	***REMOVED***
	ccb.updateCh.Put(&subConnUpdate***REMOVED***acbw: acbw***REMOVED***)
***REMOVED***

func (ccb *ccBalancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) ***REMOVED***
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok ***REMOVED***
		return
	***REMOVED***
	acbw.UpdateAddresses(addrs)
***REMOVED***

func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) ***REMOVED***
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
