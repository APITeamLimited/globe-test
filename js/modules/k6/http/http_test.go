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
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	null "gopkg.in/guregu/null.v3"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

func assertRequestMetricsEmitted(t *testing.T, samples []stats.Sample, method, url string, status int, group string) ***REMOVED***
	seenDuration := false
	seenBlocked := false
	seenConnecting := false
	seenSending := false
	seenWaiting := false
	seenReceiving := false
	for _, sample := range samples ***REMOVED***
		if sample.Tags["url"] == url ***REMOVED***
			switch sample.Metric ***REMOVED***
			case metrics.HTTPReqDuration:
				seenDuration = true
			case metrics.HTTPReqBlocked:
				seenBlocked = true
			case metrics.HTTPReqConnecting:
				seenConnecting = true
			case metrics.HTTPReqSending:
				seenSending = true
			case metrics.HTTPReqWaiting:
				seenWaiting = true
			case metrics.HTTPReqReceiving:
				seenReceiving = true
			***REMOVED***

			assert.Equal(t, strconv.Itoa(status), sample.Tags["status"])
			assert.Equal(t, method, sample.Tags["method"])
			assert.Equal(t, group, sample.Tags["group"])
		***REMOVED***
	***REMOVED***
	assert.True(t, seenDuration, "url %s didn't emit Duration", url)
	assert.True(t, seenBlocked, "url %s didn't emit Blocked", url)
	assert.True(t, seenConnecting, "url %s didn't emit Connecting", url)
	assert.True(t, seenSending, "url %s didn't emit Sending", url)
	assert.True(t, seenWaiting, "url %s didn't emit Waiting", url)
	assert.True(t, seenReceiving, "url %s didn't emit Receiving", url)
***REMOVED***

func TestRequest(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	state := &common.State***REMOVED***
		Options: lib.Options***REMOVED***
			MaxRedirects: null.IntFrom(10),
		***REMOVED***,
		Group: root,
		HTTPTransport: &http.Transport***REMOVED***
			DialContext: (netext.NewDialer(net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***)).DialContext,
		***REMOVED***,
	***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)
	rt.Set("http", common.Bind(rt, &HTTP***REMOVED******REMOVED***, &ctx))

	t.Run("Redirects", func(t *testing.T) ***REMOVED***
		t.Run("9", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.get("https://httpbin.org/redirect/9")`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("10", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.get("https://httpbin.org/redirect/10")`)
			assert.EqualError(t, err, "GoError: Get /get: stopped after 10 redirects")
		***REMOVED***)
	***REMOVED***)
	t.Run("Timeout", func(t *testing.T) ***REMOVED***
		t.Run("10s", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
				http.get("https://httpbin.org/delay/1", ***REMOVED***
					timeout: 5*1000,
				***REMOVED***)
			`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("10s", func(t *testing.T) ***REMOVED***
			startTime := time.Now()
			_, err := common.RunString(rt, `
				http.get("https://httpbin.org/delay/10", ***REMOVED***
					timeout: 1*1000,
				***REMOVED***)
			`)
			endTime := time.Now()
			assert.EqualError(t, err, "GoError: Get https://httpbin.org/delay/10: net/http: request canceled (Client.Timeout exceeded while awaiting headers)")
			assert.WithinDuration(t, startTime.Add(1*time.Second), endTime, 1*time.Second)
		***REMOVED***)
	***REMOVED***)

	t.Run("HTML", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.request("GET", "https://httpbin.org/html");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.indexOf("Herman Melville - Moby-Dick") == -1) ***REMOVED*** throw new Error("wrong body: " + res.body); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/html", 200, "")

		t.Run("html", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			if (res.html().find("h1").text() != "Herman Melville - Moby-Dick") ***REMOVED*** throw new Error("wrong title: " + res.body); ***REMOVED***
			`)
			assert.NoError(t, err)

			t.Run("shorthand", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, `
				if (res.html("h1").text() != "Herman Melville - Moby-Dick") ***REMOVED*** throw new Error("wrong title: " + res.body); ***REMOVED***
				`)
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED***)

		t.Run("group", func(t *testing.T) ***REMOVED***
			g, err := root.Group("my group")
			if assert.NoError(t, err) ***REMOVED***
				old := state.Group
				state.Group = g
				defer func() ***REMOVED*** state.Group = old ***REMOVED***()
			***REMOVED***

			state.Samples = nil
			_, err = common.RunString(rt, `
			let res = http.request("GET", "https://httpbin.org/html");
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.body.indexOf("Herman Melville - Moby-Dick") == -1) ***REMOVED*** throw new Error("wrong body: " + res.body); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/html", 200, "::my group")
		***REMOVED***)
	***REMOVED***)
	t.Run("JSON", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.request("GET", "https://httpbin.org/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.json().args.a != "1") ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
		if (res.json().args.b != "2") ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get?a=1&b=2", 200, "")

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.request("GET", "https://httpbin.org/html").json();`)
			assert.EqualError(t, err, "GoError: invalid character '<' looking for beginning of value")
		***REMOVED***)
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `http.request("", "");`)
		assert.EqualError(t, err, "GoError: Get : unsupported protocol scheme \"\"")
	***REMOVED***)
	t.Run("Unroutable", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `http.request("GET", "http://sdafsgdhfjg/");`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("Params", func(t *testing.T) ***REMOVED***
		for _, literal := range []string***REMOVED***`undefined`, `null`***REMOVED*** ***REMOVED***
			t.Run(literal, func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, fmt.Sprintf(`
				let res = http.request("GET", "https://httpbin.org/headers", null, %s);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`, literal))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", 200, "")
			***REMOVED***)
		***REMOVED***

		t.Run("headers", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				state.Samples = nil
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED*** headers: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`, literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED***
					headers: ***REMOVED*** "X-My-Header": "value" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().headers["X-My-Header"] != "value") ***REMOVED*** throw new Error("wrong X-My-Header: " + res.json().headers["X-My-Header"]); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", 200, "")
			***REMOVED***)
		***REMOVED***)

		t.Run("tags", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				t.Run(literal, func(t *testing.T) ***REMOVED***
					state.Samples = nil
					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED*** tags: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`, literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED*** tags: ***REMOVED*** tag: "value" ***REMOVED*** ***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", 200, "")
				for _, sample := range state.Samples ***REMOVED***
					assert.Equal(t, "value", sample.Tags["tag"])
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("GET", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.get("https://httpbin.org/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.json().args.a != "1") ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
		if (res.json().args.b != "2") ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get?a=1&b=2", 200, "")
	***REMOVED***)
	t.Run("HEAD", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.head("https://httpbin.org/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.length != 0) ***REMOVED*** throw new Error("HEAD responses shouldn't have a body"); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "HEAD", "https://httpbin.org/get?a=1&b=2", 200, "")
	***REMOVED***)

	postMethods := map[string]string***REMOVED***
		"POST":   "post",
		"PUT":    "put",
		"PATCH":  "patch",
		"DELETE": "del",
	***REMOVED***
	for method, fn := range postMethods ***REMOVED***
		t.Run(method, func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, fmt.Sprintf(`
			let res = http.%s("https://httpbin.org/%s", "data");
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.json().data != "data") ***REMOVED*** throw new Error("wrong data: " + res.json().data); ***REMOVED***
			if (res.json().headers["Content-Type"]) ***REMOVED*** throw new Error("content type set: " + res.json().headers["Content-Type"]); ***REMOVED***
			`, fn, strings.ToLower(method)))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), 200, "")

			t.Run("object", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, fmt.Sprintf(`
				let res = http.%s("https://httpbin.org/%s", ***REMOVED***a: "a", b: 2***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().form.a != "a") ***REMOVED*** throw new Error("wrong a=: " + res.json().form.a); ***REMOVED***
				if (res.json().form.b != "2") ***REMOVED*** throw new Error("wrong b=: " + res.json().form.b); ***REMOVED***
				if (res.json().headers["Content-Type"] != "application/x-www-form-urlencoded") ***REMOVED*** throw new Error("wrong content type: " + res.json().headers["Content-Type"]); ***REMOVED***
				`, fn, strings.ToLower(method)))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), 200, "")

				t.Run("Content-Type", func(t *testing.T) ***REMOVED***
					state.Samples = nil
					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.%s("https://httpbin.org/%s", ***REMOVED***a: "a", b: 2***REMOVED***, ***REMOVED***headers: ***REMOVED***"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"***REMOVED******REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					if (res.json().form.a != "a") ***REMOVED*** throw new Error("wrong a=: " + res.json().form.a); ***REMOVED***
					if (res.json().form.b != "2") ***REMOVED*** throw new Error("wrong b=: " + res.json().form.b); ***REMOVED***
					if (res.json().headers["Content-Type"] != "application/x-www-form-urlencoded; charset=utf-8") ***REMOVED*** throw new Error("wrong content type: " + res.json().headers["Content-Type"]); ***REMOVED***
					`, fn, strings.ToLower(method)))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), 200, "")
				***REMOVED***)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***

	t.Run("Batch", func(t *testing.T) ***REMOVED***
		t.Run("GET", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let reqs = [
				["GET", "https://httpbin.org/"],
				["GET", "https://example.com/"],
			];
			let res = http.batch(reqs);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].url != reqs[key][1]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)

			t.Run("Shorthand", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, `
				let reqs = [
					"https://httpbin.org/",
					"https://example.com/",
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
				***REMOVED***`)
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED***)
		t.Run("POST", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let res = http.batch([ ["POST", "https://httpbin.org/post", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("PUT", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let res = http.batch([ ["PUT", "https://httpbin.org/put", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
