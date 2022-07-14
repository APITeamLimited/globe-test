/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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
package ws

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	httpModule "go.k6.io/k6/js/modules/k6/http"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/metrics"
)

const statusProtocolSwitch = 101

func assertSessionMetricsEmitted(t *testing.T, sampleContainers []metrics.SampleContainer, subprotocol, url string, status int, group string) ***REMOVED*** //nolint:unparam
	seenSessions := false
	seenSessionDuration := false
	seenConnecting := false

	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			tags := sample.Tags.CloneTags()
			if tags["url"] == url ***REMOVED***
				switch sample.Metric.Name ***REMOVED***
				case metrics.WSConnectingName:
					seenConnecting = true
				case metrics.WSSessionDurationName:
					seenSessionDuration = true
				case metrics.WSSessionsName:
					seenSessions = true
				***REMOVED***

				assert.Equal(t, strconv.Itoa(status), tags["status"])
				assert.Equal(t, subprotocol, tags["subproto"])
				assert.Equal(t, group, tags["group"])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.True(t, seenConnecting, "url %s didn't emit Connecting", url)
	assert.True(t, seenSessions, "url %s didn't emit Sessions", url)
	assert.True(t, seenSessionDuration, "url %s didn't emit SessionDuration", url)
***REMOVED***

func assertMetricEmittedCount(t *testing.T, metricName string, sampleContainers []metrics.SampleContainer, url string, count int) ***REMOVED***
	t.Helper()
	actualCount := 0

	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			surl, ok := sample.Tags.Get("url")
			assert.True(t, ok)
			if surl == url && sample.Metric.Name == metricName ***REMOVED***
				actualCount++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.Equal(t, count, actualCount, "url %s emitted %s %d times, expected was %d times", url, metricName, actualCount, count)
***REMOVED***

type testState struct ***REMOVED***
	ctxPtr  *context.Context
	rt      *goja.Runtime
	tb      *httpmultibin.HTTPMultiBin
	state   *lib.State
	samples chan metrics.SampleContainer
***REMOVED***

func newTestState(t testing.TB) testState ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	samples := make(chan metrics.SampleContainer, 1000)

	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
			UserAgent: null.StringFrom("TestUserAgent"),
		***REMOVED***,
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     tb.Context,
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	return testState***REMOVED***
		rt:      rt,
		tb:      tb,
		state:   state,
		samples: samples,
	***REMOVED***
***REMOVED***

func TestSession(t *testing.T) ***REMOVED***
	// TODO: split and paralelize tests
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
		***REMOVED***,
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	t.Run("connect_ws", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("connection failed with status: " + res.status); ***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")

	t.Run("connect_wss", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = ws.connect("WSSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-echo"), statusProtocolSwitch, "")

	t.Run("open", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var opened = false;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				opened = true;
				socket.close()
			***REMOVED***)
		***REMOVED***);
		if (!opened) ***REMOVED*** throw new Error ("open event not fired"); ***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")

	t.Run("send_receive", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				socket.send("test")
			***REMOVED***)
			socket.on("message", function (data) ***REMOVED***
				if (!data=="test") ***REMOVED***
					throw new Error ("echo'd data doesn't match our message!");
				***REMOVED***
				socket.close()
			***REMOVED***);
		***REMOVED***);
		`))
		require.NoError(t, err)
	***REMOVED***)

	samplesBuf := metrics.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")
	assertMetricEmittedCount(t, metrics.WSMessagesSentName, samplesBuf, sr("WSBIN_URL/ws-echo"), 1)
	assertMetricEmittedCount(t, metrics.WSMessagesReceivedName, samplesBuf, sr("WSBIN_URL/ws-echo"), 1)

	t.Run("interval", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var counter = 0;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setInterval(function () ***REMOVED***
				counter += 1;
				if (counter > 2) ***REMOVED*** socket.close(); ***REMOVED***
			***REMOVED***, 100);
		***REMOVED***);
		if (counter < 3) ***REMOVED***throw new Error ("setInterval should have been called at least 3 times, counter=" + counter);***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")
	t.Run("bad interval", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var counter = 0;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setInterval(function () ***REMOVED***
				counter += 1;
				if (counter > 2) ***REMOVED*** socket.close(); ***REMOVED***
			***REMOVED***, -1.23);
		***REMOVED***);
		`))
		require.Error(t, err)
		require.Contains(t, err.Error(), "setInterval requires a >0 timeout parameter, received -1.23 ")
	***REMOVED***)

	t.Run("timeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var start = new Date().getTime();
		var ellapsed = new Date().getTime() - start;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setTimeout(function () ***REMOVED***
				ellapsed = new Date().getTime() - start;
				socket.close();
			***REMOVED***, 500);
		***REMOVED***);
		if (ellapsed > 3000 || ellapsed < 500) ***REMOVED***
			throw new Error ("setTimeout occurred after " + ellapsed + "ms, expected 500<T<3000");
		***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)

	t.Run("bad timeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var start = new Date().getTime();
		var ellapsed = new Date().getTime() - start;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setTimeout(function () ***REMOVED***
				ellapsed = new Date().getTime() - start;
				socket.close();
			***REMOVED***, 0);
		***REMOVED***);
		`))
		require.Error(t, err)
		require.Contains(t, err.Error(), "setTimeout requires a >0 timeout parameter, received 0.00 ")
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")

	t.Run("ping", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var pongReceived = false;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function(data) ***REMOVED***
				socket.ping();
			***REMOVED***);
			socket.on("pong", function() ***REMOVED***
				pongReceived = true;
				socket.close();
			***REMOVED***);
			socket.setTimeout(function ()***REMOVED***socket.close();***REMOVED***, 3000);
		***REMOVED***);
		if (!pongReceived) ***REMOVED***
			throw new Error ("sent ping but didn't get pong back");
		***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)

	samplesBuf = metrics.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")
	assertMetricEmittedCount(t, metrics.WSPingName, samplesBuf, sr("WSBIN_URL/ws-echo"), 1)

	t.Run("multiple_handlers", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var pongReceived = false;
		var otherPongReceived = false;

		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function(data) ***REMOVED***
				socket.ping();
			***REMOVED***);
			socket.on("pong", function() ***REMOVED***
				pongReceived = true;
				if (otherPongReceived) ***REMOVED***
					socket.close();
				***REMOVED***
			***REMOVED***);
			socket.on("pong", function() ***REMOVED***
				otherPongReceived = true;
				if (pongReceived) ***REMOVED***
					socket.close();
				***REMOVED***
			***REMOVED***);
			socket.setTimeout(function ()***REMOVED***socket.close();***REMOVED***, 3000);
		***REMOVED***);
		if (!pongReceived || !otherPongReceived) ***REMOVED***
			throw new Error ("sent ping but didn't get pong back");
		***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)

	samplesBuf = metrics.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")
	assertMetricEmittedCount(t, metrics.WSPingName, samplesBuf, sr("WSBIN_URL/ws-echo"), 1)

	t.Run("client_close", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var closed = false;
		var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
							socket.close()
			***REMOVED***)
			socket.on("close", function() ***REMOVED***
							closed = true;
			***REMOVED***)
		***REMOVED***);
		if (!closed) ***REMOVED*** throw new Error ("close event not fired"); ***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), statusProtocolSwitch, "")

	serverCloseTests := []struct ***REMOVED***
		name     string
		endpoint string
	***REMOVED******REMOVED***
		***REMOVED***"server_close_ok", "/ws-echo"***REMOVED***,
		// Ensure we correctly handle invalid WS server
		// implementations that close the connection prematurely
		// without sending a close control frame first.
		***REMOVED***"server_close_invalid", "/ws-close-invalid"***REMOVED***,
	***REMOVED***

	for _, tc := range serverCloseTests ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(fmt.Sprintf(`
			var closed = false;
			var res = ws.connect("WSBIN_URL%s", function(socket)***REMOVED***
				socket.on("open", function() ***REMOVED***
					socket.send("test");
				***REMOVED***)
				socket.on("close", function() ***REMOVED***
					closed = true;
				***REMOVED***)
			***REMOVED***);
			if (!closed) ***REMOVED*** throw new Error ("close event not fired"); ***REMOVED***
			`, tc.endpoint)))
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***

	t.Run("multi_message", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tb.Mux.HandleFunc("/ws-echo-multi", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, req, w.Header())
			if err != nil ***REMOVED***
				return
			***REMOVED***

			for ***REMOVED***
				messageType, r, e := conn.NextReader()
				if e != nil ***REMOVED***
					return
				***REMOVED***
				var wc io.WriteCloser
				wc, err = conn.NextWriter(messageType)
				if err != nil ***REMOVED***
					return
				***REMOVED***
				if _, err = io.Copy(wc, r); err != nil ***REMOVED***
					return
				***REMOVED***
				if err = wc.Close(); err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***))

		t.Run("send_receive_multiple_ws", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var msg1 = "test1"
			var msg2 = "test2"
			var msg3 = "test3"
			var allMsgsRecvd = false
			var res = ws.connect("WSBIN_URL/ws-echo-multi", (socket) => ***REMOVED***
				socket.on("open", () => ***REMOVED***
					socket.send(msg1)
				***REMOVED***)
				socket.on("message", (data) => ***REMOVED***
					if (data == msg1)***REMOVED***
						socket.send(msg2)
					***REMOVED***
					if (data == msg2)***REMOVED***
						socket.send(msg3)
					***REMOVED***
					if (data == msg3)***REMOVED***
						allMsgsRecvd = true
						socket.close()
					***REMOVED***
				***REMOVED***);
			***REMOVED***);

			if (!allMsgsRecvd) ***REMOVED***
				throw new Error ("messages 1,2,3 in sequence, was not received from server");
			***REMOVED***
			`))
			require.NoError(t, err)
		***REMOVED***)

		samplesBuf = metrics.GetBufferedSamples(samples)
		assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo-multi"), statusProtocolSwitch, "")
		assertMetricEmittedCount(t, metrics.WSMessagesSentName, samplesBuf, sr("WSBIN_URL/ws-echo-multi"), 3)
		assertMetricEmittedCount(t, metrics.WSMessagesReceivedName, samplesBuf, sr("WSBIN_URL/ws-echo-multi"), 3)

		t.Run("send_receive_multiple_wss", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var msg1 = "test1"
			var msg2 = "test2"
			var secondMsgReceived = false
			var res = ws.connect("WSSBIN_URL/ws-echo-multi", (socket) => ***REMOVED***
				socket.on("open", () => ***REMOVED***
					socket.send(msg1)
				***REMOVED***)
				socket.on("message", (data) => ***REMOVED***
					if (data == msg1)***REMOVED***
						socket.send(msg2)
					***REMOVED***
					if (data == msg2)***REMOVED***
						secondMsgReceived = true
						socket.close()
					***REMOVED***
				***REMOVED***);
			***REMOVED***);

			if (!secondMsgReceived) ***REMOVED***
				throw new Error ("second test message was not received from server!");
			***REMOVED***
			`))
			require.NoError(t, err)
		***REMOVED***)

		samplesBuf = metrics.GetBufferedSamples(samples)
		assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSSBIN_URL/ws-echo-multi"), statusProtocolSwitch, "")
		assertMetricEmittedCount(t, metrics.WSMessagesSentName, samplesBuf, sr("WSSBIN_URL/ws-echo-multi"), 2)
		assertMetricEmittedCount(t, metrics.WSMessagesReceivedName, samplesBuf, sr("WSSBIN_URL/ws-echo-multi"), 2)

		t.Run("send_receive_text_binary", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(sr(`
			var msg1 = "test1"
			var msg2 = new Uint8Array([116, 101, 115, 116, 50]); // 'test2'
			var secondMsgReceived = false
			var res = ws.connect("WSBIN_URL/ws-echo-multi", (socket) => ***REMOVED***
				socket.on("open", () => ***REMOVED***
					socket.send(msg1)
				***REMOVED***)
				socket.on("message", (data) => ***REMOVED***
					if (data == msg1)***REMOVED***
						socket.sendBinary(msg2.buffer)
					***REMOVED***
				***REMOVED***);
				socket.on("binaryMessage", (data) => ***REMOVED***
					let data2 = new Uint8Array(data)
					if(JSON.stringify(msg2) == JSON.stringify(data2))***REMOVED***
						secondMsgReceived = true
					***REMOVED***
					socket.close()
				***REMOVED***)
			***REMOVED***);

			if (!secondMsgReceived) ***REMOVED***
				throw new Error ("second test message was not received from server!");
			***REMOVED***
			`))
			require.NoError(t, err)
		***REMOVED***)

		samplesBuf = metrics.GetBufferedSamples(samples)
		assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo-multi"), statusProtocolSwitch, "")
		assertMetricEmittedCount(t, metrics.WSMessagesSentName, samplesBuf, sr("WSBIN_URL/ws-echo-multi"), 2)
		assertMetricEmittedCount(t, metrics.WSMessagesReceivedName, samplesBuf, sr("WSBIN_URL/ws-echo-multi"), 2)
	***REMOVED***)
***REMOVED***

func TestSocketSendBinary(t *testing.T) ***REMOVED*** //nolint:tparallel
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED*** //nolint:exhaustivestruct
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED*** //nolint:exhaustivestruct
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
		***REMOVED***,
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	t.Run("ok", func(t *testing.T) ***REMOVED***
		_, err = rt.RunString(sr(`
		var gotMsg = false;
		var res = ws.connect('WSBIN_URL/ws-echo', function(socket)***REMOVED***
			var data = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'

			socket.on('open', function() ***REMOVED***
				socket.sendBinary(data.buffer);
			***REMOVED***)
			socket.on('binaryMessage', function(msg) ***REMOVED***
				gotMsg = true;
				let decText = String.fromCharCode.apply(null, new Uint8Array(msg));
				decText = decodeURIComponent(escape(decText));
				if (decText !== 'hello') ***REMOVED***
					throw new Error('received unexpected binary message: ' + decText);
				***REMOVED***
				socket.close()
			***REMOVED***);
		***REMOVED***);
		if (!gotMsg) ***REMOVED***
			throw new Error("the 'binaryMessage' handler wasn't called")
		***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)

	errTestCases := []struct ***REMOVED***
		in, expErrType string
	***REMOVED******REMOVED***
		***REMOVED***"", ""***REMOVED***,
		***REMOVED***"undefined", "undefined"***REMOVED***,
		***REMOVED***"null", "null"***REMOVED***,
		***REMOVED***"true", "Boolean"***REMOVED***,
		***REMOVED***"1", "Number"***REMOVED***,
		***REMOVED***"3.14", "Number"***REMOVED***,
		***REMOVED***"'str'", "String"***REMOVED***,
		***REMOVED***"[1, 2, 3]", "Array"***REMOVED***,
		***REMOVED***"new Uint8Array([1, 2, 3])", "Object"***REMOVED***,
		***REMOVED***"Symbol('a')", "Symbol"***REMOVED***,
		***REMOVED***"function() ***REMOVED******REMOVED***", "Function"***REMOVED***,
	***REMOVED***

	for _, tc := range errTestCases ***REMOVED*** //nolint:paralleltest
		tc := tc
		t.Run(fmt.Sprintf("err_%s", tc.expErrType), func(t *testing.T) ***REMOVED***
			_, err = rt.RunString(fmt.Sprintf(sr(`
			var res = ws.connect('WSBIN_URL/ws-echo', function(socket)***REMOVED***
				socket.on('open', function() ***REMOVED***
					socket.sendBinary(%s);
				***REMOVED***)
			***REMOVED***);
		`), tc.in))
			require.Error(t, err)
			if tc.in == "" ***REMOVED***
				assert.Contains(t, err.Error(), "missing argument, expected ArrayBuffer")
			***REMOVED*** else ***REMOVED***
				assert.Contains(t, err.Error(), fmt.Sprintf("expected ArrayBuffer as argument, received: %s", tc.expErrType))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: &metrics.DefaultSystemTagSet,
		***REMOVED***,
		Samples:        samples,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	t.Run("invalid_url", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
		var res = ws.connect("INVALID", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				socket.close();
			***REMOVED***);
		***REMOVED***);
		`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("invalid_url_message_panic", func(t *testing.T) ***REMOVED***
		// Attempting to send a message to a non-existent socket shouldn't panic
		_, err := rt.RunString(`
		var res = ws.connect("INVALID", function(socket)***REMOVED***
			socket.send("new message");
		***REMOVED***);
		`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("error_in_setup", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var res = ws.connect("WSBIN_URL/ws-echo-invalid", function(socket)***REMOVED***
			throw new Error("error in setup");
		***REMOVED***);
		`))
		assert.Error(t, err)
	***REMOVED***)

	t.Run("send_after_close", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var hasError = false;
		var res = ws.connect("WSBIN_URL/ws-echo-invalid", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				socket.close();
				socket.send("test");
			***REMOVED***);

			socket.on("error", function(errorEvent) ***REMOVED***
				hasError = true;
			***REMOVED***);
		***REMOVED***);
		if (!hasError) ***REMOVED***
			throw new Error ("no error emitted for send after close");
		***REMOVED***
		`))
		require.NoError(t, err)
		assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo-invalid"), statusProtocolSwitch, "")
	***REMOVED***)

	t.Run("error on close", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
		var closed = false;
		var res = ws.connect("WSBIN_URL/ws-close", function(socket)***REMOVED***
			socket.on('open', function open() ***REMOVED***
				socket.setInterval(function timeout() ***REMOVED***
				  socket.ping();
				***REMOVED***, 1000);
			***REMOVED***);

			socket.on("ping", function() ***REMOVED***
				socket.close();
			***REMOVED***);

			socket.on("error", function(errorEvent) ***REMOVED***
				if (errorEvent == null) ***REMOVED***
					throw new Error(JSON.stringify(errorEvent));
				***REMOVED***
				if (!closed) ***REMOVED***
					closed = true;
				    socket.close();
				***REMOVED***
			***REMOVED***);
		***REMOVED***);
		`))
		require.NoError(t, err)
		assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-close"), statusProtocolSwitch, "")
	***REMOVED***)
***REMOVED***

func TestSystemTags(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)

	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	// TODO: test for actual tag values after removing the dependency on the
	// external service demos.kaazing.com (https://github.com/k6io/k6/issues/537)
	testedSystemTags := []string***REMOVED***"group", "status", "subproto", "url", "ip"***REMOVED***

	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:          root,
		Dialer:         tb.Dialer,
		Options:        lib.Options***REMOVED***SystemTags: metrics.ToSystemTagSet(testedSystemTags)***REMOVED***,
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	for _, expectedTag := range testedSystemTags ***REMOVED***
		expectedTag := expectedTag
		t.Run("only "+expectedTag, func(t *testing.T) ***REMOVED***
			state.Options.SystemTags = metrics.ToSystemTagSet([]string***REMOVED***expectedTag***REMOVED***)
			_, err := rt.RunString(sr(`
			var res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
				socket.on("open", function() ***REMOVED***
					socket.send("test")
				***REMOVED***)
				socket.on("message", function (data)***REMOVED***
					if (!data=="test") ***REMOVED***
						throw new Error ("echo'd data doesn't match our message!");
					***REMOVED***
					socket.close()
				***REMOVED***);
			***REMOVED***);
			`))
			require.NoError(t, err)

			for _, sampleContainer := range metrics.GetBufferedSamples(samples) ***REMOVED***
				for _, sample := range sampleContainer.GetSamples() ***REMOVED***
					for emittedTag := range sample.Tags.CloneTags() ***REMOVED***
						assert.Equal(t, expectedTag, emittedTag)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestTLSConfig(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	tb := httpmultibin.NewHTTPMultiBin(t)

	sr := tb.Replacer.Replace

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
				metrics.TagIP,
			),
		***REMOVED***,
		Samples:        samples,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	t.Run("insecure skip verify", func(t *testing.T) ***REMOVED***
		state.TLSConfig = &tls.Config***REMOVED***
			InsecureSkipVerify: true,
		***REMOVED***

		_, err := rt.RunString(sr(`
		var res = ws.connect("WSSBIN_URL/ws-close", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-close"), statusProtocolSwitch, "")

	t.Run("custom certificates", func(t *testing.T) ***REMOVED***
		state.TLSConfig = tb.TLSClientConfig

		_, err := rt.RunString(sr(`
			var res = ws.connect("WSSBIN_URL/ws-close", function(socket)***REMOVED***
				socket.close()
			***REMOVED***);
			if (res.status != 101) ***REMOVED***
				throw new Error("TLS connection failed with status: " + res.status);
			***REMOVED***
		`))
		require.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-close"), statusProtocolSwitch, "")
***REMOVED***

func TestReadPump(t *testing.T) ***REMOVED***
	var closeCode int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, r, w.Header())
		require.NoError(t, err)
		closeMsg := websocket.FormatCloseMessage(closeCode, "")
		_ = conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
	***REMOVED***))
	defer srv.Close()

	closeCodes := []int***REMOVED***websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr***REMOVED***

	numAsserts := 0
	srvURL := "ws://" + srv.Listener.Addr().String()

	// Ensure readPump returns the response close code sent by the server
	for _, code := range closeCodes ***REMOVED***
		code := code
		t.Run(strconv.Itoa(code), func(t *testing.T) ***REMOVED***
			closeCode = code
			conn, resp, err := websocket.DefaultDialer.Dial(srvURL, nil)
			require.NoError(t, err)
			defer func() ***REMOVED***
				_ = resp.Body.Close()
				_ = conn.Close()
			***REMOVED***()

			msgChan := make(chan *message)
			errChan := make(chan error)
			closeChan := make(chan int)
			s := &Socket***REMOVED***conn: conn***REMOVED***
			go s.readPump(msgChan, errChan, closeChan)

		readChans:
			for ***REMOVED***
				select ***REMOVED***
				case responseCode := <-closeChan:
					assert.Equal(t, code, responseCode)
					numAsserts++
					break readChans
				case <-errChan:
					continue
				case <-time.After(time.Second):
					t.Errorf("Read timed out")
					break readChans
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// Ensure all close code asserts passed
	assert.Equal(t, numAsserts, len(closeCodes))
***REMOVED***

func TestUserAgent(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	tb.Mux.HandleFunc("/ws-echo-useragent", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
		// Echo back User-Agent header if it exists
		responseHeaders := w.Header().Clone()
		if ua := req.Header.Get("User-Agent"); ua != "" ***REMOVED***
			responseHeaders.Add("Echo-User-Agent", req.Header.Get("User-Agent"))
		***REMOVED***

		conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, req, responseHeaders)
		if err != nil ***REMOVED***
			t.Fatalf("/ws-echo-useragent cannot upgrade request: %v", err)
			return
		***REMOVED***

		err = conn.Close()
		if err != nil ***REMOVED***
			t.Logf("error while closing connection in /ws-echo-useragent: %v", err)
			return
		***REMOVED***
	***REMOVED***))

	root, err := lib.NewGroup("", nil)
	require.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
			UserAgent: null.StringFrom("TestUserAgent"),
		***REMOVED***,
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	m := New().NewModuleInstance(&modulestest.VU***REMOVED***
		CtxField:     context.Background(),
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		RuntimeField: rt,
		StateField:   state,
	***REMOVED***)
	require.NoError(t, rt.Set("ws", m.Exports().Default))

	// websocket handler should echo back User-Agent as Echo-User-Agent for this test to work
	_, err = rt.RunString(sr(`
		var res = ws.connect("WSBIN_URL/ws-echo-useragent", function(socket)***REMOVED***
			socket.close()
		***REMOVED***)
		var userAgent = res.headers["Echo-User-Agent"];
		if (userAgent == undefined) ***REMOVED***
			throw new Error("user agent is not echoed back by test server");
		***REMOVED***
		if (userAgent != "TestUserAgent") ***REMOVED***
			throw new Error("incorrect user agent: " + userAgent);
		***REMOVED***
		`))
	require.NoError(t, err)

	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo-useragent"), statusProtocolSwitch, "")
***REMOVED***

func TestCompression(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("session", func(t *testing.T) ***REMOVED***
		t.Parallel()
		const text string = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed pharetra sapien. Nunc laoreet molestie ante ac gravida. Etiam interdum dui viverra posuere egestas. Pellentesque at dolor tristique, mattis turpis eget, commodo purus. Nunc orci aliquam.`

		ts := newTestState(t)
		sr := ts.tb.Replacer.Replace
		ts.tb.Mux.HandleFunc("/ws-compression", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			upgrader := websocket.Upgrader***REMOVED***
				EnableCompression: true,
				ReadBufferSize:    1024,
				WriteBufferSize:   1024,
			***REMOVED***

			conn, e := upgrader.Upgrade(w, req, w.Header())
			if e != nil ***REMOVED***
				t.Fatalf("/ws-compression cannot upgrade request: %v", e)
				return
			***REMOVED***

			// send a message and exit
			if e = conn.WriteMessage(websocket.TextMessage, []byte(text)); e != nil ***REMOVED***
				t.Logf("error while sending message in /ws-compression: %v", e)
				return
			***REMOVED***

			e = conn.Close()
			if e != nil ***REMOVED***
				t.Logf("error while closing connection in /ws-compression: %v", e)
				return
			***REMOVED***
		***REMOVED***))

		_, err := ts.rt.RunString(sr(`
		// if client supports compression, it has to send the header
		// 'Sec-Websocket-Extensions:permessage-deflate; server_no_context_takeover; client_no_context_takeover' to server.
		// if compression is negotiated successfully, server will reply with header
		// 'Sec-Websocket-Extensions:permessage-deflate; server_no_context_takeover; client_no_context_takeover'

		var params = ***REMOVED***
			"compression": "deflate"
		***REMOVED***
		var res = ws.connect("WSBIN_URL/ws-compression", params, function(socket)***REMOVED***
			socket.on('message', (data) => ***REMOVED***
				if(data != "` + text + `")***REMOVED***
					throw new Error("wrong message received from server: ", data)
				***REMOVED***
				socket.close()
			***REMOVED***)
		***REMOVED***);

		var wsExtensions = res.headers["Sec-Websocket-Extensions"].split(';').map(e => e.trim())
		if (!(wsExtensions.includes("permessage-deflate") && wsExtensions.includes("server_no_context_takeover") && wsExtensions.includes("client_no_context_takeover")))***REMOVED***
			throw new Error("websocket compression negotiation failed");
		***REMOVED***
		`))

		require.NoError(t, err)
		assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(ts.samples), "", sr("WSBIN_URL/ws-compression"), statusProtocolSwitch, "")
	***REMOVED***)

	t.Run("params", func(t *testing.T) ***REMOVED***
		t.Parallel()
		testCases := []struct ***REMOVED***
			compression   string
			expectedError string
		***REMOVED******REMOVED***
			***REMOVED***compression: ""***REMOVED***,
			***REMOVED***compression: "  "***REMOVED***,
			***REMOVED***compression: "deflate"***REMOVED***,
			***REMOVED***compression: "deflate "***REMOVED***,
			***REMOVED***
				compression:   "gzip",
				expectedError: `unsupported compression algorithm 'gzip', supported algorithm is 'deflate'`,
			***REMOVED***,
			***REMOVED***
				compression:   "deflate, gzip",
				expectedError: `unsupported compression algorithm 'deflate, gzip', supported algorithm is 'deflate'`,
			***REMOVED***,
			***REMOVED***
				compression:   "deflate, deflate",
				expectedError: `unsupported compression algorithm 'deflate, deflate', supported algorithm is 'deflate'`,
			***REMOVED***,
			***REMOVED***
				compression:   "deflate, ",
				expectedError: `unsupported compression algorithm 'deflate,', supported algorithm is 'deflate'`,
			***REMOVED***,
		***REMOVED***

		for _, testCase := range testCases ***REMOVED***
			testCase := testCase
			t.Run(testCase.compression, func(t *testing.T) ***REMOVED***
				t.Parallel()
				ts := newTestState(t)
				sr := ts.tb.Replacer.Replace
				ts.tb.Mux.HandleFunc("/ws-compression-param", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
					upgrader := websocket.Upgrader***REMOVED***
						EnableCompression: true,
						ReadBufferSize:    1024,
						WriteBufferSize:   1024,
					***REMOVED***

					conn, e := upgrader.Upgrade(w, req, w.Header())
					if e != nil ***REMOVED***
						t.Fatalf("/ws-compression-param cannot upgrade request: %v", e)
						return
					***REMOVED***

					e = conn.Close()
					if e != nil ***REMOVED***
						t.Logf("error while closing connection in /ws-compression-param: %v", e)
						return
					***REMOVED***
				***REMOVED***))

				_, err := ts.rt.RunString(sr(`
					var res = ws.connect("WSBIN_URL/ws-compression-param", ***REMOVED***"compression":"` + testCase.compression + `"***REMOVED***, function(socket)***REMOVED***
						socket.close()
					***REMOVED***);
				`))

				if testCase.expectedError == "" ***REMOVED***
					require.NoError(t, err)
				***REMOVED*** else ***REMOVED***
					require.Error(t, err)
					require.Contains(t, err.Error(), testCase.expectedError)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func clearSamples(tb *httpmultibin.HTTPMultiBin, samples chan metrics.SampleContainer) ***REMOVED***
	ctxDone := tb.Context.Done()
	for ***REMOVED***
		select ***REMOVED***
		case <-samples:
		case <-ctxDone:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCompression(b *testing.B) ***REMOVED***
	const textMessage = 1
	ts := newTestState(b)
	sr := ts.tb.Replacer.Replace
	go clearSamples(ts.tb, ts.samples)

	testCodes := []string***REMOVED***
		sr(`
		var res = ws.connect("WSBIN_URL/ws-compression", ***REMOVED***"compression":"deflate"***REMOVED***, (socket) => ***REMOVED***
			socket.on('message', (data) => ***REMOVED***
				socket.close()
			***REMOVED***)
		***REMOVED***);
		`),
		sr(`
		var res = ws.connect("WSBIN_URL/ws-compression", ***REMOVED******REMOVED***, (socket) => ***REMOVED***
			socket.on('message', (data) => ***REMOVED***
				socket.close()
			***REMOVED***)
		***REMOVED***);
		`),
	***REMOVED***

	ts.tb.Mux.HandleFunc("/ws-compression", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
		kbData := bytes.Repeat([]byte("0123456789"), 100)

		// upgrade connection, send the first (long) message, disconnect
		upgrader := websocket.Upgrader***REMOVED***
			EnableCompression: true,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
		***REMOVED***

		conn, e := upgrader.Upgrade(w, req, w.Header())

		if e != nil ***REMOVED***
			b.Fatalf("/ws-compression cannot upgrade request: %v", e)
			return
		***REMOVED***

		if e = conn.WriteMessage(textMessage, kbData); e != nil ***REMOVED***
			b.Fatalf("/ws-compression cannot write message: %v", e)
			return
		***REMOVED***

		e = conn.Close()
		if e != nil ***REMOVED***
			b.Logf("error while closing connection in /ws-compression: %v", e)
			return
		***REMOVED***
	***REMOVED***))

	b.ResetTimer()
	b.Run("compression-enabled", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			if _, err := ts.rt.RunString(testCodes[0]); err != nil ***REMOVED***
				b.Error(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	b.Run("compression-disabled", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			if _, err := ts.rt.RunString(testCodes[1]); err != nil ***REMOVED***
				b.Error(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestCookieJar(t *testing.T) ***REMOVED***
	t.Parallel()
	ts := newTestState(t)
	sr := ts.tb.Replacer.Replace

	ts.tb.Mux.HandleFunc("/ws-echo-someheader", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
		responseHeaders := w.Header().Clone()
		if sh, err := req.Cookie("someheader"); err == nil ***REMOVED***
			responseHeaders.Add("Echo-Someheader", sh.Value)
		***REMOVED***

		conn, err := (&websocket.Upgrader***REMOVED******REMOVED***).Upgrade(w, req, responseHeaders)
		if err != nil ***REMOVED***
			t.Fatalf("/ws-echo-someheader cannot upgrade request: %v", err)
		***REMOVED***

		err = conn.Close()
		if err != nil ***REMOVED***
			t.Logf("error while closing connection in /ws-echo-someheader: %v", err)
		***REMOVED***
	***REMOVED***))

	mii := &modulestest.VU***REMOVED***
		RuntimeField: ts.rt,
		InitEnvField: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
		CtxField:     context.Background(),
		StateField:   ts.state,
	***REMOVED***
	err := ts.rt.Set("http", httpModule.New().NewModuleInstance(mii).Exports().Default)
	require.NoError(t, err)
	ts.state.CookieJar, _ = cookiejar.New(nil)

	_, err = ts.rt.RunString(sr(`
		var res = ws.connect("WSBIN_URL/ws-echo-someheader", function(socket)***REMOVED***
			socket.close()
		***REMOVED***)
		var someheader = res.headers["Echo-Someheader"];
		if (someheader !== undefined) ***REMOVED***
			throw new Error("someheader is echoed back by test server even though it doesn't exist");
		***REMOVED***

		http.cookieJar().set("HTTPBIN_URL/ws-echo-someheader", "someheader", "defaultjar")
		res = ws.connect("WSBIN_URL/ws-echo-someheader", function(socket)***REMOVED***
			socket.close()
		***REMOVED***)
		someheader = res.headers["Echo-Someheader"];
		if (someheader != "defaultjar") ***REMOVED***
			throw new Error("someheader has wrong value "+ someheader + " instead of defaultjar");
		***REMOVED***

		var jar = new http.CookieJar();
		jar.set("HTTPBIN_URL/ws-echo-someheader", "someheader", "customjar")
		res = ws.connect("WSBIN_URL/ws-echo-someheader", ***REMOVED***jar: jar***REMOVED***, function(socket)***REMOVED***
			socket.close()
		***REMOVED***)
		someheader = res.headers["Echo-Someheader"];
		if (someheader != "customjar") ***REMOVED***
			throw new Error("someheader has wrong value "+ someheader + " instead of customjar");
		***REMOVED***
		`))
	require.NoError(t, err)

	assertSessionMetricsEmitted(t, metrics.GetBufferedSamples(ts.samples), "", sr("WSBIN_URL/ws-echo-someheader"), statusProtocolSwitch, "")
***REMOVED***
