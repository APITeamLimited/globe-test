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
	"context"
	"crypto/tls"
	"strconv"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

func assertSessionMetricsEmitted(t *testing.T, sampleContainers []stats.SampleContainer, subprotocol, url string, status int, group string) ***REMOVED***
	seenSessions := false
	seenSessionDuration := false
	seenConnecting := false

	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			tags := sample.Tags.CloneTags()
			if tags["url"] == url ***REMOVED***
				switch sample.Metric ***REMOVED***
				case metrics.WSConnecting:
					seenConnecting = true
				case metrics.WSSessionDuration:
					seenSessionDuration = true
				case metrics.WSSessions:
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

func assertMetricEmitted(t *testing.T, metric *stats.Metric, sampleContainers []stats.SampleContainer, url string) ***REMOVED***
	seenMetric := false

	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			surl, ok := sample.Tags.Get("url")
			assert.True(t, ok)
			if surl == url ***REMOVED***
				if sample.Metric == metric ***REMOVED***
					seenMetric = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.True(t, seenMetric, "url %s didn't emit %s", url, metric.Name)
***REMOVED***

func TestSession(t *testing.T) ***REMOVED***
	//TODO: split and paralelize tests
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: stats.ToSystemTagSet([]string***REMOVED***
				stats.TagURL.String(),
				stats.TagProto.String(),
				stats.TagStatus.String(),
				stats.TagSubProto.String(),
			***REMOVED***),
		***REMOVED***,
		Samples:   samples,
		TLSConfig: tb.TLSClientConfig,
	***REMOVED***

	ctx := context.Background()
	ctx = lib.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	t.Run("connect_ws", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("connection failed with status: " + res.status); ***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")

	t.Run("connect_wss", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = ws.connect("WSSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-echo"), 101, "")

	t.Run("open", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let opened = false;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				opened = true;
				socket.close()
			***REMOVED***)
		***REMOVED***);
		if (!opened) ***REMOVED*** throw new Error ("open event not fired"); ***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")

	t.Run("send_receive", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
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
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf := stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), 101, "")
	assertMetricEmitted(t, metrics.WSMessagesSent, samplesBuf, sr("WSBIN_URL/ws-echo"))
	assertMetricEmitted(t, metrics.WSMessagesReceived, samplesBuf, sr("WSBIN_URL/ws-echo"))

	t.Run("interval", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let counter = 0;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setInterval(function () ***REMOVED***
				counter += 1;
				if (counter > 2) ***REMOVED*** socket.close(); ***REMOVED***
			***REMOVED***, 100);
		***REMOVED***);
		if (counter < 3) ***REMOVED***throw new Error ("setInterval should have been called at least 3 times, counter=" + counter);***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")

	t.Run("timeout", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let start = new Date().getTime();
		let ellapsed = new Date().getTime() - start;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.setTimeout(function () ***REMOVED***
				ellapsed = new Date().getTime() - start;
				socket.close();
			***REMOVED***, 500);
		***REMOVED***);
		if (ellapsed > 3000 || ellapsed < 500) ***REMOVED***
			throw new Error ("setTimeout occurred after " + ellapsed + "ms, expected 500<T<3000");
		***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")

	t.Run("ping", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let pongReceived = false;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
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
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf = stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), 101, "")
	assertMetricEmitted(t, metrics.WSPing, samplesBuf, sr("WSBIN_URL/ws-echo"))

	t.Run("multiple_handlers", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let pongReceived = false;
		let otherPongReceived = false;

		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
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
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf = stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", sr("WSBIN_URL/ws-echo"), 101, "")
	assertMetricEmitted(t, metrics.WSPing, samplesBuf, sr("WSBIN_URL/ws-echo"))

	t.Run("close", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let closed = false;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
							socket.close()
			***REMOVED***)
			socket.on("close", function() ***REMOVED***
							closed = true;
			***REMOVED***)
		***REMOVED***);
		if (!closed) ***REMOVED*** throw new Error ("close event not fired"); ***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")
***REMOVED***

func TestErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: stats.ToSystemTagSet(stats.DefaultSystemTagList),
		***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = lib.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	t.Run("invalid_url", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let res = ws.connect("INVALID", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				socket.close();
			***REMOVED***);
		***REMOVED***);
		`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("invalid_url_message_panic", func(t *testing.T) ***REMOVED***
		// Attempting to send a message to a non-existent socket shouldn't panic
		_, err := common.RunString(rt, `
		let res = ws.connect("INVALID", function(socket)***REMOVED***
			socket.send("new message");
		***REMOVED***);
		`)
		assert.Error(t, err)
	***REMOVED***)

	t.Run("error_in_setup", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
			throw new Error("error in setup");
		***REMOVED***);
		`))
		assert.Error(t, err)
	***REMOVED***)

	t.Run("send_after_close", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		let hasError = false;
		let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
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
		assert.NoError(t, err)
		assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-echo"), 101, "")
	***REMOVED***)

	t.Run("error on close", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
		var closed = false;
		let res = ws.connect("WSBIN_URL/ws-close", function(socket)***REMOVED***
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
		assert.NoError(t, err)
		assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSBIN_URL/ws-close"), 101, "")
	***REMOVED***)
***REMOVED***

func TestSystemTags(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	//TODO: test for actual tag values after removing the dependency on the
	// external service demos.kaazing.com (https://github.com/loadimpact/k6/issues/537)
	testedSystemTags := []string***REMOVED***"group", "status", "subproto", "url", "ip"***REMOVED***

	samples := make(chan stats.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:     root,
		Dialer:    tb.Dialer,
		Options:   lib.Options***REMOVED***SystemTags: stats.ToSystemTagSet(testedSystemTags)***REMOVED***,
		Samples:   samples,
		TLSConfig: tb.TLSClientConfig,
	***REMOVED***

	ctx := context.Background()
	ctx = lib.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	for _, expectedTag := range testedSystemTags ***REMOVED***
		expectedTag := expectedTag
		t.Run("only "+expectedTag, func(t *testing.T) ***REMOVED***
			state.Options.SystemTags = stats.ToSystemTagSet([]string***REMOVED***expectedTag***REMOVED***)
			_, err := common.RunString(rt, sr(`
			let res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
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
			assert.NoError(t, err)

			for _, sampleContainer := range stats.GetBufferedSamples(samples) ***REMOVED***
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
	assert.NoError(t, err)

	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	sr := tb.Replacer.Replace

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group:  root,
		Dialer: tb.Dialer,
		Options: lib.Options***REMOVED***
			SystemTags: stats.ToSystemTagSet([]string***REMOVED***
				stats.TagURL.String(),
				stats.TagProto.String(),
				stats.TagStatus.String(),
				stats.TagSubProto.String(),
				stats.TagIP.String(),
			***REMOVED***),
		***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = lib.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	t.Run("insecure skip verify", func(t *testing.T) ***REMOVED***
		state.TLSConfig = &tls.Config***REMOVED***
			InsecureSkipVerify: true,
		***REMOVED***

		_, err := common.RunString(rt, sr(`
		let res = ws.connect("WSSBIN_URL/ws-close", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-close"), 101, "")

	t.Run("custom certificates", func(t *testing.T) ***REMOVED***
		state.TLSConfig = tb.TLSClientConfig

		_, err := common.RunString(rt, sr(`
			let res = ws.connect("WSSBIN_URL/ws-close", function(socket)***REMOVED***
				socket.close()
			***REMOVED***);
			if (res.status != 101) ***REMOVED***
				throw new Error("TLS connection failed with status: " + res.status);
			***REMOVED***
		`))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", sr("WSSBIN_URL/ws-close"), 101, "")
***REMOVED***
