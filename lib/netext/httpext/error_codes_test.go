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
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"

	"go.k6.io/k6/lib/netext"
)

func TestDefaultError(t *testing.T) ***REMOVED***
	t.Parallel()
	testErrorCode(t, defaultErrorCode, fmt.Errorf("random error"))
***REMOVED***

func TestHTTP2Errors(t *testing.T) ***REMOVED***
	t.Parallel()
	unknownErrorCode := 220
	connectionError := http2.ConnectionError(unknownErrorCode)
	testTable := map[errCode]error***REMOVED***
		unknownHTTP2ConnectionErrorCode + 1: new(http2.ConnectionError),
		unknownHTTP2StreamErrorCode + 1:     new(http2.StreamError),
		unknownHTTP2GoAwayErrorCode + 1:     new(http2.GoAwayError),

		unknownHTTP2ConnectionErrorCode: &connectionError,
		unknownHTTP2StreamErrorCode:     &http2.StreamError***REMOVED***Code: 220***REMOVED***,
		unknownHTTP2GoAwayErrorCode:     &http2.GoAwayError***REMOVED***ErrCode: 220***REMOVED***,
	***REMOVED***
	testMapOfErrorCodes(t, testTable)
***REMOVED***

func TestTLSErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	testTable := map[errCode]error***REMOVED***
		x509UnknownAuthorityErrorCode: new(x509.UnknownAuthorityError),
		x509HostnameErrorCode:         new(x509.HostnameError),
		defaultTLSErrorCode:           new(tls.RecordHeaderError),
	***REMOVED***
	testMapOfErrorCodes(t, testTable)
***REMOVED***

func TestDNSErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	var (
		defaultDNSError = new(net.DNSError)
		noSuchHostError = new(net.DNSError)
	)

	noSuchHostError.Err = "no such host" // defined as private in go stdlib
	testTable := map[errCode]error***REMOVED***
		defaultDNSErrorCode:    defaultDNSError,
		dnsNoSuchHostErrorCode: noSuchHostError,
	***REMOVED***
	testMapOfErrorCodes(t, testTable)
***REMOVED***

func TestBlackListedIPError(t *testing.T) ***REMOVED***
	t.Parallel()
	err := netext.BlackListedIPError***REMOVED******REMOVED***
	testErrorCode(t, blackListedIPErrorCode, err)
	errorCode, errorMsg := errorCodeForError(err)
	require.NotEqual(t, err.Error(), errorMsg)
	require.Equal(t, blackListedIPErrorCode, errorCode)
***REMOVED***

type timeoutError bool

func (t timeoutError) Timeout() bool ***REMOVED***
	return (bool)(t)
***REMOVED***

func (t timeoutError) Error() string ***REMOVED***
	return fmt.Sprintf("%t", t)
***REMOVED***

func TestUnknownNetErrno(t *testing.T) ***REMOVED***
	t.Parallel()
	err := new(net.OpError)
	err.Op = "write"
	err.Net = "tcp"
	err.Err = syscall.ENOTRECOVERABLE // Highly unlikely to actually need to do anything with this error
	expectedError := fmt.Sprintf(
		"write: unknown errno `%d` on %s with message `%s`",
		syscall.ENOTRECOVERABLE, runtime.GOOS, err.Err)
	errorCode, errorMsg := errorCodeForError(err)
	require.Equal(t, expectedError, errorMsg)
	require.Equal(t, netUnknownErrnoErrorCode, errorCode)
***REMOVED***

func TestTCPErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	var (
		nonTCPError       = &net.OpError***REMOVED***Net: "something", Err: errors.New("non tcp error")***REMOVED***
		econnreset        = &net.OpError***REMOVED***Net: "tcp", Op: "write", Err: &os.SyscallError***REMOVED***Err: syscall.ECONNRESET***REMOVED******REMOVED***
		epipeerror        = &net.OpError***REMOVED***Net: "tcp", Op: "write", Err: &os.SyscallError***REMOVED***Err: syscall.EPIPE***REMOVED******REMOVED***
		econnrefused      = &net.OpError***REMOVED***Net: "tcp", Op: "dial", Err: &os.SyscallError***REMOVED***Err: syscall.ECONNREFUSED***REMOVED******REMOVED***
		errnounknown      = &net.OpError***REMOVED***Net: "tcp", Op: "dial", Err: &os.SyscallError***REMOVED***Err: syscall.E2BIG***REMOVED******REMOVED***
		tcperror          = &net.OpError***REMOVED***Net: "tcp", Err: errors.New("tcp error")***REMOVED***
		notTimeoutedError = &net.OpError***REMOVED***Net: "tcp", Op: "dial", Err: timeoutError(false)***REMOVED***
	)

	testTable := map[errCode]error***REMOVED***
		defaultNetNonTCPErrorCode: nonTCPError,
		tcpResetByPeerErrorCode:   econnreset,
		tcpBrokenPipeErrorCode:    epipeerror,
		tcpDialRefusedErrorCode:   econnrefused,
		tcpDialUnknownErrnoCode:   errnounknown,
		defaultTCPErrorCode:       tcperror,
		tcpDialErrorCode:          notTimeoutedError,
	***REMOVED***

	testMapOfErrorCodes(t, testTable)
***REMOVED***

func testErrorCode(t *testing.T, code errCode, err error) ***REMOVED***
	t.Helper()
	result, _ := errorCodeForError(err)
	require.Equalf(t, code, result, "Wrong error code for error `%s`", err)

	result, _ = errorCodeForError(fmt.Errorf("foo: %w", err))
	require.Equalf(t, code, result, "Wrong error code for error `%s`", err)

	result, _ = errorCodeForError(&url.Error***REMOVED***Err: err***REMOVED***)
	require.Equalf(t, code, result, "Wrong error code for error `%s`", err)
***REMOVED***

func testMapOfErrorCodes(t *testing.T, testTable map[errCode]error) ***REMOVED***
	t.Helper()
	for code, err := range testTable ***REMOVED***
		testErrorCode(t, code, err)
	***REMOVED***
***REMOVED***

func TestConnReset(t *testing.T) ***REMOVED***
	t.Parallel()
	// based on https://gist.github.com/jpittis/4357d817dc425ae99fbf719828ab1800
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	addr := ln.Addr()
	ch := make(chan error, 10)

	go func() ***REMOVED***
		defer close(ch)
		// Accept one connection.
		conn, innerErr := ln.Accept()
		if innerErr != nil ***REMOVED***
			ch <- innerErr
			return
		***REMOVED***

		// Force an RST
		tcpConn, ok := conn.(*net.TCPConn)
		require.True(t, ok)
		innerErr = tcpConn.SetLinger(0)
		if innerErr != nil ***REMOVED***
			ch <- innerErr
		***REMOVED***
		time.Sleep(time.Second) // Give time for the http request to start
		_ = conn.Close()
	***REMOVED***()

	res, err := http.Get("http://" + addr.String()) //nolint:bodyclose,noctx
	require.Nil(t, res)

	code, msg := errorCodeForError(err)
	assert.Equal(t, tcpResetByPeerErrorCode, code)
	assert.Contains(t, msg, fmt.Sprintf(tcpResetByPeerErrorCodeMsg, ""))
	for err := range ch ***REMOVED***
		assert.Nil(t, err)
	***REMOVED***
***REMOVED***

func TestDnsResolve(t *testing.T) ***REMOVED***
	t.Parallel()
	// this uses the Unwrap path
	// this is not happening in our current codebase as the resolution in our code
	// happens earlier so it doesn't get wrapped, but possibly happens in other cases as well
	_, err := http.Get("http://s.com") //nolint:bodyclose,noctx
	code, msg := errorCodeForError(err)

	assert.Equal(t, dnsNoSuchHostErrorCode, code)
	assert.Equal(t, dnsNoSuchHostErrorCodeMsg, msg)
***REMOVED***
