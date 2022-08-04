package influxdb

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/metrics"
)

func benchmarkInfluxdb(b *testing.B, t time.Duration) ***REMOVED***
	metric, err := metrics.NewRegistry().NewMetric("test_gauge", metrics.Gauge)
	require.NoError(b, err)

	testOutputCycle(b, func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		for ***REMOVED***
			time.Sleep(t)
			m, _ := io.CopyN(ioutil.Discard, r.Body, 1<<18) // read 1/4 mb a time
			if m == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		rw.WriteHeader(204)
	***REMOVED***, func(tb testing.TB, c *Output) ***REMOVED***
		b = tb.(*testing.B)
		b.ResetTimer()

		samples := make(metrics.Samples, 10)
		for i := 0; i < len(samples); i++ ***REMOVED***
			samples[i] = metrics.Sample***REMOVED***
				Metric: metric,
				Time:   time.Now(),
				Tags: metrics.NewSampleTags(map[string]string***REMOVED***
					"something": "else",
					"VU":        "21",
					"else":      "something",
				***REMOVED***),
				Value: 2.0,
			***REMOVED***
		***REMOVED***

		b.ResetTimer()
		for i := 0; i < b.N; i++ ***REMOVED***
			c.AddMetricSamples([]metrics.SampleContainer***REMOVED***samples***REMOVED***)
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
