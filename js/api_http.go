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

package js

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"sync"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
)

type HTTPResponse struct ***REMOVED***
	Status int
***REMOVED***

type HTTPParams struct ***REMOVED***
	Headers map[string]string `json:"headers"`
	Tags    map[string]string `json:"tags"`
***REMOVED***

func (a JSAPI) HTTPRequest(method, url, body string, paramData string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	bodyReader := io.Reader(nil)
	if body != "" ***REMOVED***
		bodyReader = strings.NewReader(body)
	***REMOVED***
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	var params HTTPParams
	if err := json.Unmarshal([]byte(paramData), &params); err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	for key, value := range params.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	tags := map[string]string***REMOVED***
		"vu":     a.vu.IDString,
		"status": "0",
		"method": method,
		"url":    url,
		"group":  a.vu.group.Path,
	***REMOVED***
	for key, value := range params.Tags ***REMOVED***
		tags[key] = value
	***REMOVED***

	tracer := lib.Tracer***REMOVED******REMOVED***
	res, err := a.vu.HTTPClient.Do(req.WithContext(httptrace.WithClientTrace(a.vu.ctx, tracer.Trace())))
	if err != nil ***REMOVED***
		a.vu.Samples = append(a.vu.Samples, tracer.Done().Samples(tags)...)
		throw(a.vu.vm, err)
	***REMOVED***
	tags["status"] = strconv.Itoa(res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		a.vu.Samples = append(a.vu.Samples, tracer.Done().Samples(tags)...)
		throw(a.vu.vm, err)
	***REMOVED***
	_ = res.Body.Close()

	trail := tracer.Done()
	a.vu.Samples = append(a.vu.Samples, trail.Samples(tags)...)

	headers := make(map[string]string)
	for k, v := range res.Header ***REMOVED***
		headers[k] = strings.Join(v, ", ")
	***REMOVED***
	remoteAddr := strings.Split(trail.ConnRemoteAddr.String(), ":")
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"remote_ip":   remoteAddr[0],
		"remote_port": remoteAddr[1],
		"url":         res.Request.URL.String(),
		"status":      res.StatusCode,
		"body":        string(resBody),
		"headers":     headers,
		"timings": map[string]float64***REMOVED***
			"duration":   stats.D(trail.Duration),
			"blocked":    stats.D(trail.Blocked),
			"looking_up": stats.D(trail.LookingUp),
			"connecting": stats.D(trail.Connecting),
			"sending":    stats.D(trail.Sending),
			"waiting":    stats.D(trail.Waiting),
			"receiving":  stats.D(trail.Receiving),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (a JSAPI) BatchHTTPRequest(requests otto.Value) otto.Value ***REMOVED***
	obj := requests.Object()
	mutex := sync.Mutex***REMOVED******REMOVED***

	keys := obj.Keys()
	errs := make(chan interface***REMOVED******REMOVED***, len(keys))
	for _, key := range keys ***REMOVED***
		v, _ := obj.Get(key)

		var method string
		var url string
		var body string
		var params string

		o := v.Object()

		v, _ = o.Get("method")
		method = v.String()
		v, _ = o.Get("url")
		url = v.String()
		v, _ = o.Get("body")
		body = v.String()
		v, _ = o.Get("params")
		params = v.String()

		go func(tkey string) ***REMOVED***
			defer func() ***REMOVED*** errs <- recover() ***REMOVED***()
			res := a.HTTPRequest(method, url, body, params)

			mutex.Lock()
			_ = obj.Set(tkey, res)
			mutex.Unlock()
		***REMOVED***(key)
	***REMOVED***

	for i := 0; i < len(keys); i++ ***REMOVED***
		if err := <-errs; err != nil ***REMOVED***
			panic(err)
		***REMOVED***
	***REMOVED***

	return obj.Value()
***REMOVED***
