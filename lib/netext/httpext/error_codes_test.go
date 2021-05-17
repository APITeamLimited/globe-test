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
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/lib/types"
)

func TestDefaultError(t *testing.T) ***REMOVED***
	t.Parallel()
	testErrorCode(t, defaultErrorCode, fmt.Errorf("random error"))
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

func TestHTTP2StreamError(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	tb.Mux.HandleFunc("/tsr", func(rw http.ResponseWriter, req *http.Request) ***REMOVED***
		rw.Header().Set("Content-Length", "100000")
		rw.WriteHeader(200)

		rw.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 2)
		panic("expected internal error")
	***REMOVED***)
	client := http.Client***REMOVED***
		Timeout:   time.Second * 3,
		Transport: tb.HTTPTransport,
	***REMOVED***

	res, err := client.Get(tb.Replacer.Replace("HTTP2BIN_URL/tsr")) //nolint:noctx
	require.NotNil(t, res)
	require.NoError(t, err)
	_, err = ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	require.Error(t, err)

	code, msg := errorCodeForError(err)
	assert.Equal(t, unknownHTTP2StreamErrorCode+errCode(http2.ErrCodeInternal)+1, code)
	assert.Contains(t, msg, fmt.Sprintf(http2StreamErrorCodeMsg, http2.ErrCodeInternal))
***REMOVED***

func TestX509HostnameError(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	client := http.Client***REMOVED***
		Timeout:   time.Second * 3,
		Transport: tb.HTTPTransport,
	***REMOVED***
	var err error
	badHostname := "somewhere.else"
	tb.Dialer.Hosts[badHostname], err = lib.NewHostAddress(net.ParseIP(tb.Replacer.Replace("HTTPSBIN_IP")), "")
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(context.Background(), "GET", tb.Replacer.Replace("https://"+badHostname+":HTTPSBIN_PORT/get"), nil)
	require.NoError(t, err)
	res, err := client.Do(req) //nolint:bodyclose
	require.Nil(t, res)
	require.Error(t, err)

	code, msg := errorCodeForError(err)
	assert.Equal(t, x509HostnameErrorCode, code)
	assert.Contains(t, msg, x509HostnameErrorCodeMsg)
***REMOVED***

func TestX509UnknownAuthorityError(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	client := http.Client***REMOVED***
		Timeout: time.Second * 3,
		Transport: &http.Transport***REMOVED***
			DialContext: tb.HTTPTransport.DialContext,
		***REMOVED***,
	***REMOVED***
	req, err := http.NewRequestWithContext(context.Background(), "GET", tb.Replacer.Replace("HTTPSBIN_URL/get"), nil)
	require.NoError(t, err)
	res, err := client.Do(req) //nolint:bodyclose
	require.Nil(t, res)
	require.Error(t, err)

	code, msg := errorCodeForError(err)
	assert.Equal(t, x509UnknownAuthorityErrorCode, code)
	assert.Contains(t, msg, x509UnknownAuthority)
***REMOVED***

func TestDefaultTLSError(t *testing.T) ***REMOVED***
	t.Parallel()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() ***REMOVED***
		conn, err := l.Accept() //nolint:govet // the shadowing is intentional
		require.NoError(t, err)
		_, err = conn.Write([]byte("not tls header")) // we just want to get an error
		require.NoError(t, err)
		// wait so it has time to get the tls header error and not the reset socket one
		time.Sleep(time.Second)
	***REMOVED***()

	client := http.Client***REMOVED***
		Timeout: time.Second * 3,
		Transport: &http.Transport***REMOVED***
			TLSClientConfig: &tls.Config***REMOVED***
				InsecureSkipVerify: true, //nolint:gosec
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	_, err = client.Get("https://" + l.Addr().String()) //nolint:bodyclose,noctx
	require.Error(t, err)

	code, msg := errorCodeForError(err)
	assert.Equal(t, defaultTLSErrorCode, code)
	urlError := new(url.Error)
	require.ErrorAs(t, err, &urlError)
	assert.Equal(t, urlError.Err.Error(), msg)
***REMOVED***

func TestHTTP2ConnectionError(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := getHTTP2ServerWithCustomConnContext(t)

	// Pre-configure the HTTP client transport with the dialer and TLS config (incl. HTTP2 support)
	tb.Mux.HandleFunc("/tsr", func(rw http.ResponseWriter, req *http.Request) ***REMOVED***
		conn := req.Context().Value(connKey).(*tls.Conn) //nolint:forcetypeassert
		f := http2.NewFramer(conn, conn)
		require.NoError(t, f.WriteData(3213, false, []byte("something")))
	***REMOVED***)
	client := http.Client***REMOVED***
		Timeout:   time.Second * 5,
		Transport: tb.HTTPTransport,
	***REMOVED***

	_, err := client.Get(tb.Replacer.Replace("HTTP2BIN_URL/tsr")) //nolint:bodyclose,noctx
	code, msg := errorCodeForError(err)
	assert.Equal(t, unknownHTTP2ConnectionErrorCode+errCode(http2.ErrCodeProtocol)+1, code)
	assert.Equal(t, fmt.Sprintf(http2ConnectionErrorCodeMsg, http2.ErrCodeProtocol), msg)
***REMOVED***

func TestHTTP2GoAwayError(t *testing.T) ***REMOVED***
	t.Parallel()

	tb := getHTTP2ServerWithCustomConnContext(t)
	tb.Mux.HandleFunc("/tsr", func(rw http.ResponseWriter, req *http.Request) ***REMOVED***
		conn := req.Context().Value(connKey).(*tls.Conn) //nolint:forcetypeassert
		f := http2.NewFramer(conn, conn)
		require.NoError(t, f.WriteGoAway(4, http2.ErrCodeInadequateSecurity, []byte("whatever")))
		require.NoError(t, conn.CloseWrite())
	***REMOVED***)
	client := http.Client***REMOVED***
		Timeout:   time.Second * 5,
		Transport: tb.HTTPTransport,
	***REMOVED***

	_, err := client.Get(tb.Replacer.Replace("HTTP2BIN_URL/tsr")) //nolint:bodyclose,noctx

	require.Error(t, err)
	code, msg := errorCodeForError(err)
	assert.Equal(t, unknownHTTP2GoAwayErrorCode+errCode(http2.ErrCodeInadequateSecurity)+1, code)
	assert.Equal(t, fmt.Sprintf(http2GoAwayErrorCodeMsg, http2.ErrCodeInadequateSecurity), msg)
***REMOVED***

type connKeyT int32

const connKey connKeyT = 2

func getHTTP2ServerWithCustomConnContext(t *testing.T) *httpmultibin.HTTPMultiBin ***REMOVED***
	const http2Domain = "example.com"
	mux := http.NewServeMux()
	http2Srv := httptest.NewUnstartedServer(mux)
	http2Srv.EnableHTTP2 = true
	http2Srv.Config.ConnContext = func(ctx context.Context, c net.Conn) context.Context ***REMOVED***
		return context.WithValue(ctx, connKey, c)
	***REMOVED***
	http2Srv.StartTLS()
	t.Cleanup(http2Srv.Close)
	tlsConfig := httpmultibin.GetTLSClientConfig(t, http2Srv)

	http2URL, err := url.Parse(http2Srv.URL)
	require.NoError(t, err)
	http2IP := net.ParseIP(http2URL.Hostname())
	require.NotNil(t, http2IP)
	http2DomainValue, err := lib.NewHostAddress(http2IP, "")
	require.NoError(t, err)

	// Set up the dialer with shorter timeouts and the custom domains
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   2 * time.Second,
		KeepAlive: 10 * time.Second,
		DualStack: true,
	***REMOVED***, netext.NewResolver(net.LookupIP, 0, types.DNSfirst, types.DNSpreferIPv4))
	dialer.Hosts = map[string]*lib.HostAddress***REMOVED***
		http2Domain: http2DomainValue,
	***REMOVED***

	transport := &http.Transport***REMOVED***
		DialContext:     dialer.DialContext,
		TLSClientConfig: tlsConfig,
	***REMOVED***
	require.NoError(t, http2.ConfigureTransport(transport))
	return &httpmultibin.HTTPMultiBin***REMOVED***
		Mux:         mux,
		ServerHTTP2: http2Srv,
		Replacer: strings.NewReplacer(
			"HTTP2BIN_IP_URL", http2Srv.URL,
			"HTTP2BIN_DOMAIN", http2Domain,
			"HTTP2BIN_URL", fmt.Sprintf("https://%s:%s", http2Domain, http2URL.Port()),
			"HTTP2BIN_IP", http2IP.String(),
			"HTTP2BIN_PORT", http2URL.Port(),
		),
		TLSClientConfig: tlsConfig,
		Dialer:          dialer,
		HTTPTransport:   transport,
	***REMOVED***
***REMOVED***
