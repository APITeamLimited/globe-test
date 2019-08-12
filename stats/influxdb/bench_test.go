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
