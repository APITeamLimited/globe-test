/*
 *
 * Copyright 2016 gRPC authors.
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

// Package stats is for collecting and reporting various network and RPC stats.
// This package is for monitoring purpose only. All fields are read-only.
// All APIs are experimental.
package stats // import "google.golang.org/grpc/stats"

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
)

// RPCStats contains stats information about RPCs.
type RPCStats interface ***REMOVED***
	isRPCStats()
	// IsClient returns true if this RPCStats is from client side.
	IsClient() bool
***REMOVED***

// Begin contains stats when an RPC attempt begins.
// FailFast is only valid if this Begin is from client side.
type Begin struct ***REMOVED***
	// Client is true if this Begin is from client side.
	Client bool
	// BeginTime is the time when the RPC attempt begins.
	BeginTime time.Time
	// FailFast indicates if this RPC is failfast.
	FailFast bool
	// IsClientStream indicates whether the RPC is a client streaming RPC.
	IsClientStream bool
	// IsServerStream indicates whether the RPC is a server streaming RPC.
	IsServerStream bool
	// IsTransparentRetryAttempt indicates whether this attempt was initiated
	// due to transparently retrying a previous attempt.
	IsTransparentRetryAttempt bool
***REMOVED***

// IsClient indicates if the stats information is from client side.
func (s *Begin) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *Begin) isRPCStats() ***REMOVED******REMOVED***

// InPayload contains the information for an incoming payload.
type InPayload struct ***REMOVED***
	// Client is true if this InPayload is from client side.
	Client bool
	// Payload is the payload with original type.
	Payload interface***REMOVED******REMOVED***
	// Data is the serialized message payload.
	Data []byte
	// Length is the length of uncompressed data.
	Length int
	// WireLength is the length of data on wire (compressed, signed, encrypted).
	WireLength int
	// RecvTime is the time when the payload is received.
	RecvTime time.Time
***REMOVED***

// IsClient indicates if the stats information is from client side.
func (s *InPayload) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *InPayload) isRPCStats() ***REMOVED******REMOVED***

// InHeader contains stats when a header is received.
type InHeader struct ***REMOVED***
	// Client is true if this InHeader is from client side.
	Client bool
	// WireLength is the wire length of header.
	WireLength int
	// Compression is the compression algorithm used for the RPC.
	Compression string
	// Header contains the header metadata received.
	Header metadata.MD

	// The following fields are valid only if Client is false.
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
	// RemoteAddr is the remote address of the corresponding connection.
	RemoteAddr net.Addr
	// LocalAddr is the local address of the corresponding connection.
	LocalAddr net.Addr
***REMOVED***

// IsClient indicates if the stats information is from client side.
func (s *InHeader) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *InHeader) isRPCStats() ***REMOVED******REMOVED***

// InTrailer contains stats when a trailer is received.
type InTrailer struct ***REMOVED***
	// Client is true if this InTrailer is from client side.
	Client bool
	// WireLength is the wire length of trailer.
	WireLength int
	// Trailer contains the trailer metadata received from the server. This
	// field is only valid if this InTrailer is from the client side.
	Trailer metadata.MD
***REMOVED***

// IsClient indicates if the stats information is from client side.
func (s *InTrailer) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *InTrailer) isRPCStats() ***REMOVED******REMOVED***

// OutPayload contains the information for an outgoing payload.
type OutPayload struct ***REMOVED***
	// Client is true if this OutPayload is from client side.
	Client bool
	// Payload is the payload with original type.
	Payload interface***REMOVED******REMOVED***
	// Data is the serialized message payload.
	Data []byte
	// Length is the length of uncompressed data.
	Length int
	// WireLength is the length of data on wire (compressed, signed, encrypted).
	WireLength int
	// SentTime is the time when the payload is sent.
	SentTime time.Time
***REMOVED***

// IsClient indicates if this stats information is from client side.
func (s *OutPayload) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *OutPayload) isRPCStats() ***REMOVED******REMOVED***

// OutHeader contains stats when a header is sent.
type OutHeader struct ***REMOVED***
	// Client is true if this OutHeader is from client side.
	Client bool
	// Compression is the compression algorithm used for the RPC.
	Compression string
	// Header contains the header metadata sent.
	Header metadata.MD

	// The following fields are valid only if Client is true.
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
	// RemoteAddr is the remote address of the corresponding connection.
	RemoteAddr net.Addr
	// LocalAddr is the local address of the corresponding connection.
	LocalAddr net.Addr
***REMOVED***

// IsClient indicates if this stats information is from client side.
func (s *OutHeader) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *OutHeader) isRPCStats() ***REMOVED******REMOVED***

// OutTrailer contains stats when a trailer is sent.
type OutTrailer struct ***REMOVED***
	// Client is true if this OutTrailer is from client side.
	Client bool
	// WireLength is the wire length of trailer.
	//
	// Deprecated: This field is never set. The length is not known when this message is
	// emitted because the trailer fields are compressed with hpack after that.
	WireLength int
	// Trailer contains the trailer metadata sent to the client. This
	// field is only valid if this OutTrailer is from the server side.
	Trailer metadata.MD
***REMOVED***

// IsClient indicates if this stats information is from client side.
func (s *OutTrailer) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *OutTrailer) isRPCStats() ***REMOVED******REMOVED***

// End contains stats when an RPC ends.
type End struct ***REMOVED***
	// Client is true if this End is from client side.
	Client bool
	// BeginTime is the time when the RPC began.
	BeginTime time.Time
	// EndTime is the time when the RPC ends.
	EndTime time.Time
	// Trailer contains the trailer metadata received from the server. This
	// field is only valid if this End is from the client side.
	// Deprecated: use Trailer in InTrailer instead.
	Trailer metadata.MD
	// Error is the error the RPC ended with. It is an error generated from
	// status.Status and can be converted back to status.Status using
	// status.FromError if non-nil.
	Error error
***REMOVED***

// IsClient indicates if this is from client side.
func (s *End) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *End) isRPCStats() ***REMOVED******REMOVED***

// ConnStats contains stats information about connections.
type ConnStats interface ***REMOVED***
	isConnStats()
	// IsClient returns true if this ConnStats is from client side.
	IsClient() bool
***REMOVED***

// ConnBegin contains the stats of a connection when it is established.
type ConnBegin struct ***REMOVED***
	// Client is true if this ConnBegin is from client side.
	Client bool
***REMOVED***

// IsClient indicates if this is from client side.
func (s *ConnBegin) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *ConnBegin) isConnStats() ***REMOVED******REMOVED***

// ConnEnd contains the stats of a connection when it ends.
type ConnEnd struct ***REMOVED***
	// Client is true if this ConnEnd is from client side.
	Client bool
***REMOVED***

// IsClient indicates if this is from client side.
func (s *ConnEnd) IsClient() bool ***REMOVED*** return s.Client ***REMOVED***

func (s *ConnEnd) isConnStats() ***REMOVED******REMOVED***

type incomingTagsKey struct***REMOVED******REMOVED***
type outgoingTagsKey struct***REMOVED******REMOVED***

// SetTags attaches stats tagging data to the context, which will be sent in
// the outgoing RPC with the header grpc-tags-bin.  Subsequent calls to
// SetTags will overwrite the values from earlier calls.
//
// NOTE: this is provided only for backward compatibility with existing clients
// and will likely be removed in an upcoming release.  New uses should transmit
// this type of data using metadata with a different, non-reserved (i.e. does
// not begin with "grpc-") header name.
func SetTags(ctx context.Context, b []byte) context.Context ***REMOVED***
	return context.WithValue(ctx, outgoingTagsKey***REMOVED******REMOVED***, b)
***REMOVED***

// Tags returns the tags from the context for the inbound RPC.
//
// NOTE: this is provided only for backward compatibility with existing clients
// and will likely be removed in an upcoming release.  New uses should transmit
// this type of data using metadata with a different, non-reserved (i.e. does
// not begin with "grpc-") header name.
func Tags(ctx context.Context) []byte ***REMOVED***
	b, _ := ctx.Value(incomingTagsKey***REMOVED******REMOVED***).([]byte)
	return b
***REMOVED***

// SetIncomingTags attaches stats tagging data to the context, to be read by
// the application (not sent in outgoing RPCs).
//
// This is intended for gRPC-internal use ONLY.
func SetIncomingTags(ctx context.Context, b []byte) context.Context ***REMOVED***
	return context.WithValue(ctx, incomingTagsKey***REMOVED******REMOVED***, b)
***REMOVED***

// OutgoingTags returns the tags from the context for the outbound RPC.
//
// This is intended for gRPC-internal use ONLY.
func OutgoingTags(ctx context.Context) []byte ***REMOVED***
	b, _ := ctx.Value(outgoingTagsKey***REMOVED******REMOVED***).([]byte)
	return b
***REMOVED***

type incomingTraceKey struct***REMOVED******REMOVED***
type outgoingTraceKey struct***REMOVED******REMOVED***

// SetTrace attaches stats tagging data to the context, which will be sent in
// the outgoing RPC with the header grpc-trace-bin.  Subsequent calls to
// SetTrace will overwrite the values from earlier calls.
//
// NOTE: this is provided only for backward compatibility with existing clients
// and will likely be removed in an upcoming release.  New uses should transmit
// this type of data using metadata with a different, non-reserved (i.e. does
// not begin with "grpc-") header name.
func SetTrace(ctx context.Context, b []byte) context.Context ***REMOVED***
	return context.WithValue(ctx, outgoingTraceKey***REMOVED******REMOVED***, b)
***REMOVED***

// Trace returns the trace from the context for the inbound RPC.
//
// NOTE: this is provided only for backward compatibility with existing clients
// and will likely be removed in an upcoming release.  New uses should transmit
// this type of data using metadata with a different, non-reserved (i.e. does
// not begin with "grpc-") header name.
func Trace(ctx context.Context) []byte ***REMOVED***
	b, _ := ctx.Value(incomingTraceKey***REMOVED******REMOVED***).([]byte)
	return b
***REMOVED***

// SetIncomingTrace attaches stats tagging data to the context, to be read by
// the application (not sent in outgoing RPCs).  It is intended for
// gRPC-internal use.
func SetIncomingTrace(ctx context.Context, b []byte) context.Context ***REMOVED***
	return context.WithValue(ctx, incomingTraceKey***REMOVED******REMOVED***, b)
***REMOVED***

// OutgoingTrace returns the trace from the context for the outbound RPC.  It is
// intended for gRPC-internal use.
func OutgoingTrace(ctx context.Context) []byte ***REMOVED***
	b, _ := ctx.Value(outgoingTraceKey***REMOVED******REMOVED***).([]byte)
	return b
***REMOVED***
