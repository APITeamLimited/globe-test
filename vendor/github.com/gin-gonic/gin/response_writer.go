// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

type (
	ResponseWriter interface ***REMOVED***
		http.ResponseWriter
		http.Hijacker
		http.Flusher
		http.CloseNotifier

		// Returns the HTTP response status code of the current request.
		Status() int

		// Returns the number of bytes already written into the response http body.
		// See Written()
		Size() int

		// Writes the string into the response body.
		WriteString(string) (int, error)

		// Returns true if the response body was already written.
		Written() bool

		// Forces to write the http header (status code + headers).
		WriteHeaderNow()
	***REMOVED***

	responseWriter struct ***REMOVED***
		http.ResponseWriter
		size   int
		status int
	***REMOVED***
)

var _ ResponseWriter = &responseWriter***REMOVED******REMOVED***

func (w *responseWriter) reset(writer http.ResponseWriter) ***REMOVED***
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
***REMOVED***

func (w *responseWriter) WriteHeader(code int) ***REMOVED***
	if code > 0 && w.status != code ***REMOVED***
		if w.Written() ***REMOVED***
			debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		***REMOVED***
		w.status = code
	***REMOVED***
***REMOVED***

func (w *responseWriter) WriteHeaderNow() ***REMOVED***
	if !w.Written() ***REMOVED***
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	***REMOVED***
***REMOVED***

func (w *responseWriter) Write(data []byte) (n int, err error) ***REMOVED***
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
***REMOVED***

func (w *responseWriter) WriteString(s string) (n int, err error) ***REMOVED***
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
***REMOVED***

func (w *responseWriter) Status() int ***REMOVED***
	return w.status
***REMOVED***

func (w *responseWriter) Size() int ***REMOVED***
	return w.size
***REMOVED***

func (w *responseWriter) Written() bool ***REMOVED***
	return w.size != noWritten
***REMOVED***

// Implements the http.Hijacker interface
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) ***REMOVED***
	if w.size < 0 ***REMOVED***
		w.size = 0
	***REMOVED***
	return w.ResponseWriter.(http.Hijacker).Hijack()
***REMOVED***

// Implements the http.CloseNotify interface
func (w *responseWriter) CloseNotify() <-chan bool ***REMOVED***
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
***REMOVED***

// Implements the http.Flush interface
func (w *responseWriter) Flush() ***REMOVED***
	w.ResponseWriter.(http.Flusher).Flush()
***REMOVED***
