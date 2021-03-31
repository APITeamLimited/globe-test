/*
 *
 * Copyright 2017 gRPC authors.
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
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const proxyAuthHeaderKey = "Proxy-Authorization"

var (
	// The following variable will be overwritten in the tests.
	httpProxyFromEnvironment = http.ProxyFromEnvironment
)

func mapAddress(ctx context.Context, address string) (*url.URL, error) ***REMOVED***
	req := &http.Request***REMOVED***
		URL: &url.URL***REMOVED***
			Scheme: "https",
			Host:   address,
		***REMOVED***,
	***REMOVED***
	url, err := httpProxyFromEnvironment(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return url, nil
***REMOVED***

// To read a response from a net.Conn, http.ReadResponse() takes a bufio.Reader.
// It's possible that this reader reads more than what's need for the response and stores
// those bytes in the buffer.
// bufConn wraps the original net.Conn and the bufio.Reader to make sure we don't lose the
// bytes in the buffer.
type bufConn struct ***REMOVED***
	net.Conn
	r io.Reader
***REMOVED***

func (c *bufConn) Read(b []byte) (int, error) ***REMOVED***
	return c.r.Read(b)
***REMOVED***

func basicAuth(username, password string) string ***REMOVED***
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
***REMOVED***

func doHTTPConnectHandshake(ctx context.Context, conn net.Conn, backendAddr string, proxyURL *url.URL, grpcUA string) (_ net.Conn, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			conn.Close()
		***REMOVED***
	***REMOVED***()

	req := &http.Request***REMOVED***
		Method: http.MethodConnect,
		URL:    &url.URL***REMOVED***Host: backendAddr***REMOVED***,
		Header: map[string][]string***REMOVED***"User-Agent": ***REMOVED***grpcUA***REMOVED******REMOVED***,
	***REMOVED***
	if t := proxyURL.User; t != nil ***REMOVED***
		u := t.Username()
		p, _ := t.Password()
		req.Header.Add(proxyAuthHeaderKey, "Basic "+basicAuth(u, p))
	***REMOVED***

	if err := sendHTTPRequest(ctx, req, conn); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to write the HTTP request: %v", err)
	***REMOVED***

	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, req)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("reading server HTTP response: %v", err)
	***REMOVED***
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK ***REMOVED***
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to do connect handshake, status code: %s", resp.Status)
		***REMOVED***
		return nil, fmt.Errorf("failed to do connect handshake, response: %q", dump)
	***REMOVED***

	return &bufConn***REMOVED***Conn: conn, r: r***REMOVED***, nil
***REMOVED***

// proxyDial dials, connecting to a proxy first if necessary. Checks if a proxy
// is necessary, dials, does the HTTP CONNECT handshake, and returns the
// connection.
func proxyDial(ctx context.Context, addr string, grpcUA string) (conn net.Conn, err error) ***REMOVED***
	newAddr := addr
	proxyURL, err := mapAddress(ctx, addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if proxyURL != nil ***REMOVED***
		newAddr = proxyURL.Host
	***REMOVED***

	conn, err = (&net.Dialer***REMOVED******REMOVED***).DialContext(ctx, "tcp", newAddr)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if proxyURL != nil ***REMOVED***
		// proxy is disabled if proxyURL is nil.
		conn, err = doHTTPConnectHandshake(ctx, conn, addr, proxyURL, grpcUA)
	***REMOVED***
	return
***REMOVED***

func sendHTTPRequest(ctx context.Context, req *http.Request, conn net.Conn) error ***REMOVED***
	req = req.WithContext(ctx)
	if err := req.Write(conn); err != nil ***REMOVED***
		return fmt.Errorf("failed to write the HTTP request: %v", err)
	***REMOVED***
	return nil
***REMOVED***
