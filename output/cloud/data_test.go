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
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/stats"
)

func TestSampleMarshaling(t *testing.T) ***REMOVED***
	t.Parallel()

	now := time.Now()
	exptoMicroSecond := now.UnixNano() / 1000

	testCases := []struct ***REMOVED***
		s    *Sample
		json string
	***REMOVED******REMOVED***
		***REMOVED***
			&Sample***REMOVED***
				Type:   DataTypeSingle,
				Metric: metrics.VUs.Name,
				Data: &SampleDataSingle***REMOVED***
					Type:  metrics.VUs.Type,
					Time:  toMicroSecond(now),
					Tags:  stats.IntoSampleTags(&map[string]string***REMOVED***"aaa": "bbb", "ccc": "123"***REMOVED***),
					Value: 999,
				***REMOVED***,
			***REMOVED***,
			fmt.Sprintf(`***REMOVED***"type":"Point","metric":"vus","data":***REMOVED***"time":"%d","type":"gauge","tags":***REMOVED***"aaa":"bbb","ccc":"123"***REMOVED***,"value":999***REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
		***REMOVED***
			&Sample***REMOVED***
				Type:   DataTypeMap,
				Metric: "iter_li_all",
				Data: &SampleDataMap***REMOVED***
					Time: toMicroSecond(now),
					Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest"***REMOVED***),
					Values: map[string]float64***REMOVED***
						metrics.DataSent.Name:          1234.5,
						metrics.DataReceived.Name:      6789.1,
						metrics.IterationDuration.Name: stats.D(10 * time.Second),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			fmt.Sprintf(`***REMOVED***"type":"Points","metric":"iter_li_all","data":***REMOVED***"time":"%d","type":"counter","tags":***REMOVED***"test":"mest"***REMOVED***,"values":***REMOVED***"data_received":6789.1,"data_sent":1234.5,"iteration_duration":10000***REMOVED******REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
		***REMOVED***
			NewSampleFromTrail(&httpext.Trail***REMOVED***
				EndTime:        now,
				Duration:       123000,
				Blocked:        1000,
				Connecting:     2000,
				TLSHandshaking: 3000,
				Sending:        4000,
				Waiting:        5000,
				Receiving:      6000,
			***REMOVED***),
			fmt.Sprintf(`***REMOVED***"type":"Points","metric":"http_req_li_all","data":***REMOVED***"time":"%d","type":"counter","values":***REMOVED***"http_req_blocked":0.001,"http_req_connecting":0.002,"http_req_duration":0.123,"http_req_receiving":0.006,"http_req_sending":0.004,"http_req_tls_handshaking":0.003,"http_req_waiting":0.005,"http_reqs":1***REMOVED******REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
		***REMOVED***
			NewSampleFromTrail(&httpext.Trail***REMOVED***
				EndTime:        now,
				Duration:       123000,
				Blocked:        1000,
				Connecting:     2000,
				TLSHandshaking: 3000,
				Sending:        4000,
				Waiting:        5000,
				Receiving:      6000,
				Failed:         null.NewBool(false, true),
			***REMOVED***),
			fmt.Sprintf(`***REMOVED***"type":"Points","metric":"http_req_li_all","data":***REMOVED***"time":"%d","type":"counter","values":***REMOVED***"http_req_blocked":0.001,"http_req_connecting":0.002,"http_req_duration":0.123,"http_req_failed":0,"http_req_receiving":0.006,"http_req_sending":0.004,"http_req_tls_handshaking":0.003,"http_req_waiting":0.005,"http_reqs":1***REMOVED******REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
		***REMOVED***
			func() *Sample ***REMOVED***
				aggrData := &SampleDataAggregatedHTTPReqs***REMOVED***
					Time: exptoMicroSecond,
					Type: "aggregated_trend",
					Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest"***REMOVED***),
				***REMOVED***
				aggrData.Add(
					&httpext.Trail***REMOVED***
						EndTime:        now,
						Duration:       123000,
						Blocked:        1000,
						Connecting:     2000,
						TLSHandshaking: 3000,
						Sending:        4000,
						Waiting:        5000,
						Receiving:      6000,
					***REMOVED***,
				)

				aggrData.Add(
					&httpext.Trail***REMOVED***
						EndTime:        now,
						Duration:       13000,
						Blocked:        3000,
						Connecting:     1000,
						TLSHandshaking: 4000,
						Sending:        5000,
						Waiting:        8000,
						Receiving:      8000,
					***REMOVED***,
				)
				aggrData.CalcAverages()

				return &Sample***REMOVED***
					Type:   DataTypeAggregatedHTTPReqs,
					Metric: "http_req_li_all",
					Data:   aggrData,
				***REMOVED***
			***REMOVED***(),
			fmt.Sprintf(`***REMOVED***"type":"AggregatedPoints","metric":"http_req_li_all","data":***REMOVED***"time":"%d","type":"aggregated_trend","count":2,"tags":***REMOVED***"test":"mest"***REMOVED***,"values":***REMOVED***"http_req_duration":***REMOVED***"min":0.013,"max":0.123,"avg":0.068***REMOVED***,"http_req_blocked":***REMOVED***"min":0.001,"max":0.003,"avg":0.002***REMOVED***,"http_req_connecting":***REMOVED***"min":0.001,"max":0.002,"avg":0.0015***REMOVED***,"http_req_tls_handshaking":***REMOVED***"min":0.003,"max":0.004,"avg":0.0035***REMOVED***,"http_req_sending":***REMOVED***"min":0.004,"max":0.005,"avg":0.0045***REMOVED***,"http_req_waiting":***REMOVED***"min":0.005,"max":0.008,"avg":0.0065***REMOVED***,"http_req_receiving":***REMOVED***"min":0.006,"max":0.008,"avg":0.007***REMOVED******REMOVED******REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
		***REMOVED***
			func() *Sample ***REMOVED***
				aggrData := &SampleDataAggregatedHTTPReqs***REMOVED***
					Time: exptoMicroSecond,
					Type: "aggregated_trend",
					Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest"***REMOVED***),
				***REMOVED***
				aggrData.Add(
					&httpext.Trail***REMOVED***
						EndTime:        now,
						Duration:       123000,
						Blocked:        1000,
						Connecting:     2000,
						TLSHandshaking: 3000,
						Sending:        4000,
						Waiting:        5000,
						Receiving:      6000,
						Failed:         null.BoolFrom(false),
					***REMOVED***,
				)

				aggrData.Add(
					&httpext.Trail***REMOVED***
						EndTime:        now,
						Duration:       13000,
						Blocked:        3000,
						Connecting:     1000,
						TLSHandshaking: 4000,
						Sending:        5000,
						Waiting:        8000,
						Receiving:      8000,
					***REMOVED***,
				)
				aggrData.CalcAverages()

				return &Sample***REMOVED***
					Type:   DataTypeAggregatedHTTPReqs,
					Metric: "http_req_li_all",
					Data:   aggrData,
				***REMOVED***
			***REMOVED***(),
			fmt.Sprintf(`***REMOVED***"type":"AggregatedPoints","metric":"http_req_li_all","data":***REMOVED***"time":"%d","type":"aggregated_trend","count":2,"tags":***REMOVED***"test":"mest"***REMOVED***,"values":***REMOVED***"http_req_duration":***REMOVED***"min":0.013,"max":0.123,"avg":0.068***REMOVED***,"http_req_blocked":***REMOVED***"min":0.001,"max":0.003,"avg":0.002***REMOVED***,"http_req_connecting":***REMOVED***"min":0.001,"max":0.002,"avg":0.0015***REMOVED***,"http_req_tls_handshaking":***REMOVED***"min":0.003,"max":0.004,"avg":0.0035***REMOVED***,"http_req_sending":***REMOVED***"min":0.004,"max":0.005,"avg":0.0045***REMOVED***,"http_req_waiting":***REMOVED***"min":0.005,"max":0.008,"avg":0.0065***REMOVED***,"http_req_receiving":***REMOVED***"min":0.006,"max":0.008,"avg":0.007***REMOVED***,"http_req_failed":***REMOVED***"count":1,"nz_count":0***REMOVED******REMOVED******REMOVED******REMOVED***`, exptoMicroSecond),
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		sJSON, err := easyjson.Marshal(tc.s)
		if !assert.NoError(t, err) ***REMOVED***
			continue
		***REMOVED***
		t.Logf(string(sJSON))
		assert.JSONEq(t, tc.json, string(sJSON))

		var newS Sample
		assert.NoError(t, json.Unmarshal(sJSON, &newS))
		assert.Equal(t, tc.s.Type, newS.Type)
		assert.Equal(t, tc.s.Metric, newS.Metric)
		assert.IsType(t, tc.s.Data, newS.Data)
		// Cannot directly compare tc.s.Data and newS.Data (because of internal time.Time and SampleTags fields)
		newJSON, err := easyjson.Marshal(newS)
		assert.NoError(t, err)
		assert.JSONEq(t, string(sJSON), string(newJSON))
	***REMOVED***
***REMOVED***

func TestMetricAggregation(t *testing.T) ***REMOVED***
	m := AggregatedMetric***REMOVED******REMOVED***
	m.Add(1 * time.Second)
	m.Add(1 * time.Second)
	m.Add(3 * time.Second)
	m.Add(5 * time.Second)
	m.Add(10 * time.Second)
	m.Calc(5)
	assert.Equal(t, m.Min, stats.D(1*time.Second))
	assert.Equal(t, m.Max, stats.D(10*time.Second))
	assert.Equal(t, m.Avg, stats.D(4*time.Second))
***REMOVED***

// For more realistic request time distributions, import
// "gonum.org/v1/gonum/stat/distuv" and use something like this:
//
// randSrc := rand.NewSource(uint64(time.Now().UnixNano()))
// dist := distuv.LogNormal***REMOVED***Mu: 0, Sigma: 0.5, Src: randSrc***REMOVED***
//
// then set the data elements to time.Duration(dist.Rand() * multiplier)
//
// I've not used that after the initial tests because it's a big
// external dependency that's not really needed for the tests at
// this point.
func getDurations(r *rand.Rand, count int, min, multiplier float64) durations ***REMOVED***
	data := make(durations, count)
	for j := 0; j < count; j++ ***REMOVED***
		data[j] = time.Duration(min + r.Float64()*multiplier) //nolint:gosec
	***REMOVED***
	return data
***REMOVED***

func BenchmarkDurationBounds(b *testing.B) ***REMOVED***
	iqrRadius := 0.25 // If it's something different, the Q in IQR won't make much sense...
	iqrLowerCoef := 1.5
	iqrUpperCoef := 1.5

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)) //nolint:gosec
	b.Logf("Random source seeded with %d\n", seed)

	getData := func(b *testing.B, count int) durations ***REMOVED***
		b.StopTimer()
		defer b.StartTimer()
		return getDurations(r, count, 0.1*float64(time.Second), float64(time.Second))
	***REMOVED***

	for count := 100; count <= 5000; count += 500 ***REMOVED***
		count := count
		b.Run(fmt.Sprintf("Sort-no-interp-%d-elements", count), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				data := getData(b, count)
				data.SortGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef, false)
			***REMOVED***
		***REMOVED***)
		b.Run(fmt.Sprintf("Sort-with-interp-%d-elements", count), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				data := getData(b, count)
				data.SortGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef, true)
			***REMOVED***
		***REMOVED***)
		b.Run(fmt.Sprintf("Select-%d-elements", count), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				data := getData(b, count)
				data.SelectGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestQuickSelectAndBounds(t *testing.T) ***REMOVED***
	t.Parallel()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)) //nolint:gosec
	t.Logf("Random source seeded with %d\n", seed)

	mult := time.Millisecond
	for _, count := range []int***REMOVED***1, 2, 3, 4, 5, 10, 15, 20, 25, 50, 100, 250 + r.Intn(100)***REMOVED*** ***REMOVED***
		count := count
		t.Run(fmt.Sprintf("simple-%d", count), func(t *testing.T) ***REMOVED***
			data := make(durations, count)
			for i := 0; i < count; i++ ***REMOVED***
				data[i] = time.Duration(i) * mult
			***REMOVED***
			rand.Shuffle(len(data), data.Swap)
			for i := 0; i < 10; i++ ***REMOVED***
				dataCopy := make(durations, count)
				assert.Equal(t, count, copy(dataCopy, data))
				k := r.Intn(count)
				assert.Equal(t, dataCopy.quickSelect(k), time.Duration(k)*mult)
			***REMOVED***
		***REMOVED***)
		t.Run(fmt.Sprintf("random-%d", count), func(t *testing.T) ***REMOVED***
			testCases := []struct***REMOVED*** r, l, u float64 ***REMOVED******REMOVED***
				***REMOVED***0.25, 1.5, 1.5***REMOVED***, // Textbook
				***REMOVED***0.25, 1.5, 1.3***REMOVED***, // Defaults
				***REMOVED***0.1, 0.5, 0.3***REMOVED***,  // Extreme narrow
				***REMOVED***0.3, 2, 1.8***REMOVED***,    // Extreme wide
			***REMOVED***

			for tcNum, tc := range testCases ***REMOVED***
				tc := tc
				data := getDurations(r, count, 0.3*float64(time.Second), 2*float64(time.Second))
				dataForSort := make(durations, count)
				dataForSelect := make(durations, count)
				assert.Equal(t, count, copy(dataForSort, data))
				assert.Equal(t, count, copy(dataForSelect, data))
				assert.Equal(t, dataForSort, dataForSelect)

				t.Run(fmt.Sprintf("bounds-tc%d", tcNum), func(t *testing.T) ***REMOVED***
					sortMin, sortMax := dataForSort.SortGetNormalBounds(tc.r, tc.l, tc.u, false)
					selectMin, selectMax := dataForSelect.SelectGetNormalBounds(tc.r, tc.l, tc.u)
					assert.Equal(t, sortMin, selectMin)
					assert.Equal(t, sortMax, selectMax)

					k := r.Intn(count)
					assert.Equal(t, dataForSort[k], dataForSelect.quickSelect(k))
					assert.Equal(t, dataForSort[k], data.quickSelect(k))
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSortInterpolation(t *testing.T) ***REMOVED***
	t.Parallel()

	// Super contrived example to make the checks easy - 11 values from 0 to 10 seconds inclusive
	count := 11
	data := make(durations, count)
	for i := 0; i < count; i++ ***REMOVED***
		data[i] = time.Duration(i) * time.Second
	***REMOVED***

	min, max := data.SortGetNormalBounds(0.25, 1, 1, true)
	// Expected values: Q1=2.5, Q3=7.5, IQR=5, so with 1 for coefficients we can expect min=-2,5, max=12.5 seconds
	assert.Equal(t, min, -2500*time.Millisecond)
	assert.Equal(t, max, 12500*time.Millisecond)
***REMOVED***
