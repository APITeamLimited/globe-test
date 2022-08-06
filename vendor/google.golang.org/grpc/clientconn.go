/*
 *
 * Copyright 2014 gRPC authors.
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
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	iresolver "google.golang.org/grpc/internal/resolver"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/balancer/roundrobin"           // To register roundrobin.
	_ "google.golang.org/grpc/internal/resolver/dns"         // To register dns resolver.
	_ "google.golang.org/grpc/internal/resolver/passthrough" // To register passthrough resolver.
	_ "google.golang.org/grpc/internal/resolver/unix"        // To register unix resolver.
)

const (
	// minimum time to give a connection to complete
	minConnectTimeout = 20 * time.Second
	// must match grpclbName in grpclb/grpclb.go
	grpclbName = "grpclb"
)

var (
	// ErrClientConnClosing indicates that the operation is illegal because
	// the ClientConn is closing.
	//
	// Deprecated: this error should not be relied upon by users; use the status
	// code of Canceled instead.
	ErrClientConnClosing = status.Error(codes.Canceled, "grpc: the client connection is closing")
	// errConnDrain indicates that the connection starts to be drained and does not accept any new RPCs.
	errConnDrain = errors.New("grpc: the connection is drained")
	// errConnClosing indicates that the connection is closing.
	errConnClosing = errors.New("grpc: the connection is closing")
	// invalidDefaultServiceConfigErrPrefix is used to prefix the json parsing error for the default
	// service config.
	invalidDefaultServiceConfigErrPrefix = "grpc: the provided default service config is invalid"
)

// The following errors are returned from Dial and DialContext
var (
	// errNoTransportSecurity indicates that there is no transport security
	// being set for ClientConn. Users should either set one or explicitly
	// call WithInsecure DialOption to disable security.
	errNoTransportSecurity = errors.New("grpc: no transport security set (use grpc.WithTransportCredentials(insecure.NewCredentials()) explicitly or set credentials)")
	// errTransportCredsAndBundle indicates that creds bundle is used together
	// with other individual Transport Credentials.
	errTransportCredsAndBundle = errors.New("grpc: credentials.Bundle may not be used with individual TransportCredentials")
	// errNoTransportCredsInBundle indicated that the configured creds bundle
	// returned a transport credentials which was nil.
	errNoTransportCredsInBundle = errors.New("grpc: credentials.Bundle must return non-nil transport credentials")
	// errTransportCredentialsMissing indicates that users want to transmit
	// security information (e.g., OAuth2 token) which requires secure
	// connection on an insecure connection.
	errTransportCredentialsMissing = errors.New("grpc: the credentials require transport level security (use grpc.WithTransportCredentials() to set)")
)

const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultClientMaxSendMessageSize    = math.MaxInt32
	// http2IOBufSize specifies the buffer size for sending frames.
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

// Dial creates a client connection to the given target.
func Dial(target string, opts ...DialOption) (*ClientConn, error) ***REMOVED***
	return DialContext(context.Background(), target, opts...)
***REMOVED***

type defaultConfigSelector struct ***REMOVED***
	sc *ServiceConfig
***REMOVED***

func (dcs *defaultConfigSelector) SelectConfig(rpcInfo iresolver.RPCInfo) (*iresolver.RPCConfig, error) ***REMOVED***
	return &iresolver.RPCConfig***REMOVED***
		Context:      rpcInfo.Context,
		MethodConfig: getMethodConfig(dcs.sc, rpcInfo.Method),
	***REMOVED***, nil
***REMOVED***

// DialContext creates a client connection to the given target. By default, it's
// a non-blocking dial (the function won't wait for connections to be
// established, and connecting happens in the background). To make it a blocking
// dial, use WithBlock() dial option.
//
// In the non-blocking case, the ctx does not act against the connection. It
// only controls the setup steps.
//
// In the blocking case, ctx can be used to cancel or expire the pending
// connection. Once this function returns, the cancellation and expiration of
// ctx will be noop. Users should call ClientConn.Close to terminate all the
// pending operations after this function returns.
//
// The target name syntax is defined in
// https://github.com/grpc/grpc/blob/master/doc/naming.md.
// e.g. to use dns resolver, a "dns:///" prefix should be applied to the target.
func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) ***REMOVED***
	cc := &ClientConn***REMOVED***
		target:            target,
		csMgr:             &connectivityStateManager***REMOVED******REMOVED***,
		conns:             make(map[*addrConn]struct***REMOVED******REMOVED***),
		dopts:             defaultDialOptions(),
		blockingpicker:    newPickerWrapper(),
		czData:            new(channelzData),
		firstResolveEvent: grpcsync.NewEvent(),
	***REMOVED***
	cc.retryThrottler.Store((*retryThrottler)(nil))
	cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector***REMOVED***nil***REMOVED***)
	cc.ctx, cc.cancel = context.WithCancel(context.Background())

	for _, opt := range extraDialOptions ***REMOVED***
		opt.apply(&cc.dopts)
	***REMOVED***

	for _, opt := range opts ***REMOVED***
		opt.apply(&cc.dopts)
	***REMOVED***

	chainUnaryClientInterceptors(cc)
	chainStreamClientInterceptors(cc)

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			cc.Close()
		***REMOVED***
	***REMOVED***()

	pid := cc.dopts.channelzParentID
	cc.channelzID = channelz.RegisterChannel(&channelzChannel***REMOVED***cc***REMOVED***, pid, target)
	ted := &channelz.TraceEventDesc***REMOVED***
		Desc:     "Channel created",
		Severity: channelz.CtInfo,
	***REMOVED***
	if cc.dopts.channelzParentID != nil ***REMOVED***
		ted.Parent = &channelz.TraceEventDesc***REMOVED***
			Desc:     fmt.Sprintf("Nested Channel(id:%d) created", cc.channelzID.Int()),
			Severity: channelz.CtInfo,
		***REMOVED***
	***REMOVED***
	channelz.AddTraceEvent(logger, cc.channelzID, 1, ted)
	cc.csMgr.channelzID = cc.channelzID

	if cc.dopts.copts.TransportCredentials == nil && cc.dopts.copts.CredsBundle == nil ***REMOVED***
		return nil, errNoTransportSecurity
	***REMOVED***
	if cc.dopts.copts.TransportCredentials != nil && cc.dopts.copts.CredsBundle != nil ***REMOVED***
		return nil, errTransportCredsAndBundle
	***REMOVED***
	if cc.dopts.copts.CredsBundle != nil && cc.dopts.copts.CredsBundle.TransportCredentials() == nil ***REMOVED***
		return nil, errNoTransportCredsInBundle
	***REMOVED***
	transportCreds := cc.dopts.copts.TransportCredentials
	if transportCreds == nil ***REMOVED***
		transportCreds = cc.dopts.copts.CredsBundle.TransportCredentials()
	***REMOVED***
	if transportCreds.Info().SecurityProtocol == "insecure" ***REMOVED***
		for _, cd := range cc.dopts.copts.PerRPCCredentials ***REMOVED***
			if cd.RequireTransportSecurity() ***REMOVED***
				return nil, errTransportCredentialsMissing
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if cc.dopts.defaultServiceConfigRawJSON != nil ***REMOVED***
		scpr := parseServiceConfig(*cc.dopts.defaultServiceConfigRawJSON)
		if scpr.Err != nil ***REMOVED***
			return nil, fmt.Errorf("%s: %v", invalidDefaultServiceConfigErrPrefix, scpr.Err)
		***REMOVED***
		cc.dopts.defaultServiceConfig, _ = scpr.Config.(*ServiceConfig)
	***REMOVED***
	cc.mkp = cc.dopts.copts.KeepaliveParams

	if cc.dopts.copts.UserAgent != "" ***REMOVED***
		cc.dopts.copts.UserAgent += " " + grpcUA
	***REMOVED*** else ***REMOVED***
		cc.dopts.copts.UserAgent = grpcUA
	***REMOVED***

	if cc.dopts.timeout > 0 ***REMOVED***
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
		defer cancel()
	***REMOVED***
	defer func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			switch ***REMOVED***
			case ctx.Err() == err:
				conn = nil
			case err == nil || !cc.dopts.returnLastError:
				conn, err = nil, ctx.Err()
			default:
				conn, err = nil, fmt.Errorf("%v: %v", ctx.Err(), err)
			***REMOVED***
		default:
		***REMOVED***
	***REMOVED***()

	scSet := false
	if cc.dopts.scChan != nil ***REMOVED***
		// Try to get an initial service config.
		select ***REMOVED***
		case sc, ok := <-cc.dopts.scChan:
			if ok ***REMOVED***
				cc.sc = &sc
				cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector***REMOVED***&sc***REMOVED***)
				scSet = true
			***REMOVED***
		default:
		***REMOVED***
	***REMOVED***
	if cc.dopts.bs == nil ***REMOVED***
		cc.dopts.bs = backoff.DefaultExponential
	***REMOVED***

	// Determine the resolver to use.
	resolverBuilder, err := cc.parseTargetAndFindResolver()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cc.authority, err = determineAuthority(cc.parsedTarget.Endpoint, cc.target, cc.dopts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	channelz.Infof(logger, cc.channelzID, "Channel authority set to %q", cc.authority)

	if cc.dopts.scChan != nil && !scSet ***REMOVED***
		// Blocking wait for the initial service config.
		select ***REMOVED***
		case sc, ok := <-cc.dopts.scChan:
			if ok ***REMOVED***
				cc.sc = &sc
				cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector***REMOVED***&sc***REMOVED***)
			***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		***REMOVED***
	***REMOVED***
	if cc.dopts.scChan != nil ***REMOVED***
		go cc.scWatcher()
	***REMOVED***

	var credsClone credentials.TransportCredentials
	if creds := cc.dopts.copts.TransportCredentials; creds != nil ***REMOVED***
		credsClone = creds.Clone()
	***REMOVED***
	cc.balancerWrapper = newCCBalancerWrapper(cc, balancer.BuildOptions***REMOVED***
		DialCreds:        credsClone,
		CredsBundle:      cc.dopts.copts.CredsBundle,
		Dialer:           cc.dopts.copts.Dialer,
		Authority:        cc.authority,
		CustomUserAgent:  cc.dopts.copts.UserAgent,
		ChannelzParentID: cc.channelzID,
		Target:           cc.parsedTarget,
	***REMOVED***)

	// Build the resolver.
	rWrapper, err := newCCResolverWrapper(cc, resolverBuilder)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to build resolver: %v", err)
	***REMOVED***
	cc.mu.Lock()
	cc.resolverWrapper = rWrapper
	cc.mu.Unlock()

	// A blocking dial blocks until the clientConn is ready.
	if cc.dopts.block ***REMOVED***
		for ***REMOVED***
			cc.Connect()
			s := cc.GetState()
			if s == connectivity.Ready ***REMOVED***
				break
			***REMOVED*** else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure ***REMOVED***
				if err = cc.connectionError(); err != nil ***REMOVED***
					terr, ok := err.(interface ***REMOVED***
						Temporary() bool
					***REMOVED***)
					if ok && !terr.Temporary() ***REMOVED***
						return nil, err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !cc.WaitForStateChange(ctx, s) ***REMOVED***
				// ctx got timeout or canceled.
				if err = cc.connectionError(); err != nil && cc.dopts.returnLastError ***REMOVED***
					return nil, err
				***REMOVED***
				return nil, ctx.Err()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return cc, nil
***REMOVED***

// chainUnaryClientInterceptors chains all unary client interceptors into one.
func chainUnaryClientInterceptors(cc *ClientConn) ***REMOVED***
	interceptors := cc.dopts.chainUnaryInts
	// Prepend dopts.unaryInt to the chaining interceptors if it exists, since unaryInt will
	// be executed before any other chained interceptors.
	if cc.dopts.unaryInt != nil ***REMOVED***
		interceptors = append([]UnaryClientInterceptor***REMOVED***cc.dopts.unaryInt***REMOVED***, interceptors...)
	***REMOVED***
	var chainedInt UnaryClientInterceptor
	if len(interceptors) == 0 ***REMOVED***
		chainedInt = nil
	***REMOVED*** else if len(interceptors) == 1 ***REMOVED***
		chainedInt = interceptors[0]
	***REMOVED*** else ***REMOVED***
		chainedInt = func(ctx context.Context, method string, req, reply interface***REMOVED******REMOVED***, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error ***REMOVED***
			return interceptors[0](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, 0, invoker), opts...)
		***REMOVED***
	***REMOVED***
	cc.dopts.unaryInt = chainedInt
***REMOVED***

// getChainUnaryInvoker recursively generate the chained unary invoker.
func getChainUnaryInvoker(interceptors []UnaryClientInterceptor, curr int, finalInvoker UnaryInvoker) UnaryInvoker ***REMOVED***
	if curr == len(interceptors)-1 ***REMOVED***
		return finalInvoker
	***REMOVED***
	return func(ctx context.Context, method string, req, reply interface***REMOVED******REMOVED***, cc *ClientConn, opts ...CallOption) error ***REMOVED***
		return interceptors[curr+1](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, curr+1, finalInvoker), opts...)
	***REMOVED***
***REMOVED***

// chainStreamClientInterceptors chains all stream client interceptors into one.
func chainStreamClientInterceptors(cc *ClientConn) ***REMOVED***
	interceptors := cc.dopts.chainStreamInts
	// Prepend dopts.streamInt to the chaining interceptors if it exists, since streamInt will
	// be executed before any other chained interceptors.
	if cc.dopts.streamInt != nil ***REMOVED***
		interceptors = append([]StreamClientInterceptor***REMOVED***cc.dopts.streamInt***REMOVED***, interceptors...)
	***REMOVED***
	var chainedInt StreamClientInterceptor
	if len(interceptors) == 0 ***REMOVED***
		chainedInt = nil
	***REMOVED*** else if len(interceptors) == 1 ***REMOVED***
		chainedInt = interceptors[0]
	***REMOVED*** else ***REMOVED***
		chainedInt = func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error) ***REMOVED***
			return interceptors[0](ctx, desc, cc, method, getChainStreamer(interceptors, 0, streamer), opts...)
		***REMOVED***
	***REMOVED***
	cc.dopts.streamInt = chainedInt
***REMOVED***

// getChainStreamer recursively generate the chained client stream constructor.
func getChainStreamer(interceptors []StreamClientInterceptor, curr int, finalStreamer Streamer) Streamer ***REMOVED***
	if curr == len(interceptors)-1 ***REMOVED***
		return finalStreamer
	***REMOVED***
	return func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error) ***REMOVED***
		return interceptors[curr+1](ctx, desc, cc, method, getChainStreamer(interceptors, curr+1, finalStreamer), opts...)
	***REMOVED***
***REMOVED***

// connectivityStateManager keeps the connectivity.State of ClientConn.
// This struct will eventually be exported so the balancers can access it.
type connectivityStateManager struct ***REMOVED***
	mu         sync.Mutex
	state      connectivity.State
	notifyChan chan struct***REMOVED******REMOVED***
	channelzID *channelz.Identifier
***REMOVED***

// updateState updates the connectivity.State of ClientConn.
// If there's a change it notifies goroutines waiting on state change to
// happen.
func (csm *connectivityStateManager) updateState(state connectivity.State) ***REMOVED***
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.state == connectivity.Shutdown ***REMOVED***
		return
	***REMOVED***
	if csm.state == state ***REMOVED***
		return
	***REMOVED***
	csm.state = state
	channelz.Infof(logger, csm.channelzID, "Channel Connectivity change to %v", state)
	if csm.notifyChan != nil ***REMOVED***
		// There are other goroutines waiting on this channel.
		close(csm.notifyChan)
		csm.notifyChan = nil
	***REMOVED***
***REMOVED***

func (csm *connectivityStateManager) getState() connectivity.State ***REMOVED***
	csm.mu.Lock()
	defer csm.mu.Unlock()
	return csm.state
***REMOVED***

func (csm *connectivityStateManager) getNotifyChan() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.notifyChan == nil ***REMOVED***
		csm.notifyChan = make(chan struct***REMOVED******REMOVED***)
	***REMOVED***
	return csm.notifyChan
***REMOVED***

// ClientConnInterface defines the functions clients need to perform unary and
// streaming RPCs.  It is implemented by *ClientConn, and is only intended to
// be referenced by generated code.
type ClientConnInterface interface ***REMOVED***
	// Invoke performs a unary RPC and returns after the response is received
	// into reply.
	Invoke(ctx context.Context, method string, args interface***REMOVED******REMOVED***, reply interface***REMOVED******REMOVED***, opts ...CallOption) error
	// NewStream begins a streaming RPC.
	NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error)
***REMOVED***

// Assert *ClientConn implements ClientConnInterface.
var _ ClientConnInterface = (*ClientConn)(nil)

// ClientConn represents a virtual connection to a conceptual endpoint, to
// perform RPCs.
//
// A ClientConn is free to have zero or more actual connections to the endpoint
// based on configuration, load, etc. It is also free to determine which actual
// endpoints to use and may change it every RPC, permitting client-side load
// balancing.
//
// A ClientConn encapsulates a range of functionality including name
// resolution, TCP connection establishment (with retries and backoff) and TLS
// handshakes. It also handles errors on established connections by
// re-resolving the name and reconnecting.
type ClientConn struct ***REMOVED***
	ctx    context.Context    // Initialized using the background context at dial time.
	cancel context.CancelFunc // Cancelled on close.

	// The following are initialized at dial time, and are read-only after that.
	target          string               // User's dial target.
	parsedTarget    resolver.Target      // See parseTargetAndFindResolver().
	authority       string               // See determineAuthority().
	dopts           dialOptions          // Default and user specified dial options.
	channelzID      *channelz.Identifier // Channelz identifier for the channel.
	balancerWrapper *ccBalancerWrapper   // Uses gracefulswitch.balancer underneath.

	// The following provide their own synchronization, and therefore don't
	// require cc.mu to be held to access them.
	csMgr              *connectivityStateManager
	blockingpicker     *pickerWrapper
	safeConfigSelector iresolver.SafeConfigSelector
	czData             *channelzData
	retryThrottler     atomic.Value // Updated from service config.

	// firstResolveEvent is used to track whether the name resolver sent us at
	// least one update. RPCs block on this event.
	firstResolveEvent *grpcsync.Event

	// mu protects the following fields.
	// TODO: split mu so the same mutex isn't used for everything.
	mu              sync.RWMutex
	resolverWrapper *ccResolverWrapper         // Initialized in Dial; cleared in Close.
	sc              *ServiceConfig             // Latest service config received from the resolver.
	conns           map[*addrConn]struct***REMOVED******REMOVED***     // Set to nil on close.
	mkp             keepalive.ClientParameters // May be updated upon receipt of a GoAway.

	lceMu               sync.Mutex // protects lastConnectionError
	lastConnectionError error
***REMOVED***

// WaitForStateChange waits until the connectivity.State of ClientConn changes from sourceState or
// ctx expires. A true value is returned in former case and false in latter.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (cc *ClientConn) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool ***REMOVED***
	ch := cc.csMgr.getNotifyChan()
	if cc.csMgr.getState() != sourceState ***REMOVED***
		return true
	***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return false
	case <-ch:
		return true
	***REMOVED***
***REMOVED***

// GetState returns the connectivity.State of ClientConn.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a later
// release.
func (cc *ClientConn) GetState() connectivity.State ***REMOVED***
	return cc.csMgr.getState()
***REMOVED***

// Connect causes all subchannels in the ClientConn to attempt to connect if
// the channel is idle.  Does not wait for the connection attempts to begin
// before returning.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a later
// release.
func (cc *ClientConn) Connect() ***REMOVED***
	cc.balancerWrapper.exitIdle()
***REMOVED***

func (cc *ClientConn) scWatcher() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case sc, ok := <-cc.dopts.scChan:
			if !ok ***REMOVED***
				return
			***REMOVED***
			cc.mu.Lock()
			// TODO: load balance policy runtime change is ignored.
			// We may revisit this decision in the future.
			cc.sc = &sc
			cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector***REMOVED***&sc***REMOVED***)
			cc.mu.Unlock()
		case <-cc.ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// waitForResolvedAddrs blocks until the resolver has provided addresses or the
// context expires.  Returns nil unless the context expires first; otherwise
// returns a status error based on the context.
func (cc *ClientConn) waitForResolvedAddrs(ctx context.Context) error ***REMOVED***
	// This is on the RPC path, so we use a fast path to avoid the
	// more-expensive "select" below after the resolver has returned once.
	if cc.firstResolveEvent.HasFired() ***REMOVED***
		return nil
	***REMOVED***
	select ***REMOVED***
	case <-cc.firstResolveEvent.Done():
		return nil
	case <-ctx.Done():
		return status.FromContextError(ctx.Err()).Err()
	case <-cc.ctx.Done():
		return ErrClientConnClosing
	***REMOVED***
***REMOVED***

var emptyServiceConfig *ServiceConfig

func init() ***REMOVED***
	cfg := parseServiceConfig("***REMOVED******REMOVED***")
	if cfg.Err != nil ***REMOVED***
		panic(fmt.Sprintf("impossible error parsing empty service config: %v", cfg.Err))
	***REMOVED***
	emptyServiceConfig = cfg.Config.(*ServiceConfig)
***REMOVED***

func (cc *ClientConn) maybeApplyDefaultServiceConfig(addrs []resolver.Address) ***REMOVED***
	if cc.sc != nil ***REMOVED***
		cc.applyServiceConfigAndBalancer(cc.sc, nil, addrs)
		return
	***REMOVED***
	if cc.dopts.defaultServiceConfig != nil ***REMOVED***
		cc.applyServiceConfigAndBalancer(cc.dopts.defaultServiceConfig, &defaultConfigSelector***REMOVED***cc.dopts.defaultServiceConfig***REMOVED***, addrs)
	***REMOVED*** else ***REMOVED***
		cc.applyServiceConfigAndBalancer(emptyServiceConfig, &defaultConfigSelector***REMOVED***emptyServiceConfig***REMOVED***, addrs)
	***REMOVED***
***REMOVED***

func (cc *ClientConn) updateResolverState(s resolver.State, err error) error ***REMOVED***
	defer cc.firstResolveEvent.Fire()
	cc.mu.Lock()
	// Check if the ClientConn is already closed. Some fields (e.g.
	// balancerWrapper) are set to nil when closing the ClientConn, and could
	// cause nil pointer panic if we don't have this check.
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return nil
	***REMOVED***

	if err != nil ***REMOVED***
		// May need to apply the initial service config in case the resolver
		// doesn't support service configs, or doesn't provide a service config
		// with the new addresses.
		cc.maybeApplyDefaultServiceConfig(nil)

		cc.balancerWrapper.resolverError(err)

		// No addresses are valid with err set; return early.
		cc.mu.Unlock()
		return balancer.ErrBadResolverState
	***REMOVED***

	var ret error
	if cc.dopts.disableServiceConfig ***REMOVED***
		channelz.Infof(logger, cc.channelzID, "ignoring service config from resolver (%v) and applying the default because service config is disabled", s.ServiceConfig)
		cc.maybeApplyDefaultServiceConfig(s.Addresses)
	***REMOVED*** else if s.ServiceConfig == nil ***REMOVED***
		cc.maybeApplyDefaultServiceConfig(s.Addresses)
		// TODO: do we need to apply a failing LB policy if there is no
		// default, per the error handling design?
	***REMOVED*** else ***REMOVED***
		if sc, ok := s.ServiceConfig.Config.(*ServiceConfig); s.ServiceConfig.Err == nil && ok ***REMOVED***
			configSelector := iresolver.GetConfigSelector(s)
			if configSelector != nil ***REMOVED***
				if len(s.ServiceConfig.Config.(*ServiceConfig).Methods) != 0 ***REMOVED***
					channelz.Infof(logger, cc.channelzID, "method configs in service config will be ignored due to presence of config selector")
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				configSelector = &defaultConfigSelector***REMOVED***sc***REMOVED***
			***REMOVED***
			cc.applyServiceConfigAndBalancer(sc, configSelector, s.Addresses)
		***REMOVED*** else ***REMOVED***
			ret = balancer.ErrBadResolverState
			if cc.sc == nil ***REMOVED***
				// Apply the failing LB only if we haven't received valid service config
				// from the name resolver in the past.
				cc.applyFailingLB(s.ServiceConfig)
				cc.mu.Unlock()
				return ret
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var balCfg serviceconfig.LoadBalancingConfig
	if cc.sc != nil && cc.sc.lbConfig != nil ***REMOVED***
		balCfg = cc.sc.lbConfig.cfg
	***REMOVED***
	bw := cc.balancerWrapper
	cc.mu.Unlock()

	uccsErr := bw.updateClientConnState(&balancer.ClientConnState***REMOVED***ResolverState: s, BalancerConfig: balCfg***REMOVED***)
	if ret == nil ***REMOVED***
		ret = uccsErr // prefer ErrBadResolver state since any other error is
		// currently meaningless to the caller.
	***REMOVED***
	return ret
***REMOVED***

// applyFailingLB is akin to configuring an LB policy on the channel which
// always fails RPCs. Here, an actual LB policy is not configured, but an always
// erroring picker is configured, which returns errors with information about
// what was invalid in the received service config. A config selector with no
// service config is configured, and the connectivity state of the channel is
// set to TransientFailure.
//
// Caller must hold cc.mu.
func (cc *ClientConn) applyFailingLB(sc *serviceconfig.ParseResult) ***REMOVED***
	var err error
	if sc.Err != nil ***REMOVED***
		err = status.Errorf(codes.Unavailable, "error parsing service config: %v", sc.Err)
	***REMOVED*** else ***REMOVED***
		err = status.Errorf(codes.Unavailable, "illegal service config type: %T", sc.Config)
	***REMOVED***
	cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector***REMOVED***nil***REMOVED***)
	cc.blockingpicker.updatePicker(base.NewErrPicker(err))
	cc.csMgr.updateState(connectivity.TransientFailure)
***REMOVED***

func (cc *ClientConn) handleSubConnStateChange(sc balancer.SubConn, s connectivity.State, err error) ***REMOVED***
	cc.balancerWrapper.updateSubConnState(sc, s, err)
***REMOVED***

// newAddrConn creates an addrConn for addrs and adds it to cc.conns.
//
// Caller needs to make sure len(addrs) > 0.
func (cc *ClientConn) newAddrConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (*addrConn, error) ***REMOVED***
	ac := &addrConn***REMOVED***
		state:        connectivity.Idle,
		cc:           cc,
		addrs:        addrs,
		scopts:       opts,
		dopts:        cc.dopts,
		czData:       new(channelzData),
		resetBackoff: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	ac.ctx, ac.cancel = context.WithCancel(cc.ctx)
	// Track ac in cc. This needs to be done before any getTransport(...) is called.
	cc.mu.Lock()
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return nil, ErrClientConnClosing
	***REMOVED***

	var err error
	ac.channelzID, err = channelz.RegisterSubChannel(ac, cc.channelzID, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	channelz.AddTraceEvent(logger, ac.channelzID, 0, &channelz.TraceEventDesc***REMOVED***
		Desc:     "Subchannel created",
		Severity: channelz.CtInfo,
		Parent: &channelz.TraceEventDesc***REMOVED***
			Desc:     fmt.Sprintf("Subchannel(id:%d) created", ac.channelzID.Int()),
			Severity: channelz.CtInfo,
		***REMOVED***,
	***REMOVED***)

	cc.conns[ac] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	cc.mu.Unlock()
	return ac, nil
***REMOVED***

// removeAddrConn removes the addrConn in the subConn from clientConn.
// It also tears down the ac with the given error.
func (cc *ClientConn) removeAddrConn(ac *addrConn, err error) ***REMOVED***
	cc.mu.Lock()
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return
	***REMOVED***
	delete(cc.conns, ac)
	cc.mu.Unlock()
	ac.tearDown(err)
***REMOVED***

func (cc *ClientConn) channelzMetric() *channelz.ChannelInternalMetric ***REMOVED***
	return &channelz.ChannelInternalMetric***REMOVED***
		State:                    cc.GetState(),
		Target:                   cc.target,
		CallsStarted:             atomic.LoadInt64(&cc.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&cc.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&cc.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&cc.czData.lastCallStartedTime)),
	***REMOVED***
***REMOVED***

// Target returns the target string of the ClientConn.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (cc *ClientConn) Target() string ***REMOVED***
	return cc.target
***REMOVED***

func (cc *ClientConn) incrCallsStarted() ***REMOVED***
	atomic.AddInt64(&cc.czData.callsStarted, 1)
	atomic.StoreInt64(&cc.czData.lastCallStartedTime, time.Now().UnixNano())
***REMOVED***

func (cc *ClientConn) incrCallsSucceeded() ***REMOVED***
	atomic.AddInt64(&cc.czData.callsSucceeded, 1)
***REMOVED***

func (cc *ClientConn) incrCallsFailed() ***REMOVED***
	atomic.AddInt64(&cc.czData.callsFailed, 1)
***REMOVED***

// connect starts creating a transport.
// It does nothing if the ac is not IDLE.
// TODO(bar) Move this to the addrConn section.
func (ac *addrConn) connect() error ***REMOVED***
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown ***REMOVED***
		ac.mu.Unlock()
		return errConnClosing
	***REMOVED***
	if ac.state != connectivity.Idle ***REMOVED***
		ac.mu.Unlock()
		return nil
	***REMOVED***
	// Update connectivity state within the lock to prevent subsequent or
	// concurrent calls from resetting the transport more than once.
	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()

	ac.resetTransport()
	return nil
***REMOVED***

func equalAddresses(a, b []resolver.Address) bool ***REMOVED***
	if len(a) != len(b) ***REMOVED***
		return false
	***REMOVED***
	for i, v := range a ***REMOVED***
		if !v.Equal(b[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// tryUpdateAddrs tries to update ac.addrs with the new addresses list.
//
// If ac is TransientFailure, it updates ac.addrs and returns true. The updated
// addresses will be picked up by retry in the next iteration after backoff.
//
// If ac is Shutdown or Idle, it updates ac.addrs and returns true.
//
// If the addresses is the same as the old list, it does nothing and returns
// true.
//
// If ac is Connecting, it returns false. The caller should tear down the ac and
// create a new one. Note that the backoff will be reset when this happens.
//
// If ac is Ready, it checks whether current connected address of ac is in the
// new addrs list.
//  - If true, it updates ac.addrs and returns true. The ac will keep using
//    the existing connection.
//  - If false, it does nothing and returns false.
func (ac *addrConn) tryUpdateAddrs(addrs []resolver.Address) bool ***REMOVED***
	ac.mu.Lock()
	defer ac.mu.Unlock()
	channelz.Infof(logger, ac.channelzID, "addrConn: tryUpdateAddrs curAddr: %v, addrs: %v", ac.curAddr, addrs)
	if ac.state == connectivity.Shutdown ||
		ac.state == connectivity.TransientFailure ||
		ac.state == connectivity.Idle ***REMOVED***
		ac.addrs = addrs
		return true
	***REMOVED***

	if equalAddresses(ac.addrs, addrs) ***REMOVED***
		return true
	***REMOVED***

	if ac.state == connectivity.Connecting ***REMOVED***
		return false
	***REMOVED***

	// ac.state is Ready, try to find the connected address.
	var curAddrFound bool
	for _, a := range addrs ***REMOVED***
		a.ServerName = ac.cc.getServerName(a)
		if reflect.DeepEqual(ac.curAddr, a) ***REMOVED***
			curAddrFound = true
			break
		***REMOVED***
	***REMOVED***
	channelz.Infof(logger, ac.channelzID, "addrConn: tryUpdateAddrs curAddrFound: %v", curAddrFound)
	if curAddrFound ***REMOVED***
		ac.addrs = addrs
	***REMOVED***

	return curAddrFound
***REMOVED***

// getServerName determines the serverName to be used in the connection
// handshake. The default value for the serverName is the authority on the
// ClientConn, which either comes from the user's dial target or through an
// authority override specified using the WithAuthority dial option. Name
// resolvers can specify a per-address override for the serverName through the
// resolver.Address.ServerName field which is used only if the WithAuthority
// dial option was not used. The rationale is that per-address authority
// overrides specified by the name resolver can represent a security risk, while
// an override specified by the user is more dependable since they probably know
// what they are doing.
func (cc *ClientConn) getServerName(addr resolver.Address) string ***REMOVED***
	if cc.dopts.authority != "" ***REMOVED***
		return cc.dopts.authority
	***REMOVED***
	if addr.ServerName != "" ***REMOVED***
		return addr.ServerName
	***REMOVED***
	return cc.authority
***REMOVED***

func getMethodConfig(sc *ServiceConfig, method string) MethodConfig ***REMOVED***
	if sc == nil ***REMOVED***
		return MethodConfig***REMOVED******REMOVED***
	***REMOVED***
	if m, ok := sc.Methods[method]; ok ***REMOVED***
		return m
	***REMOVED***
	i := strings.LastIndex(method, "/")
	if m, ok := sc.Methods[method[:i+1]]; ok ***REMOVED***
		return m
	***REMOVED***
	return sc.Methods[""]
***REMOVED***

// GetMethodConfig gets the method config of the input method.
// If there's an exact match for input method (i.e. /service/method), we return
// the corresponding MethodConfig.
// If there isn't an exact match for the input method, we look for the service's default
// config under the service (i.e /service/) and then for the default for all services (empty string).
//
// If there is a default MethodConfig for the service, we return it.
// Otherwise, we return an empty MethodConfig.
func (cc *ClientConn) GetMethodConfig(method string) MethodConfig ***REMOVED***
	// TODO: Avoid the locking here.
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return getMethodConfig(cc.sc, method)
***REMOVED***

func (cc *ClientConn) healthCheckConfig() *healthCheckConfig ***REMOVED***
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.sc == nil ***REMOVED***
		return nil
	***REMOVED***
	return cc.sc.healthCheckConfig
***REMOVED***

func (cc *ClientConn) getTransport(ctx context.Context, failfast bool, method string) (transport.ClientTransport, func(balancer.DoneInfo), error) ***REMOVED***
	return cc.blockingpicker.pick(ctx, failfast, balancer.PickInfo***REMOVED***
		Ctx:            ctx,
		FullMethodName: method,
	***REMOVED***)
***REMOVED***

func (cc *ClientConn) applyServiceConfigAndBalancer(sc *ServiceConfig, configSelector iresolver.ConfigSelector, addrs []resolver.Address) ***REMOVED***
	if sc == nil ***REMOVED***
		// should never reach here.
		return
	***REMOVED***
	cc.sc = sc
	if configSelector != nil ***REMOVED***
		cc.safeConfigSelector.UpdateConfigSelector(configSelector)
	***REMOVED***

	if cc.sc.retryThrottling != nil ***REMOVED***
		newThrottler := &retryThrottler***REMOVED***
			tokens: cc.sc.retryThrottling.MaxTokens,
			max:    cc.sc.retryThrottling.MaxTokens,
			thresh: cc.sc.retryThrottling.MaxTokens / 2,
			ratio:  cc.sc.retryThrottling.TokenRatio,
		***REMOVED***
		cc.retryThrottler.Store(newThrottler)
	***REMOVED*** else ***REMOVED***
		cc.retryThrottler.Store((*retryThrottler)(nil))
	***REMOVED***

	var newBalancerName string
	if cc.sc != nil && cc.sc.lbConfig != nil ***REMOVED***
		newBalancerName = cc.sc.lbConfig.name
	***REMOVED*** else ***REMOVED***
		var isGRPCLB bool
		for _, a := range addrs ***REMOVED***
			if a.Type == resolver.GRPCLB ***REMOVED***
				isGRPCLB = true
				break
			***REMOVED***
		***REMOVED***
		if isGRPCLB ***REMOVED***
			newBalancerName = grpclbName
		***REMOVED*** else if cc.sc != nil && cc.sc.LB != nil ***REMOVED***
			newBalancerName = *cc.sc.LB
		***REMOVED*** else ***REMOVED***
			newBalancerName = PickFirstBalancerName
		***REMOVED***
	***REMOVED***
	cc.balancerWrapper.switchTo(newBalancerName)
***REMOVED***

func (cc *ClientConn) resolveNow(o resolver.ResolveNowOptions) ***REMOVED***
	cc.mu.RLock()
	r := cc.resolverWrapper
	cc.mu.RUnlock()
	if r == nil ***REMOVED***
		return
	***REMOVED***
	go r.resolveNow(o)
***REMOVED***

// ResetConnectBackoff wakes up all subchannels in transient failure and causes
// them to attempt another connection immediately.  It also resets the backoff
// times used for subsequent attempts regardless of the current state.
//
// In general, this function should not be used.  Typical service or network
// outages result in a reasonable client reconnection strategy by default.
// However, if a previously unavailable network becomes available, this may be
// used to trigger an immediate reconnect.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (cc *ClientConn) ResetConnectBackoff() ***REMOVED***
	cc.mu.Lock()
	conns := cc.conns
	cc.mu.Unlock()
	for ac := range conns ***REMOVED***
		ac.resetConnectBackoff()
	***REMOVED***
***REMOVED***

// Close tears down the ClientConn and all underlying connections.
func (cc *ClientConn) Close() error ***REMOVED***
	defer cc.cancel()

	cc.mu.Lock()
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return ErrClientConnClosing
	***REMOVED***
	conns := cc.conns
	cc.conns = nil
	cc.csMgr.updateState(connectivity.Shutdown)

	rWrapper := cc.resolverWrapper
	cc.resolverWrapper = nil
	bWrapper := cc.balancerWrapper
	cc.mu.Unlock()

	// The order of closing matters here since the balancer wrapper assumes the
	// picker is closed before it is closed.
	cc.blockingpicker.close()
	if bWrapper != nil ***REMOVED***
		bWrapper.close()
	***REMOVED***
	if rWrapper != nil ***REMOVED***
		rWrapper.close()
	***REMOVED***

	for ac := range conns ***REMOVED***
		ac.tearDown(ErrClientConnClosing)
	***REMOVED***
	ted := &channelz.TraceEventDesc***REMOVED***
		Desc:     "Channel deleted",
		Severity: channelz.CtInfo,
	***REMOVED***
	if cc.dopts.channelzParentID != nil ***REMOVED***
		ted.Parent = &channelz.TraceEventDesc***REMOVED***
			Desc:     fmt.Sprintf("Nested channel(id:%d) deleted", cc.channelzID.Int()),
			Severity: channelz.CtInfo,
		***REMOVED***
	***REMOVED***
	channelz.AddTraceEvent(logger, cc.channelzID, 0, ted)
	// TraceEvent needs to be called before RemoveEntry, as TraceEvent may add
	// trace reference to the entity being deleted, and thus prevent it from being
	// deleted right away.
	channelz.RemoveEntry(cc.channelzID)

	return nil
***REMOVED***

// addrConn is a network connection to a given address.
type addrConn struct ***REMOVED***
	ctx    context.Context
	cancel context.CancelFunc

	cc     *ClientConn
	dopts  dialOptions
	acbw   balancer.SubConn
	scopts balancer.NewSubConnOptions

	// transport is set when there's a viable transport (note: ac state may not be READY as LB channel
	// health checking may require server to report healthy to set ac to READY), and is reset
	// to nil when the current transport should no longer be used to create a stream (e.g. after GoAway
	// is received, transport is closed, ac has been torn down).
	transport transport.ClientTransport // The current transport.

	mu      sync.Mutex
	curAddr resolver.Address   // The current address.
	addrs   []resolver.Address // All addresses that the resolver resolved to.

	// Use updateConnectivityState for updating addrConn's connectivity state.
	state connectivity.State

	backoffIdx   int // Needs to be stateful for resetConnectBackoff.
	resetBackoff chan struct***REMOVED******REMOVED***

	channelzID *channelz.Identifier
	czData     *channelzData
***REMOVED***

// Note: this requires a lock on ac.mu.
func (ac *addrConn) updateConnectivityState(s connectivity.State, lastErr error) ***REMOVED***
	if ac.state == s ***REMOVED***
		return
	***REMOVED***
	ac.state = s
	channelz.Infof(logger, ac.channelzID, "Subchannel Connectivity change to %v", s)
	ac.cc.handleSubConnStateChange(ac.acbw, s, lastErr)
***REMOVED***

// adjustParams updates parameters used to create transports upon
// receiving a GoAway.
func (ac *addrConn) adjustParams(r transport.GoAwayReason) ***REMOVED***
	switch r ***REMOVED***
	case transport.GoAwayTooManyPings:
		v := 2 * ac.dopts.copts.KeepaliveParams.Time
		ac.cc.mu.Lock()
		if v > ac.cc.mkp.Time ***REMOVED***
			ac.cc.mkp.Time = v
		***REMOVED***
		ac.cc.mu.Unlock()
	***REMOVED***
***REMOVED***

func (ac *addrConn) resetTransport() ***REMOVED***
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown ***REMOVED***
		ac.mu.Unlock()
		return
	***REMOVED***

	addrs := ac.addrs
	backoffFor := ac.dopts.bs.Backoff(ac.backoffIdx)
	// This will be the duration that dial gets to finish.
	dialDuration := minConnectTimeout
	if ac.dopts.minConnectTimeout != nil ***REMOVED***
		dialDuration = ac.dopts.minConnectTimeout()
	***REMOVED***

	if dialDuration < backoffFor ***REMOVED***
		// Give dial more time as we keep failing to connect.
		dialDuration = backoffFor
	***REMOVED***
	// We can potentially spend all the time trying the first address, and
	// if the server accepts the connection and then hangs, the following
	// addresses will never be tried.
	//
	// The spec doesn't mention what should be done for multiple addresses.
	// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md#proposed-backoff-algorithm
	connectDeadline := time.Now().Add(dialDuration)

	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()

	if err := ac.tryAllAddrs(addrs, connectDeadline); err != nil ***REMOVED***
		ac.cc.resolveNow(resolver.ResolveNowOptions***REMOVED******REMOVED***)
		// After exhausting all addresses, the addrConn enters
		// TRANSIENT_FAILURE.
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown ***REMOVED***
			ac.mu.Unlock()
			return
		***REMOVED***
		ac.updateConnectivityState(connectivity.TransientFailure, err)

		// Backoff.
		b := ac.resetBackoff
		ac.mu.Unlock()

		timer := time.NewTimer(backoffFor)
		select ***REMOVED***
		case <-timer.C:
			ac.mu.Lock()
			ac.backoffIdx++
			ac.mu.Unlock()
		case <-b:
			timer.Stop()
		case <-ac.ctx.Done():
			timer.Stop()
			return
		***REMOVED***

		ac.mu.Lock()
		if ac.state != connectivity.Shutdown ***REMOVED***
			ac.updateConnectivityState(connectivity.Idle, err)
		***REMOVED***
		ac.mu.Unlock()
		return
	***REMOVED***
	// Success; reset backoff.
	ac.mu.Lock()
	ac.backoffIdx = 0
	ac.mu.Unlock()
***REMOVED***

// tryAllAddrs tries to creates a connection to the addresses, and stop when at
// the first successful one. It returns an error if no address was successfully
// connected, or updates ac appropriately with the new transport.
func (ac *addrConn) tryAllAddrs(addrs []resolver.Address, connectDeadline time.Time) error ***REMOVED***
	var firstConnErr error
	for _, addr := range addrs ***REMOVED***
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown ***REMOVED***
			ac.mu.Unlock()
			return errConnClosing
		***REMOVED***

		ac.cc.mu.RLock()
		ac.dopts.copts.KeepaliveParams = ac.cc.mkp
		ac.cc.mu.RUnlock()

		copts := ac.dopts.copts
		if ac.scopts.CredsBundle != nil ***REMOVED***
			copts.CredsBundle = ac.scopts.CredsBundle
		***REMOVED***
		ac.mu.Unlock()

		channelz.Infof(logger, ac.channelzID, "Subchannel picks a new address %q to connect", addr.Addr)

		err := ac.createTransport(addr, copts, connectDeadline)
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
		if firstConnErr == nil ***REMOVED***
			firstConnErr = err
		***REMOVED***
		ac.cc.updateConnectionError(err)
	***REMOVED***

	// Couldn't connect to any address.
	return firstConnErr
***REMOVED***

// createTransport creates a connection to addr. It returns an error if the
// address was not successfully connected, or updates ac appropriately with the
// new transport.
func (ac *addrConn) createTransport(addr resolver.Address, copts transport.ConnectOptions, connectDeadline time.Time) error ***REMOVED***
	// TODO: Delete prefaceReceived and move the logic to wait for it into the
	// transport.
	prefaceReceived := grpcsync.NewEvent()
	connClosed := grpcsync.NewEvent()

	addr.ServerName = ac.cc.getServerName(addr)
	hctx, hcancel := context.WithCancel(ac.ctx)
	hcStarted := false // protected by ac.mu

	onClose := func() ***REMOVED***
		ac.mu.Lock()
		defer ac.mu.Unlock()
		defer connClosed.Fire()
		defer hcancel()
		if !hcStarted || hctx.Err() != nil ***REMOVED***
			// We didn't start the health check or set the state to READY, so
			// no need to do anything else here.
			//
			// OR, we have already cancelled the health check context, meaning
			// we have already called onClose once for this transport.  In this
			// case it would be dangerous to clear the transport and update the
			// state, since there may be a new transport in this addrConn.
			return
		***REMOVED***
		ac.transport = nil
		// Refresh the name resolver
		ac.cc.resolveNow(resolver.ResolveNowOptions***REMOVED******REMOVED***)
		if ac.state != connectivity.Shutdown ***REMOVED***
			ac.updateConnectivityState(connectivity.Idle, nil)
		***REMOVED***
	***REMOVED***

	onGoAway := func(r transport.GoAwayReason) ***REMOVED***
		ac.mu.Lock()
		ac.adjustParams(r)
		ac.mu.Unlock()
		onClose()
	***REMOVED***

	connectCtx, cancel := context.WithDeadline(ac.ctx, connectDeadline)
	defer cancel()
	copts.ChannelzParentID = ac.channelzID

	newTr, err := transport.NewClientTransport(connectCtx, ac.cc.ctx, addr, copts, func() ***REMOVED*** prefaceReceived.Fire() ***REMOVED***, onGoAway, onClose)
	if err != nil ***REMOVED***
		// newTr is either nil, or closed.
		hcancel()
		channelz.Warningf(logger, ac.channelzID, "grpc: addrConn.createTransport failed to connect to %s. Err: %v", addr, err)
		return err
	***REMOVED***

	select ***REMOVED***
	case <-connectCtx.Done():
		// We didn't get the preface in time.
		// The error we pass to Close() is immaterial since there are no open
		// streams at this point, so no trailers with error details will be sent
		// out. We just need to pass a non-nil error.
		newTr.Close(transport.ErrConnClosing)
		if connectCtx.Err() == context.DeadlineExceeded ***REMOVED***
			err := errors.New("failed to receive server preface within timeout")
			channelz.Warningf(logger, ac.channelzID, "grpc: addrConn.createTransport failed to connect to %s: %v", addr, err)
			return err
		***REMOVED***
		return nil
	case <-prefaceReceived.Done():
		// We got the preface - huzzah! things are good.
		ac.mu.Lock()
		defer ac.mu.Unlock()
		if connClosed.HasFired() ***REMOVED***
			// onClose called first; go idle but do nothing else.
			if ac.state != connectivity.Shutdown ***REMOVED***
				ac.updateConnectivityState(connectivity.Idle, nil)
			***REMOVED***
			return nil
		***REMOVED***
		if ac.state == connectivity.Shutdown ***REMOVED***
			// This can happen if the subConn was removed while in `Connecting`
			// state. tearDown() would have set the state to `Shutdown`, but
			// would not have closed the transport since ac.transport would not
			// been set at that point.
			//
			// We run this in a goroutine because newTr.Close() calls onClose()
			// inline, which requires locking ac.mu.
			//
			// The error we pass to Close() is immaterial since there are no open
			// streams at this point, so no trailers with error details will be sent
			// out. We just need to pass a non-nil error.
			go newTr.Close(transport.ErrConnClosing)
			return nil
		***REMOVED***
		ac.curAddr = addr
		ac.transport = newTr
		hcStarted = true
		ac.startHealthCheck(hctx) // Will set state to READY if appropriate.
		return nil
	case <-connClosed.Done():
		// The transport has already closed.  If we received the preface, too,
		// this is not an error.
		select ***REMOVED***
		case <-prefaceReceived.Done():
			return nil
		default:
			return errors.New("connection closed before server preface received")
		***REMOVED***
	***REMOVED***
***REMOVED***

// startHealthCheck starts the health checking stream (RPC) to watch the health
// stats of this connection if health checking is requested and configured.
//
// LB channel health checking is enabled when all requirements below are met:
// 1. it is not disabled by the user with the WithDisableHealthCheck DialOption
// 2. internal.HealthCheckFunc is set by importing the grpc/health package
// 3. a service config with non-empty healthCheckConfig field is provided
// 4. the load balancer requests it
//
// It sets addrConn to READY if the health checking stream is not started.
//
// Caller must hold ac.mu.
func (ac *addrConn) startHealthCheck(ctx context.Context) ***REMOVED***
	var healthcheckManagingState bool
	defer func() ***REMOVED***
		if !healthcheckManagingState ***REMOVED***
			ac.updateConnectivityState(connectivity.Ready, nil)
		***REMOVED***
	***REMOVED***()

	if ac.cc.dopts.disableHealthCheck ***REMOVED***
		return
	***REMOVED***
	healthCheckConfig := ac.cc.healthCheckConfig()
	if healthCheckConfig == nil ***REMOVED***
		return
	***REMOVED***
	if !ac.scopts.HealthCheckEnabled ***REMOVED***
		return
	***REMOVED***
	healthCheckFunc := ac.cc.dopts.healthCheckFunc
	if healthCheckFunc == nil ***REMOVED***
		// The health package is not imported to set health check function.
		//
		// TODO: add a link to the health check doc in the error message.
		channelz.Error(logger, ac.channelzID, "Health check is requested but health check function is not set.")
		return
	***REMOVED***

	healthcheckManagingState = true

	// Set up the health check helper functions.
	currentTr := ac.transport
	newStream := func(method string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		ac.mu.Lock()
		if ac.transport != currentTr ***REMOVED***
			ac.mu.Unlock()
			return nil, status.Error(codes.Canceled, "the provided transport is no longer valid to use")
		***REMOVED***
		ac.mu.Unlock()
		return newNonRetryClientStream(ctx, &StreamDesc***REMOVED***ServerStreams: true***REMOVED***, method, currentTr, ac)
	***REMOVED***
	setConnectivityState := func(s connectivity.State, lastErr error) ***REMOVED***
		ac.mu.Lock()
		defer ac.mu.Unlock()
		if ac.transport != currentTr ***REMOVED***
			return
		***REMOVED***
		ac.updateConnectivityState(s, lastErr)
	***REMOVED***
	// Start the health checking stream.
	go func() ***REMOVED***
		err := ac.cc.dopts.healthCheckFunc(ctx, newStream, setConnectivityState, healthCheckConfig.ServiceName)
		if err != nil ***REMOVED***
			if status.Code(err) == codes.Unimplemented ***REMOVED***
				channelz.Error(logger, ac.channelzID, "Subchannel health check is unimplemented at server side, thus health check is disabled")
			***REMOVED*** else ***REMOVED***
				channelz.Errorf(logger, ac.channelzID, "HealthCheckFunc exits with unexpected error %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (ac *addrConn) resetConnectBackoff() ***REMOVED***
	ac.mu.Lock()
	close(ac.resetBackoff)
	ac.backoffIdx = 0
	ac.resetBackoff = make(chan struct***REMOVED******REMOVED***)
	ac.mu.Unlock()
***REMOVED***

// getReadyTransport returns the transport if ac's state is READY or nil if not.
func (ac *addrConn) getReadyTransport() transport.ClientTransport ***REMOVED***
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.state == connectivity.Ready ***REMOVED***
		return ac.transport
	***REMOVED***
	return nil
***REMOVED***

// tearDown starts to tear down the addrConn.
//
// Note that tearDown doesn't remove ac from ac.cc.conns, so the addrConn struct
// will leak. In most cases, call cc.removeAddrConn() instead.
func (ac *addrConn) tearDown(err error) ***REMOVED***
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown ***REMOVED***
		ac.mu.Unlock()
		return
	***REMOVED***
	curTr := ac.transport
	ac.transport = nil
	// We have to set the state to Shutdown before anything else to prevent races
	// between setting the state and logic that waits on context cancellation / etc.
	ac.updateConnectivityState(connectivity.Shutdown, nil)
	ac.cancel()
	ac.curAddr = resolver.Address***REMOVED******REMOVED***
	if err == errConnDrain && curTr != nil ***REMOVED***
		// GracefulClose(...) may be executed multiple times when
		// i) receiving multiple GoAway frames from the server; or
		// ii) there are concurrent name resolver/Balancer triggered
		// address removal and GoAway.
		// We have to unlock and re-lock here because GracefulClose => Close => onClose, which requires locking ac.mu.
		ac.mu.Unlock()
		curTr.GracefulClose()
		ac.mu.Lock()
	***REMOVED***
	channelz.AddTraceEvent(logger, ac.channelzID, 0, &channelz.TraceEventDesc***REMOVED***
		Desc:     "Subchannel deleted",
		Severity: channelz.CtInfo,
		Parent: &channelz.TraceEventDesc***REMOVED***
			Desc:     fmt.Sprintf("Subchannel(id:%d) deleted", ac.channelzID.Int()),
			Severity: channelz.CtInfo,
		***REMOVED***,
	***REMOVED***)
	// TraceEvent needs to be called before RemoveEntry, as TraceEvent may add
	// trace reference to the entity being deleted, and thus prevent it from
	// being deleted right away.
	channelz.RemoveEntry(ac.channelzID)
	ac.mu.Unlock()
***REMOVED***

func (ac *addrConn) getState() connectivity.State ***REMOVED***
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.state
***REMOVED***

func (ac *addrConn) ChannelzMetric() *channelz.ChannelInternalMetric ***REMOVED***
	ac.mu.Lock()
	addr := ac.curAddr.Addr
	ac.mu.Unlock()
	return &channelz.ChannelInternalMetric***REMOVED***
		State:                    ac.getState(),
		Target:                   addr,
		CallsStarted:             atomic.LoadInt64(&ac.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&ac.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&ac.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&ac.czData.lastCallStartedTime)),
	***REMOVED***
***REMOVED***

func (ac *addrConn) incrCallsStarted() ***REMOVED***
	atomic.AddInt64(&ac.czData.callsStarted, 1)
	atomic.StoreInt64(&ac.czData.lastCallStartedTime, time.Now().UnixNano())
***REMOVED***

func (ac *addrConn) incrCallsSucceeded() ***REMOVED***
	atomic.AddInt64(&ac.czData.callsSucceeded, 1)
***REMOVED***

func (ac *addrConn) incrCallsFailed() ***REMOVED***
	atomic.AddInt64(&ac.czData.callsFailed, 1)
***REMOVED***

type retryThrottler struct ***REMOVED***
	max    float64
	thresh float64
	ratio  float64

	mu     sync.Mutex
	tokens float64 // TODO(dfawley): replace with atomic and remove lock.
***REMOVED***

// throttle subtracts a retry token from the pool and returns whether a retry
// should be throttled (disallowed) based upon the retry throttling policy in
// the service config.
func (rt *retryThrottler) throttle() bool ***REMOVED***
	if rt == nil ***REMOVED***
		return false
	***REMOVED***
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens--
	if rt.tokens < 0 ***REMOVED***
		rt.tokens = 0
	***REMOVED***
	return rt.tokens <= rt.thresh
***REMOVED***

func (rt *retryThrottler) successfulRPC() ***REMOVED***
	if rt == nil ***REMOVED***
		return
	***REMOVED***
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens += rt.ratio
	if rt.tokens > rt.max ***REMOVED***
		rt.tokens = rt.max
	***REMOVED***
***REMOVED***

type channelzChannel struct ***REMOVED***
	cc *ClientConn
***REMOVED***

func (c *channelzChannel) ChannelzMetric() *channelz.ChannelInternalMetric ***REMOVED***
	return c.cc.channelzMetric()
***REMOVED***

// ErrClientConnTimeout indicates that the ClientConn cannot establish the
// underlying connections within the specified timeout.
//
// Deprecated: This error is never returned by grpc and should not be
// referenced by users.
var ErrClientConnTimeout = errors.New("grpc: timed out when dialing")

func (cc *ClientConn) getResolver(scheme string) resolver.Builder ***REMOVED***
	for _, rb := range cc.dopts.resolvers ***REMOVED***
		if scheme == rb.Scheme() ***REMOVED***
			return rb
		***REMOVED***
	***REMOVED***
	return resolver.Get(scheme)
***REMOVED***

func (cc *ClientConn) updateConnectionError(err error) ***REMOVED***
	cc.lceMu.Lock()
	cc.lastConnectionError = err
	cc.lceMu.Unlock()
***REMOVED***

func (cc *ClientConn) connectionError() error ***REMOVED***
	cc.lceMu.Lock()
	defer cc.lceMu.Unlock()
	return cc.lastConnectionError
***REMOVED***

func (cc *ClientConn) parseTargetAndFindResolver() (resolver.Builder, error) ***REMOVED***
	channelz.Infof(logger, cc.channelzID, "original dial target is: %q", cc.target)

	var rb resolver.Builder
	parsedTarget, err := parseTarget(cc.target)
	if err != nil ***REMOVED***
		channelz.Infof(logger, cc.channelzID, "dial target %q parse failed: %v", cc.target, err)
	***REMOVED*** else ***REMOVED***
		channelz.Infof(logger, cc.channelzID, "parsed dial target is: %+v", parsedTarget)
		rb = cc.getResolver(parsedTarget.Scheme)
		if rb != nil ***REMOVED***
			cc.parsedTarget = parsedTarget
			return rb, nil
		***REMOVED***
	***REMOVED***

	// We are here because the user's dial target did not contain a scheme or
	// specified an unregistered scheme. We should fallback to the default
	// scheme, except when a custom dialer is specified in which case, we should
	// always use passthrough scheme.
	defScheme := resolver.GetDefaultScheme()
	channelz.Infof(logger, cc.channelzID, "fallback to scheme %q", defScheme)
	canonicalTarget := defScheme + ":///" + cc.target

	parsedTarget, err = parseTarget(canonicalTarget)
	if err != nil ***REMOVED***
		channelz.Infof(logger, cc.channelzID, "dial target %q parse failed: %v", canonicalTarget, err)
		return nil, err
	***REMOVED***
	channelz.Infof(logger, cc.channelzID, "parsed dial target is: %+v", parsedTarget)
	rb = cc.getResolver(parsedTarget.Scheme)
	if rb == nil ***REMOVED***
		return nil, fmt.Errorf("could not get resolver for default scheme: %q", parsedTarget.Scheme)
	***REMOVED***
	cc.parsedTarget = parsedTarget
	return rb, nil
***REMOVED***

// parseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
func parseTarget(target string) (resolver.Target, error) ***REMOVED***
	u, err := url.Parse(target)
	if err != nil ***REMOVED***
		return resolver.Target***REMOVED******REMOVED***, err
	***REMOVED***
	// For targets of the form "[scheme]://[authority]/endpoint, the endpoint
	// value returned from url.Parse() contains a leading "/". Although this is
	// in accordance with RFC 3986, we do not want to break existing resolver
	// implementations which expect the endpoint without the leading "/". So, we
	// end up stripping the leading "/" here. But this will result in an
	// incorrect parsing for something like "unix:///path/to/socket". Since we
	// own the "unix" resolver, we can workaround in the unix resolver by using
	// the `URL` field instead of the `Endpoint` field.
	endpoint := u.Path
	if endpoint == "" ***REMOVED***
		endpoint = u.Opaque
	***REMOVED***
	endpoint = strings.TrimPrefix(endpoint, "/")
	return resolver.Target***REMOVED***
		Scheme:    u.Scheme,
		Authority: u.Host,
		Endpoint:  endpoint,
		URL:       *u,
	***REMOVED***, nil
***REMOVED***

// Determine channel authority. The order of precedence is as follows:
// - user specified authority override using `WithAuthority` dial option
// - creds' notion of server name for the authentication handshake
// - endpoint from dial target of the form "scheme://[authority]/endpoint"
func determineAuthority(endpoint, target string, dopts dialOptions) (string, error) ***REMOVED***
	// Historically, we had two options for users to specify the serverName or
	// authority for a channel. One was through the transport credentials
	// (either in its constructor, or through the OverrideServerName() method).
	// The other option (for cases where WithInsecure() dial option was used)
	// was to use the WithAuthority() dial option.
	//
	// A few things have changed since:
	// - `insecure` package with an implementation of the `TransportCredentials`
	//   interface for the insecure case
	// - WithAuthority() dial option support for secure credentials
	authorityFromCreds := ""
	if creds := dopts.copts.TransportCredentials; creds != nil && creds.Info().ServerName != "" ***REMOVED***
		authorityFromCreds = creds.Info().ServerName
	***REMOVED***
	authorityFromDialOption := dopts.authority
	if (authorityFromCreds != "" && authorityFromDialOption != "") && authorityFromCreds != authorityFromDialOption ***REMOVED***
		return "", fmt.Errorf("ClientConn's authority from transport creds %q and dial option %q don't match", authorityFromCreds, authorityFromDialOption)
	***REMOVED***

	switch ***REMOVED***
	case authorityFromDialOption != "":
		return authorityFromDialOption, nil
	case authorityFromCreds != "":
		return authorityFromCreds, nil
	case strings.HasPrefix(target, "unix:") || strings.HasPrefix(target, "unix-abstract:"):
		// TODO: remove when the unix resolver implements optional interface to
		// return channel authority.
		return "localhost", nil
	case strings.HasPrefix(endpoint, ":"):
		return "localhost" + endpoint, nil
	default:
		// TODO: Define an optional interface on the resolver builder to return
		// the channel authority given the user's dial target. For resolvers
		// which don't implement this interface, we will use the endpoint from
		// "scheme://authority/endpoint" as the default authority.
		return endpoint, nil
	***REMOVED***
***REMOVED***
