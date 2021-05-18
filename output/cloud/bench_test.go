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

package cloud

import (
	"bytes"
	"compress/gzip"
	json "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
)

func BenchmarkAggregateHTTP(b *testing.B) ***REMOVED***
	out, err := newOutput(output.Params***REMOVED***
		Logger:     testutils.NewLogger(b),
		JSONConfig: json.RawMessage(`***REMOVED***"noCompress": true, "aggregationCalcInterval": "200ms","aggregationPeriod": "200ms"***REMOVED***`),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(b, err)
	now := time.Now()
	out.referenceID = "something"
	containersCount := 500000

	for _, tagCount := range []int***REMOVED***1, 5, 35, 315, 3645***REMOVED*** ***REMOVED***
		tagCount := tagCount
		b.Run(fmt.Sprintf("tags:%d", tagCount), func(b *testing.B) ***REMOVED***
			b.ResetTimer()
			for s := 0; s < b.N; s++ ***REMOVED***
				b.StopTimer()
				container := make([]stats.SampleContainer, containersCount)
				for i := 1; i <= containersCount; i++ ***REMOVED***
					status := "200"
					if i%tagCount%7 == 6 ***REMOVED***
						status = "404"
					***REMOVED*** else if i%tagCount%7 == 5 ***REMOVED***
						status = "500"
					***REMOVED***

					tags := generateTags(i, tagCount, map[string]string***REMOVED***"status": status***REMOVED***)
					container[i-1] = generateHTTPExtTrail(now, time.Duration(i), tags)
				***REMOVED***
				out.AddMetricSamples(container)
				b.StartTimer()
				out.aggregateHTTPTrails(time.Millisecond * 200)
				out.bufferSamples = nil
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func generateTags(i, tagCount int, additionals ...map[string]string) *stats.SampleTags ***REMOVED***
	res := map[string]string***REMOVED***
		"test": "mest", "a": "b",
		"custom": fmt.Sprintf("group%d", i%tagCount%9),
		"group":  fmt.Sprintf("group%d", i%tagCount%5),
		"url":    fmt.Sprintf("something%d", i%tagCount%11),
		"name":   fmt.Sprintf("else%d", i%tagCount%11),
	***REMOVED***
	for _, a := range additionals ***REMOVED***
		for k, v := range a ***REMOVED***
			res[k] = v
		***REMOVED***
	***REMOVED***

	return stats.IntoSampleTags(&res)
***REMOVED***

func BenchmarkMetricMarshal(b *testing.B) ***REMOVED***
	for _, count := range []int***REMOVED***10000, 100000, 500000***REMOVED*** ***REMOVED***
		count := count
		b.Run(fmt.Sprintf("%d", count), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				b.StopTimer()
				s := generateSamples(count)
				b.StartTimer()
				r, err := easyjson.Marshal(samples(s))
				require.NoError(b, err)
				b.SetBytes(int64(len(r)))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkMetricMarshalWriter(b *testing.B) ***REMOVED***
	for _, count := range []int***REMOVED***10000, 100000, 500000***REMOVED*** ***REMOVED***
		count := count
		b.Run(fmt.Sprintf("%d", count), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				b.StopTimer()
				s := generateSamples(count)
				b.StartTimer()
				n, err := easyjson.MarshalToWriter(samples(s), ioutil.Discard)
				require.NoError(b, err)
				b.SetBytes(int64(n))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkMetricMarshalGzip(b *testing.B) ***REMOVED***
	for _, count := range []int***REMOVED***10000, 100000, 500000***REMOVED*** ***REMOVED***
		for name, level := range map[string]int***REMOVED***
			"bestcompression": gzip.BestCompression,
			"default":         gzip.DefaultCompression,
			"bestspeed":       gzip.BestSpeed,
		***REMOVED*** ***REMOVED***
			count := count
			level := level
			b.Run(fmt.Sprintf("%d_%s", count, name), func(b *testing.B) ***REMOVED***
				s := generateSamples(count)
				r, err := easyjson.Marshal(samples(s))
				require.NoError(b, err)
				b.ResetTimer()
				for i := 0; i < b.N; i++ ***REMOVED***
					b.StopTimer()
					var buf bytes.Buffer
					buf.Grow(len(r) / 5)
					g, err := gzip.NewWriterLevel(&buf, level)
					require.NoError(b, err)
					b.StartTimer()
					n, err := g.Write(r)
					require.NoError(b, err)
					b.SetBytes(int64(n))
					b.ReportMetric(float64(len(r))/float64(buf.Len()), "ratio")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMetricMarshalGzipAll(b *testing.B) ***REMOVED***
	for _, count := range []int***REMOVED***10000, 100000, 500000***REMOVED*** ***REMOVED***
		for name, level := range map[string]int***REMOVED***
			"bestspeed": gzip.BestSpeed,
		***REMOVED*** ***REMOVED***
			count := count
			level := level
			b.Run(fmt.Sprintf("%d_%s", count, name), func(b *testing.B) ***REMOVED***
				for i := 0; i < b.N; i++ ***REMOVED***
					b.StopTimer()

					s := generateSamples(count)
					var buf bytes.Buffer
					g, err := gzip.NewWriterLevel(&buf, level)
					require.NoError(b, err)
					b.StartTimer()

					r, err := easyjson.Marshal(samples(s))
					require.NoError(b, err)
					buf.Grow(len(r) / 5)
					n, err := g.Write(r)
					require.NoError(b, err)
					b.SetBytes(int64(n))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMetricMarshalGzipAllWriter(b *testing.B) ***REMOVED***
	for _, count := range []int***REMOVED***10000, 100000, 500000***REMOVED*** ***REMOVED***
		for name, level := range map[string]int***REMOVED***
			"bestspeed": gzip.BestSpeed,
		***REMOVED*** ***REMOVED***
			count := count
			level := level
			b.Run(fmt.Sprintf("%d_%s", count, name), func(b *testing.B) ***REMOVED***
				var buf bytes.Buffer
				for i := 0; i < b.N; i++ ***REMOVED***
					b.StopTimer()
					buf.Reset()

					s := generateSamples(count)
					g, err := gzip.NewWriterLevel(&buf, level)
					require.NoError(b, err)
					pr, pw := io.Pipe()
					b.StartTimer()

					go func() ***REMOVED***
						_, _ = easyjson.MarshalToWriter(samples(s), pw)
						_ = pw.Close()
					***REMOVED***()
					n, err := io.Copy(g, pr)
					require.NoError(b, err)
					b.SetBytes(n)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func generateSamples(count int) []*Sample ***REMOVED***
	samples := make([]*Sample, count)
	now := time.Now()
	for i := range samples ***REMOVED***
		tags := generateTags(i, 200)
		switch i % 3 ***REMOVED***
		case 0:
			samples[i] = &Sample***REMOVED***
				Type:   DataTypeSingle,
				Metric: "something",
				Data: &SampleDataSingle***REMOVED***
					Time:  toMicroSecond(now),
					Type:  stats.Counter,
					Tags:  tags,
					Value: float64(i),
				***REMOVED***,
			***REMOVED***
		case 1:
			aggrData := &SampleDataAggregatedHTTPReqs***REMOVED***
				Time: toMicroSecond(now),
				Type: "aggregated_trend",
				Tags: tags,
			***REMOVED***
			trail := generateHTTPExtTrail(now, time.Duration(i), tags)
			aggrData.Add(trail)
			aggrData.Add(trail)
			aggrData.Add(trail)
			aggrData.Add(trail)
			aggrData.Add(trail)
			aggrData.CalcAverages()

			samples[i] = &Sample***REMOVED***
				Type:   DataTypeAggregatedHTTPReqs,
				Metric: "something",
				Data:   aggrData,
			***REMOVED***
		default:
			samples[i] = NewSampleFromTrail(generateHTTPExtTrail(now, time.Duration(i), tags))
		***REMOVED***
	***REMOVED***

	return samples
***REMOVED***

func generateHTTPExtTrail(now time.Time, i time.Duration, tags *stats.SampleTags) *httpext.Trail ***REMOVED***
	return &httpext.Trail***REMOVED***
		Blocked:        i % 200 * 100 * time.Millisecond,
		Connecting:     i % 200 * 200 * time.Millisecond,
		TLSHandshaking: i % 200 * 300 * time.Millisecond,
		Sending:        i % 200 * 400 * time.Millisecond,
		Waiting:        500 * time.Millisecond,
		Receiving:      600 * time.Millisecond,
		EndTime:        now.Add(i % 100 * 100),
		ConnDuration:   500 * time.Millisecond,
		Duration:       i % 150 * 1500 * time.Millisecond,
		Tags:           tags,
	***REMOVED***
***REMOVED***

func BenchmarkHTTPPush(b *testing.B) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(b)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprint(w, `***REMOVED***
			"reference_id": "fake",
		***REMOVED***`)
		require.NoError(b, err)
	***REMOVED***))
	tb.Mux.HandleFunc("/v1/metrics/fake",
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			_, err := io.Copy(ioutil.Discard, r.Body)
			assert.NoError(b, err)
		***REMOVED***,
	)

	out, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(b),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***
			"host": "%s",
			"noCompress": false,
			"aggregationCalcInterval": "200ms",
			"aggregationPeriod": "200ms"
		***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(b, err)
	out.referenceID = "fake"
	assert.False(b, out.config.NoCompress.Bool)

	for _, count := range []int***REMOVED***1000, 5000, 50000, 100000, 250000***REMOVED*** ***REMOVED***
		count := count
		b.Run(fmt.Sprintf("count:%d", count), func(b *testing.B) ***REMOVED***
			samples := generateSamples(count)
			b.ResetTimer()
			for s := 0; s < b.N; s++ ***REMOVED***
				b.StopTimer()
				toSend := append([]*Sample***REMOVED******REMOVED***, samples...)
				b.StartTimer()
				require.NoError(b, out.client.PushMetric("fake", toSend))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkNewSampleFromTrail(b *testing.B) ***REMOVED***
	tags := generateTags(1, 200)
	now := time.Now()
	trail := &httpext.Trail***REMOVED***
		Blocked:        200 * 100 * time.Millisecond,
		Connecting:     200 * 200 * time.Millisecond,
		TLSHandshaking: 200 * 300 * time.Millisecond,
		Sending:        200 * 400 * time.Millisecond,
		Waiting:        500 * time.Millisecond,
		Receiving:      600 * time.Millisecond,
		EndTime:        now,
		ConnDuration:   500 * time.Millisecond,
		Duration:       150 * 1500 * time.Millisecond,
		Tags:           tags,
	***REMOVED***

	b.Run("no failed", func(b *testing.B) ***REMOVED***
		for s := 0; s < b.N; s++ ***REMOVED***
			_ = NewSampleFromTrail(trail)
		***REMOVED***
	***REMOVED***)
	trail.Failed = null.BoolFrom(true)

	b.Run("failed", func(b *testing.B) ***REMOVED***
		for s := 0; s < b.N; s++ ***REMOVED***
			_ = NewSampleFromTrail(trail)
		***REMOVED***
	***REMOVED***)
***REMOVED***
