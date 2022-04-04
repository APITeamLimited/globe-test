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

package httpext

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mccutchen/go-httpbin/httpbin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/metrics"
)

const traceDelay = 100 * time.Millisecond

func getTestTracer(t *testing.T) (*Tracer, *httptrace.ClientTrace) ***REMOVED***
	tracer := &Tracer***REMOVED******REMOVED***
	ct := tracer.Trace()
	if runtime.GOOS == "windows" ***REMOVED***
		// HACK: Time resolution is not as accurate on Windows, see:
		//  https://github.com/golang/go/issues/8687
		//  https://github.com/golang/go/issues/41087
		// Which seems to be causing some metrics to have a value of 0,
		// since e.g. ConnectStart and ConnectDone could register the same time.
		// So we force delays in the ClientTrace event handlers
		// to hopefully reduce the chances of this happening.
		ct = &httptrace.ClientTrace***REMOVED***
			ConnectStart: func(a, n string) ***REMOVED***
				t.Logf("called ConnectStart at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.ConnectStart(a, n)
			***REMOVED***,
			ConnectDone: func(a, n string, e error) ***REMOVED***
				t.Logf("called ConnectDone at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.ConnectDone(a, n, e)
			***REMOVED***,
			GetConn: func(h string) ***REMOVED***
				t.Logf("called GetConn at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.GetConn(h)
			***REMOVED***,
			GotConn: func(i httptrace.GotConnInfo) ***REMOVED***
				t.Logf("called GotConn at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.GotConn(i)
			***REMOVED***,
			TLSHandshakeStart: func() ***REMOVED***
				t.Logf("called TLSHandshakeStart at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.TLSHandshakeStart()
			***REMOVED***,
			TLSHandshakeDone: func(s tls.ConnectionState, e error) ***REMOVED***
				t.Logf("called TLSHandshakeDone at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.TLSHandshakeDone(s, e)
			***REMOVED***,
			WroteRequest: func(i httptrace.WroteRequestInfo) ***REMOVED***
				t.Logf("called WroteRequest at\t\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.WroteRequest(i)
			***REMOVED***,
			GotFirstResponseByte: func() ***REMOVED***
				t.Logf("called GotFirstResponseByte at\t%v\n", now())
				time.Sleep(traceDelay)
				tracer.GotFirstResponseByte()
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return tracer, ct
***REMOVED***

func TestTracer(t *testing.T) ***REMOVED*** //nolint:tparallel
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	transport, ok := srv.Client().Transport.(*http.Transport)
	assert.True(t, ok)
	transport.DialContext = netext.NewDialer(
		net.Dialer***REMOVED******REMOVED***,
		netext.NewResolver(net.LookupIP, 0, types.DNSfirst, types.DNSpreferIPv4),
	).DialContext

	var prev int64
	assertLaterOrZero := func(t *testing.T, val int64, canBeZero bool) ***REMOVED***
		if canBeZero && val == 0 ***REMOVED***
			return
		***REMOVED***
		if prev > val ***REMOVED***
			_, file, line, _ := runtime.Caller(1)
			t.Errorf("Expected %d to be greater or equal to %d (from %s:%d)", val, prev, file, line)
			return
		***REMOVED***
		prev = val
	***REMOVED***
	builtinMetrics := metrics.RegisterBuiltinMetrics(metrics.NewRegistry())

	for tnum, isReuse := range []bool***REMOVED***false, true, true***REMOVED*** ***REMOVED*** //nolint:paralleltest
		t.Run(fmt.Sprintf("Test #%d", tnum), func(t *testing.T) ***REMOVED***
			// Do not enable parallel testing, test relies on sequential execution
			req, err := http.NewRequest("GET", srv.URL+"/get", nil)
			require.NoError(t, err)

			tracer, ct := getTestTracer(t)
			res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), ct)))
			require.NoError(t, err)

			_, err = io.Copy(ioutil.Discard, res.Body)
			assert.NoError(t, err)
			assert.NoError(t, res.Body.Close())
			if runtime.GOOS == "windows" ***REMOVED***
				time.Sleep(traceDelay)
			***REMOVED***
			trail := tracer.Done()
			trail.SaveSamples(builtinMetrics, metrics.IntoSampleTags(&map[string]string***REMOVED***"tag": "value"***REMOVED***))
			samples := trail.GetSamples()

			assertLaterOrZero(t, tracer.getConn, isReuse)
			assertLaterOrZero(t, tracer.connectStart, isReuse)
			assertLaterOrZero(t, tracer.connectDone, isReuse)
			assertLaterOrZero(t, tracer.tlsHandshakeStart, isReuse)
			assertLaterOrZero(t, tracer.tlsHandshakeDone, isReuse)
			assertLaterOrZero(t, tracer.gotConn, false)
			assertLaterOrZero(t, tracer.wroteRequest, false)
			assertLaterOrZero(t, tracer.gotFirstResponseByte, false)
			assertLaterOrZero(t, now(), false)

			assert.Equal(t, strings.TrimPrefix(srv.URL, "https://"), trail.ConnRemoteAddr.String())

			assert.Len(t, samples, 8)
			seenMetrics := map[*metrics.Metric]bool***REMOVED******REMOVED***
			for i, s := range samples ***REMOVED***
				assert.NotContains(t, seenMetrics, s.Metric)
				seenMetrics[s.Metric] = true

				assert.False(t, s.Time.IsZero())
				assert.Equal(t, map[string]string***REMOVED***"tag": "value"***REMOVED***, s.Tags.CloneTags())

				switch s.Metric ***REMOVED***
				case builtinMetrics.HTTPReqs:
					assert.Equal(t, 1.0, s.Value)
					assert.Equal(t, 0, i, "`HTTPReqs` is reported before the other HTTP builtinMetrics")
				case builtinMetrics.HTTPReqConnecting, builtinMetrics.HTTPReqTLSHandshaking:
					if isReuse ***REMOVED***
						assert.Equal(t, 0.0, s.Value)
						break
					***REMOVED***
					fallthrough
				case builtinMetrics.HTTPReqDuration, builtinMetrics.HTTPReqBlocked, builtinMetrics.HTTPReqSending, builtinMetrics.HTTPReqWaiting, builtinMetrics.HTTPReqReceiving:
					assert.True(t, s.Value > 0.0, "%s is <= 0", s.Metric.Name)
				default:
					t.Errorf("unexpected metric: %s", s.Metric.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

type failingConn struct ***REMOVED***
	net.Conn
***REMOVED***

var failOnConnWrite = false

func (c failingConn) Write(b []byte) (int, error) ***REMOVED***
	if failOnConnWrite ***REMOVED***
		failOnConnWrite = false
		return 0, errors.New("write error")
	***REMOVED***

	return c.Conn.Write(b)
***REMOVED***

func TestTracerNegativeHttpSendingValues(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	transport, ok := srv.Client().Transport.(*http.Transport)
	assert.True(t, ok)

	dialer := &net.Dialer***REMOVED******REMOVED***
	transport.DialContext = func(ctx context.Context, proto, addr string) (net.Conn, error) ***REMOVED***
		conn, err := dialer.DialContext(ctx, proto, addr)
		return failingConn***REMOVED***conn***REMOVED***, err
	***REMOVED***

	req, err := http.NewRequest("GET", srv.URL+"/get", nil)
	require.NoError(t, err)

	***REMOVED***
		tracer := &Tracer***REMOVED******REMOVED***
		res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), tracer.Trace())))
		require.NoError(t, err)
		_, err = io.Copy(ioutil.Discard, res.Body)
		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		tracer.Done()
	***REMOVED***

	// make the next connection write fail
	failOnConnWrite = true

	***REMOVED***
		tracer := &Tracer***REMOVED******REMOVED***
		res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), tracer.Trace())))
		require.NoError(t, err)
		_, err = io.Copy(ioutil.Discard, res.Body)
		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		trail := tracer.Done()
		builtinMetrics := metrics.RegisterBuiltinMetrics(metrics.NewRegistry())
		trail.SaveSamples(builtinMetrics, nil)

		require.True(t, trail.Sending > 0)
	***REMOVED***
***REMOVED***

func TestTracerError(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	tracer := &Tracer***REMOVED******REMOVED***
	req, err := http.NewRequest("GET", srv.URL+"/get", nil)
	require.NoError(t, err)

	_, err = http.DefaultTransport.RoundTrip(
		req.WithContext(
			httptrace.WithClientTrace(
				context.Background(),
				tracer.Trace())))

	assert.Error(t, err)
***REMOVED***

func TestCancelledRequest(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	t.Cleanup(srv.Close)

	cancelTest := func(t *testing.T) ***REMOVED***
		tracer := &Tracer***REMOVED******REMOVED***
		req, err := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/delay/1", nil)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(httptrace.WithClientTrace(req.Context(), tracer.Trace()))
		req = req.WithContext(ctx)
		go func() ***REMOVED***
			time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond) //nolint:gosec
			cancel()
		***REMOVED***()

		resp, err := srv.Client().Transport.RoundTrip(req) //nolint:bodyclose
		_ = tracer.Done()
		if resp == nil && err == nil ***REMOVED***
			t.Errorf("Expected either a RoundTrip response or error but got %#v and %#v", resp, err)
		***REMOVED***
	***REMOVED***

	// This Run will not return until the parallel subtests complete.
	t.Run("group", func(t *testing.T) ***REMOVED***
		t.Parallel()
		for i := 0; i < 200; i++ ***REMOVED***
			t.Run(fmt.Sprintf("TestCancelledRequest_%d", i),
				func(t *testing.T) ***REMOVED***
					t.Parallel()
					cancelTest(t)
				***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
