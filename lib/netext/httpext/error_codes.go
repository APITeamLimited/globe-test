/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package httpext

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"os"
	"runtime"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"

	"github.com/loadimpact/k6/lib/netext"
)

// TODO: maybe rename the type errorCode, so we can have errCode variables? and
// also the constants would probably be better of if `ErrorCode` was a prefix,
// not a suffix - they would be much easier for auto-autocompletion at least...

type errCode uint32

const (
	// non specific
	defaultErrorCode          errCode = 1000
	defaultNetNonTCPErrorCode errCode = 1010
	// DNS errors
	defaultDNSErrorCode    errCode = 1100
	dnsNoSuchHostErrorCode errCode = 1101
	blackListedIPErrorCode errCode = 1110
	// tcp errors
	defaultTCPErrorCode      errCode = 1200
	tcpBrokenPipeErrorCode   errCode = 1201
	netUnknownErrnoErrorCode errCode = 1202
	tcpDialErrorCode         errCode = 1210
	tcpDialTimeoutErrorCode  errCode = 1211
	tcpDialRefusedErrorCode  errCode = 1212
	tcpDialUnknownErrnoCode  errCode = 1213
	tcpResetByPeerErrorCode  errCode = 1220
	// TLS errors
	defaultTLSErrorCode           errCode = 1300
	x509UnknownAuthorityErrorCode errCode = 1310
	x509HostnameErrorCode         errCode = 1311

	// HTTP2 errors
	// defaultHTTP2ErrorCode errCode = 1600 // commented because of golint
	// HTTP2 GoAway errors
	unknownHTTP2GoAwayErrorCode errCode = 1610
	// errors till 1611 + 13 are other HTTP2 GoAway errors with a specific errCode

	// HTTP2 Stream errors
	unknownHTTP2StreamErrorCode errCode = 1630
	// errors till 1631 + 13 are other HTTP2 Stream errors with a specific errCode

	// HTTP2 Connection errors
	unknownHTTP2ConnectionErrorCode errCode = 1650
	// errors till 1651 + 13 are other HTTP2 Connection errors with a specific errCode

	// Custom k6 content errors, i.e. when the magic fails
	//defaultContentError errCode = 1700 // reserved for future use
	responseDecompressionErrorCode errCode = 1701
)

const (
	tcpResetByPeerErrorCodeMsg  = "write: connection reset by peer"
	tcpDialTimeoutErrorCodeMsg  = "dial: i/o timeout"
	tcpDialRefusedErrorCodeMsg  = "dial: connection refused"
	tcpBrokenPipeErrorCodeMsg   = "write: broken pipe"
	netUnknownErrnoErrorCodeMsg = "%s: unknown errno `%d` on %s with message `%s`"
	dnsNoSuchHostErrorCodeMsg   = "lookup: no such host"
	blackListedIPErrorCodeMsg   = "ip is blacklisted"
	http2GoAwayErrorCodeMsg     = "http2: received GoAway with http2 ErrCode %s"
	http2StreamErrorCodeMsg     = "http2: stream error with http2 ErrCode %s"
	http2ConnectionErrorCodeMsg = "http2: connection error with http2 ErrCode %s"
	x509HostnameErrorCodeMsg    = "x509: certificate doesn't match hostname"
	x509UnknownAuthority        = "x509: unknown authority"
)

func http2ErrCodeOffset(code http2.ErrCode) errCode ***REMOVED***
	if code > http2.ErrCodeHTTP11Required ***REMOVED***
		return 0
	***REMOVED***
	return 1 + errCode(code)
***REMOVED***

// errorCodeForError returns the errorCode and a specific error message for given error.
func errorCodeForError(err error) (errCode, string) ***REMOVED***
	switch e := errors.Cause(err).(type) ***REMOVED***
	case K6Error:
		return e.Code, e.Message
	case *net.DNSError:
		switch e.Err ***REMOVED***
		case "no such host": // defined as private in the go stdlib
			return dnsNoSuchHostErrorCode, dnsNoSuchHostErrorCodeMsg
		default:
			return defaultDNSErrorCode, err.Error()
		***REMOVED***
	case netext.BlackListedIPError:
		return blackListedIPErrorCode, blackListedIPErrorCodeMsg
	case *http2.GoAwayError:
		return unknownHTTP2GoAwayErrorCode + http2ErrCodeOffset(e.ErrCode),
			fmt.Sprintf(http2GoAwayErrorCodeMsg, e.ErrCode)
	case *http2.StreamError:
		return unknownHTTP2StreamErrorCode + http2ErrCodeOffset(e.Code),
			fmt.Sprintf(http2StreamErrorCodeMsg, e.Code)
	case *http2.ConnectionError:
		return unknownHTTP2ConnectionErrorCode + http2ErrCodeOffset(http2.ErrCode(*e)),
			fmt.Sprintf(http2ConnectionErrorCodeMsg, http2.ErrCode(*e))
	case *net.OpError:
		if e.Net != "tcp" && e.Net != "tcp6" ***REMOVED***
			// TODO: figure out how this happens
			return defaultNetNonTCPErrorCode, err.Error()
		***REMOVED***
		if e.Op == "write" ***REMOVED***
			if sErr, ok := e.Err.(*os.SyscallError); ok ***REMOVED***
				switch sErr.Err ***REMOVED***
				case syscall.ECONNRESET:
					return tcpResetByPeerErrorCode, tcpResetByPeerErrorCodeMsg
				case syscall.EPIPE:
					return tcpBrokenPipeErrorCode, tcpBrokenPipeErrorCodeMsg
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if e.Op == "dial" ***REMOVED***
			if e.Timeout() ***REMOVED***
				return tcpDialTimeoutErrorCode, tcpDialTimeoutErrorCodeMsg
			***REMOVED***
			if iErr, ok := e.Err.(*os.SyscallError); ok ***REMOVED***
				if errno, ok := iErr.Err.(syscall.Errno); ok ***REMOVED***
					if errno == syscall.ECONNREFUSED ||
						// 10061 is some connection refused like thing on windows
						// TODO: fix by moving to x/sys instead of syscall after
						// https://github.com/golang/go/issues/31360 gets resolved
						(errno == 10061 && runtime.GOOS == "windows") ***REMOVED***
						return tcpDialRefusedErrorCode, tcpDialRefusedErrorCodeMsg
					***REMOVED***
					return tcpDialUnknownErrnoCode,
						fmt.Sprintf("dial: unknown errno %d error with msg `%s`", errno, iErr.Err)
				***REMOVED***
			***REMOVED***
			return tcpDialErrorCode, err.Error()
		***REMOVED***
		switch inErr := e.Err.(type) ***REMOVED***
		case syscall.Errno:
			return netUnknownErrnoErrorCode,
				fmt.Sprintf(netUnknownErrnoErrorCodeMsg,
					e.Op, (int)(inErr), runtime.GOOS, inErr.Error())
		default:
			return defaultTCPErrorCode, err.Error()
		***REMOVED***

	case *x509.UnknownAuthorityError:
		return x509UnknownAuthorityErrorCode, x509UnknownAuthority
	case *x509.HostnameError:
		return x509HostnameErrorCode, x509HostnameErrorCodeMsg
	case *tls.RecordHeaderError:
		return defaultTLSErrorCode, err.Error()
	case *url.Error:
		return errorCodeForError(e.Err)
	default:
		return defaultErrorCode, err.Error()
	***REMOVED***
***REMOVED***

// K6Error is a helper struct that enhances Go errors with custom k6-specific
// error-codes and more user-readable error messages.
type K6Error struct ***REMOVED***
	Code          errCode
	Message       string
	OriginalError error
***REMOVED***

// NewK6Error is the constructor for K6Error
func NewK6Error(code errCode, msg string, originalErr error) K6Error ***REMOVED***
	return K6Error***REMOVED***code, msg, originalErr***REMOVED***
***REMOVED***

// Error implements the `error` interface, so K6Errors are normal Go errors.
func (k6Err K6Error) Error() string ***REMOVED***
	return k6Err.Message
***REMOVED***

// Unwrap implements the `xerrors.Wrapper` interface, so K6Errors are a bit
// future-proof Go 2 errors.
func (k6Err K6Error) Unwrap() error ***REMOVED***
	return k6Err.OriginalError
***REMOVED***
