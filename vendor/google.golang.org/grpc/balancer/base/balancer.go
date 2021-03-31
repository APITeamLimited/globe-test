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

package base

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var logger = grpclog.Component("balancer")

type baseBuilder struct ***REMOVED***
	name          string
	pickerBuilder PickerBuilder
	config        Config
***REMOVED***

func (bb *baseBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer ***REMOVED***
	bal := &baseBalancer***REMOVED***
		cc:            cc,
		pickerBuilder: bb.pickerBuilder,

		subConns: make(map[resolver.Address]balancer.SubConn),
		scStates: make(map[balancer.SubConn]connectivity.State),
		csEvltr:  &balancer.ConnectivityStateEvaluator***REMOVED******REMOVED***,
		config:   bb.config,
	***REMOVED***
	// Initialize picker to a picker that always returns
	// ErrNoSubConnAvailable, because when state of a SubConn changes, we
	// may call UpdateState with this picker.
	bal.picker = NewErrPicker(balancer.ErrNoSubConnAvailable)
	return bal
***REMOVED***

func (bb *baseBuilder) Name() string ***REMOVED***
	return bb.name
***REMOVED***

type baseBalancer struct ***REMOVED***
	cc            balancer.ClientConn
	pickerBuilder PickerBuilder

	csEvltr *balancer.ConnectivityStateEvaluator
	state   connectivity.State

	subConns map[resolver.Address]balancer.SubConn // `attributes` is stripped from the keys of this map (the addresses)
	scStates map[balancer.SubConn]connectivity.State
	picker   balancer.Picker
	config   Config

	resolverErr error // the last error reported by the resolver; cleared on successful resolution
	connErr     error // the last connection error; cleared upon leaving TransientFailure
***REMOVED***

func (b *baseBalancer) ResolverError(err error) ***REMOVED***
	b.resolverErr = err
	if len(b.subConns) == 0 ***REMOVED***
		b.state = connectivity.TransientFailure
	***REMOVED***

	if b.state != connectivity.TransientFailure ***REMOVED***
		// The picker will not change since the balancer does not currently
		// report an error.
		return
	***REMOVED***
	b.regeneratePicker()
	b.cc.UpdateState(balancer.State***REMOVED***
		ConnectivityState: b.state,
		Picker:            b.picker,
	***REMOVED***)
***REMOVED***

func (b *baseBalancer) UpdateClientConnState(s balancer.ClientConnState) error ***REMOVED***
	// TODO: handle s.ResolverState.ServiceConfig?
	if logger.V(2) ***REMOVED***
		logger.Info("base.baseBalancer: got new ClientConn state: ", s)
	***REMOVED***
	// Successful resolution; clear resolver error and ensure we return nil.
	b.resolverErr = nil
	// addrsSet is the set converted from addrs, it's used for quick lookup of an address.
	addrsSet := make(map[resolver.Address]struct***REMOVED******REMOVED***)
	for _, a := range s.ResolverState.Addresses ***REMOVED***
		// Strip attributes from addresses before using them as map keys. So
		// that when two addresses only differ in attributes pointers (but with
		// the same attribute content), they are considered the same address.
		//
		// Note that this doesn't handle the case where the attribute content is
		// different. So if users want to set different attributes to create
		// duplicate connections to the same backend, it doesn't work. This is
		// fine for now, because duplicate is done by setting Metadata today.
		//
		// TODO: read attributes to handle duplicate connections.
		aNoAttrs := a
		aNoAttrs.Attributes = nil
		addrsSet[aNoAttrs] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		if sc, ok := b.subConns[aNoAttrs]; !ok ***REMOVED***
			// a is a new address (not existing in b.subConns).
			//
			// When creating SubConn, the original address with attributes is
			// passed through. So that connection configurations in attributes
			// (like creds) will be used.
			sc, err := b.cc.NewSubConn([]resolver.Address***REMOVED***a***REMOVED***, balancer.NewSubConnOptions***REMOVED***HealthCheckEnabled: b.config.HealthCheck***REMOVED***)
			if err != nil ***REMOVED***
				logger.Warningf("base.baseBalancer: failed to create new SubConn: %v", err)
				continue
			***REMOVED***
			b.subConns[aNoAttrs] = sc
			b.scStates[sc] = connectivity.Idle
			sc.Connect()
		***REMOVED*** else ***REMOVED***
			// Always update the subconn's address in case the attributes
			// changed.
			//
			// The SubConn does a reflect.DeepEqual of the new and old
			// addresses. So this is a noop if the current address is the same
			// as the old one (including attributes).
			sc.UpdateAddresses([]resolver.Address***REMOVED***a***REMOVED***)
		***REMOVED***
	***REMOVED***
	for a, sc := range b.subConns ***REMOVED***
		// a was removed by resolver.
		if _, ok := addrsSet[a]; !ok ***REMOVED***
			b.cc.RemoveSubConn(sc)
			delete(b.subConns, a)
			// Keep the state of this sc in b.scStates until sc's state becomes Shutdown.
			// The entry will be deleted in UpdateSubConnState.
		***REMOVED***
	***REMOVED***
	// If resolver state contains no addresses, return an error so ClientConn
	// will trigger re-resolve. Also records this as an resolver error, so when
	// the overall state turns transient failure, the error message will have
	// the zero address information.
	if len(s.ResolverState.Addresses) == 0 ***REMOVED***
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	***REMOVED***
	return nil
***REMOVED***

// mergeErrors builds an error from the last connection error and the last
// resolver error.  Must only be called if b.state is TransientFailure.
func (b *baseBalancer) mergeErrors() error ***REMOVED***
	// connErr must always be non-nil unless there are no SubConns, in which
	// case resolverErr must be non-nil.
	if b.connErr == nil ***REMOVED***
		return fmt.Errorf("last resolver error: %v", b.resolverErr)
	***REMOVED***
	if b.resolverErr == nil ***REMOVED***
		return fmt.Errorf("last connection error: %v", b.connErr)
	***REMOVED***
	return fmt.Errorf("last connection error: %v; last resolver error: %v", b.connErr, b.resolverErr)
***REMOVED***

// regeneratePicker takes a snapshot of the balancer, and generates a picker
// from it. The picker is
//  - errPicker if the balancer is in TransientFailure,
//  - built by the pickerBuilder with all READY SubConns otherwise.
func (b *baseBalancer) regeneratePicker() ***REMOVED***
	if b.state == connectivity.TransientFailure ***REMOVED***
		b.picker = NewErrPicker(b.mergeErrors())
		return
	***REMOVED***
	readySCs := make(map[balancer.SubConn]SubConnInfo)

	// Filter out all ready SCs from full subConn map.
	for addr, sc := range b.subConns ***REMOVED***
		if st, ok := b.scStates[sc]; ok && st == connectivity.Ready ***REMOVED***
			readySCs[sc] = SubConnInfo***REMOVED***Address: addr***REMOVED***
		***REMOVED***
	***REMOVED***
	b.picker = b.pickerBuilder.Build(PickerBuildInfo***REMOVED***ReadySCs: readySCs***REMOVED***)
***REMOVED***

func (b *baseBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) ***REMOVED***
	s := state.ConnectivityState
	if logger.V(2) ***REMOVED***
		logger.Infof("base.baseBalancer: handle SubConn state change: %p, %v", sc, s)
	***REMOVED***
	oldS, ok := b.scStates[sc]
	if !ok ***REMOVED***
		if logger.V(2) ***REMOVED***
			logger.Infof("base.baseBalancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		***REMOVED***
		return
	***REMOVED***
	if oldS == connectivity.TransientFailure && s == connectivity.Connecting ***REMOVED***
		// Once a subconn enters TRANSIENT_FAILURE, ignore subsequent
		// CONNECTING transitions to prevent the aggregated state from being
		// always CONNECTING when many backends exist but are all down.
		return
	***REMOVED***
	b.scStates[sc] = s
	switch s ***REMOVED***
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		// When an address was removed by resolver, b called RemoveSubConn but
		// kept the sc's state in scStates. Remove state for this sc here.
		delete(b.scStates, sc)
	case connectivity.TransientFailure:
		// Save error to be reported via picker.
		b.connErr = state.ConnectionError
	***REMOVED***

	b.state = b.csEvltr.RecordTransition(oldS, s)

	// Regenerate picker when one of the following happens:
	//  - this sc entered or left ready
	//  - the aggregated state of balancer is TransientFailure
	//    (may need to update error message)
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		b.state == connectivity.TransientFailure ***REMOVED***
		b.regeneratePicker()
	***REMOVED***

	b.cc.UpdateState(balancer.State***REMOVED***ConnectivityState: b.state, Picker: b.picker***REMOVED***)
***REMOVED***

// Close is a nop because base balancer doesn't have internal state to clean up,
// and it doesn't need to call RemoveSubConn for the SubConns.
func (b *baseBalancer) Close() ***REMOVED***
***REMOVED***

// NewErrPicker returns a Picker that always returns err on Pick().
func NewErrPicker(err error) balancer.Picker ***REMOVED***
	return &errPicker***REMOVED***err: err***REMOVED***
***REMOVED***

// NewErrPickerV2 is temporarily defined for backward compatibility reasons.
//
// Deprecated: use NewErrPicker instead.
var NewErrPickerV2 = NewErrPicker

type errPicker struct ***REMOVED***
	err error // Pick() always returns this err.
***REMOVED***

func (p *errPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) ***REMOVED***
	return balancer.PickResult***REMOVED******REMOVED***, p.err
***REMOVED***
