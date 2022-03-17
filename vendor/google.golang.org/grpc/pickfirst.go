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
	state connectivity.State
	cc    balancer.ClientConn
	sc    balancer.SubConn
***REMOVED***

func (b *pickfirstBalancer) ResolverError(err error) ***REMOVED***
	switch b.state ***REMOVED***
	case connectivity.TransientFailure, connectivity.Idle, connectivity.Connecting:
		// Set a failing picker if we don't have a good picker.
		b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: connectivity.TransientFailure,
			Picker: &picker***REMOVED***err: fmt.Errorf("name resolver error: %v", err)***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if logger.V(2) ***REMOVED***
		logger.Infof("pickfirstBalancer: ResolverError called with error %v", err)
	***REMOVED***
***REMOVED***

func (b *pickfirstBalancer) UpdateClientConnState(cs balancer.ClientConnState) error ***REMOVED***
	if len(cs.ResolverState.Addresses) == 0 ***REMOVED***
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	***REMOVED***
	if b.sc == nil ***REMOVED***
		var err error
		b.sc, err = b.cc.NewSubConn(cs.ResolverState.Addresses, balancer.NewSubConnOptions***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			if logger.V(2) ***REMOVED***
				logger.Errorf("pickfirstBalancer: failed to NewSubConn: %v", err)
			***REMOVED***
			b.state = connectivity.TransientFailure
			b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: connectivity.TransientFailure,
				Picker: &picker***REMOVED***err: fmt.Errorf("error creating connection: %v", err)***REMOVED***,
			***REMOVED***)
			return balancer.ErrBadResolverState
		***REMOVED***
		b.state = connectivity.Idle
		b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: connectivity.Idle, Picker: &picker***REMOVED***result: balancer.PickResult***REMOVED***SubConn: b.sc***REMOVED******REMOVED******REMOVED***)
		b.sc.Connect()
	***REMOVED*** else ***REMOVED***
		b.cc.UpdateAddresses(b.sc, cs.ResolverState.Addresses)
		b.sc.Connect()
	***REMOVED***
	return nil
***REMOVED***

func (b *pickfirstBalancer) UpdateSubConnState(sc balancer.SubConn, s balancer.SubConnState) ***REMOVED***
	if logger.V(2) ***REMOVED***
		logger.Infof("pickfirstBalancer: UpdateSubConnState: %p, %v", sc, s)
	***REMOVED***
	if b.sc != sc ***REMOVED***
		if logger.V(2) ***REMOVED***
			logger.Infof("pickfirstBalancer: ignored state change because sc is not recognized")
		***REMOVED***
		return
	***REMOVED***
	b.state = s.ConnectivityState
	if s.ConnectivityState == connectivity.Shutdown ***REMOVED***
		b.sc = nil
		return
	***REMOVED***

	switch s.ConnectivityState ***REMOVED***
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: s.ConnectivityState, Picker: &picker***REMOVED***result: balancer.PickResult***REMOVED***SubConn: sc***REMOVED******REMOVED******REMOVED***)
	case connectivity.Connecting:
		b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: s.ConnectivityState, Picker: &picker***REMOVED***err: balancer.ErrNoSubConnAvailable***REMOVED******REMOVED***)
	case connectivity.Idle:
		b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: s.ConnectivityState, Picker: &idlePicker***REMOVED***sc: sc***REMOVED******REMOVED***)
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State***REMOVED***
			ConnectivityState: s.ConnectivityState,
			Picker:            &picker***REMOVED***err: s.ConnectionError***REMOVED***,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (b *pickfirstBalancer) Close() ***REMOVED***
***REMOVED***

func (b *pickfirstBalancer) ExitIdle() ***REMOVED***
	if b.sc != nil && b.state == connectivity.Idle ***REMOVED***
		b.sc.Connect()
	***REMOVED***
***REMOVED***

type picker struct ***REMOVED***
	result balancer.PickResult
	err    error
***REMOVED***

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) ***REMOVED***
	return p.result, p.err
***REMOVED***

// idlePicker is used when the SubConn is IDLE and kicks the SubConn into
// CONNECTING when Pick is called.
type idlePicker struct ***REMOVED***
	sc balancer.SubConn
***REMOVED***

func (i *idlePicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) ***REMOVED***
	i.sc.Connect()
	return balancer.PickResult***REMOVED******REMOVED***, balancer.ErrNoSubConnAvailable
***REMOVED***

func init() ***REMOVED***
	balancer.Register(newPickfirstBuilder())
***REMOVED***
