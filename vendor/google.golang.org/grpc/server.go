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
	"io"
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/trace"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/binarylog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

const (
	defaultServerMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultServerMaxSendMessageSize    = math.MaxInt32

	// Server transports are tracked in a map which is keyed on listener
	// address. For regular gRPC traffic, connections are accepted in Serve()
	// through a call to Accept(), and we use the actual listener address as key
	// when we add it to the map. But for connections received through
	// ServeHTTP(), we do not have a listener and hence use this dummy value.
	listenerAddressForServeHTTP = "listenerAddressForServeHTTP"
)

func init() ***REMOVED***
	internal.GetServerCredentials = func(srv *Server) credentials.TransportCredentials ***REMOVED***
		return srv.opts.creds
	***REMOVED***
	internal.DrainServerTransports = func(srv *Server, addr string) ***REMOVED***
		srv.drainServerTransports(addr)
	***REMOVED***
***REMOVED***

var statusOK = status.New(codes.OK, "")
var logger = grpclog.Component("core")

type methodHandler func(srv interface***REMOVED******REMOVED***, ctx context.Context, dec func(interface***REMOVED******REMOVED***) error, interceptor UnaryServerInterceptor) (interface***REMOVED******REMOVED***, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct ***REMOVED***
	MethodName string
	Handler    methodHandler
***REMOVED***

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct ***REMOVED***
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface***REMOVED******REMOVED***
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    interface***REMOVED******REMOVED***
***REMOVED***

// serviceInfo wraps information about a service. It is very similar to
// ServiceDesc and is constructed from it for internal purposes.
type serviceInfo struct ***REMOVED***
	// Contains the implementation for the methods in this service.
	serviceImpl interface***REMOVED******REMOVED***
	methods     map[string]*MethodDesc
	streams     map[string]*StreamDesc
	mdata       interface***REMOVED******REMOVED***
***REMOVED***

type serverWorkerData struct ***REMOVED***
	st     transport.ServerTransport
	wg     *sync.WaitGroup
	stream *transport.Stream
***REMOVED***

// Server is a gRPC server to serve RPC requests.
type Server struct ***REMOVED***
	opts serverOptions

	mu  sync.Mutex // guards following
	lis map[net.Listener]bool
	// conns contains all active server transports. It is a map keyed on a
	// listener address with the value being the set of active transports
	// belonging to that listener.
	conns    map[string]map[transport.ServerTransport]bool
	serve    bool
	drain    bool
	cv       *sync.Cond              // signaled when connections close for GracefulStop
	services map[string]*serviceInfo // service name -> service info
	events   trace.EventLog

	quit               *grpcsync.Event
	done               *grpcsync.Event
	channelzRemoveOnce sync.Once
	serveWG            sync.WaitGroup // counts active Serve goroutines for GracefulStop

	channelzID int64 // channelz unique identification number
	czData     *channelzData

	serverWorkerChannels []chan *serverWorkerData
***REMOVED***

type serverOptions struct ***REMOVED***
	creds                 credentials.TransportCredentials
	codec                 baseCodec
	cp                    Compressor
	dc                    Decompressor
	unaryInt              UnaryServerInterceptor
	streamInt             StreamServerInterceptor
	chainUnaryInts        []UnaryServerInterceptor
	chainStreamInts       []StreamServerInterceptor
	inTapHandle           tap.ServerInHandle
	statsHandler          stats.Handler
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	unknownStreamDesc     *StreamDesc
	keepaliveParams       keepalive.ServerParameters
	keepalivePolicy       keepalive.EnforcementPolicy
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	connectionTimeout     time.Duration
	maxHeaderListSize     *uint32
	headerTableSize       *uint32
	numServerWorkers      uint32
***REMOVED***

var defaultServerOptions = serverOptions***REMOVED***
	maxReceiveMessageSize: defaultServerMaxReceiveMessageSize,
	maxSendMessageSize:    defaultServerMaxSendMessageSize,
	connectionTimeout:     120 * time.Second,
	writeBufferSize:       defaultWriteBufSize,
	readBufferSize:        defaultReadBufSize,
***REMOVED***

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type ServerOption interface ***REMOVED***
	apply(*serverOptions)
***REMOVED***

// EmptyServerOption does not alter the server configuration. It can be embedded
// in another structure to build custom server options.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type EmptyServerOption struct***REMOVED******REMOVED***

func (EmptyServerOption) apply(*serverOptions) ***REMOVED******REMOVED***

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type funcServerOption struct ***REMOVED***
	f func(*serverOptions)
***REMOVED***

func (fdo *funcServerOption) apply(do *serverOptions) ***REMOVED***
	fdo.f(do)
***REMOVED***

func newFuncServerOption(f func(*serverOptions)) *funcServerOption ***REMOVED***
	return &funcServerOption***REMOVED***
		f: f,
	***REMOVED***
***REMOVED***

// WriteBufferSize determines how much data can be batched before doing a write on the wire.
// The corresponding memory allocation for this buffer will be twice the size to keep syscalls low.
// The default value for this buffer is 32KB.
// Zero will disable the write buffer such that each write will be on underlying connection.
// Note: A Send call may not directly translate to a write.
func WriteBufferSize(s int) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.writeBufferSize = s
	***REMOVED***)
***REMOVED***

// ReadBufferSize lets you set the size of read buffer, this determines how much data can be read at most
// for one read syscall.
// The default value for this buffer is 32KB.
// Zero will disable read buffer for a connection so data framer can access the underlying
// conn directly.
func ReadBufferSize(s int) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.readBufferSize = s
	***REMOVED***)
***REMOVED***

// InitialWindowSize returns a ServerOption that sets window size for stream.
// The lower bound for window size is 64K and any value smaller than that will be ignored.
func InitialWindowSize(s int32) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.initialWindowSize = s
	***REMOVED***)
***REMOVED***

// InitialConnWindowSize returns a ServerOption that sets window size for a connection.
// The lower bound for window size is 64K and any value smaller than that will be ignored.
func InitialConnWindowSize(s int32) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.initialConnWindowSize = s
	***REMOVED***)
***REMOVED***

// KeepaliveParams returns a ServerOption that sets keepalive and max-age parameters for the server.
func KeepaliveParams(kp keepalive.ServerParameters) ServerOption ***REMOVED***
	if kp.Time > 0 && kp.Time < time.Second ***REMOVED***
		logger.Warning("Adjusting keepalive ping interval to minimum period of 1s")
		kp.Time = time.Second
	***REMOVED***

	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.keepaliveParams = kp
	***REMOVED***)
***REMOVED***

// KeepaliveEnforcementPolicy returns a ServerOption that sets keepalive enforcement policy for the server.
func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.keepalivePolicy = kep
	***REMOVED***)
***REMOVED***

// CustomCodec returns a ServerOption that sets a codec for message marshaling and unmarshaling.
//
// This will override any lookups by content-subtype for Codecs registered with RegisterCodec.
//
// Deprecated: register codecs using encoding.RegisterCodec. The server will
// automatically use registered codecs based on the incoming requests' headers.
// See also
// https://github.com/grpc/grpc-go/blob/master/Documentation/encoding.md#using-a-codec.
// Will be supported throughout 1.x.
func CustomCodec(codec Codec) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.codec = codec
	***REMOVED***)
***REMOVED***

// ForceServerCodec returns a ServerOption that sets a codec for message
// marshaling and unmarshaling.
//
// This will override any lookups by content-subtype for Codecs registered
// with RegisterCodec.
//
// See Content-Type on
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details. Also see the documentation on RegisterCodec and
// CallContentSubtype for more details on the interaction between encoding.Codec
// and content-subtype.
//
// This function is provided for advanced users; prefer to register codecs
// using encoding.RegisterCodec.
// The server will automatically use registered codecs based on the incoming
// requests' headers. See also
// https://github.com/grpc/grpc-go/blob/master/Documentation/encoding.md#using-a-codec.
// Will be supported throughout 1.x.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func ForceServerCodec(codec encoding.Codec) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.codec = codec
	***REMOVED***)
***REMOVED***

// RPCCompressor returns a ServerOption that sets a compressor for outbound
// messages.  For backward compatibility, all outbound messages will be sent
// using this compressor, regardless of incoming message compression.  By
// default, server messages will be sent using the same compressor with which
// request messages were sent.
//
// Deprecated: use encoding.RegisterCompressor instead. Will be supported
// throughout 1.x.
func RPCCompressor(cp Compressor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.cp = cp
	***REMOVED***)
***REMOVED***

// RPCDecompressor returns a ServerOption that sets a decompressor for inbound
// messages.  It has higher priority than decompressors registered via
// encoding.RegisterCompressor.
//
// Deprecated: use encoding.RegisterCompressor instead. Will be supported
// throughout 1.x.
func RPCDecompressor(dc Decompressor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.dc = dc
	***REMOVED***)
***REMOVED***

// MaxMsgSize returns a ServerOption to set the max message size in bytes the server can receive.
// If this is not set, gRPC uses the default limit.
//
// Deprecated: use MaxRecvMsgSize instead. Will be supported throughout 1.x.
func MaxMsgSize(m int) ServerOption ***REMOVED***
	return MaxRecvMsgSize(m)
***REMOVED***

// MaxRecvMsgSize returns a ServerOption to set the max message size in bytes the server can receive.
// If this is not set, gRPC uses the default 4MB.
func MaxRecvMsgSize(m int) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.maxReceiveMessageSize = m
	***REMOVED***)
***REMOVED***

// MaxSendMsgSize returns a ServerOption to set the max message size in bytes the server can send.
// If this is not set, gRPC uses the default `math.MaxInt32`.
func MaxSendMsgSize(m int) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.maxSendMessageSize = m
	***REMOVED***)
***REMOVED***

// MaxConcurrentStreams returns a ServerOption that will apply a limit on the number
// of concurrent streams to each ServerTransport.
func MaxConcurrentStreams(n uint32) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.maxConcurrentStreams = n
	***REMOVED***)
***REMOVED***

// Creds returns a ServerOption that sets credentials for server connections.
func Creds(c credentials.TransportCredentials) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.creds = c
	***REMOVED***)
***REMOVED***

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the
// server. Only one unary interceptor can be installed. The construction of multiple
// interceptors (e.g., chaining) can be implemented at the caller.
func UnaryInterceptor(i UnaryServerInterceptor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		if o.unaryInt != nil ***REMOVED***
			panic("The unary server interceptor was already set and may not be reset.")
		***REMOVED***
		o.unaryInt = i
	***REMOVED***)
***REMOVED***

// ChainUnaryInterceptor returns a ServerOption that specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All unary interceptors added by this method will be chained.
func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.chainUnaryInts = append(o.chainUnaryInts, interceptors...)
	***REMOVED***)
***REMOVED***

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the
// server. Only one stream interceptor can be installed.
func StreamInterceptor(i StreamServerInterceptor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		if o.streamInt != nil ***REMOVED***
			panic("The stream server interceptor was already set and may not be reset.")
		***REMOVED***
		o.streamInt = i
	***REMOVED***)
***REMOVED***

// ChainStreamInterceptor returns a ServerOption that specifies the chained interceptor
// for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All stream interceptors added by this method will be chained.
func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.chainStreamInts = append(o.chainStreamInts, interceptors...)
	***REMOVED***)
***REMOVED***

// InTapHandle returns a ServerOption that sets the tap handle for all the server
// transport to be created. Only one can be installed.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func InTapHandle(h tap.ServerInHandle) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		if o.inTapHandle != nil ***REMOVED***
			panic("The tap handle was already set and may not be reset.")
		***REMOVED***
		o.inTapHandle = h
	***REMOVED***)
***REMOVED***

// StatsHandler returns a ServerOption that sets the stats handler for the server.
func StatsHandler(h stats.Handler) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.statsHandler = h
	***REMOVED***)
***REMOVED***

// UnknownServiceHandler returns a ServerOption that allows for adding a custom
// unknown service handler. The provided method is a bidi-streaming RPC service
// handler that will be invoked instead of returning the "unimplemented" gRPC
// error whenever a request is received for an unregistered service or method.
// The handling function and stream interceptor (if set) have full access to
// the ServerStream, including its Context.
func UnknownServiceHandler(streamHandler StreamHandler) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.unknownStreamDesc = &StreamDesc***REMOVED***
			StreamName: "unknown_service_handler",
			Handler:    streamHandler,
			// We need to assume that the users of the streamHandler will want to use both.
			ClientStreams: true,
			ServerStreams: true,
		***REMOVED***
	***REMOVED***)
***REMOVED***

// ConnectionTimeout returns a ServerOption that sets the timeout for
// connection establishment (up to and including HTTP/2 handshaking) for all
// new connections.  If this is not set, the default is 120 seconds.  A zero or
// negative value will result in an immediate timeout.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func ConnectionTimeout(d time.Duration) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.connectionTimeout = d
	***REMOVED***)
***REMOVED***

// MaxHeaderListSize returns a ServerOption that sets the max (uncompressed) size
// of header list that the server is prepared to accept.
func MaxHeaderListSize(s uint32) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.maxHeaderListSize = &s
	***REMOVED***)
***REMOVED***

// HeaderTableSize returns a ServerOption that sets the size of dynamic
// header table for stream.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func HeaderTableSize(s uint32) ServerOption ***REMOVED***
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.headerTableSize = &s
	***REMOVED***)
***REMOVED***

// NumStreamWorkers returns a ServerOption that sets the number of worker
// goroutines that should be used to process incoming streams. Setting this to
// zero (default) will disable workers and spawn a new goroutine for each
// stream.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func NumStreamWorkers(numServerWorkers uint32) ServerOption ***REMOVED***
	// TODO: If/when this API gets stabilized (i.e. stream workers become the
	// only way streams are processed), change the behavior of the zero value to
	// a sane default. Preliminary experiments suggest that a value equal to the
	// number of CPUs available is most performant; requires thorough testing.
	return newFuncServerOption(func(o *serverOptions) ***REMOVED***
		o.numServerWorkers = numServerWorkers
	***REMOVED***)
***REMOVED***

// serverWorkerResetThreshold defines how often the stack must be reset. Every
// N requests, by spawning a new goroutine in its place, a worker can reset its
// stack so that large stacks don't live in memory forever. 2^16 should allow
// each goroutine stack to live for at least a few seconds in a typical
// workload (assuming a QPS of a few thousand requests/sec).
const serverWorkerResetThreshold = 1 << 16

// serverWorkers blocks on a *transport.Stream channel forever and waits for
// data to be fed by serveStreams. This allows different requests to be
// processed by the same goroutine, removing the need for expensive stack
// re-allocations (see the runtime.morestack problem [1]).
//
// [1] https://github.com/golang/go/issues/18138
func (s *Server) serverWorker(ch chan *serverWorkerData) ***REMOVED***
	// To make sure all server workers don't reset at the same time, choose a
	// random number of iterations before resetting.
	threshold := serverWorkerResetThreshold + grpcrand.Intn(serverWorkerResetThreshold)
	for completed := 0; completed < threshold; completed++ ***REMOVED***
		data, ok := <-ch
		if !ok ***REMOVED***
			return
		***REMOVED***
		s.handleStream(data.st, data.stream, s.traceInfo(data.st, data.stream))
		data.wg.Done()
	***REMOVED***
	go s.serverWorker(ch)
***REMOVED***

// initServerWorkers creates worker goroutines and channels to process incoming
// connections to reduce the time spent overall on runtime.morestack.
func (s *Server) initServerWorkers() ***REMOVED***
	s.serverWorkerChannels = make([]chan *serverWorkerData, s.opts.numServerWorkers)
	for i := uint32(0); i < s.opts.numServerWorkers; i++ ***REMOVED***
		s.serverWorkerChannels[i] = make(chan *serverWorkerData)
		go s.serverWorker(s.serverWorkerChannels[i])
	***REMOVED***
***REMOVED***

func (s *Server) stopServerWorkers() ***REMOVED***
	for i := uint32(0); i < s.opts.numServerWorkers; i++ ***REMOVED***
		close(s.serverWorkerChannels[i])
	***REMOVED***
***REMOVED***

// NewServer creates a gRPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(opt ...ServerOption) *Server ***REMOVED***
	opts := defaultServerOptions
	for _, o := range opt ***REMOVED***
		o.apply(&opts)
	***REMOVED***
	s := &Server***REMOVED***
		lis:      make(map[net.Listener]bool),
		opts:     opts,
		conns:    make(map[string]map[transport.ServerTransport]bool),
		services: make(map[string]*serviceInfo),
		quit:     grpcsync.NewEvent(),
		done:     grpcsync.NewEvent(),
		czData:   new(channelzData),
	***REMOVED***
	chainUnaryServerInterceptors(s)
	chainStreamServerInterceptors(s)
	s.cv = sync.NewCond(&s.mu)
	if EnableTracing ***REMOVED***
		_, file, line, _ := runtime.Caller(1)
		s.events = trace.NewEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
	***REMOVED***

	if s.opts.numServerWorkers > 0 ***REMOVED***
		s.initServerWorkers()
	***REMOVED***

	if channelz.IsOn() ***REMOVED***
		s.channelzID = channelz.RegisterServer(&channelzServer***REMOVED***s***REMOVED***, "")
	***REMOVED***
	return s
***REMOVED***

// printf records an event in s's event log, unless s has been stopped.
// REQUIRES s.mu is held.
func (s *Server) printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if s.events != nil ***REMOVED***
		s.events.Printf(format, a...)
	***REMOVED***
***REMOVED***

// errorf records an error in s's event log, unless s has been stopped.
// REQUIRES s.mu is held.
func (s *Server) errorf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if s.events != nil ***REMOVED***
		s.events.Errorf(format, a...)
	***REMOVED***
***REMOVED***

// ServiceRegistrar wraps a single method that supports service registration. It
// enables users to pass concrete types other than grpc.Server to the service
// registration methods exported by the IDL generated code.
type ServiceRegistrar interface ***REMOVED***
	// RegisterService registers a service and its implementation to the
	// concrete type implementing this interface.  It may not be called
	// once the server has started serving.
	// desc describes the service and its methods and handlers. impl is the
	// service implementation which is passed to the method handlers.
	RegisterService(desc *ServiceDesc, impl interface***REMOVED******REMOVED***)
***REMOVED***

// RegisterService registers a service and its implementation to the gRPC
// server. It is called from the IDL generated code. This must be called before
// invoking Serve. If ss is non-nil (for legacy code), its type is checked to
// ensure it implements sd.HandlerType.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface***REMOVED******REMOVED***) ***REMOVED***
	if ss != nil ***REMOVED***
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) ***REMOVED***
			logger.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		***REMOVED***
	***REMOVED***
	s.register(sd, ss)
***REMOVED***

func (s *Server) register(sd *ServiceDesc, ss interface***REMOVED******REMOVED***) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	s.printf("RegisterService(%q)", sd.ServiceName)
	if s.serve ***REMOVED***
		logger.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
	***REMOVED***
	if _, ok := s.services[sd.ServiceName]; ok ***REMOVED***
		logger.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	***REMOVED***
	info := &serviceInfo***REMOVED***
		serviceImpl: ss,
		methods:     make(map[string]*MethodDesc),
		streams:     make(map[string]*StreamDesc),
		mdata:       sd.Metadata,
	***REMOVED***
	for i := range sd.Methods ***REMOVED***
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	***REMOVED***
	for i := range sd.Streams ***REMOVED***
		d := &sd.Streams[i]
		info.streams[d.StreamName] = d
	***REMOVED***
	s.services[sd.ServiceName] = info
***REMOVED***

// MethodInfo contains the information of an RPC including its method name and type.
type MethodInfo struct ***REMOVED***
	// Name is the method name only, without the service name or package name.
	Name string
	// IsClientStream indicates whether the RPC is a client streaming RPC.
	IsClientStream bool
	// IsServerStream indicates whether the RPC is a server streaming RPC.
	IsServerStream bool
***REMOVED***

// ServiceInfo contains unary RPC method info, streaming RPC method info and metadata for a service.
type ServiceInfo struct ***REMOVED***
	Methods []MethodInfo
	// Metadata is the metadata specified in ServiceDesc when registering service.
	Metadata interface***REMOVED******REMOVED***
***REMOVED***

// GetServiceInfo returns a map from service names to ServiceInfo.
// Service names include the package names, in the form of <package>.<service>.
func (s *Server) GetServiceInfo() map[string]ServiceInfo ***REMOVED***
	ret := make(map[string]ServiceInfo)
	for n, srv := range s.services ***REMOVED***
		methods := make([]MethodInfo, 0, len(srv.methods)+len(srv.streams))
		for m := range srv.methods ***REMOVED***
			methods = append(methods, MethodInfo***REMOVED***
				Name:           m,
				IsClientStream: false,
				IsServerStream: false,
			***REMOVED***)
		***REMOVED***
		for m, d := range srv.streams ***REMOVED***
			methods = append(methods, MethodInfo***REMOVED***
				Name:           m,
				IsClientStream: d.ClientStreams,
				IsServerStream: d.ServerStreams,
			***REMOVED***)
		***REMOVED***

		ret[n] = ServiceInfo***REMOVED***
			Methods:  methods,
			Metadata: srv.mdata,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

// ErrServerStopped indicates that the operation is now illegal because of
// the server being stopped.
var ErrServerStopped = errors.New("grpc: the server has been stopped")

type listenSocket struct ***REMOVED***
	net.Listener
	channelzID int64
***REMOVED***

func (l *listenSocket) ChannelzMetric() *channelz.SocketInternalMetric ***REMOVED***
	return &channelz.SocketInternalMetric***REMOVED***
		SocketOptions: channelz.GetSocketOption(l.Listener),
		LocalAddr:     l.Listener.Addr(),
	***REMOVED***
***REMOVED***

func (l *listenSocket) Close() error ***REMOVED***
	err := l.Listener.Close()
	if channelz.IsOn() ***REMOVED***
		channelz.RemoveEntry(l.channelzID)
	***REMOVED***
	return err
***REMOVED***

// Serve accepts incoming connections on the listener lis, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when lis.Accept fails with fatal errors.  lis will be closed when
// this method returns.
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s *Server) Serve(lis net.Listener) error ***REMOVED***
	s.mu.Lock()
	s.printf("serving")
	s.serve = true
	if s.lis == nil ***REMOVED***
		// Serve called after Stop or GracefulStop.
		s.mu.Unlock()
		lis.Close()
		return ErrServerStopped
	***REMOVED***

	s.serveWG.Add(1)
	defer func() ***REMOVED***
		s.serveWG.Done()
		if s.quit.HasFired() ***REMOVED***
			// Stop or GracefulStop called; block until done and return nil.
			<-s.done.Done()
		***REMOVED***
	***REMOVED***()

	ls := &listenSocket***REMOVED***Listener: lis***REMOVED***
	s.lis[ls] = true

	if channelz.IsOn() ***REMOVED***
		ls.channelzID = channelz.RegisterListenSocket(ls, s.channelzID, lis.Addr().String())
	***REMOVED***
	s.mu.Unlock()

	defer func() ***REMOVED***
		s.mu.Lock()
		if s.lis != nil && s.lis[ls] ***REMOVED***
			ls.Close()
			delete(s.lis, ls)
		***REMOVED***
		s.mu.Unlock()
	***REMOVED***()

	var tempDelay time.Duration // how long to sleep on accept failure

	for ***REMOVED***
		rawConn, err := lis.Accept()
		if err != nil ***REMOVED***
			if ne, ok := err.(interface ***REMOVED***
				Temporary() bool
			***REMOVED***); ok && ne.Temporary() ***REMOVED***
				if tempDelay == 0 ***REMOVED***
					tempDelay = 5 * time.Millisecond
				***REMOVED*** else ***REMOVED***
					tempDelay *= 2
				***REMOVED***
				if max := 1 * time.Second; tempDelay > max ***REMOVED***
					tempDelay = max
				***REMOVED***
				s.mu.Lock()
				s.printf("Accept error: %v; retrying in %v", err, tempDelay)
				s.mu.Unlock()
				timer := time.NewTimer(tempDelay)
				select ***REMOVED***
				case <-timer.C:
				case <-s.quit.Done():
					timer.Stop()
					return nil
				***REMOVED***
				continue
			***REMOVED***
			s.mu.Lock()
			s.printf("done serving; Accept = %v", err)
			s.mu.Unlock()

			if s.quit.HasFired() ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***
		tempDelay = 0
		// Start a new goroutine to deal with rawConn so we don't stall this Accept
		// loop goroutine.
		//
		// Make sure we account for the goroutine so GracefulStop doesn't nil out
		// s.conns before this conn can be added.
		s.serveWG.Add(1)
		go func() ***REMOVED***
			s.handleRawConn(lis.Addr().String(), rawConn)
			s.serveWG.Done()
		***REMOVED***()
	***REMOVED***
***REMOVED***

// handleRawConn forks a goroutine to handle a just-accepted connection that
// has not had any I/O performed on it yet.
func (s *Server) handleRawConn(lisAddr string, rawConn net.Conn) ***REMOVED***
	if s.quit.HasFired() ***REMOVED***
		rawConn.Close()
		return
	***REMOVED***
	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))

	// Finish handshaking (HTTP2)
	st := s.newHTTP2Transport(rawConn)
	rawConn.SetDeadline(time.Time***REMOVED******REMOVED***)
	if st == nil ***REMOVED***
		return
	***REMOVED***

	if !s.addConn(lisAddr, st) ***REMOVED***
		return
	***REMOVED***
	go func() ***REMOVED***
		s.serveStreams(st)
		s.removeConn(lisAddr, st)
	***REMOVED***()
***REMOVED***

func (s *Server) drainServerTransports(addr string) ***REMOVED***
	s.mu.Lock()
	conns := s.conns[addr]
	for st := range conns ***REMOVED***
		st.Drain()
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// newHTTP2Transport sets up a http/2 transport (using the
// gRPC http2 server transport in transport/http2_server.go).
func (s *Server) newHTTP2Transport(c net.Conn) transport.ServerTransport ***REMOVED***
	config := &transport.ServerConfig***REMOVED***
		MaxStreams:            s.opts.maxConcurrentStreams,
		ConnectionTimeout:     s.opts.connectionTimeout,
		Credentials:           s.opts.creds,
		InTapHandle:           s.opts.inTapHandle,
		StatsHandler:          s.opts.statsHandler,
		KeepaliveParams:       s.opts.keepaliveParams,
		KeepalivePolicy:       s.opts.keepalivePolicy,
		InitialWindowSize:     s.opts.initialWindowSize,
		InitialConnWindowSize: s.opts.initialConnWindowSize,
		WriteBufferSize:       s.opts.writeBufferSize,
		ReadBufferSize:        s.opts.readBufferSize,
		ChannelzParentID:      s.channelzID,
		MaxHeaderListSize:     s.opts.maxHeaderListSize,
		HeaderTableSize:       s.opts.headerTableSize,
	***REMOVED***
	st, err := transport.NewServerTransport(c, config)
	if err != nil ***REMOVED***
		s.mu.Lock()
		s.errorf("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		s.mu.Unlock()
		// ErrConnDispatched means that the connection was dispatched away from
		// gRPC; those connections should be left open.
		if err != credentials.ErrConnDispatched ***REMOVED***
			// Don't log on ErrConnDispatched and io.EOF to prevent log spam.
			if err != io.EOF ***REMOVED***
				channelz.Warning(logger, s.channelzID, "grpc: Server.Serve failed to create ServerTransport: ", err)
			***REMOVED***
			c.Close()
		***REMOVED***
		return nil
	***REMOVED***

	return st
***REMOVED***

func (s *Server) serveStreams(st transport.ServerTransport) ***REMOVED***
	defer st.Close()
	var wg sync.WaitGroup

	var roundRobinCounter uint32
	st.HandleStreams(func(stream *transport.Stream) ***REMOVED***
		wg.Add(1)
		if s.opts.numServerWorkers > 0 ***REMOVED***
			data := &serverWorkerData***REMOVED***st: st, wg: &wg, stream: stream***REMOVED***
			select ***REMOVED***
			case s.serverWorkerChannels[atomic.AddUint32(&roundRobinCounter, 1)%s.opts.numServerWorkers] <- data:
			default:
				// If all stream workers are busy, fallback to the default code path.
				go func() ***REMOVED***
					s.handleStream(st, stream, s.traceInfo(st, stream))
					wg.Done()
				***REMOVED***()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			go func() ***REMOVED***
				defer wg.Done()
				s.handleStream(st, stream, s.traceInfo(st, stream))
			***REMOVED***()
		***REMOVED***
	***REMOVED***, func(ctx context.Context, method string) context.Context ***REMOVED***
		if !EnableTracing ***REMOVED***
			return ctx
		***REMOVED***
		tr := trace.New("grpc.Recv."+methodFamily(method), method)
		return trace.NewContext(ctx, tr)
	***REMOVED***)
	wg.Wait()
***REMOVED***

var _ http.Handler = (*Server)(nil)

// ServeHTTP implements the Go standard library's http.Handler
// interface by responding to the gRPC request r, by looking up
// the requested gRPC method in the gRPC server s.
//
// The provided HTTP request must have arrived on an HTTP/2
// connection. When using the Go standard library's server,
// practically this means that the Request must also have arrived
// over TLS.
//
// To share one port (such as 443 for https) between gRPC and an
// existing http.Handler, use a root http.Handler such as:
//
//   if r.ProtoMajor == 2 && strings.HasPrefix(
//   	r.Header.Get("Content-Type"), "application/grpc") ***REMOVED***
//   	grpcServer.ServeHTTP(w, r)
//   ***REMOVED*** else ***REMOVED***
//   	yourMux.ServeHTTP(w, r)
//   ***REMOVED***
//
// Note that ServeHTTP uses Go's HTTP/2 server implementation which is totally
// separate from grpc-go's HTTP/2 server. Performance and features may vary
// between the two paths. ServeHTTP does not support some gRPC features
// available through grpc-go's HTTP/2 server.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	st, err := transport.NewServerHandlerTransport(w, r, s.opts.statsHandler)
	if err != nil ***REMOVED***
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	if !s.addConn(listenerAddressForServeHTTP, st) ***REMOVED***
		return
	***REMOVED***
	defer s.removeConn(listenerAddressForServeHTTP, st)
	s.serveStreams(st)
***REMOVED***

// traceInfo returns a traceInfo and associates it with stream, if tracing is enabled.
// If tracing is not enabled, it returns nil.
func (s *Server) traceInfo(st transport.ServerTransport, stream *transport.Stream) (trInfo *traceInfo) ***REMOVED***
	if !EnableTracing ***REMOVED***
		return nil
	***REMOVED***
	tr, ok := trace.FromContext(stream.Context())
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	trInfo = &traceInfo***REMOVED***
		tr: tr,
		firstLine: firstLine***REMOVED***
			client:     false,
			remoteAddr: st.RemoteAddr(),
		***REMOVED***,
	***REMOVED***
	if dl, ok := stream.Context().Deadline(); ok ***REMOVED***
		trInfo.firstLine.deadline = time.Until(dl)
	***REMOVED***
	return trInfo
***REMOVED***

func (s *Server) addConn(addr string, st transport.ServerTransport) bool ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil ***REMOVED***
		st.Close()
		return false
	***REMOVED***
	if s.drain ***REMOVED***
		// Transport added after we drained our existing conns: drain it
		// immediately.
		st.Drain()
	***REMOVED***

	if s.conns[addr] == nil ***REMOVED***
		// Create a map entry if this is the first connection on this listener.
		s.conns[addr] = make(map[transport.ServerTransport]bool)
	***REMOVED***
	s.conns[addr][st] = true
	return true
***REMOVED***

func (s *Server) removeConn(addr string, st transport.ServerTransport) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	conns := s.conns[addr]
	if conns != nil ***REMOVED***
		delete(conns, st)
		if len(conns) == 0 ***REMOVED***
			// If the last connection for this address is being removed, also
			// remove the map entry corresponding to the address. This is used
			// in GracefulStop() when waiting for all connections to be closed.
			delete(s.conns, addr)
		***REMOVED***
		s.cv.Broadcast()
	***REMOVED***
***REMOVED***

func (s *Server) channelzMetric() *channelz.ServerInternalMetric ***REMOVED***
	return &channelz.ServerInternalMetric***REMOVED***
		CallsStarted:             atomic.LoadInt64(&s.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&s.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&s.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&s.czData.lastCallStartedTime)),
	***REMOVED***
***REMOVED***

func (s *Server) incrCallsStarted() ***REMOVED***
	atomic.AddInt64(&s.czData.callsStarted, 1)
	atomic.StoreInt64(&s.czData.lastCallStartedTime, time.Now().UnixNano())
***REMOVED***

func (s *Server) incrCallsSucceeded() ***REMOVED***
	atomic.AddInt64(&s.czData.callsSucceeded, 1)
***REMOVED***

func (s *Server) incrCallsFailed() ***REMOVED***
	atomic.AddInt64(&s.czData.callsFailed, 1)
***REMOVED***

func (s *Server) sendResponse(t transport.ServerTransport, stream *transport.Stream, msg interface***REMOVED******REMOVED***, cp Compressor, opts *transport.Options, comp encoding.Compressor) error ***REMOVED***
	data, err := encode(s.getCodec(stream.ContentSubtype()), msg)
	if err != nil ***REMOVED***
		channelz.Error(logger, s.channelzID, "grpc: server failed to encode response: ", err)
		return err
	***REMOVED***
	compData, err := compress(data, cp, comp)
	if err != nil ***REMOVED***
		channelz.Error(logger, s.channelzID, "grpc: server failed to compress response: ", err)
		return err
	***REMOVED***
	hdr, payload := msgHeader(data, compData)
	// TODO(dfawley): should we be checking len(data) instead?
	if len(payload) > s.opts.maxSendMessageSize ***REMOVED***
		return status.Errorf(codes.ResourceExhausted, "grpc: trying to send message larger than max (%d vs. %d)", len(payload), s.opts.maxSendMessageSize)
	***REMOVED***
	err = t.Write(stream, hdr, payload, opts)
	if err == nil && s.opts.statsHandler != nil ***REMOVED***
		s.opts.statsHandler.HandleRPC(stream.Context(), outPayload(false, msg, data, payload, time.Now()))
	***REMOVED***
	return err
***REMOVED***

// chainUnaryServerInterceptors chains all unary server interceptors into one.
func chainUnaryServerInterceptors(s *Server) ***REMOVED***
	// Prepend opts.unaryInt to the chaining interceptors if it exists, since unaryInt will
	// be executed before any other chained interceptors.
	interceptors := s.opts.chainUnaryInts
	if s.opts.unaryInt != nil ***REMOVED***
		interceptors = append([]UnaryServerInterceptor***REMOVED***s.opts.unaryInt***REMOVED***, s.opts.chainUnaryInts...)
	***REMOVED***

	var chainedInt UnaryServerInterceptor
	if len(interceptors) == 0 ***REMOVED***
		chainedInt = nil
	***REMOVED*** else if len(interceptors) == 1 ***REMOVED***
		chainedInt = interceptors[0]
	***REMOVED*** else ***REMOVED***
		chainedInt = chainUnaryInterceptors(interceptors)
	***REMOVED***

	s.opts.unaryInt = chainedInt
***REMOVED***

func chainUnaryInterceptors(interceptors []UnaryServerInterceptor) UnaryServerInterceptor ***REMOVED***
	return func(ctx context.Context, req interface***REMOVED******REMOVED***, info *UnaryServerInfo, handler UnaryHandler) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		// the struct ensures the variables are allocated together, rather than separately, since we
		// know they should be garbage collected together. This saves 1 allocation and decreases
		// time/call by about 10% on the microbenchmark.
		var state struct ***REMOVED***
			i    int
			next UnaryHandler
		***REMOVED***
		state.next = func(ctx context.Context, req interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
			if state.i == len(interceptors)-1 ***REMOVED***
				return interceptors[state.i](ctx, req, info, handler)
			***REMOVED***
			state.i++
			return interceptors[state.i-1](ctx, req, info, state.next)
		***REMOVED***
		return state.next(ctx, req)
	***REMOVED***
***REMOVED***

func (s *Server) processUnaryRPC(t transport.ServerTransport, stream *transport.Stream, info *serviceInfo, md *MethodDesc, trInfo *traceInfo) (err error) ***REMOVED***
	sh := s.opts.statsHandler
	if sh != nil || trInfo != nil || channelz.IsOn() ***REMOVED***
		if channelz.IsOn() ***REMOVED***
			s.incrCallsStarted()
		***REMOVED***
		var statsBegin *stats.Begin
		if sh != nil ***REMOVED***
			beginTime := time.Now()
			statsBegin = &stats.Begin***REMOVED***
				BeginTime:      beginTime,
				IsClientStream: false,
				IsServerStream: false,
			***REMOVED***
			sh.HandleRPC(stream.Context(), statsBegin)
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&trInfo.firstLine, false)
		***REMOVED***
		// The deferred error handling for tracing, stats handler and channelz are
		// combined into one function to reduce stack usage -- a defer takes ~56-64
		// bytes on the stack, so overflowing the stack will require a stack
		// re-allocation, which is expensive.
		//
		// To maintain behavior similar to separate deferred statements, statements
		// should be executed in the reverse order. That is, tracing first, stats
		// handler second, and channelz last. Note that panics *within* defers will
		// lead to different behavior, but that's an acceptable compromise; that
		// would be undefined behavior territory anyway.
		defer func() ***REMOVED***
			if trInfo != nil ***REMOVED***
				if err != nil && err != io.EOF ***REMOVED***
					trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
					trInfo.tr.SetError()
				***REMOVED***
				trInfo.tr.Finish()
			***REMOVED***

			if sh != nil ***REMOVED***
				end := &stats.End***REMOVED***
					BeginTime: statsBegin.BeginTime,
					EndTime:   time.Now(),
				***REMOVED***
				if err != nil && err != io.EOF ***REMOVED***
					end.Error = toRPCErr(err)
				***REMOVED***
				sh.HandleRPC(stream.Context(), end)
			***REMOVED***

			if channelz.IsOn() ***REMOVED***
				if err != nil && err != io.EOF ***REMOVED***
					s.incrCallsFailed()
				***REMOVED*** else ***REMOVED***
					s.incrCallsSucceeded()
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	binlog := binarylog.GetMethodLogger(stream.Method())
	if binlog != nil ***REMOVED***
		ctx := stream.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		logEntry := &binarylog.ClientHeader***REMOVED***
			Header:     md,
			MethodName: stream.Method(),
			PeerAddr:   nil,
		***REMOVED***
		if deadline, ok := ctx.Deadline(); ok ***REMOVED***
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 ***REMOVED***
				logEntry.Timeout = 0
			***REMOVED***
		***REMOVED***
		if a := md[":authority"]; len(a) > 0 ***REMOVED***
			logEntry.Authority = a[0]
		***REMOVED***
		if peer, ok := peer.FromContext(ctx); ok ***REMOVED***
			logEntry.PeerAddr = peer.Addr
		***REMOVED***
		binlog.Log(logEntry)
	***REMOVED***

	// comp and cp are used for compression.  decomp and dc are used for
	// decompression.  If comp and decomp are both set, they are the same;
	// however they are kept separate to ensure that at most one of the
	// compressor/decompressor variable pairs are set for use later.
	var comp, decomp encoding.Compressor
	var cp Compressor
	var dc Decompressor

	// If dc is set and matches the stream's compression, use it.  Otherwise, try
	// to find a matching registered compressor for decomp.
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc ***REMOVED***
		dc = s.opts.dc
	***REMOVED*** else if rc != "" && rc != encoding.Identity ***REMOVED***
		decomp = encoding.GetCompressor(rc)
		if decomp == nil ***REMOVED***
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			t.WriteStatus(stream, st)
			return st.Err()
		***REMOVED***
	***REMOVED***

	// If cp is set, use it.  Otherwise, attempt to compress the response using
	// the incoming message compression method.
	//
	// NOTE: this needs to be ahead of all handling, https://github.com/grpc/grpc-go/issues/686.
	if s.opts.cp != nil ***REMOVED***
		cp = s.opts.cp
		stream.SetSendCompress(cp.Type())
	***REMOVED*** else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity ***REMOVED***
		// Legacy compressor not specified; attempt to respond with same encoding.
		comp = encoding.GetCompressor(rc)
		if comp != nil ***REMOVED***
			stream.SetSendCompress(rc)
		***REMOVED***
	***REMOVED***

	var payInfo *payloadInfo
	if sh != nil || binlog != nil ***REMOVED***
		payInfo = &payloadInfo***REMOVED******REMOVED***
	***REMOVED***
	d, err := recvAndDecompress(&parser***REMOVED***r: stream***REMOVED***, stream, dc, s.opts.maxReceiveMessageSize, payInfo, decomp)
	if err != nil ***REMOVED***
		if e := t.WriteStatus(stream, status.Convert(err)); e != nil ***REMOVED***
			channelz.Warningf(logger, s.channelzID, "grpc: Server.processUnaryRPC failed to write status %v", e)
		***REMOVED***
		return err
	***REMOVED***
	if channelz.IsOn() ***REMOVED***
		t.IncrMsgRecv()
	***REMOVED***
	df := func(v interface***REMOVED******REMOVED***) error ***REMOVED***
		if err := s.getCodec(stream.ContentSubtype()).Unmarshal(d, v); err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
		***REMOVED***
		if sh != nil ***REMOVED***
			sh.HandleRPC(stream.Context(), &stats.InPayload***REMOVED***
				RecvTime:   time.Now(),
				Payload:    v,
				WireLength: payInfo.wireLength + headerLen,
				Data:       d,
				Length:     len(d),
			***REMOVED***)
		***REMOVED***
		if binlog != nil ***REMOVED***
			binlog.Log(&binarylog.ClientMessage***REMOVED***
				Message: d,
			***REMOVED***)
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&payload***REMOVED***sent: false, msg: v***REMOVED***, true)
		***REMOVED***
		return nil
	***REMOVED***
	ctx := NewContextWithServerTransportStream(stream.Context(), stream)
	reply, appErr := md.Handler(info.serviceImpl, ctx, df, s.opts.unaryInt)
	if appErr != nil ***REMOVED***
		appStatus, ok := status.FromError(appErr)
		if !ok ***REMOVED***
			// Convert non-status application error to a status error with code
			// Unknown, but handle context errors specifically.
			appStatus = status.FromContextError(appErr)
			appErr = appStatus.Err()
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			trInfo.tr.SetError()
		***REMOVED***
		if e := t.WriteStatus(stream, appStatus); e != nil ***REMOVED***
			channelz.Warningf(logger, s.channelzID, "grpc: Server.processUnaryRPC failed to write status: %v", e)
		***REMOVED***
		if binlog != nil ***REMOVED***
			if h, _ := stream.Header(); h.Len() > 0 ***REMOVED***
				// Only log serverHeader if there was header. Otherwise it can
				// be trailer only.
				binlog.Log(&binarylog.ServerHeader***REMOVED***
					Header: h,
				***REMOVED***)
			***REMOVED***
			binlog.Log(&binarylog.ServerTrailer***REMOVED***
				Trailer: stream.Trailer(),
				Err:     appErr,
			***REMOVED***)
		***REMOVED***
		return appErr
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyLog(stringer("OK"), false)
	***REMOVED***
	opts := &transport.Options***REMOVED***Last: true***REMOVED***

	if err := s.sendResponse(t, stream, reply, cp, opts, comp); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			// The entire stream is done (for unary RPC only).
			return err
		***REMOVED***
		if sts, ok := status.FromError(err); ok ***REMOVED***
			if e := t.WriteStatus(stream, sts); e != nil ***REMOVED***
				channelz.Warningf(logger, s.channelzID, "grpc: Server.processUnaryRPC failed to write status: %v", e)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch st := err.(type) ***REMOVED***
			case transport.ConnectionError:
				// Nothing to do here.
			default:
				panic(fmt.Sprintf("grpc: Unexpected error (%T) from sendResponse: %v", st, st))
			***REMOVED***
		***REMOVED***
		if binlog != nil ***REMOVED***
			h, _ := stream.Header()
			binlog.Log(&binarylog.ServerHeader***REMOVED***
				Header: h,
			***REMOVED***)
			binlog.Log(&binarylog.ServerTrailer***REMOVED***
				Trailer: stream.Trailer(),
				Err:     appErr,
			***REMOVED***)
		***REMOVED***
		return err
	***REMOVED***
	if binlog != nil ***REMOVED***
		h, _ := stream.Header()
		binlog.Log(&binarylog.ServerHeader***REMOVED***
			Header: h,
		***REMOVED***)
		binlog.Log(&binarylog.ServerMessage***REMOVED***
			Message: reply,
		***REMOVED***)
	***REMOVED***
	if channelz.IsOn() ***REMOVED***
		t.IncrMsgSent()
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyLog(&payload***REMOVED***sent: true, msg: reply***REMOVED***, true)
	***REMOVED***
	// TODO: Should we be logging if writing status failed here, like above?
	// Should the logging be in WriteStatus?  Should we ignore the WriteStatus
	// error or allow the stats handler to see it?
	err = t.WriteStatus(stream, statusOK)
	if binlog != nil ***REMOVED***
		binlog.Log(&binarylog.ServerTrailer***REMOVED***
			Trailer: stream.Trailer(),
			Err:     appErr,
		***REMOVED***)
	***REMOVED***
	return err
***REMOVED***

// chainStreamServerInterceptors chains all stream server interceptors into one.
func chainStreamServerInterceptors(s *Server) ***REMOVED***
	// Prepend opts.streamInt to the chaining interceptors if it exists, since streamInt will
	// be executed before any other chained interceptors.
	interceptors := s.opts.chainStreamInts
	if s.opts.streamInt != nil ***REMOVED***
		interceptors = append([]StreamServerInterceptor***REMOVED***s.opts.streamInt***REMOVED***, s.opts.chainStreamInts...)
	***REMOVED***

	var chainedInt StreamServerInterceptor
	if len(interceptors) == 0 ***REMOVED***
		chainedInt = nil
	***REMOVED*** else if len(interceptors) == 1 ***REMOVED***
		chainedInt = interceptors[0]
	***REMOVED*** else ***REMOVED***
		chainedInt = chainStreamInterceptors(interceptors)
	***REMOVED***

	s.opts.streamInt = chainedInt
***REMOVED***

func chainStreamInterceptors(interceptors []StreamServerInterceptor) StreamServerInterceptor ***REMOVED***
	return func(srv interface***REMOVED******REMOVED***, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error ***REMOVED***
		// the struct ensures the variables are allocated together, rather than separately, since we
		// know they should be garbage collected together. This saves 1 allocation and decreases
		// time/call by about 10% on the microbenchmark.
		var state struct ***REMOVED***
			i    int
			next StreamHandler
		***REMOVED***
		state.next = func(srv interface***REMOVED******REMOVED***, ss ServerStream) error ***REMOVED***
			if state.i == len(interceptors)-1 ***REMOVED***
				return interceptors[state.i](srv, ss, info, handler)
			***REMOVED***
			state.i++
			return interceptors[state.i-1](srv, ss, info, state.next)
		***REMOVED***
		return state.next(srv, ss)
	***REMOVED***
***REMOVED***

func (s *Server) processStreamingRPC(t transport.ServerTransport, stream *transport.Stream, info *serviceInfo, sd *StreamDesc, trInfo *traceInfo) (err error) ***REMOVED***
	if channelz.IsOn() ***REMOVED***
		s.incrCallsStarted()
	***REMOVED***
	sh := s.opts.statsHandler
	var statsBegin *stats.Begin
	if sh != nil ***REMOVED***
		beginTime := time.Now()
		statsBegin = &stats.Begin***REMOVED***
			BeginTime:      beginTime,
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		***REMOVED***
		sh.HandleRPC(stream.Context(), statsBegin)
	***REMOVED***
	ctx := NewContextWithServerTransportStream(stream.Context(), stream)
	ss := &serverStream***REMOVED***
		ctx:                   ctx,
		t:                     t,
		s:                     stream,
		p:                     &parser***REMOVED***r: stream***REMOVED***,
		codec:                 s.getCodec(stream.ContentSubtype()),
		maxReceiveMessageSize: s.opts.maxReceiveMessageSize,
		maxSendMessageSize:    s.opts.maxSendMessageSize,
		trInfo:                trInfo,
		statsHandler:          sh,
	***REMOVED***

	if sh != nil || trInfo != nil || channelz.IsOn() ***REMOVED***
		// See comment in processUnaryRPC on defers.
		defer func() ***REMOVED***
			if trInfo != nil ***REMOVED***
				ss.mu.Lock()
				if err != nil && err != io.EOF ***REMOVED***
					ss.trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
					ss.trInfo.tr.SetError()
				***REMOVED***
				ss.trInfo.tr.Finish()
				ss.trInfo.tr = nil
				ss.mu.Unlock()
			***REMOVED***

			if sh != nil ***REMOVED***
				end := &stats.End***REMOVED***
					BeginTime: statsBegin.BeginTime,
					EndTime:   time.Now(),
				***REMOVED***
				if err != nil && err != io.EOF ***REMOVED***
					end.Error = toRPCErr(err)
				***REMOVED***
				sh.HandleRPC(stream.Context(), end)
			***REMOVED***

			if channelz.IsOn() ***REMOVED***
				if err != nil && err != io.EOF ***REMOVED***
					s.incrCallsFailed()
				***REMOVED*** else ***REMOVED***
					s.incrCallsSucceeded()
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	ss.binlog = binarylog.GetMethodLogger(stream.Method())
	if ss.binlog != nil ***REMOVED***
		md, _ := metadata.FromIncomingContext(ctx)
		logEntry := &binarylog.ClientHeader***REMOVED***
			Header:     md,
			MethodName: stream.Method(),
			PeerAddr:   nil,
		***REMOVED***
		if deadline, ok := ctx.Deadline(); ok ***REMOVED***
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 ***REMOVED***
				logEntry.Timeout = 0
			***REMOVED***
		***REMOVED***
		if a := md[":authority"]; len(a) > 0 ***REMOVED***
			logEntry.Authority = a[0]
		***REMOVED***
		if peer, ok := peer.FromContext(ss.Context()); ok ***REMOVED***
			logEntry.PeerAddr = peer.Addr
		***REMOVED***
		ss.binlog.Log(logEntry)
	***REMOVED***

	// If dc is set and matches the stream's compression, use it.  Otherwise, try
	// to find a matching registered compressor for decomp.
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc ***REMOVED***
		ss.dc = s.opts.dc
	***REMOVED*** else if rc != "" && rc != encoding.Identity ***REMOVED***
		ss.decomp = encoding.GetCompressor(rc)
		if ss.decomp == nil ***REMOVED***
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			t.WriteStatus(ss.s, st)
			return st.Err()
		***REMOVED***
	***REMOVED***

	// If cp is set, use it.  Otherwise, attempt to compress the response using
	// the incoming message compression method.
	//
	// NOTE: this needs to be ahead of all handling, https://github.com/grpc/grpc-go/issues/686.
	if s.opts.cp != nil ***REMOVED***
		ss.cp = s.opts.cp
		stream.SetSendCompress(s.opts.cp.Type())
	***REMOVED*** else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity ***REMOVED***
		// Legacy compressor not specified; attempt to respond with same encoding.
		ss.comp = encoding.GetCompressor(rc)
		if ss.comp != nil ***REMOVED***
			stream.SetSendCompress(rc)
		***REMOVED***
	***REMOVED***

	ss.ctx = newContextWithRPCInfo(ss.ctx, false, ss.codec, ss.cp, ss.comp)

	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
	***REMOVED***
	var appErr error
	var server interface***REMOVED******REMOVED***
	if info != nil ***REMOVED***
		server = info.serviceImpl
	***REMOVED***
	if s.opts.streamInt == nil ***REMOVED***
		appErr = sd.Handler(server, ss)
	***REMOVED*** else ***REMOVED***
		info := &StreamServerInfo***REMOVED***
			FullMethod:     stream.Method(),
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		***REMOVED***
		appErr = s.opts.streamInt(server, ss, info, sd.Handler)
	***REMOVED***
	if appErr != nil ***REMOVED***
		appStatus, ok := status.FromError(appErr)
		if !ok ***REMOVED***
			// Convert non-status application error to a status error with code
			// Unknown, but handle context errors specifically.
			appStatus = status.FromContextError(appErr)
			appErr = appStatus.Err()
		***REMOVED***
		if trInfo != nil ***REMOVED***
			ss.mu.Lock()
			ss.trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			ss.trInfo.tr.SetError()
			ss.mu.Unlock()
		***REMOVED***
		t.WriteStatus(ss.s, appStatus)
		if ss.binlog != nil ***REMOVED***
			ss.binlog.Log(&binarylog.ServerTrailer***REMOVED***
				Trailer: ss.s.Trailer(),
				Err:     appErr,
			***REMOVED***)
		***REMOVED***
		// TODO: Should we log an error from WriteStatus here and below?
		return appErr
	***REMOVED***
	if trInfo != nil ***REMOVED***
		ss.mu.Lock()
		ss.trInfo.tr.LazyLog(stringer("OK"), false)
		ss.mu.Unlock()
	***REMOVED***
	err = t.WriteStatus(ss.s, statusOK)
	if ss.binlog != nil ***REMOVED***
		ss.binlog.Log(&binarylog.ServerTrailer***REMOVED***
			Trailer: ss.s.Trailer(),
			Err:     appErr,
		***REMOVED***)
	***REMOVED***
	return err
***REMOVED***

func (s *Server) handleStream(t transport.ServerTransport, stream *transport.Stream, trInfo *traceInfo) ***REMOVED***
	sm := stream.Method()
	if sm != "" && sm[0] == '/' ***REMOVED***
		sm = sm[1:]
	***REMOVED***
	pos := strings.LastIndex(sm, "/")
	if pos == -1 ***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&fmtStringer***REMOVED***"Malformed method name %q", []interface***REMOVED******REMOVED******REMOVED***sm***REMOVED******REMOVED***, true)
			trInfo.tr.SetError()
		***REMOVED***
		errDesc := fmt.Sprintf("malformed method name: %q", stream.Method())
		if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil ***REMOVED***
			if trInfo != nil ***REMOVED***
				trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
				trInfo.tr.SetError()
			***REMOVED***
			channelz.Warningf(logger, s.channelzID, "grpc: Server.handleStream failed to write status: %v", err)
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.Finish()
		***REMOVED***
		return
	***REMOVED***
	service := sm[:pos]
	method := sm[pos+1:]

	srv, knownService := s.services[service]
	if knownService ***REMOVED***
		if md, ok := srv.methods[method]; ok ***REMOVED***
			s.processUnaryRPC(t, stream, srv, md, trInfo)
			return
		***REMOVED***
		if sd, ok := srv.streams[method]; ok ***REMOVED***
			s.processStreamingRPC(t, stream, srv, sd, trInfo)
			return
		***REMOVED***
	***REMOVED***
	// Unknown service, or known server unknown method.
	if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil ***REMOVED***
		s.processStreamingRPC(t, stream, nil, unknownDesc, trInfo)
		return
	***REMOVED***
	var errDesc string
	if !knownService ***REMOVED***
		errDesc = fmt.Sprintf("unknown service %v", service)
	***REMOVED*** else ***REMOVED***
		errDesc = fmt.Sprintf("unknown method %v for service %v", method, service)
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyPrintf("%s", errDesc)
		trInfo.tr.SetError()
	***REMOVED***
	if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil ***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
			trInfo.tr.SetError()
		***REMOVED***
		channelz.Warningf(logger, s.channelzID, "grpc: Server.handleStream failed to write status: %v", err)
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.Finish()
	***REMOVED***
***REMOVED***

// The key to save ServerTransportStream in the context.
type streamKey struct***REMOVED******REMOVED***

// NewContextWithServerTransportStream creates a new context from ctx and
// attaches stream to it.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func NewContextWithServerTransportStream(ctx context.Context, stream ServerTransportStream) context.Context ***REMOVED***
	return context.WithValue(ctx, streamKey***REMOVED******REMOVED***, stream)
***REMOVED***

// ServerTransportStream is a minimal interface that a transport stream must
// implement. This can be used to mock an actual transport stream for tests of
// handler code that use, for example, grpc.SetHeader (which requires some
// stream to be in context).
//
// See also NewContextWithServerTransportStream.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ServerTransportStream interface ***REMOVED***
	Method() string
	SetHeader(md metadata.MD) error
	SendHeader(md metadata.MD) error
	SetTrailer(md metadata.MD) error
***REMOVED***

// ServerTransportStreamFromContext returns the ServerTransportStream saved in
// ctx. Returns nil if the given context has no stream associated with it
// (which implies it is not an RPC invocation context).
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func ServerTransportStreamFromContext(ctx context.Context) ServerTransportStream ***REMOVED***
	s, _ := ctx.Value(streamKey***REMOVED******REMOVED***).(ServerTransportStream)
	return s
***REMOVED***

// Stop stops the gRPC server. It immediately closes all open
// connections and listeners.
// It cancels all active RPCs on the server side and the corresponding
// pending RPCs on the client side will get notified by connection
// errors.
func (s *Server) Stop() ***REMOVED***
	s.quit.Fire()

	defer func() ***REMOVED***
		s.serveWG.Wait()
		s.done.Fire()
	***REMOVED***()

	s.channelzRemoveOnce.Do(func() ***REMOVED***
		if channelz.IsOn() ***REMOVED***
			channelz.RemoveEntry(s.channelzID)
		***REMOVED***
	***REMOVED***)

	s.mu.Lock()
	listeners := s.lis
	s.lis = nil
	conns := s.conns
	s.conns = nil
	// interrupt GracefulStop if Stop and GracefulStop are called concurrently.
	s.cv.Broadcast()
	s.mu.Unlock()

	for lis := range listeners ***REMOVED***
		lis.Close()
	***REMOVED***
	for _, cs := range conns ***REMOVED***
		for st := range cs ***REMOVED***
			st.Close()
		***REMOVED***
	***REMOVED***
	if s.opts.numServerWorkers > 0 ***REMOVED***
		s.stopServerWorkers()
	***REMOVED***

	s.mu.Lock()
	if s.events != nil ***REMOVED***
		s.events.Finish()
		s.events = nil
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// GracefulStop stops the gRPC server gracefully. It stops the server from
// accepting new connections and RPCs and blocks until all the pending RPCs are
// finished.
func (s *Server) GracefulStop() ***REMOVED***
	s.quit.Fire()
	defer s.done.Fire()

	s.channelzRemoveOnce.Do(func() ***REMOVED***
		if channelz.IsOn() ***REMOVED***
			channelz.RemoveEntry(s.channelzID)
		***REMOVED***
	***REMOVED***)
	s.mu.Lock()
	if s.conns == nil ***REMOVED***
		s.mu.Unlock()
		return
	***REMOVED***

	for lis := range s.lis ***REMOVED***
		lis.Close()
	***REMOVED***
	s.lis = nil
	if !s.drain ***REMOVED***
		for _, conns := range s.conns ***REMOVED***
			for st := range conns ***REMOVED***
				st.Drain()
			***REMOVED***
		***REMOVED***
		s.drain = true
	***REMOVED***

	// Wait for serving threads to be ready to exit.  Only then can we be sure no
	// new conns will be created.
	s.mu.Unlock()
	s.serveWG.Wait()
	s.mu.Lock()

	for len(s.conns) != 0 ***REMOVED***
		s.cv.Wait()
	***REMOVED***
	s.conns = nil
	if s.events != nil ***REMOVED***
		s.events.Finish()
		s.events = nil
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// contentSubtype must be lowercase
// cannot return nil
func (s *Server) getCodec(contentSubtype string) baseCodec ***REMOVED***
	if s.opts.codec != nil ***REMOVED***
		return s.opts.codec
	***REMOVED***
	if contentSubtype == "" ***REMOVED***
		return encoding.GetCodec(proto.Name)
	***REMOVED***
	codec := encoding.GetCodec(contentSubtype)
	if codec == nil ***REMOVED***
		return encoding.GetCodec(proto.Name)
	***REMOVED***
	return codec
***REMOVED***

// SetHeader sets the header metadata.
// When called multiple times, all the provided metadata will be merged.
// All the metadata will be sent out when one of the following happens:
//  - grpc.SendHeader() is called;
//  - The first response is sent out;
//  - An RPC status is sent out (error or success).
func SetHeader(ctx context.Context, md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	return stream.SetHeader(md)
***REMOVED***

// SendHeader sends header metadata. It may be called at most once.
// The provided md and headers set by SetHeader() will be sent.
func SendHeader(ctx context.Context, md metadata.MD) error ***REMOVED***
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	if err := stream.SendHeader(md); err != nil ***REMOVED***
		return toRPCErr(err)
	***REMOVED***
	return nil
***REMOVED***

// SetTrailer sets the trailer metadata that will be sent when an RPC returns.
// When called more than once, all the provided metadata will be merged.
func SetTrailer(ctx context.Context, md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	return stream.SetTrailer(md)
***REMOVED***

// Method returns the method string for the server context.  The returned
// string is in the format of "/service/method".
func Method(ctx context.Context) (string, bool) ***REMOVED***
	s := ServerTransportStreamFromContext(ctx)
	if s == nil ***REMOVED***
		return "", false
	***REMOVED***
	return s.Method(), true
***REMOVED***

type channelzServer struct ***REMOVED***
	s *Server
***REMOVED***

func (c *channelzServer) ChannelzMetric() *channelz.ServerInternalMetric ***REMOVED***
	return c.s.channelzMetric()
***REMOVED***
