/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package core

import (
	"context"
	"testing"
	"time"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

type testErrorWithString string

func (e testErrorWithString) Error() string  ***REMOVED*** return string(e) ***REMOVED***
func (e testErrorWithString) String() string ***REMOVED*** return string(e) ***REMOVED***

// Apply a null logger to the engine and return the hook.
func applyNullLogger(e *Engine) *logtest.Hook ***REMOVED***
	logger, hook := logtest.NewNullLogger()
	e.SetLogger(logger)
	return hook
***REMOVED***

// Wrapper around newEngine that applies a null logger.
func newTestEngine(ex lib.Executor, opts lib.Options) (*Engine, error, *logtest.Hook) ***REMOVED***
	e, err := NewEngine(ex, opts)
	if err != nil ***REMOVED***
		return e, err, nil
	***REMOVED***
	hook := applyNullLogger(e)
	return e, nil, hook
***REMOVED***

func L(r lib.Runner) lib.Executor ***REMOVED***
	return local.New(r)
***REMOVED***

func LF(fn func(ctx context.Context) ([]stats.Sample, error)) lib.Executor ***REMOVED***
	return L(lib.RunnerFunc(fn))
***REMOVED***

func TestNewEngine(t *testing.T) ***REMOVED***
	_, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Duration: lib.NullDurationFrom(10 * time.Second),
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second)***REMOVED***)
		***REMOVED***
		assert.Equal(t, lib.NullDurationFrom(10*time.Second), e.Executor.GetEndTime())

		t.Run("Infinite", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Duration: lib.NullDurationFrom(0)***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, []lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED***, e.Stages)
			assert.Equal(t, lib.NullDuration***REMOVED******REMOVED***, e.Executor.GetEndTime())
		***REMOVED***)
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Stages: []lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
		assert.Equal(t, lib.NullDurationFrom(10*time.Second), e.Executor.GetEndTime())
	***REMOVED***)
	t.Run("Stages/Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Duration: lib.NullDurationFrom(60 * time.Second),
			Stages: []lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
		assert.Equal(t, lib.NullDurationFrom(10*time.Second), e.Executor.GetEndTime())
	***REMOVED***)
	t.Run("Iterations", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Iterations: null.IntFrom(100)***REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, null.IntFrom(100), e.Executor.GetEndIterations())
	***REMOVED***)
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.Executor.GetVUsMax())
			assert.Equal(t, int64(0), e.Executor.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(0), e.Executor.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "can't raise vu count (to 10) above vu cap (0)")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "can't raise vu count (to 10) above vu cap (1)")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(1), e.Executor.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(10), e.Executor.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.Executor.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.Executor.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			assert.NoError(t, err)
			assert.True(t, e.Executor.IsPaused())
		***REMOVED***)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric": ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		assert.Contains(t, e.thresholds, "my_metric")

		t.Run("submetrics", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Thresholds: map[string]stats.Thresholds***REMOVED***
					"my_metric***REMOVED***tag:value***REMOVED***": ***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***)
			assert.NoError(t, err)
			assert.Contains(t, e.thresholds, "my_metric***REMOVED***tag:value***REMOVED***")
			assert.Contains(t, e.submetrics, "my_metric")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		duration := 100 * time.Millisecond
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		startTime := time.Now()
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("exits with iterations", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			VUs:        null.IntFrom(10),
			VUsMax:     null.IntFrom(10),
			Iterations: null.IntFrom(100),
		***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.Run(context.Background()))
		assert.Equal(t, int64(100), e.Executor.GetIterations())
	***REMOVED***)
	t.Run("exits with duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			VUs:      null.IntFrom(10),
			VUsMax:   null.IntFrom(10),
			Duration: lib.NullDurationFrom(1 * time.Second),
		***REMOVED***)
		assert.NoError(t, err)
		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background()))
		assert.True(t, time.Now().After(startTime.Add(1*time.Second)))
	***REMOVED***)
	t.Run("exits with stages", func(t *testing.T) ***REMOVED***
		testdata := map[string]struct ***REMOVED***
			Duration time.Duration
			Stages   []lib.Stage
		***REMOVED******REMOVED***
			"none": ***REMOVED******REMOVED***,
			"one": ***REMOVED***
				1 * time.Second,
				[]lib.Stage***REMOVED******REMOVED***Duration: lib.NullDurationFrom(1 * time.Second)***REMOVED******REMOVED***,
			***REMOVED***,
			"two": ***REMOVED***
				2 * time.Second,
				[]lib.Stage***REMOVED***
					***REMOVED***Duration: lib.NullDurationFrom(1 * time.Second)***REMOVED***,
					***REMOVED***Duration: lib.NullDurationFrom(1 * time.Second)***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			"two/targeted": ***REMOVED***
				2 * time.Second,
				[]lib.Stage***REMOVED***
					***REMOVED***Duration: lib.NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
					***REMOVED***Duration: lib.NullDurationFrom(1 * time.Second), Target: null.IntFrom(10)***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
					VUs:    null.IntFrom(10),
					VUsMax: null.IntFrom(10),
				***REMOVED***)
				assert.NoError(t, err)

				e.Stages = data.Stages
				startTime := time.Now()
				assert.NoError(t, e.Run(context.Background()))
				assert.WithinDuration(t,
					startTime.Add(data.Duration),
					startTime.Add(e.Executor.GetTime()),
					100*TickRate,
				)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("collects samples", func(t *testing.T) ***REMOVED***
		testMetric := stats.New("test_metric", stats.Trend)

		signalChan := make(chan interface***REMOVED******REMOVED***)
		var e *Engine
		e, err, _ := newTestEngine(LF(func(ctx context.Context) (samples []stats.Sample, err error) ***REMOVED***
			samples = append(samples, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 1***REMOVED***)
			close(signalChan)
			<-ctx.Done()
			samples = append(samples, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 2***REMOVED***)
			return samples, err
		***REMOVED***), lib.Options***REMOVED***
			VUs:        null.IntFrom(1),
			VUsMax:     null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		c := &dummy.Collector***REMOVED******REMOVED***
		e.Collector = c

		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		go func() ***REMOVED*** errC <- e.Run(ctx) ***REMOVED***()
		<-signalChan
		cancel()
		assert.NoError(t, <-errC)

		found := 0
		for _, s := range c.Samples ***REMOVED***
			if s.Metric != testMetric ***REMOVED***
				continue
			***REMOVED***
			found++
			assert.Equal(t, 1.0, s.Value, "wrong value")
		***REMOVED***
		assert.Equal(t, 1, found, "wrong number of samples")
	***REMOVED***)
***REMOVED***

func TestEngineAtTime(t *testing.T) ***REMOVED***
	e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	assert.NoError(t, e.Run(ctx))
***REMOVED***

func TestEngineCollector(t *testing.T) ***REMOVED***
	testMetric := stats.New("test_metric", stats.Trend)

	e, err, _ := newTestEngine(LF(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
		return []stats.Sample***REMOVED******REMOVED***Metric: testMetric***REMOVED******REMOVED***, nil
	***REMOVED***), lib.Options***REMOVED***VUs: null.IntFrom(1), VUsMax: null.IntFrom(1), Iterations: null.IntFrom(1)***REMOVED***)
	assert.NoError(t, err)

	c := &dummy.Collector***REMOVED******REMOVED***
	e.Collector = c

	assert.NoError(t, e.Run(context.Background()))

	cSamples := []stats.Sample***REMOVED******REMOVED***
	for _, sample := range c.Samples ***REMOVED***
		if sample.Metric == testMetric ***REMOVED***
			cSamples = append(cSamples, sample)
		***REMOVED***
	***REMOVED***
	metric := e.Metrics["test_metric"]
	if assert.NotNil(t, metric) ***REMOVED***
		sink := metric.Sink.(*stats.TrendSink)
		if assert.NotNil(t, sink) ***REMOVED***
			numCollectorSamples := len(cSamples)
			numEngineSamples := len(sink.Values)
			assert.Equal(t, numEngineSamples, numCollectorSamples)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEngine_processSamples(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	t.Run("metric", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
	***REMOVED***)
	t.Run("submetric", func(t *testing.T) ***REMOVED***
		ths, err := stats.NewThresholds([]string***REMOVED***`1+1==2`***REMOVED***)
		assert.NoError(t, err)

		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric***REMOVED***a:1***REMOVED***": ths,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)

		sms := e.submetrics["my_metric"]
		assert.Len(t, sms, 1)
		assert.Equal(t, "my_metric***REMOVED***a:1***REMOVED***", sms[0].Name)
		assert.EqualValues(t, map[string]string***REMOVED***"a": "1"***REMOVED***, sms[0].Tags)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric***REMOVED***a:1***REMOVED***"].Sink)
	***REMOVED***)
***REMOVED***

func TestEngine_processThresholds(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	testdata := map[string]struct ***REMOVED***
		pass bool
		ths  map[string][]string
	***REMOVED******REMOVED***
		"passing": ***REMOVED***true, map[string][]string***REMOVED***"my_metric": ***REMOVED***"1+1==2"***REMOVED******REMOVED******REMOVED***,
		"failing": ***REMOVED***false, map[string][]string***REMOVED***"my_metric": ***REMOVED***"1+1==3"***REMOVED******REMOVED******REMOVED***,

		"submetric,match,passing":   ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:1***REMOVED***": ***REMOVED***"1+1==2"***REMOVED******REMOVED******REMOVED***,
		"submetric,match,failing":   ***REMOVED***false, map[string][]string***REMOVED***"my_metric***REMOVED***a:1***REMOVED***": ***REMOVED***"1+1==3"***REMOVED******REMOVED******REMOVED***,
		"submetric,nomatch,passing": ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:2***REMOVED***": ***REMOVED***"1+1==2"***REMOVED******REMOVED******REMOVED***,
		"submetric,nomatch,failing": ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:2***REMOVED***": ***REMOVED***"1+1==3"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			thresholds := make(map[string]stats.Thresholds, len(data.ths))
			for m, srcs := range data.ths ***REMOVED***
				ths, err := stats.NewThresholds(srcs)
				assert.NoError(t, err)
				thresholds[m] = ths
			***REMOVED***

			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
			assert.NoError(t, err)

			e.processSamples(
				stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
			)
			e.processThresholds()

			assert.Equal(t, data.pass, !e.IsTainted())
		***REMOVED***)
	***REMOVED***
***REMOVED***
