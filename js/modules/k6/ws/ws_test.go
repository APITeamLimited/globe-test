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
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/testutils"
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

func makeWsProto(s string) string ***REMOVED***
	return "ws" + strings.TrimPrefix(s, "http")
***REMOVED***

func TestSession(t *testing.T) ***REMOVED***
	//TODO: split and paralelize tests

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	***REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &common.State***REMOVED***
		Group:  root,
		Dialer: dialer,
		Options: lib.Options***REMOVED***
			SystemTags: lib.GetTagSet("url", "proto", "status", "subproto"),
		***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	t.Run("connect_ws", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("connection failed with status: " + res.status); ***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("connect_wss", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let res = ws.connect("wss://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "wss://demos.kaazing.com/echo", 101, "")

	t.Run("open", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let opened = false;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
				opened = true;
				socket.close()
			***REMOVED***)
		***REMOVED***);
		if (!opened) ***REMOVED*** throw new Error ("open event not fired"); ***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("send_receive", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
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
		`)
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf := stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSMessagesSent, samplesBuf, "ws://demos.kaazing.com/echo")
	assertMetricEmitted(t, metrics.WSMessagesReceived, samplesBuf, "ws://demos.kaazing.com/echo")

	t.Run("interval", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let counter = 0;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.setInterval(function () ***REMOVED***
				counter += 1;
				if (counter > 2) ***REMOVED*** socket.close(); ***REMOVED***
			***REMOVED***, 100);
		***REMOVED***);
		if (counter < 3) ***REMOVED***throw new Error ("setInterval should have been called at least 3 times, counter=" + counter);***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("timeout", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let start = new Date().getTime();
		let ellapsed = new Date().getTime() - start;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.setTimeout(function () ***REMOVED***
				ellapsed = new Date().getTime() - start;
				socket.close();
			***REMOVED***, 500);
		***REMOVED***);
		if (ellapsed > 2000 || ellapsed < 500) ***REMOVED***
			throw new Error ("setTimeout occurred after " + ellapsed + "ms, expected 500<T<2000");
		***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("ping", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let pongReceived = false;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
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
		`)
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf = stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSPing, samplesBuf, "ws://demos.kaazing.com/echo")

	t.Run("multiple_handlers", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let pongReceived = false;
		let otherPongReceived = false;

		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
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
		`)
		assert.NoError(t, err)
	***REMOVED***)

	samplesBuf = stats.GetBufferedSamples(samples)
	assertSessionMetricsEmitted(t, samplesBuf, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSPing, samplesBuf, "ws://demos.kaazing.com/echo")

	t.Run("close", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let closed = false;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.on("open", function() ***REMOVED***
							socket.close()
			***REMOVED***)
			socket.on("close", function() ***REMOVED***
							closed = true;
			***REMOVED***)
		***REMOVED***);
		if (!closed) ***REMOVED*** throw new Error ("close event not fired"); ***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")
***REMOVED***

func TestErrors(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	***REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &common.State***REMOVED***
		Group:  root,
		Dialer: dialer,
		Options: lib.Options***REMOVED***
			SystemTags: lib.GetTagSet(lib.DefaultSystemTagList...),
		***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
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

	t.Run("send_after_close", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let hasError = false;
		let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
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
		`)
		assert.NoError(t, err)
		assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", "ws://demos.kaazing.com/echo", 101, "")
	***REMOVED***)
***REMOVED***

func TestSystemTags(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	***REMOVED***)

	//TODO: test for actual tag values after removing the dependency on the
	// external service demos.kaazing.com (https://github.com/loadimpact/k6/issues/537)
	testedSystemTags := []string***REMOVED***"group", "status", "subproto", "url", "ip"***REMOVED***

	samples := make(chan stats.SampleContainer, 1000)
	state := &common.State***REMOVED***
		Group:   root,
		Dialer:  dialer,
		Options: lib.Options***REMOVED***SystemTags: lib.GetTagSet(testedSystemTags...)***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	for _, expectedTag := range testedSystemTags ***REMOVED***
		t.Run("only "+expectedTag, func(t *testing.T) ***REMOVED***
			state.Options.SystemTags = map[string]bool***REMOVED***
				expectedTag: true,
			***REMOVED***
			_, err := common.RunString(rt, `
			let res = ws.connect("ws://demos.kaazing.com/echo", function(socket)***REMOVED***
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
			`)
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

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	***REMOVED***)
	samples := make(chan stats.SampleContainer, 1000)
	state := &common.State***REMOVED***
		Group:  root,
		Dialer: dialer,
		Options: lib.Options***REMOVED***
			SystemTags: lib.GetTagSet("url", "proto", "status", "subproto", "ip"),
		***REMOVED***,
		Samples: samples,
	***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	tb := testutils.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	url := makeWsProto(tb.ServerHTTPS.URL) + "/ws-close"

	t.Run("insecure skip verify", func(t *testing.T) ***REMOVED***
		state.TLSConfig = &tls.Config***REMOVED***
			InsecureSkipVerify: true,
		***REMOVED***

		_, err := common.RunString(rt, fmt.Sprintf(`
		let res = ws.connect("%s", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`, url))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", url, 101, "")

	t.Run("custom certificates", func(t *testing.T) ***REMOVED***
		state.TLSConfig = tb.TLSClientConfig

		_, err := common.RunString(rt, fmt.Sprintf(`
		let res = ws.connect("%s", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`, url))
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, stats.GetBufferedSamples(samples), "", url, 101, "")
***REMOVED***
