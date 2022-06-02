// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package http2 implements the HTTP/2 protocol.
//
// This package is low-level and intended to be used directly by very
// few people. Most users will use it indirectly through the automatic
// use by the net/http package (from Go 1.6 and later).
// For use in earlier Go versions see ConfigureServer. (Transport support
// requires Go 1.6 or later)
//
// See https://http2.github.io/ for more information on HTTP/2.
//
// See https://http2.golang.org/ for a test server running this code.
package http2 // import "golang.org/x/net/http2"

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/http/httpguts"
)

var (
	VerboseLogs    bool
	logFrameWrites bool
	logFrameReads  bool
	inTests        bool
)

func init() ***REMOVED***
	e := os.Getenv("GODEBUG")
	if strings.Contains(e, "http2debug=1") ***REMOVED***
		VerboseLogs = true
	***REMOVED***
	if strings.Contains(e, "http2debug=2") ***REMOVED***
		VerboseLogs = true
		logFrameWrites = true
		logFrameReads = true
	***REMOVED***
***REMOVED***

const (
	// ClientPreface is the string that must be sent by new
	// connections from clients.
	ClientPreface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	// SETTINGS_MAX_FRAME_SIZE default
	// http://http2.github.io/http2-spec/#rfc.section.6.5.2
	initialMaxFrameSize = 16384

	// NextProtoTLS is the NPN/ALPN protocol negotiated during
	// HTTP/2's TLS setup.
	NextProtoTLS = "h2"

	// http://http2.github.io/http2-spec/#SettingValues
	initialHeaderTableSize = 4096

	initialWindowSize = 65535 // 6.9.2 Initial Flow Control Window Size

	defaultMaxReadFrameSize = 1 << 20
)

var (
	clientPreface = []byte(ClientPreface)
)

type streamState int

// HTTP/2 stream states.
//
// See http://tools.ietf.org/html/rfc7540#section-5.1.
//
// For simplicity, the server code merges "reserved (local)" into
// "half-closed (remote)". This is one less state transition to track.
// The only downside is that we send PUSH_PROMISEs slightly less
// liberally than allowable. More discussion here:
// https://lists.w3.org/Archives/Public/ietf-http-wg/2016JulSep/0599.html
//
// "reserved (remote)" is omitted since the client code does not
// support server push.
const (
	stateIdle streamState = iota
	stateOpen
	stateHalfClosedLocal
	stateHalfClosedRemote
	stateClosed
)

var stateName = [...]string***REMOVED***
	stateIdle:             "Idle",
	stateOpen:             "Open",
	stateHalfClosedLocal:  "HalfClosedLocal",
	stateHalfClosedRemote: "HalfClosedRemote",
	stateClosed:           "Closed",
***REMOVED***

func (st streamState) String() string ***REMOVED***
	return stateName[st]
***REMOVED***

// Setting is a setting parameter: which setting it is, and its value.
type Setting struct ***REMOVED***
	// ID is which setting is being set.
	// See http://http2.github.io/http2-spec/#SettingValues
	ID SettingID

	// Val is the value.
	Val uint32
***REMOVED***

func (s Setting) String() string ***REMOVED***
	return fmt.Sprintf("[%v = %d]", s.ID, s.Val)
***REMOVED***

// Valid reports whether the setting is valid.
func (s Setting) Valid() error ***REMOVED***
	// Limits and error codes from 6.5.2 Defined SETTINGS Parameters
	switch s.ID ***REMOVED***
	case SettingEnablePush:
		if s.Val != 1 && s.Val != 0 ***REMOVED***
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
	case SettingInitialWindowSize:
		if s.Val > 1<<31-1 ***REMOVED***
			return ConnectionError(ErrCodeFlowControl)
		***REMOVED***
	case SettingMaxFrameSize:
		if s.Val < 16384 || s.Val > 1<<24-1 ***REMOVED***
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// A SettingID is an HTTP/2 setting as defined in
// http://http2.github.io/http2-spec/#iana-settings
type SettingID uint16

const (
	SettingHeaderTableSize      SettingID = 0x1
	SettingEnablePush           SettingID = 0x2
	SettingMaxConcurrentStreams SettingID = 0x3
	SettingInitialWindowSize    SettingID = 0x4
	SettingMaxFrameSize         SettingID = 0x5
	SettingMaxHeaderListSize    SettingID = 0x6
)

var settingName = map[SettingID]string***REMOVED***
	SettingHeaderTableSize:      "HEADER_TABLE_SIZE",
	SettingEnablePush:           "ENABLE_PUSH",
	SettingMaxConcurrentStreams: "MAX_CONCURRENT_STREAMS",
	SettingInitialWindowSize:    "INITIAL_WINDOW_SIZE",
	SettingMaxFrameSize:         "MAX_FRAME_SIZE",
	SettingMaxHeaderListSize:    "MAX_HEADER_LIST_SIZE",
***REMOVED***

func (s SettingID) String() string ***REMOVED***
	if v, ok := settingName[s]; ok ***REMOVED***
		return v
	***REMOVED***
	return fmt.Sprintf("UNKNOWN_SETTING_%d", uint16(s))
***REMOVED***

// validWireHeaderFieldName reports whether v is a valid header field
// name (key). See httpguts.ValidHeaderName for the base rules.
//
// Further, http2 says:
//
//	"Just as in HTTP/1.x, header field names are strings of ASCII
//	characters that are compared in a case-insensitive
//	fashion. However, header field names MUST be converted to
//	lowercase prior to their encoding in HTTP/2. "
func validWireHeaderFieldName(v string) bool ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return false
	***REMOVED***
	for _, r := range v ***REMOVED***
		if !httpguts.IsTokenRune(r) ***REMOVED***
			return false
		***REMOVED***
		if 'A' <= r && r <= 'Z' ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func httpCodeString(code int) string ***REMOVED***
	switch code ***REMOVED***
	case 200:
		return "200"
	case 404:
		return "404"
	***REMOVED***
	return strconv.Itoa(code)
***REMOVED***

// from pkg io
type stringWriter interface ***REMOVED***
	WriteString(s string) (n int, err error)
***REMOVED***

// A gate lets two goroutines coordinate their activities.
type gate chan struct***REMOVED******REMOVED***

func (g gate) Done() ***REMOVED*** g <- struct***REMOVED******REMOVED******REMOVED******REMOVED*** ***REMOVED***
func (g gate) Wait() ***REMOVED*** <-g ***REMOVED***

// A closeWaiter is like a sync.WaitGroup but only goes 1 to 0 (open to closed).
type closeWaiter chan struct***REMOVED******REMOVED***

// Init makes a closeWaiter usable.
// It exists because so a closeWaiter value can be placed inside a
// larger struct and have the Mutex and Cond's memory in the same
// allocation.
func (cw *closeWaiter) Init() ***REMOVED***
	*cw = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

// Close marks the closeWaiter as closed and unblocks any waiters.
func (cw closeWaiter) Close() ***REMOVED***
	close(cw)
***REMOVED***

// Wait waits for the closeWaiter to become closed.
func (cw closeWaiter) Wait() ***REMOVED***
	<-cw
***REMOVED***

// bufferedWriter is a buffered writer that writes to w.
// Its buffered writer is lazily allocated as needed, to minimize
// idle memory usage with many connections.
type bufferedWriter struct ***REMOVED***
	_  incomparable
	w  io.Writer     // immutable
	bw *bufio.Writer // non-nil when data is buffered
***REMOVED***

func newBufferedWriter(w io.Writer) *bufferedWriter ***REMOVED***
	return &bufferedWriter***REMOVED***w: w***REMOVED***
***REMOVED***

// bufWriterPoolBufferSize is the size of bufio.Writer's
// buffers created using bufWriterPool.
//
// TODO: pick a less arbitrary value? this is a bit under
// (3 x typical 1500 byte MTU) at least. Other than that,
// not much thought went into it.
const bufWriterPoolBufferSize = 4 << 10

var bufWriterPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return bufio.NewWriterSize(nil, bufWriterPoolBufferSize)
	***REMOVED***,
***REMOVED***

func (w *bufferedWriter) Available() int ***REMOVED***
	if w.bw == nil ***REMOVED***
		return bufWriterPoolBufferSize
	***REMOVED***
	return w.bw.Available()
***REMOVED***

func (w *bufferedWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if w.bw == nil ***REMOVED***
		bw := bufWriterPool.Get().(*bufio.Writer)
		bw.Reset(w.w)
		w.bw = bw
	***REMOVED***
	return w.bw.Write(p)
***REMOVED***

func (w *bufferedWriter) Flush() error ***REMOVED***
	bw := w.bw
	if bw == nil ***REMOVED***
		return nil
	***REMOVED***
	err := bw.Flush()
	bw.Reset(nil)
	bufWriterPool.Put(bw)
	w.bw = nil
	return err
***REMOVED***

func mustUint31(v int32) uint32 ***REMOVED***
	if v < 0 || v > 2147483647 ***REMOVED***
		panic("out of range")
	***REMOVED***
	return uint32(v)
***REMOVED***

// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 7230, section 3.3.
func bodyAllowedForStatus(status int) bool ***REMOVED***
	switch ***REMOVED***
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	***REMOVED***
	return true
***REMOVED***

type httpError struct ***REMOVED***
	_       incomparable
	msg     string
	timeout bool
***REMOVED***

func (e *httpError) Error() string   ***REMOVED*** return e.msg ***REMOVED***
func (e *httpError) Timeout() bool   ***REMOVED*** return e.timeout ***REMOVED***
func (e *httpError) Temporary() bool ***REMOVED*** return true ***REMOVED***

var errTimeout error = &httpError***REMOVED***msg: "http2: timeout awaiting response headers", timeout: true***REMOVED***

type connectionStater interface ***REMOVED***
	ConnectionState() tls.ConnectionState
***REMOVED***

var sorterPool = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(sorter) ***REMOVED******REMOVED***

type sorter struct ***REMOVED***
	v []string // owned by sorter
***REMOVED***

func (s *sorter) Len() int           ***REMOVED*** return len(s.v) ***REMOVED***
func (s *sorter) Swap(i, j int)      ***REMOVED*** s.v[i], s.v[j] = s.v[j], s.v[i] ***REMOVED***
func (s *sorter) Less(i, j int) bool ***REMOVED*** return s.v[i] < s.v[j] ***REMOVED***

// Keys returns the sorted keys of h.
//
// The returned slice is only valid until s used again or returned to
// its pool.
func (s *sorter) Keys(h http.Header) []string ***REMOVED***
	keys := s.v[:0]
	for k := range h ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	s.v = keys
	sort.Sort(s)
	return keys
***REMOVED***

func (s *sorter) SortStrings(ss []string) ***REMOVED***
	// Our sorter works on s.v, which sorter owns, so
	// stash it away while we sort the user's buffer.
	save := s.v
	s.v = ss
	sort.Sort(s)
	s.v = save
***REMOVED***

// validPseudoPath reports whether v is a valid :path pseudo-header
// value. It must be either:
//
//   - a non-empty string starting with '/'
//   - the string '*', for OPTIONS requests.
//
// For now this is only used a quick check for deciding when to clean
// up Opaque URLs before sending requests from the Transport.
// See golang.org/issue/16847
//
// We used to enforce that the path also didn't start with "//", but
// Google's GFE accepts such paths and Chrome sends them, so ignore
// that part of the spec. See golang.org/issue/19103.
func validPseudoPath(v string) bool ***REMOVED***
	return (len(v) > 0 && v[0] == '/') || v == "*"
***REMOVED***

// incomparable is a zero-width, non-comparable type. Adding it to a struct
// makes that struct also non-comparable, and generally doesn't add
// any size (as long as it's first).
type incomparable [0]func()
