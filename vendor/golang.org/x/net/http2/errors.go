// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"errors"
	"fmt"
)

// An ErrCode is an unsigned 32-bit error code as defined in the HTTP/2 spec.
type ErrCode uint32

const (
	ErrCodeNo                 ErrCode = 0x0
	ErrCodeProtocol           ErrCode = 0x1
	ErrCodeInternal           ErrCode = 0x2
	ErrCodeFlowControl        ErrCode = 0x3
	ErrCodeSettingsTimeout    ErrCode = 0x4
	ErrCodeStreamClosed       ErrCode = 0x5
	ErrCodeFrameSize          ErrCode = 0x6
	ErrCodeRefusedStream      ErrCode = 0x7
	ErrCodeCancel             ErrCode = 0x8
	ErrCodeCompression        ErrCode = 0x9
	ErrCodeConnect            ErrCode = 0xa
	ErrCodeEnhanceYourCalm    ErrCode = 0xb
	ErrCodeInadequateSecurity ErrCode = 0xc
	ErrCodeHTTP11Required     ErrCode = 0xd
)

var errCodeName = map[ErrCode]string***REMOVED***
	ErrCodeNo:                 "NO_ERROR",
	ErrCodeProtocol:           "PROTOCOL_ERROR",
	ErrCodeInternal:           "INTERNAL_ERROR",
	ErrCodeFlowControl:        "FLOW_CONTROL_ERROR",
	ErrCodeSettingsTimeout:    "SETTINGS_TIMEOUT",
	ErrCodeStreamClosed:       "STREAM_CLOSED",
	ErrCodeFrameSize:          "FRAME_SIZE_ERROR",
	ErrCodeRefusedStream:      "REFUSED_STREAM",
	ErrCodeCancel:             "CANCEL",
	ErrCodeCompression:        "COMPRESSION_ERROR",
	ErrCodeConnect:            "CONNECT_ERROR",
	ErrCodeEnhanceYourCalm:    "ENHANCE_YOUR_CALM",
	ErrCodeInadequateSecurity: "INADEQUATE_SECURITY",
	ErrCodeHTTP11Required:     "HTTP_1_1_REQUIRED",
***REMOVED***

func (e ErrCode) String() string ***REMOVED***
	if s, ok := errCodeName[e]; ok ***REMOVED***
		return s
	***REMOVED***
	return fmt.Sprintf("unknown error code 0x%x", uint32(e))
***REMOVED***

func (e ErrCode) stringToken() string ***REMOVED***
	if s, ok := errCodeName[e]; ok ***REMOVED***
		return s
	***REMOVED***
	return fmt.Sprintf("ERR_UNKNOWN_%d", uint32(e))
***REMOVED***

// ConnectionError is an error that results in the termination of the
// entire connection.
type ConnectionError ErrCode

func (e ConnectionError) Error() string ***REMOVED*** return fmt.Sprintf("connection error: %s", ErrCode(e)) ***REMOVED***

// StreamError is an error that only affects one stream within an
// HTTP/2 connection.
type StreamError struct ***REMOVED***
	StreamID uint32
	Code     ErrCode
	Cause    error // optional additional detail
***REMOVED***

// errFromPeer is a sentinel error value for StreamError.Cause to
// indicate that the StreamError was sent from the peer over the wire
// and wasn't locally generated in the Transport.
var errFromPeer = errors.New("received from peer")

func streamError(id uint32, code ErrCode) StreamError ***REMOVED***
	return StreamError***REMOVED***StreamID: id, Code: code***REMOVED***
***REMOVED***

func (e StreamError) Error() string ***REMOVED***
	if e.Cause != nil ***REMOVED***
		return fmt.Sprintf("stream error: stream ID %d; %v; %v", e.StreamID, e.Code, e.Cause)
	***REMOVED***
	return fmt.Sprintf("stream error: stream ID %d; %v", e.StreamID, e.Code)
***REMOVED***

// 6.9.1 The Flow Control Window
// "If a sender receives a WINDOW_UPDATE that causes a flow control
// window to exceed this maximum it MUST terminate either the stream
// or the connection, as appropriate. For streams, [...]; for the
// connection, a GOAWAY frame with a FLOW_CONTROL_ERROR code."
type goAwayFlowError struct***REMOVED******REMOVED***

func (goAwayFlowError) Error() string ***REMOVED*** return "connection exceeded flow control window size" ***REMOVED***

// connError represents an HTTP/2 ConnectionError error code, along
// with a string (for debugging) explaining why.
//
// Errors of this type are only returned by the frame parser functions
// and converted into ConnectionError(Code), after stashing away
// the Reason into the Framer's errDetail field, accessible via
// the (*Framer).ErrorDetail method.
type connError struct ***REMOVED***
	Code   ErrCode // the ConnectionError error code
	Reason string  // additional reason
***REMOVED***

func (e connError) Error() string ***REMOVED***
	return fmt.Sprintf("http2: connection error: %v: %v", e.Code, e.Reason)
***REMOVED***

type pseudoHeaderError string

func (e pseudoHeaderError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid pseudo-header %q", string(e))
***REMOVED***

type duplicatePseudoHeaderError string

func (e duplicatePseudoHeaderError) Error() string ***REMOVED***
	return fmt.Sprintf("duplicate pseudo-header %q", string(e))
***REMOVED***

type headerFieldNameError string

func (e headerFieldNameError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid header field name %q", string(e))
***REMOVED***

type headerFieldValueError string

func (e headerFieldValueError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid header field value %q", string(e))
***REMOVED***

var (
	errMixPseudoHeaderTypes = errors.New("mix of request and response pseudo headers")
	errPseudoAfterRegular   = errors.New("pseudo header field after regular")
)
