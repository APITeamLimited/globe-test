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
	"context"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

var (
	typeString                     = reflect.TypeOf("")
	typeURLTag                     = reflect.TypeOf(URLTag***REMOVED******REMOVED***)
	typeMapKeyStringValueInterface = reflect.TypeOf(map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
)

const SSL_3_0 = "ssl3.0"
const TLS_1_0 = "tls1.0"
const TLS_1_1 = "tls1.1"
const TLS_1_2 = "tls1.2"
const OCSP_STATUS_GOOD = "good"
const OCSP_STATUS_REVOKED = "revoked"
const OCSP_STATUS_SERVER_FAILED = "server_failed"
const OCSP_STATUS_UNKNOWN = "unknown"
const OCSP_REASON_UNSPECIFIED = "unspecified"
const OCSP_REASON_KEY_COMPROMISE = "key_compromise"
const OCSP_REASON_CA_COMPROMISE = "ca_compromise"
const OCSP_REASON_AFFILIATION_CHANGED = "affiliation_changed"
const OCSP_REASON_SUPERSEDED = "superseded"
const OCSP_REASON_CESSATION_OF_OPERATION = "cessation_of_operation"
const OCSP_REASON_CERTIFICATE_HOLD = "certificate_hold"
const OCSP_REASON_REMOVE_FROM_CRL = "remove_from_crl"
const OCSP_REASON_PRIVILEGE_WITHDRAWN = "privilege_withdrawn"
const OCSP_REASON_AA_COMPROMISE = "aa_compromise"

type HTTPCookie struct ***REMOVED***
	Name, Value, Domain, Path string
	HttpOnly, Secure          bool
	MaxAge                    int
	Expires                   int64
***REMOVED***

type HTTPRequestCookie struct ***REMOVED***
	Name, Value string
	Replace     bool
***REMOVED***

type HTTP struct ***REMOVED***
	SSL_3_0                            string `js:"SSL_3_0"`
	TLS_1_0                            string `js:"TLS_1_0"`
	TLS_1_1                            string `js:"TLS_1_1"`
	TLS_1_2                            string `js:"TLS_1_2"`
	OCSP_STATUS_GOOD                   string `js:"OCSP_STATUS_GOOD"`
	OCSP_STATUS_REVOKED                string `js:"OCSP_STATUS_REVOKED"`
	OCSP_STATUS_SERVER_FAILED          string `js:"OCSP_STATUS_SERVER_FAILED"`
	OCSP_STATUS_UNKNOWN                string `js:"OCSP_STATUS_UNKNOWN"`
	OCSP_REASON_UNSPECIFIED            string `js:"OCSP_REASON_UNSPECIFIED"`
	OCSP_REASON_KEY_COMPROMISE         string `js:"OCSP_REASON_KEY_COMPROMISE"`
	OCSP_REASON_CA_COMPROMISE          string `js:"OCSP_REASON_CA_COMPROMISE"`
	OCSP_REASON_AFFILIATION_CHANGED    string `js:"OCSP_REASON_AFFILIATION_CHANGED"`
	OCSP_REASON_SUPERSEDED             string `js:"OCSP_REASON_SUPERSEDED"`
	OCSP_REASON_CESSATION_OF_OPERATION string `js:"OCSP_REASON_CESSATION_OF_OPERATION"`
	OCSP_REASON_CERTIFICATE_HOLD       string `js:"OCSP_REASON_CERTIFICATE_HOLD"`
	OCSP_REASON_REMOVE_FROM_CRL        string `js:"OCSP_REASON_REMOVE_FROM_CRL"`
	OCSP_REASON_PRIVILEGE_WITHDRAWN    string `js:"OCSP_REASON_PRIVILEGE_WITHDRAWN"`
	OCSP_REASON_AA_COMPROMISE          string `js:"OCSP_REASON_AA_COMPROMISE"`
***REMOVED***

func New() *HTTP ***REMOVED***
	return &HTTP***REMOVED***
		SSL_3_0:                            SSL_3_0,
		TLS_1_0:                            TLS_1_0,
		TLS_1_1:                            TLS_1_1,
		TLS_1_2:                            TLS_1_2,
		OCSP_STATUS_GOOD:                   OCSP_STATUS_GOOD,
		OCSP_STATUS_REVOKED:                OCSP_STATUS_REVOKED,
		OCSP_STATUS_SERVER_FAILED:          OCSP_STATUS_SERVER_FAILED,
		OCSP_STATUS_UNKNOWN:                OCSP_STATUS_UNKNOWN,
		OCSP_REASON_UNSPECIFIED:            OCSP_REASON_UNSPECIFIED,
		OCSP_REASON_KEY_COMPROMISE:         OCSP_REASON_KEY_COMPROMISE,
		OCSP_REASON_CA_COMPROMISE:          OCSP_REASON_CA_COMPROMISE,
		OCSP_REASON_AFFILIATION_CHANGED:    OCSP_REASON_AFFILIATION_CHANGED,
		OCSP_REASON_SUPERSEDED:             OCSP_REASON_SUPERSEDED,
		OCSP_REASON_CESSATION_OF_OPERATION: OCSP_REASON_CESSATION_OF_OPERATION,
		OCSP_REASON_CERTIFICATE_HOLD:       OCSP_REASON_CERTIFICATE_HOLD,
		OCSP_REASON_REMOVE_FROM_CRL:        OCSP_REASON_REMOVE_FROM_CRL,
		OCSP_REASON_PRIVILEGE_WITHDRAWN:    OCSP_REASON_PRIVILEGE_WITHDRAWN,
		OCSP_REASON_AA_COMPROMISE:          OCSP_REASON_AA_COMPROMISE,
	***REMOVED***
***REMOVED***

func (*HTTP) XCookieJar(ctx *context.Context) *HTTPCookieJar ***REMOVED***
	return newCookieJar(ctx)
***REMOVED***

func (*HTTP) CookieJar(ctx context.Context) *HTTPCookieJar ***REMOVED***
	state := common.GetState(ctx)
	return &HTTPCookieJar***REMOVED***state.CookieJar, &ctx***REMOVED***
***REMOVED***

func (*HTTP) setRequestCookies(req *http.Request, jar *cookiejar.Jar, reqCookies map[string]*HTTPRequestCookie) ***REMOVED***
	jarCookies := make(map[string][]*http.Cookie)
	for _, c := range jar.Cookies(req.URL) ***REMOVED***
		jarCookies[c.Name] = append(jarCookies[c.Name], c)
	***REMOVED***
	for key, reqCookie := range reqCookies ***REMOVED***
		if jc := jarCookies[key]; jc != nil && reqCookie.Replace ***REMOVED***
			jarCookies[key] = []*http.Cookie***REMOVED******REMOVED***Name: key, Value: reqCookie.Value***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			jarCookies[key] = append(jarCookies[key], &http.Cookie***REMOVED***Name: key, Value: reqCookie.Value***REMOVED***)
		***REMOVED***
	***REMOVED***
	for _, cookies := range jarCookies ***REMOVED***
		for _, c := range cookies ***REMOVED***
			req.AddCookie(c)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (h *HTTP) request(ctx context.Context, rt *goja.Runtime, state *common.State, method string, url goja.Value, args ...goja.Value) (*HTTPResponse, []stats.Sample, error) ***REMOVED***
	var bodyReader io.Reader
	var contentType string
	if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) ***REMOVED***
		var data map[string]goja.Value
		if rt.ExportTo(args[0], &data) == nil ***REMOVED***
			bodyQuery := make(neturl.Values, len(data))
			for k, v := range data ***REMOVED***
				bodyQuery.Set(k, v.String())
			***REMOVED***
			bodyReader = bytes.NewBufferString(bodyQuery.Encode())
			contentType = "application/x-www-form-urlencoded"
		***REMOVED*** else ***REMOVED***
			bodyReader = bytes.NewBufferString(args[0].String())
		***REMOVED***
	***REMOVED***

	// The provided URL can be either a string (or at least something stringable) or a URLTag.
	var urlStr string
	var nameTag string
	switch v := url.Export().(type) ***REMOVED***
	case URLTag:
		urlStr = v.URL
		nameTag = v.Name
	default:
		urlStr = url.String()
		nameTag = urlStr
	***REMOVED***

	req, err := http.NewRequest(method, urlStr, bodyReader)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if contentType != "" ***REMOVED***
		req.Header.Set("Content-Type", contentType)
	***REMOVED***
	if userAgent := state.Options.UserAgent; userAgent.Valid ***REMOVED***
		req.Header.Set("User-Agent", userAgent.String)
	***REMOVED***

	tags := map[string]string***REMOVED***
		"proto":  "",
		"status": "0",
		"method": method,
		"url":    urlStr,
		"name":   nameTag,
		"group":  state.Group.Path,
	***REMOVED***
	redirects := state.Options.MaxRedirects
	timeout := 60 * time.Second
	throw := state.Options.Throw.Bool

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
						req.Header.Set(key, headers.Get(key).String())
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
				case "timeout":
					timeout = time.Duration(params.Get(k).ToFloat() * float64(time.Millisecond))
				case "throw":
					throw = params.Get(k).ToBoolean()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if activeJar != nil ***REMOVED***
		h.setRequestCookies(req, activeJar, reqCookies)
	***REMOVED***

	resp := &HTTPResponse***REMOVED***
		ctx: ctx,
		URL: urlStr,
	***REMOVED***
	client := http.Client***REMOVED***
		Transport: state.HTTPTransport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if activeJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					activeJar.SetCookies(req.URL, respCookies)
					h.setRequestCookies(req, activeJar, reqCookies)
				***REMOVED***
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
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	tracer := netext.Tracer***REMOVED******REMOVED***
	res, resErr := client.Do(req.WithContext(netext.WithTracer(ctx, &tracer)))
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
		Duration:   stats.D(trail.Duration),
		Blocked:    stats.D(trail.Blocked),
		Connecting: stats.D(trail.Connecting),
		Sending:    stats.D(trail.Sending),
		Waiting:    stats.D(trail.Waiting),
		Receiving:  stats.D(trail.Receiving),
	***REMOVED***

	if resErr != nil ***REMOVED***
		resp.Error = resErr.Error()
		tags["error"] = resp.Error
	***REMOVED*** else ***REMOVED***
		if activeJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				activeJar.SetCookies(req.URL, rc)
			***REMOVED***
		***REMOVED***

		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
		resp.Proto = res.Proto
		tags["url"] = resp.URL
		tags["status"] = strconv.Itoa(resp.Status)
		tags["proto"] = resp.Proto

		if res.TLS != nil ***REMOVED***
			resp.setTLSInfo(res.TLS)
			tags["tls_version"] = resp.TLSVersion
			tags["ocsp_status"] = resp.OCSP.Status
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
	return resp, trail.Samples(tags), nil
***REMOVED***

func (http *HTTP) Request(ctx context.Context, method string, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	res, samples, err := http.request(ctx, rt, state, method, url, args...)
	state.Samples = append(state.Samples, samples...)
	return res, err
***REMOVED***

func (http *HTTP) Get(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "GET", url, args...)
***REMOVED***

func (http *HTTP) Head(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "HEAD", url, args...)
***REMOVED***

func (http *HTTP) Post(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "POST", url, args...)
***REMOVED***

func (http *HTTP) Put(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PUT", url, args...)
***REMOVED***

func (http *HTTP) Patch(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PATCH", url, args...)
***REMOVED***

func (http *HTTP) Del(ctx context.Context, url goja.Value, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "DELETE", url, args...)
***REMOVED***

func (http *HTTP) Batch(ctx context.Context, reqsV goja.Value) (goja.Value, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	errs := make(chan error)
	retval := rt.NewObject()
	mutex := sync.Mutex***REMOVED******REMOVED***

	reqs := reqsV.ToObject(rt)
	keys := reqs.Keys()
	for _, k := range keys ***REMOVED***
		k := k
		v := reqs.Get(k)

		var method string
		var url goja.Value
		var args []goja.Value

		// Shorthand: "http://example.com/" -> ["GET", "http://example.com/"]
		switch v.ExportType() ***REMOVED***
		case typeString, typeURLTag:
			method = "GET"
			url = v
		default:
			obj := v.ToObject(rt)
			objkeys := obj.Keys()
			for _, objk := range objkeys ***REMOVED***
				objv := obj.Get(objk)
				switch objk ***REMOVED***
				case "0", "method":
					method = strings.ToUpper(objv.String())
					if method == "GET" || method == "HEAD" ***REMOVED***
						args = []goja.Value***REMOVED***goja.Undefined()***REMOVED***
					***REMOVED***
				case "1", "url":
					url = objv
				default:
					args = append(args, objv)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		go func() ***REMOVED***
			res, samples, err := http.request(ctx, rt, state, method, url, args...)
			if err != nil ***REMOVED***
				errs <- err
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

func (http *HTTP) Url(parts []string, pieces ...string) URLTag ***REMOVED***
	var tag URLTag
	for i, part := range parts ***REMOVED***
		tag.Name += part
		tag.URL += part
		if i < len(pieces) ***REMOVED***
			tag.Name += "$***REMOVED******REMOVED***"
			tag.URL += pieces[i]
		***REMOVED***
	***REMOVED***
	return tag
***REMOVED***
