package httpwrap

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"net/http/httptrace"
	"time"
)

// A Tracer uses Go 1.7's new "net/http/httptrace" package to collect detailed request metrics.
// Cheers, love, the cavalry's here.
type Tracer struct ***REMOVED***
	// Duration of the full request.
	Duration time.Duration
	// Time between the start of the request until the first response byte is obtained.
	TimeToFirstByte time.Duration
	// Time between the request is sent and the first byte is obtained.
	TimeWaiting time.Duration

	// Timings for various parts of the request cycle.
	TimeForDNS          time.Duration
	TimeForConnect      time.Duration
	TimeForWriteHeaders time.Duration
	TimeForWriteBody    time.Duration

	// Non-timing related connection info.
	// TODO: Find a way to report this; stats currently only handles float64s.
	ConnAddr     net.Addr
	ConnReused   bool
	ConnWasIdle  bool
	ConnIdleTime time.Duration

	// Reference points.
	startTimeDNS        time.Time
	startTimeConnect    time.Time
	endTimeConnect      time.Time
	endTimeWriteHeaders time.Time
	endTimeWriteBody    time.Time
***REMOVED***

// MakeClientTrace makes a ClientTrace for use with the httptrace package.
func (t *Tracer) MakeClientTrace() httptrace.ClientTrace ***REMOVED***
	return httptrace.ClientTrace***REMOVED***
		DNSStart:             t.dnsStart,
		DNSDone:              t.dnsDone,
		ConnectStart:         t.connectStart,
		ConnectDone:          t.connectDone,
		GotConn:              t.gotConn,
		WroteHeaders:         t.wroteHeaders,
		WroteRequest:         t.wroteRequest,
		GotFirstResponseByte: t.gotFirstResponseByte,
	***REMOVED***
***REMOVED***

// RequestDone tells the tracer that the request has been fully completed, and is needed to fully
// compute timings. Should not be needed in the future: https://github.com/golang/go/issues/16400
func (t *Tracer) RequestDone() ***REMOVED***
	log.Debug("Request Done")
	if t.startTimeConnect.IsZero() ***REMOVED***
		t.Duration = 0
		return
	***REMOVED***
	t.Duration = time.Since(t.startTimeConnect)
***REMOVED***

func (t *Tracer) dnsStart(info httptrace.DNSStartInfo) ***REMOVED***
	log.Debug("DNS Start")
	t.startTimeDNS = time.Now()
***REMOVED***

func (t *Tracer) dnsDone(info httptrace.DNSDoneInfo) ***REMOVED***
	log.Debug("DNS Done")
	t.TimeForDNS = time.Since(t.startTimeDNS)
***REMOVED***

func (t *Tracer) connectStart(network, addr string) ***REMOVED***
	log.Debug("Connect Start")
	// Dual-stack dials will call this multiple times, then discard all but the first successful
	// connection. For our purposes, connection time is the time between the FIRST outgoing
	// connection attempt, to the FIRST successful connection.
	if !t.startTimeConnect.IsZero() ***REMOVED***
		log.Debug("-> Duplicate!")
		return
	***REMOVED***
	t.startTimeConnect = time.Now()
	t.TimeForConnect = 0
***REMOVED***

func (t *Tracer) connectDone(network, addr string, err error) ***REMOVED***
	log.Debug("Connect Done")
	// Discard all but the first successful connection. See ConnectStart().
	if t.TimeForConnect != 0 ***REMOVED***
		log.Debug("-> Duplicate!")
		return
	***REMOVED***
	t.endTimeConnect = time.Now()
	t.TimeForConnect = t.endTimeConnect.Sub(t.startTimeConnect)
***REMOVED***

func (t *Tracer) gotConn(info httptrace.GotConnInfo) ***REMOVED***
	log.Debug("Got Conn")
	if info.Reused ***REMOVED***
		t.startTimeConnect = time.Now()
		t.endTimeConnect = t.startTimeConnect
		t.TimeForConnect = 0
	***REMOVED***

	t.ConnAddr = info.Conn.RemoteAddr()
	t.ConnReused = info.Reused
	t.ConnWasIdle = info.WasIdle
	t.ConnIdleTime = info.IdleTime
***REMOVED***

func (t *Tracer) wroteHeaders() ***REMOVED***
	log.Debug("Wrote Headers")
	t.endTimeWriteHeaders = time.Now()
	t.TimeForWriteHeaders = t.endTimeWriteHeaders.Sub(t.endTimeConnect)
***REMOVED***

func (t *Tracer) wroteRequest(info httptrace.WroteRequestInfo) ***REMOVED***
	log.Debug("Wrote Request")
	t.endTimeWriteBody = time.Now()
	t.TimeForWriteBody = t.endTimeWriteBody.Sub(t.endTimeWriteHeaders)
***REMOVED***

func (t *Tracer) gotFirstResponseByte() ***REMOVED***
	log.Debug("Got First Response Byte")
	t.TimeToFirstByte = time.Since(t.startTimeConnect)
	t.TimeWaiting = time.Since(t.endTimeWriteBody)
***REMOVED***
