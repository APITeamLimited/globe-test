/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package output

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/metrics"
	"go.k6.io/k6/stats"
)

func TestSampleBufferBasics(t *testing.T) ***REMOVED***
	t.Parallel()

	registry := metrics.NewRegistry()
	metric, err := registry.NewMetric("my_metric", stats.Rate)
	require.NoError(t, err)

	single := stats.Sample***REMOVED***
		Time:   time.Now(),
		Metric: metric,
		Value:  float64(123),
		Tags:   stats.NewSampleTags(map[string]string***REMOVED***"tag1": "val1"***REMOVED***),
	***REMOVED***
	connected := stats.ConnectedSamples***REMOVED***Samples: []stats.Sample***REMOVED***single, single***REMOVED***, Time: single.Time***REMOVED***
	buffer := SampleBuffer***REMOVED******REMOVED***

	assert.Empty(t, buffer.GetBufferedSamples())
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***single, single***REMOVED***)
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***single, connected, single***REMOVED***)
	assert.Equal(t, []stats.SampleContainer***REMOVED***single, single, single, connected, single***REMOVED***, buffer.GetBufferedSamples())
	assert.Empty(t, buffer.GetBufferedSamples())

	// Verify some internals
	assert.Equal(t, cap(buffer.buffer), 5)
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***single, connected***REMOVED***)
	buffer.AddMetricSamples(nil)
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED******REMOVED***)
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***single***REMOVED***)
	assert.Equal(t, []stats.SampleContainer***REMOVED***single, connected, single***REMOVED***, buffer.GetBufferedSamples())
	assert.Equal(t, cap(buffer.buffer), 4)
	buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***single***REMOVED***)
	assert.Equal(t, []stats.SampleContainer***REMOVED***single***REMOVED***, buffer.GetBufferedSamples())
	assert.Equal(t, cap(buffer.buffer), 3)
	assert.Empty(t, buffer.GetBufferedSamples())
***REMOVED***

func TestSampleBufferConcurrently(t *testing.T) ***REMOVED***
	t.Parallel()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)) //nolint:gosec
	t.Logf("Random source seeded with %d\n", seed)

	registry := metrics.NewRegistry()
	metric, err := registry.NewMetric("my_metric", stats.Gauge)
	require.NoError(t, err)

	producersCount := 50 + r.Intn(50)
	sampleCount := 10 + r.Intn(10)
	sleepModifier := 10 + r.Intn(10)
	buffer := SampleBuffer***REMOVED******REMOVED***

	wg := make(chan struct***REMOVED******REMOVED***)
	fillBuffer := func() ***REMOVED***
		for i := 0; i < sampleCount; i++ ***REMOVED***
			buffer.AddMetricSamples([]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: metric,
				Value:  float64(i),
				Tags:   stats.NewSampleTags(map[string]string***REMOVED***"tag1": "val1"***REMOVED***),
			***REMOVED******REMOVED***)
			time.Sleep(time.Duration(i*sleepModifier) * time.Microsecond)
		***REMOVED***
		wg <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	for i := 0; i < producersCount; i++ ***REMOVED***
		go fillBuffer()
	***REMOVED***

	timer := time.NewTicker(5 * time.Millisecond)
	timeout := time.After(5 * time.Second)
	defer timer.Stop()
	readSamples := make([]stats.SampleContainer, 0, sampleCount*producersCount)
	finishedProducers := 0
loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-timer.C:
			readSamples = append(readSamples, buffer.GetBufferedSamples()...)
		case <-wg:
			finishedProducers++
			if finishedProducers == producersCount ***REMOVED***
				readSamples = append(readSamples, buffer.GetBufferedSamples()...)
				break loop
			***REMOVED***
		case <-timeout:
			t.Fatalf("test timed out")
		***REMOVED***
	***REMOVED***
	assert.Equal(t, sampleCount*producersCount, len(readSamples))
	for _, s := range readSamples ***REMOVED***
		require.NotNil(t, s)
		ss := s.GetSamples()
		require.Len(t, ss, 1)
		assert.Equal(t, "my_metric", ss[0].Metric.Name)
	***REMOVED***
***REMOVED***

func TestPeriodicFlusherBasics(t *testing.T) ***REMOVED***
	t.Parallel()

	f, err := NewPeriodicFlusher(-1*time.Second, func() ***REMOVED******REMOVED***)
	assert.Error(t, err)
	assert.Nil(t, f)
	f, err = NewPeriodicFlusher(0, func() ***REMOVED******REMOVED***)
	assert.Error(t, err)
	assert.Nil(t, f)

	count := 0
	wg := &sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	f, err = NewPeriodicFlusher(100*time.Millisecond, func() ***REMOVED***
		count++
		if count == 2 ***REMOVED***
			wg.Done()
		***REMOVED***
	***REMOVED***)
	assert.NotNil(t, f)
	assert.Nil(t, err)
	wg.Wait()
	f.Stop()
	assert.Equal(t, 3, count)
***REMOVED***

func TestPeriodicFlusherConcurrency(t *testing.T) ***REMOVED***
	t.Parallel()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)) //nolint:gosec
	randStops := 10 + r.Intn(10)
	t.Logf("Random source seeded with %d\n", seed)

	count := 0
	wg := &sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	f, err := NewPeriodicFlusher(1000*time.Microsecond, func() ***REMOVED***
		// Sleep intentionally may be longer than the flush period. Also, this
		// should never happen concurrently, so it's intentionally not locked.
		time.Sleep(time.Duration(700+r.Intn(1000)) * time.Microsecond)
		count++
		if count == 100 ***REMOVED***
			wg.Done()
		***REMOVED***
	***REMOVED***)
	assert.NotNil(t, f)
	assert.Nil(t, err)
	wg.Wait()

	stopWG := &sync.WaitGroup***REMOVED******REMOVED***
	stopWG.Add(randStops)
	for i := 0; i < randStops; i++ ***REMOVED***
		go func() ***REMOVED***
			f.Stop()
			stopWG.Done()
		***REMOVED***()
	***REMOVED***
	stopWG.Wait()
	assert.True(t, count >= 101) // due to the short intervals, we might not get exactly 101
***REMOVED***
