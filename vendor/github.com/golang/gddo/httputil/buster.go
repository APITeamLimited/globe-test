// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputil

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type busterWriter struct ***REMOVED***
	headerMap http.Header
	status    int
	io.Writer
***REMOVED***

func (bw *busterWriter) Header() http.Header ***REMOVED***
	return bw.headerMap
***REMOVED***

func (bw *busterWriter) WriteHeader(status int) ***REMOVED***
	bw.status = status
***REMOVED***

// CacheBusters maintains a cache of cache busting tokens for static resources served by Handler.
type CacheBusters struct ***REMOVED***
	Handler http.Handler

	mu     sync.Mutex
	tokens map[string]string
***REMOVED***

func sanitizeTokenRune(r rune) rune ***REMOVED***
	if r <= ' ' || r >= 127 ***REMOVED***
		return -1
	***REMOVED***
	// Convert percent encoding reserved characters to '-'.
	if strings.ContainsRune("!#$&'()*+,/:;=?@[]", r) ***REMOVED***
		return '-'
	***REMOVED***
	return r
***REMOVED***

// Get returns the cache busting token for path. If the token is not already
// cached, Get issues a HEAD request on handler and uses the response ETag and
// Last-Modified headers to compute a token.
func (cb *CacheBusters) Get(path string) string ***REMOVED***
	cb.mu.Lock()
	if cb.tokens == nil ***REMOVED***
		cb.tokens = make(map[string]string)
	***REMOVED***
	token, ok := cb.tokens[path]
	cb.mu.Unlock()
	if ok ***REMOVED***
		return token
	***REMOVED***

	w := busterWriter***REMOVED***
		Writer:    ioutil.Discard,
		headerMap: make(http.Header),
	***REMOVED***
	r := &http.Request***REMOVED***URL: &url.URL***REMOVED***Path: path***REMOVED***, Method: "HEAD"***REMOVED***
	cb.Handler.ServeHTTP(&w, r)

	if w.status == 200 ***REMOVED***
		token = w.headerMap.Get("Etag")
		if token == "" ***REMOVED***
			token = w.headerMap.Get("Last-Modified")
		***REMOVED***
		token = strings.Trim(token, `" `)
		token = strings.Map(sanitizeTokenRune, token)
	***REMOVED***

	cb.mu.Lock()
	cb.tokens[path] = token
	cb.mu.Unlock()

	return token
***REMOVED***

// AppendQueryParam appends the token as a query parameter to path.
func (cb *CacheBusters) AppendQueryParam(path string, name string) string ***REMOVED***
	token := cb.Get(path)
	if token == "" ***REMOVED***
		return path
	***REMOVED***
	return path + "?" + name + "=" + token
***REMOVED***
