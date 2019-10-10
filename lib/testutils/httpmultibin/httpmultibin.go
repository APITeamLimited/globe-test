/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

// Package httpmultibin is indended only for use in tests, do not import in production code!
package httpmultibin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
	"github.com/klauspost/compress/zstd"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/mccutchen/go-httpbin/httpbin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

// GetTLSClientConfig returns a TLS config that trusts the supplied
// httptest.Server certificate as well as all the system root certificates
func GetTLSClientConfig(t testing.TB, srv *httptest.Server) *tls.Config ***REMOVED***
	var err error

	certs := x509.NewCertPool()

	if runtime.GOOS != "windows" ***REMOVED***
		certs, err = x509.SystemCertPool()
		require.NoError(t, err)
	***REMOVED***

	for _, c := range srv.TLS.Certificates ***REMOVED***
		roots, err := x509.ParseCertificates(c.Certificate[len(c.Certificate)-1])
		require.NoError(t, err)
		for _, root := range roots ***REMOVED***
			certs.AddCert(root)
		***REMOVED***
	***REMOVED***
	return &tls.Config***REMOVED***
		RootCAs:            certs,
		InsecureSkipVerify: false,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	***REMOVED***
***REMOVED***

const httpDomain = "httpbin.local"

// We have to use example.com if we want a real HTTPS domain with a valid
// certificate because the default httptest certificate is for example.com:
// https://golang.org/src/net/http/internal/testcert.go?s=399:410#L10
const httpsDomain = "example.com"

// HTTPMultiBin can be used as a local alternative of httpbin.org. It offers both http and https servers, as well as real domains
type HTTPMultiBin struct ***REMOVED***
	Mux             *http.ServeMux
	ServerHTTP      *httptest.Server
	ServerHTTPS     *httptest.Server
	Replacer        *strings.Replacer
	TLSClientConfig *tls.Config
	Dialer          *netext.Dialer
	HTTPTransport   *http.Transport
	Context         context.Context
	Cleanup         func()
***REMOVED***

type jsonBody struct ***REMOVED***
	Header      http.Header `json:"headers"`
	Compression string      `json:"compression"`
***REMOVED***

func websocketEchoHandler(w http.ResponseWriter, req *http.Request) ***REMOVED***
	conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, req, w.Header())
	if err != nil ***REMOVED***
		return
	***REMOVED***

	mt, message, err := conn.ReadMessage()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = conn.WriteMessage(mt, message)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = conn.Close()
	if err != nil ***REMOVED***
		return
	***REMOVED***
***REMOVED***

func websocketCloserHandler(w http.ResponseWriter, req *http.Request) ***REMOVED***
	conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, req, w.Header())
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_ = conn.Close()
***REMOVED***

func writeJSON(w io.Writer, v interface***REMOVED******REMOVED***) error ***REMOVED***
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return errors.Wrap(e.Encode(v), "failed to encode JSON")
***REMOVED***

func getEncodedHandler(t testing.TB, compressionType httpext.CompressionType) http.Handler ***REMOVED***
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) ***REMOVED***
		var (
			encoding string
			err      error
			encw     io.WriteCloser
		)

		switch compressionType ***REMOVED***
		case httpext.CompressionTypeBr:
			encw = brotli.NewWriter(rw)
			encoding = "br"
		case httpext.CompressionTypeZstd:
			encw, _ = zstd.NewWriter(rw)
			encoding = "zstd"
		***REMOVED***

		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Add("Content-Encoding", encoding)
		data := jsonBody***REMOVED***
			Header:      req.Header,
			Compression: encoding,
		***REMOVED***
		err = writeJSON(encw, data)
		_ = encw.Close()
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***)
***REMOVED***

func getZstdBrHandler(t testing.TB) http.Handler ***REMOVED***
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) ***REMOVED***
		encoding := "zstd, br"
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Add("Content-Encoding", encoding)
		data := jsonBody***REMOVED***
			Header:      req.Header,
			Compression: encoding,
		***REMOVED***

		bw := brotli.NewWriter(rw)
		zw, _ := zstd.NewWriter(bw)
		defer func() ***REMOVED***
			_ = zw.Close()
			_ = bw.Close()
		***REMOVED***()

		require.NoError(t, writeJSON(zw, data))
	***REMOVED***)
***REMOVED***

// NewHTTPMultiBin returns a fully configured and running HTTPMultiBin
func NewHTTPMultiBin(t testing.TB) *HTTPMultiBin ***REMOVED***
	// Create a http.ServeMux and set the httpbin handler as the default
	mux := http.NewServeMux()
	mux.Handle("/brotli", getEncodedHandler(t, httpext.CompressionTypeBr))
	mux.HandleFunc("/ws-echo", websocketEchoHandler)
	mux.HandleFunc("/ws-close", websocketCloserHandler)
	mux.Handle("/zstd", getEncodedHandler(t, httpext.CompressionTypeZstd))
	mux.Handle("/zstd-br", getZstdBrHandler(t))
	mux.Handle("/", httpbin.New().Handler())

	// Initialize the HTTP server and get its details
	httpSrv := httptest.NewServer(mux)
	httpURL, err := url.Parse(httpSrv.URL)
	require.NoError(t, err)
	httpIP := net.ParseIP(httpURL.Hostname())
	require.NotNil(t, httpIP)

	// Initialize the HTTPS server and get its details and tls config
	httpsSrv := httptest.NewTLSServer(mux)
	httpsURL, err := url.Parse(httpsSrv.URL)
	require.NoError(t, err)
	httpsIP := net.ParseIP(httpsURL.Hostname())
	require.NotNil(t, httpsIP)
	tlsConfig := GetTLSClientConfig(t, httpsSrv)

	// Set up the dialer with shorter timeouts and the custom domains
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   2 * time.Second,
		KeepAlive: 10 * time.Second,
		DualStack: true,
	***REMOVED***)
	dialer.Hosts = map[string]net.IP***REMOVED***
		httpDomain:  httpIP,
		httpsDomain: httpsIP,
	***REMOVED***

	// Pre-configure the HTTP client transport with the dialer and TLS config (incl. HTTP2 support)
	transport := &http.Transport***REMOVED***
		DialContext:     dialer.DialContext,
		TLSClientConfig: tlsConfig,
	***REMOVED***
	require.NoError(t, http2.ConfigureTransport(transport))

	ctx, ctxCancel := context.WithCancel(context.Background())
	return &HTTPMultiBin***REMOVED***
		Mux:         mux,
		ServerHTTP:  httpSrv,
		ServerHTTPS: httpsSrv,
		Replacer: strings.NewReplacer(
			"HTTPBIN_IP_URL", httpSrv.URL,
			"HTTPBIN_DOMAIN", httpDomain,
			"HTTPBIN_URL", fmt.Sprintf("http://%s:%s", httpDomain, httpURL.Port()),
			"WSBIN_URL", fmt.Sprintf("ws://%s:%s", httpDomain, httpURL.Port()),
			"HTTPBIN_IP", httpIP.String(),
			"HTTPBIN_PORT", httpURL.Port(),
			"HTTPSBIN_IP_URL", httpsSrv.URL,
			"HTTPSBIN_DOMAIN", httpsDomain,
			"HTTPSBIN_URL", fmt.Sprintf("https://%s:%s", httpsDomain, httpsURL.Port()),
			"WSSBIN_URL", fmt.Sprintf("wss://%s:%s", httpsDomain, httpsURL.Port()),
			"HTTPSBIN_IP", httpsIP.String(),
			"HTTPSBIN_PORT", httpsURL.Port(),
		),
		TLSClientConfig: tlsConfig,
		Dialer:          dialer,
		HTTPTransport:   transport,
		Context:         ctx,
		Cleanup: func() ***REMOVED***
			httpsSrv.Close()
			httpSrv.Close()
			ctxCancel()
		***REMOVED***,
	***REMOVED***
***REMOVED***
