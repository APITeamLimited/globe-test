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
	"errors"
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
)

// PickFirstBalancerName is the name of the pick_first balancer.
const PickFirstBalancerName = "pick_first"

func newPickfirstBuilder() balancer.Builder ***REMOVED***
	return &pickfirstBuilder***REMOVED******REMOVED***
***REMOVED***

type pickfirstBuilder struct***REMOVED******REMOVED***

func (*pickfirstBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer ***REMOVED***
	return &pickfirstBalancer***REMOVED***cc: cc***REMOVED***
***REMOVED***

func (*pickfirstBuilder) Name() string ***REMOVED***
	return PickFirstBalancerName
***REMOVED***

type pickfirstBalancer struct ***REMOVED***
	state   connectivity.State
	cc      balancer.ClientConn
	subConn balancer.SubConn
***REMOVED***

func (b *pickfirstBalancer) ResolverError(err error) ***REMOVED***
	if logger.V(2) ***REMOVED***
		logger.Infof("pickfirstBalancer: ResolverError called with error %v", err)
	***REMOVED***
	if b.subConn == nil ***REMOVED***
		b.state = connectivity.TransientFailure
	***REMOVED***

	if b.state != connectivity.TransientFailure ***REMOVED***
		// The picker will not change since the balancer does not currently
		// report an error.
		return
	***REMOVED***
	b.cc.UpdateState(balancer.State***REMOVED***
		ConnectivityState: connectivity.TransientFailure,
		Picker:            &picker***REMOVED***err: fmt.Errorf("name resolver error: %v", err)***REMOVED***,
	***REMOVED***)
***REMOVED***

func (b *pickfirstBalancer) UpdateClientConnState(state balancer.ClientConnState) error ***REMOVED***
	if len(state.ResolverState.Addresses) == 0 ***REMOVED***
		// The resolver reported an empty address list. Treat it like an error by
		// calling b.ResolverError.
		if b.subConn != nil ***REMOVED***
			// Remove the old subConn. All addresses were removed, so it is no longer
			// valid.
			b.cc.RemoveSubConn(b.subConn)
			b.subConn = nil
		***REMOVED***
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	***REMOVED***

	if b.subConn != nil ***REMOVED***
		b.cc.UpdateAddresses(b.subConn, state.ResolverState.Addresses)
		return nil
	***REMOVED***

	subConn, err := b.cc.NewSubConn(state.ResolverState.Addresses, balancer.NewSubConnOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		if logger.V(2) ***REMOVED***
			logger.Errorf("pickfirstBalancer: failed to NewSubConn: %v", err)
		***REMOVED***
		b.state = connectivity.TransientFailure
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: connectivity.TransientFailure,
			Picker:            &picker***REMOVED***err: fmt.Errorf("error creating connection: %v", err)***REMOVED***,
		***REMOVED***)
		return balancer.ErrBadResolverState
	***REMOVED***
	b.subConn = subConn
	b.state = connectivity.Idle
	b.cc.UpdateState(balancer.State***REMOVED***
		ConnectivityState: connectivity.Idle,
		Picker:            &picker***REMOVED***result: balancer.PickResult***REMOVED***SubConn: b.subConn***REMOVED******REMOVED***,
	***REMOVED***)
	b.subConn.Connect()
	return nil
***REMOVED***

func (b *pickfirstBalancer) UpdateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) ***REMOVED***
	if logger.V(2) ***REMOVED***
		logger.Infof("pickfirstBalancer: UpdateSubConnState: %p, %v", subConn, state)
	***REMOVED***
	if b.subConn != subConn ***REMOVED***
		if logger.V(2) ***REMOVED***
			logger.Infof("pickfirstBalancer: ignored state change because subConn is not recognized")
		***REMOVED***
		return
	***REMOVED***
	b.state = state.ConnectivityState
	if state.ConnectivityState == connectivity.Shutdown ***REMOVED***
		b.subConn = nil
		return
	***REMOVED***

	switch state.ConnectivityState ***REMOVED***
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker***REMOVED***result: balancer.PickResult***REMOVED***SubConn: subConn***REMOVED******REMOVED***,
		***REMOVED***)
	case connectivity.Connecting:
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker***REMOVED***err: balancer.ErrNoSubConnAvailable***REMOVED***,
		***REMOVED***)
	case connectivity.Idle:
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: state.ConnectivityState,
			Picker:            &idlePicker***REMOVED***subConn: subConn***REMOVED***,
		***REMOVED***)
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker***REMOVED***err: state.ConnectionError***REMOVED***,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (b *pickfirstBalancer) Close() ***REMOVED***
***REMOVED***

func (b *pickfirstBalancer) ExitIdle() ***REMOVED***
	if b.subConn != nil && b.state == connectivity.Idle ***REMOVED***
		b.subConn.Connect()
	***REMOVED***
***REMOVED***

type picker struct ***REMOVED***
	result balancer.PickResult
	err    error
***REMOVED***

func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) ***REMOVED***
	return p.result, p.err
***REMOVED***

// idlePicker is used when the SubConn is IDLE and kicks the SubConn into
// CONNECTING when Pick is called.
type idlePicker struct ***REMOVED***
	subConn balancer.SubConn
***REMOVED***

func (i *idlePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) ***REMOVED***
	i.subConn.Connect()
	return balancer.PickResult***REMOVED******REMOVED***, balancer.ErrNoSubConnAvailable
***REMOVED***

func init() ***REMOVED***
	balancer.Register(newPickfirstBuilder())
***REMOVED***
