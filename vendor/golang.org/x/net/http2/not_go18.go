// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.8

package http2

import (
	"io"
	"net/http"
)

func configureServer18(h1 *http.Server, h2 *Server) error ***REMOVED***
	// No IdleTimeout to sync prior to Go 1.8.
	return nil
***REMOVED***

func shouldLogPanic(panicValue interface***REMOVED******REMOVED***) bool ***REMOVED***
	return panicValue != nil
***REMOVED***

func reqGetBody(req *http.Request) func() (io.ReadCloser, error) ***REMOVED***
	return nil
***REMOVED***

func reqBodyIsNoBody(io.ReadCloser) bool ***REMOVED*** return false ***REMOVED***

func go18httpNoBody() io.ReadCloser ***REMOVED*** return nil ***REMOVED*** // for tests only
