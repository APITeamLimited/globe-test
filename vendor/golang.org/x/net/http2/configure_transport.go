// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.6

package http2

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func configureTransport(t1 *http.Transport) (*Transport, error) ***REMOVED***
	connPool := new(clientConnPool)
	t2 := &Transport***REMOVED***
		ConnPool: noDialClientConnPool***REMOVED***connPool***REMOVED***,
		t1:       t1,
	***REMOVED***
	connPool.t = t2
	if err := registerHTTPSProtocol(t1, noDialH2RoundTripper***REMOVED***t2***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t1.TLSClientConfig == nil ***REMOVED***
		t1.TLSClientConfig = new(tls.Config)
	***REMOVED***
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "h2") ***REMOVED***
		t1.TLSClientConfig.NextProtos = append([]string***REMOVED***"h2"***REMOVED***, t1.TLSClientConfig.NextProtos...)
	***REMOVED***
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "http/1.1") ***REMOVED***
		t1.TLSClientConfig.NextProtos = append(t1.TLSClientConfig.NextProtos, "http/1.1")
	***REMOVED***
	upgradeFn := func(authority string, c *tls.Conn) http.RoundTripper ***REMOVED***
		addr := authorityAddr("https", authority)
		if used, err := connPool.addConnIfNeeded(addr, t2, c); err != nil ***REMOVED***
			go c.Close()
			return erringRoundTripper***REMOVED***err***REMOVED***
		***REMOVED*** else if !used ***REMOVED***
			// Turns out we don't need this c.
			// For example, two goroutines made requests to the same host
			// at the same time, both kicking off TCP dials. (since protocol
			// was unknown)
			go c.Close()
		***REMOVED***
		return t2
	***REMOVED***
	if m := t1.TLSNextProto; len(m) == 0 ***REMOVED***
		t1.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper***REMOVED***
			"h2": upgradeFn,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		m["h2"] = upgradeFn
	***REMOVED***
	return t2, nil
***REMOVED***

// registerHTTPSProtocol calls Transport.RegisterProtocol but
// converting panics into errors.
func registerHTTPSProtocol(t *http.Transport, rt http.RoundTripper) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if e := recover(); e != nil ***REMOVED***
			err = fmt.Errorf("%v", e)
		***REMOVED***
	***REMOVED***()
	t.RegisterProtocol("https", rt)
	return nil
***REMOVED***

// noDialH2RoundTripper is a RoundTripper which only tries to complete the request
// if there's already has a cached connection to the host.
type noDialH2RoundTripper struct***REMOVED*** t *Transport ***REMOVED***

func (rt noDialH2RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	res, err := rt.t.RoundTrip(req)
	if isNoCachedConnError(err) ***REMOVED***
		return nil, http.ErrSkipAltProtocol
	***REMOVED***
	return res, err
***REMOVED***
