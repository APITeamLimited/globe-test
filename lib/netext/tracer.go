/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package netext

import (
	"crypto/tls"
	"net"
	"net/http/httptrace"
	"sync"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
)

// A Trail represents detailed information about an HTTP request.
// You'd typically get one from a Tracer.
type Trail struct ***REMOVED***
	StartTime time.Time
	EndTime   time.Time

	// Total connect time (Connecting + TLSHandshaking)
	ConnDuration time.Duration

	// Total request duration, excluding DNS lookup and connect time.
	Duration time.Duration

	Blocked        time.Duration // Waiting to acquire a connection.
	Connecting     time.Duration // Connecting to remote host.
	TLSHandshaking time.Duration // Executing TLS handshake.
	Sending        time.Duration // Writing request.
	Waiting        time.Duration // Waiting for first byte.
	Receiving      time.Duration // Receiving response.

	// Detailed connection information.
	ConnReused     bool
	ConnRemoteAddr net.Addr
	Errors         []error

	// Populated by SaveSamples()
	Tags    *stats.SampleTags
	Samples []stats.Sample
***REMOVED***

// SaveSamples populates the Trail's sample slice so they're accesible via GetSamples()
func (tr *Trail) SaveSamples(tags *stats.SampleTags) ***REMOVED***
	tr.Tags = tags
	tr.Samples = []stats.Sample***REMOVED***
		***REMOVED***Metric: metrics.HTTPReqs, Time: tr.EndTime, Tags: tags, Value: 1***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqDuration, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Duration)***REMOVED***,

		***REMOVED***Metric: metrics.HTTPReqBlocked, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Blocked)***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqConnecting, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Connecting)***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqTLSHandshaking, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.TLSHandshaking)***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqSending, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Sending)***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqWaiting, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Waiting)***REMOVED***,
		***REMOVED***Metric: metrics.HTTPReqReceiving, Time: tr.EndTime, Tags: tags, Value: stats.D(tr.Receiving)***REMOVED***,
	***REMOVED***
***REMOVED***

// GetSamples implements the stats.SampleContainer interface.
func (tr *Trail) GetSamples() []stats.Sample ***REMOVED***
	return tr.Samples
***REMOVED***

// GetTags implements the stats.ConnectedSampleContainer interface.
func (tr *Trail) GetTags() *stats.SampleTags ***REMOVED***
	return tr.Tags
***REMOVED***

// GetTime implements the stats.ConnectedSampleContainer interface.
func (tr *Trail) GetTime() time.Time ***REMOVED***
	return tr.EndTime
***REMOVED***

// Ensure that interfaces are implemented correctly
var _ stats.ConnectedSampleContainer = &Trail***REMOVED******REMOVED***

// A Tracer wraps "net/http/httptrace" to collect granular timings for HTTP requests.
// Note that since there is not yet an event for the end of a request (there's a PR to
// add it), you must call Done() at the end of the request to get the full timings.
// It's NOT safe to reuse Tracers between requests.
// Cheers, love, the cavalry's here.
type Tracer struct ***REMOVED***
	getConn              int64
	connectStart         int64
	connectDone          int64
	tlsHandshakeStart    int64
	tlsHandshakeDone     int64
	gotConn              int64
	wroteRequest         int64
	gotFirstResponseByte int64

	connReused     bool
	connRemoteAddr net.Addr

	protoErrorsMutex sync.Mutex
	protoErrors      []error
***REMOVED***

// Trace returns a premade ClientTrace that calls all of the Tracer's hooks.
func (t *Tracer) Trace() *httptrace.ClientTrace ***REMOVED***
	return &httptrace.ClientTrace***REMOVED***
		GetConn:              t.GetConn,
		ConnectStart:         t.ConnectStart,
		ConnectDone:          t.ConnectDone,
		TLSHandshakeStart:    t.TLSHandshakeStart,
		TLSHandshakeDone:     t.TLSHandshakeDone,
		GotConn:              t.GotConn,
		WroteRequest:         t.WroteRequest,
		GotFirstResponseByte: t.GotFirstResponseByte,
	***REMOVED***
***REMOVED***

// Add an error in a thread-safe way
func (t *Tracer) addError(err error) ***REMOVED***
	t.protoErrorsMutex.Lock()
	defer t.protoErrorsMutex.Unlock()
	t.protoErrors = append(t.protoErrors, err)
***REMOVED***

func now() int64 ***REMOVED***
	return time.Now().UnixNano()
***REMOVED***

// GetConn is called before a connection is created or
// retrieved from an idle pool. The hostPort is the
// "host:port" of the target or proxy. GetConn is called even
// if there's already an idle cached connection available.
//
// Keep in mind that GetConn won't be called if a connection
// is reused though, for example when there's a redirect.
// If it's called, it will be called before all other hooks.
func (t *Tracer) GetConn(hostPort string) ***REMOVED***
	t.getConn = now()
***REMOVED***

// ConnectStart is called when a new connection's Dial begins.
// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
// enabled, this may be called multiple times.
//
// If the connection is reused, this won't be called. Otherwise,
// it will be called after GetConn() and before ConnectDone().
func (t *Tracer) ConnectStart(network, addr string) ***REMOVED***
	// If using dual-stack dialing, it's possible to get this
	// multiple times, so the atomic compareAndSwap ensures
	// that only the first call's time is recorded
	atomic.CompareAndSwapInt64(&t.connectStart, 0, now())
***REMOVED***

// ConnectDone is called when a new connection's Dial
// completes. The provided err indicates whether the
// connection completedly successfully.
// If net.Dialer.DualStack ("Happy Eyeballs") support is
// enabled, this may be called multiple times.
//
// If the connection is reused, this won't be called. Otherwise,
// it will be called after ConnectStart() and before either
// TLSHandshakeStart() (for TLS connections) or GotConn().
func (t *Tracer) ConnectDone(network, addr string, err error) ***REMOVED***
	// If using dual-stack dialing, it's possible to get this
	// multiple times, so the atomic compareAndSwap ensures
	// that only the first call's time is recorded
	if err == nil ***REMOVED***
		atomic.CompareAndSwapInt64(&t.connectDone, 0, now())
	***REMOVED*** else ***REMOVED***
		t.addError(err)
	***REMOVED***
***REMOVED***

// TLSHandshakeStart is called when the TLS handshake is started. When
// connecting to a HTTPS site via a HTTP proxy, the handshake happens after
// the CONNECT request is processed by the proxy.
//
// If the connection is reused, this won't be called. Otherwise,
// it will be called after ConnectDone() and before TLSHandshakeDone().
func (t *Tracer) TLSHandshakeStart() ***REMOVED***
	atomic.CompareAndSwapInt64(&t.tlsHandshakeStart, 0, now())
***REMOVED***

// TLSHandshakeDone is called after the TLS handshake with either the
// successful handshake's connection state, or a non-nil error on handshake
// failure.
//
// If the connection is reused, this won't be called. Otherwise,
// it will be called after TLSHandshakeStart() and before GotConn().
// If the request was cancelled, this could be called after the
// RoundTrip() method has returned.
func (t *Tracer) TLSHandshakeDone(state tls.ConnectionState, err error) ***REMOVED***
	if err == nil ***REMOVED***
		atomic.CompareAndSwapInt64(&t.tlsHandshakeDone, 0, now())
	***REMOVED*** else ***REMOVED***
		t.addError(err)
	***REMOVED***
***REMOVED***

// GotConn is called after a successful connection is
// obtained. There is no hook for failure to obtain a
// connection; instead, use the error from Transport.RoundTrip.
//
// This is the fist hook called for reused connections. For new
// connections, it's called either after TLSHandshakeDone()
// (for TLS connections) or after ConnectDone()
func (t *Tracer) GotConn(info httptrace.GotConnInfo) ***REMOVED***
	now := now()

	// This shouldn't be called multiple times so no synchronization here,
	// it's better for the race detector to panic if we're wrong.
	t.gotConn = now
	t.connReused = info.Reused
	t.connRemoteAddr = info.Conn.RemoteAddr()

	// The Go stdlib's http module can start connecting to a remote server, only
	// to abandon that connection even before it was fully established and reuse
	// a recently freed already existing connection.
	// We overwrite the different timestamps here, so the other callbacks don't
	// put incorrect values in them (they use CompareAndSwap)
	_, isConnTLS := info.Conn.(*tls.Conn)
	if info.Reused ***REMOVED***
		atomic.SwapInt64(&t.connectStart, now)
		atomic.SwapInt64(&t.connectDone, now)
		if isConnTLS ***REMOVED***
			atomic.SwapInt64(&t.tlsHandshakeStart, now)
			atomic.SwapInt64(&t.tlsHandshakeDone, now)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// There's a bug in the Go stdlib where an HTTP/2 connection can be reused
		// but the httptrace.GotConnInfo struct will contain a false Reused property...
		// That's probably from a previously made connection that was abandoned and
		// directly put in the connection pool in favor of a just-freed already
		// established connection...
		//
		// Using CompareAndSwap here because the HTTP/2 roundtripper has retries and
		// it's possible this isn't actually the first request attempt...
		atomic.CompareAndSwapInt64(&t.connectStart, 0, now)
		atomic.CompareAndSwapInt64(&t.connectDone, 0, now)
		if isConnTLS ***REMOVED***
			atomic.CompareAndSwapInt64(&t.tlsHandshakeStart, 0, now)
			atomic.CompareAndSwapInt64(&t.tlsHandshakeDone, 0, now)
		***REMOVED***
	***REMOVED***
***REMOVED***

// WroteRequest is called with the result of writing the
// request and any body. It may be called multiple times
// in the case of retried requests.
func (t *Tracer) WroteRequest(info httptrace.WroteRequestInfo) ***REMOVED***
	atomic.StoreInt64(&t.wroteRequest, now())

	if info.Err != nil ***REMOVED***
		t.addError(info.Err)
	***REMOVED***
***REMOVED***

// GotFirstResponseByte is called when the first byte of the response
// headers is available.
// If the request was cancelled, this could be called after the
// RoundTrip() method has returned.
func (t *Tracer) GotFirstResponseByte() ***REMOVED***
	atomic.CompareAndSwapInt64(&t.gotFirstResponseByte, 0, now())
***REMOVED***

// Done calculates all metrics and should be called when the request is finished.
func (t *Tracer) Done() *Trail ***REMOVED***
	done := time.Now()

	trail := Trail***REMOVED***
		ConnReused:     t.connReused,
		ConnRemoteAddr: t.connRemoteAddr,
	***REMOVED***

	if t.gotConn != 0 && t.getConn != 0 && t.gotConn > t.getConn ***REMOVED***
		trail.Blocked = time.Duration(t.gotConn - t.getConn)
	***REMOVED***

	// It's possible for some of the methods of httptrace.ClientTrace to
	// actually be called after the http.Client or http.RoundTripper have
	// already returned our result and we've called Done(). This happens
	// mostly for cancelled requests, but we have to use atomics here as
	// well (or use global Tracer locking) so we can avoid data races.
	connectStart := atomic.LoadInt64(&t.connectStart)
	connectDone := atomic.LoadInt64(&t.connectDone)
	tlsHandshakeStart := atomic.LoadInt64(&t.tlsHandshakeStart)
	tlsHandshakeDone := atomic.LoadInt64(&t.tlsHandshakeDone)
	gotConn := atomic.LoadInt64(&t.gotConn)
	wroteRequest := atomic.LoadInt64(&t.wroteRequest)
	gotFirstResponseByte := atomic.LoadInt64(&t.gotFirstResponseByte)

	if connectDone != 0 && connectStart != 0 ***REMOVED***
		trail.Connecting = time.Duration(connectDone - connectStart)
	***REMOVED***
	if tlsHandshakeDone != 0 && tlsHandshakeStart != 0 ***REMOVED***
		trail.TLSHandshaking = time.Duration(tlsHandshakeDone - tlsHandshakeStart)
	***REMOVED***
	if wroteRequest != 0 ***REMOVED***
		if tlsHandshakeDone != 0 ***REMOVED***
			// If the request was sent over TLS, we need to use
			// TLS Handshake Done time to calculate sending duration
			trail.Sending = time.Duration(wroteRequest - tlsHandshakeDone)
		***REMOVED*** else if connectDone != 0 ***REMOVED***
			// Otherwise, use the end of the normal connection
			trail.Sending = time.Duration(wroteRequest - connectDone)
		***REMOVED*** else ***REMOVED***
			// Finally, this handles the strange HTTP/2 case where the GotConn() hook
			// gets called first, but with Reused=false
			trail.Sending = time.Duration(wroteRequest - gotConn)
		***REMOVED***

		if gotFirstResponseByte != 0 && gotFirstResponseByte > wroteRequest ***REMOVED***
			// For some requests, especially HTTP/2, the server starts responding before the
			// client has finished sending the full request
			trail.Waiting = time.Duration(gotFirstResponseByte - wroteRequest)
		***REMOVED***
	***REMOVED***
	if gotFirstResponseByte != 0 ***REMOVED***
		trail.Receiving = done.Sub(time.Unix(0, gotFirstResponseByte))
	***REMOVED***

	// Calculate total times using adjusted values.
	trail.EndTime = done
	trail.ConnDuration = trail.Connecting + trail.TLSHandshaking
	trail.Duration = trail.Sending + trail.Waiting + trail.Receiving
	trail.StartTime = trail.EndTime.Add(-trail.Duration)

	t.protoErrorsMutex.Lock()
	defer t.protoErrorsMutex.Unlock()
	if len(t.protoErrors) > 0 ***REMOVED***
		trail.Errors = append([]error***REMOVED******REMOVED***, t.protoErrors...)
	***REMOVED***

	return &trail
***REMOVED***
