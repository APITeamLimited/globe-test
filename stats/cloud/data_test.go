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

	"github.com/loadimpact/k6/lib/netext"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampMarshaling(t *testing.T) ***REMOVED***
	t.Parallel()

	oldTimeFormat, err := time.Parse(
		time.RFC3339,
		//1521806137415652223 as a unix nanosecond timestamp
		"2018-03-23T13:55:37.415652223+02:00",
	)
	require.NoError(t, err)

	testCases := []struct ***REMOVED***
		t   time.Time
		exp string
	***REMOVED******REMOVED***
		***REMOVED***oldTimeFormat, `"1521806137415652"`***REMOVED***,
		***REMOVED***time.Unix(1521806137, 415652223), `"1521806137415652"`***REMOVED***,
		***REMOVED***time.Unix(1521806137, 0), `"1521806137000000"`***REMOVED***,
		***REMOVED***time.Unix(0, 0), `"0"`***REMOVED***,
		***REMOVED***time.Unix(0, 1), `"0"`***REMOVED***,
		***REMOVED***time.Unix(0, 1000), `"1"`***REMOVED***,
		***REMOVED***time.Unix(1, 0), `"1000000"`***REMOVED***,
	***REMOVED***

	for i, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) ***REMOVED***
			res, err := json.Marshal(Timestamp(tc.t))
			require.NoError(t, err)
			assert.Equal(t, string(res), tc.exp)

			var rev Timestamp
			require.NoError(t, json.Unmarshal(res, &rev))

			assert.Truef(
				t,
				rev.Equal(Timestamp(tc.t)),
				"Expected the difference to be under a microsecond, but is %s (%d and %d)",
				tc.t.Sub(time.Time(rev)),
				tc.t.UnixNano(),
				time.Time(rev).UnixNano(),
			)

			assert.False(t, Timestamp(time.Now()).Equal(Timestamp(tc.t)))
		***REMOVED***)
	***REMOVED***

	var expErr Timestamp
	assert.Error(t, json.Unmarshal([]byte(`1234`), &expErr))
	assert.Error(t, json.Unmarshal([]byte(`"1234a"`), &expErr))
***REMOVED***

func TestSampleMarshaling(t *testing.T) ***REMOVED***
	t.Parallel()

	now := time.Now()
	expTimestamp := now.UnixNano() / 1000

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
					Time:  Timestamp(now),
					Tags:  stats.IntoSampleTags(&map[string]string***REMOVED***"aaa": "bbb", "ccc": "123"***REMOVED***),
					Value: 999,
				***REMOVED***,
			***REMOVED***,
			fmt.Sprintf(`***REMOVED***"type":"Point","metric":"vus","data":***REMOVED***"time":"%d","type":"gauge","tags":***REMOVED***"aaa":"bbb","ccc":"123"***REMOVED***,"value":999***REMOVED******REMOVED***`, expTimestamp),
		***REMOVED***,
		***REMOVED***
			&Sample***REMOVED***
				Type:   DataTypeMap,
				Metric: "iter_li_all",
				Data: &SampleDataMap***REMOVED***
					Time: Timestamp(now),
					Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest"***REMOVED***),
					Values: map[string]float64***REMOVED***
						metrics.DataSent.Name:          1234.5,
						metrics.DataReceived.Name:      6789.1,
						metrics.IterationDuration.Name: stats.D(10 * time.Second),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			fmt.Sprintf(`***REMOVED***"type":"Points","metric":"iter_li_all","data":***REMOVED***"time":"%d","type":"counter","tags":***REMOVED***"test":"mest"***REMOVED***,"values":***REMOVED***"data_received":6789.1,"data_sent":1234.5,"iteration_duration":10000***REMOVED******REMOVED******REMOVED***`, expTimestamp),
		***REMOVED***,
		***REMOVED***
			NewSampleFromTrail(&netext.Trail***REMOVED***
				EndTime:        now,
				Duration:       123000,
				Blocked:        1000,
				Connecting:     2000,
				TLSHandshaking: 3000,
				Sending:        4000,
				Waiting:        5000,
				Receiving:      6000,
			***REMOVED***),
			fmt.Sprintf(`***REMOVED***"type":"Points","metric":"http_req_li_all","data":***REMOVED***"time":"%d","type":"counter","values":***REMOVED***"http_req_blocked":0.001,"http_req_connecting":0.002,"http_req_duration":0.123,"http_req_receiving":0.006,"http_req_sending":0.004,"http_req_tls_handshaking":0.003,"http_req_waiting":0.005,"http_reqs":1***REMOVED******REMOVED******REMOVED***`, expTimestamp),
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		sJSON, err := json.Marshal(tc.s)
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
		newJSON, err := json.Marshal(newS)
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
func getDurations(count int, min, multiplier float64) durations ***REMOVED***
	data := make(durations, count)
	for j := 0; j < count; j++ ***REMOVED***
		data[j] = time.Duration(min + rand.Float64()*multiplier)
	***REMOVED***
	return data
***REMOVED***
func BenchmarkDurationBounds(b *testing.B) ***REMOVED***
	iqrRadius := 0.25 // If it's something different, the Q in IQR won't make much sense...
	iqrLowerCoef := 1.5
	iqrUpperCoef := 1.5

	getData := func(b *testing.B, count int) durations ***REMOVED***
		b.StopTimer()
		defer b.StartTimer()
		return getDurations(count, 0.1*float64(time.Second), float64(time.Second))
	***REMOVED***

	for count := 100; count <= 5000; count += 500 ***REMOVED***
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
	mult := time.Millisecond
	for _, count := range []int***REMOVED***1, 2, 3, 4, 5, 10, 15, 20, 25, 50, 100, 250 + rand.Intn(100)***REMOVED*** ***REMOVED***
		count := count
		t.Run(fmt.Sprintf("simple-%d", count), func(t *testing.T) ***REMOVED***
			t.Parallel()
			data := make(durations, count)
			for i := 0; i < count; i++ ***REMOVED***
				data[i] = time.Duration(i) * mult
			***REMOVED***
			rand.Shuffle(len(data), data.Swap)
			for i := 0; i < 10; i++ ***REMOVED***
				dataCopy := make(durations, count)
				assert.Equal(t, count, copy(dataCopy, data))
				k := rand.Intn(count)
				assert.Equal(t, dataCopy.quickSelect(k), time.Duration(k)*mult)
			***REMOVED***
		***REMOVED***)
		t.Run(fmt.Sprintf("random-%d", count), func(t *testing.T) ***REMOVED***
			t.Parallel()

			testCases := []struct***REMOVED*** r, l, u float64 ***REMOVED******REMOVED***
				***REMOVED***0.25, 1.5, 1.5***REMOVED***, // Textbook
				***REMOVED***0.25, 1.5, 1.3***REMOVED***, // Defaults
				***REMOVED***0.1, 0.5, 0.3***REMOVED***,  // Extreme narrow
				***REMOVED***0.3, 2, 1.8***REMOVED***,    // Extreme wide
			***REMOVED***

			for tcNum, tc := range testCases ***REMOVED***
				tc := tc
				data := getDurations(count, 0.3*float64(time.Second), 2*float64(time.Second))
				dataForSort := make(durations, count)
				dataForSelect := make(durations, count)
				assert.Equal(t, count, copy(dataForSort, data))
				assert.Equal(t, count, copy(dataForSelect, data))
				assert.Equal(t, dataForSort, dataForSelect)

				t.Run(fmt.Sprintf("bounds-tc%d", tcNum), func(t *testing.T) ***REMOVED***
					t.Parallel()
					sortMin, sortMax := dataForSort.SortGetNormalBounds(tc.r, tc.l, tc.u, false)
					selectMin, selectMax := dataForSelect.SelectGetNormalBounds(tc.r, tc.l, tc.u)
					assert.Equal(t, sortMin, selectMin)
					assert.Equal(t, sortMax, selectMax)

					k := rand.Intn(count)
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
