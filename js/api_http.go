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
	"github.com/loadimpact/k6/lib"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"time"
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
		"vu":       a.vu.IDString,
		"status":   "0",
		"method":   method,
		"url":      url,
		"group_id": a.vu.group.ID,
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
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"status":  res.StatusCode,
		"body":    string(resBody),
		"headers": headers,
		"timings": map[string]float64***REMOVED***
			"duration":   float64(trail.Duration) / float64(time.Millisecond),
			"blocked":    float64(trail.Blocked) / float64(time.Millisecond),
			"looking_up": float64(trail.LookingUp) / float64(time.Millisecond),
			"connecting": float64(trail.Connecting) / float64(time.Millisecond),
			"sending":    float64(trail.Sending) / float64(time.Millisecond),
			"waiting":    float64(trail.Waiting) / float64(time.Millisecond),
			"receiving":  float64(trail.Receiving) / float64(time.Millisecond),
		***REMOVED***,
	***REMOVED***
***REMOVED***
