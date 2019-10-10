/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
)

func getSampleChecker(t *testing.T, expSamples <-chan []Sample) http.HandlerFunc ***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		receivedSamples := []Sample***REMOVED******REMOVED***
		assert.NoError(t, json.Unmarshal(body, &receivedSamples))

		expSamples := <-expSamples
		if !assert.Len(t, receivedSamples, len(expSamples)) ***REMOVED***
			return
		***REMOVED***

		for i, expSample := range expSamples ***REMOVED***
			receivedSample := receivedSamples[i]
			assert.Equal(t, expSample.Metric, receivedSample.Metric)
			assert.Equal(t, expSample.Type, receivedSample.Type)

			if callbackCheck, ok := expSample.Data.(func(interface***REMOVED******REMOVED***)); ok ***REMOVED***
				callbackCheck(receivedSample.Data)
				continue
			***REMOVED***

			if !assert.IsType(t, expSample.Data, receivedSample.Data) ***REMOVED***
				continue
			***REMOVED***

			switch expData := expSample.Data.(type) ***REMOVED***
			case *SampleDataSingle:
				receivedData, ok := receivedSample.Data.(*SampleDataSingle)
				assert.True(t, ok)
				assert.True(t, expData.Tags.IsEqual(receivedData.Tags))
				assert.True(t, expData.Time.Equal(receivedData.Time))
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Value, receivedData.Value)
			case *SampleDataMap:
				receivedData, ok := receivedSample.Data.(*SampleDataMap)
				assert.True(t, ok)
				assert.True(t, expData.Tags.IsEqual(receivedData.Tags))
				assert.True(t, expData.Time.Equal(receivedData.Time))
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Values, receivedData.Values)
			case *SampleDataAggregatedHTTPReqs:
				receivedData, ok := receivedSample.Data.(*SampleDataAggregatedHTTPReqs)
				assert.True(t, ok)
				assert.True(t, expData.Tags.IsEqual(receivedData.Tags))
				assert.True(t, expData.Time.Equal(receivedData.Time))
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Values, receivedData.Values)
			default:
				t.Errorf("Unknown data type %#v", expData)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func skewTrail(t httpext.Trail, minCoef, maxCoef float64) httpext.Trail ***REMOVED***
	coef := minCoef + rand.Float64()*(maxCoef-minCoef)
	addJitter := func(d *time.Duration) ***REMOVED***
		*d = time.Duration(float64(*d) * coef)
	***REMOVED***
	addJitter(&t.Blocked)
	addJitter(&t.Connecting)
	addJitter(&t.TLSHandshaking)
	addJitter(&t.Sending)
	addJitter(&t.Waiting)
	addJitter(&t.Receiving)
	t.ConnDuration = t.Connecting + t.TLSHandshaking
	t.Duration = t.Sending + t.Waiting + t.Receiving
	t.StartTime = t.EndTime.Add(-t.Duration)
	return t
***REMOVED***

func TestCloudCollector(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprintf(w, `***REMOVED***
			"reference_id": "123",
			"config": ***REMOVED***
				"metricPushInterval": "10ms",
				"aggregationPeriod": "30ms",
				"aggregationCalcInterval": "40ms",
				"aggregationWaitPeriod": "5ms"
			***REMOVED***
		***REMOVED***`)
		require.NoError(t, err)
	***REMOVED***))
	defer tb.Cleanup()

	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		Host:       null.StringFrom(tb.ServerHTTP.URL),
		NoCompress: null.BoolFrom(true),
	***REMOVED***)
	collector, err := New(config, script, options, []lib.ExecutionStep***REMOVED******REMOVED***, "1.0")
	require.NoError(t, err)

	assert.True(t, collector.config.Host.Valid)
	assert.Equal(t, tb.ServerHTTP.URL, collector.config.Host.String)
	assert.True(t, collector.config.NoCompress.Valid)
	assert.True(t, collector.config.NoCompress.Bool)
	assert.False(t, collector.config.MetricPushInterval.Valid)
	assert.False(t, collector.config.AggregationPeriod.Valid)
	assert.False(t, collector.config.AggregationWaitPeriod.Valid)

	require.NoError(t, collector.Init())
	assert.Equal(t, "123", collector.referenceID)
	assert.True(t, collector.config.MetricPushInterval.Valid)
	assert.Equal(t, types.Duration(10*time.Millisecond), collector.config.MetricPushInterval.Duration)
	assert.True(t, collector.config.AggregationPeriod.Valid)
	assert.Equal(t, types.Duration(30*time.Millisecond), collector.config.AggregationPeriod.Duration)
	assert.True(t, collector.config.AggregationWaitPeriod.Valid)
	assert.Equal(t, types.Duration(5*time.Millisecond), collector.config.AggregationWaitPeriod.Duration)

	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)

	expSamples := make(chan []Sample)
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", collector.referenceID), getSampleChecker(t, expSamples))

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()

	collector.Collect([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   tags,
		Value:  1.0,
	***REMOVED******REMOVED***)
	expSamples <- []Sample***REMOVED******REMOVED***
		Type:   DataTypeSingle,
		Metric: metrics.VUs.Name,
		Data: &SampleDataSingle***REMOVED***
			Type:  metrics.VUs.Type,
			Time:  Timestamp(now),
			Tags:  tags,
			Value: 1.0,
		***REMOVED***,
	***REMOVED******REMOVED***

	simpleTrail := httpext.Trail***REMOVED***
		Blocked:        100 * time.Millisecond,
		Connecting:     200 * time.Millisecond,
		TLSHandshaking: 300 * time.Millisecond,
		Sending:        400 * time.Millisecond,
		Waiting:        500 * time.Millisecond,
		Receiving:      600 * time.Millisecond,

		EndTime:      now,
		ConnDuration: 500 * time.Millisecond,
		Duration:     1500 * time.Millisecond,
		Tags:         tags,
	***REMOVED***
	collector.Collect([]stats.SampleContainer***REMOVED***&simpleTrail***REMOVED***)
	expSamples <- []Sample***REMOVED****NewSampleFromTrail(&simpleTrail)***REMOVED***

	smallSkew := 0.05

	trails := []stats.SampleContainer***REMOVED******REMOVED***
	for i := int64(0); i < collector.config.AggregationMinSamples.Int64; i++ ***REMOVED***
		similarTrail := skewTrail(simpleTrail, 1.0, 1.0+smallSkew)
		trails = append(trails, &similarTrail)
	***REMOVED***

	checkAggrMetric := func(normal time.Duration, aggr AggregatedMetric) ***REMOVED***
		assert.True(t, aggr.Min <= aggr.Avg)
		assert.True(t, aggr.Avg <= aggr.Max)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Min), smallSkew)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Avg), smallSkew)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Max), smallSkew)
	***REMOVED***

	outlierTrail := skewTrail(simpleTrail, 2.0+smallSkew, 3.0+smallSkew)
	trails = append(trails, &outlierTrail)
	collector.Collect(trails)
	expSamples <- []Sample***REMOVED***
		*NewSampleFromTrail(&outlierTrail),
		***REMOVED***
			Type:   DataTypeAggregatedHTTPReqs,
			Metric: "http_req_li_all",
			Data: func(data interface***REMOVED******REMOVED***) ***REMOVED***
				aggrData, ok := data.(*SampleDataAggregatedHTTPReqs)
				assert.True(t, ok)
				assert.True(t, aggrData.Tags.IsEqual(tags))
				assert.Equal(t, collector.config.AggregationMinSamples.Int64, int64(aggrData.Count))
				assert.Equal(t, "aggregated_trend", aggrData.Type)
				assert.InDelta(t, now.UnixNano(), time.Time(aggrData.Time).UnixNano(), float64(collector.config.AggregationPeriod.Duration))

				checkAggrMetric(simpleTrail.Duration, aggrData.Values.Duration)
				checkAggrMetric(simpleTrail.Blocked, aggrData.Values.Blocked)
				checkAggrMetric(simpleTrail.Connecting, aggrData.Values.Connecting)
				checkAggrMetric(simpleTrail.TLSHandshaking, aggrData.Values.TLSHandshaking)
				checkAggrMetric(simpleTrail.Sending, aggrData.Values.Sending)
				checkAggrMetric(simpleTrail.Waiting, aggrData.Values.Waiting)
				checkAggrMetric(simpleTrail.Receiving, aggrData.Values.Receiving)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	cancel()
	wg.Wait()
***REMOVED***

func TestCloudCollectorMaxPerPacket(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	var maxMetricSamplesPerPackage = 20
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprintf(w, `***REMOVED***
			"reference_id": "12",
			"config": ***REMOVED***
				"metricPushInterval": "200ms",
				"aggregationPeriod": "100ms",
				"maxMetricSamplesPerPackage": %d,
				"aggregationCalcInterval": "100ms",
				"aggregationWaitPeriod": "100ms"
			***REMOVED***
		***REMOVED***`, maxMetricSamplesPerPackage)
		require.NoError(t, err)
	***REMOVED***))
	defer tb.Cleanup()

	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		Host:       null.StringFrom(tb.ServerHTTP.URL),
		NoCompress: null.BoolFrom(true),
	***REMOVED***)
	collector, err := New(config, script, options, []lib.ExecutionStep***REMOVED******REMOVED***, "1.0")
	require.NoError(t, err)
	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)
	var gotTheLimit = false
	var m sync.Mutex

	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", collector.referenceID),
		func(_ http.ResponseWriter, r *http.Request) ***REMOVED***
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			receivedSamples := []Sample***REMOVED******REMOVED***
			assert.NoError(t, json.Unmarshal(body, &receivedSamples))
			assert.True(t, len(receivedSamples) <= maxMetricSamplesPerPackage)
			if len(receivedSamples) == maxMetricSamplesPerPackage ***REMOVED***
				m.Lock()
				gotTheLimit = true
				m.Unlock()
			***REMOVED***
		***REMOVED***)

	require.NoError(t, collector.Init())
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()

	collector.Collect([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	for j := time.Duration(1); j <= 200; j++ ***REMOVED***
		var container = make([]stats.SampleContainer, 0, 500)
		for i := time.Duration(1); i <= 50; i++ ***REMOVED***
			container = append(container, &httpext.Trail***REMOVED***
				Blocked:        i % 200 * 100 * time.Millisecond,
				Connecting:     i % 200 * 200 * time.Millisecond,
				TLSHandshaking: i % 200 * 300 * time.Millisecond,
				Sending:        i * i * 400 * time.Millisecond,
				Waiting:        500 * time.Millisecond,
				Receiving:      600 * time.Millisecond,

				EndTime:      now.Add(i * 100),
				ConnDuration: 500 * time.Millisecond,
				Duration:     j * i * 1500 * time.Millisecond,
				Tags:         stats.NewSampleTags(tags.CloneTags()),
			***REMOVED***)
		***REMOVED***
		collector.Collect(container)
	***REMOVED***

	cancel()
	wg.Wait()
	require.True(t, gotTheLimit)
***REMOVED***

func TestCloudCollectorStopSendingMetric(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprint(w, `***REMOVED***
			"reference_id": "12",
			"config": ***REMOVED***
				"metricPushInterval": "200ms",
				"aggregationPeriod": "100ms",
				"maxMetricSamplesPerPackage": 20,
				"aggregationCalcInterval": "100ms",
				"aggregationWaitPeriod": "100ms"
			***REMOVED***
		***REMOVED***`)
		require.NoError(t, err)
	***REMOVED***))
	defer tb.Cleanup()

	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		Host:       null.StringFrom(tb.ServerHTTP.URL),
		NoCompress: null.BoolFrom(true),
	***REMOVED***)
	collector, err := New(config, script, options, "1.0")
	require.NoError(t, err)
	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)

	count := 1
	max := 5
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", collector.referenceID),
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			count++
			if count == max ***REMOVED***
				type payload struct ***REMOVED***
					Error ErrorResponse `json:"error"`
				***REMOVED***
				res := &payload***REMOVED******REMOVED***
				res.Error = ErrorResponse***REMOVED***Code: 4***REMOVED***
				w.Header().Set("Content-Type", "application/json")
				data, err := json.Marshal(res)
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write(data)
				return
			***REMOVED***
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			receivedSamples := []Sample***REMOVED******REMOVED***
			assert.NoError(t, json.Unmarshal(body, &receivedSamples))
		***REMOVED***)

	require.NoError(t, collector.Init())
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()

	collector.Collect([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	for j := time.Duration(1); j <= 200; j++ ***REMOVED***
		var container = make([]stats.SampleContainer, 0, 500)
		for i := time.Duration(1); i <= 50; i++ ***REMOVED***
			container = append(container, &httpext.Trail***REMOVED***
				Blocked:        i % 200 * 100 * time.Millisecond,
				Connecting:     i % 200 * 200 * time.Millisecond,
				TLSHandshaking: i % 200 * 300 * time.Millisecond,
				Sending:        i * i * 400 * time.Millisecond,
				Waiting:        500 * time.Millisecond,
				Receiving:      600 * time.Millisecond,

				EndTime:      now.Add(i * 100),
				ConnDuration: 500 * time.Millisecond,
				Duration:     j * i * 1500 * time.Millisecond,
				Tags:         stats.NewSampleTags(tags.CloneTags()),
			***REMOVED***)
		***REMOVED***
		collector.Collect(container)
	***REMOVED***

	cancel()
	wg.Wait()
	require.Equal(t, lib.RunStatusQueued, collector.runStatus)
	_, ok := <-collector.stopSendingMetricsCh
	require.False(t, ok)
	require.Equal(t, max, count)

	nBufferSamples := len(collector.bufferSamples)
	nBufferHTTPTrails := len(collector.bufferHTTPTrails)
	collector.Collect([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	if nBufferSamples != len(collector.bufferSamples) || nBufferHTTPTrails != len(collector.bufferHTTPTrails) ***REMOVED***
		t.Errorf("Collector still collects data after stop sending metrics")
	***REMOVED***
***REMOVED***

func TestCloudCollectorAggregationPeriodZeroNoBlock(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprintf(w, `***REMOVED***
			"reference_id": "123",
			"config": ***REMOVED***
				"metricPushInterval": "10ms",
				"aggregationPeriod": "0ms",
				"aggregationCalcInterval": "40ms",
				"aggregationWaitPeriod": "5ms"
			***REMOVED***
		***REMOVED***`)
		require.NoError(t, err)
	***REMOVED***))
	defer tb.Cleanup()

	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		Host:       null.StringFrom(tb.ServerHTTP.URL),
		NoCompress: null.BoolFrom(true),
	***REMOVED***)
	collector, err := New(config, script, options, "1.0")
	require.NoError(t, err)

	assert.True(t, collector.config.Host.Valid)
	assert.Equal(t, tb.ServerHTTP.URL, collector.config.Host.String)
	assert.True(t, collector.config.NoCompress.Valid)
	assert.True(t, collector.config.NoCompress.Bool)
	assert.False(t, collector.config.MetricPushInterval.Valid)
	assert.False(t, collector.config.AggregationPeriod.Valid)
	assert.False(t, collector.config.AggregationWaitPeriod.Valid)

	require.NoError(t, collector.Init())
	assert.Equal(t, "123", collector.referenceID)
	assert.True(t, collector.config.MetricPushInterval.Valid)
	assert.Equal(t, types.Duration(10*time.Millisecond), collector.config.MetricPushInterval.Duration)
	assert.True(t, collector.config.AggregationPeriod.Valid)
	assert.Equal(t, types.Duration(0), collector.config.AggregationPeriod.Duration)
	assert.True(t, collector.config.AggregationWaitPeriod.Valid)
	assert.Equal(t, types.Duration(5*time.Millisecond), collector.config.AggregationWaitPeriod.Duration)

	expSamples := make(chan []Sample)
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", collector.referenceID), getSampleChecker(t, expSamples))

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()

	cancel()
	wg.Wait()
	require.Equal(t, lib.RunStatusQueued, collector.runStatus)
***REMOVED***
