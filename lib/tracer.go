package lib

import (
	"net"
	"net/http/httptrace"
	"time"
)

// A Trail represents detailed information about an HTTP request.
// You'd typically get one from a Tracer.
type Trail struct ***REMOVED***
	// Total request duration, excluding DNS lookup and connect time.
	Duration time.Duration

	Blocked    time.Duration // Waiting to acquire a connection.
	LookingUp  time.Duration // Looking up DNS records.
	Connecting time.Duration // Connecting to remote host.
	Sending    time.Duration // Writing request.
	Waiting    time.Duration // Waiting for first byte.
	Receiving  time.Duration // Receiving response.

	// Detailed connection information.
	ConnReused     bool
	ConnRemoteAddr net.Addr
***REMOVED***

// A Tracer wraps "net/http/httptrace" to collect granular timings for HTTP requests.
// Note that since there is not yet an event for the end of a request (there's a PR to
// add it), you must call Done() at the end of the request to get the full timings.
// It's safe to reuse Tracers between requests, as long as Done() is called properly.
// Cheers, love, the cavalry's here.
type Tracer struct ***REMOVED***
	getConn              time.Time
	gotConn              time.Time
	gotFirstResponseByte time.Time
	dnsStart             time.Time
	dnsDone              time.Time
	connectStart         time.Time
	connectDone          time.Time
	wroteRequest         time.Time

	connReused     bool
	connRemoteAddr net.Addr
***REMOVED***

// Trace() returns a premade ClientTrace that calls all of the Tracer's hooks.
func (t *Tracer) Trace() *httptrace.ClientTrace ***REMOVED***
	return &httptrace.ClientTrace***REMOVED***
		GetConn:              t.GetConn,
		GotConn:              t.GotConn,
		GotFirstResponseByte: t.GotFirstResponseByte,
		DNSStart:             t.DNSStart,
		DNSDone:              t.DNSDone,
		ConnectStart:         t.ConnectStart,
		ConnectDone:          t.ConnectDone,
		WroteRequest:         t.WroteRequest,
	***REMOVED***
***REMOVED***

// Call when the request is finished. Calculates metrics and resets the tracer.
func (t *Tracer) Done() Trail ***REMOVED***
	done := time.Now()
	trail := Trail***REMOVED***
		Duration:   done.Sub(t.getConn),
		Blocked:    t.gotConn.Sub(t.getConn),
		LookingUp:  t.dnsDone.Sub(t.dnsStart),
		Connecting: t.connectDone.Sub(t.connectStart),
		Sending:    t.wroteRequest.Sub(t.connectDone),
		Waiting:    t.gotFirstResponseByte.Sub(t.wroteRequest),
		Receiving:  done.Sub(t.gotFirstResponseByte),

		ConnReused:     t.connReused,
		ConnRemoteAddr: t.connRemoteAddr,
	***REMOVED***

	*t = Tracer***REMOVED******REMOVED***
	return trail
***REMOVED***

// GetConn event hook.
func (t *Tracer) GetConn(hostPort string) ***REMOVED***
	t.getConn = time.Now()
***REMOVED***

// GotConn event hook.
func (t *Tracer) GotConn(info httptrace.GotConnInfo) ***REMOVED***
	t.gotConn = time.Now()
	t.connReused = info.Reused
	t.connRemoteAddr = info.Conn.RemoteAddr()

	if t.connReused ***REMOVED***
		t.connectStart = t.gotConn
		t.connectDone = t.gotConn
	***REMOVED***
***REMOVED***

// GotFirstResponseByte hook.
func (t *Tracer) GotFirstResponseByte() ***REMOVED***
	t.gotFirstResponseByte = time.Now()
***REMOVED***

// DNSStart hook.
func (t *Tracer) DNSStart(info httptrace.DNSStartInfo) ***REMOVED***
	t.dnsStart = time.Now()
***REMOVED***

// DNSDone hook.
func (t *Tracer) DNSDone(info httptrace.DNSDoneInfo) ***REMOVED***
	t.dnsDone = time.Now()
***REMOVED***

// ConnectStart hook.
func (t *Tracer) ConnectStart(network, addr string) ***REMOVED***
	// If using dual-stack dialing, it's possible to get this multiple times.
	if !t.connectStart.IsZero() ***REMOVED***
		return
	***REMOVED***
	t.connectStart = time.Now()
***REMOVED***

// ConnectDone hook.
func (t *Tracer) ConnectDone(network, addr string, err error) ***REMOVED***
	// If using dual-stack dialing, it's possible to get this multiple times.
	if !t.connectDone.IsZero() ***REMOVED***
		return
	***REMOVED***
	t.connectDone = time.Now()
***REMOVED***

// WroteRequest hook.
func (t *Tracer) WroteRequest(info httptrace.WroteRequestInfo) ***REMOVED***
	t.wroteRequest = time.Now()
***REMOVED***
