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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

// ccResolverWrapper is a wrapper on top of cc for resolvers.
// It implements resolver.ClientConn interface.
type ccResolverWrapper struct ***REMOVED***
	cc         *ClientConn
	resolverMu sync.Mutex
	resolver   resolver.Resolver
	done       *grpcsync.Event
	curState   resolver.State

	incomingMu sync.Mutex // Synchronizes all the incoming calls.
***REMOVED***

// newCCResolverWrapper uses the resolver.Builder to build a Resolver and
// returns a ccResolverWrapper object which wraps the newly built resolver.
func newCCResolverWrapper(cc *ClientConn, rb resolver.Builder) (*ccResolverWrapper, error) ***REMOVED***
	ccr := &ccResolverWrapper***REMOVED***
		cc:   cc,
		done: grpcsync.NewEvent(),
	***REMOVED***

	var credsClone credentials.TransportCredentials
	if creds := cc.dopts.copts.TransportCredentials; creds != nil ***REMOVED***
		credsClone = creds.Clone()
	***REMOVED***
	rbo := resolver.BuildOptions***REMOVED***
		DisableServiceConfig: cc.dopts.disableServiceConfig,
		DialCreds:            credsClone,
		CredsBundle:          cc.dopts.copts.CredsBundle,
		Dialer:               cc.dopts.copts.Dialer,
	***REMOVED***

	var err error
	// We need to hold the lock here while we assign to the ccr.resolver field
	// to guard against a data race caused by the following code path,
	// rb.Build-->ccr.ReportError-->ccr.poll-->ccr.resolveNow, would end up
	// accessing ccr.resolver which is being assigned here.
	ccr.resolverMu.Lock()
	defer ccr.resolverMu.Unlock()
	ccr.resolver, err = rb.Build(cc.parsedTarget, ccr, rbo)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return ccr, nil
***REMOVED***

func (ccr *ccResolverWrapper) resolveNow(o resolver.ResolveNowOptions) ***REMOVED***
	ccr.resolverMu.Lock()
	if !ccr.done.HasFired() ***REMOVED***
		ccr.resolver.ResolveNow(o)
	***REMOVED***
	ccr.resolverMu.Unlock()
***REMOVED***

func (ccr *ccResolverWrapper) close() ***REMOVED***
	ccr.resolverMu.Lock()
	ccr.resolver.Close()
	ccr.done.Fire()
	ccr.resolverMu.Unlock()
***REMOVED***

func (ccr *ccResolverWrapper) UpdateState(s resolver.State) error ***REMOVED***
	ccr.incomingMu.Lock()
	defer ccr.incomingMu.Unlock()
	if ccr.done.HasFired() ***REMOVED***
		return nil
	***REMOVED***
	channelz.Infof(logger, ccr.cc.channelzID, "ccResolverWrapper: sending update to cc: %v", s)
	if channelz.IsOn() ***REMOVED***
		ccr.addChannelzTraceEvent(s)
	***REMOVED***
	ccr.curState = s
	if err := ccr.cc.updateResolverState(ccr.curState, nil); err == balancer.ErrBadResolverState ***REMOVED***
		return balancer.ErrBadResolverState
	***REMOVED***
	return nil
***REMOVED***

func (ccr *ccResolverWrapper) ReportError(err error) ***REMOVED***
	ccr.incomingMu.Lock()
	defer ccr.incomingMu.Unlock()
	if ccr.done.HasFired() ***REMOVED***
		return
	***REMOVED***
	channelz.Warningf(logger, ccr.cc.channelzID, "ccResolverWrapper: reporting error to cc: %v", err)
	ccr.cc.updateResolverState(resolver.State***REMOVED******REMOVED***, err)
***REMOVED***

// NewAddress is called by the resolver implementation to send addresses to gRPC.
func (ccr *ccResolverWrapper) NewAddress(addrs []resolver.Address) ***REMOVED***
	ccr.incomingMu.Lock()
	defer ccr.incomingMu.Unlock()
	if ccr.done.HasFired() ***REMOVED***
		return
	***REMOVED***
	channelz.Infof(logger, ccr.cc.channelzID, "ccResolverWrapper: sending new addresses to cc: %v", addrs)
	if channelz.IsOn() ***REMOVED***
		ccr.addChannelzTraceEvent(resolver.State***REMOVED***Addresses: addrs, ServiceConfig: ccr.curState.ServiceConfig***REMOVED***)
	***REMOVED***
	ccr.curState.Addresses = addrs
	ccr.cc.updateResolverState(ccr.curState, nil)
***REMOVED***

// NewServiceConfig is called by the resolver implementation to send service
// configs to gRPC.
func (ccr *ccResolverWrapper) NewServiceConfig(sc string) ***REMOVED***
	ccr.incomingMu.Lock()
	defer ccr.incomingMu.Unlock()
	if ccr.done.HasFired() ***REMOVED***
		return
	***REMOVED***
	channelz.Infof(logger, ccr.cc.channelzID, "ccResolverWrapper: got new service config: %v", sc)
	if ccr.cc.dopts.disableServiceConfig ***REMOVED***
		channelz.Info(logger, ccr.cc.channelzID, "Service config lookups disabled; ignoring config")
		return
	***REMOVED***
	scpr := parseServiceConfig(sc)
	if scpr.Err != nil ***REMOVED***
		channelz.Warningf(logger, ccr.cc.channelzID, "ccResolverWrapper: error parsing service config: %v", scpr.Err)
		return
	***REMOVED***
	if channelz.IsOn() ***REMOVED***
		ccr.addChannelzTraceEvent(resolver.State***REMOVED***Addresses: ccr.curState.Addresses, ServiceConfig: scpr***REMOVED***)
	***REMOVED***
	ccr.curState.ServiceConfig = scpr
	ccr.cc.updateResolverState(ccr.curState, nil)
***REMOVED***

func (ccr *ccResolverWrapper) ParseServiceConfig(scJSON string) *serviceconfig.ParseResult ***REMOVED***
	return parseServiceConfig(scJSON)
***REMOVED***

func (ccr *ccResolverWrapper) addChannelzTraceEvent(s resolver.State) ***REMOVED***
	var updates []string
	var oldSC, newSC *ServiceConfig
	var oldOK, newOK bool
	if ccr.curState.ServiceConfig != nil ***REMOVED***
		oldSC, oldOK = ccr.curState.ServiceConfig.Config.(*ServiceConfig)
	***REMOVED***
	if s.ServiceConfig != nil ***REMOVED***
		newSC, newOK = s.ServiceConfig.Config.(*ServiceConfig)
	***REMOVED***
	if oldOK != newOK || (oldOK && newOK && oldSC.rawJSONString != newSC.rawJSONString) ***REMOVED***
		updates = append(updates, "service config updated")
	***REMOVED***
	if len(ccr.curState.Addresses) > 0 && len(s.Addresses) == 0 ***REMOVED***
		updates = append(updates, "resolver returned an empty address list")
	***REMOVED*** else if len(ccr.curState.Addresses) == 0 && len(s.Addresses) > 0 ***REMOVED***
		updates = append(updates, "resolver returned new addresses")
	***REMOVED***
	channelz.AddTraceEvent(logger, ccr.cc.channelzID, 0, &channelz.TraceEventDesc***REMOVED***
		Desc:     fmt.Sprintf("Resolver state updated: %+v (%v)", s, strings.Join(updates, "; ")),
		Severity: channelz.CtInfo,
	***REMOVED***)
***REMOVED***
