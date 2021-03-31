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

package transport

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/status"
)

const (
	// http2MaxFrameLen specifies the max length of a HTTP2 frame.
	http2MaxFrameLen = 16384 // 16KB frame
	// http://http2.github.io/http2-spec/#SettingValues
	http2InitHeaderTableSize = 4096
	// baseContentType is the base content-type for gRPC.  This is a valid
	// content-type on it's own, but can also include a content-subtype such as
	// "proto" as a suffix after "+" or ";".  See
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
	// for more details.

)

var (
	clientPreface   = []byte(http2.ClientPreface)
	http2ErrConvTab = map[http2.ErrCode]codes.Code***REMOVED***
		http2.ErrCodeNo:                 codes.Internal,
		http2.ErrCodeProtocol:           codes.Internal,
		http2.ErrCodeInternal:           codes.Internal,
		http2.ErrCodeFlowControl:        codes.ResourceExhausted,
		http2.ErrCodeSettingsTimeout:    codes.Internal,
		http2.ErrCodeStreamClosed:       codes.Internal,
		http2.ErrCodeFrameSize:          codes.Internal,
		http2.ErrCodeRefusedStream:      codes.Unavailable,
		http2.ErrCodeCancel:             codes.Canceled,
		http2.ErrCodeCompression:        codes.Internal,
		http2.ErrCodeConnect:            codes.Internal,
		http2.ErrCodeEnhanceYourCalm:    codes.ResourceExhausted,
		http2.ErrCodeInadequateSecurity: codes.PermissionDenied,
		http2.ErrCodeHTTP11Required:     codes.Internal,
	***REMOVED***
	// HTTPStatusConvTab is the HTTP status code to gRPC error code conversion table.
	HTTPStatusConvTab = map[int]codes.Code***REMOVED***
		// 400 Bad Request - INTERNAL.
		http.StatusBadRequest: codes.Internal,
		// 401 Unauthorized  - UNAUTHENTICATED.
		http.StatusUnauthorized: codes.Unauthenticated,
		// 403 Forbidden - PERMISSION_DENIED.
		http.StatusForbidden: codes.PermissionDenied,
		// 404 Not Found - UNIMPLEMENTED.
		http.StatusNotFound: codes.Unimplemented,
		// 429 Too Many Requests - UNAVAILABLE.
		http.StatusTooManyRequests: codes.Unavailable,
		// 502 Bad Gateway - UNAVAILABLE.
		http.StatusBadGateway: codes.Unavailable,
		// 503 Service Unavailable - UNAVAILABLE.
		http.StatusServiceUnavailable: codes.Unavailable,
		// 504 Gateway timeout - UNAVAILABLE.
		http.StatusGatewayTimeout: codes.Unavailable,
	***REMOVED***
	logger = grpclog.Component("transport")
)

type parsedHeaderData struct ***REMOVED***
	encoding string
	// statusGen caches the stream status received from the trailer the server
	// sent.  Client side only.  Do not access directly.  After all trailers are
	// parsed, use the status method to retrieve the status.
	statusGen *status.Status
	// rawStatusCode and rawStatusMsg are set from the raw trailer fields and are not
	// intended for direct access outside of parsing.
	rawStatusCode *int
	rawStatusMsg  string
	httpStatus    *int
	// Server side only fields.
	timeoutSet bool
	timeout    time.Duration
	method     string
	// key-value metadata map from the peer.
	mdata          map[string][]string
	statsTags      []byte
	statsTrace     []byte
	contentSubtype string

	// isGRPC field indicates whether the peer is speaking gRPC (otherwise HTTP).
	//
	// We are in gRPC mode (peer speaking gRPC) if:
	// 	* We are client side and have already received a HEADER frame that indicates gRPC peer.
	//  * The header contains valid  a content-type, i.e. a string starts with "application/grpc"
	// And we should handle error specific to gRPC.
	//
	// Otherwise (i.e. a content-type string starts without "application/grpc", or does not exist), we
	// are in HTTP fallback mode, and should handle error specific to HTTP.
	isGRPC         bool
	grpcErr        error
	httpErr        error
	contentTypeErr string
***REMOVED***

// decodeState configures decoding criteria and records the decoded data.
type decodeState struct ***REMOVED***
	// whether decoding on server side or not
	serverSide bool

	// Records the states during HPACK decoding. It will be filled with info parsed from HTTP HEADERS
	// frame once decodeHeader function has been invoked and returned.
	data parsedHeaderData
***REMOVED***

// isReservedHeader checks whether hdr belongs to HTTP2 headers
// reserved by gRPC protocol. Any other headers are classified as the
// user-specified metadata.
func isReservedHeader(hdr string) bool ***REMOVED***
	if hdr != "" && hdr[0] == ':' ***REMOVED***
		return true
	***REMOVED***
	switch hdr ***REMOVED***
	case "content-type",
		"user-agent",
		"grpc-message-type",
		"grpc-encoding",
		"grpc-message",
		"grpc-status",
		"grpc-timeout",
		"grpc-status-details-bin",
		// Intentionally exclude grpc-previous-rpc-attempts and
		// grpc-retry-pushback-ms, which are "reserved", but their API
		// intentionally works via metadata.
		"te":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// isWhitelistedHeader checks whether hdr should be propagated into metadata
// visible to users, even though it is classified as "reserved", above.
func isWhitelistedHeader(hdr string) bool ***REMOVED***
	switch hdr ***REMOVED***
	case ":authority", "user-agent":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (d *decodeState) status() *status.Status ***REMOVED***
	if d.data.statusGen == nil ***REMOVED***
		// No status-details were provided; generate status using code/msg.
		d.data.statusGen = status.New(codes.Code(int32(*(d.data.rawStatusCode))), d.data.rawStatusMsg)
	***REMOVED***
	return d.data.statusGen
***REMOVED***

const binHdrSuffix = "-bin"

func encodeBinHeader(v []byte) string ***REMOVED***
	return base64.RawStdEncoding.EncodeToString(v)
***REMOVED***

func decodeBinHeader(v string) ([]byte, error) ***REMOVED***
	if len(v)%4 == 0 ***REMOVED***
		// Input was padded, or padding was not necessary.
		return base64.StdEncoding.DecodeString(v)
	***REMOVED***
	return base64.RawStdEncoding.DecodeString(v)
***REMOVED***

func encodeMetadataHeader(k, v string) string ***REMOVED***
	if strings.HasSuffix(k, binHdrSuffix) ***REMOVED***
		return encodeBinHeader(([]byte)(v))
	***REMOVED***
	return v
***REMOVED***

func decodeMetadataHeader(k, v string) (string, error) ***REMOVED***
	if strings.HasSuffix(k, binHdrSuffix) ***REMOVED***
		b, err := decodeBinHeader(v)
		return string(b), err
	***REMOVED***
	return v, nil
***REMOVED***

func (d *decodeState) decodeHeader(frame *http2.MetaHeadersFrame) (http2.ErrCode, error) ***REMOVED***
	// frame.Truncated is set to true when framer detects that the current header
	// list size hits MaxHeaderListSize limit.
	if frame.Truncated ***REMOVED***
		return http2.ErrCodeFrameSize, status.Error(codes.Internal, "peer header list size exceeded limit")
	***REMOVED***

	for _, hf := range frame.Fields ***REMOVED***
		d.processHeaderField(hf)
	***REMOVED***

	if d.data.isGRPC ***REMOVED***
		if d.data.grpcErr != nil ***REMOVED***
			return http2.ErrCodeProtocol, d.data.grpcErr
		***REMOVED***
		if d.serverSide ***REMOVED***
			return http2.ErrCodeNo, nil
		***REMOVED***
		if d.data.rawStatusCode == nil && d.data.statusGen == nil ***REMOVED***
			// gRPC status doesn't exist.
			// Set rawStatusCode to be unknown and return nil error.
			// So that, if the stream has ended this Unknown status
			// will be propagated to the user.
			// Otherwise, it will be ignored. In which case, status from
			// a later trailer, that has StreamEnded flag set, is propagated.
			code := int(codes.Unknown)
			d.data.rawStatusCode = &code
		***REMOVED***
		return http2.ErrCodeNo, nil
	***REMOVED***

	// HTTP fallback mode
	if d.data.httpErr != nil ***REMOVED***
		return http2.ErrCodeProtocol, d.data.httpErr
	***REMOVED***

	var (
		code = codes.Internal // when header does not include HTTP status, return INTERNAL
		ok   bool
	)

	if d.data.httpStatus != nil ***REMOVED***
		code, ok = HTTPStatusConvTab[*(d.data.httpStatus)]
		if !ok ***REMOVED***
			code = codes.Unknown
		***REMOVED***
	***REMOVED***

	return http2.ErrCodeProtocol, status.Error(code, d.constructHTTPErrMsg())
***REMOVED***

// constructErrMsg constructs error message to be returned in HTTP fallback mode.
// Format: HTTP status code and its corresponding message + content-type error message.
func (d *decodeState) constructHTTPErrMsg() string ***REMOVED***
	var errMsgs []string

	if d.data.httpStatus == nil ***REMOVED***
		errMsgs = append(errMsgs, "malformed header: missing HTTP status")
	***REMOVED*** else ***REMOVED***
		errMsgs = append(errMsgs, fmt.Sprintf("%s: HTTP status code %d", http.StatusText(*(d.data.httpStatus)), *d.data.httpStatus))
	***REMOVED***

	if d.data.contentTypeErr == "" ***REMOVED***
		errMsgs = append(errMsgs, "transport: missing content-type field")
	***REMOVED*** else ***REMOVED***
		errMsgs = append(errMsgs, d.data.contentTypeErr)
	***REMOVED***

	return strings.Join(errMsgs, "; ")
***REMOVED***

func (d *decodeState) addMetadata(k, v string) ***REMOVED***
	if d.data.mdata == nil ***REMOVED***
		d.data.mdata = make(map[string][]string)
	***REMOVED***
	d.data.mdata[k] = append(d.data.mdata[k], v)
***REMOVED***

func (d *decodeState) processHeaderField(f hpack.HeaderField) ***REMOVED***
	switch f.Name ***REMOVED***
	case "content-type":
		contentSubtype, validContentType := grpcutil.ContentSubtype(f.Value)
		if !validContentType ***REMOVED***
			d.data.contentTypeErr = fmt.Sprintf("transport: received the unexpected content-type %q", f.Value)
			return
		***REMOVED***
		d.data.contentSubtype = contentSubtype
		// TODO: do we want to propagate the whole content-type in the metadata,
		// or come up with a way to just propagate the content-subtype if it was set?
		// ie ***REMOVED***"content-type": "application/grpc+proto"***REMOVED*** or ***REMOVED***"content-subtype": "proto"***REMOVED***
		// in the metadata?
		d.addMetadata(f.Name, f.Value)
		d.data.isGRPC = true
	case "grpc-encoding":
		d.data.encoding = f.Value
	case "grpc-status":
		code, err := strconv.Atoi(f.Value)
		if err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed grpc-status: %v", err)
			return
		***REMOVED***
		d.data.rawStatusCode = &code
	case "grpc-message":
		d.data.rawStatusMsg = decodeGrpcMessage(f.Value)
	case "grpc-status-details-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
			return
		***REMOVED***
		s := &spb.Status***REMOVED******REMOVED***
		if err := proto.Unmarshal(v, s); err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
			return
		***REMOVED***
		d.data.statusGen = status.FromProto(s)
	case "grpc-timeout":
		d.data.timeoutSet = true
		var err error
		if d.data.timeout, err = decodeTimeout(f.Value); err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed time-out: %v", err)
		***REMOVED***
	case ":path":
		d.data.method = f.Value
	case ":status":
		code, err := strconv.Atoi(f.Value)
		if err != nil ***REMOVED***
			d.data.httpErr = status.Errorf(codes.Internal, "transport: malformed http-status: %v", err)
			return
		***REMOVED***
		d.data.httpStatus = &code
	case "grpc-tags-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed grpc-tags-bin: %v", err)
			return
		***REMOVED***
		d.data.statsTags = v
		d.addMetadata(f.Name, string(v))
	case "grpc-trace-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil ***REMOVED***
			d.data.grpcErr = status.Errorf(codes.Internal, "transport: malformed grpc-trace-bin: %v", err)
			return
		***REMOVED***
		d.data.statsTrace = v
		d.addMetadata(f.Name, string(v))
	default:
		if isReservedHeader(f.Name) && !isWhitelistedHeader(f.Name) ***REMOVED***
			break
		***REMOVED***
		v, err := decodeMetadataHeader(f.Name, f.Value)
		if err != nil ***REMOVED***
			if logger.V(logLevel) ***REMOVED***
				logger.Errorf("Failed to decode metadata header (%q, %q): %v", f.Name, f.Value, err)
			***REMOVED***
			return
		***REMOVED***
		d.addMetadata(f.Name, v)
	***REMOVED***
***REMOVED***

type timeoutUnit uint8

const (
	hour        timeoutUnit = 'H'
	minute      timeoutUnit = 'M'
	second      timeoutUnit = 'S'
	millisecond timeoutUnit = 'm'
	microsecond timeoutUnit = 'u'
	nanosecond  timeoutUnit = 'n'
)

func timeoutUnitToDuration(u timeoutUnit) (d time.Duration, ok bool) ***REMOVED***
	switch u ***REMOVED***
	case hour:
		return time.Hour, true
	case minute:
		return time.Minute, true
	case second:
		return time.Second, true
	case millisecond:
		return time.Millisecond, true
	case microsecond:
		return time.Microsecond, true
	case nanosecond:
		return time.Nanosecond, true
	default:
	***REMOVED***
	return
***REMOVED***

func decodeTimeout(s string) (time.Duration, error) ***REMOVED***
	size := len(s)
	if size < 2 ***REMOVED***
		return 0, fmt.Errorf("transport: timeout string is too short: %q", s)
	***REMOVED***
	if size > 9 ***REMOVED***
		// Spec allows for 8 digits plus the unit.
		return 0, fmt.Errorf("transport: timeout string is too long: %q", s)
	***REMOVED***
	unit := timeoutUnit(s[size-1])
	d, ok := timeoutUnitToDuration(unit)
	if !ok ***REMOVED***
		return 0, fmt.Errorf("transport: timeout unit is not recognized: %q", s)
	***REMOVED***
	t, err := strconv.ParseInt(s[:size-1], 10, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	const maxHours = math.MaxInt64 / int64(time.Hour)
	if d == time.Hour && t > maxHours ***REMOVED***
		// This timeout would overflow math.MaxInt64; clamp it.
		return time.Duration(math.MaxInt64), nil
	***REMOVED***
	return d * time.Duration(t), nil
***REMOVED***

const (
	spaceByte   = ' '
	tildeByte   = '~'
	percentByte = '%'
)

// encodeGrpcMessage is used to encode status code in header field
// "grpc-message". It does percent encoding and also replaces invalid utf-8
// characters with Unicode replacement character.
//
// It checks to see if each individual byte in msg is an allowable byte, and
// then either percent encoding or passing it through. When percent encoding,
// the byte is converted into hexadecimal notation with a '%' prepended.
func encodeGrpcMessage(msg string) string ***REMOVED***
	if msg == "" ***REMOVED***
		return ""
	***REMOVED***
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		c := msg[i]
		if !(c >= spaceByte && c <= tildeByte && c != percentByte) ***REMOVED***
			return encodeGrpcMessageUnchecked(msg)
		***REMOVED***
	***REMOVED***
	return msg
***REMOVED***

func encodeGrpcMessageUnchecked(msg string) string ***REMOVED***
	var buf bytes.Buffer
	for len(msg) > 0 ***REMOVED***
		r, size := utf8.DecodeRuneInString(msg)
		for _, b := range []byte(string(r)) ***REMOVED***
			if size > 1 ***REMOVED***
				// If size > 1, r is not ascii. Always do percent encoding.
				buf.WriteString(fmt.Sprintf("%%%02X", b))
				continue
			***REMOVED***

			// The for loop is necessary even if size == 1. r could be
			// utf8.RuneError.
			//
			// fmt.Sprintf("%%%02X", utf8.RuneError) gives "%FFFD".
			if b >= spaceByte && b <= tildeByte && b != percentByte ***REMOVED***
				buf.WriteByte(b)
			***REMOVED*** else ***REMOVED***
				buf.WriteString(fmt.Sprintf("%%%02X", b))
			***REMOVED***
		***REMOVED***
		msg = msg[size:]
	***REMOVED***
	return buf.String()
***REMOVED***

// decodeGrpcMessage decodes the msg encoded by encodeGrpcMessage.
func decodeGrpcMessage(msg string) string ***REMOVED***
	if msg == "" ***REMOVED***
		return ""
	***REMOVED***
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		if msg[i] == percentByte && i+2 < lenMsg ***REMOVED***
			return decodeGrpcMessageUnchecked(msg)
		***REMOVED***
	***REMOVED***
	return msg
***REMOVED***

func decodeGrpcMessageUnchecked(msg string) string ***REMOVED***
	var buf bytes.Buffer
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		c := msg[i]
		if c == percentByte && i+2 < lenMsg ***REMOVED***
			parsed, err := strconv.ParseUint(msg[i+1:i+3], 16, 8)
			if err != nil ***REMOVED***
				buf.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				buf.WriteByte(byte(parsed))
				i += 2
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(c)
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

type bufWriter struct ***REMOVED***
	buf       []byte
	offset    int
	batchSize int
	conn      net.Conn
	err       error

	onFlush func()
***REMOVED***

func newBufWriter(conn net.Conn, batchSize int) *bufWriter ***REMOVED***
	return &bufWriter***REMOVED***
		buf:       make([]byte, batchSize*2),
		batchSize: batchSize,
		conn:      conn,
	***REMOVED***
***REMOVED***

func (w *bufWriter) Write(b []byte) (n int, err error) ***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***
	if w.batchSize == 0 ***REMOVED*** // Buffer has been disabled.
		return w.conn.Write(b)
	***REMOVED***
	for len(b) > 0 ***REMOVED***
		nn := copy(w.buf[w.offset:], b)
		b = b[nn:]
		w.offset += nn
		n += nn
		if w.offset >= w.batchSize ***REMOVED***
			err = w.Flush()
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

func (w *bufWriter) Flush() error ***REMOVED***
	if w.err != nil ***REMOVED***
		return w.err
	***REMOVED***
	if w.offset == 0 ***REMOVED***
		return nil
	***REMOVED***
	if w.onFlush != nil ***REMOVED***
		w.onFlush()
	***REMOVED***
	_, w.err = w.conn.Write(w.buf[:w.offset])
	w.offset = 0
	return w.err
***REMOVED***

type framer struct ***REMOVED***
	writer *bufWriter
	fr     *http2.Framer
***REMOVED***

func newFramer(conn net.Conn, writeBufferSize, readBufferSize int, maxHeaderListSize uint32) *framer ***REMOVED***
	if writeBufferSize < 0 ***REMOVED***
		writeBufferSize = 0
	***REMOVED***
	var r io.Reader = conn
	if readBufferSize > 0 ***REMOVED***
		r = bufio.NewReaderSize(r, readBufferSize)
	***REMOVED***
	w := newBufWriter(conn, writeBufferSize)
	f := &framer***REMOVED***
		writer: w,
		fr:     http2.NewFramer(w, r),
	***REMOVED***
	f.fr.SetMaxReadFrameSize(http2MaxFrameLen)
	// Opt-in to Frame reuse API on framer to reduce garbage.
	// Frames aren't safe to read from after a subsequent call to ReadFrame.
	f.fr.SetReuseFrames()
	f.fr.MaxHeaderListSize = maxHeaderListSize
	f.fr.ReadMetaHeaders = hpack.NewDecoder(http2InitHeaderTableSize, nil)
	return f
***REMOVED***

// parseDialTarget returns the network and address to pass to dialer.
func parseDialTarget(target string) (string, string) ***REMOVED***
	net := "tcp"
	m1 := strings.Index(target, ":")
	m2 := strings.Index(target, ":/")
	// handle unix:addr which will fail with url.Parse
	if m1 >= 0 && m2 < 0 ***REMOVED***
		if n := target[0:m1]; n == "unix" ***REMOVED***
			return n, target[m1+1:]
		***REMOVED***
	***REMOVED***
	if m2 >= 0 ***REMOVED***
		t, err := url.Parse(target)
		if err != nil ***REMOVED***
			return net, target
		***REMOVED***
		scheme := t.Scheme
		addr := t.Path
		if scheme == "unix" ***REMOVED***
			if addr == "" ***REMOVED***
				addr = t.Host
			***REMOVED***
			return scheme, addr
		***REMOVED***
	***REMOVED***
	return net, target
***REMOVED***
