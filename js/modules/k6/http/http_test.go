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
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/oxtoacart/bpool"
	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func assertRequestMetricsEmitted(t *testing.T, samples []stats.Sample, method, url, name string, status int, group string) ***REMOVED***
	if name == "" ***REMOVED***
		name = url
	***REMOVED***

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
			assert.Equal(t, name, sample.Tags["name"])
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

	logger := log.New()
	logger.Level = log.DebugLevel
	logger.Out = ioutil.Discard

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	state := &common.State***REMOVED***
		Options: lib.Options***REMOVED***
			MaxRedirects: null.IntFrom(10),
			UserAgent:    null.StringFrom("TestUserAgent"),
			Throw:        null.BoolFrom(true),
		***REMOVED***,
		Logger: logger,
		Group:  root,
		HTTPTransport: &http.Transport***REMOVED***
			DialContext: (netext.NewDialer(net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***)).DialContext,
		***REMOVED***,
		BPool: bpool.NewBufferPool(1),
	***REMOVED***

	ctx := new(context.Context)
	*ctx = context.Background()
	*ctx = common.WithState(*ctx, state)
	*ctx = common.WithRuntime(*ctx, rt)
	rt.Set("http", common.Bind(rt, New(), ctx))

	t.Run("Redirects", func(t *testing.T) ***REMOVED***
		t.Run("10", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.get("https://httpbin.org/redirect/10")`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("11", func(t *testing.T) ***REMOVED***
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			_, err := common.RunString(rt, `
			let res = http.get("https://httpbin.org/redirect/11");
			if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "https://httpbin.org/relative-redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
			`)
			assert.NoError(t, err)

			logEntry := hook.LastEntry()
			if assert.NotNil(t, logEntry) ***REMOVED***
				assert.Equal(t, log.WarnLevel, logEntry.Level)
				assert.Equal(t, "Possible redirect loop, 302 response returned last, 10 redirects followed; pass ***REMOVED*** redirects: n ***REMOVED*** in request params to silence this", logEntry.Data["error"])
				assert.Equal(t, "https://httpbin.org/redirect/11", logEntry.Data["url"])
				assert.Equal(t, "Redirect Limit", logEntry.Message)
			***REMOVED***
		***REMOVED***)
		t.Run("requestScopeRedirects", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let res = http.get("https://httpbin.org/redirect/1", ***REMOVED***redirects: 3***REMOVED***);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "https://httpbin.org/get") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("requestScopeNoRedirects", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
			let res = http.get("https://httpbin.org/redirect/1", ***REMOVED***redirects: 0***REMOVED***);
			if (res.status != 302) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			if (res.url != "https://httpbin.org/redirect/1") ***REMOVED*** throw new Error("incorrect URL: " + res.url) ***REMOVED***
			if (res.headers["Location"] != "/get") ***REMOVED*** throw new Error("incorrect Location header: " + res.headers["Location"]) ***REMOVED***
			`)
			assert.NoError(t, err)
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
			hook := logtest.NewLocal(state.Logger)
			defer hook.Reset()

			startTime := time.Now()
			_, err := common.RunString(rt, `
				http.get("https://httpbin.org/delay/10", ***REMOVED***
					timeout: 1*1000,
				***REMOVED***)
			`)
			endTime := time.Now()
			assert.EqualError(t, err, "GoError: Get https://httpbin.org/delay/10: net/http: request canceled (Client.Timeout exceeded while awaiting headers)")
			assert.WithinDuration(t, startTime.Add(1*time.Second), endTime, 1*time.Second)

			logEntry := hook.LastEntry()
			if assert.NotNil(t, logEntry) ***REMOVED***
				assert.Equal(t, log.WarnLevel, logEntry.Level)
				assert.EqualError(t, logEntry.Data["error"].(error), "Get https://httpbin.org/delay/10: net/http: request canceled (Client.Timeout exceeded while awaiting headers)")
				assert.Equal(t, "Request Failed", logEntry.Message)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("UserAgent", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
			let res = http.get("http://httpbin.org/user-agent");
			if (res.json()['user-agent'] != "TestUserAgent") ***REMOVED***
				throw new Error("incorrect user agent: " + res.json()['user-agent'])
			***REMOVED***
		`)
		assert.NoError(t, err)

		t.Run("Override", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
				let res = http.get("http://httpbin.org/user-agent", ***REMOVED***
					headers: ***REMOVED*** "User-Agent": "OtherUserAgent" ***REMOVED***,
				***REMOVED***);
				if (res.json()['user-agent'] != "OtherUserAgent") ***REMOVED***
					throw new Error("incorrect user agent: " + res.json()['user-agent'])
				***REMOVED***
			`)
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

		_, err := common.RunString(rt, `http.get("https://httpbin.org/get/");`)
		assert.Error(t, err)
		assert.Nil(t, hook.LastEntry())
	***REMOVED***)

	t.Run("HTTP/2", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.request("GET", "https://http2.akamai.com/demo");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
		if (res.proto != "HTTP/2.0") ***REMOVED*** throw new Error("wrong proto: " + res.proto) ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://http2.akamai.com/demo", "", 200, "")
		for _, sample := range state.Samples ***REMOVED***
			assert.Equal(t, "HTTP/2.0", sample.Tags["proto"])
		***REMOVED***
	***REMOVED***)

	t.Run("HTML", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.request("GET", "https://httpbin.org/html");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.indexOf("Herman Melville - Moby-Dick") == -1) ***REMOVED*** throw new Error("wrong body: " + res.body); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/html", "", 200, "")

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
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/html", "", 200, "::my group")
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
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get?a=1&b=2", "", 200, "")

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `http.request("GET", "https://httpbin.org/html").json();`)
			assert.EqualError(t, err, "GoError: invalid character '<' looking for beginning of value")
		***REMOVED***)
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
				assertRequestMetricsEmitted(t, state.Samples, "GET", versionTest.URL, "", 200, "")
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
				assertRequestMetricsEmitted(t, state.Samples, "GET", cipherSuiteTest.URL, "", 200, "")
			***REMOVED***)
		***REMOVED***
		t.Run("ocsp_stapled_good", func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, `
			let res = http.request("GET", "https://stackoverflow.com/");
			if (res.ocsp.status != http.OCSP_STATUS_GOOD) ***REMOVED*** throw new Error("wrong ocsp stapled response status: " + res.ocsp.status); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://stackoverflow.com/", "", 200, "")
		***REMOVED***)
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		hook := logtest.NewLocal(state.Logger)
		defer hook.Reset()

		_, err := common.RunString(rt, `http.request("", "");`)
		assert.EqualError(t, err, "GoError: Get : unsupported protocol scheme \"\"")

		logEntry := hook.LastEntry()
		if assert.NotNil(t, logEntry) ***REMOVED***
			assert.Equal(t, log.WarnLevel, logEntry.Level)
			assert.EqualError(t, logEntry.Data["error"].(error), "Get : unsupported protocol scheme \"\"")
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
				assert.Equal(t, log.WarnLevel, logEntry.Level)
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
				state.Samples = nil
				_, err := common.RunString(rt, fmt.Sprintf(`
				let res = http.request("GET", "https://httpbin.org/headers", null, %s);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`, literal))
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", "", 200, "")
			***REMOVED***)
		***REMOVED***

		t.Run("cookies", func(t *testing.T) ***REMOVED***
			t.Run("access", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/cookies/set?key=value", null);
				if (res.cookies.key[0].value != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.cookies.key[0]); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "https://httpbin.org/cookies/set?key=value", 200, "")
			***REMOVED***)

			t.Run("setting", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [***REMOVED*** name: "key", value: "value" ***REMOVED***] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.cookies.key[0]); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)

			t.Run("settingSimple", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: ***REMOVED*** key: "value" ***REMOVED*** ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED*** throw new Error("wrong cookie value: " + res.cookies.key[0]); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)

			t.Run("domain", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let cookie = ***REMOVED*** name: "key", value: "value", domain: "httpbin.org" ***REMOVED***;
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.cookies.key[0]);
				***REMOVED***
				cookie = ***REMOVED*** name: "key2", value: "value2", domain: "example.com" ***REMOVED***;
				res = http.request("GET", "http://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.cookies.key[0]);
				***REMOVED***
				if (res.cookies.key2 != undefined) ***REMOVED***
					throw new Error("cookie 'key2' unexpectedly found");
				***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)

			t.Run("path", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let cookie = ***REMOVED*** name: "key", value: "value", path: "/cookies" ***REMOVED***;
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.cookies.key[0]);
				***REMOVED***
				cookie = ***REMOVED*** name: "key2", value: "value2", path: "/some-other-path" ***REMOVED***;
				res = http.request("GET", "http://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.cookies.key[0]);
				***REMOVED***
				if (res.cookies.key2 != undefined) ***REMOVED***
					throw new Error("cookie 'key2' unexpectedly found");
				***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)

			t.Run("expires", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let cookie = ***REMOVED*** name: "key", value: "value", expires: "Sun, 24 Jul 1983 17:01:02 GMT" ***REMOVED***;
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key != undefined) ***REMOVED***
					throw new Error("cookie 'key' unexpectedly found");
				***REMOVED***
				cookie.expires = "Sat, 24 Jul 2083 17:01:02 GMT";
				res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("cookie 'key' not found");
				***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)

			t.Run("secure", func(t *testing.T) ***REMOVED***
				cookieJar, err := cookiejar.New(nil)
				assert.NoError(t, err)
				state.CookieJar = cookieJar
				state.Samples = nil
				_, err = common.RunString(rt, `
				let cookie = ***REMOVED*** name: "key", value: "value", secure: true ***REMOVED***;
				let res = http.request("GET", "https://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (res.cookies.key[0].value != "value") ***REMOVED***
					throw new Error("wrong cookie value: " + res.cookies.key[0]);
				***REMOVED***
				res = http.request("GET", "http://httpbin.org/cookies", null, ***REMOVED*** cookies: [cookie] ***REMOVED***);
				if (Object.keys(res.cookies).length != 0) ***REMOVED***
					throw new Error("no cookies should've been sent");
				***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/cookies", "", 200, "")
			***REMOVED***)
		***REMOVED***)

		t.Run("headers", func(t *testing.T) ***REMOVED***
			for _, literal := range []string***REMOVED***`null`, `undefined`***REMOVED*** ***REMOVED***
				state.Samples = nil
				t.Run(literal, func(t *testing.T) ***REMOVED***
					_, err := common.RunString(rt, fmt.Sprintf(`
					let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED*** headers: %s ***REMOVED***);
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
					`, literal))
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", "", 200, "")
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
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", "", 200, "")
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
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", "", 200, "")
				***REMOVED***)
			***REMOVED***

			t.Run("object", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let res = http.request("GET", "https://httpbin.org/headers", null, ***REMOVED*** tags: ***REMOVED*** tag: "value" ***REMOVED*** ***REMOVED***);
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/headers", "", 200, "")
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
		assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get?a=1&b=2", "", 200, "")

		t.Run("Tagged", func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, `
			let a = "1";
			let b = "2";
			let res = http.get(http.url`+"`"+`https://httpbin.org/get?a=$***REMOVED***a***REMOVED***&b=$***REMOVED***b***REMOVED***`+"`"+`);
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.json().args.a != a) ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
			if (res.json().args.b != b) ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
			`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get?a=1&b=2", "https://httpbin.org/get?a=$***REMOVED******REMOVED***&b=$***REMOVED******REMOVED***", 200, "")
		***REMOVED***)
	***REMOVED***)
	t.Run("HEAD", func(t *testing.T) ***REMOVED***
		state.Samples = nil
		_, err := common.RunString(rt, `
		let res = http.head("https://httpbin.org/get?a=1&b=2");
		if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.body.length != 0) ***REMOVED*** throw new Error("HEAD responses shouldn't have a body"); ***REMOVED***
		`)
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, state.Samples, "HEAD", "https://httpbin.org/get?a=1&b=2", "", 200, "")
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
			assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), "", 200, "")

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
				assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), "", 200, "")

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
					assertRequestMetricsEmitted(t, state.Samples, method, "https://httpbin.org/"+strings.ToLower(method), "", 200, "")
				***REMOVED***)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***

	t.Run("Batch", func(t *testing.T) ***REMOVED***
		t.Run("GET", func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, `
			let reqs = [
				["GET", "https://httpbin.org/"],
				["GET", "https://now.httpbin.org/"],
			];
			let res = http.batch(reqs);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].url != reqs[key][1]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/", "", 200, "")
			assertRequestMetricsEmitted(t, state.Samples, "GET", "https://now.httpbin.org/", "", 200, "")

			t.Run("Tagged", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let fragment = "get";
				let reqs = [
					["GET", http.url`+"`"+`https://httpbin.org/$***REMOVED***fragment***REMOVED***`+"`"+`],
					["GET", http.url`+"`"+`https://now.httpbin.org/`+"`"+`],
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key][1].url) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
				***REMOVED***`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get", "https://httpbin.org/$***REMOVED******REMOVED***", 200, "")
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://now.httpbin.org/", "", 200, "")
			***REMOVED***)

			t.Run("Shorthand", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let reqs = [
					"https://httpbin.org/",
					"https://now.httpbin.org/",
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key]) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
				***REMOVED***`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/", "", 200, "")
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://now.httpbin.org/", "", 200, "")

				t.Run("Tagged", func(t *testing.T) ***REMOVED***
					state.Samples = nil
					_, err := common.RunString(rt, `
					let fragment = "get";
					let reqs = [
						http.url`+"`"+`https://httpbin.org/$***REMOVED***fragment***REMOVED***`+"`"+`,
						http.url`+"`"+`https://now.httpbin.org/`+"`"+`,
					];
					let res = http.batch(reqs);
					for (var key in res) ***REMOVED***
						if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
						if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
					***REMOVED***`)
					assert.NoError(t, err)
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/get", "https://httpbin.org/$***REMOVED******REMOVED***", 200, "")
					assertRequestMetricsEmitted(t, state.Samples, "GET", "https://now.httpbin.org/", "", 200, "")
				***REMOVED***)
			***REMOVED***)

			t.Run("ObjectForm", func(t *testing.T) ***REMOVED***
				state.Samples = nil
				_, err := common.RunString(rt, `
				let reqs = [
					***REMOVED*** url: "https://httpbin.org/", method: "GET" ***REMOVED***,
					***REMOVED*** method: "GET", url: "https://now.httpbin.org/" ***REMOVED***,
				];
				let res = http.batch(reqs);
				for (var key in res) ***REMOVED***
					if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
					if (res[key].url != reqs[key].url) ***REMOVED*** throw new Error("wrong url: " + res[key].url); ***REMOVED***
				***REMOVED***`)
				assert.NoError(t, err)
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://httpbin.org/", "", 200, "")
				assertRequestMetricsEmitted(t, state.Samples, "GET", "https://now.httpbin.org/", "", 200, "")
			***REMOVED***)
		***REMOVED***)
		t.Run("POST", func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, `
			let res = http.batch([ ["POST", "https://httpbin.org/post", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "POST", "https://httpbin.org/post", "", 200, "")
		***REMOVED***)
		t.Run("PUT", func(t *testing.T) ***REMOVED***
			state.Samples = nil
			_, err := common.RunString(rt, `
			let res = http.batch([ ["PUT", "https://httpbin.org/put", ***REMOVED*** key: "value" ***REMOVED***] ]);
			for (var key in res) ***REMOVED***
				if (res[key].status != 200) ***REMOVED*** throw new Error("wrong status: " + res[key].status); ***REMOVED***
				if (res[key].json().form.key != "value") ***REMOVED*** throw new Error("wrong form: " + JSON.stringify(res[key].json().form)); ***REMOVED***
			***REMOVED***`)
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, state.Samples, "PUT", "https://httpbin.org/put", "", 200, "")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestTagURL(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	rt.Set("http", common.Bind(rt, New(), nil))

	testdata := map[string]URLTag***REMOVED***
		`http://httpbin.org/anything/`:               ***REMOVED***URL: "http://httpbin.org/anything/", Name: "http://httpbin.org/anything/"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***`:         ***REMOVED***URL: "http://httpbin.org/anything/2", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/`:        ***REMOVED***URL: "http://httpbin.org/anything/2/", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***`:  ***REMOVED***URL: "http://httpbin.org/anything/2/3", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***/`: ***REMOVED***URL: "http://httpbin.org/anything/2/3/", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***/"***REMOVED***,
	***REMOVED***
	for expr, tag := range testdata ***REMOVED***
		t.Run("expr="+expr, func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, "http.url`"+expr+"`")
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, tag, v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
