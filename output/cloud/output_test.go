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
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/cloudapi"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
)

func tagEqual(expected, got *stats.SampleTags) bool ***REMOVED***
	expectedMap := expected.CloneTags()
	gotMap := got.CloneTags()

	if len(expectedMap) != len(gotMap) ***REMOVED***
		return false
	***REMOVED***

	for k, v := range gotMap ***REMOVED***
		if k == "url" ***REMOVED***
			if expectedMap["name"] != v ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else if expectedMap[k] != v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

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
				assert.Equal(t, expData.Time, receivedData.Time)
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Value, receivedData.Value)
			case *SampleDataMap:
				receivedData, ok := receivedSample.Data.(*SampleDataMap)
				assert.True(t, ok)
				assert.True(t, tagEqual(expData.Tags, receivedData.Tags))
				assert.Equal(t, expData.Time, receivedData.Time)
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Values, receivedData.Values)
			case *SampleDataAggregatedHTTPReqs:
				receivedData, ok := receivedSample.Data.(*SampleDataAggregatedHTTPReqs)
				assert.True(t, ok)
				assert.True(t, expData.Tags.IsEqual(receivedData.Tags))
				assert.Equal(t, expData.Time, receivedData.Time)
				assert.Equal(t, expData.Type, receivedData.Type)
				assert.Equal(t, expData.Values, receivedData.Values)
			default:
				t.Errorf("Unknown data type %#v", expData)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func skewTrail(r *rand.Rand, t httpext.Trail, minCoef, maxCoef float64) httpext.Trail ***REMOVED***
	coef := minCoef + r.Float64()*(maxCoef-minCoef)
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
	return t
***REMOVED***

func TestCloudOutput(t *testing.T) ***REMOVED***
	t.Parallel()

	getTestRunner := func(minSamples int) func(t *testing.T) ***REMOVED***
		return func(t *testing.T) ***REMOVED***
			t.Parallel()
			runCloudOutputTestCase(t, minSamples)
		***REMOVED***
	***REMOVED***

	for tcNum, minSamples := range []int***REMOVED***60, 75, 100***REMOVED*** ***REMOVED***
		t.Run(fmt.Sprintf("tc%d_minSamples%d", tcNum, minSamples), getTestRunner(minSamples))
	***REMOVED***
***REMOVED***

func runCloudOutputTestCase(t *testing.T, minSamples int) ***REMOVED***
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)) //nolint:gosec
	t.Logf("Random source seeded with %d\n", seed)

	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprintf(w, `***REMOVED***
			"reference_id": "123",
			"config": ***REMOVED***
				"metricPushInterval": "10ms",
				"aggregationPeriod": "30ms",
				"aggregationCalcInterval": "40ms",
				"aggregationWaitPeriod": "5ms",
				"aggregationMinSamples": %d
			***REMOVED***
		***REMOVED***`, minSamples)
		require.NoError(t, err)
	***REMOVED***))
	defer tb.Cleanup()

	out, err := newOutput(output.Params***REMOVED***
		Logger:     testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***"host": "%s", "noCompress": true***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	assert.True(t, out.config.Host.Valid)
	assert.Equal(t, tb.ServerHTTP.URL, out.config.Host.String)
	assert.True(t, out.config.NoCompress.Valid)
	assert.True(t, out.config.NoCompress.Bool)
	assert.False(t, out.config.MetricPushInterval.Valid)
	assert.False(t, out.config.AggregationPeriod.Valid)
	assert.False(t, out.config.AggregationWaitPeriod.Valid)

	require.NoError(t, out.Start())
	assert.Equal(t, "123", out.referenceID)
	assert.True(t, out.config.MetricPushInterval.Valid)
	assert.Equal(t, types.Duration(10*time.Millisecond), out.config.MetricPushInterval.Duration)
	assert.True(t, out.config.AggregationPeriod.Valid)
	assert.Equal(t, types.Duration(30*time.Millisecond), out.config.AggregationPeriod.Duration)
	assert.True(t, out.config.AggregationWaitPeriod.Valid)
	assert.Equal(t, types.Duration(5*time.Millisecond), out.config.AggregationWaitPeriod.Duration)

	now := time.Now()
	tagMap := map[string]string***REMOVED***"test": "mest", "a": "b", "name": "name", "url": "url"***REMOVED***
	tags := stats.IntoSampleTags(&tagMap)
	expectedTagMap := tags.CloneTags()
	expectedTagMap["url"], _ = tags.Get("name")
	expectedTags := stats.IntoSampleTags(&expectedTagMap)

	expSamples := make(chan []Sample)
	defer close(expSamples)
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", out.referenceID), getSampleChecker(t, expSamples))
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/tests/%s", out.referenceID), func(rw http.ResponseWriter, _ *http.Request) ***REMOVED***
		rw.WriteHeader(http.StatusOK) // silence a test warning
	***REMOVED***)

	out.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
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
			Time:  toMicroSecond(now),
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
	out.AddMetricSamples([]stats.SampleContainer***REMOVED***&simpleTrail***REMOVED***)
	expSamples <- []Sample***REMOVED****NewSampleFromTrail(&simpleTrail)***REMOVED***

	smallSkew := 0.02

	trails := []stats.SampleContainer***REMOVED******REMOVED***
	durations := make([]time.Duration, len(trails))
	for i := int64(0); i < out.config.AggregationMinSamples.Int64; i++ ***REMOVED***
		similarTrail := skewTrail(r, simpleTrail, 1.0, 1.0+smallSkew)
		trails = append(trails, &similarTrail)
		durations = append(durations, similarTrail.Duration)
	***REMOVED***
	sort.Slice(durations, func(i, j int) bool ***REMOVED*** return durations[i] < durations[j] ***REMOVED***)
	t.Logf("Sorted durations: %#v", durations) // Useful to debug any failures, doesn't get in the way otherwise

	checkAggrMetric := func(normal time.Duration, aggr AggregatedMetric) ***REMOVED***
		assert.True(t, aggr.Min <= aggr.Avg)
		assert.True(t, aggr.Avg <= aggr.Max)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Min), smallSkew)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Avg), smallSkew)
		assert.InEpsilon(t, normal, stats.ToD(aggr.Max), smallSkew)
	***REMOVED***

	outlierTrail := skewTrail(r, simpleTrail, 2.0+smallSkew, 3.0+smallSkew)
	trails = append(trails, &outlierTrail)
	out.AddMetricSamples(trails)
	expSamples <- []Sample***REMOVED***
		*NewSampleFromTrail(&outlierTrail),
		***REMOVED***
			Type:   DataTypeAggregatedHTTPReqs,
			Metric: "http_req_li_all",
			Data: func(data interface***REMOVED******REMOVED***) ***REMOVED***
				aggrData, ok := data.(*SampleDataAggregatedHTTPReqs)
				assert.True(t, ok)
				assert.True(t, aggrData.Tags.IsEqual(expectedTags))
				assert.Equal(t, out.config.AggregationMinSamples.Int64, int64(aggrData.Count))
				assert.Equal(t, "aggregated_trend", aggrData.Type)
				assert.InDelta(t, now.UnixNano(), aggrData.Time*1000, float64(out.config.AggregationPeriod.Duration))

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

	require.NoError(t, out.Stop())
***REMOVED***

func TestCloudOutputMaxPerPacket(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	maxMetricSamplesPerPackage := 20
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
	tb.Mux.HandleFunc("/v1/tests/12", func(rw http.ResponseWriter, _ *http.Request) ***REMOVED*** rw.WriteHeader(http.StatusOK) ***REMOVED***)
	defer tb.Cleanup()

	out, err := newOutput(output.Params***REMOVED***
		Logger:     testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***"host": "%s", "noCompress": true***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)
	require.NoError(t, err)
	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)
	gotTheLimit := false
	var m sync.Mutex

	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", out.referenceID),
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

	require.NoError(t, out.Start())

	out.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	for j := time.Duration(1); j <= 200; j++ ***REMOVED***
		container := make([]stats.SampleContainer, 0, 500)
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
		out.AddMetricSamples(container)
	***REMOVED***

	require.NoError(t, out.Stop())
	require.True(t, gotTheLimit)
***REMOVED***

func TestCloudOutputStopSendingMetric(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		data := &cloudapi.TestRun***REMOVED******REMOVED***
		err = json.Unmarshal(body, &data)
		require.NoError(t, err)
		assert.Equal(t, "my-custom-name", data.Name)

		_, err = fmt.Fprint(resp, `***REMOVED***
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
	tb.Mux.HandleFunc("/v1/tests/12", func(rw http.ResponseWriter, _ *http.Request) ***REMOVED*** rw.WriteHeader(http.StatusOK) ***REMOVED***)
	defer tb.Cleanup()

	out, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***
			"host": "%s", "noCompress": true,
			"maxMetricSamplesPerPackage": 50,
			"name": "something-that-should-be-overwritten"
		***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
			External: map[string]json.RawMessage***REMOVED***
				"loadimpact": json.RawMessage(`***REMOVED***"name": "my-custom-name"***REMOVED***`),
			***REMOVED***,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)
	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)

	count := 1
	max := 5
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", out.referenceID),
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			count++
			if count == max ***REMOVED***
				type payload struct ***REMOVED***
					Error cloudapi.ErrorResponse `json:"error"`
				***REMOVED***
				res := &payload***REMOVED******REMOVED***
				res.Error = cloudapi.ErrorResponse***REMOVED***Code: 4***REMOVED***
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

	require.NoError(t, out.Start())

	out.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	for j := time.Duration(1); j <= 200; j++ ***REMOVED***
		container := make([]stats.SampleContainer, 0, 500)
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
		out.AddMetricSamples(container)
	***REMOVED***

	require.NoError(t, out.Stop())

	require.Equal(t, lib.RunStatusQueued, out.runStatus)
	select ***REMOVED***
	case <-out.stopSendingMetrics:
		// all is fine
	default:
		t.Fatal("sending metrics wasn't stopped")
	***REMOVED***
	require.Equal(t, max, count)

	nBufferSamples := len(out.bufferSamples)
	nBufferHTTPTrails := len(out.bufferHTTPTrails)
	out.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.VUs,
		Tags:   stats.NewSampleTags(tags.CloneTags()),
		Value:  1.0,
	***REMOVED******REMOVED***)
	if nBufferSamples != len(out.bufferSamples) || nBufferHTTPTrails != len(out.bufferHTTPTrails) ***REMOVED***
		t.Errorf("Output still collects data after stop sending metrics")
	***REMOVED***
***REMOVED***

func TestCloudOutputRequireScriptName(t *testing.T) ***REMOVED***
	t.Parallel()
	_, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: ""***REMOVED***,
	***REMOVED***)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "script name not set")
***REMOVED***

func TestCloudOutputAggregationPeriodZeroNoBlock(t *testing.T) ***REMOVED***
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
	tb.Mux.HandleFunc("/v1/tests/123", func(rw http.ResponseWriter, _ *http.Request) ***REMOVED*** rw.WriteHeader(http.StatusOK) ***REMOVED***)
	defer tb.Cleanup()

	out, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***
			"host": "%s", "noCompress": true,
			"maxMetricSamplesPerPackage": 50
		***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	assert.True(t, out.config.Host.Valid)
	assert.Equal(t, tb.ServerHTTP.URL, out.config.Host.String)
	assert.True(t, out.config.NoCompress.Valid)
	assert.True(t, out.config.NoCompress.Bool)
	assert.False(t, out.config.MetricPushInterval.Valid)
	assert.False(t, out.config.AggregationPeriod.Valid)
	assert.False(t, out.config.AggregationWaitPeriod.Valid)

	require.NoError(t, out.Start())
	assert.Equal(t, "123", out.referenceID)
	assert.True(t, out.config.MetricPushInterval.Valid)
	assert.Equal(t, types.Duration(10*time.Millisecond), out.config.MetricPushInterval.Duration)
	assert.True(t, out.config.AggregationPeriod.Valid)
	assert.Equal(t, types.Duration(0), out.config.AggregationPeriod.Duration)
	assert.True(t, out.config.AggregationWaitPeriod.Valid)
	assert.Equal(t, types.Duration(5*time.Millisecond), out.config.AggregationWaitPeriod.Duration)

	expSamples := make(chan []Sample)
	defer close(expSamples)
	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", out.referenceID), getSampleChecker(t, expSamples))

	require.NoError(t, out.Stop())
	require.Equal(t, lib.RunStatusQueued, out.runStatus)
***REMOVED***

func TestCloudOutputPushRefID(t *testing.T) ***REMOVED***
	t.Parallel()
	expSamples := make(chan []Sample)
	defer close(expSamples)

	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	failHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		t.Errorf("%s should not have been called at all", r.RequestURI)
	***REMOVED***)
	tb.Mux.HandleFunc("/v1/tests", failHandler)
	tb.Mux.HandleFunc("/v1/tests/333", failHandler)
	tb.Mux.HandleFunc("/v1/metrics/333", getSampleChecker(t, expSamples))

	out, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***
			"host": "%s", "noCompress": true,
			"metricPushInterval": "10ms",
			"aggregationPeriod": "0ms",
			"pushRefID": "333"
		***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	assert.Equal(t, "333", out.config.PushRefID.String)
	require.NoError(t, out.Start())
	assert.Equal(t, "333", out.referenceID)

	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b"***REMOVED***)

	out.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
		Time:   now,
		Metric: metrics.HTTPReqDuration,
		Tags:   tags,
		Value:  123.45,
	***REMOVED******REMOVED***)
	exp := []Sample***REMOVED******REMOVED***
		Type:   DataTypeSingle,
		Metric: metrics.HTTPReqDuration.Name,
		Data: &SampleDataSingle***REMOVED***
			Type:  metrics.HTTPReqDuration.Type,
			Time:  toMicroSecond(now),
			Tags:  tags,
			Value: 123.45,
		***REMOVED***,
	***REMOVED******REMOVED***

	select ***REMOVED***
	case expSamples <- exp:
	case <-time.After(5 * time.Second):
		t.Error("test timeout")
	***REMOVED***

	require.NoError(t, out.Stop())
***REMOVED***

func TestCloudOutputRecvIterLIAllIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		data := &cloudapi.TestRun***REMOVED******REMOVED***
		err = json.Unmarshal(body, &data)
		require.NoError(t, err)
		assert.Equal(t, "script.js", data.Name)

		_, err = fmt.Fprintf(resp, `***REMOVED***"reference_id": "123"***REMOVED***`)
		require.NoError(t, err)
	***REMOVED***))
	tb.Mux.HandleFunc("/v1/tests/123", func(rw http.ResponseWriter, _ *http.Request) ***REMOVED*** rw.WriteHeader(http.StatusOK) ***REMOVED***)
	defer tb.Cleanup()

	out, err := newOutput(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***
			"host": "%s", "noCompress": true,
			"maxMetricSamplesPerPackage": 50
		***REMOVED***`, tb.ServerHTTP.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "path/to/script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	gotIterations := false
	var m sync.Mutex
	expValues := map[string]float64***REMOVED***
		"data_received":      100,
		"data_sent":          200,
		"iteration_duration": 60000,
		"iterations":         1,
	***REMOVED***

	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", out.referenceID),
		func(_ http.ResponseWriter, r *http.Request) ***REMOVED***
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)

			receivedSamples := []Sample***REMOVED******REMOVED***
			assert.NoError(t, json.Unmarshal(body, &receivedSamples))

			assert.Len(t, receivedSamples, 1)
			assert.Equal(t, "iter_li_all", receivedSamples[0].Metric)
			assert.Equal(t, DataTypeMap, receivedSamples[0].Type)
			data, ok := receivedSamples[0].Data.(*SampleDataMap)
			assert.True(t, ok)
			assert.Equal(t, expValues, data.Values)

			m.Lock()
			gotIterations = true
			m.Unlock()
		***REMOVED***)

	require.NoError(t, out.Start())

	now := time.Now()
	simpleNetTrail := netext.NetTrail***REMOVED***
		BytesRead:     100,
		BytesWritten:  200,
		FullIteration: true,
		StartTime:     now.Add(-time.Minute),
		EndTime:       now,
		Samples: []stats.Sample***REMOVED***
			***REMOVED***
				Time:   now,
				Metric: metrics.DataSent,
				Value:  float64(200),
			***REMOVED***,
			***REMOVED***
				Time:   now,
				Metric: metrics.DataReceived,
				Value:  float64(100),
			***REMOVED***,
			***REMOVED***
				Time:   now,
				Metric: metrics.Iterations,
				Value:  1,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	out.AddMetricSamples([]stats.SampleContainer***REMOVED***&simpleNetTrail***REMOVED***)
	require.NoError(t, out.Stop())
	require.True(t, gotIterations)
***REMOVED***

func TestNewName(t *testing.T) ***REMOVED***
	t.Parallel()
	mustParse := func(u string) *url.URL ***REMOVED***
		result, err := url.Parse(u)
		require.NoError(t, err)
		return result
	***REMOVED***

	cases := []struct ***REMOVED***
		url      *url.URL
		expected string
	***REMOVED******REMOVED***
		***REMOVED***
			url: &url.URL***REMOVED***
				Opaque: "github.com/loadimpact/k6/samples/http_get.js",
			***REMOVED***,
			expected: "http_get.js",
		***REMOVED***,
		***REMOVED***
			url:      mustParse("http://github.com/loadimpact/k6/samples/http_get.js"),
			expected: "http_get.js",
		***REMOVED***,
		***REMOVED***
			url:      mustParse("file://home/user/k6/samples/http_get.js"),
			expected: "http_get.js",
		***REMOVED***,
		***REMOVED***
			url:      mustParse("file://C:/home/user/k6/samples/http_get.js"),
			expected: "http_get.js",
		***REMOVED***,
	***REMOVED***

	for _, testCase := range cases ***REMOVED***
		testCase := testCase

		t.Run(testCase.url.String(), func(t *testing.T) ***REMOVED***
			out, err := newOutput(output.Params***REMOVED***
				Logger: testutils.NewLogger(t),
				ScriptOptions: lib.Options***REMOVED***
					Duration:   types.NullDurationFrom(1 * time.Second),
					SystemTags: &stats.DefaultSystemTagSet,
				***REMOVED***,
				ScriptPath: testCase.url,
			***REMOVED***)
			require.NoError(t, err)
			require.Equal(t, out.config.Name.String, testCase.expected)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestPublishMetric(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		g, err := gzip.NewReader(r.Body)

		require.NoError(t, err)
		var buf bytes.Buffer
		_, err = io.Copy(&buf, g) //nolint:gosec
		require.NoError(t, err)
		byteCount, err := strconv.Atoi(r.Header.Get("x-payload-byte-count"))
		require.NoError(t, err)
		require.Equal(t, buf.Len(), byteCount)

		samplesCount, err := strconv.Atoi(r.Header.Get("x-payload-sample-count"))
		require.NoError(t, err)
		var samples []*Sample
		err = json.Unmarshal(buf.Bytes(), &samples)
		require.NoError(t, err)
		require.Equal(t, len(samples), samplesCount)

		_, err = fmt.Fprintf(w, "")
		require.NoError(t, err)
	***REMOVED***))
	defer server.Close()

	out, err := newOutput(output.Params***REMOVED***
		Logger:     testutils.NewLogger(t),
		JSONConfig: json.RawMessage(fmt.Sprintf(`***REMOVED***"host": "%s", "noCompress": false***REMOVED***`, server.URL)),
		ScriptOptions: lib.Options***REMOVED***
			Duration:   types.NullDurationFrom(1 * time.Second),
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		ScriptPath: &url.URL***REMOVED***Path: "script.js"***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	samples := []*Sample***REMOVED***
		***REMOVED***
			Type:   "Point",
			Metric: "metric",
			Data: &SampleDataSingle***REMOVED***
				Type:  1,
				Time:  toMicroSecond(time.Now()),
				Value: 1.2,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err = out.client.PushMetric("1", samples)

	assert.Nil(t, err)
***REMOVED***
