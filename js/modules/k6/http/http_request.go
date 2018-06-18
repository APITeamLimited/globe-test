/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package http

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	digest "github.com/Soontao/goHttpDigestClient"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

type HTTPRequest struct ***REMOVED***
	Method  string
	URL     string
	Headers map[string][]string
	Body    string
	Cookies map[string][]*HTTPRequestCookie
***REMOVED***

func (h *HTTP) Get(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return h.Request(ctx, HTTP_METHOD_GET, url, args...)
***REMOVED***

func (h *HTTP) Head(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return h.Request(ctx, HTTP_METHOD_HEAD, url, args...)
***REMOVED***

func (h *HTTP) Post(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return h.Request(ctx, HTTP_METHOD_POST, url, args...)
***REMOVED***

func (h *HTTP) Put(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return h.Request(ctx, HTTP_METHOD_PUT, url, args...)
***REMOVED***

func (h *HTTP) Patch(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return h.Request(ctx, HTTP_METHOD_PATCH, url, args...)
***REMOVED***

func (h *HTTP) Del(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return h.Request(ctx, HTTP_METHOD_DELETE, url, args...)
***REMOVED***

func (h *HTTP) Options(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return h.Request(ctx, HTTP_METHOD_OPTIONS, url, args...)
***REMOVED***

func (h *HTTP) Request(ctx context.Context, method string, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	u, err := ToURL(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var body interface***REMOVED******REMOVED***
	var params goja.Value

	if len(args) > 0 ***REMOVED***
		body = args[0].Export()
	***REMOVED***
	if len(args) > 1 ***REMOVED***
		params = args[1]
	***REMOVED***

	req, err := h.parseRequest(ctx, method, u, body, params)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return h.request(ctx, req)
***REMOVED***

type parsedHTTPRequest struct ***REMOVED***
	url           *URL
	body          *bytes.Buffer
	req           *http.Request
	timeout       time.Duration
	auth          string
	throw         bool
	redirects     null.Int
	activeJar     *cookiejar.Jar
	cookies       map[string]*HTTPRequestCookie
	mergedCookies map[string][]*HTTPRequestCookie
	tags          map[string]string
***REMOVED***

func (h *HTTP) parseRequest(ctx context.Context, method string, reqURL URL, body interface***REMOVED******REMOVED***, params goja.Value) (*parsedHTTPRequest, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	result := &parsedHTTPRequest***REMOVED***
		url: &reqURL,
		req: &http.Request***REMOVED***
			Method: method,
			URL:    reqURL.URL,
			Header: make(http.Header),
		***REMOVED***,
		timeout:   60 * time.Second,
		throw:     state.Options.Throw.Bool,
		redirects: state.Options.MaxRedirects,
		cookies:   make(map[string]*HTTPRequestCookie),
		tags:      make(map[string]string),
	***REMOVED***

	formatFormVal := func(v interface***REMOVED******REMOVED***) string ***REMOVED***
		//TODO: handle/warn about unsupported/nested values
		return fmt.Sprintf("%v", v)
	***REMOVED***

	handleObjectBody := func(data map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
		if !requestContainsFile(data) ***REMOVED***
			bodyQuery := make(url.Values, len(data))
			for k, v := range data ***REMOVED***
				bodyQuery.Set(k, formatFormVal(v))
			***REMOVED***
			result.body = bytes.NewBufferString(bodyQuery.Encode())
			result.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			return nil
		***REMOVED***

		// handling multipart request
		result.body = &bytes.Buffer***REMOVED******REMOVED***
		mpw := multipart.NewWriter(result.body)

		// For parameters of type common.FileData, created with open(file, "b"),
		// we write the file boundary to the body buffer.
		// Otherwise parameters are treated as standard form field.
		for k, v := range data ***REMOVED***
			switch ve := v.(type) ***REMOVED***
			case FileData:
				// writing our own part to handle receiving
				// different content-type than the default application/octet-stream
				h := make(textproto.MIMEHeader)
				escapedFilename := escapeQuotes(ve.Filename)
				h.Set("Content-Disposition",
					fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
						k, escapedFilename))
				h.Set("Content-Type", ve.ContentType)

				// this writer will be closed either by the next part or
				// the call to mpw.Close()
				fw, err := mpw.CreatePart(h)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				if _, err := fw.Write(ve.Data); err != nil ***REMOVED***
					return err
				***REMOVED***
			default:
				fw, err := mpw.CreateFormField(k)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				if _, err := fw.Write([]byte(formatFormVal(v))); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := mpw.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***

		result.req.Header.Set("Content-Type", mpw.FormDataContentType())
		return nil
	***REMOVED***

	if body != nil ***REMOVED***
		switch data := body.(type) ***REMOVED***
		case map[string]goja.Value:
			//TODO: fix forms submission and serialization in k6/html before fixing this..
			newData := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
			for k, v := range data ***REMOVED***
				newData[k] = v.Export()
			***REMOVED***
			if err := handleObjectBody(newData); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case map[string]interface***REMOVED******REMOVED***:
			if err := handleObjectBody(data); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case string:
			result.body = bytes.NewBufferString(data)
		case []byte:
			result.body = bytes.NewBuffer(data)
		default:
			return nil, fmt.Errorf("Unknown request body type %T", body)
		***REMOVED***
	***REMOVED***

	if result.body != nil ***REMOVED***
		result.req.Body = ioutil.NopCloser(result.body)
		result.req.ContentLength = int64(result.body.Len())
	***REMOVED***

	if userAgent := state.Options.UserAgent; userAgent.String != "" ***REMOVED***
		result.req.Header.Set("User-Agent", userAgent.String)
	***REMOVED***

	if state.CookieJar != nil ***REMOVED***
		result.activeJar = state.CookieJar
	***REMOVED***

	// TODO: ditch goja.Value, reflections and Object and use a simple go map and type assertions?
	if params != nil && !goja.IsUndefined(params) && !goja.IsNull(params) ***REMOVED***
		params := params.ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "cookies":
				cookiesV := params.Get(k)
				if goja.IsUndefined(cookiesV) || goja.IsNull(cookiesV) ***REMOVED***
					continue
				***REMOVED***
				cookies := cookiesV.ToObject(rt)
				if cookies == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range cookies.Keys() ***REMOVED***
					cookieV := cookies.Get(key)
					if goja.IsUndefined(cookieV) || goja.IsNull(cookieV) ***REMOVED***
						continue
					***REMOVED***
					switch cookieV.ExportType() ***REMOVED***
					case reflect.TypeOf(map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***):
						result.cookies[key] = &HTTPRequestCookie***REMOVED***Name: key, Value: "", Replace: false***REMOVED***
						cookie := cookieV.ToObject(rt)
						for _, attr := range cookie.Keys() ***REMOVED***
							switch strings.ToLower(attr) ***REMOVED***
							case "replace":
								result.cookies[key].Replace = cookie.Get(attr).ToBoolean()
							case "value":
								result.cookies[key].Value = cookie.Get(attr).String()
							***REMOVED***
						***REMOVED***
					default:
						result.cookies[key] = &HTTPRequestCookie***REMOVED***Name: key, Value: cookieV.String(), Replace: false***REMOVED***
					***REMOVED***
				***REMOVED***
			case "headers":
				headersV := params.Get(k)
				if goja.IsUndefined(headersV) || goja.IsNull(headersV) ***REMOVED***
					continue
				***REMOVED***
				headers := headersV.ToObject(rt)
				if headers == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range headers.Keys() ***REMOVED***
					str := headers.Get(key).String()
					switch strings.ToLower(key) ***REMOVED***
					case "host":
						result.req.Host = str
					default:
						result.req.Header.Set(key, str)
					***REMOVED***
				***REMOVED***
			case "jar":
				jarV := params.Get(k)
				if goja.IsUndefined(jarV) || goja.IsNull(jarV) ***REMOVED***
					continue
				***REMOVED***
				switch v := jarV.Export().(type) ***REMOVED***
				case *HTTPCookieJar:
					result.activeJar = v.jar
				***REMOVED***
			case "redirects":
				result.redirects = null.IntFrom(params.Get(k).ToInteger())
			case "tags":
				tagsV := params.Get(k)
				if goja.IsUndefined(tagsV) || goja.IsNull(tagsV) ***REMOVED***
					continue
				***REMOVED***
				tagObj := tagsV.ToObject(rt)
				if tagObj == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range tagObj.Keys() ***REMOVED***
					result.tags[key] = tagObj.Get(key).String()
				***REMOVED***
			case "auth":
				result.auth = params.Get(k).String()
			case "timeout":
				result.timeout = time.Duration(params.Get(k).ToFloat() * float64(time.Millisecond))
			case "throw":
				result.throw = params.Get(k).ToBoolean()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if result.activeJar != nil ***REMOVED***
		result.mergedCookies = h.mergeCookies(result.req, result.activeJar, result.cookies)
		h.setRequestCookies(result.req, result.mergedCookies)
	***REMOVED***

	return result, nil
***REMOVED***

// request() shouldn't mess with the goja runtime or other thread-unsafe
// things because it's called concurrently by Batch()
func (h *HTTP) request(ctx context.Context, preq *parsedHTTPRequest) (*HTTPResponse, error) ***REMOVED***
	state := common.GetState(ctx)

	respReq := &HTTPRequest***REMOVED***
		Method:  preq.req.Method,
		URL:     preq.req.URL.String(),
		Cookies: preq.mergedCookies,
		Headers: preq.req.Header,
	***REMOVED***
	if preq.body != nil ***REMOVED***
		respReq.Body = preq.body.String()
	***REMOVED***

	tags := state.Options.RunTags.CloneTags()
	for k, v := range preq.tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	if state.Options.SystemTags["method"] ***REMOVED***
		tags["method"] = preq.req.Method
	***REMOVED***
	if state.Options.SystemTags["url"] ***REMOVED***
		tags["url"] = preq.url.URLString
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := tags["name"]; !ok && state.Options.SystemTags["name"] ***REMOVED***
		tags["name"] = preq.url.Name
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

	resp := &HTTPResponse***REMOVED***ctx: ctx, URL: preq.url.URLString, Request: *respReq***REMOVED***
	client := http.Client***REMOVED***
		Transport: state.HTTPTransport,
		Timeout:   preq.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			h.debugResponse(state, req.Response, "RedirectResponse")

			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if preq.activeJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					preq.activeJar.SetCookies(req.URL, respCookies)
				***REMOVED***
				req.Header.Del("Cookie")
				mergedCookies := h.mergeCookies(req, preq.activeJar, preq.cookies)

				h.setRequestCookies(req, mergedCookies)
			***REMOVED***

			if l := len(via); int64(l) > preq.redirects.Int64 ***REMOVED***
				if !preq.redirects.Valid ***REMOVED***
					url := req.URL
					if l > 0 ***REMOVED***
						url = via[0].URL
					***REMOVED***
					state.Logger.WithFields(log.Fields***REMOVED***"url": url.String()***REMOVED***).Warnf("Stopped after %d redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED*** in request params or set global maxRedirects to silence this", l)
				***REMOVED***
				return http.ErrUseLastResponse
			***REMOVED***
			h.debugRequest(state, req, "RedirectRequest")
			return nil
		***REMOVED***,
	***REMOVED***

	// if digest authentication option is passed, make an initial request to get the authentication params to compute the authorization header
	if preq.auth == "digest" ***REMOVED***
		username := preq.url.URL.User.Username()
		password, _ := preq.url.URL.User.Password()

		// removing user from URL to avoid sending the authorization header fo basic auth
		preq.req.URL.User = nil

		tracer := netext.Tracer***REMOVED******REMOVED***
		h.debugRequest(state, preq.req, "DigestRequest")
		res, err := client.Do(preq.req.WithContext(netext.WithTracer(ctx, &tracer)))
		h.debugRequest(state, preq.req, "DigestResponse")
		if err != nil ***REMOVED***
			// Do *not* log errors about the contex being cancelled.
			select ***REMOVED***
			case <-ctx.Done():
			default:
				state.Logger.WithField("error", res).Warn("Digest request failed")
			***REMOVED***

			if preq.throw ***REMOVED***
				return nil, err
			***REMOVED***

			resp.Error = err.Error()
			return resp, nil
		***REMOVED***

		if res.StatusCode == http.StatusUnauthorized ***REMOVED***
			body := ""
			if b, err := ioutil.ReadAll(res.Body); err == nil ***REMOVED***
				body = string(b)
			***REMOVED***

			challenge := digest.GetChallengeFromHeader(&res.Header)
			challenge.ComputeResponse(preq.req.Method, preq.req.URL.RequestURI(), body, username, password)
			authorization := challenge.ToAuthorizationStr()
			preq.req.Header.Set(digest.KEY_AUTHORIZATION, authorization)
		***REMOVED***
		trail := tracer.Done()

		if state.Options.SystemTags["ip"] && trail.ConnRemoteAddr != nil ***REMOVED***
			if ip, _, err := net.SplitHostPort(trail.ConnRemoteAddr.String()); err == nil ***REMOVED***
				tags["ip"] = ip
			***REMOVED***
		***REMOVED***
		trail.SaveSamples(stats.NewSampleTags(tags))
		delete(tags, "ip")
		state.Samples <- trail
	***REMOVED***

	if preq.auth == "ntlm" ***REMOVED***
		ctx = netext.WithAuth(ctx, "ntlm")
	***REMOVED***

	tracer := netext.Tracer***REMOVED******REMOVED***
	h.debugRequest(state, preq.req, "Request")
	res, resErr := client.Do(preq.req.WithContext(netext.WithTracer(ctx, &tracer)))
	h.debugResponse(state, res, "Response")
	if resErr == nil && res != nil ***REMOVED***
		switch res.Header.Get("Content-Encoding") ***REMOVED***
		case "deflate":
			res.Body, resErr = zlib.NewReader(res.Body)
		case "gzip":
			res.Body, resErr = gzip.NewReader(res.Body)
		***REMOVED***
	***REMOVED***
	if resErr == nil && res != nil ***REMOVED***
		buf := state.BPool.Get()
		buf.Reset()
		defer state.BPool.Put(buf)
		_, err := io.Copy(buf, res.Body)
		if err != nil && err != io.EOF ***REMOVED***
			resErr = err
		***REMOVED***
		resp.Body = buf.String()
		_ = res.Body.Close()
	***REMOVED***
	trail := tracer.Done()
	if trail.ConnRemoteAddr != nil ***REMOVED***
		remoteHost, remotePortStr, _ := net.SplitHostPort(trail.ConnRemoteAddr.String())
		remotePort, _ := strconv.Atoi(remotePortStr)
		resp.RemoteIP = remoteHost
		resp.RemotePort = remotePort
	***REMOVED***
	resp.Timings = HTTPResponseTimings***REMOVED***
		Duration:       stats.D(trail.Duration),
		Blocked:        stats.D(trail.Blocked),
		Connecting:     stats.D(trail.Connecting),
		TLSHandshaking: stats.D(trail.TLSHandshaking),
		Sending:        stats.D(trail.Sending),
		Waiting:        stats.D(trail.Waiting),
		Receiving:      stats.D(trail.Receiving),
	***REMOVED***

	if resErr != nil ***REMOVED***
		resp.Error = resErr.Error()
		if state.Options.SystemTags["error"] ***REMOVED***
			tags["error"] = resp.Error
		***REMOVED***

		//TODO: expand/replace this so we can recognize the different non-HTTP
		// errors, probably by using a type switch for resErr
		if state.Options.SystemTags["status"] ***REMOVED***
			tags["status"] = "0"
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if preq.activeJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				preq.activeJar.SetCookies(res.Request.URL, rc)
			***REMOVED***
		***REMOVED***

		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
		resp.Proto = res.Proto

		if state.Options.SystemTags["url"] ***REMOVED***
			tags["url"] = resp.URL
		***REMOVED***
		if state.Options.SystemTags["status"] ***REMOVED***
			tags["status"] = strconv.Itoa(resp.Status)
		***REMOVED***
		if state.Options.SystemTags["proto"] ***REMOVED***
			tags["proto"] = resp.Proto
		***REMOVED***

		if res.TLS != nil ***REMOVED***
			resp.setTLSInfo(res.TLS)
			if state.Options.SystemTags["tls_version"] ***REMOVED***
				tags["tls_version"] = resp.TLSVersion
			***REMOVED***
			if state.Options.SystemTags["ocsp_status"] ***REMOVED***
				tags["ocsp_status"] = resp.OCSP.Status
			***REMOVED***
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
				HttpOnly: c.HttpOnly,
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

		if preq.throw ***REMOVED***
			return nil, resErr
		***REMOVED***
	***REMOVED***

	if state.Options.SystemTags["ip"] && trail.ConnRemoteAddr != nil ***REMOVED***
		if ip, _, err := net.SplitHostPort(trail.ConnRemoteAddr.String()); err == nil ***REMOVED***
			tags["ip"] = ip
		***REMOVED***
	***REMOVED***
	trail.SaveSamples(stats.IntoSampleTags(&tags))
	state.Samples <- trail
	return resp, nil
***REMOVED***

func (h *HTTP) Batch(ctx context.Context, reqsV goja.Value) (goja.Value, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	// Return values; retval must be guarded by the mutex.
	var mutex sync.Mutex
	retval := rt.NewObject()
	errs := make(chan error)

	// Concurrency limits.
	globalLimiter := NewSlotLimiter(int(state.Options.Batch.Int64))
	perHostLimiter := NewMultiSlotLimiter(int(state.Options.BatchPerHost.Int64))

	parseBatchRequest := func(key string, val goja.Value) (result *parsedHTTPRequest, err error) ***REMOVED***
		method := HTTP_METHOD_GET
		var ok bool
		var reqURL URL
		var body interface***REMOVED******REMOVED***
		var params goja.Value

		switch data := val.Export().(type) ***REMOVED***
		case []interface***REMOVED******REMOVED***:
			// Handling of ["GET", "http://example.com/"]
			dataLen := len(data)
			if dataLen < 2 ***REMOVED***
				return nil, fmt.Errorf("Invalid batch request '%#v'", data)
			***REMOVED***
			method, ok = data[0].(string)
			if !ok ***REMOVED***
				return nil, fmt.Errorf("Invalid method type '%#v'", data[0])
			***REMOVED***
			reqURL, err = ToURL(data[1])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if dataLen > 2 ***REMOVED***
				body = data[2]
			***REMOVED***
			if dataLen > 3 ***REMOVED***
				params = rt.ToValue(data[3])
			***REMOVED***

		case map[string]interface***REMOVED******REMOVED***:
			// Handling of ***REMOVED***method: "GET", url: "http://test.loadimpact.com"***REMOVED***
			if murl, ok := data["url"]; !ok ***REMOVED***
				return nil, fmt.Errorf("Batch request %s doesn't have an url key", key)
			***REMOVED*** else if reqURL, err = ToURL(murl); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			body = data["body"] // It's fine if it's missing, the map lookup will return

			if newMethod, ok := data["method"]; ok ***REMOVED***
				if method, ok = newMethod.(string); !ok ***REMOVED***
					return nil, fmt.Errorf("Invalid method type '%#v'", newMethod)
				***REMOVED***
				method = strings.ToUpper(method)
				if method == HTTP_METHOD_GET || method == HTTP_METHOD_HEAD ***REMOVED***
					body = nil
				***REMOVED***
			***REMOVED***

			if p, ok := data["params"]; ok ***REMOVED***
				params = rt.ToValue(p)
			***REMOVED***

		default:
			// Handling of "http://example.com/" or http.url`http://example.com/***REMOVED***$id***REMOVED***`
			reqURL, err = ToURL(data)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		return h.parseRequest(ctx, method, reqURL, body, params)
	***REMOVED***

	reqs := reqsV.ToObject(rt)
	keys := reqs.Keys()
	parsedReqs := map[string]*parsedHTTPRequest***REMOVED******REMOVED***
	for _, key := range keys ***REMOVED***
		parsedReq, err := parseBatchRequest(key, reqs.Get(key))
		if err != nil ***REMOVED***
			return retval, err
		***REMOVED***
		parsedReqs[key] = parsedReq
	***REMOVED***

	for k, pr := range parsedReqs ***REMOVED***
		go func(key string, parsedReq *parsedHTTPRequest) ***REMOVED***
			globalLimiter.Begin()
			defer globalLimiter.End()

			if hl := perHostLimiter.Slot(parsedReq.url.URL.Host); hl != nil ***REMOVED***
				hl.Begin()
				defer hl.End()
			***REMOVED***

			res, err := h.request(ctx, parsedReq)
			if err != nil ***REMOVED***
				errs <- err
				return
			***REMOVED***

			mutex.Lock()
			_ = retval.Set(key, res)
			mutex.Unlock()

			errs <- nil
		***REMOVED***(k, pr)
	***REMOVED***

	var err error
	for range keys ***REMOVED***
		if e := <-errs; e != nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	return retval, err
***REMOVED***

func requestContainsFile(data map[string]interface***REMOVED******REMOVED***) bool ***REMOVED***
	for _, v := range data ***REMOVED***
		switch v.(type) ***REMOVED***
		case FileData:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
