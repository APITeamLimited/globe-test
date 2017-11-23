// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package http2

import (
	"crypto/tls"
	"io"
	"net/http"
)

func cloneTLSConfig(c *tls.Config) *tls.Config ***REMOVED***
	c2 := c.Clone()
	c2.GetClientCertificate = c.GetClientCertificate // golang.org/issue/19264
	return c2
***REMOVED***

var _ http.Pusher = (*responseWriter)(nil)

// Push implements http.Pusher.
func (w *responseWriter) Push(target string, opts *http.PushOptions) error ***REMOVED***
	internalOpts := pushOptions***REMOVED******REMOVED***
	if opts != nil ***REMOVED***
		internalOpts.Method = opts.Method
		internalOpts.Header = opts.Header
	***REMOVED***
	return w.push(target, internalOpts)
***REMOVED***

func configureServer18(h1 *http.Server, h2 *Server) error ***REMOVED***
	if h2.IdleTimeout == 0 ***REMOVED***
		if h1.IdleTimeout != 0 ***REMOVED***
			h2.IdleTimeout = h1.IdleTimeout
		***REMOVED*** else ***REMOVED***
			h2.IdleTimeout = h1.ReadTimeout
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func shouldLogPanic(panicValue interface***REMOVED******REMOVED***) bool ***REMOVED***
	return panicValue != nil && panicValue != http.ErrAbortHandler
***REMOVED***

func reqGetBody(req *http.Request) func() (io.ReadCloser, error) ***REMOVED***
	return req.GetBody
***REMOVED***

func reqBodyIsNoBody(body io.ReadCloser) bool ***REMOVED***
	return body == http.NoBody
***REMOVED***

func go18httpNoBody() io.ReadCloser ***REMOVED*** return http.NoBody ***REMOVED*** // for tests only
