// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputil

import (
	"bytes"
	"net/http"
	"strconv"
)

// ResponseBuffer is the current response being composed by its owner.
// It implements http.ResponseWriter and io.WriterTo.
type ResponseBuffer struct ***REMOVED***
	buf    bytes.Buffer
	status int
	header http.Header
***REMOVED***

// Write implements the http.ResponseWriter interface.
func (rb *ResponseBuffer) Write(p []byte) (int, error) ***REMOVED***
	return rb.buf.Write(p)
***REMOVED***

// WriteHeader implements the http.ResponseWriter interface.
func (rb *ResponseBuffer) WriteHeader(status int) ***REMOVED***
	rb.status = status
***REMOVED***

// Header implements the http.ResponseWriter interface.
func (rb *ResponseBuffer) Header() http.Header ***REMOVED***
	if rb.header == nil ***REMOVED***
		rb.header = make(http.Header)
	***REMOVED***
	return rb.header
***REMOVED***

// WriteTo implements the io.WriterTo interface.
func (rb *ResponseBuffer) WriteTo(w http.ResponseWriter) error ***REMOVED***
	for k, v := range rb.header ***REMOVED***
		w.Header()[k] = v
	***REMOVED***
	if rb.buf.Len() > 0 ***REMOVED***
		w.Header().Set("Content-Length", strconv.Itoa(rb.buf.Len()))
	***REMOVED***
	if rb.status != 0 ***REMOVED***
		w.WriteHeader(rb.status)
	***REMOVED***
	if rb.buf.Len() > 0 ***REMOVED***
		if _, err := w.Write(rb.buf.Bytes()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
