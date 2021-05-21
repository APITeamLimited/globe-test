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
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/dop251/goja"
	"github.com/klauspost/compress/zstd"
	"github.com/mccutchen/go-httpbin/httpbin"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/compiler"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/stats"
)

// runES6String Runs an ES6 string in the given runtime. Use this rather than writing ES5 in tests.
func runES6String(tb testing.TB, rt *goja.Runtime, src string) (goja.Value, error) ***REMOVED***
	var err error
	c := compiler.New(testutils.NewLogger(tb)) // TODO drop it ? maybe we will drop babel and this will be less needed
	src, _, err = c.Transform(src, "__string__")
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	return rt.RunString(src)
***REMOVED***

func TestRunES6String(t *testing.T) ***REMOVED***
	t.Run("Valid", func(t *testing.T) ***REMOVED***
		_, err := runES6String(t, goja.New(), `let a = 1;`)
		assert.NoError(t, err)
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := runES6String(t, goja.New(), `let a = #;`)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "SyntaxError: __string__: Unexpected character '#' (1:8)\n> 1 | let a = #;\n")
	***REMOVED***)
***REMOVED***

// TODO replace this with the Single version
func assertRequestMetricsEmitted(t *testing.T, sampleContainers []stats.SampleContainer, method, url, name string, status int, group string) ***REMOVED***
	if name == "" ***REMOVED***
		name = url
	***REMOVED***

	seenDuration := false
	seenBlocked := false
	seenConnecting := false
	seenTLSHandshaking := false
	seenSending := false
	seenWaiting := false
	seenReceiving := false
	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			tags := sample.Tags.CloneTags()
			if tags["url"] == url ***REMOVED***
				switch sample.Metric ***REMOVED***
				case metrics.HTTPReqDuration:
					seenDuration = true
				case metrics.HTTPReqBlocked:
					seenBlocked = true
				case metrics.HTTPReqConnecting:
					seenConnecting = true
				case metrics.HTTPReqTLSHandshaking:
					seenTLSHandshaking = true
				case metrics.HTTPReqSending:
					seenSending = true
				case metrics.HTTPReqWaiting:
					seenWaiting = true
				case metrics.HTTPReqReceiving:
					seenReceiving = true
				***REMOVED***

				assert.Equal(t, strconv.Itoa(status), tags["status"])
				assert.Equal(t, method, tags["method"])
				assert.Equal(t, group, tags["group"])
				assert.Equal(t, name, tags["name"])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.True(t, seenDuration, "url %s didn't emit Duration", url)
	assert.True(t, seenBlocked, "url %s didn't emit Blocked", url)
	assert.True(t, seenConnecting, "url %s didn't emit Connecting", url)
	assert.True(t, seenTLSHandshaking, "url %s didn't emit TLSHandshaking", url)
	assert.True(t, seenSending, "url %s didn't emit Sending", url)
	assert.True(t, seenWaiting, "url %s didn't emit Waiting", url)
	assert.True(t, seenReceiving, "url %s didn't emit Receiving", url)
***REMOVED***

func assertRequestMetricsEmittedSingle(t *testing.T, sampleContainer stats.SampleContainer, expectedTags map[string]string, metrics []*stats.Metric, callback func(sample stats.Sample)) ***REMOVED***
	t.Helper()

	metricMap := make(map[string]bool, len(metrics))
	for _, m := range metrics ***REMOVED***
		metricMap[m.Name] = false
	***REMOVED***
	for _, sample := range sampleContainer.GetSamples() ***REMOVED***
		tags := sample.Tags.CloneTags()
		v, ok := metricMap[sample.Metric.Name]
		assert.True(t, ok, "unexpected metric %s", sample.Metric.Name)
		assert.False(t, v, "second metric %s", sample.Metric.Name)
		metricMap[sample.Metric.Name] = true
		assert.EqualValues(t, expectedTags, tags, "%s", tags)
		if callback != nil ***REMOVED***
			callback(sample)
		***REMOVED***
	***REMOVED***
	for k, v := range metricMap ***REMOVED***
		assert.True(t, v, "didn't emit %s", k)
	***REMOVED***
***REMOVED***

func newRuntime(
	t testing.TB,
) (*httpmultibin.HTTPMultiBin, *lib.State, chan stats.SampleContainer, *goja.Runtime, *context.Context) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	options := lib.Options***REMOVED***
		MaxRedirects: null.IntFrom(10),
		UserAgent:    null.StringFrom("TestUserAgent"),
		Throw:        null.BoolFrom(true),
		SystemTags:   &stats.DefaultSystemTagSet,
		Batch:        null.IntFrom(20),
		BatchPerHost: null.IntFrom(20),
		// HTTPDebug:    null.StringFrom("full"),
	***REMOVED***
	samples := make(chan stats.SampleContainer, 1000)

	state := &lib.State***REMOVED***
		Options:   options,
		Logger:    logger,
		Group:     root,
		TLSConfig: tb.TLSClientConfig,
		Transport: tb.HTTPTransport,
		BPool:     bpool.NewBufferPool(1),
		Samples:   samples,
		Tags:      map[string]string***REMOVED***"group": root.Path***REMOVED***,
	***REMOVED***

	ctx := new(context.Context)
	*ctx = lib.WithState(tb.Context, state)
	*ctx = common.WithRuntime(*ctx, rt)
	rt.Set("http", common.Bind(rt, new(GlobalHTTP).NewModuleInstancePerVU(), ctx))

	return tb, state, samples, rt, ctx
***REMOVED***

func TestRequestAndBatch(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, samples, rt, ctx := newRuntime(t)
	sr := tb.Replacer.Replace

	// Handle paths with custom logic
	tb.Mux.HandleFunc("/digest-auth/failure", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(2 * time.Second)
	***REMOVED***))

	t.Run("Redirects", func(t *testing.T) ***REMOVED***
		t.Run("tracing", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.get("HTTPBIN_URL/redirect/9");
			`))
			assert.NoError(t, err)
			bufSamples := stats.GetBufferedSamples(samples)

			reqsCount := 0
			for _, container := range bufSamples ***REMOVED***
				for _, sample := range container.GetSamples() ***REMOVED***
					if sample.Metric.Name == "http_reqs" ***REMOVED***
						reqsCount++
					***REMOVED***
				***REMOVED***
			***REMOVED***

			assert.Equal(t, 10, reqsCount)
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/redirect/9"), sr("HTTPBIN_URL/redirect/9"), 302, "")
			for i := 8; i > 0; i-- ***REMOVED***
				url := sr(fmt.Sprintf("HTTPBIN_URL/relative-redirect/%d", i))
				assertRequestMetricsEmitted(t, bufSamples, "GET", url, url, 302, "")
			***REMOVED***
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get"), sr("HTTPBIN_URL/get"), 200, "")
		***REMOVED***)

		t.Run("10", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`http.get("HTTPBIN_URL/redirect/10")`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("11", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.get("HTTPBIN_URL/redirect/11");
			if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/relative-redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
			`))
			assert.NoError(t, err)

			t.Run("Unset Max", func(t *testing.T) ***REMOVED***
				hook := logtest.NewLocal(state.Logger)
				defer hook.Reset()

				oldOpts := state.Options
				defer func() ***REMOVED*** state.Options = oldOpts ***REMOVED***()
				state.Options.MaxRedirects = null.NewInt(10, false)

				_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/redirect/11");
				if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
				if (res.url != "HTTPBIN_URL/relative-redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
				if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
				`))
				assert.NoError(t, err)

				logEntry := hook.LastEntry()
				if assert.NotNil(t, logEntry) ***REMOVED***
					assert.Equal(t, logrus.WarnLevel, logEntry.Level)
					assert.Equal(t, sr("HTTPBIN_URL/redirect/11"), logEntry.Data["url"])
					assert.Equal(t, "Stopped after 11 redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED*** in request params or set global maxRedirects to silence this", logEntry.Message)
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("requestScopeRedirects", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.get("HTTPBIN_URL/redirect/1", ***REMOVED***redirects: 3***REMOVED***);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/get") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("requestScopeNoRedirects", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.get("HTTPBIN_URL/redirect/1", ***REMOVED***redirects: 0***REMOVED***);
			if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)

		t.Run("post body", func(t *testing.T) ***REMOVED***
			tb.Mux.HandleFunc("/post-redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				require.Equal(t, r.Method, "POST")
				_, _ = io.Copy(ioutil.Discard, r.Body)
				http.Redirect(w, r, sr("HTTPBIN_URL/post"), http.StatusPermanentRedirect)
			***REMOVED***))
			_, err := rt.RunString(sr(`
			var res = http.post("HTTPBIN_URL/post-redirect", "pesho", ***REMOVED***redirects: 1***REMOVED***);

			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/post") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.json().data != "pesho") ***REMOVED*** throw new Error("incorrect data : " + res.json().data) ***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("Timeout", func(t *testing.T) ***REMOVED***
		t.Run("10s", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				http.get("HTTPBIN_URL/delay/1", ***REMOVED***
					timeout: 5*1000,
				***REMOVED***)
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("10s", func(t *testing.T) ***REMOVED***
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			startTime := time.Now()
			_, err := rt.RunString(sr(`
				http.get("HTTPBIN_URL/delay/10", ***REMOVED***
					timeout: 1*1000,
				***REMOVED***)
			`))
			endTime := time.Now()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "request timeout")
			assert.WithinDuration(t, startTime.Add(1*time.Second), endTime, 2*time.Second)

			logEntry := hook.LastEntry()
			assert.Nil(t, logEntry)
		***REMOVED***)
		t.Run("10s", func(t *testing.T) ***REMOVED***
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			startTime := time.Now()
			_, err := rt.RunString(sr(`
				http.get("HTTPBIN_URL/delay/10", ***REMOVED***
					timeout: "1s",
				***REMOVED***)
			`))
			endTime := time.Now()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "request timeout")
			assert.WithinDuration(t, startTime.Add(1*time.Second), endTime, 2*time.Second)

			logEntry := hook.LastEntry()
			assert.Nil(t, logEntry)
		***REMOVED***)
	***REMOVED***)
	t.Run("UserAgent", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			var res = http.get("HTTPBIN_URL/headers");
			var headers = res.json()["headers"];
			if (headers['User-Agent'] != "TestUserAgent") ***REMOVED***
				throw new Error("incorrect user agent: " + headers['User-Agent'])
			***REMOVED***
		`))
		assert.NoError(t, err)

		t.Run("Override", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/headers", ***REMOVED***
					headers: ***REMOVED*** "User-Agent": "OtherUserAgent" ***REMOVED***,
				***REMOVED***);
				var headers = res.json()["headers"];
				if (headers['User-Agent'] != "OtherUserAgent") ***REMOVED***
					throw new Error("incorrect user agent: " + headers['User-Agent'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)

		t.Run("Override empty", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/headers", ***REMOVED***
					headers: ***REMOVED*** "User-Agent": "" ***REMOVED***,
				***REMOVED***);
				var headers = res.json()["headers"]
				if (typeof headers['User-Agent'] !== 'undefined') ***REMOVED***
					throw new Error("not undefined user agent: " +  headers['User-Agent'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)

		t.Run("empty", func(t *testing.T) ***REMOVED***
			oldUserAgent := state.Options.UserAgent
			defer func() ***REMOVED***
				state.Options.UserAgent = oldUserAgent
			***REMOVED***()

			state.Options.UserAgent = null.NewString("", true)
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/headers");
				var headers = res.json()["headers"]
				if (typeof headers['User-Agent'] !== 'undefined') ***REMOVED***
					throw new Error("not undefined user agent: " + headers['User-Agent'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)

		t.Run("default", func(t *testing.T) ***REMOVED***
			oldUserAgent := state.Options.UserAgent
			defer func() ***REMOVED***
				state.Options.UserAgent = oldUserAgent
			***REMOVED***()

			state.Options.UserAgent = null.NewString("Default one", false)
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/headers");
				var headers = res.json()["headers"]
				if (headers['User-Agent'] != "Default one") ***REMOVED***
					throw new Error("incorrect user agent: " + headers['User-Agent'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("Compression", func(t *testing.T) ***REMOVED***
		t.Run("gzip", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPSBIN_IP_URL/gzip");
				if (res.json()['gzipped'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['gzipped'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("deflate", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/deflate");
				if (res.json()['deflated'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['deflated'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("zstd", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPSBIN_IP_URL/zstd");
				if (res.json()['compression'] != 'zstd') ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['compression'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("brotli", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPSBIN_IP_URL/brotli");
				if (res.json()['compression'] != 'br') ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['compression'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("zstd-br", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.get("HTTPSBIN_IP_URL/zstd-br");
				if (res.json()['compression'] != 'zstd, br') ***REMOVED***
					throw new Error("unexpected compression: " + res.json()['compression'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("custom compression", func(t *testing.T) ***REMOVED***
			// We should not try to decode it
			tb.Mux.HandleFunc("/customcompression", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				w.Header().Set("Content-Encoding", "custom")
				_, err := w.Write([]byte(`***REMOVED***"custom": true***REMOVED***`))
				assert.NoError(t, err)
			***REMOVED***))

			_, err := rt.RunString(sr(`
				var res = http.get("HTTPBIN_URL/customcompression");
				if (res.json()["custom"] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.body)
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("CompressionWithAcceptEncodingHeader", func(t *testing.T) ***REMOVED***
		t.Run("gzip", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var params = ***REMOVED*** headers: ***REMOVED*** "Accept-Encoding": "gzip" ***REMOVED*** ***REMOVED***;
				var res = http.get("HTTPBIN_URL/gzip", params);
				if (res.json()['gzipped'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['gzipped'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("deflate", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var params = ***REMOVED*** headers: ***REMOVED*** "Accept-Encoding": "deflate" ***REMOVED*** ***REMOVED***;
				var res = http.get("HTTPBIN_URL/deflate", params);
				if (res.json()['deflated'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['deflated'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("Cancelled", func(t *testing.T) ***REMOVED***
		hook := logtest.NewLocal(state.Logger)
		defer hook.Reset()

		oldctx := *ctx
		newctx, cancel := context.WithCancel(oldctx)
		cancel()
		*ctx = newctx
		defer func() ***REMOVED*** *ctx = oldctx ***REMOVED***()

		_, err := rt.RunString(sr(`http.get("HTTPBIN_URL/get/");`))
		assert.Error(t, err)
		assert.Nil(t, hook.LastEntry())
	***REMOVED***)
	t.Run("HTTP/2", func(t *testing.T) ***REMOVED***
		stats.GetBufferedSamples(samples) // Clean up buffered samples from previous tests
		_, err := rt.RunString(sr(`
		var res = http.request("GET", "HTTP2BIN_URL/get");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
		if (res.proto != "HTTP/2.0") ***REMOVED*** throw new Error("wrong proto: " + res.proto) ***REMOVED***
		`))
		assert.NoError(t, err)

		bufSamples := stats.GetBufferedSamples(samples)
		assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTP2BIN_URL/get"), "", 200, "")
		for _, sampleC := range bufSamples ***REMOVED***
			for _, sample := range sampleC.GetSamples() ***REMOVED***
				proto, ok := sample.Tags.Get("proto")
				assert.True(t, ok)
				assert.Equal(t, "HTTP/2.0", proto)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		hook := logtest.NewLocal(state.Logger)
		defer hook.Reset()

		_, err := rt.RunString(`http.request("", "");`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported protocol scheme")

		logEntry := hook.LastEntry()
		assert.Nil(t, logEntry)

		t.Run("throw=false", func(t *testing.T) ***REMOVED***
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			_, err := rt.RunString(`
				var res = http.request("GET", "some://example.com", null, ***REMOVED*** throw: false ***REMOVED***);
				if (res.error.search('unsupported protocol scheme "some"')  == -1) ***REMOVED***
					throw new Error("wrong error:" + res.error);
				***REMOVED***
				throw new Error("another error");
			`)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "another error")

			logEntry := hook.LastEntry()
			if assert.NotNil(t, logEntry) ***REMOVED***
				assert.Equal(t, logrus.WarnLevel, logEntry.Level)
				assert.Contains(t, logEntry.Data["error"].(error).Error(), "unsupported protocol scheme")
				assert.Equal(t, "Request Failed", logEntry.Message)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("InvalidURL", func(t *testing.T) ***REMOVED***
		t.Parallel()
		js := `http.request("GET", "https:// test.k6.io");`
		t.Run("throw=true", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(js)
			require.Error(t, err)
			assert.Contains(t, err.Error(),
				`invalid URL: parse "https:// test.k6.io": invalid character " " in host name`)
		***REMOVED***)

		t.Run("throw=false", func(t *testing.T) ***REMOVED***
			// Disable throw only on this sub-test
			state.Options.Throw.Bool = false
			defer func() ***REMOVED*** state.Options.Throw.Bool = true ***REMOVED***()

			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			_, err := rt.RunString(js)
			require.NoError(t, err)

			logEntry := hook.LastEntry()
			require.NotNil(t, logEntry)
			assert.Equal(t, logrus.WarnLevel, logEntry.Level)
			assert.Contains(t, logEntry.Data["error"].(error).Error(),
				`invalid URL: parse "https:// test.k6.io": invalid character " " in host name`)
			assert.Equal(t, "Request Failed", logEntry.Message)
		***REMOVED***)
	***REMOVED***)

	t.Run("Unroutable", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`http.request("GET", "http://sdafsgdhfjg/");`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("Params", func(t *testing.T) ***REMOVED***
		for _, literal := range []string***REMOVED***`undefined`, `null`***REMOVED*** ***REMOVED***
			t.Run(literal, func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(fmt.Sprintf(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, %s);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`), literal))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)
		***REMOVED***

		t.Run("cookies", func(t *testing.T) ***REMOVED***
			t.Run("access", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/cookies/set?key=value", null, ***REMOVED*** redirects: 0 ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.cookies.key[0].value); ***REMOVED***
				var props = ["name", "value", "domain", "path", "expires", "max_age", "secure", "http_only"];
				var cookie = res.cookies.key[0];
				for (var i = 0; i < props.length; i++) ***REMOVED***
					if (cookie[props[i]] === undefined) ***REMOVED***
						throw new Error("cookie property not found: " + props[i]);
					***REMOVED***
				***REMOVED***
				if (Object.keys(cookie).length != props.length) ***REMOVED***
					throw new Error("cookie has more properties than expected: " + JSON.stringify(Object.keys(cookie)));
				***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies/set?key=value"), "", 302, "")
			***REMOVED***)

			t.Run("vuJar", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				var res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key2: "value2" ***REMOVED*** ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				if (res.json().key2 != "value2") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key2); ***REMOVED***
				var jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key[0] != "value") ***REMOVED*** throw new Error("wrong cookie value in jar"); ***REMOVED***
				if (jarCookies.key2 != undefined) ***REMOVED*** throw new Error("unexpected cookie in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("requestScope", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key: "value" ***REMOVED*** ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				var jar = http.cookieJar();
				var jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key != undefined) ***REMOVED*** throw new Error("unexpected cookie in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("requestScopeReplace", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				var res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key: ***REMOVED*** value: "replaced", replace: true ***REMOVED*** ***REMOVED*** ***REMOVED***);
				if (res.json().key != "replaced") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				var jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key[0] != "value") ***REMOVED*** throw new Error("wrong cookie value in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("redirect", func(t *testing.T) ***REMOVED***
				t.Run("set cookie after redirect", func(t *testing.T) ***REMOVED***
					// TODO figure out a way to remove this ?
					tb.Mux.HandleFunc("/set-cookie-without-redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
						cookie := http.Cookie***REMOVED***
							Name:   "key-foo",
							Value:  "value-bar",
							Path:   "/",
							Domain: sr("HTTPSBIN_DOMAIN"),
						***REMOVED***

						http.SetCookie(w, &cookie)
						w.WriteHeader(200)
					***REMOVED***))
					cookieJar, err := cookiejar.New(nil)
					require.NoError(t, err)
					state.CookieJar = cookieJar
					_, err = rt.RunString(sr(`
						var res = http.request("GET", "HTTPBIN_URL/redirect-to?url=HTTPSBIN_URL/set-cookie-without-redirect");
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`))
					require.NoError(t, err)

					redirectURL, err := url.Parse(sr("HTTPSBIN_URL"))
					require.NoError(t, err)
					require.Len(t, cookieJar.Cookies(redirectURL), 1)
					require.Equal(t, "key-foo", cookieJar.Cookies(redirectURL)[0].Name)
					require.Equal(t, "value-bar", cookieJar.Cookies(redirectURL)[0].Value)

					assertRequestMetricsEmitted(
						t,
						stats.GetBufferedSamples(samples),
						"GET",
						sr("HTTPSBIN_URL/set-cookie-without-redirect"),
						sr("HTTPSBIN_URL/set-cookie-without-redirect"),
						200,
						"",
					)
				***REMOVED***)
				t.Run("set cookie before redirect", func(t *testing.T) ***REMOVED***
					cookieJar, err := cookiejar.New(nil)
					require.NoError(t, err)
					state.CookieJar = cookieJar
					_, err = rt.RunString(sr(`
						var res = http.request("GET", "HTTPSBIN_URL/cookies/set?key=value");
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`))
					require.NoError(t, err)

					redirectURL, err := url.Parse(sr("HTTPSBIN_URL/cookies"))
					require.NoError(t, err)

					require.Len(t, cookieJar.Cookies(redirectURL), 1)
					require.Equal(t, "key", cookieJar.Cookies(redirectURL)[0].Name)
					require.Equal(t, "value", cookieJar.Cookies(redirectURL)[0].Value)

					assertRequestMetricsEmitted(
						t,
						stats.GetBufferedSamples(samples),
						"GET",
						sr("HTTPSBIN_URL/cookies"),
						sr("HTTPSBIN_URL/cookies"),
						200,
						"",
					)
				***REMOVED***)
				t.Run("set cookie after redirect and before second redirect", func(t *testing.T) ***REMOVED***
					cookieJar, err := cookiejar.New(nil)
					require.NoError(t, err)
					state.CookieJar = cookieJar

					// TODO figure out a way to remove this ?
					tb.Mux.HandleFunc("/set-cookie-and-redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
						cookie := http.Cookie***REMOVED***
							Name:   "key-foo",
							Value:  "value-bar",
							Path:   "/set-cookie-and-redirect",
							Domain: sr("HTTPSBIN_DOMAIN"),
						***REMOVED***

						http.SetCookie(w, &cookie)
						http.Redirect(w, r, sr("HTTPBIN_IP_URL/get"), http.StatusMovedPermanently)
					***REMOVED***))

					_, err = rt.RunString(sr(`
						var res = http.request("GET", "HTTPBIN_IP_URL/redirect-to?url=HTTPSBIN_URL/set-cookie-and-redirect");
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`))
					require.NoError(t, err)

					redirectURL, err := url.Parse(sr("HTTPSBIN_URL/set-cookie-and-redirect"))
					require.NoError(t, err)

					require.Len(t, cookieJar.Cookies(redirectURL), 1)
					require.Equal(t, "key-foo", cookieJar.Cookies(redirectURL)[0].Name)
					require.Equal(t, "value-bar", cookieJar.Cookies(redirectURL)[0].Value)

					for _, cookieLessURL := range []string***REMOVED***"HTTPSBIN_URL", "HTTPBIN_IP_URL/redirect-to", "HTTPBIN_IP_URL/get"***REMOVED*** ***REMOVED***
						redirectURL, err = url.Parse(sr(cookieLessURL))
						require.NoError(t, err)
						require.Empty(t, cookieJar.Cookies(redirectURL))
					***REMOVED***

					assertRequestMetricsEmitted(
						t,
						stats.GetBufferedSamples(samples),
						"GET",
						sr("HTTPBIN_IP_URL/get"),
						sr("HTTPBIN_IP_URL/get"),
						200,
						"",
					)
				***REMOVED***)
			***REMOVED***)

			t.Run("domain", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** domain: "HTTPBIN_DOMAIN" ***REMOVED***);
				var res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("wrong cookie value 1: " + res.json().key);
				***REMOVED***
				jar.set("HTTPBIN_URL/cookies", "key2", "value2", ***REMOVED*** domain: "example.com" ***REMOVED***);
				res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("wrong cookie value 2: " + res.json().key);
				***REMOVED***
				if (res.json().key2 != undefined) ***REMOVED***
					throw new Error("cookie 'key2' unexpectedly found");
				***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("path", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** path: "/cookies" ***REMOVED***);
				var res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.json().key);
				***REMOVED***
				jar.set("HTTPBIN_URL/cookies", "key2", "value2", ***REMOVED*** path: "/some-other-path" ***REMOVED***);
				res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.json().key);
				***REMOVED***
				if (res.json().key2 != undefined) ***REMOVED***
					throw new Error("cookie 'key2' unexpectedly found");
				***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("expires", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** expires: "Sun, 24 Jul 1983 17:01:02 GMT" ***REMOVED***);
				var res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != undefined) ***REMOVED***
					throw new Error("cookie 'key' unexpectedly found");
				***REMOVED***
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** expires: "Sat, 24 Jul 2083 17:01:02 GMT" ***REMOVED***);
				res = http.request("GET", "HTTPBIN_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("cookie 'key' not found");
				***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("secure", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = http.cookieJar();
				jar.set("HTTPSBIN_IP_URL/cookies", "key", "value", ***REMOVED*** secure: true ***REMOVED***);
				var res = http.request("GET", "HTTPSBIN_IP_URL/cookies");
				if (res.json().key != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.json().key);
				***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPSBIN_IP_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("localJar", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = rt.RunString(sr(`
				var jar = new http.CookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				var res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key2: "value2" ***REMOVED***, jar: jar ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				if (res.json().key2 != "value2") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key2); ***REMOVED***
				var jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key[0] != "value") ***REMOVED*** throw new Error("wrong cookie value in jar: " + jarCookies.key[0]); ***REMOVED***
				if (jarCookies.key2 != undefined) ***REMOVED*** throw new Error("unexpected cookie in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)
		***REMOVED***)

		t.Run("auth", func(t *testing.T) ***REMOVED***
			t.Run("basic", func(t *testing.T) ***REMOVED***
				url := sr("http://bob:pass@HTTPBIN_IP:HTTPBIN_PORT/basic-auth/bob/pass")
				urlExpected := sr("http://****:****@HTTPBIN_IP:HTTPBIN_PORT/basic-auth/bob/pass")

				_, err := rt.RunString(fmt.Sprintf(`
				var res = http.request("GET", "%s", null, ***REMOVED******REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`, url))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", urlExpected, urlExpected, 200, "")
			***REMOVED***)
			t.Run("digest", func(t *testing.T) ***REMOVED***
				t.Run("success", func(t *testing.T) ***REMOVED***
					url := sr("http://bob:pass@HTTPBIN_IP:HTTPBIN_PORT/digest-auth/auth/bob/pass")
					urlRaw := sr("HTTPBIN_IP_URL/digest-auth/auth/bob/pass")

					_, err := rt.RunString(fmt.Sprintf(`
					var res = http.request("GET", "%s", null, ***REMOVED*** auth: "digest" ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					if (res.error_code != 0) ***REMOVED*** throw new Error("wrong error code: " + res.error_code); ***REMOVED***
					`, url))
					assert.NoError(t, err)

					sampleContainers := stats.GetBufferedSamples(samples)
					assertRequestMetricsEmitted(t, sampleContainers[0:1], "GET",
						urlRaw, urlRaw, 401, "")
					assertRequestMetricsEmitted(t, sampleContainers[1:2], "GET",
						urlRaw, urlRaw, 200, "")
				***REMOVED***)
				t.Run("failure", func(t *testing.T) ***REMOVED***
					url := sr("http://bob:pass@HTTPBIN_IP:HTTPBIN_PORT/digest-auth/failure")

					_, err := rt.RunString(fmt.Sprintf(`
					var res = http.request("GET", "%s", null, ***REMOVED*** auth: "digest", timeout: 1, throw: false ***REMOVED***);
					`, url))
					assert.NoError(t, err)
				***REMOVED***)
			***REMOVED***)
		***REMOVED***)

		t.Run("headers", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := rt.RunString(fmt.Sprintf(sr(`
					var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** headers: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`), literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "X-My-Header": "value" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().headers["X-My-Header"] != "value") ***REMOVED*** throw new Error("wrong X-My-Header: " + res.json().headers["X-My-Header"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)

			t.Run("Host", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "Host": "HTTPBIN_DOMAIN" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().headers["Host"] != "HTTPBIN_DOMAIN") ***REMOVED*** throw new Error("wrong Host: " + res.json().headers["Host"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)

			t.Run("response_request", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "host": "HTTPBIN_DOMAIN" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.request.headers["Host"] != "HTTPBIN_DOMAIN") ***REMOVED*** throw new Error("wrong Host: " + res.request.headers["Host"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)

			t.Run("differentHost", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var custHost = 'k6.io';
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "host": custHost ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.request.headers["Host"] != custHost) ***REMOVED*** throw new Error("wrong Host: " + res.request.headers["Host"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)
		***REMOVED***)

		t.Run("tags", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := rt.RunString(fmt.Sprintf(sr(`
					var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`), literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("name/none", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
					var res = http.request("GET", "HTTPBIN_URL/headers");
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET",
					sr("HTTPBIN_URL/headers"), sr("HTTPBIN_URL/headers"), 200, "")
			***REMOVED***)

			t.Run("name/request", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
					var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: ***REMOVED*** name: "myReq" ***REMOVED******REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET",
					sr("HTTPBIN_URL/headers"), "myReq", 200, "")
			***REMOVED***)

			t.Run("name/template", func(t *testing.T) ***REMOVED***
				_, err := runES6String(t, rt, "http.get(http.url`"+sr(`HTTPBIN_URL/anything/$***REMOVED***1+1***REMOVED***`)+"`);")
				assert.NoError(t, err)
				// There's no /anything endpoint in the go-httpbin library we're using, hence the 404,
				// but it doesn't matter for this test.
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET",
					sr("HTTPBIN_URL/anything/2"), sr("HTTPBIN_URL/anything/$***REMOVED******REMOVED***"), 404, "")
			***REMOVED***)

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: ***REMOVED*** tag: "value" ***REMOVED*** ***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
				for _, sampleC := range bufSamples ***REMOVED***
					for _, sample := range sampleC.GetSamples() ***REMOVED***
						tagValue, ok := sample.Tags.Get("tag")
						assert.True(t, ok)
						assert.Equal(t, "value", tagValue)
					***REMOVED***
				***REMOVED***
			***REMOVED***)

			t.Run("tags-precedence", func(t *testing.T) ***REMOVED***
				oldTags := state.Tags
				defer func() ***REMOVED*** state.Tags = oldTags ***REMOVED***()
				state.Tags = map[string]string***REMOVED***"runtag1": "val1", "runtag2": "val2"***REMOVED***

				_, err := rt.RunString(sr(`
				var res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: ***REMOVED*** method: "test", name: "myName", runtag1: "fromreq" ***REMOVED*** ***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`))
				assert.NoError(t, err)

				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/headers"), "myName", 200, "")
				for _, sampleC := range bufSamples ***REMOVED***
					for _, sample := range sampleC.GetSamples() ***REMOVED***
						tagValue, ok := sample.Tags.Get("method")
						assert.True(t, ok)
						assert.Equal(t, "GET", tagValue)

						tagValue, ok = sample.Tags.Get("name")
						assert.True(t, ok)
						assert.Equal(t, "myName", tagValue)

						tagValue, ok = sample.Tags.Get("runtag1")
						assert.True(t, ok)
						assert.Equal(t, "fromreq", tagValue)

						tagValue, ok = sample.Tags.Get("runtag2")
						assert.True(t, ok)
						assert.Equal(t, "val2", tagValue)
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("GET", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = http.get("HTTPBIN_URL/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.json().args.a != "1") ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
		if (res.json().args.b != "2") ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get?a=1&b=2"), "", 200, "")

		t.Run("Tagged", func(t *testing.T) ***REMOVED***
			_, err := runES6String(t, rt, `
			var a = "1";
			var b = "2";
			var res = http.get(http.url`+"`"+sr(`HTTPBIN_URL/get?a=$***REMOVED***a***REMOVED***&b=$***REMOVED***b***REMOVED***`)+"`"+`);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.json().args.a != a) ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
			if (res.json().args.b != b) ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get?a=1&b=2"), sr("HTTPBIN_URL/get?a=$***REMOVED******REMOVED***&b=$***REMOVED******REMOVED***"), 200, "")
		***REMOVED***)
	***REMOVED***)
	t.Run("HEAD", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = http.head("HTTPBIN_URL/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.length != 0) ***REMOVED*** throw new Error("HEAD responses shouldn't have a body"); ***REMOVED***
		if (!res.headers["Content-Length"]) ***REMOVED*** throw new Error("Missing or invalid Content-Length header!"); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "HEAD", sr("HTTPBIN_URL/get?a=1&b=2"), "", 200, "")
	***REMOVED***)

	t.Run("OPTIONS", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = http.options("HTTPBIN_URL/?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (!res.headers["Access-Control-Allow-Methods"]) ***REMOVED*** throw new Error("Missing Access-Control-Allow-Methods header!"); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "OPTIONS", sr("HTTPBIN_URL/?a=1&b=2"), "", 200, "")
	***REMOVED***)

	// DELETE HTTP requests shouldn't usually send a request body, they should use url parameters instead; references:
	// https://golang.org/pkg/net/http/#Request.ParseForm
	// https://stackoverflow.com/questions/299628/is-an-entity-body-allowed-for-an-http-delete-request
	// https://tools.ietf.org/html/rfc7231#section-4.3.5
	t.Run("DELETE", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = http.del("HTTPBIN_URL/delete?test=mest");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.json().args.test != "mest") ***REMOVED*** throw new Error("wrong args: " + JSON.stringify(res.json().args)); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "DELETE", sr("HTTPBIN_URL/delete?test=mest"), "", 200, "")
	***REMOVED***)

	postMethods := map[string]string***REMOVED***
		"POST":  "post",
		"PUT":   "put",
		"PATCH": "patch",
	***REMOVED***
	for method, fn := range postMethods ***REMOVED***
		t.Run(method, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(fmt.Sprintf(sr(`
				var res = http.%s("HTTPBIN_URL/%s", "data");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().data != "data") ***REMOVED*** throw new Error("wrong data: " + res.json().data); ***REMOVED***
				if (res.json().headers["Content-Type"]) ***REMOVED*** throw new Error("content type set: " + res.json().headers["Content-Type"]); ***REMOVED***
				`), fn, strings.ToLower(method)))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), method, sr("HTTPBIN_URL/")+strings.ToLower(method), "", 200, "")

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(fmt.Sprintf(sr(`
				var res = http.%s("HTTPBIN_URL/%s", ***REMOVED***a: "a", b: 2***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().form.a != "a") ***REMOVED*** throw new Error("wrong a=: " + res.json().form.a); ***REMOVED***
				if (res.json().form.b != "2") ***REMOVED*** throw new Error("wrong b=: " + res.json().form.b); ***REMOVED***
				if (res.json().headers["Content-Type"] != "application/x-www-form-urlencoded") ***REMOVED*** throw new Error("wrong content type: " + res.json().headers["Content-Type"]); ***REMOVED***
				`), fn, strings.ToLower(method)))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), method, sr("HTTPBIN_URL/")+strings.ToLower(method), "", 200, "")
				t.Run("Content-Type", func(t *testing.T) ***REMOVED***
					_, err := rt.RunString(fmt.Sprintf(sr(`
						var res = http.%s("HTTPBIN_URL/%s", ***REMOVED***a: "a", b: 2***REMOVED***, ***REMOVED***headers: ***REMOVED***"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"***REMOVED******REMOVED***);
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
						if (res.json().form.a != "a") ***REMOVED*** throw new Error("wrong a=: " + res.json().form.a); ***REMOVED***
						if (res.json().form.b != "2") ***REMOVED*** throw new Error("wrong b=: " + res.json().form.b); ***REMOVED***
						if (res.json().headers["Content-Type"] != "application/x-www-form-urlencoded; charset=utf-8") ***REMOVED*** throw new Error("wrong content type: " + res.json().headers["Content-Type"]); ***REMOVED***
						`), fn, strings.ToLower(method)))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), method, sr("HTTPBIN_URL/")+strings.ToLower(method), "", 200, "")
				***REMOVED***)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***

	t.Run("Batch", func(t *testing.T) ***REMOVED***
		t.Run("error", func(t *testing.T) ***REMOVED***
			invalidURLerr := `invalid URL: parse "https:// invalidurl.com": invalid character " " in host name`
			testCases := []struct ***REMOVED***
				name, code, expErr string
				throw              bool
			***REMOVED******REMOVED***
				***REMOVED***
					name: "invalid arg", code: `"https://somevalidurl.com"`,
					expErr: `invalid http.batch() argument type string`, throw: true,
				***REMOVED***,
				***REMOVED***
					name: "invalid URL short", code: `["https:// invalidurl.com"]`,
					expErr: invalidURLerr, throw: true,
				***REMOVED***,
				***REMOVED***
					name: "invalid URL short no throw", code: `["https:// invalidurl.com"]`,
					expErr: invalidURLerr, throw: false,
				***REMOVED***,
				***REMOVED***
					name: "invalid URL array", code: `[ ["GET", "https:// invalidurl.com"] ]`,
					expErr: invalidURLerr, throw: true,
				***REMOVED***,
				***REMOVED***
					name: "invalid URL array no throw", code: `[ ["GET", "https:// invalidurl.com"] ]`,
					expErr: invalidURLerr, throw: false,
				***REMOVED***,
				***REMOVED***
					name: "invalid URL object", code: `[ ***REMOVED***method: "GET", url: "https:// invalidurl.com"***REMOVED*** ]`,
					expErr: invalidURLerr, throw: true,
				***REMOVED***,
				***REMOVED***
					name: "invalid object no throw", code: `[ ***REMOVED***method: "GET", url: "https:// invalidurl.com"***REMOVED*** ]`,
					expErr: invalidURLerr, throw: false,
				***REMOVED***,
				***REMOVED***
					name: "object no url key", code: `[ ***REMOVED***method: "GET"***REMOVED*** ]`,
					expErr: `batch request 0 doesn't have a url key`, throw: true,
				***REMOVED***,
			***REMOVED***

			for _, tc := range testCases ***REMOVED***
				tc := tc
				t.Run(tc.name, func(t *testing.T) ***REMOVED*** //nolint:paralleltest
					oldThrow := state.Options.Throw.Bool
					state.Options.Throw.Bool = tc.throw
					defer func() ***REMOVED*** state.Options.Throw.Bool = oldThrow ***REMOVED***()

					hook := logtest.NewLocal(state.Logger)
					defer hook.Reset()

					_, err := rt.RunString(fmt.Sprintf(`http.batch(%s);`, tc.code))
					if tc.throw ***REMOVED***
						require.Error(t, err)
						assert.Contains(t, err.Error(), tc.expErr)
					***REMOVED*** else ***REMOVED***
						logEntry := hook.LastEntry()
						require.NotNil(t, logEntry)
						assert.Equal(t, logrus.WarnLevel, logEntry.Level)
						assert.Contains(t, logEntry.Data["error"].(error).Error(), tc.expErr)
						assert.Equal(t, "Request Failed", logEntry.Message)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
		t.Run("GET", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
		***REMOVED***
			let reqs = [
				["GET", "HTTPBIN_URL/"],
				["GET", "HTTPBIN_IP_URL/"],
			];
			let res = http.batch(reqs);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].url != reqs[key][1]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
			***REMOVED***
		***REMOVED***`))
			require.NoError(t, err)
			bufSamples := stats.GetBufferedSamples(samples)
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")

			t.Run("Tagged", func(t *testing.T) ***REMOVED***
				_, err := runES6String(t, rt, sr(`
			***REMOVED***
				let fragment = "get";
				let reqs = [
					["GET", http.url`+"`"+`HTTPBIN_URL/$***REMOVED***fragment***REMOVED***`+"`"+`],
					["GET", http.url`+"`"+`HTTPBIN_IP_URL/`+"`"+`],
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key][1].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key][1].url); ***REMOVED***
				***REMOVED***
			***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get"), sr("HTTPBIN_URL/$***REMOVED******REMOVED***"), 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
			***REMOVED***)

			t.Run("Shorthand", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
			***REMOVED***
				let reqs = [
					"HTTPBIN_URL/",
					"HTTPBIN_IP_URL/",
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key]) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url); ***REMOVED***
				***REMOVED***
			***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")

				t.Run("Tagged", func(t *testing.T) ***REMOVED***
					_, err := runES6String(t, rt, sr(`
				***REMOVED***
					let fragment = "get";
					let reqs = [
						http.url`+"`"+`HTTPBIN_URL/$***REMOVED***fragment***REMOVED***`+"`"+`,
						http.url`+"`"+`HTTPBIN_IP_URL/`+"`"+`,
					];
					let res = http.batch(reqs);
					for (var key in res) ***REMOVED***
						if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
						if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key].url); ***REMOVED***
					***REMOVED***
				***REMOVED***`))
					assert.NoError(t, err)
					bufSamples := stats.GetBufferedSamples(samples)
					assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get"), sr("HTTPBIN_URL/$***REMOVED******REMOVED***"), 200, "")
					assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
				***REMOVED***)
			***REMOVED***)

			t.Run("ObjectForm", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
			***REMOVED***
				let reqs = [
					***REMOVED*** method: "GET", url: "HTTPBIN_URL/" ***REMOVED***,
					***REMOVED*** url: "HTTPBIN_IP_URL/", method: "GET"***REMOVED***,
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key].url); ***REMOVED***
				***REMOVED***
			***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
			***REMOVED***)

			t.Run("ObjectKeys", func(t *testing.T) ***REMOVED***
				_, err := rt.RunString(sr(`
				var reqs = ***REMOVED***
					shorthand: "HTTPBIN_URL/get?r=shorthand",
					arr: ["GET", "HTTPBIN_URL/get?r=arr", null, ***REMOVED***tags: ***REMOVED***name: 'arr'***REMOVED******REMOVED***],
					obj1: ***REMOVED*** method: "GET", url: "HTTPBIN_URL/get?r=obj1" ***REMOVED***,
					obj2: ***REMOVED*** url: "HTTPBIN_URL/get?r=obj2", params: ***REMOVED***tags: ***REMOVED***name: 'obj2'***REMOVED******REMOVED***, method: "GET"***REMOVED***,
				***REMOVED***;
				var res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].json().args.r != key) ***REMOVED*** throw new Error("wrong request id: " + key); ***REMOVED***
				***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get?r=shorthand"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get?r=arr"), "arr", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get?r=obj1"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get?r=obj2"), "obj2", 200, "")
			***REMOVED***)

			t.Run("BodyAndParams", func(t *testing.T) ***REMOVED***
				testStr := "testbody"
				rt.Set("someStrFile", testStr)
				rt.Set("someBinFile", []byte(testStr))

				_, err := rt.RunString(sr(`
					var reqs = [
						["POST", "HTTPBIN_URL/post", "testbody"],
						["POST", "HTTPBIN_URL/post", someStrFile],
						["POST", "HTTPBIN_URL/post", someBinFile],
						***REMOVED***
							method: "POST",
							url: "HTTPBIN_URL/post",
							test: "test1",
							body: "testbody",
						***REMOVED***, ***REMOVED***
							body: someBinFile,
							url: "HTTPBIN_IP_URL/post",
							params: ***REMOVED*** tags: ***REMOVED*** name: "myname" ***REMOVED*** ***REMOVED***,
							method: "POST",
						***REMOVED***, ***REMOVED***
							method: "POST",
							url: "HTTPBIN_IP_URL/post",
							body: ***REMOVED***
								hello: "world!",
							***REMOVED***,
							params: ***REMOVED***
								tags: ***REMOVED*** name: "myname" ***REMOVED***,
								headers: ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					];
					var res = http.batch(reqs);
					for (var key in res) ***REMOVED***
						if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
						if (res[key].json().data != "testbody" && res[key].json().form.hello != "world!") ***REMOVED*** throw new Error("wrong response for " + key + ": " + res[key].body); ***REMOVED***
					***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "POST", sr("HTTPBIN_URL/post"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "POST", sr("HTTPBIN_IP_URL/post"), "myname", 200, "")
			***REMOVED***)
		***REMOVED***)
		t.Run("POST", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.batch([ ["POST", "HTTPBIN_URL/post", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + key + ": " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)
		t.Run("PUT", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var res = http.batch([ ["PUT", "HTTPBIN_URL/put", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + key + ": " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "PUT", sr("HTTPBIN_URL/put"), "", 200, "")
		***REMOVED***)
	***REMOVED***)

	t.Run("HTTPRequest", func(t *testing.T) ***REMOVED***
		t.Run("EmptyBody", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var reqUrl = "HTTPBIN_URL/cookies"
				var res = http.get(reqUrl);
				var jar = new http.CookieJar();

				jar.set("HTTPBIN_URL/cookies", "key", "value");
				res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key2: "value2" ***REMOVED***, jar: jar ***REMOVED***);

				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***

				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.request["method"] !== "GET") ***REMOVED*** throw new Error("http request method was not \"GET\": " + JSON.stringify(res.request)) ***REMOVED***
				if (res.request["body"].length != 0) ***REMOVED*** throw new Error("http request body was not null: " + JSON.stringify(res.request["body"])) ***REMOVED***
				if (res.request["url"] != reqUrl) ***REMOVED***
					throw new Error("wrong http request url: " + JSON.stringify(res.request))
				***REMOVED***
				if (res.request["cookies"]["key2"][0].name != "key2") ***REMOVED*** throw new Error("wrong http request cookies: " + JSON.stringify(JSON.stringify(res.request["cookies"]["key2"]))) ***REMOVED***
				if (res.request["headers"]["User-Agent"][0] != "TestUserAgent") ***REMOVED*** throw new Error("wrong http request headers: " + JSON.stringify(res.request)) ***REMOVED***
				`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("NonEmptyBody", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
				var res = http.post("HTTPBIN_URL/post", ***REMOVED***a: "a", b: 2***REMOVED***, ***REMOVED***headers: ***REMOVED***"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"***REMOVED******REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.request["body"] != "a=a&b=2") ***REMOVED*** throw new Error("http request body was not set properly: " + JSON.stringify(res.request))***REMOVED***
				`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestRequestArrayBufferBody(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, _, _, rt, _ := newRuntime(t) //nolint: dogsled
	sr := tb.Replacer.Replace

	tb.Mux.HandleFunc("/post-arraybuffer", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		require.Equal(t, "POST", r.Method)
		var in bytes.Buffer
		_, err := io.Copy(&in, r.Body)
		require.NoError(t, err)
		_, err = w.Write(in.Bytes())
		require.NoError(t, err)
	***REMOVED***))

	testCases := []struct ***REMOVED***
		arr, expected string
	***REMOVED******REMOVED***
		***REMOVED***"Uint8Array", "104,101,108,108,111"***REMOVED***,
		***REMOVED***"Uint16Array", "104,0,101,0,108,0,108,0,111,0"***REMOVED***,
		***REMOVED***"Uint32Array", "104,0,0,0,101,0,0,0,108,0,0,0,108,0,0,0,111,0,0,0"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.arr, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(fmt.Sprintf(`
			var arr = new %[1]s([104, 101, 108, 108, 111]); // "hello"
			var res = http.post("HTTPBIN_URL/post-arraybuffer", arr.buffer, ***REMOVED*** responseType: 'binary' ***REMOVED***);

			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***

			var resTyped = new Uint8Array(res.body);
			var exp = new %[1]s([%[2]s]);
			if (exp.length !== resTyped.length) ***REMOVED***
				throw new Error(
					"incorrect data length: expected " + exp.length + ", received " + resTypedLength)
			***REMOVED***
			for (var i = 0; i < exp.length; i++) ***REMOVED***
				if (exp[i] !== resTyped[i])	***REMOVED***
					throw new Error(
						"incorrect data at index " + i + ": expected " + exp[i] + ", received " + resTyped[i])
				***REMOVED***
			***REMOVED***
			`, tc.arr, tc.expected)))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestRequestCompression(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	state.Logger.AddHook(&logHook)

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)

	text := `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit.
	Maecenas sed pharetra sapien. Nunc laoreet molestie ante ac gravida.
	Etiam interdum dui viverra posuere egestas. Pellentesque at dolor tristique,
	mattis turpis eget, commodo purus. Nunc orci aliquam.`

	decompress := func(algo string, input io.Reader) io.Reader ***REMOVED***
		switch algo ***REMOVED***
		case "br":
			w := brotli.NewReader(input)
			return w
		case "gzip":
			w, err := gzip.NewReader(input)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			return w
		case "deflate":
			w, err := zlib.NewReader(input)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			return w
		case "zstd":
			w, err := zstd.NewReader(input)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			return w
		default:
			t.Fatal("unknown algorithm " + algo)
		***REMOVED***
		return nil // unreachable
	***REMOVED***

	var (
		expectedEncoding string
		actualEncoding   string
	)
	tb.Mux.HandleFunc("/compressed-text", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		require.Equal(t, expectedEncoding, r.Header.Get("Content-Encoding"))

		expectedLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
		require.NoError(t, err)
		algos := strings.Split(actualEncoding, ", ")
		compressedBuf := new(bytes.Buffer)
		n, err := io.Copy(compressedBuf, r.Body)
		require.Equal(t, int(n), expectedLength)
		require.NoError(t, err)
		var prev io.Reader = compressedBuf

		if expectedEncoding != "" ***REMOVED***
			for i := len(algos) - 1; i >= 0; i-- ***REMOVED***
				prev = decompress(algos[i], prev)
			***REMOVED***
		***REMOVED***

		var buf bytes.Buffer
		_, err = io.Copy(&buf, prev)
		require.NoError(t, err)
		require.Equal(t, text, buf.String())
	***REMOVED***))

	testCases := []struct ***REMOVED***
		name          string
		compression   string
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***compression: ""***REMOVED***,
		***REMOVED***compression: "  "***REMOVED***,
		***REMOVED***compression: "gzip"***REMOVED***,
		***REMOVED***compression: "gzip, gzip"***REMOVED***,
		***REMOVED***compression: "gzip,   gzip "***REMOVED***,
		***REMOVED***compression: "gzip,gzip"***REMOVED***,
		***REMOVED***compression: "gzip, gzip, gzip, gzip, gzip, gzip, gzip"***REMOVED***,
		***REMOVED***compression: "deflate"***REMOVED***,
		***REMOVED***compression: "deflate, gzip"***REMOVED***,
		***REMOVED***compression: "gzip,deflate, gzip"***REMOVED***,
		***REMOVED***compression: "zstd"***REMOVED***,
		***REMOVED***compression: "zstd, gzip, deflate"***REMOVED***,
		***REMOVED***compression: "br"***REMOVED***,
		***REMOVED***compression: "br, gzip, deflate"***REMOVED***,
		***REMOVED***
			compression:   "George",
			expectedError: `unknown compression algorithm George`,
		***REMOVED***,
		***REMOVED***
			compression:   "gzip, George",
			expectedError: `unknown compression algorithm George`,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.compression, func(t *testing.T) ***REMOVED***
			algos := strings.Split(testCase.compression, ",")
			for i, algo := range algos ***REMOVED***
				algos[i] = strings.TrimSpace(algo)
			***REMOVED***
			expectedEncoding = strings.Join(algos, ", ")
			actualEncoding = expectedEncoding
			_, err := runES6String(t, rt, tb.Replacer.Replace(`
		http.post("HTTPBIN_URL/compressed-text", `+"`"+text+"`"+`,  ***REMOVED***"compression": "`+testCase.compression+`"***REMOVED***);
	`))
			if testCase.expectedError == "" ***REMOVED***
				require.NoError(t, err)
			***REMOVED*** else ***REMOVED***
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.expectedError)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	t.Run("custom set header", func(t *testing.T) ***REMOVED***
		expectedEncoding = "not, valid"
		actualEncoding = "gzip, deflate"

		logHook.Drain()
		t.Run("encoding", func(t *testing.T) ***REMOVED***
			_, err := runES6String(t, rt, tb.Replacer.Replace(`
				http.post("HTTPBIN_URL/compressed-text", `+"`"+text+"`"+`,
					***REMOVED***"compression": "`+actualEncoding+`",
					 "headers": ***REMOVED***"Content-Encoding": "`+expectedEncoding+`"***REMOVED***
					***REMOVED***
				);
			`))
			require.NoError(t, err)
			require.NotEmpty(t, logHook.Drain())
		***REMOVED***)

		t.Run("encoding and length", func(t *testing.T) ***REMOVED***
			_, err := runES6String(t, rt, tb.Replacer.Replace(`
				http.post("HTTPBIN_URL/compressed-text", `+"`"+text+"`"+`,
					***REMOVED***"compression": "`+actualEncoding+`",
					 "headers": ***REMOVED***"Content-Encoding": "`+expectedEncoding+`",
								 "Content-Length": "12"***REMOVED***
					***REMOVED***
				);
			`))
			require.NoError(t, err)
			require.NotEmpty(t, logHook.Drain())
		***REMOVED***)

		expectedEncoding = actualEncoding
		t.Run("correct encoding", func(t *testing.T) ***REMOVED***
			_, err := runES6String(t, rt, tb.Replacer.Replace(`
				http.post("HTTPBIN_URL/compressed-text", `+"`"+text+"`"+`,
					***REMOVED***"compression": "`+actualEncoding+`",
					 "headers": ***REMOVED***"Content-Encoding": "`+actualEncoding+`"***REMOVED***
					***REMOVED***
				);
			`))
			require.NoError(t, err)
			require.Empty(t, logHook.Drain())
		***REMOVED***)

		// TODO: move to some other test?
		t.Run("correct length", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(tb.Replacer.Replace(
				`http.post("HTTPBIN_URL/post", "0123456789", ***REMOVED*** "headers": ***REMOVED***"Content-Length": "10"***REMOVED******REMOVED***);`,
			))
			require.NoError(t, err)
			require.Empty(t, logHook.Drain())
		***REMOVED***)

		t.Run("content-length is set", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(tb.Replacer.Replace(`
				var resp = http.post("HTTPBIN_URL/post", "0123456789");
				if (resp.json().headers["Content-Length"][0] != "10") ***REMOVED***
					throw new Error("content-length not set: " + JSON.stringify(resp.json().headers));
				***REMOVED***
			`))
			require.NoError(t, err)
			require.Empty(t, logHook.Drain())
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestResponseTypes(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)

	text := `•?((¯°·._.• ţ€$ţɨɲǥ µɲɨȼ๏ď€ ɨɲ Ќ6 •._.·°¯))؟•`
	textLen := len(text)
	tb.Mux.HandleFunc("/get-text", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		n, err := w.Write([]byte(text))
		assert.NoError(t, err)
		assert.Equal(t, textLen, n)
	***REMOVED***))
	tb.Mux.HandleFunc("/compare-text", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, text, string(body))
	***REMOVED***))

	binaryLen := 300
	binary := make([]byte, binaryLen)
	for i := 0; i < binaryLen; i++ ***REMOVED***
		binary[i] = byte(i)
	***REMOVED***
	tb.Mux.HandleFunc("/get-bin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		n, err := w.Write(binary)
		assert.NoError(t, err)
		assert.Equal(t, binaryLen, n)
	***REMOVED***))
	tb.Mux.HandleFunc("/compare-bin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.True(t, bytes.Equal(binary, body))
	***REMOVED***))

	replace := func(s string) string ***REMOVED***
		return strings.NewReplacer(
			"EXP_TEXT", text,
			"EXP_BIN_LEN", strconv.Itoa(binaryLen),
		).Replace(tb.Replacer.Replace(s))
	***REMOVED***

	_, err := rt.RunString(replace(`
		var expText = "EXP_TEXT";
		var expBinLength = EXP_BIN_LEN;

		// Check default behaviour with a unicode text
		var respTextImplicit = http.get("HTTPBIN_URL/get-text").body;
		if (respTextImplicit !== expText) ***REMOVED***
			throw new Error("default response body should be '" + expText + "' but was '" + respTextImplicit + "'");
		***REMOVED***
		http.post("HTTPBIN_URL/compare-text", respTextImplicit);

		// Check discarding of responses
		var respNone = http.get("HTTPBIN_URL/get-text", ***REMOVED*** responseType: "none" ***REMOVED***).body;
		if (respNone != null) ***REMOVED***
			throw new Error("none response body should be null but was " + respNone);
		***REMOVED***

		// Check binary transmission of the text response as well
		var respBin = http.get("HTTPBIN_URL/get-text", ***REMOVED*** responseType: "binary" ***REMOVED***);

		// Convert a UTF-8 ArrayBuffer to a JS string
		var respBinText = String.fromCharCode.apply(null, new Uint8Array(respBin.body));
		var strConv = decodeURIComponent(escape(respBinText));
		if (strConv !== expText) ***REMOVED***
			throw new Error("converted response body should be '" + expText + "' but was '" + strConv + "'");
		***REMOVED***
		http.post("HTTPBIN_URL/compare-text", respBin.body);

		// Check binary response
		var respBin = http.get("HTTPBIN_URL/get-bin", ***REMOVED*** responseType: "binary" ***REMOVED***);
		var respBinTyped = new Uint8Array(respBin.body);
		if (expBinLength !== respBinTyped.length) ***REMOVED***
			throw new Error("response body length should be '" + expBinLength
							+ "' but was '" + respBinTyped.length + "'");
		***REMOVED***
		for(var i = 0; i < respBinTyped.length; i++) ***REMOVED***
			if (respBinTyped[i] !== i%256) ***REMOVED***
				throw new Error("expected value " + (i%256) + " to be at position "
								+ i + " but it was " + respBinTyped[i]);
			***REMOVED***
		***REMOVED***
		http.post("HTTPBIN_URL/compare-bin", respBin.body);

		// Check ArrayBuffer response
		var respBin = http.get("HTTPBIN_URL/get-bin", ***REMOVED*** responseType: "binary" ***REMOVED***).body;
		if (respBin.byteLength !== expBinLength) ***REMOVED***
			throw new Error("response body length should be '" + expBinLength + "' but was '" + respBin.byteLength + "'");
		***REMOVED***

		// Check ArrayBuffer responses with http.batch()
		var responses = http.batch([
			["GET", "HTTPBIN_URL/get-bin", null, ***REMOVED*** responseType: "binary" ***REMOVED***],
			["GET", "HTTPBIN_URL/get-bin", null, ***REMOVED*** responseType: "binary" ***REMOVED***],
		]);
		if (responses.length != 2) ***REMOVED***
			throw new Error("expected 2 responses, received " + responses.length);
		***REMOVED***
		for (var i = 0; i < responses.length; i++) ***REMOVED***
			if (responses[i].body.byteLength !== expBinLength) ***REMOVED***
				throw new Error("response body length should be '"
					+ expBinLength + "' but was '" + responses[i].body.byteLength + "'");
			***REMOVED***
		***REMOVED***
	`))
	assert.NoError(t, err)

	// Verify that if we enable discardResponseBodies globally, the default value is none
	state.Options.DiscardResponseBodies = null.BoolFrom(true)

	_, err = rt.RunString(replace(`
		var expText = "EXP_TEXT";

		// Check default behaviour
		var respDefault = http.get("HTTPBIN_URL/get-text").body;
		if (respDefault !== null) ***REMOVED***
			throw new Error("default response body should be discarded and null but was " + respDefault);
		***REMOVED***

		// Check explicit text response
		var respTextExplicit = http.get("HTTPBIN_URL/get-text", ***REMOVED*** responseType: "text" ***REMOVED***).body;
		if (respTextExplicit !== expText) ***REMOVED***
			throw new Error("text response body should be '" + expText + "' but was '" + respTextExplicit + "'");
		***REMOVED***
		http.post("HTTPBIN_URL/compare-text", respTextExplicit);
	`))
	assert.NoError(t, err)
***REMOVED***

func checkErrorCode(t testing.TB, tags *stats.SampleTags, code int, msg string) ***REMOVED***
	errorMsg, ok := tags.Get("error")
	if msg == "" ***REMOVED***
		assert.False(t, ok)
	***REMOVED*** else ***REMOVED***
		assert.Contains(t, errorMsg, msg)
	***REMOVED***
	errorCodeStr, ok := tags.Get("error_code")
	if code == 0 ***REMOVED***
		assert.False(t, ok)
	***REMOVED*** else ***REMOVED***
		errorCode, err := strconv.Atoi(errorCodeStr)
		assert.NoError(t, err)
		assert.Equal(t, code, errorCode)
	***REMOVED***
***REMOVED***

func TestErrorCodes(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, samples, rt, _ := newRuntime(t)
	state.Options.Throw = null.BoolFrom(false)
	sr := tb.Replacer.Replace

	// Handple paths with custom logic
	tb.Mux.HandleFunc("/no-location-redirect", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) ***REMOVED***
		w.WriteHeader(302)
	***REMOVED***))
	tb.Mux.HandleFunc("/bad-location-redirect", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) ***REMOVED***
		w.Header().Set("Location", "h\t:/") // \n is forbidden
		w.WriteHeader(302)
	***REMOVED***))

	testCases := []struct ***REMOVED***
		name                string
		status              int
		moreSamples         int
		expectedErrorCode   int
		expectedErrorMsg    string
		expectedScriptError string
		script              string
	***REMOVED******REMOVED***
		***REMOVED***
			name:              "Unroutable",
			expectedErrorCode: 1101,
			expectedErrorMsg:  "lookup: no such host",
			script:            `var res = http.request("GET", "http://sdafsgdhfjg.com/");`,
		***REMOVED***,

		***REMOVED***
			name:              "404",
			status:            404,
			expectedErrorCode: 1404,
			script:            `var res = http.request("GET", "HTTPBIN_URL/status/404");`,
		***REMOVED***,
		***REMOVED***
			name:              "Unroutable redirect",
			expectedErrorCode: 1101,
			expectedErrorMsg:  "lookup: no such host",
			moreSamples:       1,
			script:            `var res = http.request("GET", "HTTPBIN_URL/redirect-to?url=http://dafsgdhfjg.com/");`,
		***REMOVED***,
		***REMOVED***
			name:              "Non location redirect",
			expectedErrorCode: 1000,
			expectedErrorMsg:  "302 response missing Location header",
			script:            `var res = http.request("GET", "HTTPBIN_URL/no-location-redirect");`,
		***REMOVED***,
		***REMOVED***
			name:              "Bad location redirect",
			expectedErrorCode: 1000,
			expectedErrorMsg:  "failed to parse Location header \"h\\t:/\": ",
			script:            `var res = http.request("GET", "HTTPBIN_URL/bad-location-redirect");`,
		***REMOVED***,
		***REMOVED***
			name:              "Missing protocol",
			expectedErrorCode: 1000,
			expectedErrorMsg:  `unsupported protocol scheme ""`,
			script:            `var res = http.request("GET", "dafsgdhfjg.com/");`,
		***REMOVED***,
		***REMOVED***
			name:        "Too many redirects",
			status:      302,
			moreSamples: 2,
			script: `
			var res = http.get("HTTPBIN_URL/relative-redirect/3", ***REMOVED***redirects: 2***REMOVED***);
			if (res.url != "HTTPBIN_URL/relative-redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***`,
		***REMOVED***,
		***REMOVED***
			name:              "Connection refused redirect",
			status:            0,
			moreSamples:       1,
			expectedErrorMsg:  `dial: connection refused`,
			expectedErrorCode: 1212,
			script: `
			var res = http.get("HTTPBIN_URL/redirect-to?url=http%3A%2F%2F127.0.0.1%3A1%2Fpesho");
			if (res.url != "http://127.0.0.1:1/pesho") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***`,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		// clear the Samples
		stats.GetBufferedSamples(samples)
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(testCase.script + "\n" + fmt.Sprintf(`
			if (res.status != %d) ***REMOVED*** throw new Error("wrong status: "+ res.status);***REMOVED***
			if (res.error.indexOf(%q, 0) === -1) ***REMOVED*** throw new Error("wrong error: '" + res.error + "'");***REMOVED***
			if (res.error_code != %d) ***REMOVED*** throw new Error("wrong error_code: "+ res.error_code);***REMOVED***
			`, testCase.status, testCase.expectedErrorMsg, testCase.expectedErrorCode)))
			if testCase.expectedScriptError == "" ***REMOVED***
				require.NoError(t, err)
			***REMOVED*** else ***REMOVED***
				require.Error(t, err)
				require.Equal(t, err.Error(), testCase.expectedScriptError)
			***REMOVED***
			cs := stats.GetBufferedSamples(samples)
			assert.Len(t, cs, 1+testCase.moreSamples)
			for _, c := range cs[len(cs)-1:] ***REMOVED***
				assert.NotZero(t, len(c.GetSamples()))
				for _, sample := range c.GetSamples() ***REMOVED***
					checkErrorCode(t, sample.GetTags(), testCase.expectedErrorCode, testCase.expectedErrorMsg)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestResponseWaitingAndReceivingTimings(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)

	tb.Mux.HandleFunc("/slow-response", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		flusher, ok := w.(http.Flusher)
		require.True(t, ok)

		time.Sleep(1200 * time.Millisecond)
		n, err := w.Write([]byte("1st bytes!"))
		assert.NoError(t, err)
		assert.Equal(t, 10, n)

		flusher.Flush()
		time.Sleep(1200 * time.Millisecond)

		n, err = w.Write([]byte("2nd bytes!"))
		assert.NoError(t, err)
		assert.Equal(t, 10, n)
	***REMOVED***))

	_, err := rt.RunString(tb.Replacer.Replace(`
		var resp = http.get("HTTPBIN_URL/slow-response");

		if (resp.timings.waiting < 1000) ***REMOVED***
			throw new Error("expected waiting time to be over 1000ms but was " + resp.timings.waiting);
		***REMOVED***

		if (resp.timings.receiving < 1000) ***REMOVED***
			throw new Error("expected receiving time to be over 1000ms but was " + resp.timings.receiving);
		***REMOVED***

		if (resp.body !== "1st bytes!2nd bytes!") ***REMOVED***
			throw new Error("wrong response body: " + resp.body);
		***REMOVED***
	`))
	assert.NoError(t, err)
***REMOVED***

func TestResponseTimingsWhenTimeout(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	// We expect a failed request
	state.Options.Throw = null.BoolFrom(false)

	_, err := rt.RunString(tb.Replacer.Replace(`
		var resp = http.get("HTTPBIN_URL/delay/10", ***REMOVED*** timeout: 2500 ***REMOVED***);

		if (resp.timings.waiting < 2000) ***REMOVED***
			throw new Error("expected waiting time to be over 2000ms but was " + resp.timings.waiting);
		***REMOVED***

		if (resp.timings.duration < 2000) ***REMOVED***
			throw new Error("expected duration time to be over 2000ms but was " + resp.timings.duration);
		***REMOVED***
	`))
	assert.NoError(t, err)
***REMOVED***

func TestNoResponseBodyMangling(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)

	_, err := rt.RunString(tb.Replacer.Replace(`
	    var batchSize = 100;

		var requests = [];

		for (var i = 0; i < batchSize; i++) ***REMOVED***
			requests.push(["GET", "HTTPBIN_URL/get?req=" + i, null, ***REMOVED*** responseType: (i % 2 ? "binary" : "text") ***REMOVED***]);
		***REMOVED***

		var responses = http.batch(requests);

		for (var i = 0; i < batchSize; i++) ***REMOVED***
			var reqNumber = parseInt(responses[i].json().args.req[0], 10);
			if (i !== reqNumber) ***REMOVED***
				throw new Error("Response " + i + " has " + reqNumber + ", expected " + i)
			***REMOVED***
		***REMOVED***
	`))
	assert.NoError(t, err)
***REMOVED***

func TestRedirectMetricTags(t *testing.T) ***REMOVED***
	tb, _, samples, rt, _ := newRuntime(t)

	tb.Mux.HandleFunc("/redirect/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		http.Redirect(w, r, "/get", http.StatusMovedPermanently)
	***REMOVED***))

	sr := tb.Replacer.Replace
	script := sr(`
		http.post("HTTPBIN_URL/redirect/post", ***REMOVED***data: "some data"***REMOVED***);
	`)

	_, err := rt.RunString(script)
	require.NoError(t, err)

	require.Len(t, samples, 2)

	checkTags := func(sc stats.SampleContainer, expTags map[string]string) ***REMOVED***
		allSamples := sc.GetSamples()
		assert.Len(t, allSamples, 9)
		for _, s := range allSamples ***REMOVED***
			assert.Equal(t, expTags, s.Tags.CloneTags())
		***REMOVED***
	***REMOVED***
	expPOSTtags := map[string]string***REMOVED***
		"group":             "",
		"method":            "POST",
		"url":               sr("HTTPBIN_URL/redirect/post"),
		"name":              sr("HTTPBIN_URL/redirect/post"),
		"status":            "301",
		"proto":             "HTTP/1.1",
		"expected_response": "true",
	***REMOVED***
	expGETtags := map[string]string***REMOVED***
		"group":             "",
		"method":            "GET",
		"url":               sr("HTTPBIN_URL/get"),
		"name":              sr("HTTPBIN_URL/get"),
		"status":            "200",
		"proto":             "HTTP/1.1",
		"expected_response": "true",
	***REMOVED***
	checkTags(<-samples, expPOSTtags)
	checkTags(<-samples, expGETtags)
***REMOVED***

func BenchmarkHandlingOfResponseBodies(b *testing.B) ***REMOVED***
	tb, state, samples, rt, _ := newRuntime(b)

	state.BPool = bpool.NewBufferPool(100)

	go func() ***REMOVED***
		ctxDone := tb.Context.Done()
		for ***REMOVED***
			select ***REMOVED***
			case <-samples:
			case <-ctxDone:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	mbData := bytes.Repeat([]byte("0123456789"), 100000)
	tb.Mux.HandleFunc("/1mbdata", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		_, err := resp.Write(mbData)
		if err != nil ***REMOVED***
			b.Error(err)
		***REMOVED***
	***REMOVED***))

	testCodeTemplate := tb.Replacer.Replace(`
		http.get("HTTPBIN_URL/", ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***);
		http.post("HTTPBIN_URL/post", ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***);
		http.batch([
			["GET", "HTTPBIN_URL/gzip", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/gzip", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/deflate", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/deflate", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/redirect/5", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***], // 6 requests
			["GET", "HTTPBIN_URL/get", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/html", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/bytes/100000", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/image/png", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/image/jpeg", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/image/jpeg", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/image/webp", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/image/svg", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/forms/post", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/bytes/100000", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
			["GET", "HTTPBIN_URL/stream-bytes/100000", null, ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***],
		]);
		http.get("HTTPBIN_URL/get", ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***);
		http.get("HTTPBIN_URL/get", ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***);
		http.get("HTTPBIN_URL/1mbdata", ***REMOVED*** responseType: "TEST_RESPONSE_TYPE" ***REMOVED***);
	`)

	testResponseType := func(responseType string) func(b *testing.B) ***REMOVED***
		testCode := strings.Replace(testCodeTemplate, "TEST_RESPONSE_TYPE", responseType, -1)
		return func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				_, err := rt.RunString(testCode)
				if err != nil ***REMOVED***
					b.Error(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	b.ResetTimer()
	b.Run("text", testResponseType("text"))
	b.Run("binary", testResponseType("binary"))
	b.Run("none", testResponseType("none"))
***REMOVED***

func TestErrorsWithDecompression(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, _, rt, _ := newRuntime(t)

	state.Options.Throw = null.BoolFrom(false)

	tb.Mux.HandleFunc("/broken-archive", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		enc := r.URL.Query()["encoding"][0]
		w.Header().Set("Content-Encoding", enc)
		_, _ = fmt.Fprintf(w, "Definitely not %s, but it's all cool...", enc)
	***REMOVED***))

	_, err := rt.RunString(tb.Replacer.Replace(`
		function handleResponseEncodingError (encoding) ***REMOVED***
			var resp = http.get("HTTPBIN_URL/broken-archive?encoding=" + encoding);
			if (resp.error_code != 1701) ***REMOVED***
				throw new Error("Expected error_code 1701 for '" + encoding +"', but got " + resp.error_code);
			***REMOVED***
		***REMOVED***

		["gzip", "deflate", "br", "zstd"].forEach(handleResponseEncodingError);
	`))
	assert.NoError(t, err)
***REMOVED***

func TestRequestAndBatchTLS(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip()
	***REMOVED***
	t.Parallel()
	_, _, samples, rt, _ := newRuntime(t) //nolint:dogsled

	t.Run("cert_expired", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`http.get("https://expired.badssl.com/");`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "x509: certificate has expired or is not yet valid")
	***REMOVED***)
	tlsVersionTests := []struct ***REMOVED***
		Name, URL, Version string
	***REMOVED******REMOVED***
		***REMOVED***Name: "tls10", URL: "https://tls-v1-0.badssl.com:1010/", Version: "http.TLS_1_0"***REMOVED***,
		***REMOVED***Name: "tls11", URL: "https://tls-v1-1.badssl.com:1011/", Version: "http.TLS_1_1"***REMOVED***,
		***REMOVED***Name: "tls12", URL: "https://badssl.com/", Version: "http.TLS_1_2"***REMOVED***,
	***REMOVED***
	for _, versionTest := range tlsVersionTests ***REMOVED***
		versionTest := versionTest
		t.Run(versionTest.Name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(fmt.Sprintf(`
					var res = http.get("%s");
					if (res.tls_version != %s) ***REMOVED*** throw new Error("wrong TLS version: " + res.tls_version); ***REMOVED***
				`, versionTest.URL, versionTest.Version))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", versionTest.URL, "", 200, "")
		***REMOVED***)
	***REMOVED***
	tlsCipherSuiteTests := []struct ***REMOVED***
		Name, URL, CipherSuite string
	***REMOVED******REMOVED***
		***REMOVED***Name: "cipher_suite_cbc", URL: "https://cbc.badssl.com/", CipherSuite: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"***REMOVED***,
		***REMOVED***Name: "cipher_suite_ecc384", URL: "https://ecc384.badssl.com/", CipherSuite: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"***REMOVED***,
	***REMOVED***
	for _, cipherSuiteTest := range tlsCipherSuiteTests ***REMOVED***
		cipherSuiteTest := cipherSuiteTest
		t.Run(cipherSuiteTest.Name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(fmt.Sprintf(`
					var res = http.get("%s");
					if (res.tls_cipher_suite != "%s") ***REMOVED*** throw new Error("wrong TLS cipher suite: " + res.tls_cipher_suite); ***REMOVED***
				`, cipherSuiteTest.URL, cipherSuiteTest.CipherSuite))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", cipherSuiteTest.URL, "", 200, "")
		***REMOVED***)
	***REMOVED***
	t.Run("ocsp_stapled_good", func(t *testing.T) ***REMOVED***
		website := "https://www.wikipedia.org/"
		_, err := rt.RunString(fmt.Sprintf(`
			var res = http.request("GET", "%s");
			if (res.ocsp.status != http.OCSP_STATUS_GOOD) ***REMOVED*** throw new Error("wrong ocsp stapled response status: " + res.ocsp.status); ***REMOVED***
			`, website))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", website, "", 200, "")
	***REMOVED***)
***REMOVED***

func TestDigestAuthWithBody(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, samples, rt, _ := newRuntime(t)

	state.Options.Throw = null.BoolFrom(true)
	state.Options.HTTPDebug = null.StringFrom("full")

	tb.Mux.HandleFunc("/digest-auth-with-post/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		require.Equal(t, "POST", r.Method)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "super secret body", string(body))
		httpbin.New().DigestAuth(w, r) // this doesn't read the body
	***REMOVED***))

	urlWithCreds := tb.Replacer.Replace(
		"http://testuser:testpwd@HTTPBIN_IP:HTTPBIN_PORT/digest-auth-with-post/auth/testuser/testpwd",
	)

	_, err := rt.RunString(fmt.Sprintf(`
		var res = http.post(%q, "super secret body", ***REMOVED*** auth: "digest" ***REMOVED***);
		if (res.status !== 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.error_code !== 0) ***REMOVED*** throw new Error("wrong error code: " + res.error_code); ***REMOVED***
	`, urlWithCreds))
	require.NoError(t, err)

	urlRaw := tb.Replacer.Replace(
		"http://HTTPBIN_IP:HTTPBIN_PORT/digest-auth-with-post/auth/testuser/testpwd")

	sampleContainers := stats.GetBufferedSamples(samples)
	assertRequestMetricsEmitted(t, sampleContainers[0:1], "POST", urlRaw, urlRaw, 401, "")
	assertRequestMetricsEmitted(t, sampleContainers[1:2], "POST", urlRaw, urlRaw, 200, "")
***REMOVED***

func TestBinaryResponseWithStatus0(t *testing.T) ***REMOVED***
	t.Parallel()
	_, state, _, rt, _ := newRuntime(t) //nolint:dogsled
	state.Options.Throw = null.BoolFrom(false)
	_, err := rt.RunString(`
		var res = http.get("https://asdajkdahdqiuwhejkasdnakjdnadasdlkas.com", ***REMOVED*** responseType: "binary" ***REMOVED***);
		if (res.status !== 0) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body !== null) ***REMOVED*** throw new Error("wrong body: " + JSON.stringify(res.body)); ***REMOVED***
	`)
	require.NoError(t, err)
***REMOVED***
