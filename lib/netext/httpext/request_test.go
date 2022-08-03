/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/mccutchen/go-httpbin/httpbin"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
)

type reader func([]byte) (int, error)

func (r reader) Read(a []byte) (int, error) ***REMOVED***
	return ((func([]byte) (int, error))(r))(a)
***REMOVED***

const (
	badReadMsg  = "bad read error for test"
	badCloseMsg = "bad close error for test"
)

func badReadBody() io.Reader ***REMOVED***
	return reader(func(_ []byte) (int, error) ***REMOVED***
		return 0, errors.New(badReadMsg)
	***REMOVED***)
***REMOVED***

type closer func() error

func (c closer) Close() error ***REMOVED***
	return ((func() error)(c))()
***REMOVED***

func badCloseBody() io.ReadCloser ***REMOVED***
	return struct ***REMOVED***
		io.Reader
		io.Closer
	***REMOVED******REMOVED***
		Reader: reader(func(_ []byte) (int, error) ***REMOVED***
			return 0, io.EOF
		***REMOVED***),
		Closer: closer(func() error ***REMOVED***
			return errors.New(badCloseMsg)
		***REMOVED***),
	***REMOVED***
***REMOVED***

func TestCompressionBodyError(t *testing.T) ***REMOVED***
	t.Parallel()
	algos := []CompressionType***REMOVED***CompressionTypeGzip***REMOVED***
	t.Run("bad read body", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, _, err := compressBody(algos, ioutil.NopCloser(badReadBody()))
		require.Error(t, err)
		require.Equal(t, err.Error(), badReadMsg)
	***REMOVED***)

	t.Run("bad close body", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, _, err := compressBody(algos, badCloseBody())
		require.Error(t, err)
		require.Equal(t, err.Error(), badCloseMsg)
	***REMOVED***)
***REMOVED***

func TestMakeRequestError(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("bad compression algorithm body", func(t *testing.T) ***REMOVED***
		t.Parallel()
		req, err := http.NewRequest("GET", "https://wont.be.used", nil)

		require.NoError(t, err)
		badCompressionType := CompressionType(13)
		require.False(t, badCompressionType.IsACompressionType())
		preq := &ParsedHTTPRequest***REMOVED***
			Req:          req,
			Body:         new(bytes.Buffer),
			Compressions: []CompressionType***REMOVED***badCompressionType***REMOVED***,
		***REMOVED***
		state := &lib.State***REMOVED***
			Transport: http.DefaultTransport,
			Logger:    logrus.New(),
			Tags:      lib.NewTagMap(nil),
		***REMOVED***
		_, err = MakeRequest(ctx, state, preq)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unknown compressionType CompressionType(13)")
	***REMOVED***)

	t.Run("invalid upgrade response", func(t *testing.T) ***REMOVED***
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			w.Header().Add("Connection", "Upgrade")
			w.Header().Add("Upgrade", "h2c")
			w.WriteHeader(http.StatusSwitchingProtocols)
		***REMOVED***))
		defer srv.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		logger := logrus.New()
		logger.Level = logrus.DebugLevel
		state := &lib.State***REMOVED***
			Transport: srv.Client().Transport,
			Logger:    logger,
			Tags:      lib.NewTagMap(nil),
		***REMOVED***
		req, _ := http.NewRequest("GET", srv.URL, nil)
		preq := &ParsedHTTPRequest***REMOVED***
			Req:     req,
			URL:     &URL***REMOVED***u: req.URL***REMOVED***,
			Body:    new(bytes.Buffer),
			Timeout: 10 * time.Second,
		***REMOVED***

		res, err := MakeRequest(ctx, state, preq)

		assert.Nil(t, res)
		assert.EqualError(t, err, "unsupported response status: 101 Switching Protocols")
	***REMOVED***)
***REMOVED***

func TestResponseStatus(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("response status", func(t *testing.T) ***REMOVED***
		t.Parallel()
		testCases := []struct ***REMOVED***
			name                     string
			statusCode               int
			statusCodeExpected       int
			statusCodeStringExpected string
		***REMOVED******REMOVED***
			***REMOVED***"status 200", 200, 200, "200 OK"***REMOVED***,
			***REMOVED***"status 201", 201, 201, "201 Created"***REMOVED***,
			***REMOVED***"status 202", 202, 202, "202 Accepted"***REMOVED***,
			***REMOVED***"status 203", 203, 203, "203 Non-Authoritative Information"***REMOVED***,
			***REMOVED***"status 204", 204, 204, "204 No Content"***REMOVED***,
			***REMOVED***"status 205", 205, 205, "205 Reset Content"***REMOVED***,
		***REMOVED***

		for _, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(tc.name, func(t *testing.T) ***REMOVED***
				t.Parallel()
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
					w.WriteHeader(tc.statusCode)
				***REMOVED***))
				defer server.Close()
				logger := logrus.New()
				logger.Level = logrus.DebugLevel
				samples := make(chan<- metrics.SampleContainer, 1)
				registry := metrics.NewRegistry()
				state := &lib.State***REMOVED***
					Transport:      server.Client().Transport,
					Logger:         logger,
					Samples:        samples,
					BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
					Tags:           lib.NewTagMap(nil),
				***REMOVED***
				req, err := http.NewRequest("GET", server.URL, nil)
				require.NoError(t, err)

				preq := &ParsedHTTPRequest***REMOVED***
					Req:          req,
					URL:          &URL***REMOVED***u: req.URL***REMOVED***,
					Body:         new(bytes.Buffer),
					Timeout:      10 * time.Second,
					ResponseType: ResponseTypeNone,
				***REMOVED***

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				res, err := MakeRequest(ctx, state, preq)
				require.NoError(t, err)
				assert.Equal(t, tc.statusCodeExpected, res.Status)
				assert.Equal(t, tc.statusCodeStringExpected, res.StatusText)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestURL(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Clean", func(t *testing.T) ***REMOVED***
		t.Parallel()
		testCases := []struct ***REMOVED***
			url      string
			expected string
		***REMOVED******REMOVED***
			***REMOVED***"https://example.com/", "https://example.com/"***REMOVED***,
			***REMOVED***"https://example.com/$***REMOVED******REMOVED***", "https://example.com/$***REMOVED******REMOVED***"***REMOVED***,
			***REMOVED***"https://user@example.com/", "https://****@example.com/"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/", "https://****:****@example.com/"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/path?a=1&b=2", "https://****:****@example.com/path?a=1&b=2"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***", "https://****:****@example.com/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
			***REMOVED***"@malformed/url", "@malformed/url"***REMOVED***,
			***REMOVED***"not a url", "not a url"***REMOVED***,
		***REMOVED***

		for _, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(tc.url, func(t *testing.T) ***REMOVED***
				t.Parallel()
				u, err := url.Parse(tc.url)
				require.NoError(t, err)
				ut := URL***REMOVED***u: u, URL: tc.url***REMOVED***
				require.Equal(t, tc.expected, ut.Clean())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestMakeRequestTimeoutInTheMiddle(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Add("Content-Length", "100000")
		w.WriteHeader(200)
		if f, ok := w.(http.Flusher); ok ***REMOVED***
			f.Flush()
		***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan metrics.SampleContainer, 10)
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	registry := metrics.NewRegistry()
	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			SystemTags: &metrics.DefaultSystemTagSet,
		***REMOVED***,
		Transport:      srv.Client().Transport,
		Samples:        samples,
		Logger:         logger,
		BPool:          bpool.NewBufferPool(100),
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***
	req, _ := http.NewRequest("GET", srv.URL, nil)
	preq := &ParsedHTTPRequest***REMOVED***
		Req:              req,
		URL:              &URL***REMOVED***u: req.URL, URL: srv.URL***REMOVED***,
		Body:             new(bytes.Buffer),
		Timeout:          50 * time.Millisecond,
		ResponseCallback: func(i int) bool ***REMOVED*** return i == 0 ***REMOVED***,
	***REMOVED***

	res, err := MakeRequest(ctx, state, preq)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, samples, 1)
	sampleCont := <-samples
	allSamples := sampleCont.GetSamples()
	require.Len(t, allSamples, 9)
	expTags := map[string]string***REMOVED***
		"error":             "request timeout",
		"error_code":        "1050",
		"status":            "0",
		"expected_response": "true", // we wait for status code 0
		"method":            "GET",
		"url":               srv.URL,
		"name":              srv.URL,
	***REMOVED***
	for _, s := range allSamples ***REMOVED***
		assert.Equal(t, expTags, s.Tags.CloneTags())
	***REMOVED***
***REMOVED***

func BenchmarkWrapDecompressionError(b *testing.B) ***REMOVED***
	err := errors.New("error")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ ***REMOVED***
		_ = wrapDecompressionError(err)
	***REMOVED***
***REMOVED***

func TestTrailFailed(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	t.Cleanup(srv.Close)

	testCases := map[string]struct ***REMOVED***
		responseCallback func(int) bool
		failed           null.Bool
	***REMOVED******REMOVED***
		"null responsecallback": ***REMOVED***responseCallback: nil, failed: null.NewBool(false, false)***REMOVED***,
		"unexpected response":   ***REMOVED***responseCallback: func(int) bool ***REMOVED*** return false ***REMOVED***, failed: null.NewBool(true, true)***REMOVED***,
		"expected response":     ***REMOVED***responseCallback: func(int) bool ***REMOVED*** return true ***REMOVED***, failed: null.NewBool(false, true)***REMOVED***,
	***REMOVED***
	for name, testCase := range testCases ***REMOVED***
		responseCallback := testCase.responseCallback
		failed := testCase.failed

		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			samples := make(chan metrics.SampleContainer, 10)
			logger := logrus.New()
			logger.Level = logrus.DebugLevel
			registry := metrics.NewRegistry()
			state := &lib.State***REMOVED***
				Options: lib.Options***REMOVED***
					SystemTags: &metrics.DefaultSystemTagSet,
				***REMOVED***,
				Transport:      srv.Client().Transport,
				Samples:        samples,
				Logger:         logger,
				BPool:          bpool.NewBufferPool(2),
				BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
				Tags:           lib.NewTagMap(nil),
			***REMOVED***
			req, _ := http.NewRequest("GET", srv.URL, nil)
			preq := &ParsedHTTPRequest***REMOVED***
				Req:              req,
				URL:              &URL***REMOVED***u: req.URL, URL: srv.URL***REMOVED***,
				Body:             new(bytes.Buffer),
				Timeout:          10 * time.Millisecond,
				ResponseCallback: responseCallback,
			***REMOVED***
			res, err := MakeRequest(ctx, state, preq)

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Len(t, samples, 1)
			sample := <-samples
			trail := sample.(*Trail)
			require.Equal(t, failed, trail.Failed)

			var httpReqFailedSampleValue null.Bool
			for _, s := range sample.GetSamples() ***REMOVED***
				if s.Metric.Name == metrics.HTTPReqFailedName ***REMOVED***
					httpReqFailedSampleValue.Valid = true
					if s.Value == 1.0 ***REMOVED***
						httpReqFailedSampleValue.Bool = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
			require.Equal(t, failed, httpReqFailedSampleValue)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMakeRequestDialTimeout(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skipf("dial timeout doesn't get returned on windows") // or we don't match it correctly
	***REMOVED***
	t.Parallel()
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	addr := ln.Addr()
	defer func() ***REMOVED***
		require.NoError(t, ln.Close())
	***REMOVED***()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan metrics.SampleContainer, 10)
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	registry := metrics.NewRegistry()
	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			SystemTags: &metrics.DefaultSystemTagSet,
		***REMOVED***,
		Transport: &http.Transport***REMOVED***
			DialContext: (&net.Dialer***REMOVED***
				Timeout: 1 * time.Microsecond,
			***REMOVED***).DialContext,
		***REMOVED***,
		Samples:        samples,
		Logger:         logger,
		BPool:          bpool.NewBufferPool(100),
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***

	req, _ := http.NewRequest("GET", "http://"+addr.String(), nil)
	preq := &ParsedHTTPRequest***REMOVED***
		Req:              req,
		URL:              &URL***REMOVED***u: req.URL, URL: req.URL.String()***REMOVED***,
		Body:             new(bytes.Buffer),
		Timeout:          500 * time.Millisecond,
		ResponseCallback: func(i int) bool ***REMOVED*** return i == 0 ***REMOVED***,
	***REMOVED***

	res, err := MakeRequest(ctx, state, preq)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, samples, 1)
	sampleCont := <-samples
	allSamples := sampleCont.GetSamples()
	require.Len(t, allSamples, 9)
	expTags := map[string]string***REMOVED***
		"error":             "dial: i/o timeout",
		"error_code":        "1211",
		"status":            "0",
		"expected_response": "true", // we wait for status code 0
		"method":            "GET",
		"url":               req.URL.String(),
		"name":              req.URL.String(),
	***REMOVED***
	for _, s := range allSamples ***REMOVED***
		assert.Equal(t, expTags, s.Tags.CloneTags())
	***REMOVED***
***REMOVED***

func TestMakeRequestTimeoutInTheBegining(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan metrics.SampleContainer, 10)
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	registry := metrics.NewRegistry()
	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			SystemTags: &metrics.DefaultSystemTagSet,
		***REMOVED***,
		Transport:      srv.Client().Transport,
		Samples:        samples,
		Logger:         logger,
		BPool:          bpool.NewBufferPool(100),
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Tags:           lib.NewTagMap(nil),
	***REMOVED***
	req, _ := http.NewRequest("GET", srv.URL, nil)
	preq := &ParsedHTTPRequest***REMOVED***
		Req:              req,
		URL:              &URL***REMOVED***u: req.URL, URL: srv.URL***REMOVED***,
		Body:             new(bytes.Buffer),
		Timeout:          50 * time.Millisecond,
		ResponseCallback: func(i int) bool ***REMOVED*** return i == 0 ***REMOVED***,
	***REMOVED***

	res, err := MakeRequest(ctx, state, preq)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, samples, 1)
	sampleCont := <-samples
	allSamples := sampleCont.GetSamples()
	require.Len(t, allSamples, 9)
	expTags := map[string]string***REMOVED***
		"error":             "request timeout",
		"error_code":        "1050",
		"status":            "0",
		"expected_response": "true", // we wait for status code 0
		"method":            "GET",
		"url":               srv.URL,
		"name":              srv.URL,
	***REMOVED***
	for _, s := range allSamples ***REMOVED***
		assert.Equal(t, expTags, s.Tags.CloneTags())
	***REMOVED***
***REMOVED***
