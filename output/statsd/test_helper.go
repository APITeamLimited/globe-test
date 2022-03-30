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

package statsd

import (
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/metrics"
)

type getOutputFn func(
	logger logrus.FieldLogger,
	addr, namespace null.String,
	bufferSize null.Int,
	pushInterval types.NullDuration,
) (*Output, error)

//nolint:funlen
func baseTest(t *testing.T,
	getOutput getOutputFn,
	checkResult func(t *testing.T, samples []metrics.SampleContainer, expectedOutput, output string),
) ***REMOVED***
	t.Helper()
	testNamespace := "testing.things." // to be dynamic

	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	require.NoError(t, err)
	listener, err := net.ListenUDP("udp", addr) // we want to listen on a random port
	require.NoError(t, err)
	ch := make(chan string, 20)
	end := make(chan struct***REMOVED******REMOVED***)
	defer close(end)

	go func() ***REMOVED***
		defer close(ch)
		var buf [4096]byte
		for ***REMOVED***
			select ***REMOVED***
			case <-end:
				return
			default:
				n, _, err := listener.ReadFromUDP(buf[:])
				require.NoError(t, err)
				ch <- string(buf[:n])
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	pushInterval := types.NullDurationFrom(time.Millisecond * 10)
	collector, err := getOutput(
		testutils.NewLogger(t),
		null.StringFrom(listener.LocalAddr().String()),
		null.StringFrom(testNamespace),
		null.IntFrom(5),
		pushInterval,
	)
	require.NoError(t, err)
	require.NoError(t, collector.Start())
	defer func() ***REMOVED***
		require.NoError(t, collector.Stop())
	***REMOVED***()
	newSample := func(m *metrics.Metric, value float64, tags map[string]string) metrics.Sample ***REMOVED***
		return metrics.Sample***REMOVED***
			Time:   time.Now(),
			Metric: m, Value: value, Tags: metrics.IntoSampleTags(&tags),
		***REMOVED***
	***REMOVED***

	registry := metrics.NewRegistry()
	myCounter, err := registry.NewMetric("my_counter", metrics.Counter)
	require.NoError(t, err)
	myGauge, err := registry.NewMetric("my_gauge", metrics.Gauge)
	require.NoError(t, err)
	myTrend, err := registry.NewMetric("my_trend", metrics.Trend)
	require.NoError(t, err)
	myRate, err := registry.NewMetric("my_rate", metrics.Rate)
	require.NoError(t, err)
	myCheck, err := registry.NewMetric("my_check", metrics.Rate)
	require.NoError(t, err)

	testMatrix := []struct ***REMOVED***
		input  []metrics.SampleContainer
		output string
	***REMOVED******REMOVED***
		***REMOVED***
			input: []metrics.SampleContainer***REMOVED***
				newSample(myCounter, 12, map[string]string***REMOVED***
					"tag1": "value1",
					"tag3": "value3",
				***REMOVED***),
			***REMOVED***,
			output: "testing.things.my_counter:12|c",
		***REMOVED***,
		***REMOVED***
			input: []metrics.SampleContainer***REMOVED***
				newSample(myGauge, 13, map[string]string***REMOVED***
					"tag1": "value1",
					"tag3": "value3",
				***REMOVED***),
			***REMOVED***,
			output: "testing.things.my_gauge:13.000000|g",
		***REMOVED***,
		***REMOVED***
			input: []metrics.SampleContainer***REMOVED***
				newSample(myTrend, 14, map[string]string***REMOVED***
					"tag1": "value1",
					"tag3": "value3",
				***REMOVED***),
			***REMOVED***,
			output: "testing.things.my_trend:14.000000|ms",
		***REMOVED***,
		***REMOVED***
			input: []metrics.SampleContainer***REMOVED***
				newSample(myRate, 15, map[string]string***REMOVED***
					"tag1": "value1",
					"tag3": "value3",
				***REMOVED***),
			***REMOVED***,
			output: "testing.things.my_rate:15|c",
		***REMOVED***,
		***REMOVED***
			input: []metrics.SampleContainer***REMOVED***
				newSample(myCheck, 16, map[string]string***REMOVED***
					"tag1":  "value1",
					"tag3":  "value3",
					"check": "max<100",
				***REMOVED***),
				newSample(myCheck, 0, map[string]string***REMOVED***
					"tag1":  "value1",
					"tag3":  "value3",
					"check": "max>100",
				***REMOVED***),
			***REMOVED***,
			output: "testing.things.check.max<100.pass:1|c\ntesting.things.check.max>100.fail:1|c",
		***REMOVED***,
	***REMOVED***
	for _, test := range testMatrix ***REMOVED***
		collector.AddMetricSamples(test.input)
		time.Sleep((time.Duration)(pushInterval.Duration))
		output := <-ch
		checkResult(t, test.input, test.output, output)
	***REMOVED***
***REMOVED***
