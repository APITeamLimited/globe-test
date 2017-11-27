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
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

func assertSessionMetricsEmitted(t *testing.T, samples []stats.Sample, subprotocol, url string, status int, group string) ***REMOVED***
	seenSessions := false
	seenSessionDuration := false
	seenConnecting := false

	for _, sample := range samples ***REMOVED***
		if sample.Tags["url"] == url ***REMOVED***
			switch sample.Metric ***REMOVED***
			case metrics.WSConnecting:
				seenConnecting = true
			case metrics.WSSessionDuration:
				seenSessionDuration = true
			case metrics.WSSessions:
				seenSessions = true
			***REMOVED***

			assert.Equal(t, strconv.Itoa(status), sample.Tags["status"])
			assert.Equal(t, subprotocol, sample.Tags["subprotocol"])
			assert.Equal(t, group, sample.Tags["group"])
		***REMOVED***
	***REMOVED***
	assert.True(t, seenConnecting, "url %s didn't emit Connecting", url)
	assert.True(t, seenSessions, "url %s didn't emit Sessions", url)
	assert.True(t, seenSessionDuration, "url %s didn't emit SessionDuration", url)
***REMOVED***

func assertMetricEmitted(t *testing.T, metric *stats.Metric, samples []stats.Sample, url string) ***REMOVED***
	seenMetric := false

	for _, sample := range samples ***REMOVED***
		if sample.Tags["url"] == url ***REMOVED***
			if sample.Metric == metric ***REMOVED***
				seenMetric = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.True(t, seenMetric, "url %s didn't emit %s", url, metric.Name)
***REMOVED***
func TestSession(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	dialer := netext.NewDialer(net.Dialer***REMOVED***
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	***REMOVED***)
	state := &common.State***REMOVED***Group: root, Dialer: dialer***REMOVED***

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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("connect_wss", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		let res = ws.connect("wss://demos.kaazing.com/echo", function(socket)***REMOVED***
			socket.close()
		***REMOVED***);
		if (res.status != 101) ***REMOVED*** throw new Error("TLS connection failed with status: " + res.status); ***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)
	assertSessionMetricsEmitted(t, state.Samples, "", "wss://demos.kaazing.com/echo", 101, "")

	t.Run("open", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("send_receive", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSMessagesSent, state.Samples, "ws://demos.kaazing.com/echo")
	assertMetricEmitted(t, metrics.WSMessagesReceived, state.Samples, "ws://demos.kaazing.com/echo")

	t.Run("interval", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("timeout", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")

	t.Run("ping", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSPing, state.Samples, "ws://demos.kaazing.com/echo")

	t.Run("multiple_handlers", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")
	assertMetricEmitted(t, metrics.WSPing, state.Samples, "ws://demos.kaazing.com/echo")

	t.Run("close", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
	assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")
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
	state := &common.State***REMOVED***Group: root, Dialer: dialer***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)

	rt.Set("ws", common.Bind(rt, New(), &ctx))

	t.Run("invalid_url", func(t *testing.T) ***REMOVED***
		state.Samples = nil
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
		state.Samples = nil
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
		assertSessionMetricsEmitted(t, state.Samples, "", "ws://demos.kaazing.com/echo", 101, "")
	***REMOVED***)
***REMOVED***
