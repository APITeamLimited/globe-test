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
	neturl "net/url"
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

func (http *HTTP) Get(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, HTTP_METHOD_GET, url, args...)
***REMOVED***

func (http *HTTP) Head(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, HTTP_METHOD_HEAD, url, args...)
***REMOVED***

func (http *HTTP) Post(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, HTTP_METHOD_POST, url, args...)
***REMOVED***

func (http *HTTP) Put(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, HTTP_METHOD_PUT, url, args...)
***REMOVED***

func (http *HTTP) Patch(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, HTTP_METHOD_PATCH, url, args...)
***REMOVED***

func (http *HTTP) Del(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, HTTP_METHOD_DELETE, url, args...)
***REMOVED***

func (http *HTTP) Options(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, HTTP_METHOD_OPTIONS, url, args...)
***REMOVED***

func (http *HTTP) Request(ctx context.Context, method string, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	u, err := ToURL(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	res, samples, err := http.request(ctx, rt, state, method, u, args...)
	state.Samples = append(state.Samples, samples...)
	return res, err
***REMOVED***

func (h *HTTP) request(ctx context.Context, rt *goja.Runtime, state *common.State, method string, url URL, args ...goja.Value) (*HTTPResponse, []stats.SampleContainer, error) ***REMOVED***
	var bodyBuf *bytes.Buffer
	var contentType string
	if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) ***REMOVED***
		var data map[string]goja.Value
		if rt.ExportTo(args[0], &data) == nil ***REMOVED***
			// handling multipart request
			if requestContainsFile(data) ***REMOVED***
				bodyBuf = &bytes.Buffer***REMOVED******REMOVED***
				mpw := multipart.NewWriter(bodyBuf)

				// For parameters of type common.FileData, created with open(file, "b"),
				// we write the file boundary to the body buffer.
				// Otherwise parameters are treated as standard form field.
				for k, v := range data ***REMOVED***
					switch ve := v.Export().(type) ***REMOVED***
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
							return nil, nil, err
						***REMOVED***

						if _, err := fw.Write(ve.Data); err != nil ***REMOVED***
							return nil, nil, err
						***REMOVED***
					default:
						fw, err := mpw.CreateFormField(k)
						if err != nil ***REMOVED***
							return nil, nil, err
						***REMOVED***

						if _, err := fw.Write([]byte(v.String())); err != nil ***REMOVED***
							return nil, nil, err
						***REMOVED***
					***REMOVED***
				***REMOVED***

				if err := mpw.Close(); err != nil ***REMOVED***
					return nil, nil, err
				***REMOVED***

				contentType = mpw.FormDataContentType()
			***REMOVED*** else ***REMOVED***
				bodyQuery := make(neturl.Values, len(data))
				for k, v := range data ***REMOVED***
					bodyQuery.Set(k, v.String())
				***REMOVED***
				bodyBuf = bytes.NewBufferString(bodyQuery.Encode())
				contentType = "application/x-www-form-urlencoded"
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			bodyBuf = bytes.NewBufferString(args[0].String())
		***REMOVED***
	***REMOVED***

	req := &http.Request***REMOVED***
		Method: method,
		URL:    url.URL,
		Header: make(http.Header),
	***REMOVED***
	respReq := &HTTPRequest***REMOVED***
		Method: req.Method,
		URL:    req.URL.String(),
	***REMOVED***
	if bodyBuf != nil ***REMOVED***
		req.Body = ioutil.NopCloser(bodyBuf)
		req.ContentLength = int64(bodyBuf.Len())
		respReq.Body = bodyBuf.String()
	***REMOVED***
	if contentType != "" ***REMOVED***
		req.Header.Set("Content-Type", contentType)
	***REMOVED***
	if userAgent := state.Options.UserAgent; userAgent.String != "" ***REMOVED***
		req.Header.Set("User-Agent", userAgent.String)
	***REMOVED***

	tags := state.Options.RunTags.CloneTags()
	if state.Options.SystemTags["method"] ***REMOVED***
		tags["method"] = method
	***REMOVED***
	if state.Options.SystemTags["url"] ***REMOVED***
		tags["url"] = url.URLString
	***REMOVED***
	if state.Options.SystemTags["name"] ***REMOVED***
		tags["name"] = url.Name
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

	redirects := state.Options.MaxRedirects
	timeout := 60 * time.Second
	throw := state.Options.Throw.Bool
	auth := ""

	var activeJar *cookiejar.Jar
	if state.CookieJar != nil ***REMOVED***
		activeJar = state.CookieJar
	***REMOVED***
	reqCookies := make(map[string]*HTTPRequestCookie)

	if len(args) > 1 ***REMOVED***
		paramsV := args[1]
		if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
			params := paramsV.ToObject(rt)
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
						case typeMapKeyStringValueInterface:
							reqCookies[key] = &HTTPRequestCookie***REMOVED***Name: key, Value: "", Replace: false***REMOVED***
							cookie := cookieV.ToObject(rt)
							for _, attr := range cookie.Keys() ***REMOVED***
								switch strings.ToLower(attr) ***REMOVED***
								case "replace":
									reqCookies[key].Replace = cookie.Get(attr).ToBoolean()
								case "value":
									reqCookies[key].Value = cookie.Get(attr).String()
								***REMOVED***
							***REMOVED***
						default:
							reqCookies[key] = &HTTPRequestCookie***REMOVED***Name: key, Value: cookieV.String(), Replace: false***REMOVED***
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
							req.Host = str
						default:
							req.Header.Set(key, str)
						***REMOVED***
					***REMOVED***
				case "jar":
					jarV := params.Get(k)
					if goja.IsUndefined(jarV) || goja.IsNull(jarV) ***REMOVED***
						continue
					***REMOVED***
					switch v := jarV.Export().(type) ***REMOVED***
					case *HTTPCookieJar:
						activeJar = v.jar
					***REMOVED***
				case "redirects":
					redirects = null.IntFrom(params.Get(k).ToInteger())
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
						tags[key] = tagObj.Get(key).String()
					***REMOVED***
				case "auth":
					auth = params.Get(k).String()
				case "timeout":
					timeout = time.Duration(params.Get(k).ToFloat() * float64(time.Millisecond))
				case "throw":
					throw = params.Get(k).ToBoolean()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if activeJar != nil ***REMOVED***
		mergedCookies := h.mergeCookies(req, activeJar, reqCookies)
		respReq.Cookies = mergedCookies
		h.setRequestCookies(req, mergedCookies)
	***REMOVED***

	// Check rate limit *after* we've prepared a request; no need to wait with that part.
	if rpsLimit := state.RPSLimit; rpsLimit != nil ***REMOVED***
		if err := rpsLimit.Wait(ctx); err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***

	respReq.Headers = req.Header

	resp := &HTTPResponse***REMOVED***ctx: ctx, URL: url.URLString, Request: *respReq***REMOVED***
	client := http.Client***REMOVED***
		Transport: state.HTTPTransport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			h.debugResponse(state, req.Response, "RedirectResponse")

			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if activeJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					activeJar.SetCookies(req.URL, respCookies)
				***REMOVED***
				req.Header.Del("Cookie")
				mergedCookies := h.mergeCookies(req, activeJar, reqCookies)

				h.setRequestCookies(req, mergedCookies)
			***REMOVED***

			if l := len(via); int64(l) > redirects.Int64 ***REMOVED***
				if !redirects.Valid ***REMOVED***
					url := req.URL
					if l > 0 ***REMOVED***
						url = via[0].URL
					***REMOVED***
					state.Logger.WithFields(log.Fields***REMOVED***"url": url.String()***REMOVED***).Warnf("Stopped after %d redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED*** in request params or set global maxRedirects to silence this", l)
				***REMOVED***
				return http.ErrUseLastResponse
			***REMOVED*** else ***REMOVED***
				h.debugRequest(state, req, "RedirectRequest")
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	statsSamples := []stats.SampleContainer***REMOVED******REMOVED***
	// if digest authentication option is passed, make an initial request to get the authentication params to compute the authorization header
	if auth == "digest" ***REMOVED***
		username := url.URL.User.Username()
		password, _ := url.URL.User.Password()

		// removing user from URL to avoid sending the authorization header fo basic auth
		req.URL.User = nil

		tracer := netext.Tracer***REMOVED******REMOVED***
		h.debugRequest(state, req, "DigestRequest")
		res, err := client.Do(req.WithContext(netext.WithTracer(ctx, &tracer)))
		h.debugRequest(state, req, "DigestResponse")
		if err != nil ***REMOVED***
			// Do *not* log errors about the contex being cancelled.
			select ***REMOVED***
			case <-ctx.Done():
			default:
				state.Logger.WithField("error", res).Warn("Digest request failed")
			***REMOVED***

			if throw ***REMOVED***
				return nil, nil, err
			***REMOVED***

			resp.Error = err.Error()
			return resp, statsSamples, nil
		***REMOVED***

		if res.StatusCode == http.StatusUnauthorized ***REMOVED***
			body := ""
			if b, err := ioutil.ReadAll(res.Body); err == nil ***REMOVED***
				body = string(b)
			***REMOVED***

			challenge := digest.GetChallengeFromHeader(&res.Header)
			challenge.ComputeResponse(req.Method, req.URL.RequestURI(), body, username, password)
			authorization := challenge.ToAuthorizationStr()
			req.Header.Set(digest.KEY_AUTHORIZATION, authorization)
		***REMOVED***
		trail := tracer.Done()
		trail.SaveSamples(stats.NewSampleTags(tags))
		statsSamples = append(statsSamples, trail)
	***REMOVED***

	if auth == "ntlm" ***REMOVED***
		ctx = netext.WithAuth(ctx, "ntlm")
	***REMOVED***

	tracer := netext.Tracer***REMOVED******REMOVED***
	h.debugRequest(state, req, "Request")
	res, resErr := client.Do(req.WithContext(netext.WithTracer(ctx, &tracer)))
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
		if activeJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				activeJar.SetCookies(res.Request.URL, rc)
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

		if throw ***REMOVED***
			return nil, nil, resErr
		***REMOVED***
	***REMOVED***

	trail.SaveSamples(stats.IntoSampleTags(&tags))
	statsSamples = append(statsSamples, trail)
	return resp, statsSamples, nil
***REMOVED***

func (http *HTTP) Batch(ctx context.Context, reqsV goja.Value) (goja.Value, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	// Return values; retval must be guarded by the mutex.
	var mutex sync.Mutex
	retval := rt.NewObject()
	errs := make(chan error)

	// Concurrency limits.
	globalLimiter := NewSlotLimiter(int(state.Options.Batch.Int64))
	perHostLimiter := NewMultiSlotLimiter(int(state.Options.BatchPerHost.Int64))

	reqs := reqsV.ToObject(rt)
	keys := reqs.Keys()
	for _, k := range keys ***REMOVED***
		k := k
		v := reqs.Get(k)

		method := HTTP_METHOD_GET
		var url URL
		var args []goja.Value

		// Shorthand: "http://example.com/" -> ["GET", "http://example.com/"]
		switch v.ExportType() ***REMOVED***
		case typeURL:
			url = v.Export().(URL)
		case typeString:
			u, err := ToURL(v)
			if err != nil ***REMOVED***
				return goja.Undefined(), err
			***REMOVED***
			url = u
		default:
			obj := v.ToObject(rt)
			objkeys := obj.Keys()
			for _, objk := range objkeys ***REMOVED***
				objv := obj.Get(objk)
				switch objk ***REMOVED***
				case "0", "method":
					method = strings.ToUpper(objv.String())
					if method == HTTP_METHOD_GET || method == HTTP_METHOD_HEAD ***REMOVED***
						args = []goja.Value***REMOVED***goja.Undefined()***REMOVED***
					***REMOVED***
				case "1", "url":
					u, err := ToURL(objv)
					if err != nil ***REMOVED***
						return goja.Undefined(), err
					***REMOVED***
					url = u
				default:
					args = append(args, objv)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		go func() ***REMOVED***
			globalLimiter.Begin()
			defer globalLimiter.End()

			if hl := perHostLimiter.Slot(url.URL.Host); hl != nil ***REMOVED***
				hl.Begin()
				defer hl.End()
			***REMOVED***

			res, samples, err := http.request(ctx, rt, state, method, url, args...)
			if err != nil ***REMOVED***
				errs <- err
				return
			***REMOVED***

			mutex.Lock()
			_ = retval.Set(k, res)
			state.Samples = append(state.Samples, samples...)
			mutex.Unlock()

			errs <- nil
		***REMOVED***()
	***REMOVED***

	var err error
	for range keys ***REMOVED***
		if e := <-errs; e != nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	return retval, err
***REMOVED***

func requestContainsFile(data map[string]goja.Value) bool ***REMOVED***
	for _, v := range data ***REMOVED***
		switch v.Export().(type) ***REMOVED***
		case FileData:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
