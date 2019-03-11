/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package httpext

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	ntlmssp "github.com/Azure/go-ntlmssp"
	digest "github.com/Soontao/goHttpDigestClient"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

// HTTPRequestCookie is a representation of a cookie used for request objects
type HTTPRequestCookie struct ***REMOVED***
	Name, Value string
	Replace     bool
***REMOVED***

// A URL wraps net.URL, and preserves the template (if any) the URL was constructed from.
type URL struct ***REMOVED***
	u    *url.URL
	Name string // http://example.com/thing/$***REMOVED******REMOVED***/
	URL  string // http://example.com/thing/1234/
***REMOVED***

// NewURL returns a new URL for the provided url and name. The error is returned if the url provided
// can't be parsed
func NewURL(urlString, name string) (URL, error) ***REMOVED***
	u, err := url.Parse(urlString)
	return URL***REMOVED***u: u, Name: name, URL: urlString***REMOVED***, err
***REMOVED***

// GetURL returns the internal url.URL
func (u URL) GetURL() *url.URL ***REMOVED***
	return u.u
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
	URL          *URL
	Body         *bytes.Buffer
	Req          *http.Request
	Timeout      time.Duration
	Auth         string
	Throw        bool
	ResponseType ResponseType
	Redirects    null.Int
	ActiveJar    *cookiejar.Jar
	Cookies      map[string]*HTTPRequestCookie
	Tags         map[string]string
***REMOVED***

func stdCookiesToHTTPRequestCookies(cookies []*http.Cookie) map[string][]*HTTPRequestCookie ***REMOVED***
	var result = make(map[string][]*HTTPRequestCookie, len(cookies))
	for _, cookie := range cookies ***REMOVED***
		result[cookie.Name] = append(result[cookie.Name],
			&HTTPRequestCookie***REMOVED***Name: cookie.Name, Value: cookie.Value***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

// MakeRequest makes http request for tor the provided ParsedHTTPRequest
//TODO break this function up
//nolint: gocyclo
func MakeRequest(ctx context.Context, preq *ParsedHTTPRequest) (*Response, error) ***REMOVED***
	state := lib.GetState(ctx)

	respReq := &Request***REMOVED***
		Method:  preq.Req.Method,
		URL:     preq.Req.URL.String(),
		Cookies: stdCookiesToHTTPRequestCookies(preq.Req.Cookies()),
		Headers: preq.Req.Header,
	***REMOVED***
	if preq.Body != nil ***REMOVED***
		respReq.Body = preq.Body.String()
	***REMOVED***

	tags := state.Options.RunTags.CloneTags()
	for k, v := range preq.Tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	if state.Options.SystemTags["method"] ***REMOVED***
		tags["method"] = preq.Req.Method
	***REMOVED***
	if state.Options.SystemTags["url"] ***REMOVED***
		tags["url"] = preq.URL.URL
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := tags["name"]; !ok && state.Options.SystemTags["name"] ***REMOVED***
		tags["name"] = preq.URL.Name
	***REMOVED***
	if state.Options.SystemTags["group"] ***REMOVED***
		tags["group"] = state.Group.Path
	***REMOVED***
	if state.Options.SystemTags["vu"] ***REMOVED***
		tags["vu"] = strconv.FormatInt(state.Vu, 10)
	***REMOVED***
	if state.Options.SystemTags["iter"] ***REMOVED***
		tags["iter"] = strconv.FormatInt(state.Iteration, 10)
	***REMOVED***

	// Check rate limit *after* we've prepared a request; no need to wait with that part.
	if rpsLimit := state.RPSLimit; rpsLimit != nil ***REMOVED***
		if err := rpsLimit.Wait(ctx); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	tracerTransport := newTransport(state.Transport, state.Samples, &state.Options, tags)
	var transport http.RoundTripper = tracerTransport
	if preq.Auth == "ntlm" ***REMOVED***
		transport = ntlmssp.Negotiator***REMOVED***
			RoundTripper: tracerTransport,
		***REMOVED***
	***REMOVED***

	resp := &Response***REMOVED***ctx: ctx, URL: preq.URL.URL, Request: *respReq***REMOVED***
	client := http.Client***REMOVED***
		Transport: transport,
		Timeout:   preq.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			debugResponse(state, req.Response, "RedirectResponse")

			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if preq.ActiveJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					preq.ActiveJar.SetCookies(req.URL, respCookies)
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
					state.Logger.WithFields(log.Fields***REMOVED***"url": url.String()***REMOVED***).Warnf(
						"Stopped after %d redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED***"+
							" in request params or set global maxRedirects to silence this", l)
				***REMOVED***
				return http.ErrUseLastResponse
			***REMOVED***
			debugRequest(state, req, "RedirectRequest")
			return nil
		***REMOVED***,
	***REMOVED***

	// if digest authentication option is passed, make an initial request
	// to get the authentication params to compute the authorization header
	if preq.Auth == "digest" ***REMOVED***
		username := preq.URL.u.User.Username()
		password, _ := preq.URL.u.User.Password()

		// removing user from URL to avoid sending the authorization header fo basic auth
		preq.Req.URL.User = nil

		debugRequest(state, preq.Req, "DigestRequest")
		res, err := client.Do(preq.Req.WithContext(ctx))
		debugRequest(state, preq.Req, "DigestResponse")
		resp.Error = tracerTransport.errorMsg
		resp.ErrorCode = int(tracerTransport.errorCode)
		if err != nil ***REMOVED***
			// Do *not* log errors about the contex being cancelled.
			select ***REMOVED***
			case <-ctx.Done():
			default:
				state.Logger.WithField("error", res).Warn("Digest request failed")
			***REMOVED***

			if preq.Throw || resp.Error == "" ***REMOVED***
				return nil, err
			***REMOVED***

			return resp, nil
		***REMOVED***

		if res.StatusCode == http.StatusUnauthorized ***REMOVED***
			body := ""
			if b, err := ioutil.ReadAll(res.Body); err == nil ***REMOVED***
				body = string(b)
			***REMOVED***

			challenge := digest.GetChallengeFromHeader(&res.Header)
			challenge.ComputeResponse(preq.Req.Method, preq.Req.URL.RequestURI(), body, username, password)
			authorization := challenge.ToAuthorizationStr()
			preq.Req.Header.Set(digest.KEY_AUTHORIZATION, authorization)
		***REMOVED***
	***REMOVED***

	debugRequest(state, preq.Req, "Request")
	res, resErr := client.Do(preq.Req.WithContext(ctx))
	debugResponse(state, res, "Response")
	resp.Error = tracerTransport.errorMsg
	resp.ErrorCode = int(tracerTransport.errorCode)
	if resErr == nil && res != nil ***REMOVED***
		switch res.Header.Get("Content-Encoding") ***REMOVED***
		case "deflate":
			res.Body, resErr = zlib.NewReader(res.Body)
		case "gzip":
			res.Body, resErr = gzip.NewReader(res.Body)
		***REMOVED***
	***REMOVED***
	if resErr == nil && res != nil ***REMOVED***
		if preq.ResponseType == ResponseTypeNone ***REMOVED***
			_, err := io.Copy(ioutil.Discard, res.Body)
			if err != nil && err != io.EOF ***REMOVED***
				resErr = err
			***REMOVED***
			resp.Body = nil
		***REMOVED*** else ***REMOVED***
			// Binary or string
			buf := state.BPool.Get()
			buf.Reset()
			defer state.BPool.Put(buf)
			_, err := io.Copy(buf, res.Body)
			if err != nil && err != io.EOF ***REMOVED***
				resErr = err
			***REMOVED***

			switch preq.ResponseType ***REMOVED***
			case ResponseTypeText:
				resp.Body = buf.String()
			case ResponseTypeBinary:
				resp.Body = buf.Bytes()
			default:
				resErr = fmt.Errorf("unknown responseType %s", preq.ResponseType)
			***REMOVED***
		***REMOVED***
		_ = res.Body.Close()
	***REMOVED***

	trail := tracerTransport.GetTrail()

	if trail.ConnRemoteAddr != nil ***REMOVED***
		remoteHost, remotePortStr, _ := net.SplitHostPort(trail.ConnRemoteAddr.String())
		remotePort, _ := strconv.Atoi(remotePortStr)
		resp.RemoteIP = remoteHost
		resp.RemotePort = remotePort
	***REMOVED***
	resp.Timings = ResponseTimings***REMOVED***
		Duration:       stats.D(trail.Duration),
		Blocked:        stats.D(trail.Blocked),
		Connecting:     stats.D(trail.Connecting),
		TLSHandshaking: stats.D(trail.TLSHandshaking),
		Sending:        stats.D(trail.Sending),
		Waiting:        stats.D(trail.Waiting),
		Receiving:      stats.D(trail.Receiving),
	***REMOVED***

	if resErr == nil ***REMOVED***
		if preq.ActiveJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				preq.ActiveJar.SetCookies(res.Request.URL, rc)
			***REMOVED***
		***REMOVED***

		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
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
		// Do *not* log errors about the contex being cancelled.
		select ***REMOVED***
		case <-ctx.Done():
		default:
			state.Logger.WithField("error", resErr).Warn("Request Failed")
		***REMOVED***

		if preq.Throw || resp.Error == "" ***REMOVED***
			return nil, resErr
		***REMOVED***
	***REMOVED***

	return resp, nil
***REMOVED***

// SetRequestCookies sets the cookies of the requests getting those cookies both from the jar and
// from the reqCookies map. The Replace field of the HTTPRequestCookie will be taken into account
func SetRequestCookies(req *http.Request, jar *cookiejar.Jar, reqCookies map[string]*HTTPRequestCookie) ***REMOVED***
	var replacedCookies = make(map[string]struct***REMOVED******REMOVED***)
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

func debugRequest(state *lib.State, req *http.Request, description string) ***REMOVED***
	if state.Options.HttpDebug.String != "" ***REMOVED***
		dump, err := httputil.DumpRequestOut(req, state.Options.HttpDebug.String == "full")
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		logDump(description, dump)
	***REMOVED***
***REMOVED***
func logDump(description string, dump []byte) ***REMOVED***
	fmt.Printf("%s:\n%s\n", description, dump)
***REMOVED***
