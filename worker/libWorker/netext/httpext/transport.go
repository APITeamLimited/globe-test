package httpext

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

// transport is an implementation of http.RoundTripper that will measure and emit
// different metrics for each roundtrip
type transport struct ***REMOVED***
	ctx              context.Context
	state            *libWorker.State
	tags             map[string]string
	responseCallback func(int) bool

	lastRequest     *unfinishedRequest
	lastRequestLock *sync.Mutex
***REMOVED***

// unfinishedRequest stores the request and the raw result returned from the
// underlying http.RoundTripper, but before its body has been read
type unfinishedRequest struct ***REMOVED***
	ctx      context.Context
	tracer   *Tracer
	request  *http.Request
	response *http.Response
	err      error
***REMOVED***

// finishedRequest is produced once the request has been finalized; it is
// triggered either by a subsequent RoundTrip, or for the last request in the
// chain - by the MakeRequest function manually calling the transport method
// processLastSavedRequest(), after reading the HTTP response body.
type finishedRequest struct ***REMOVED***
	*unfinishedRequest
	trail     *Trail
	tlsInfo   netext.TLSInfo
	errorCode errCode
	errorMsg  string
***REMOVED***

var _ http.RoundTripper = &transport***REMOVED******REMOVED***

// newTransport returns a new http.RoundTripper implementation that wraps around
// the provided state's Transport. It uses a httpext.Tracer to measure all HTTP
// requests made through it and annotates and emits the recorded metric samples
// through the state.Samples channel.
func newTransport(
	ctx context.Context,
	state *libWorker.State,
	tags map[string]string,
	responseCallback func(int) bool,
) *transport ***REMOVED***
	return &transport***REMOVED***
		ctx:              ctx,
		state:            state,
		tags:             tags,
		responseCallback: responseCallback,
		lastRequestLock:  new(sync.Mutex),
	***REMOVED***
***REMOVED***

// Helper method to finish the tracer trail, assemble the tag values and emits
// the metric samples for the supplied unfinished request.
//nolint:nestif,funlen
func (t *transport) measureAndEmitMetrics(unfReq *unfinishedRequest) *finishedRequest ***REMOVED***
	trail := unfReq.tracer.Done()

	tags := map[string]string***REMOVED******REMOVED***
	for k, v := range t.tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	result := &finishedRequest***REMOVED***
		unfinishedRequest: unfReq,
		trail:             trail,
	***REMOVED***

	enabledTags := t.state.Options.SystemTags
	urlEnabled := enabledTags.Has(workerMetrics.TagURL)
	var setName bool
	if _, ok := tags["name"]; !ok && enabledTags.Has(workerMetrics.TagName) ***REMOVED***
		setName = true
	***REMOVED***
	if urlEnabled || setName ***REMOVED***
		cleanURL := URL***REMOVED***u: unfReq.request.URL, URL: unfReq.request.URL.String()***REMOVED***.Clean()
		if urlEnabled ***REMOVED***
			tags["url"] = cleanURL
		***REMOVED***
		if setName ***REMOVED***
			tags["name"] = cleanURL
		***REMOVED***
	***REMOVED***

	if enabledTags.Has(workerMetrics.TagMethod) ***REMOVED***
		tags["method"] = unfReq.request.Method
	***REMOVED***

	if unfReq.err != nil ***REMOVED***
		result.errorCode, result.errorMsg = errorCodeForError(unfReq.err)
		if enabledTags.Has(workerMetrics.TagError) ***REMOVED***
			tags["error"] = result.errorMsg
		***REMOVED***

		if enabledTags.Has(workerMetrics.TagErrorCode) ***REMOVED***
			tags["error_code"] = strconv.Itoa(int(result.errorCode))
		***REMOVED***

		if enabledTags.Has(workerMetrics.TagStatus) ***REMOVED***
			tags["status"] = "0"
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if enabledTags.Has(workerMetrics.TagStatus) ***REMOVED***
			tags["status"] = strconv.Itoa(unfReq.response.StatusCode)
		***REMOVED***
		if unfReq.response.StatusCode >= 400 ***REMOVED***
			if enabledTags.Has(workerMetrics.TagErrorCode) ***REMOVED***
				result.errorCode = errCode(1000 + unfReq.response.StatusCode)
				tags["error_code"] = strconv.Itoa(int(result.errorCode))
			***REMOVED***
		***REMOVED***
		if enabledTags.Has(workerMetrics.TagProto) ***REMOVED***
			tags["proto"] = unfReq.response.Proto
		***REMOVED***

		if unfReq.response.TLS != nil ***REMOVED***
			tlsInfo, oscp := netext.ParseTLSConnState(unfReq.response.TLS)
			if enabledTags.Has(workerMetrics.TagTLSVersion) ***REMOVED***
				tags["tls_version"] = tlsInfo.Version
			***REMOVED***
			if enabledTags.Has(workerMetrics.TagOCSPStatus) ***REMOVED***
				tags["ocsp_status"] = oscp.Status
			***REMOVED***
			result.tlsInfo = tlsInfo
		***REMOVED***
	***REMOVED***
	if enabledTags.Has(workerMetrics.TagIP) && trail.ConnRemoteAddr != nil ***REMOVED***
		if ip, _, err := net.SplitHostPort(trail.ConnRemoteAddr.String()); err == nil ***REMOVED***
			tags["ip"] = ip
		***REMOVED***
	***REMOVED***
	var failed float64
	if t.responseCallback != nil ***REMOVED***
		var statusCode int
		if unfReq.err == nil ***REMOVED***
			statusCode = unfReq.response.StatusCode
		***REMOVED***
		expected := t.responseCallback(statusCode)
		if !expected ***REMOVED***
			failed = 1
		***REMOVED***

		if enabledTags.Has(workerMetrics.TagExpectedResponse) ***REMOVED***
			tags[workerMetrics.TagExpectedResponse.String()] = strconv.FormatBool(expected)
		***REMOVED***
	***REMOVED***

	finalTags := workerMetrics.IntoSampleTags(&tags)
	builtinMetrics := t.state.BuiltinMetrics
	trail.SaveSamples(builtinMetrics, finalTags)
	if t.responseCallback != nil ***REMOVED***
		trail.Failed.Valid = true
		if failed == 1 ***REMOVED***
			trail.Failed.Bool = true
		***REMOVED***
		trail.Samples = append(trail.Samples,
			workerMetrics.Sample***REMOVED***
				Metric: builtinMetrics.HTTPReqFailed, Time: trail.EndTime, Tags: finalTags, Value: failed,
			***REMOVED***,
		)
	***REMOVED***
	workerMetrics.PushIfNotDone(t.ctx, t.state.Samples, trail)

	return result
***REMOVED***

func (t *transport) saveCurrentRequest(currentRequest *unfinishedRequest) ***REMOVED***
	t.lastRequestLock.Lock()
	unprocessedRequest := t.lastRequest
	t.lastRequest = currentRequest
	t.lastRequestLock.Unlock()

	if unprocessedRequest != nil ***REMOVED***
		// This shouldn't happen, since we have one transport per request, but just in case...
		t.state.Logger.Warnf("TracerTransport: unexpected unprocessed request for %s", unprocessedRequest.request.URL)
		t.measureAndEmitMetrics(unprocessedRequest)
	***REMOVED***
***REMOVED***

func (t *transport) processLastSavedRequest(lastErr error) *finishedRequest ***REMOVED***
	t.lastRequestLock.Lock()
	unprocessedRequest := t.lastRequest
	t.lastRequest = nil
	t.lastRequestLock.Unlock()

	if unprocessedRequest != nil ***REMOVED***
		// We don't want to overwrite any previous errors, but if there were
		// none and we (i.e. the MakeRequest() function) have one, save it
		// before we emit the workerMetrics.
		if unprocessedRequest.err == nil && lastErr != nil ***REMOVED***
			unprocessedRequest.err = lastErr
		***REMOVED***

		return t.measureAndEmitMetrics(unprocessedRequest)
	***REMOVED***
	return nil
***REMOVED***

// RoundTrip is the implementation of http.RoundTripper
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	t.processLastSavedRequest(nil)

	ctx := req.Context()
	tracer := &Tracer***REMOVED******REMOVED***
	reqWithTracer := req.WithContext(httptrace.WithClientTrace(ctx, tracer.Trace()))
	resp, err := t.state.Transport.RoundTrip(reqWithTracer)

	var netError net.Error
	if errors.As(err, &netError) && netError.Timeout() ***REMOVED***
		var netOpError *net.OpError
		if errors.As(err, &netOpError) && netOpError.Op == "dial" ***REMOVED***
			err = NewK6Error(tcpDialTimeoutErrorCode, tcpDialTimeoutErrorCodeMsg, netError)
		***REMOVED*** else ***REMOVED***
			err = NewK6Error(requestTimeoutErrorCode, requestTimeoutErrorCodeMsg, netError)
		***REMOVED***
	***REMOVED***

	t.saveCurrentRequest(&unfinishedRequest***REMOVED***
		ctx:      ctx,
		tracer:   tracer,
		request:  req,
		response: resp,
		err:      err,
	***REMOVED***)

	return resp, err
***REMOVED***
