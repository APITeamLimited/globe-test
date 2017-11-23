package negroni

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface ***REMOVED***
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func(ResponseWriter))
***REMOVED***

type beforeFunc func(ResponseWriter)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter ***REMOVED***
	nrw := &responseWriter***REMOVED***
		ResponseWriter: rw,
	***REMOVED***

	if _, ok := rw.(http.CloseNotifier); ok ***REMOVED***
		return &responseWriterCloseNotifer***REMOVED***nrw***REMOVED***
	***REMOVED***

	return nrw
***REMOVED***

type responseWriter struct ***REMOVED***
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
***REMOVED***

func (rw *responseWriter) WriteHeader(s int) ***REMOVED***
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
***REMOVED***

func (rw *responseWriter) Write(b []byte) (int, error) ***REMOVED***
	if !rw.Written() ***REMOVED***
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	***REMOVED***
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
***REMOVED***

func (rw *responseWriter) Status() int ***REMOVED***
	return rw.status
***REMOVED***

func (rw *responseWriter) Size() int ***REMOVED***
	return rw.size
***REMOVED***

func (rw *responseWriter) Written() bool ***REMOVED***
	return rw.status != 0
***REMOVED***

func (rw *responseWriter) Before(before func(ResponseWriter)) ***REMOVED***
	rw.beforeFuncs = append(rw.beforeFuncs, before)
***REMOVED***

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) ***REMOVED***
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok ***REMOVED***
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	***REMOVED***
	return hijacker.Hijack()
***REMOVED***

func (rw *responseWriter) callBefore() ***REMOVED***
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- ***REMOVED***
		rw.beforeFuncs[i](rw)
	***REMOVED***
***REMOVED***

func (rw *responseWriter) Flush() ***REMOVED***
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok ***REMOVED***
		if !rw.Written() ***REMOVED***
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		***REMOVED***
		flusher.Flush()
	***REMOVED***
***REMOVED***

type responseWriterCloseNotifer struct ***REMOVED***
	*responseWriter
***REMOVED***

func (rw *responseWriterCloseNotifer) CloseNotify() <-chan bool ***REMOVED***
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
***REMOVED***
