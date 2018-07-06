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
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ThomsonReutersEikon/go-ntlm/ntlm"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/stats"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

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

func newRuntime(t *testing.T) (*testutils.HTTPMultiBin, *common.State, chan stats.SampleContainer, *goja.Runtime, *context.Context) ***REMOVED***
	tb := testutils.NewHTTPMultiBin(t)

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	samples := make(chan stats.SampleContainer, 1000)

	state := &common.State***REMOVED***
		Options: lib.Options***REMOVED***
			MaxRedirects: null.IntFrom(10),
			UserAgent:    null.StringFrom("TestUserAgent"),
			Throw:        null.BoolFrom(true),
			SystemTags:   lib.GetTagSet(lib.DefaultSystemTagList...),
			//HttpDebug:    null.StringFrom("full"),
		***REMOVED***,
		Logger:        logger,
		Group:         root,
		TLSConfig:     tb.TLSClientConfig,
		HTTPTransport: netext.NewHTTPTransport(tb.HTTPTransport),
		BPool:         bpool.NewBufferPool(1),
		Samples:       samples,
	***REMOVED***

	ctx := new(context.Context)
	*ctx = context.Background()
	*ctx = common.WithState(*ctx, state)
	*ctx = common.WithRuntime(*ctx, rt)
	rt.Set("http", common.Bind(rt, New(), ctx))

	return tb, state, samples, rt, ctx
***REMOVED***

func TestRequestAndBatch(t *testing.T) ***REMOVED***
	tb, state, samples, rt, ctx := newRuntime(t)
	defer tb.Cleanup()
	sr := tb.Replacer.Replace

	// Handple paths with custom logic
	tb.Mux.HandleFunc("/ntlm", http.HandlerFunc(ntlmHandler("bob", "pass")))
	tb.Mux.HandleFunc("/digest-auth/failure", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(2 * time.Second)
	***REMOVED***))
	tb.Mux.HandleFunc("/set-cookie-before-redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		cookie := http.Cookie***REMOVED***
			Name:   "key-foo",
			Value:  "value-bar",
			Path:   "/",
			Domain: sr("HTTPBIN_DOMAIN"),
		***REMOVED***

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, sr("HTTPBIN_URL/get"), 301)
	***REMOVED***))

	t.Run("Redirects", func(t *testing.T) ***REMOVED***
		t.Run("10", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`http.get("HTTPBIN_URL/redirect/10")`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("11", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
			let res = http.get("HTTPBIN_URL/redirect/11");
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

				_, err := common.RunString(rt, sr(`
				let res = http.get("HTTPBIN_URL/redirect/11");
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
			_, err := common.RunString(rt, sr(`
			let res = http.get("HTTPBIN_URL/redirect/1", ***REMOVED***redirects: 3***REMOVED***);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/get") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("requestScopeNoRedirects", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
			let res = http.get("HTTPBIN_URL/redirect/1", ***REMOVED***redirects: 0***REMOVED***);
			if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "HTTPBIN_URL/redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("Timeout", func(t *testing.T) ***REMOVED***
		t.Run("10s", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
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
			_, err := common.RunString(rt, sr(`
				http.get("HTTPBIN_URL/delay/10", ***REMOVED***
					timeout: 1*1000,
				***REMOVED***)
			`))
			endTime := time.Now()
			assert.EqualError(t, err, sr("GoError: Get HTTPBIN_URL/delay/10: net/http: request canceled (Client.Timeout exceeded while awaiting headers)"))
			assert.WithinDuration(t, startTime.Add(1*time.Second), endTime, 1*time.Second)

			logEntry := hook.LastEntry()
			if assert.NotNil(t, logEntry) ***REMOVED***
				assert.Equal(t, logrus.WarnLevel, logEntry.Level)
				assert.EqualError(t, logEntry.Data["error"].(error), sr("Get HTTPBIN_URL/delay/10: net/http: request canceled (Client.Timeout exceeded while awaiting headers)"))
				assert.Equal(t, "Request Failed", logEntry.Message)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("UserAgent", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
			let res = http.get("HTTPBIN_URL/user-agent");
			if (res.json()['user-agent'] != "TestUserAgent") ***REMOVED***
				throw new Error("incorrect user agent: " + res.json()['user-agent'])
			***REMOVED***
		`))
		assert.NoError(t, err)

		t.Run("Override", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.get("HTTPBIN_URL/user-agent", ***REMOVED***
					headers: ***REMOVED*** "User-Agent": "OtherUserAgent" ***REMOVED***,
				***REMOVED***);
				if (res.json()['user-agent'] != "OtherUserAgent") ***REMOVED***
					throw new Error("incorrect user agent: " + res.json()['user-agent'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("Compression", func(t *testing.T) ***REMOVED***
		t.Run("gzip", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.get("HTTPSBIN_IP_URL/gzip");
				if (res.json()['gzipped'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['gzipped'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("deflate", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.get("HTTPBIN_URL/deflate");
				if (res.json()['deflated'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['deflated'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
	t.Run("CompressionWithAcceptEncodingHeader", func(t *testing.T) ***REMOVED***
		t.Run("gzip", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let params = ***REMOVED*** headers: ***REMOVED*** "Accept-Encoding": "gzip" ***REMOVED*** ***REMOVED***;
				let res = http.get("HTTPBIN_URL/gzip", params);
				if (res.json()['gzipped'] != true) ***REMOVED***
					throw new Error("unexpected body data: " + res.json()['gzipped'])
				***REMOVED***
			`))
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("deflate", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let params = ***REMOVED*** headers: ***REMOVED*** "Accept-Encoding": "deflate" ***REMOVED*** ***REMOVED***;
				let res = http.get("HTTPBIN_URL/deflate", params);
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

		_, err := common.RunString(rt, sr(`http.get("HTTPBIN_URL/get/");`))
		assert.Error(t, err)
		assert.Nil(t, hook.LastEntry())
	***REMOVED***)
	t.Run("HTTP/2", func(t *testing.T) ***REMOVED***
		stats.GetBufferedSamples(samples) // Clean up buffered samples from previous tests

		_, err := common.RunString(rt, `
		let res = http.request("GET", "https://http2.akamai.com/demo");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
		if (res.proto != "HTTP/2.0") ***REMOVED*** throw new Error("wrong proto: " + res.proto) ***REMOVED***
		`)
		assert.NoError(t, err)

		bufSamples := stats.GetBufferedSamples(samples)
		assertRequestMetricsEmitted(t, bufSamples, "GET", "https://http2.akamai.com/demo", "", 200, "")
		for _, sampleC := range bufSamples ***REMOVED***
			for _, sample := range sampleC.GetSamples() ***REMOVED***
				proto, ok := sample.Tags.Get("proto")
				assert.True(t, ok)
				assert.Equal(t, "HTTP/2.0", proto)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	t.Run("TLS", func(t *testing.T) ***REMOVED***
		t.Run("cert_expired", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.get("https://expired.badssl.com/");`)
			assert.EqualError(t, err, "GoError: Get https://expired.badssl.com/: x509: certificate has expired or is not yet valid")
		***REMOVED***)
		tlsVersionTests := []struct ***REMOVED***
			Name, URL, Version string
		***REMOVED******REMOVED***
			***REMOVED***Name: "tls10", URL: "https://tls-v1-0.badssl.com:1010/", Version: "http.TLS_1_0"***REMOVED***,
			***REMOVED***Name: "tls11", URL: "https://tls-v1-1.badssl.com:1011/", Version: "http.TLS_1_1"***REMOVED***,
			***REMOVED***Name: "tls12", URL: "https://badssl.com/", Version: "http.TLS_1_2"***REMOVED***,
		***REMOVED***
		for _, versionTest := range tlsVersionTests ***REMOVED***
			t.Run(versionTest.Name, func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.get("%s");
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
			t.Run(cipherSuiteTest.Name, func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.get("%s");
					if (res.tls_cipher_suite != "%s") ***REMOVED*** throw new Error("wrong TLS cipher suite: " + res.tls_cipher_suite); ***REMOVED***
				`, cipherSuiteTest.URL, cipherSuiteTest.CipherSuite))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", cipherSuiteTest.URL, "", 200, "")
			***REMOVED***)
		***REMOVED***
		t.Run("ocsp_stapled_good", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let res = http.request("GET", "https://stackoverflow.com/");
			if (res.ocsp.status != http.OCSP_STATUS_GOOD) ***REMOVED*** throw new Error("wrong ocsp stapled response status: " + res.ocsp.status); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", "https://stackoverflow.com/", "", 200, "")
		***REMOVED***)
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		hook := logtest.NewLocal(state.Logger)
		defer hook.Reset()

		_, err := common.RunString(rt, `http.request("", "");`)
		assert.EqualError(t, err, "GoError: Get : unsupported protocol scheme \"\"")

		logEntry := hook.LastEntry()
		if assert.NotNil(t, logEntry) ***REMOVED***
			assert.Equal(t, logrus.WarnLevel, logEntry.Level)
			assert.Equal(t, "Get : unsupported protocol scheme \"\"", logEntry.Data["error"].(error).Error())
			assert.Equal(t, "Request Failed", logEntry.Message)
		***REMOVED***

		t.Run("throw=false", func(t *testing.T) ***REMOVED***
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			_, err := common.RunString(rt, `
				let res = http.request("", "", ***REMOVED*** throw: false ***REMOVED***);
				throw new Error(res.error);
			`)
			assert.EqualError(t, err, "GoError: Get : unsupported protocol scheme \"\"")

			logEntry := hook.LastEntry()
			if assert.NotNil(t, logEntry) ***REMOVED***
				assert.Equal(t, logrus.WarnLevel, logEntry.Level)
				assert.EqualError(t, logEntry.Data["error"].(error), "Get : unsupported protocol scheme \"\"")
				assert.Equal(t, "Request Failed", logEntry.Message)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Unroutable", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `http.request("GET", "http://sdafsgdhfjg/");`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("Params", func(t *testing.T) ***REMOVED***
		for _, literal := range []string***REMOVED***`undefined`, `null`***REMOVED*** ***REMOVED***
			t.Run(literal, func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, fmt.Sprintf(sr(`
				let res = http.request("GET", "HTTPBIN_URL/headers", null, %s);
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
				_, err = common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/cookies/set?key=value", null, ***REMOVED*** redirects: 0 ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.cookies.key[0].value); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies/set?key=value"), "", 302, "")
			***REMOVED***)

			t.Run("vuJar", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				let res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key2: "value2" ***REMOVED*** ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				if (res.json().key2 != "value2") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key2); ***REMOVED***
				let jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
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
				_, err = common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key: "value" ***REMOVED*** ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				let jar = http.cookieJar();
				let jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key != undefined) ***REMOVED*** throw new Error("unexpected cookie in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("requestScopeReplace", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				let res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key: ***REMOVED*** value: "replaced", replace: true ***REMOVED*** ***REMOVED*** ***REMOVED***);
				if (res.json().key != "replaced") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				let jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
				if (jarCookies.key[0] != "value") ***REMOVED*** throw new Error("wrong cookie value in jar"); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/cookies"), "", 200, "")
			***REMOVED***)

			t.Run("redirect", func(t *testing.T) ***REMOVED***
				t.Run("set cookie before redirect", func(t *testing.T) ***REMOVED***
					cookieJar, err := cookiejar.New(nil)
					assert.NoError(t, err)
					state.CookieJar = cookieJar
					_, err = common.RunString(rt, sr(`
						let res = http.request("GET", "HTTPBIN_URL/set-cookie-before-redirect");
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`))
					assert.NoError(t, err)

					redirectUrl, err := url.Parse(sr("HTTPBIN_URL"))
					assert.NoError(t, err)
					require.Len(t, cookieJar.Cookies(redirectUrl), 1)
					assert.Equal(t, "key-foo", cookieJar.Cookies(redirectUrl)[0].Name)
					assert.Equal(t, "value-bar", cookieJar.Cookies(redirectUrl)[0].Value)

					assertRequestMetricsEmitted(
						t,
						stats.GetBufferedSamples(samples),
						"GET",
						sr("HTTPBIN_URL/get"),
						sr("HTTPBIN_URL/set-cookie-before-redirect"),
						200,
						"",
					)
				***REMOVED***)
				t.Run("set cookie after redirect", func(t *testing.T) ***REMOVED***
					cookieJar, err := cookiejar.New(nil)
					assert.NoError(t, err)
					state.CookieJar = cookieJar
					_, err = common.RunString(rt, sr(`
						let res = http.request("GET", "HTTPBIN_URL/redirect-to?url=HTTPSBIN_URL/cookies/set?key=value");
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`))
					assert.NoError(t, err)

					redirectUrl, err := url.Parse(sr("HTTPSBIN_URL"))
					assert.NoError(t, err)

					require.Len(t, cookieJar.Cookies(redirectUrl), 1)
					assert.Equal(t, "key", cookieJar.Cookies(redirectUrl)[0].Name)
					assert.Equal(t, "value", cookieJar.Cookies(redirectUrl)[0].Value)

					assertRequestMetricsEmitted(
						t,
						stats.GetBufferedSamples(samples),
						"GET",
						sr("HTTPSBIN_URL/cookies"),
						sr("HTTPBIN_URL/redirect-to?url=HTTPSBIN_URL/cookies/set?key=value"),
						200,
						"",
					)
				***REMOVED***)

			***REMOVED***)

			t.Run("domain", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** domain: "HTTPBIN_DOMAIN" ***REMOVED***);
				let res = http.request("GET", "HTTPBIN_URL/cookies");
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
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** path: "/cookies" ***REMOVED***);
				let res = http.request("GET", "HTTPBIN_URL/cookies");
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
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value", ***REMOVED*** expires: "Sun, 24 Jul 1983 17:01:02 GMT" ***REMOVED***);
				let res = http.request("GET", "HTTPBIN_URL/cookies");
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
				_, err = common.RunString(rt, sr(`
				let jar = http.cookieJar();
				jar.set("HTTPSBIN_IP_URL/cookies", "key", "value", ***REMOVED*** secure: true ***REMOVED***);
				let res = http.request("GET", "HTTPSBIN_IP_URL/cookies");
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
				_, err = common.RunString(rt, sr(`
				let jar = new http.CookieJar();
				jar.set("HTTPBIN_URL/cookies", "key", "value");
				let res = http.request("GET", "HTTPBIN_URL/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key2: "value2" ***REMOVED***, jar: jar ***REMOVED***);
				if (res.json().key != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key); ***REMOVED***
				if (res.json().key2 != "value2") ***REMOVED*** throw new Error("wrong cookie value: " + res.json().key2); ***REMOVED***
				let jarCookies = jar.cookiesForURL("HTTPBIN_URL/cookies");
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

				_, err := common.RunString(rt, fmt.Sprintf(`
				let res = http.request("GET", "%s", null, ***REMOVED******REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`, url))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", url, "", 200, "")
			***REMOVED***)
			t.Run("digest", func(t *testing.T) ***REMOVED***
				t.Run("success", func(t *testing.T) ***REMOVED***
					url := sr("http://bob:pass@HTTPBIN_IP:HTTPBIN_PORT/digest-auth/auth/bob/pass")

					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.request("GET", "%s", null, ***REMOVED*** auth: "digest" ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`, url))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_IP_URL/digest-auth/auth/bob/pass"), url, 200, "")
				***REMOVED***)
				t.Run("failure", func(t *testing.T) ***REMOVED***
					url := sr("http://bob:pass@HTTPBIN_IP:HTTPBIN_PORT/digest-auth/failure")

					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.request("GET", "%s", null, ***REMOVED*** auth: "digest", timeout: 1, throw: false ***REMOVED***);
					`, url))
					assert.NoError(t, err)
				***REMOVED***)
			***REMOVED***)
			t.Run("ntlm", func(t *testing.T) ***REMOVED***
				t.Run("success auth", func(t *testing.T) ***REMOVED***
					url := strings.Replace(tb.ServerHTTP.URL+"/ntlm", "http://", "http://bob:pass@", -1)
					_, err := common.RunString(rt, fmt.Sprintf(`
						let res = http.request("GET", "%s", null, ***REMOVED*** auth: "ntlm" ***REMOVED***);
						if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
						`, url))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", url, url, 200, "")
				***REMOVED***)
				t.Run("failed auth", func(t *testing.T) ***REMOVED***
					url := strings.Replace(tb.ServerHTTP.URL+"/ntlm", "http://", "http://other:otherpass@", -1)
					_, err := common.RunString(rt, fmt.Sprintf(`
						let res = http.request("GET", "%s", null, ***REMOVED*** auth: "ntlm" ***REMOVED***);
						if (res.status != 401) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
						`, url))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", url, url, 401, "")
				***REMOVED***)
			***REMOVED***)
		***REMOVED***)

		t.Run("headers", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, fmt.Sprintf(sr(`
					let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** headers: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`), literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "X-My-Header": "value" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().headers["X-My-Header"] != "value") ***REMOVED*** throw new Error("wrong X-My-Header: " + res.json().headers["X-My-Header"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)

			t.Run("Host", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED***
					headers: ***REMOVED*** "Host": "HTTPBIN_DOMAIN" ***REMOVED***,
				***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().headers["Host"] != "HTTPBIN_DOMAIN") ***REMOVED*** throw new Error("wrong Host: " + res.json().headers["Host"]); ***REMOVED***
				`))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
			***REMOVED***)
		***REMOVED***)

		t.Run("tags", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, fmt.Sprintf(sr(`
					let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`), literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/headers"), "", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: ***REMOVED*** tag: "value" ***REMOVED*** ***REMOVED***);
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
				oldOpts := state.Options
				defer func() ***REMOVED*** state.Options = oldOpts ***REMOVED***()
				state.Options.RunTags = stats.IntoSampleTags(&map[string]string***REMOVED***"runtag1": "val1", "runtag2": "val2"***REMOVED***)

				_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/headers", null, ***REMOVED*** tags: ***REMOVED*** method: "test", name: "myName", runtag1: "fromreq" ***REMOVED*** ***REMOVED***);
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
		_, err := common.RunString(rt, sr(`
		let res = http.get("HTTPBIN_URL/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.json().args.a != "1") ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
		if (res.json().args.b != "2") ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get?a=1&b=2"), "", 200, "")

		t.Run("Tagged", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let a = "1";
			let b = "2";
			let res = http.get(http.url`+"`"+sr(`HTTPBIN_URL/get?a=$***REMOVED***a***REMOVED***&b=$***REMOVED***b***REMOVED***`)+"`"+`);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.json().args.a != a) ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
			if (res.json().args.b != b) ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get?a=1&b=2"), sr("HTTPBIN_URL/get?a=$***REMOVED******REMOVED***&b=$***REMOVED******REMOVED***"), 200, "")
		***REMOVED***)
	***REMOVED***)
	t.Run("HEAD", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = http.head("HTTPBIN_URL/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.length != 0) ***REMOVED*** throw new Error("HEAD responses shouldn't have a body"); ***REMOVED***
		if (!res.headers["Content-Length"]) ***REMOVED*** throw new Error("Missing or invalid Content-Length header!"); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "HEAD", sr("HTTPBIN_URL/get?a=1&b=2"), "", 200, "")
	***REMOVED***)

	t.Run("OPTIONS", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = http.options("HTTPBIN_URL/?a=1&b=2");
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
		_, err := common.RunString(rt, sr(`
		let res = http.del("HTTPBIN_URL/delete?test=mest");
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
			_, err := common.RunString(rt, fmt.Sprintf(sr(`
				let res = http.%s("HTTPBIN_URL/%s", "data");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().data != "data") ***REMOVED*** throw new Error("wrong data: " + res.json().data); ***REMOVED***
				if (res.json().headers["Content-Type"]) ***REMOVED*** throw new Error("content type set: " + res.json().headers["Content-Type"]); ***REMOVED***
				`), fn, strings.ToLower(method)))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), method, sr("HTTPBIN_URL/")+strings.ToLower(method), "", 200, "")

			t.Run("object", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, fmt.Sprintf(sr(`
				let res = http.%s("HTTPBIN_URL/%s", ***REMOVED***a: "a", b: 2***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.json().form.a != "a") ***REMOVED*** throw new Error("wrong a=: " + res.json().form.a); ***REMOVED***
				if (res.json().form.b != "2") ***REMOVED*** throw new Error("wrong b=: " + res.json().form.b); ***REMOVED***
				if (res.json().headers["Content-Type"] != "application/x-www-form-urlencoded") ***REMOVED*** throw new Error("wrong content type: " + res.json().headers["Content-Type"]); ***REMOVED***
				`), fn, strings.ToLower(method)))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), method, sr("HTTPBIN_URL/")+strings.ToLower(method), "", 200, "")
				t.Run("Content-Type", func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, fmt.Sprintf(sr(`
						let res = http.%s("HTTPBIN_URL/%s", ***REMOVED***a: "a", b: 2***REMOVED***, ***REMOVED***headers: ***REMOVED***"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"***REMOVED******REMOVED***);
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
		t.Run("GET", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
			let reqs = [
				["GET", "HTTPBIN_URL/"],
				["GET", "HTTPBIN_IP_URL/"],
			];
			let res = http.batch(reqs);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].url != reqs[key][1]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
			***REMOVED***`))
			assert.NoError(t, err)
			bufSamples := stats.GetBufferedSamples(samples)
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
			assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")

			t.Run("Tagged", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let fragment = "get";
				let reqs = [
					["GET", http.url`+"`"+`HTTPBIN_URL/$***REMOVED***fragment***REMOVED***`+"`"+`],
					["GET", http.url`+"`"+`HTTPBIN_IP_URL/`+"`"+`],
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key][1].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key][1].url); ***REMOVED***
				***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get"), sr("HTTPBIN_URL/$***REMOVED******REMOVED***"), 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
			***REMOVED***)

			t.Run("Shorthand", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let reqs = [
					"HTTPBIN_URL/",
					"HTTPBIN_IP_URL/",
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key]) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url); ***REMOVED***
				***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")

				t.Run("Tagged", func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, sr(`
					let fragment = "get";
					let reqs = [
						http.url`+"`"+`HTTPBIN_URL/$***REMOVED***fragment***REMOVED***`+"`"+`,
						http.url`+"`"+`HTTPBIN_IP_URL/`+"`"+`,
					];
					let res = http.batch(reqs);
					for (var key in res) ***REMOVED***
						if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
						if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key].url); ***REMOVED***
					***REMOVED***`))
					assert.NoError(t, err)
					bufSamples := stats.GetBufferedSamples(samples)
					assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/get"), sr("HTTPBIN_URL/$***REMOVED******REMOVED***"), 200, "")
					assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
				***REMOVED***)
			***REMOVED***)

			t.Run("ObjectForm", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let reqs = [
					***REMOVED*** method: "GET", url: "HTTPBIN_URL/" ***REMOVED***,
					***REMOVED*** url: "HTTPBIN_IP_URL/", method: "GET"***REMOVED***,
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + key + ": " + res[key].url + " != " + reqs[key].url); ***REMOVED***
				***REMOVED***`))
				assert.NoError(t, err)
				bufSamples := stats.GetBufferedSamples(samples)
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_URL/"), "", 200, "")
				assertRequestMetricsEmitted(t, bufSamples, "GET", sr("HTTPBIN_IP_URL/"), "", 200, "")
			***REMOVED***)

			t.Run("ObjectKeys", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
				let reqs = ***REMOVED***
					shorthand: "HTTPBIN_URL/get?r=shorthand",
					arr: ["GET", "HTTPBIN_URL/get?r=arr", null, ***REMOVED***tags: ***REMOVED***name: 'arr'***REMOVED******REMOVED***],
					obj1: ***REMOVED*** method: "GET", url: "HTTPBIN_URL/get?r=obj1" ***REMOVED***,
					obj2: ***REMOVED*** url: "HTTPBIN_URL/get?r=obj2", params: ***REMOVED***tags: ***REMOVED***name: 'obj2'***REMOVED******REMOVED***, method: "GET"***REMOVED***,
				***REMOVED***;
				let res = http.batch(reqs);
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

				_, err := common.RunString(rt, sr(`
					let reqs = [
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
					let res = http.batch(reqs);
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
			_, err := common.RunString(rt, sr(`
			let res = http.batch([ ["POST", "HTTPBIN_URL/post", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + key + ": " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + key + ": " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)
		t.Run("PUT", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
			let res = http.batch([ ["PUT", "HTTPBIN_URL/put", ***REMOVED*** key: "value" ***REMOVED***] ]);
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
			_, err := common.RunString(rt, sr(`
				let reqUrl = "HTTPBIN_URL/cookies"
				let res = http.get(reqUrl);
				let jar = new http.CookieJar();

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
			_, err := common.RunString(rt, sr(`
				let res = http.post("HTTPBIN_URL/post", ***REMOVED***a: "a", b: 2***REMOVED***, ***REMOVED***headers: ***REMOVED***"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"***REMOVED******REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.request["body"] != "a=a&b=2") ***REMOVED*** throw new Error("http request body was not set properly: " + JSON.stringify(res.request))***REMOVED***
				`))
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
func TestSystemTags(t *testing.T) ***REMOVED***
	tb, state, samples, rt, _ := newRuntime(t)
	defer tb.Cleanup()

	// Handple paths with custom logic
	tb.Mux.HandleFunc("/wrong-redirect", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Add("Location", "%")
		w.WriteHeader(http.StatusTemporaryRedirect)
	***REMOVED***)

	httpGet := fmt.Sprintf(`http.get("%s");`, tb.ServerHTTP.URL)
	httpsGet := fmt.Sprintf(`http.get("%s");`, tb.ServerHTTPS.URL)

	httpURL, err := url.Parse(tb.ServerHTTP.URL)
	require.NoError(t, err)

	testedSystemTags := []struct***REMOVED*** tag, code, expVal string ***REMOVED******REMOVED***
		***REMOVED***"proto", httpGet, "HTTP/1.1"***REMOVED***,
		***REMOVED***"status", httpGet, "200"***REMOVED***,
		***REMOVED***"method", httpGet, "GET"***REMOVED***,
		***REMOVED***"url", httpGet, tb.ServerHTTP.URL***REMOVED***,
		***REMOVED***"url", httpsGet, tb.ServerHTTPS.URL***REMOVED***,
		***REMOVED***"ip", httpGet, httpURL.Hostname()***REMOVED***,
		***REMOVED***"name", httpGet, tb.ServerHTTP.URL***REMOVED***,
		***REMOVED***"group", httpGet, ""***REMOVED***,
		***REMOVED***"vu", httpGet, "0"***REMOVED***,
		***REMOVED***"iter", httpGet, "0"***REMOVED***,
		***REMOVED***"tls_version", httpsGet, "tls1.2"***REMOVED***,
		***REMOVED***"ocsp_status", httpsGet, "unknown"***REMOVED***,
		***REMOVED***
			"error",
			tb.Replacer.Replace(`http.get("HTTPBIN_IP_URL/wrong-redirect");`),
			tb.Replacer.Replace(`Get HTTPBIN_IP_URL/wrong-redirect: failed to parse Location header "%": parse %: invalid URL escape "%"`),
		***REMOVED***,
	***REMOVED***

	state.Options.Throw = null.BoolFrom(false)

	for num, tc := range testedSystemTags ***REMOVED***
		t.Run(fmt.Sprintf("TC %d with only %s", num, tc.tag), func(t *testing.T) ***REMOVED***
			state.Options.SystemTags = lib.GetTagSet(tc.tag)
			_, err := common.RunString(rt, tc.code)
			assert.NoError(t, err)

			bufSamples := stats.GetBufferedSamples(samples)
			assert.NotEmpty(t, bufSamples)
			for _, sampleC := range bufSamples ***REMOVED***
				for _, sample := range sampleC.GetSamples() ***REMOVED***
					assert.NotEmpty(t, sample.Tags)
					for emittedTag, emittedVal := range sample.Tags.CloneTags() ***REMOVED***
						assert.Equal(t, tc.tag, emittedTag)
						assert.Equal(t, tc.expVal, emittedVal)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Simple NTLM mock handler
func ntlmHandler(username, password string) func(w http.ResponseWriter, r *http.Request) ***REMOVED***
	challenges := make(map[string]*ntlm.ChallengeMessage)
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Make sure there is some kind of authentication
		if r.Header.Get("Authorization") == "" ***REMOVED***
			w.Header().Set("WWW-Authenticate", "NTLM")
			w.WriteHeader(401)
			return
		***REMOVED***

		// Parse the proxy authorization header
		auth := r.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		authType := parts[0]
		authPayload := parts[1]

		// Filter out unsupported authentication methods
		if authType != "NTLM" ***REMOVED***
			w.Header().Set("WWW-Authenticate", "NTLM")
			w.WriteHeader(401)
			return
		***REMOVED***

		// Decode base64 auth data and get NTLM message type
		rawAuthPayload, _ := base64.StdEncoding.DecodeString(authPayload)
		ntlmMessageType := binary.LittleEndian.Uint32(rawAuthPayload[8:12])

		// Handle NTLM negotiate message
		if ntlmMessageType == 1 ***REMOVED***
			session, err := ntlm.CreateServerSession(ntlm.Version2, ntlm.ConnectionOrientedMode)
			if err != nil ***REMOVED***
				return
			***REMOVED***

			session.SetUserInfo(username, password, "")

			challenge, err := session.GenerateChallengeMessage()
			if err != nil ***REMOVED***
				return
			***REMOVED***

			challenges[r.RemoteAddr] = challenge

			authPayload := base64.StdEncoding.EncodeToString(challenge.Bytes())

			w.Header().Set("WWW-Authenticate", "NTLM "+authPayload)
			w.WriteHeader(401)

			return
		***REMOVED***

		if ntlmMessageType == 3 ***REMOVED***
			challenge := challenges[r.RemoteAddr]
			if challenge == nil ***REMOVED***
				w.Header().Set("WWW-Authenticate", "NTLM")
				w.WriteHeader(401)
				return
			***REMOVED***

			msg, err := ntlm.ParseAuthenticateMessage(rawAuthPayload, 2)
			if err != nil ***REMOVED***
				msg2, err := ntlm.ParseAuthenticateMessage(rawAuthPayload, 1)

				if err != nil ***REMOVED***
					return
				***REMOVED***

				session, err := ntlm.CreateServerSession(ntlm.Version1, ntlm.ConnectionOrientedMode)
				if err != nil ***REMOVED***
					return
				***REMOVED***

				session.SetServerChallenge(challenge.ServerChallenge)
				session.SetUserInfo(username, password, "")

				err = session.ProcessAuthenticateMessage(msg2)
				if err != nil ***REMOVED***
					w.Header().Set("WWW-Authenticate", "NTLM")
					w.WriteHeader(401)
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				session, err := ntlm.CreateServerSession(ntlm.Version2, ntlm.ConnectionOrientedMode)
				if err != nil ***REMOVED***
					return
				***REMOVED***

				session.SetServerChallenge(challenge.ServerChallenge)
				session.SetUserInfo(username, password, "")

				err = session.ProcessAuthenticateMessage(msg)
				if err != nil ***REMOVED***
					w.Header().Set("WWW-Authenticate", "NTLM")
					w.WriteHeader(401)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***

		data := "authenticated"
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		if _, err := fmt.Fprint(w, data); err != nil ***REMOVED***
			panic(err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***
