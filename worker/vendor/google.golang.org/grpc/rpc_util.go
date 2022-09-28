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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// Compressor defines the interface gRPC uses to compress a message.
//
// Deprecated: use package encoding.
type Compressor interface ***REMOVED***
	// Do compresses p into w.
	Do(w io.Writer, p []byte) error
	// Type returns the compression algorithm the Compressor uses.
	Type() string
***REMOVED***

type gzipCompressor struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewGZIPCompressor creates a Compressor based on GZIP.
//
// Deprecated: use package encoding/gzip.
func NewGZIPCompressor() Compressor ***REMOVED***
	c, _ := NewGZIPCompressorWithLevel(gzip.DefaultCompression)
	return c
***REMOVED***

// NewGZIPCompressorWithLevel is like NewGZIPCompressor but specifies the gzip compression level instead
// of assuming DefaultCompression.
//
// The error returned will be nil if the level is valid.
//
// Deprecated: use package encoding/gzip.
func NewGZIPCompressorWithLevel(level int) (Compressor, error) ***REMOVED***
	if level < gzip.DefaultCompression || level > gzip.BestCompression ***REMOVED***
		return nil, fmt.Errorf("grpc: invalid compression level: %d", level)
	***REMOVED***
	return &gzipCompressor***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				w, err := gzip.NewWriterLevel(ioutil.Discard, level)
				if err != nil ***REMOVED***
					panic(err)
				***REMOVED***
				return w
			***REMOVED***,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (c *gzipCompressor) Do(w io.Writer, p []byte) error ***REMOVED***
	z := c.pool.Get().(*gzip.Writer)
	defer c.pool.Put(z)
	z.Reset(w)
	if _, err := z.Write(p); err != nil ***REMOVED***
		return err
	***REMOVED***
	return z.Close()
***REMOVED***

func (c *gzipCompressor) Type() string ***REMOVED***
	return "gzip"
***REMOVED***

// Decompressor defines the interface gRPC uses to decompress a message.
//
// Deprecated: use package encoding.
type Decompressor interface ***REMOVED***
	// Do reads the data from r and uncompress them.
	Do(r io.Reader) ([]byte, error)
	// Type returns the compression algorithm the Decompressor uses.
	Type() string
***REMOVED***

type gzipDecompressor struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewGZIPDecompressor creates a Decompressor based on GZIP.
//
// Deprecated: use package encoding/gzip.
func NewGZIPDecompressor() Decompressor ***REMOVED***
	return &gzipDecompressor***REMOVED******REMOVED***
***REMOVED***

func (d *gzipDecompressor) Do(r io.Reader) ([]byte, error) ***REMOVED***
	var z *gzip.Reader
	switch maybeZ := d.pool.Get().(type) ***REMOVED***
	case nil:
		newZ, err := gzip.NewReader(r)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		z = newZ
	case *gzip.Reader:
		z = maybeZ
		if err := z.Reset(r); err != nil ***REMOVED***
			d.pool.Put(z)
			return nil, err
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		z.Close()
		d.pool.Put(z)
	***REMOVED***()
	return ioutil.ReadAll(z)
***REMOVED***

func (d *gzipDecompressor) Type() string ***REMOVED***
	return "gzip"
***REMOVED***

// callInfo contains all related configuration and information about an RPC.
type callInfo struct ***REMOVED***
	compressorType        string
	failFast              bool
	maxReceiveMessageSize *int
	maxSendMessageSize    *int
	creds                 credentials.PerRPCCredentials
	contentSubtype        string
	codec                 baseCodec
	maxRetryRPCBufferSize int
***REMOVED***

func defaultCallInfo() *callInfo ***REMOVED***
	return &callInfo***REMOVED***
		failFast:              true,
		maxRetryRPCBufferSize: 256 * 1024, // 256KB
	***REMOVED***
***REMOVED***

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface ***REMOVED***
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*callInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*callInfo, *csAttempt)
***REMOVED***

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct***REMOVED******REMOVED***

func (EmptyCallOption) before(*callInfo) error      ***REMOVED*** return nil ***REMOVED***
func (EmptyCallOption) after(*callInfo, *csAttempt) ***REMOVED******REMOVED***

// Header returns a CallOptions that retrieves the header metadata
// for a unary RPC.
func Header(md *metadata.MD) CallOption ***REMOVED***
	return HeaderCallOption***REMOVED***HeaderAddr: md***REMOVED***
***REMOVED***

// HeaderCallOption is a CallOption for collecting response header metadata.
// The metadata field will be populated *after* the RPC completes.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type HeaderCallOption struct ***REMOVED***
	HeaderAddr *metadata.MD
***REMOVED***

func (o HeaderCallOption) before(c *callInfo) error ***REMOVED*** return nil ***REMOVED***
func (o HeaderCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED***
	*o.HeaderAddr, _ = attempt.s.Header()
***REMOVED***

// Trailer returns a CallOptions that retrieves the trailer metadata
// for a unary RPC.
func Trailer(md *metadata.MD) CallOption ***REMOVED***
	return TrailerCallOption***REMOVED***TrailerAddr: md***REMOVED***
***REMOVED***

// TrailerCallOption is a CallOption for collecting response trailer metadata.
// The metadata field will be populated *after* the RPC completes.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type TrailerCallOption struct ***REMOVED***
	TrailerAddr *metadata.MD
***REMOVED***

func (o TrailerCallOption) before(c *callInfo) error ***REMOVED*** return nil ***REMOVED***
func (o TrailerCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED***
	*o.TrailerAddr = attempt.s.Trailer()
***REMOVED***

// Peer returns a CallOption that retrieves peer information for a unary RPC.
// The peer field will be populated *after* the RPC completes.
func Peer(p *peer.Peer) CallOption ***REMOVED***
	return PeerCallOption***REMOVED***PeerAddr: p***REMOVED***
***REMOVED***

// PeerCallOption is a CallOption for collecting the identity of the remote
// peer. The peer field will be populated *after* the RPC completes.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type PeerCallOption struct ***REMOVED***
	PeerAddr *peer.Peer
***REMOVED***

func (o PeerCallOption) before(c *callInfo) error ***REMOVED*** return nil ***REMOVED***
func (o PeerCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED***
	if x, ok := peer.FromContext(attempt.s.Context()); ok ***REMOVED***
		*o.PeerAddr = *x
	***REMOVED***
***REMOVED***

// WaitForReady configures the action to take when an RPC is attempted on broken
// connections or unreachable servers. If waitForReady is false and the
// connection is in the TRANSIENT_FAILURE state, the RPC will fail
// immediately. Otherwise, the RPC client will block the call until a
// connection is available (or the call is canceled or times out) and will
// retry the call if it fails due to a transient error.  gRPC will not retry if
// data was written to the wire unless the server indicates it did not process
// the data.  Please refer to
// https://github.com/grpc/grpc/blob/master/doc/wait-for-ready.md.
//
// By default, RPCs don't "wait for ready".
func WaitForReady(waitForReady bool) CallOption ***REMOVED***
	return FailFastCallOption***REMOVED***FailFast: !waitForReady***REMOVED***
***REMOVED***

// FailFast is the opposite of WaitForReady.
//
// Deprecated: use WaitForReady.
func FailFast(failFast bool) CallOption ***REMOVED***
	return FailFastCallOption***REMOVED***FailFast: failFast***REMOVED***
***REMOVED***

// FailFastCallOption is a CallOption for indicating whether an RPC should fail
// fast or not.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type FailFastCallOption struct ***REMOVED***
	FailFast bool
***REMOVED***

func (o FailFastCallOption) before(c *callInfo) error ***REMOVED***
	c.failFast = o.FailFast
	return nil
***REMOVED***
func (o FailFastCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// MaxCallRecvMsgSize returns a CallOption which sets the maximum message size
// in bytes the client can receive.
func MaxCallRecvMsgSize(bytes int) CallOption ***REMOVED***
	return MaxRecvMsgSizeCallOption***REMOVED***MaxRecvMsgSize: bytes***REMOVED***
***REMOVED***

// MaxRecvMsgSizeCallOption is a CallOption that indicates the maximum message
// size in bytes the client can receive.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type MaxRecvMsgSizeCallOption struct ***REMOVED***
	MaxRecvMsgSize int
***REMOVED***

func (o MaxRecvMsgSizeCallOption) before(c *callInfo) error ***REMOVED***
	c.maxReceiveMessageSize = &o.MaxRecvMsgSize
	return nil
***REMOVED***
func (o MaxRecvMsgSizeCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// MaxCallSendMsgSize returns a CallOption which sets the maximum message size
// in bytes the client can send.
func MaxCallSendMsgSize(bytes int) CallOption ***REMOVED***
	return MaxSendMsgSizeCallOption***REMOVED***MaxSendMsgSize: bytes***REMOVED***
***REMOVED***

// MaxSendMsgSizeCallOption is a CallOption that indicates the maximum message
// size in bytes the client can send.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type MaxSendMsgSizeCallOption struct ***REMOVED***
	MaxSendMsgSize int
***REMOVED***

func (o MaxSendMsgSizeCallOption) before(c *callInfo) error ***REMOVED***
	c.maxSendMessageSize = &o.MaxSendMsgSize
	return nil
***REMOVED***
func (o MaxSendMsgSizeCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// PerRPCCredentials returns a CallOption that sets credentials.PerRPCCredentials
// for a call.
func PerRPCCredentials(creds credentials.PerRPCCredentials) CallOption ***REMOVED***
	return PerRPCCredsCallOption***REMOVED***Creds: creds***REMOVED***
***REMOVED***

// PerRPCCredsCallOption is a CallOption that indicates the per-RPC
// credentials to use for the call.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type PerRPCCredsCallOption struct ***REMOVED***
	Creds credentials.PerRPCCredentials
***REMOVED***

func (o PerRPCCredsCallOption) before(c *callInfo) error ***REMOVED***
	c.creds = o.Creds
	return nil
***REMOVED***
func (o PerRPCCredsCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// UseCompressor returns a CallOption which sets the compressor used when
// sending the request.  If WithCompressor is also set, UseCompressor has
// higher priority.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func UseCompressor(name string) CallOption ***REMOVED***
	return CompressorCallOption***REMOVED***CompressorType: name***REMOVED***
***REMOVED***

// CompressorCallOption is a CallOption that indicates the compressor to use.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type CompressorCallOption struct ***REMOVED***
	CompressorType string
***REMOVED***

func (o CompressorCallOption) before(c *callInfo) error ***REMOVED***
	c.compressorType = o.CompressorType
	return nil
***REMOVED***
func (o CompressorCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// CallContentSubtype returns a CallOption that will set the content-subtype
// for a call. For example, if content-subtype is "json", the Content-Type over
// the wire will be "application/grpc+json". The content-subtype is converted
// to lowercase before being included in Content-Type. See Content-Type on
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details.
//
// If ForceCodec is not also used, the content-subtype will be used to look up
// the Codec to use in the registry controlled by RegisterCodec. See the
// documentation on RegisterCodec for details on registration. The lookup of
// content-subtype is case-insensitive. If no such Codec is found, the call
// will result in an error with code codes.Internal.
//
// If ForceCodec is also used, that Codec will be used for all request and
// response messages, with the content-subtype set to the given contentSubtype
// here for requests.
func CallContentSubtype(contentSubtype string) CallOption ***REMOVED***
	return ContentSubtypeCallOption***REMOVED***ContentSubtype: strings.ToLower(contentSubtype)***REMOVED***
***REMOVED***

// ContentSubtypeCallOption is a CallOption that indicates the content-subtype
// used for marshaling messages.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ContentSubtypeCallOption struct ***REMOVED***
	ContentSubtype string
***REMOVED***

func (o ContentSubtypeCallOption) before(c *callInfo) error ***REMOVED***
	c.contentSubtype = o.ContentSubtype
	return nil
***REMOVED***
func (o ContentSubtypeCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// ForceCodec returns a CallOption that will set codec to be used for all
// request and response messages for a call. The result of calling Name() will
// be used as the content-subtype after converting to lowercase, unless
// CallContentSubtype is also used.
//
// See Content-Type on
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details. Also see the documentation on RegisterCodec and
// CallContentSubtype for more details on the interaction between Codec and
// content-subtype.
//
// This function is provided for advanced users; prefer to use only
// CallContentSubtype to select a registered codec instead.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func ForceCodec(codec encoding.Codec) CallOption ***REMOVED***
	return ForceCodecCallOption***REMOVED***Codec: codec***REMOVED***
***REMOVED***

// ForceCodecCallOption is a CallOption that indicates the codec used for
// marshaling messages.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ForceCodecCallOption struct ***REMOVED***
	Codec encoding.Codec
***REMOVED***

func (o ForceCodecCallOption) before(c *callInfo) error ***REMOVED***
	c.codec = o.Codec
	return nil
***REMOVED***
func (o ForceCodecCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// CallCustomCodec behaves like ForceCodec, but accepts a grpc.Codec instead of
// an encoding.Codec.
//
// Deprecated: use ForceCodec instead.
func CallCustomCodec(codec Codec) CallOption ***REMOVED***
	return CustomCodecCallOption***REMOVED***Codec: codec***REMOVED***
***REMOVED***

// CustomCodecCallOption is a CallOption that indicates the codec used for
// marshaling messages.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type CustomCodecCallOption struct ***REMOVED***
	Codec Codec
***REMOVED***

func (o CustomCodecCallOption) before(c *callInfo) error ***REMOVED***
	c.codec = o.Codec
	return nil
***REMOVED***
func (o CustomCodecCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// MaxRetryRPCBufferSize returns a CallOption that limits the amount of memory
// used for buffering this RPC's requests for retry purposes.
//
// Experimental
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func MaxRetryRPCBufferSize(bytes int) CallOption ***REMOVED***
	return MaxRetryRPCBufferSizeCallOption***REMOVED***bytes***REMOVED***
***REMOVED***

// MaxRetryRPCBufferSizeCallOption is a CallOption indicating the amount of
// memory to be used for caching this RPC for retry purposes.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type MaxRetryRPCBufferSizeCallOption struct ***REMOVED***
	MaxRetryRPCBufferSize int
***REMOVED***

func (o MaxRetryRPCBufferSizeCallOption) before(c *callInfo) error ***REMOVED***
	c.maxRetryRPCBufferSize = o.MaxRetryRPCBufferSize
	return nil
***REMOVED***
func (o MaxRetryRPCBufferSizeCallOption) after(c *callInfo, attempt *csAttempt) ***REMOVED******REMOVED***

// The format of the payload: compressed or not?
type payloadFormat uint8

const (
	compressionNone payloadFormat = 0 // no compression
	compressionMade payloadFormat = 1 // compressed
)

// parser reads complete gRPC messages from the underlying reader.
type parser struct ***REMOVED***
	// r is the underlying reader.
	// See the comment on recvMsg for the permissible
	// error types.
	r io.Reader

	// The header of a gRPC message. Find more detail at
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	header [5]byte
***REMOVED***

// recvMsg reads a complete gRPC message from the stream.
//
// It returns the message and its payload (compression/encoding)
// format. The caller owns the returned msg memory.
//
// If there is an error, possible values are:
//   * io.EOF, when no messages remain
//   * io.ErrUnexpectedEOF
//   * of type transport.ConnectionError
//   * an error from the status package
// No other error values or types must be returned, which also means
// that the underlying io.Reader must not return an incompatible
// error.
func (p *parser) recvMsg(maxReceiveMessageSize int) (pf payloadFormat, msg []byte, err error) ***REMOVED***
	if _, err := p.r.Read(p.header[:]); err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***

	pf = payloadFormat(p.header[0])
	length := binary.BigEndian.Uint32(p.header[1:])

	if length == 0 ***REMOVED***
		return pf, nil, nil
	***REMOVED***
	if int64(length) > int64(maxInt) ***REMOVED***
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max length allowed on current machine (%d vs. %d)", length, maxInt)
	***REMOVED***
	if int(length) > maxReceiveMessageSize ***REMOVED***
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", length, maxReceiveMessageSize)
	***REMOVED***
	// TODO(bradfitz,zhaoq): garbage. reuse buffer after proto decoding instead
	// of making it for each message:
	msg = make([]byte, int(length))
	if _, err := p.r.Read(msg); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			err = io.ErrUnexpectedEOF
		***REMOVED***
		return 0, nil, err
	***REMOVED***
	return pf, msg, nil
***REMOVED***

// encode serializes msg and returns a buffer containing the message, or an
// error if it is too large to be transmitted by grpc.  If msg is nil, it
// generates an empty message.
func encode(c baseCodec, msg interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if msg == nil ***REMOVED*** // NOTE: typed nils will not be caught by this check
		return nil, nil
	***REMOVED***
	b, err := c.Marshal(msg)
	if err != nil ***REMOVED***
		return nil, status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
	***REMOVED***
	if uint(len(b)) > math.MaxUint32 ***REMOVED***
		return nil, status.Errorf(codes.ResourceExhausted, "grpc: message too large (%d bytes)", len(b))
	***REMOVED***
	return b, nil
***REMOVED***

// compress returns the input bytes compressed by compressor or cp.  If both
// compressors are nil, returns nil.
//
// TODO(dfawley): eliminate cp parameter by wrapping Compressor in an encoding.Compressor.
func compress(in []byte, cp Compressor, compressor encoding.Compressor) ([]byte, error) ***REMOVED***
	if compressor == nil && cp == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	wrapErr := func(err error) error ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: error while compressing: %v", err.Error())
	***REMOVED***
	cbuf := &bytes.Buffer***REMOVED******REMOVED***
	if compressor != nil ***REMOVED***
		z, err := compressor.Compress(cbuf)
		if err != nil ***REMOVED***
			return nil, wrapErr(err)
		***REMOVED***
		if _, err := z.Write(in); err != nil ***REMOVED***
			return nil, wrapErr(err)
		***REMOVED***
		if err := z.Close(); err != nil ***REMOVED***
			return nil, wrapErr(err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := cp.Do(cbuf, in); err != nil ***REMOVED***
			return nil, wrapErr(err)
		***REMOVED***
	***REMOVED***
	return cbuf.Bytes(), nil
***REMOVED***

const (
	payloadLen = 1
	sizeLen    = 4
	headerLen  = payloadLen + sizeLen
)

// msgHeader returns a 5-byte header for the message being transmitted and the
// payload, which is compData if non-nil or data otherwise.
func msgHeader(data, compData []byte) (hdr []byte, payload []byte) ***REMOVED***
	hdr = make([]byte, headerLen)
	if compData != nil ***REMOVED***
		hdr[0] = byte(compressionMade)
		data = compData
	***REMOVED*** else ***REMOVED***
		hdr[0] = byte(compressionNone)
	***REMOVED***

	// Write length of payload into buf
	binary.BigEndian.PutUint32(hdr[payloadLen:], uint32(len(data)))
	return hdr, data
***REMOVED***

func outPayload(client bool, msg interface***REMOVED******REMOVED***, data, payload []byte, t time.Time) *stats.OutPayload ***REMOVED***
	return &stats.OutPayload***REMOVED***
		Client:     client,
		Payload:    msg,
		Data:       data,
		Length:     len(data),
		WireLength: len(payload) + headerLen,
		SentTime:   t,
	***REMOVED***
***REMOVED***

func checkRecvPayload(pf payloadFormat, recvCompress string, haveCompressor bool) *status.Status ***REMOVED***
	switch pf ***REMOVED***
	case compressionNone:
	case compressionMade:
		if recvCompress == "" || recvCompress == encoding.Identity ***REMOVED***
			return status.New(codes.Internal, "grpc: compressed flag set with identity or empty encoding")
		***REMOVED***
		if !haveCompressor ***REMOVED***
			return status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", recvCompress)
		***REMOVED***
	default:
		return status.Newf(codes.Internal, "grpc: received unexpected payload format %d", pf)
	***REMOVED***
	return nil
***REMOVED***

type payloadInfo struct ***REMOVED***
	wireLength        int // The compressed length got from wire.
	uncompressedBytes []byte
***REMOVED***

func recvAndDecompress(p *parser, s *transport.Stream, dc Decompressor, maxReceiveMessageSize int, payInfo *payloadInfo, compressor encoding.Compressor) ([]byte, error) ***REMOVED***
	pf, d, err := p.recvMsg(maxReceiveMessageSize)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if payInfo != nil ***REMOVED***
		payInfo.wireLength = len(d)
	***REMOVED***

	if st := checkRecvPayload(pf, s.RecvCompress(), compressor != nil || dc != nil); st != nil ***REMOVED***
		return nil, st.Err()
	***REMOVED***

	var size int
	if pf == compressionMade ***REMOVED***
		// To match legacy behavior, if the decompressor is set by WithDecompressor or RPCDecompressor,
		// use this decompressor as the default.
		if dc != nil ***REMOVED***
			d, err = dc.Do(bytes.NewReader(d))
			size = len(d)
		***REMOVED*** else ***REMOVED***
			d, size, err = decompress(compressor, d, maxReceiveMessageSize)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, status.Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
		***REMOVED***
		if size > maxReceiveMessageSize ***REMOVED***
			// TODO: Revisit the error code. Currently keep it consistent with java
			// implementation.
			return nil, status.Errorf(codes.ResourceExhausted, "grpc: received message after decompression larger than max (%d vs. %d)", size, maxReceiveMessageSize)
		***REMOVED***
	***REMOVED***
	return d, nil
***REMOVED***

// Using compressor, decompress d, returning data and size.
// Optionally, if data will be over maxReceiveMessageSize, just return the size.
func decompress(compressor encoding.Compressor, d []byte, maxReceiveMessageSize int) ([]byte, int, error) ***REMOVED***
	dcReader, err := compressor.Decompress(bytes.NewReader(d))
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	if sizer, ok := compressor.(interface ***REMOVED***
		DecompressedSize(compressedBytes []byte) int
	***REMOVED***); ok ***REMOVED***
		if size := sizer.DecompressedSize(d); size >= 0 ***REMOVED***
			if size > maxReceiveMessageSize ***REMOVED***
				return nil, size, nil
			***REMOVED***
			// size is used as an estimate to size the buffer, but we
			// will read more data if available.
			// +MinRead so ReadFrom will not reallocate if size is correct.
			buf := bytes.NewBuffer(make([]byte, 0, size+bytes.MinRead))
			bytesRead, err := buf.ReadFrom(io.LimitReader(dcReader, int64(maxReceiveMessageSize)+1))
			return buf.Bytes(), int(bytesRead), err
		***REMOVED***
	***REMOVED***
	// Read from LimitReader with limit max+1. So if the underlying
	// reader is over limit, the result will be bigger than max.
	d, err = ioutil.ReadAll(io.LimitReader(dcReader, int64(maxReceiveMessageSize)+1))
	return d, len(d), err
***REMOVED***

// For the two compressor parameters, both should not be set, but if they are,
// dc takes precedence over compressor.
// TODO(dfawley): wrap the old compressor/decompressor using the new API?
func recv(p *parser, c baseCodec, s *transport.Stream, dc Decompressor, m interface***REMOVED******REMOVED***, maxReceiveMessageSize int, payInfo *payloadInfo, compressor encoding.Compressor) error ***REMOVED***
	d, err := recvAndDecompress(p, s, dc, maxReceiveMessageSize, payInfo, compressor)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := c.Unmarshal(d, m); err != nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: failed to unmarshal the received message %v", err)
	***REMOVED***
	if payInfo != nil ***REMOVED***
		payInfo.uncompressedBytes = d
	***REMOVED***
	return nil
***REMOVED***

// Information about RPC
type rpcInfo struct ***REMOVED***
	failfast      bool
	preloaderInfo *compressorInfo
***REMOVED***

// Information about Preloader
// Responsible for storing codec, and compressors
// If stream (s) has  context s.Context which stores rpcInfo that has non nil
// pointers to codec, and compressors, then we can use preparedMsg for Async message prep
// and reuse marshalled bytes
type compressorInfo struct ***REMOVED***
	codec baseCodec
	cp    Compressor
	comp  encoding.Compressor
***REMOVED***

type rpcInfoContextKey struct***REMOVED******REMOVED***

func newContextWithRPCInfo(ctx context.Context, failfast bool, codec baseCodec, cp Compressor, comp encoding.Compressor) context.Context ***REMOVED***
	return context.WithValue(ctx, rpcInfoContextKey***REMOVED******REMOVED***, &rpcInfo***REMOVED***
		failfast: failfast,
		preloaderInfo: &compressorInfo***REMOVED***
			codec: codec,
			cp:    cp,
			comp:  comp,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func rpcInfoFromContext(ctx context.Context) (s *rpcInfo, ok bool) ***REMOVED***
	s, ok = ctx.Value(rpcInfoContextKey***REMOVED******REMOVED***).(*rpcInfo)
	return
***REMOVED***

// Code returns the error code for err if it was produced by the rpc system.
// Otherwise, it returns codes.Unknown.
//
// Deprecated: use status.Code instead.
func Code(err error) codes.Code ***REMOVED***
	return status.Code(err)
***REMOVED***

// ErrorDesc returns the error description of err if it was produced by the rpc system.
// Otherwise, it returns err.Error() or empty string when err is nil.
//
// Deprecated: use status.Convert and Message method instead.
func ErrorDesc(err error) string ***REMOVED***
	return status.Convert(err).Message()
***REMOVED***

// Errorf returns an error containing an error code and a description;
// Errorf returns nil if c is OK.
//
// Deprecated: use status.Errorf instead.
func Errorf(c codes.Code, format string, a ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return status.Errorf(c, format, a...)
***REMOVED***

// toRPCErr converts an error into an error from the status package.
func toRPCErr(err error) error ***REMOVED***
	switch err ***REMOVED***
	case nil, io.EOF:
		return err
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case io.ErrUnexpectedEOF:
		return status.Error(codes.Internal, err.Error())
	***REMOVED***

	switch e := err.(type) ***REMOVED***
	case transport.ConnectionError:
		return status.Error(codes.Unavailable, e.Desc)
	case *transport.NewStreamError:
		return toRPCErr(e.Err)
	***REMOVED***

	if _, ok := status.FromError(err); ok ***REMOVED***
		return err
	***REMOVED***

	return status.Error(codes.Unknown, err.Error())
***REMOVED***

// setCallInfoCodec should only be called after CallOptions have been applied.
func setCallInfoCodec(c *callInfo) error ***REMOVED***
	if c.codec != nil ***REMOVED***
		// codec was already set by a CallOption; use it, but set the content
		// subtype if it is not set.
		if c.contentSubtype == "" ***REMOVED***
			// c.codec is a baseCodec to hide the difference between grpc.Codec and
			// encoding.Codec (Name vs. String method name).  We only support
			// setting content subtype from encoding.Codec to avoid a behavior
			// change with the deprecated version.
			if ec, ok := c.codec.(encoding.Codec); ok ***REMOVED***
				c.contentSubtype = strings.ToLower(ec.Name())
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	if c.contentSubtype == "" ***REMOVED***
		// No codec specified in CallOptions; use proto by default.
		c.codec = encoding.GetCodec(proto.Name)
		return nil
	***REMOVED***

	// c.contentSubtype is already lowercased in CallContentSubtype
	c.codec = encoding.GetCodec(c.contentSubtype)
	if c.codec == nil ***REMOVED***
		return status.Errorf(codes.Internal, "no codec registered for content-subtype %s", c.contentSubtype)
	***REMOVED***
	return nil
***REMOVED***

// channelzData is used to store channelz related data for ClientConn, addrConn and Server.
// These fields cannot be embedded in the original structs (e.g. ClientConn), since to do atomic
// operation on int64 variable on 32-bit machine, user is responsible to enforce memory alignment.
// Here, by grouping those int64 fields inside a struct, we are enforcing the alignment.
type channelzData struct ***REMOVED***
	callsStarted   int64
	callsFailed    int64
	callsSucceeded int64
	// lastCallStartedTime stores the timestamp that last call starts. It is of int64 type instead of
	// time.Time since it's more costly to atomically update time.Time variable than int64 variable.
	lastCallStartedTime int64
***REMOVED***

// The SupportPackageIsVersion variables are referenced from generated protocol
// buffer files to ensure compatibility with the gRPC version used.  The latest
// support package version is 7.
//
// Older versions are kept for compatibility.
//
// These constants should not be referenced from any other code.
const (
	SupportPackageIsVersion3 = true
	SupportPackageIsVersion4 = true
	SupportPackageIsVersion5 = true
	SupportPackageIsVersion6 = true
	SupportPackageIsVersion7 = true
)

const grpcUA = "grpc-go/" + Version
