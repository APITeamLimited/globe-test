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

func decodeGRPCStatusDetails(rawDetails string) (*status.Status, error) ***REMOVED***
	v, err := decodeBinHeader(rawDetails)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	st := &spb.Status***REMOVED******REMOVED***
	if err = proto.Unmarshal(v, st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return status.FromProto(st), nil
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
