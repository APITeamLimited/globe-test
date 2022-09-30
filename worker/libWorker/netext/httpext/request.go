package httpext

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/Azure/go-ntlmssp"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

// HTTPRequestCookie is a representation of a cookie used for request objects
type HTTPRequestCookie struct ***REMOVED***
	Name, Value string
	Replace     bool
***REMOVED***

// Request represent an http request
type Request struct ***REMOVED***
	Method  string                          `json:"method"`
	URL     string                          `json:"url"`
	Headers map[string][]string             `json:"headers"`
	Body    string                          `json:"body"`
	Cookies map[string][]*HTTPRequestCookie `json:"cookies"`
***REMOVED***

// ParsedHTTPRequest a represantion of a request after it has been parsed from a user script
type ParsedHTTPRequest struct ***REMOVED***
	URL              *URL
	Body             *bytes.Buffer
	Req              *http.Request
	Timeout          time.Duration
	Auth             string
	Throw            bool
	ResponseType     ResponseType
	ResponseCallback func(int) bool
	Compressions     []CompressionType
	Redirects        null.Int
	ActiveJar        *cookiejar.Jar
	Cookies          map[string]*HTTPRequestCookie
	Tags             [][2]string
***REMOVED***

// Matches non-compliant io.Closer implementations (e.g. zstd.Decoder)
type ncloser interface ***REMOVED***
	Close()
***REMOVED***

type readCloser struct ***REMOVED***
	io.Reader
***REMOVED***

// Close readers with differing Close() implementations
func (r readCloser) Close() error ***REMOVED***
	var err error
	switch v := r.Reader.(type) ***REMOVED***
	case io.Closer:
		err = v.Close()
	case ncloser:
		v.Close()
	***REMOVED***
	return err
***REMOVED***

func stdCookiesToHTTPRequestCookies(cookies []*http.Cookie) map[string][]*HTTPRequestCookie ***REMOVED***
	result := make(map[string][]*HTTPRequestCookie, len(cookies))
	for _, cookie := range cookies ***REMOVED***
		result[cookie.Name] = append(result[cookie.Name],
			&HTTPRequestCookie***REMOVED***Name: cookie.Name, Value: cookie.Value***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

// TODO: move as a response method? or constructor?
func updateK6Response(k6Response *Response, finishedReq *finishedRequest) ***REMOVED***
	k6Response.ErrorCode = int(finishedReq.errorCode)
	k6Response.Error = finishedReq.errorMsg
	trail := finishedReq.trail

	if trail.ConnRemoteAddr != nil ***REMOVED***
		remoteHost, remotePortStr, _ := net.SplitHostPort(trail.ConnRemoteAddr.String())
		remotePort, _ := strconv.Atoi(remotePortStr)
		k6Response.RemoteIP = remoteHost
		k6Response.RemotePort = remotePort
	***REMOVED***
	k6Response.Timings = ResponseTimings***REMOVED***
		Duration:       workerMetrics.D(trail.Duration),
		Blocked:        workerMetrics.D(trail.Blocked),
		Connecting:     workerMetrics.D(trail.Connecting),
		TLSHandshaking: workerMetrics.D(trail.TLSHandshaking),
		Sending:        workerMetrics.D(trail.Sending),
		Waiting:        workerMetrics.D(trail.Waiting),
		Receiving:      workerMetrics.D(trail.Receiving),
	***REMOVED***
***REMOVED***

// MakeRequest makes http request for tor the provided ParsedHTTPRequest.
//
// TODO: split apart...
//nolint:cyclop, gocyclo, funlen, gocognit
func MakeRequest(ctx context.Context, state *libWorker.State, preq *ParsedHTTPRequest) (*Response, error) ***REMOVED***
	respReq := &Request***REMOVED***
		Method:  preq.Req.Method,
		URL:     preq.Req.URL.String(),
		Cookies: stdCookiesToHTTPRequestCookies(preq.Req.Cookies()),
		Headers: preq.Req.Header,
	***REMOVED***

	if preq.Body != nil ***REMOVED***
		// TODO: maybe hide this behind of flag in order for this to not happen for big post/puts?
		// should we set this after the compression? what will be the point ?
		respReq.Body = preq.Body.String()

		if len(preq.Compressions) > 0 ***REMOVED***
			compressedBody, contentEncoding, err := compressBody(preq.Compressions, ioutil.NopCloser(preq.Body))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			preq.Body = compressedBody

			currentContentEncoding := preq.Req.Header.Get("Content-Encoding")
			if currentContentEncoding == "" ***REMOVED***
				preq.Req.Header.Set("Content-Encoding", contentEncoding)
			***REMOVED*** else if currentContentEncoding != contentEncoding ***REMOVED***
				state.Logger.Warningf(
					"There's a mismatch between the desired `compression` the manually set `Content-Encoding` header "+
						"in the %s request for '%s', the custom header has precedence and won't be overwritten. "+
						"This may result in invalid data being sent to the server.", preq.Req.Method, preq.Req.URL,
				)
			***REMOVED***
		***REMOVED***

		preq.Req.ContentLength = int64(preq.Body.Len()) // This will make Go set the content-length header
		preq.Req.GetBody = func() (io.ReadCloser, error) ***REMOVED***
			//  using `Bytes()` should reuse the same buffer and as such help with the memory usage. We
			//  should not be writing to it any way so there shouldn't be way to corrupt it (?)
			return ioutil.NopCloser(bytes.NewBuffer(preq.Body.Bytes())), nil
		***REMOVED***
		// as per the documentation using GetBody still requires setting the Body.
		preq.Req.Body, _ = preq.Req.GetBody()
	***REMOVED***

	if contentLengthHeader := preq.Req.Header.Get("Content-Length"); contentLengthHeader != "" ***REMOVED***
		// The content-length header was set by the user, delete it (since Go
		// will set it automatically) and warn if there were differences
		preq.Req.Header.Del("Content-Length")
		length, err := strconv.Atoi(contentLengthHeader)
		if err != nil || preq.Req.ContentLength != int64(length) ***REMOVED***
			state.Logger.Warnf(
				"The specified Content-Length header %q in the %s request for %s "+
					"doesn't match the actual request body length of %d, so it will be ignored!",
				contentLengthHeader, preq.Req.Method, preq.Req.URL, preq.Req.ContentLength,
			)
		***REMOVED***
	***REMOVED***

	tags := state.CloneTags()
	// Override any global tags with request-specific ones.
	for _, tag := range preq.Tags ***REMOVED***
		tags[tag[0]] = tag[1]
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand,
	// and the Name was generated from a tagged template string (via http.url).
	if _, ok := tags["name"]; !ok && state.Options.SystemTags.Has(workerMetrics.TagName) &&
		preq.URL.Name != "" && preq.URL.Name != preq.URL.Clean() ***REMOVED***
		tags["name"] = preq.URL.Name
	***REMOVED***

	// Check rate limit *after* we've prepared a request; no need to wait with that part.
	if rpsLimit := state.RPSLimit; rpsLimit != nil ***REMOVED***
		if err := rpsLimit.Wait(ctx); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	tracerTransport := newTransport(ctx, state, tags, preq.ResponseCallback)
	var transport http.RoundTripper = tracerTransport

	// Combine tags with common log fields
	combinedLogFields := map[string]interface***REMOVED******REMOVED******REMOVED***"source": "http-debug", "vu": state.VUID, "iter": state.Iteration***REMOVED***
	for k, v := range tags ***REMOVED***
		if _, present := combinedLogFields[k]; !present ***REMOVED***
			combinedLogFields[k] = v
		***REMOVED***
	***REMOVED***

	if state.Options.HTTPDebug.String != "" ***REMOVED***
		transport = httpDebugTransport***REMOVED***
			originalTransport: transport,
			httpDebugOption:   state.Options.HTTPDebug.String,
			logger:            state.Logger.WithFields(combinedLogFields),
		***REMOVED***
	***REMOVED***

	if preq.Auth == "digest" ***REMOVED***
		// Until digest authentication is refactored, the first response will always
		// be a 401 error, so we expect that.
		if tracerTransport.responseCallback != nil ***REMOVED***
			originalResponseCallback := tracerTransport.responseCallback
			tracerTransport.responseCallback = func(status int) bool ***REMOVED***
				tracerTransport.responseCallback = originalResponseCallback
				return status == 401
			***REMOVED***
		***REMOVED***
		transport = digestTransport***REMOVED***originalTransport: transport***REMOVED***
	***REMOVED*** else if preq.Auth == "ntlm" ***REMOVED***
		// The first response of NTLM auth may be a 401 error.
		if tracerTransport.responseCallback != nil ***REMOVED***
			originalResponseCallback := tracerTransport.responseCallback
			tracerTransport.responseCallback = func(status int) bool ***REMOVED***
				tracerTransport.responseCallback = originalResponseCallback
				// ntlm is connection-level based so we could've already authorized the connection and to now reuse it
				return status == 401 || originalResponseCallback(status)
			***REMOVED***
		***REMOVED***
		transport = ntlmssp.Negotiator***REMOVED***RoundTripper: transport***REMOVED***
	***REMOVED***

	resp := &Response***REMOVED***URL: preq.URL.URL, Request: respReq***REMOVED***
	client := http.Client***REMOVED***
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			resp.URL = req.URL.String()

			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if preq.ActiveJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					preq.ActiveJar.SetCookies(via[len(via)-1].URL, respCookies)
				***REMOVED***
				req.Header.Del("Cookie")
				SetRequestCookies(req, preq.ActiveJar, preq.Cookies)
			***REMOVED***

			if l := len(via); int64(l) > preq.Redirects.Int64 ***REMOVED***
				if !preq.Redirects.Valid ***REMOVED***
					url := req.URL
					if l > 0 ***REMOVED***
						url = via[0].URL
					***REMOVED***
					state.Logger.WithFields(logrus.Fields***REMOVED***"url": url.String()***REMOVED***).Warnf(
						"Stopped after %d redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED***"+
							" in request params or set global maxRedirects to silence this", l)
				***REMOVED***
				return http.ErrUseLastResponse
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	reqCtx, cancelFunc := context.WithTimeout(ctx, preq.Timeout)
	defer cancelFunc()
	mreq := preq.Req.WithContext(reqCtx)
	res, resErr := client.Do(mreq)

	// TODO(imiric): It would be safer to check for a writeable
	// response body here instead of status code, but those are
	// wrapped in a read-only body when using client timeouts and are
	// unusable until https://github.com/golang/go/issues/31391 is fixed.
	if res != nil && res.StatusCode == http.StatusSwitchingProtocols ***REMOVED***
		_ = res.Body.Close()
		return nil, fmt.Errorf("unsupported response status: %s", res.Status)
	***REMOVED***

	if resErr == nil ***REMOVED***
		resp.Body, resErr = readResponseBody(state, preq.ResponseType, res, resErr)
		if resErr != nil && errors.Is(resErr, context.DeadlineExceeded) ***REMOVED***
			// TODO This can be more specific that the timeout happened in the middle of the reading of the body
			resErr = NewK6Error(requestTimeoutErrorCode, requestTimeoutErrorCodeMsg, resErr)
		***REMOVED***
	***REMOVED***
	finishedReq := tracerTransport.processLastSavedRequest(wrapDecompressionError(resErr))
	if finishedReq != nil ***REMOVED***
		updateK6Response(resp, finishedReq)
	***REMOVED***

	if resErr == nil ***REMOVED***
		if preq.ActiveJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				preq.ActiveJar.SetCookies(res.Request.URL, rc)
			***REMOVED***
		***REMOVED***

		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
		resp.StatusText = res.Status
		resp.Proto = res.Proto

		if res.TLS != nil ***REMOVED***
			resp.setTLSInfo(res.TLS)
		***REMOVED***

		resp.Headers = make(map[string]string, len(res.Header))
		for k, vs := range res.Header ***REMOVED***
			resp.Headers[k] = strings.Join(vs, ", ")
		***REMOVED***

		resCookies := res.Cookies()
		resp.Cookies = make(map[string][]*HTTPCookie, len(resCookies))
		for _, c := range resCookies ***REMOVED***
			resp.Cookies[c.Name] = append(resp.Cookies[c.Name], &HTTPCookie***REMOVED***
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Path:     c.Path,
				HTTPOnly: c.HttpOnly,
				Secure:   c.Secure,
				MaxAge:   c.MaxAge,
				Expires:  c.Expires.UnixNano() / 1000000,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	if resErr != nil ***REMOVED***
		if preq.Throw ***REMOVED*** // if we are going to throw, we shouldn't log it
			return nil, resErr
		***REMOVED***

		// Do *not* log errors about the context being cancelled.
		select ***REMOVED***
		case <-ctx.Done():
		default:
			state.Logger.WithField("error", resErr).Warn("Request Failed")
		***REMOVED***
	***REMOVED***

	return resp, nil
***REMOVED***

// SetRequestCookies sets the cookies of the requests getting those cookies both from the jar and
// from the reqCookies map. The Replace field of the HTTPRequestCookie will be taken into account
func SetRequestCookies(req *http.Request, jar *cookiejar.Jar, reqCookies map[string]*HTTPRequestCookie) ***REMOVED***
	replacedCookies := make(map[string]struct***REMOVED******REMOVED***)
	for key, reqCookie := range reqCookies ***REMOVED***
		req.AddCookie(&http.Cookie***REMOVED***Name: key, Value: reqCookie.Value***REMOVED***)
		if reqCookie.Replace ***REMOVED***
			replacedCookies[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	for _, c := range jar.Cookies(req.URL) ***REMOVED***
		if _, ok := replacedCookies[c.Name]; !ok ***REMOVED***
			req.AddCookie(&http.Cookie***REMOVED***Name: c.Name, Value: c.Value***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***
