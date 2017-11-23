package echo

import (
	"bufio"
	"net"
	"net/http"
)

type (
	// Response wraps an http.ResponseWriter and implements its interface to be used
	// by an HTTP handler to construct an HTTP response.
	// See: https://golang.org/pkg/net/http/#ResponseWriter
	Response struct ***REMOVED***
		echo        *Echo
		beforeFuncs []func()
		Writer      http.ResponseWriter
		Status      int
		Size        int64
		Committed   bool
	***REMOVED***
)

// NewResponse creates a new instance of Response.
func NewResponse(w http.ResponseWriter, e *Echo) (r *Response) ***REMOVED***
	return &Response***REMOVED***Writer: w, echo: e***REMOVED***
***REMOVED***

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *Response) Header() http.Header ***REMOVED***
	return r.Writer.Header()
***REMOVED***

// Before registers a function which is called just before the response is written.
func (r *Response) Before(fn func()) ***REMOVED***
	r.beforeFuncs = append(r.beforeFuncs, fn)
***REMOVED***

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(code int) ***REMOVED***
	if r.Committed ***REMOVED***
		r.echo.Logger.Warn("response already committed")
		return
	***REMOVED***
	for _, fn := range r.beforeFuncs ***REMOVED***
		fn()
	***REMOVED***
	r.Status = code
	r.Writer.WriteHeader(code)
	r.Committed = true
***REMOVED***

// Write writes the data to the connection as part of an HTTP reply.
func (r *Response) Write(b []byte) (n int, err error) ***REMOVED***
	if !r.Committed ***REMOVED***
		r.WriteHeader(http.StatusOK)
	***REMOVED***
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	return
***REMOVED***

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() ***REMOVED***
	r.Writer.(http.Flusher).Flush()
***REMOVED***

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) ***REMOVED***
	return r.Writer.(http.Hijacker).Hijack()
***REMOVED***

// CloseNotify implements the http.CloseNotifier interface to allow detecting
// when the underlying connection has gone away.
// This mechanism can be used to cancel long operations on the server if the
// client has disconnected before the response is ready.
// See [http.CloseNotifier](https://golang.org/pkg/net/http/#CloseNotifier)
func (r *Response) CloseNotify() <-chan bool ***REMOVED***
	return r.Writer.(http.CloseNotifier).CloseNotify()
***REMOVED***

func (r *Response) reset(w http.ResponseWriter) ***REMOVED***
	r.Writer = w
	r.Size = 0
	r.Status = http.StatusOK
	r.Committed = false
***REMOVED***
