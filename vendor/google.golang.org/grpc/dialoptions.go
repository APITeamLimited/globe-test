/*
 *
 * Copyright 2018 gRPC authors.
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
	"net"
	"time"

	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/channelz"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/internal"
	internalbackoff "google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
)

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct ***REMOVED***
	unaryInt  UnaryClientInterceptor
	streamInt StreamClientInterceptor

	chainUnaryInts  []UnaryClientInterceptor
	chainStreamInts []StreamClientInterceptor

	cp                          Compressor
	dc                          Decompressor
	bs                          internalbackoff.Strategy
	block                       bool
	returnLastError             bool
	timeout                     time.Duration
	scChan                      <-chan ServiceConfig
	authority                   string
	copts                       transport.ConnectOptions
	callOptions                 []CallOption
	channelzParentID            *channelz.Identifier
	disableServiceConfig        bool
	disableRetry                bool
	disableHealthCheck          bool
	healthCheckFunc             internal.HealthChecker
	minConnectTimeout           func() time.Duration
	defaultServiceConfig        *ServiceConfig // defaultServiceConfig is parsed from defaultServiceConfigRawJSON.
	defaultServiceConfigRawJSON *string
	resolvers                   []resolver.Builder
***REMOVED***

// DialOption configures how we set up the connection.
type DialOption interface ***REMOVED***
	apply(*dialOptions)
***REMOVED***

// EmptyDialOption does not alter the dial configuration. It can be embedded in
// another structure to build custom dial options.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type EmptyDialOption struct***REMOVED******REMOVED***

func (EmptyDialOption) apply(*dialOptions) ***REMOVED******REMOVED***

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcDialOption struct ***REMOVED***
	f func(*dialOptions)
***REMOVED***

func (fdo *funcDialOption) apply(do *dialOptions) ***REMOVED***
	fdo.f(do)
***REMOVED***

func newFuncDialOption(f func(*dialOptions)) *funcDialOption ***REMOVED***
	return &funcDialOption***REMOVED***
		f: f,
	***REMOVED***
***REMOVED***

// WithWriteBufferSize determines how much data can be batched before doing a
// write on the wire. The corresponding memory allocation for this buffer will
// be twice the size to keep syscalls low. The default value for this buffer is
// 32KB.
//
// Zero will disable the write buffer such that each write will be on underlying
// connection. Note: A Send call may not directly translate to a write.
func WithWriteBufferSize(s int) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.WriteBufferSize = s
	***REMOVED***)
***REMOVED***

// WithReadBufferSize lets you set the size of read buffer, this determines how
// much data can be read at most for each read syscall.
//
// The default value for this buffer is 32KB. Zero will disable read buffer for
// a connection so data framer can access the underlying conn directly.
func WithReadBufferSize(s int) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.ReadBufferSize = s
	***REMOVED***)
***REMOVED***

// WithInitialWindowSize returns a DialOption which sets the value for initial
// window size on a stream. The lower bound for window size is 64K and any value
// smaller than that will be ignored.
func WithInitialWindowSize(s int32) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.InitialWindowSize = s
	***REMOVED***)
***REMOVED***

// WithInitialConnWindowSize returns a DialOption which sets the value for
// initial window size on a connection. The lower bound for window size is 64K
// and any value smaller than that will be ignored.
func WithInitialConnWindowSize(s int32) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.InitialConnWindowSize = s
	***REMOVED***)
***REMOVED***

// WithMaxMsgSize returns a DialOption which sets the maximum message size the
// client can receive.
//
// Deprecated: use WithDefaultCallOptions(MaxCallRecvMsgSize(s)) instead.  Will
// be supported throughout 1.x.
func WithMaxMsgSize(s int) DialOption ***REMOVED***
	return WithDefaultCallOptions(MaxCallRecvMsgSize(s))
***REMOVED***

// WithDefaultCallOptions returns a DialOption which sets the default
// CallOptions for calls over the connection.
func WithDefaultCallOptions(cos ...CallOption) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.callOptions = append(o.callOptions, cos...)
	***REMOVED***)
***REMOVED***

// WithCodec returns a DialOption which sets a codec for message marshaling and
// unmarshaling.
//
// Deprecated: use WithDefaultCallOptions(ForceCodec(_)) instead.  Will be
// supported throughout 1.x.
func WithCodec(c Codec) DialOption ***REMOVED***
	return WithDefaultCallOptions(CallCustomCodec(c))
***REMOVED***

// WithCompressor returns a DialOption which sets a Compressor to use for
// message compression. It has lower priority than the compressor set by the
// UseCompressor CallOption.
//
// Deprecated: use UseCompressor instead.  Will be supported throughout 1.x.
func WithCompressor(cp Compressor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.cp = cp
	***REMOVED***)
***REMOVED***

// WithDecompressor returns a DialOption which sets a Decompressor to use for
// incoming message decompression.  If incoming response messages are encoded
// using the decompressor's Type(), it will be used.  Otherwise, the message
// encoding will be used to look up the compressor registered via
// encoding.RegisterCompressor, which will then be used to decompress the
// message.  If no compressor is registered for the encoding, an Unimplemented
// status error will be returned.
//
// Deprecated: use encoding.RegisterCompressor instead.  Will be supported
// throughout 1.x.
func WithDecompressor(dc Decompressor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.dc = dc
	***REMOVED***)
***REMOVED***

// WithServiceConfig returns a DialOption which has a channel to read the
// service configuration.
//
// Deprecated: service config should be received through name resolver or via
// WithDefaultServiceConfig, as specified at
// https://github.com/grpc/grpc/blob/master/doc/service_config.md.  Will be
// removed in a future 1.x release.
func WithServiceConfig(c <-chan ServiceConfig) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.scChan = c
	***REMOVED***)
***REMOVED***

// WithConnectParams configures the ClientConn to use the provided ConnectParams
// for creating and maintaining connections to servers.
//
// The backoff configuration specified as part of the ConnectParams overrides
// all defaults specified in
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md. Consider
// using the backoff.DefaultConfig as a base, in cases where you want to
// override only a subset of the backoff configuration.
func WithConnectParams(p ConnectParams) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.bs = internalbackoff.Exponential***REMOVED***Config: p.Backoff***REMOVED***
		o.minConnectTimeout = func() time.Duration ***REMOVED***
			return p.MinConnectTimeout
		***REMOVED***
	***REMOVED***)
***REMOVED***

// WithBackoffMaxDelay configures the dialer to use the provided maximum delay
// when backing off after failed connection attempts.
//
// Deprecated: use WithConnectParams instead. Will be supported throughout 1.x.
func WithBackoffMaxDelay(md time.Duration) DialOption ***REMOVED***
	return WithBackoffConfig(BackoffConfig***REMOVED***MaxDelay: md***REMOVED***)
***REMOVED***

// WithBackoffConfig configures the dialer to use the provided backoff
// parameters after connection failures.
//
// Deprecated: use WithConnectParams instead. Will be supported throughout 1.x.
func WithBackoffConfig(b BackoffConfig) DialOption ***REMOVED***
	bc := backoff.DefaultConfig
	bc.MaxDelay = b.MaxDelay
	return withBackoff(internalbackoff.Exponential***REMOVED***Config: bc***REMOVED***)
***REMOVED***

// withBackoff sets the backoff strategy used for connectRetryNum after a failed
// connection attempt.
//
// This can be exported if arbitrary backoff strategies are allowed by gRPC.
func withBackoff(bs internalbackoff.Strategy) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.bs = bs
	***REMOVED***)
***REMOVED***

// WithBlock returns a DialOption which makes callers of Dial block until the
// underlying connection is up. Without this, Dial returns immediately and
// connecting the server happens in background.
func WithBlock() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.block = true
	***REMOVED***)
***REMOVED***

// WithReturnConnectionError returns a DialOption which makes the client connection
// return a string containing both the last connection error that occurred and
// the context.DeadlineExceeded error.
// Implies WithBlock()
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithReturnConnectionError() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.block = true
		o.returnLastError = true
	***REMOVED***)
***REMOVED***

// WithInsecure returns a DialOption which disables transport security for this
// ClientConn. Under the hood, it uses insecure.NewCredentials().
//
// Note that using this DialOption with per-RPC credentials (through
// WithCredentialsBundle or WithPerRPCCredentials) which require transport
// security is incompatible and will cause grpc.Dial() to fail.
//
// Deprecated: use WithTransportCredentials and insecure.NewCredentials()
// instead. Will be supported throughout 1.x.
func WithInsecure() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.TransportCredentials = insecure.NewCredentials()
	***REMOVED***)
***REMOVED***

// WithNoProxy returns a DialOption which disables the use of proxies for this
// ClientConn. This is ignored if WithDialer or WithContextDialer are used.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithNoProxy() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.UseProxy = false
	***REMOVED***)
***REMOVED***

// WithTransportCredentials returns a DialOption which configures a connection
// level security credentials (e.g., TLS/SSL). This should not be used together
// with WithCredentialsBundle.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.TransportCredentials = creds
	***REMOVED***)
***REMOVED***

// WithPerRPCCredentials returns a DialOption which sets credentials and places
// auth state on each outbound RPC.
func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
	***REMOVED***)
***REMOVED***

// WithCredentialsBundle returns a DialOption to set a credentials bundle for
// the ClientConn.WithCreds. This should not be used together with
// WithTransportCredentials.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithCredentialsBundle(b credentials.Bundle) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.CredsBundle = b
	***REMOVED***)
***REMOVED***

// WithTimeout returns a DialOption that configures a timeout for dialing a
// ClientConn initially. This is valid if and only if WithBlock() is present.
//
// Deprecated: use DialContext instead of Dial and context.WithTimeout
// instead.  Will be supported throughout 1.x.
func WithTimeout(d time.Duration) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.timeout = d
	***REMOVED***)
***REMOVED***

// WithContextDialer returns a DialOption that sets a dialer to create
// connections. If FailOnNonTempDialError() is set to true, and an error is
// returned by f, gRPC checks the error's Temporary() method to decide if it
// should try to reconnect to the network address.
func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.Dialer = f
	***REMOVED***)
***REMOVED***

func init() ***REMOVED***
	internal.WithHealthCheckFunc = withHealthCheckFunc
***REMOVED***

// WithDialer returns a DialOption that specifies a function to use for dialing
// network addresses. If FailOnNonTempDialError() is set to true, and an error
// is returned by f, gRPC checks the error's Temporary() method to decide if it
// should try to reconnect to the network address.
//
// Deprecated: use WithContextDialer instead.  Will be supported throughout
// 1.x.
func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption ***REMOVED***
	return WithContextDialer(
		func(ctx context.Context, addr string) (net.Conn, error) ***REMOVED***
			if deadline, ok := ctx.Deadline(); ok ***REMOVED***
				return f(addr, time.Until(deadline))
			***REMOVED***
			return f(addr, 0)
		***REMOVED***)
***REMOVED***

// WithStatsHandler returns a DialOption that specifies the stats handler for
// all the RPCs and underlying network connections in this ClientConn.
func WithStatsHandler(h stats.Handler) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.StatsHandler = h
	***REMOVED***)
***REMOVED***

// FailOnNonTempDialError returns a DialOption that specifies if gRPC fails on
// non-temporary dial errors. If f is true, and dialer returns a non-temporary
// error, gRPC will fail the connection to the network address and won't try to
// reconnect. The default value of FailOnNonTempDialError is false.
//
// FailOnNonTempDialError only affects the initial dial, and does not do
// anything useful unless you are also using WithBlock().
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func FailOnNonTempDialError(f bool) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.FailOnNonTempDialError = f
	***REMOVED***)
***REMOVED***

// WithUserAgent returns a DialOption that specifies a user agent string for all
// the RPCs.
func WithUserAgent(s string) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.UserAgent = s
	***REMOVED***)
***REMOVED***

// WithKeepaliveParams returns a DialOption that specifies keepalive parameters
// for the client transport.
func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption ***REMOVED***
	if kp.Time < internal.KeepaliveMinPingTime ***REMOVED***
		logger.Warningf("Adjusting keepalive ping interval to minimum period of %v", internal.KeepaliveMinPingTime)
		kp.Time = internal.KeepaliveMinPingTime
	***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.KeepaliveParams = kp
	***REMOVED***)
***REMOVED***

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for
// unary RPCs.
func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.unaryInt = f
	***REMOVED***)
***REMOVED***

// WithChainUnaryInterceptor returns a DialOption that specifies the chained
// interceptor for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithUnaryInterceptor will always be prepended to the chain.
func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.chainUnaryInts = append(o.chainUnaryInts, interceptors...)
	***REMOVED***)
***REMOVED***

// WithStreamInterceptor returns a DialOption that specifies the interceptor for
// streaming RPCs.
func WithStreamInterceptor(f StreamClientInterceptor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.streamInt = f
	***REMOVED***)
***REMOVED***

// WithChainStreamInterceptor returns a DialOption that specifies the chained
// interceptor for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithStreamInterceptor will always be prepended to the chain.
func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.chainStreamInts = append(o.chainStreamInts, interceptors...)
	***REMOVED***)
***REMOVED***

// WithAuthority returns a DialOption that specifies the value to be used as the
// :authority pseudo-header and as the server name in authentication handshake.
func WithAuthority(a string) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.authority = a
	***REMOVED***)
***REMOVED***

// WithChannelzParentID returns a DialOption that specifies the channelz ID of
// current ClientConn's parent. This function is used in nested channel creation
// (e.g. grpclb dial).
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithChannelzParentID(id *channelz.Identifier) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.channelzParentID = id
	***REMOVED***)
***REMOVED***

// WithDisableServiceConfig returns a DialOption that causes gRPC to ignore any
// service config provided by the resolver and provides a hint to the resolver
// to not fetch service configs.
//
// Note that this dial option only disables service config from resolver. If
// default service config is provided, gRPC will use the default service config.
func WithDisableServiceConfig() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.disableServiceConfig = true
	***REMOVED***)
***REMOVED***

// WithDefaultServiceConfig returns a DialOption that configures the default
// service config, which will be used in cases where:
//
// 1. WithDisableServiceConfig is also used, or
//
// 2. The name resolver does not provide a service config or provides an
// invalid service config.
//
// The parameter s is the JSON representation of the default service config.
// For more information about service configs, see:
// https://github.com/grpc/grpc/blob/master/doc/service_config.md
// For a simple example of usage, see:
// examples/features/load_balancing/client/main.go
func WithDefaultServiceConfig(s string) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.defaultServiceConfigRawJSON = &s
	***REMOVED***)
***REMOVED***

// WithDisableRetry returns a DialOption that disables retries, even if the
// service config enables them.  This does not impact transparent retries, which
// will happen automatically if no data is written to the wire or if the RPC is
// unprocessed by the remote server.
//
// Retry support is currently enabled by default, but may be disabled by
// setting the environment variable "GRPC_GO_RETRY" to "off".
func WithDisableRetry() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.disableRetry = true
	***REMOVED***)
***REMOVED***

// WithMaxHeaderListSize returns a DialOption that specifies the maximum
// (uncompressed) size of header list that the client is prepared to accept.
func WithMaxHeaderListSize(s uint32) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.copts.MaxHeaderListSize = &s
	***REMOVED***)
***REMOVED***

// WithDisableHealthCheck disables the LB channel health checking for all
// SubConns of this ClientConn.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithDisableHealthCheck() DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.disableHealthCheck = true
	***REMOVED***)
***REMOVED***

// withHealthCheckFunc replaces the default health check function with the
// provided one. It makes tests easier to change the health check function.
//
// For testing purpose only.
func withHealthCheckFunc(f internal.HealthChecker) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.healthCheckFunc = f
	***REMOVED***)
***REMOVED***

func defaultDialOptions() dialOptions ***REMOVED***
	return dialOptions***REMOVED***
		healthCheckFunc: internal.HealthCheckFunc,
		copts: transport.ConnectOptions***REMOVED***
			WriteBufferSize: defaultWriteBufSize,
			ReadBufferSize:  defaultReadBufSize,
			UseProxy:        true,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// withGetMinConnectDeadline specifies the function that clientconn uses to
// get minConnectDeadline. This can be used to make connection attempts happen
// faster/slower.
//
// For testing purpose only.
func withMinConnectDeadline(f func() time.Duration) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.minConnectTimeout = f
	***REMOVED***)
***REMOVED***

// WithResolvers allows a list of resolver implementations to be registered
// locally with the ClientConn without needing to be globally registered via
// resolver.Register.  They will be matched against the scheme used for the
// current Dial only, and will take precedence over the global registry.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func WithResolvers(rs ...resolver.Builder) DialOption ***REMOVED***
	return newFuncDialOption(func(o *dialOptions) ***REMOVED***
		o.resolvers = append(o.resolvers, rs...)
	***REMOVED***)
***REMOVED***
