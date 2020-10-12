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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
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
	algos := []CompressionType***REMOVED***CompressionTypeGzip***REMOVED***
	t.Run("bad read body", func(t *testing.T) ***REMOVED***
		_, _, err := compressBody(algos, ioutil.NopCloser(badReadBody()))
		require.Error(t, err)
		require.Equal(t, err.Error(), badReadMsg)
	***REMOVED***)

	t.Run("bad close body", func(t *testing.T) ***REMOVED***
		_, _, err := compressBody(algos, badCloseBody())
		require.Error(t, err)
		require.Equal(t, err.Error(), badCloseMsg)
	***REMOVED***)
***REMOVED***

func TestMakeRequestError(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("bad compression algorithm body", func(t *testing.T) ***REMOVED***
		req, err := http.NewRequest("GET", "https://wont.be.used", nil)

		require.NoError(t, err)
		badCompressionType := CompressionType(13)
		require.False(t, badCompressionType.IsACompressionType())
		preq := &ParsedHTTPRequest***REMOVED***
			Req:          req,
			Body:         new(bytes.Buffer),
			Compressions: []CompressionType***REMOVED***badCompressionType***REMOVED***,
		***REMOVED***
		_, err = MakeRequest(ctx, preq)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unknown compressionType CompressionType(13)")
	***REMOVED***)

	t.Run("invalid upgrade response", func(t *testing.T) ***REMOVED***
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
			Options:   lib.Options***REMOVED***RunTags: &stats.SampleTags***REMOVED******REMOVED******REMOVED***,
			Transport: srv.Client().Transport,
			Logger:    logger,
		***REMOVED***
		ctx = lib.WithState(ctx, state)
		req, _ := http.NewRequest("GET", srv.URL, nil)
		preq := &ParsedHTTPRequest***REMOVED***
			Req:     req,
			URL:     &URL***REMOVED***u: req.URL***REMOVED***,
			Body:    new(bytes.Buffer),
			Timeout: 10 * time.Second,
		***REMOVED***

		res, err := MakeRequest(ctx, preq)

		assert.Nil(t, res)
		assert.EqualError(t, err, "unsupported response status: 101 Switching Protocols")
	***REMOVED***)
***REMOVED***

func TestResponseStatus(t *testing.T) ***REMOVED***
	t.Run("response status", func(t *testing.T) ***REMOVED***
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
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
					w.WriteHeader(tc.statusCode)
				***REMOVED***))
				defer server.Close()
				logger := logrus.New()
				logger.Level = logrus.DebugLevel
				state := &lib.State***REMOVED***
					Options:   lib.Options***REMOVED***RunTags: &stats.SampleTags***REMOVED******REMOVED******REMOVED***,
					Transport: server.Client().Transport,
					Logger:    logger,
					Samples:   make(chan<- stats.SampleContainer, 1),
				***REMOVED***
				ctx := lib.WithState(context.Background(), state)
				req, err := http.NewRequest("GET", server.URL, nil)
				require.NoError(t, err)

				preq := &ParsedHTTPRequest***REMOVED***
					Req:          req,
					URL:          &URL***REMOVED***u: req.URL***REMOVED***,
					Body:         new(bytes.Buffer),
					Timeout:      10 * time.Second,
					ResponseType: ResponseTypeNone,
				***REMOVED***

				res, err := MakeRequest(ctx, preq)
				require.NoError(t, err)
				assert.Equal(t, tc.statusCodeExpected, res.Status)
				assert.Equal(t, tc.statusCodeStringExpected, res.StatusText)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestURL(t *testing.T) ***REMOVED***
	t.Run("Clean", func(t *testing.T) ***REMOVED***
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
				u, err := url.Parse(tc.url)
				require.NoError(t, err)
				ut := URL***REMOVED***u: u, URL: tc.url***REMOVED***
				require.Equal(t, tc.expected, ut.Clean())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestMakeRequestTimeout(t *testing.T) ***REMOVED***
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan stats.SampleContainer, 10)
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			RunTags:    &stats.SampleTags***REMOVED******REMOVED***,
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		Transport: srv.Client().Transport,
		Samples:   samples,
		Logger:    logger,
	***REMOVED***
	ctx = lib.WithState(ctx, state)
	req, _ := http.NewRequest("GET", srv.URL, nil)
	preq := &ParsedHTTPRequest***REMOVED***
		Req:     req,
		URL:     &URL***REMOVED***u: req.URL, URL: srv.URL***REMOVED***,
		Body:    new(bytes.Buffer),
		Timeout: 10 * time.Millisecond,
	***REMOVED***

	res, err := MakeRequest(ctx, preq)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, samples, 1)
	sampleCont := <-samples
	allSamples := sampleCont.GetSamples()
	require.Len(t, allSamples, 8)
	expTags := map[string]string***REMOVED***
		"error":      "context deadline exceeded",
		"error_code": "1000",
		"status":     "0",
		"method":     "GET",
		"url":        srv.URL,
		"name":       srv.URL,
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
