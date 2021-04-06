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

package influxdb

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
)

func TestBadConcurrentWrites(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := testutils.NewLogger(t)
	t.Run("0", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, err := New(output.Params***REMOVED***
			Logger:         logger,
			ConfigArgument: "?concurrentWrites=0",
		***REMOVED***)
		require.Error(t, err)
		require.Equal(t, err.Error(), "influxdb's ConcurrentWrites must be a positive number")
	***REMOVED***)

	t.Run("-2", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, err := New(output.Params***REMOVED***
			Logger:         logger,
			ConfigArgument: "?concurrentWrites=-2",
		***REMOVED***)
		require.Error(t, err)
		require.Equal(t, err.Error(), "influxdb's ConcurrentWrites must be a positive number")
	***REMOVED***)

	t.Run("2", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, err := New(output.Params***REMOVED***
			Logger:         logger,
			ConfigArgument: "?concurrentWrites=2",
		***REMOVED***)
		require.NoError(t, err)
	***REMOVED***)
***REMOVED***

func testOutputCycle(t testing.TB, handler http.HandlerFunc, body func(testing.TB, *Output)) ***REMOVED***
	s := &http.Server***REMOVED***
		Addr:           ":",
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	***REMOVED***
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() ***REMOVED***
		_ = l.Close()
	***REMOVED***()

	defer func() ***REMOVED***
		require.NoError(t, s.Shutdown(context.Background()))
	***REMOVED***()

	go func() ***REMOVED***
		require.Equal(t, http.ErrServerClosed, s.Serve(l))
	***REMOVED***()

	c, err := newOutput(output.Params***REMOVED***
		Logger:         testutils.NewLogger(t),
		ConfigArgument: "http://" + l.Addr().String(),
	***REMOVED***)
	require.NoError(t, err)

	require.NoError(t, c.Start())
	body(t, c)

	require.NoError(t, c.Stop())
***REMOVED***

func TestOutput(t *testing.T) ***REMOVED***
	t.Parallel()
	var samplesRead int
	defer func() ***REMOVED***
		require.Equal(t, samplesRead, 20)
	***REMOVED***()
	testOutputCycle(t, func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		b := bytes.NewBuffer(nil)
		_, _ = io.Copy(b, r.Body)
		for ***REMOVED***
			s, err := b.ReadString('\n')
			if len(s) > 0 ***REMOVED***
				samplesRead++
			***REMOVED***
			if err != nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		rw.WriteHeader(204)
	***REMOVED***, func(tb testing.TB, c *Output) ***REMOVED***
		samples := make(stats.Samples, 10)
		for i := 0; i < len(samples); i++ ***REMOVED***
			samples[i] = stats.Sample***REMOVED***
				Metric: stats.New("testGauge", stats.Gauge),
				Time:   time.Now(),
				Tags: stats.NewSampleTags(map[string]string***REMOVED***
					"something": "else",
					"VU":        "21",
					"else":      "something",
				***REMOVED***),
				Value: 2.0,
			***REMOVED***
		***REMOVED***
		c.AddMetricSamples([]stats.SampleContainer***REMOVED***samples***REMOVED***)
		c.AddMetricSamples([]stats.SampleContainer***REMOVED***samples***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestExtractTagsToValues(t *testing.T) ***REMOVED***
	t.Parallel()
	o, err := newOutput(output.Params***REMOVED***
		Logger:         testutils.NewLogger(t),
		ConfigArgument: "?tagsAsFields=stringField&tagsAsFields=stringField2:string&tagsAsFields=boolField:bool&tagsAsFields=floatField:float&tagsAsFields=intField:int",
	***REMOVED***)
	require.NoError(t, err)
	tags := map[string]string***REMOVED***
		"stringField":  "string",
		"stringField2": "string2",
		"boolField":    "true",
		"floatField":   "3.14",
		"intField":     "12345",
	***REMOVED***
	values := o.extractTagsToValues(tags, map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***)

	require.Equal(t, "string", values["stringField"])
	require.Equal(t, "string2", values["stringField2"])
	require.Equal(t, true, values["boolField"])
	require.Equal(t, 3.14, values["floatField"])
	require.Equal(t, int64(12345), values["intField"])
***REMOVED***
