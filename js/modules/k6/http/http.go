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
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	neturl "net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules/k6/html"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
)

type HTTPResponseTimings struct ***REMOVED***
	Duration, Blocked, LookingUp, Connecting, Sending, Waiting, Receiving float64
***REMOVED***

type HTTPResponse struct ***REMOVED***
	ctx context.Context

	RemoteIP   string
	RemotePort int
	URL        string
	Status     int
	Headers    map[string]string
	Body       string
	Timings    HTTPResponseTimings
	Error      string

	cachedJSON goja.Value
***REMOVED***

func (res *HTTPResponse) Json() goja.Value ***REMOVED***
	if res.cachedJSON == nil ***REMOVED***
		var v interface***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(res.Body), &v); err != nil ***REMOVED***
			common.Throw(common.GetRuntime(res.ctx), err)
		***REMOVED***
		res.cachedJSON = common.GetRuntime(res.ctx).ToValue(v)
	***REMOVED***
	return res.cachedJSON
***REMOVED***

func (res *HTTPResponse) Html(selector ...string) html.Selection ***REMOVED***
	sel, err := html.HTML***REMOVED******REMOVED***.ParseHTML(res.ctx, res.Body)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(res.ctx), err)
	***REMOVED***
	if len(selector) > 0 ***REMOVED***
		sel = sel.Find(selector[0])
	***REMOVED***
	return sel
***REMOVED***

type HTTP struct***REMOVED******REMOVED***

func (*HTTP) Request(ctx context.Context, method, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

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

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if contentType != "" ***REMOVED***
		req.Header.Set("Content-Type", contentType)
	***REMOVED***
	if userAgent := state.Options.UserAgent; userAgent.Valid ***REMOVED***
		req.Header.Set("User-Agent", userAgent.String)
	***REMOVED***

	tags := map[string]string***REMOVED***
		"status": "0",
		"method": method,
		"url":    url,
		"group":  state.Group.Path,
	***REMOVED***
	timeout := 60 * time.Second
	throw := state.Options.Throw.Bool

	if len(args) > 1 ***REMOVED***
		paramsV := args[1]
		if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
			params := paramsV.ToObject(rt)
			for _, k := range params.Keys() ***REMOVED***
				switch k ***REMOVED***
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

	resp := &HTTPResponse***REMOVED***
		ctx: ctx,
		URL: url,
	***REMOVED***
	client := http.Client***REMOVED***
		Transport: state.HTTPTransport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			max := int(state.Options.MaxRedirects.Int64)
			if len(via) >= max ***REMOVED***
				return errors.Errorf("stopped after %d redirects", max)
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	tracer := netext.Tracer***REMOVED******REMOVED***
	res, resErr := client.Do(req.WithContext(netext.WithTracer(ctx, &tracer)))
	if res != nil ***REMOVED***
		body, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		resp.Body = string(body)
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
		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
		tags["url"] = resp.URL
		tags["status"] = strconv.Itoa(resp.Status)

		resp.Headers = make(map[string]string, len(res.Header))
		for k, vs := range res.Header ***REMOVED***
			resp.Headers[k] = strings.Join(vs, ", ")
		***REMOVED***
	***REMOVED***

	state.Samples = append(state.Samples, trail.Samples(tags)...)
	if resErr != nil ***REMOVED***
		// Do *not* log errors about the contex being cancelled.
		select ***REMOVED***
		case <-ctx.Done():
		default:
			state.Logger.WithField("error", resErr).Warn("Request Failed")
		***REMOVED***

		if throw ***REMOVED***
			return nil, resErr
		***REMOVED***
	***REMOVED***
	return resp, nil
***REMOVED***

func (http *HTTP) Get(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "GET", url, args...)
***REMOVED***

func (http *HTTP) Head(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "HEAD", url, args...)
***REMOVED***

func (http *HTTP) Post(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "POST", url, args...)
***REMOVED***

func (http *HTTP) Put(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PUT", url, args...)
***REMOVED***

func (http *HTTP) Patch(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PATCH", url, args...)
***REMOVED***

func (http *HTTP) Del(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "DELETE", url, args...)
***REMOVED***

func (http *HTTP) Batch(ctx context.Context, reqsV goja.Value) (goja.Value, error) ***REMOVED***
	rt := common.GetRuntime(ctx)

	errs := make(chan error)
	retval := rt.NewObject()
	mutex := sync.Mutex***REMOVED******REMOVED***

	reqs := reqsV.ToObject(rt)
	keys := reqs.Keys()
	for _, k := range keys ***REMOVED***
		k := k
		v := reqs.Get(k)

		var method, url string
		var args []goja.Value

		// Shorthand: "http://example.com/" -> ["GET", "http://example.com/"]
		if v.ExportType().Kind() == reflect.String ***REMOVED***
			method = "GET"
			url = v.String()
		***REMOVED*** else ***REMOVED***
			obj := v.ToObject(rt)
			objkeys := obj.Keys()
			for i, objk := range objkeys ***REMOVED***
				objv := obj.Get(objk)
				switch i ***REMOVED***
				case 0:
					method = strings.ToUpper(objv.String())
					if method == "GET" || method == "HEAD" ***REMOVED***
						args = []goja.Value***REMOVED***goja.Undefined()***REMOVED***
					***REMOVED***
				case 1:
					url = objv.String()
				default:
					args = append(args, objv)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		go func() ***REMOVED***
			res, err := http.Request(ctx, method, url, args...)
			if err != nil ***REMOVED***
				errs <- err
			***REMOVED***
			mutex.Lock()
			_ = retval.Set(k, res)
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
