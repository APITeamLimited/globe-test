// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type netDialerFunc func(network, addr string) (net.Conn, error)

func (fn netDialerFunc) Dial(network, addr string) (net.Conn, error) ***REMOVED***
	return fn(network, addr)
***REMOVED***

func init() ***REMOVED***
	proxy_RegisterDialerType("http", func(proxyURL *url.URL, forwardDialer proxy_Dialer) (proxy_Dialer, error) ***REMOVED***
		return &httpProxyDialer***REMOVED***proxyURL: proxyURL, forwardDial: forwardDialer.Dial***REMOVED***, nil
	***REMOVED***)
***REMOVED***

type httpProxyDialer struct ***REMOVED***
	proxyURL    *url.URL
	forwardDial func(network, addr string) (net.Conn, error)
***REMOVED***

func (hpd *httpProxyDialer) Dial(network string, addr string) (net.Conn, error) ***REMOVED***
	hostPort, _ := hostPortNoPort(hpd.proxyURL)
	conn, err := hpd.forwardDial(network, hostPort)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	connectHeader := make(http.Header)
	if user := hpd.proxyURL.User; user != nil ***REMOVED***
		proxyUser := user.Username()
		if proxyPassword, passwordSet := user.Password(); passwordSet ***REMOVED***
			credential := base64.StdEncoding.EncodeToString([]byte(proxyUser + ":" + proxyPassword))
			connectHeader.Set("Proxy-Authorization", "Basic "+credential)
		***REMOVED***
	***REMOVED***

	connectReq := &http.Request***REMOVED***
		Method: "CONNECT",
		URL:    &url.URL***REMOVED***Opaque: addr***REMOVED***,
		Host:   addr,
		Header: connectHeader,
	***REMOVED***

	if err := connectReq.Write(conn); err != nil ***REMOVED***
		conn.Close()
		return nil, err
	***REMOVED***

	// Read response. It's OK to use and discard buffered reader here becaue
	// the remote server does not speak until spoken to.
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil ***REMOVED***
		conn.Close()
		return nil, err
	***REMOVED***

	if resp.StatusCode != 200 ***REMOVED***
		conn.Close()
		f := strings.SplitN(resp.Status, " ", 2)
		return nil, errors.New(f[1])
	***REMOVED***
	return conn, nil
***REMOVED***
