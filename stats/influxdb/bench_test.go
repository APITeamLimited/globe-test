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
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"
)

func benchmarkInfluxdb(b *testing.B, t time.Duration) ***REMOVED***
	testCollectorCycle(b, func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		for ***REMOVED***
			time.Sleep(t)
			m, _ := io.CopyN(ioutil.Discard, r.Body, 1<<18) // read 1/4 mb a time
			if m == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		rw.WriteHeader(204)
	***REMOVED***, func(tb testing.TB, c *Collector) ***REMOVED***
		b = tb.(*testing.B)
		b.ResetTimer()

		var samples = make(stats.Samples, 10)
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

		b.ResetTimer()
		for i := 0; i < b.N; i++ ***REMOVED***
			c.Collect([]stats.SampleContainer***REMOVED***samples***REMOVED***)
			time.Sleep(time.Nanosecond * 20)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkInfluxdb1Second(b *testing.B) ***REMOVED***
	benchmarkInfluxdb(b, time.Second)
***REMOVED***

func BenchmarkInfluxdb2Second(b *testing.B) ***REMOVED***
	benchmarkInfluxdb(b, 2*time.Second)
***REMOVED***

func BenchmarkInfluxdb100Milliseconds(b *testing.B) ***REMOVED***
	benchmarkInfluxdb(b, 100*time.Millisecond)
***REMOVED***
