package influxdb

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestBadConcurrentWrites(t *testing.T) ***REMOVED***
	c := NewConfig()
	t.Run("0", func(t *testing.T) ***REMOVED***
		c.ConcurrentWrites = null.IntFrom(0)
		_, err := New(*c)
		require.Error(t, err)
		require.Equal(t, err.Error(), "influxdb's ConcurrentWrites must be a positive number")
	***REMOVED***)

	t.Run("-2", func(t *testing.T) ***REMOVED***
		c.ConcurrentWrites = null.IntFrom(-2)
		_, err := New(*c)
		require.Error(t, err)
		require.Equal(t, err.Error(), "influxdb's ConcurrentWrites must be a positive number")
	***REMOVED***)

	t.Run("2", func(t *testing.T) ***REMOVED***
		c.ConcurrentWrites = null.IntFrom(2)
		_, err := New(*c)
		require.NoError(t, err)
	***REMOVED***)
***REMOVED***

func testCollectorCycle(t testing.TB, handler http.HandlerFunc, body func(testing.TB, *Collector)) ***REMOVED***
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

	config := NewConfig()
	config.Addr = null.StringFrom("http://" + l.Addr().String())
	c, err := New(*config)
	require.NoError(t, err)

	require.NoError(t, c.Init())
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		c.Run(ctx)
	***REMOVED***()

	body(t, c)

	cancel()
	wg.Wait()
***REMOVED***
func TestCollector(t *testing.T) ***REMOVED***
	var samplesRead int
	defer func() ***REMOVED***
		require.Equal(t, samplesRead, 20)
	***REMOVED***()
	testCollectorCycle(t, func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		var b = bytes.NewBuffer(nil)
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
	***REMOVED***, func(tb testing.TB, c *Collector) ***REMOVED***
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
		c.Collect([]stats.SampleContainer***REMOVED***samples***REMOVED***)
		c.Collect([]stats.SampleContainer***REMOVED***samples***REMOVED***)
	***REMOVED***)

***REMOVED***
