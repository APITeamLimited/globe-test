// Copyright 2015 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// This file implements a http.RoundTripper that authenticates
// requests issued against api.github.com endpoint.

package httputil

import (
	"net/http"
	"net/url"
)

// AuthTransport is an implementation of http.RoundTripper that authenticates
// with the GitHub API.
//
// When both a token and client credentials are set, the latter is preferred.
type AuthTransport struct ***REMOVED***
	UserAgent          string
	GithubToken        string
	GithubClientID     string
	GithubClientSecret string
	Base               http.RoundTripper
***REMOVED***

// RoundTrip implements the http.RoundTripper interface.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	var reqCopy *http.Request
	if t.UserAgent != "" ***REMOVED***
		reqCopy = copyRequest(req)
		reqCopy.Header.Set("User-Agent", t.UserAgent)
	***REMOVED***
	if req.URL.Host == "api.github.com" && req.URL.Scheme == "https" ***REMOVED***
		switch ***REMOVED***
		case t.GithubClientID != "" && t.GithubClientSecret != "":
			if reqCopy == nil ***REMOVED***
				reqCopy = copyRequest(req)
			***REMOVED***
			if reqCopy.URL.RawQuery == "" ***REMOVED***
				reqCopy.URL.RawQuery = "client_id=" + t.GithubClientID + "&client_secret=" + t.GithubClientSecret
			***REMOVED*** else ***REMOVED***
				reqCopy.URL.RawQuery += "&client_id=" + t.GithubClientID + "&client_secret=" + t.GithubClientSecret
			***REMOVED***
		case t.GithubToken != "":
			if reqCopy == nil ***REMOVED***
				reqCopy = copyRequest(req)
			***REMOVED***
			reqCopy.Header.Set("Authorization", "token "+t.GithubToken)
		***REMOVED***
	***REMOVED***
	if reqCopy != nil ***REMOVED***
		return t.base().RoundTrip(reqCopy)
	***REMOVED***
	return t.base().RoundTrip(req)
***REMOVED***

// CancelRequest cancels an in-flight request by closing its connection.
func (t *AuthTransport) CancelRequest(req *http.Request) ***REMOVED***
	type canceler interface ***REMOVED***
		CancelRequest(req *http.Request)
	***REMOVED***
	if cr, ok := t.base().(canceler); ok ***REMOVED***
		cr.CancelRequest(req)
	***REMOVED***
***REMOVED***

func (t *AuthTransport) base() http.RoundTripper ***REMOVED***
	if t.Base != nil ***REMOVED***
		return t.Base
	***REMOVED***
	return http.DefaultTransport
***REMOVED***

func copyRequest(req *http.Request) *http.Request ***REMOVED***
	req2 := new(http.Request)
	*req2 = *req
	req2.URL = new(url.URL)
	*req2.URL = *req.URL
	req2.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header ***REMOVED***
		req2.Header[k] = append([]string(nil), s...)
	***REMOVED***
	return req2
***REMOVED***
