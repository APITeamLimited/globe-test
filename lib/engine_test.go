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

package lib

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
	"github.com/pkg/errors"
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
	e.Logger = logger
	return hook
***REMOVED***

// Wrapper around newEngine that applies a null logger.
func newTestEngine(r Runner, opts Options) (*Engine, error, *logtest.Hook) ***REMOVED***
	e, err := NewEngine(r, opts)
	if err != nil ***REMOVED***
		return e, err, nil
	***REMOVED***
	hook := applyNullLogger(e)
	return e, nil, hook
***REMOVED***

// Helper for asserting the number of active/dead VUs.
func assertActiveVUs(t *testing.T, e *Engine, active, dead int) ***REMOVED***
	e.lock.Lock()
	defer e.lock.Unlock()

	var numActive, numDead int
	var lastWasDead bool
	for _, vu := range e.vuEntries ***REMOVED***
		if vu.Cancel != nil ***REMOVED***
			numActive++
			assert.False(t, lastWasDead, "living vu in dead zone")
		***REMOVED*** else ***REMOVED***
			numDead++
			lastWasDead = true
		***REMOVED***
	***REMOVED***
	assert.Equal(t, active, numActive, "wrong number of active vus")
	assert.Equal(t, dead, numDead, "wrong number of dead vus")
***REMOVED***

func TestNewEngine(t *testing.T) ***REMOVED***
	_, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***
			Duration: NullDurationFrom(10 * time.Second),
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], Stage***REMOVED***Duration: NullDurationFrom(10 * time.Second)***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***
			Stages: []Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], Stage***REMOVED***Duration: NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Stages/Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***
			Duration: NullDurationFrom(60 * time.Second),
			Stages: []Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Stages, 1) ***REMOVED***
			assert.Equal(t, e.Stages[0], Stage***REMOVED***Duration: NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(1), e.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(10), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			assert.NoError(t, err)
			assert.True(t, e.IsPaused())
		***REMOVED***)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric": ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		assert.Contains(t, e.thresholds, "my_metric")

		t.Run("submetrics", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
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
		startTime := time.Now()
		duration := 100 * time.Millisecond
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("terminates subctx", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		subctx := e.subctx
		select ***REMOVED***
		case <-subctx.Done():
			assert.Fail(t, "context is already terminated")
		default:
		***REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.NoError(t, e.Run(ctx))

		assert.NotEqual(t, subctx, e.subctx, "subcontext not changed")
		select ***REMOVED***
		case <-subctx.Done():
		default:
			assert.Fail(t, "context was not terminated")
		***REMOVED***
	***REMOVED***)
	t.Run("exits with stages", func(t *testing.T) ***REMOVED***
		testdata := map[string]struct ***REMOVED***
			Duration time.Duration
			Stages   []Stage
		***REMOVED******REMOVED***
			"none": ***REMOVED******REMOVED***,
			"one": ***REMOVED***
				1 * time.Second,
				[]Stage***REMOVED******REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED******REMOVED***,
			***REMOVED***,
			"two": ***REMOVED***
				2 * time.Second,
				[]Stage***REMOVED***
					***REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED***,
					***REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			"two/targeted": ***REMOVED***
				2 * time.Second,
				[]Stage***REMOVED***
					***REMOVED***Duration: NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
					***REMOVED***Duration: NullDurationFrom(1 * time.Second), Target: null.IntFrom(10)***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(10)***REMOVED***)
				assert.NoError(t, err)

				e.Stages = data.Stages
				startTime := time.Now()
				assert.NoError(t, e.Run(context.Background()))
				assert.WithinDuration(t,
					startTime.Add(data.Duration),
					startTime.Add(e.AtTime()),
					100*TickRate,
				)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("collects samples", func(t *testing.T) ***REMOVED***
		testMetric := stats.New("test_metric", stats.Trend)

		errors := map[string]error***REMOVED***
			"nil":   nil,
			"error": errors.New("error"),
		***REMOVED***
		for name, reterr := range errors ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return []stats.Sample***REMOVED******REMOVED***Metric: testMetric, Value: 1.0***REMOVED******REMOVED***, reterr
				***REMOVED***), Options***REMOVED***VUsMax: null.IntFrom(1), VUs: null.IntFrom(1)***REMOVED***)
				assert.NoError(t, err)

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				assert.NoError(t, e.Run(ctx))
				cancel()

				e.lock.Lock()
				defer e.lock.Unlock()

				if !assert.True(t, e.numIterations > 0, "no iterations performed") ***REMOVED***
					return
				***REMOVED***
				sink := e.Metrics["test_metric"].Sink.(*stats.TrendSink)
				assert.True(t, len(sink.Values) > int(float64(e.numIterations)*0.99), "more than 1%% of iterations missed")
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestEngineIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	ch := make(chan error)
	go func() ***REMOVED*** ch <- e.Run(ctx) ***REMOVED***()
	runtime.Gosched()
	time.Sleep(1 * time.Millisecond)
	assert.True(t, e.IsRunning())

	cancel()
	runtime.Gosched()
	time.Sleep(1 * time.Millisecond)
	assert.False(t, e.IsRunning())

	assert.NoError(t, <-ch)
***REMOVED***

func TestEngineTotalTime(t *testing.T) ***REMOVED***
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		for _, d := range []time.Duration***REMOVED***0, 1 * time.Second, 10 * time.Second***REMOVED*** ***REMOVED***
			t.Run(d.String(), func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(nil, Options***REMOVED***Duration: NullDurationFrom(d)***REMOVED***)
				assert.NoError(t, err)

				assert.Len(t, e.Stages, 1)
				assert.Equal(t, Stage***REMOVED***Duration: NullDurationFrom(d)***REMOVED***, e.Stages[0])
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		testdata := map[string]struct ***REMOVED***
			Duration time.Duration
			Stages   []Stage
		***REMOVED******REMOVED***
			"nil":        ***REMOVED***0, nil***REMOVED***,
			"empty":      ***REMOVED***0, []Stage***REMOVED******REMOVED******REMOVED***,
			"1,infinite": ***REMOVED***0, []Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			"2,infinite": ***REMOVED***0, []Stage***REMOVED******REMOVED***Duration: NullDurationFrom(10 * time.Second)***REMOVED***, ***REMOVED******REMOVED******REMOVED******REMOVED***,
			"1,finite":   ***REMOVED***10 * time.Second, []Stage***REMOVED******REMOVED***Duration: NullDurationFrom(10 * time.Second)***REMOVED******REMOVED******REMOVED***,
			"2,finite": ***REMOVED***15 * time.Second, []Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(10 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(nil, Options***REMOVED***Stages: data.Stages***REMOVED***)
				assert.NoError(t, err)
				assert.Equal(t, data.Duration, e.TotalTime())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestEngineAtTime(t *testing.T) ***REMOVED***
	e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	assert.NoError(t, e.Run(ctx))
***REMOVED***

func TestEngineSetPaused(t *testing.T) ***REMOVED***
	t.Run("offline", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.False(t, e.IsPaused())

		e.SetPaused(true)
		assert.True(t, e.IsPaused())

		e.SetPaused(false)
		assert.False(t, e.IsPaused())
	***REMOVED***)

	t.Run("running", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			return nil, nil
		***REMOVED***), Options***REMOVED***VUsMax: null.IntFrom(1), VUs: null.IntFrom(1)***REMOVED***)
		assert.NoError(t, err)
		assert.False(t, e.IsPaused())

		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan error)
		go func() ***REMOVED*** ch <- e.Run(ctx) ***REMOVED***()
		time.Sleep(1 * time.Millisecond)
		assert.True(t, e.IsRunning())

		// The iteration counter and time should increase over time when not paused...
		iterationSampleA1 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleA1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleA2 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleA2 := e.AtTime()
		assert.True(t, iterationSampleA2 > iterationSampleA1, "iteration counter did not increase")
		assert.True(t, atTimeSampleA2 > atTimeSampleA1, "timer did not increase")

		// ...stop increasing when you pause... (sleep to ensure outstanding VUs finish)
		e.SetPaused(true)
		assert.True(t, e.IsPaused(), "engine did not pause")
		time.Sleep(1 * time.Millisecond)
		iterationSampleB1 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleB1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleB2 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleB2 := e.AtTime()
		assert.Equal(t, iterationSampleB1, iterationSampleB2, "iteration counter changed while paused")
		assert.Equal(t, atTimeSampleB1, atTimeSampleB2, "timer changed while paused")

		// ...and resume when you unpause.
		e.SetPaused(false)
		assert.False(t, e.IsPaused(), "engine did not unpause")
		iterationSampleC1 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleC1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleC2 := atomic.LoadInt64(&e.numIterations)
		atTimeSampleC2 := e.AtTime()
		assert.True(t, iterationSampleC2 > iterationSampleC1, "iteration counter did not increase after unpause")
		assert.True(t, atTimeSampleC2 > atTimeSampleC1, "timer did not increase after unpause")

		cancel()
		assert.NoError(t, <-ch)
	***REMOVED***)

	t.Run("exit", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			return nil, nil
		***REMOVED***), Options***REMOVED***VUsMax: null.IntFrom(1), VUs: null.IntFrom(1)***REMOVED***)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan error)
		go func() ***REMOVED*** ch <- e.Run(ctx) ***REMOVED***()
		time.Sleep(1 * time.Millisecond)
		assert.True(t, e.IsRunning())

		e.SetPaused(true)
		assert.True(t, e.IsPaused())
		cancel()
		time.Sleep(1 * time.Millisecond)
		assert.False(t, e.IsPaused())
		assert.False(t, e.IsRunning())

		assert.NoError(t, <-ch)
	***REMOVED***)
***REMOVED***

func TestEngineSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUsMax(10))
		assert.Equal(t, int64(10), e.GetVUsMax())
		assert.Len(t, e.vuEntries, 10)
		for _, vu := range e.vuEntries ***REMOVED***
			assert.Nil(t, vu.Cancel)
		***REMOVED***

		t.Run("higher", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(15))
			assert.Equal(t, int64(15), e.GetVUsMax())
			assert.Len(t, e.vuEntries, 15)
			for _, vu := range e.vuEntries ***REMOVED***
				assert.Nil(t, vu.Cancel)
			***REMOVED***
		***REMOVED***)

		t.Run("lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(5))
			assert.Equal(t, int64(5), e.GetVUsMax())
			assert.Len(t, e.vuEntries, 5)
			for _, vu := range e.vuEntries ***REMOVED***
				assert.Nil(t, vu.Cancel)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("set negative", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(-1), "vus-max can't be negative")
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set too low", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***
			VUsMax: null.IntFrom(10),
			VUs:    null.IntFrom(10),
		***REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(5), "can't reduce vus-max below vus")
		assert.Len(t, e.vuEntries, 10)
	***REMOVED***)
***REMOVED***

func TestEngineSetVUs(t *testing.T) ***REMOVED***
	assertVUIDSequence := func(t *testing.T, e *Engine, ids []int64) ***REMOVED***
		actualIDs := make([]int64, len(ids))
		for i := range ids ***REMOVED***
			actualIDs[i] = e.vuEntries[i].VU.(*RunnerFuncVU).ID
		***REMOVED***
		assert.Equal(t, ids, actualIDs)
	***REMOVED***

	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Equal(t, int64(0), e.GetVUs())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(RunnerFunc(nil), Options***REMOVED***VUsMax: null.IntFrom(15)***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUs(10))
		assert.Equal(t, int64(10), e.GetVUs())
		assertActiveVUs(t, e, 10, 5)
		assertVUIDSequence(t, e, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10***REMOVED***)

		t.Run("negative", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(-1), "vus can't be negative")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
			assertVUIDSequence(t, e, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10***REMOVED***)
		***REMOVED***)

		t.Run("too high", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(20), "more vus than allocated requested")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
			assertVUIDSequence(t, e, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10***REMOVED***)
		***REMOVED***)

		t.Run("lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(5))
			assert.Equal(t, int64(5), e.GetVUs())
			assertActiveVUs(t, e, 5, 10)
			assertVUIDSequence(t, e, []int64***REMOVED***1, 2, 3, 4, 5***REMOVED***)
		***REMOVED***)

		t.Run("higher", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(15))
			assert.Equal(t, int64(15), e.GetVUs())
			assertActiveVUs(t, e, 15, 0)
			assertVUIDSequence(t, e, []int64***REMOVED***1, 2, 3, 4, 5, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngine_runVUOnceKeepsCounters(t *testing.T) ***REMOVED***
	e, err, hook := newTestEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), e.numIterations)
	assert.Equal(t, int64(0), e.numErrors)

	t.Run("success", func(t *testing.T) ***REMOVED***
		hook.Reset()
		e.numIterations = 0
		e.numErrors = 0
		e.runVUOnce(context.Background(), &vuEntry***REMOVED***
			VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
				return nil, nil
			***REMOVED***).VU(),
		***REMOVED***)
		assert.Equal(t, int64(1), e.numIterations)
		assert.Equal(t, int64(0), e.numErrors)
		assert.False(t, e.IsTainted(), "test is tainted")
	***REMOVED***)
	t.Run("error", func(t *testing.T) ***REMOVED***
		hook.Reset()
		e.numIterations = 0
		e.numErrors = 0
		e.runVUOnce(context.Background(), &vuEntry***REMOVED***
			VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
				return nil, errors.New("this is an error")
			***REMOVED***).VU(),
		***REMOVED***)
		assert.Equal(t, int64(1), e.numIterations)
		assert.Equal(t, int64(1), e.numErrors)
		assert.Equal(t, "this is an error", hook.LastEntry().Data["error"].(error).Error())

		t.Run("string", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numErrors = 0
			e.runVUOnce(context.Background(), &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, testErrorWithString("this is an error")
				***REMOVED***).VU(),
			***REMOVED***)
			assert.Equal(t, int64(1), e.numIterations)
			assert.Equal(t, int64(1), e.numErrors)

			entry := hook.LastEntry()
			assert.Equal(t, "this is an error", entry.Message)
			assert.Empty(t, entry.Data)
		***REMOVED***)
	***REMOVED***)
	t.Run("cancelled", func(t *testing.T) ***REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		t.Run("success", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numErrors = 0
			e.runVUOnce(ctx, &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, nil
				***REMOVED***).VU(),
			***REMOVED***)
			assert.Equal(t, int64(0), e.numIterations)
			assert.Equal(t, int64(0), e.numErrors)
		***REMOVED***)
		t.Run("error", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numErrors = 0
			e.runVUOnce(ctx, &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, errors.New("this is an error")
				***REMOVED***).VU(),
			***REMOVED***)
			assert.Equal(t, int64(0), e.numIterations)
			assert.Equal(t, int64(0), e.numErrors)
			assert.Nil(t, hook.LastEntry())

			t.Run("string", func(t *testing.T) ***REMOVED***
				hook.Reset()
				e.numIterations = 0
				e.numErrors = 0
				e.runVUOnce(ctx, &vuEntry***REMOVED***
					VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
						return nil, testErrorWithString("this is an error")
					***REMOVED***).VU(),
				***REMOVED***)
				assert.Equal(t, int64(0), e.numIterations)
				assert.Equal(t, int64(0), e.numErrors)
				assert.Nil(t, hook.LastEntry())
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngine_processStages(t *testing.T) ***REMOVED***
	type checkpoint struct ***REMOVED***
		D    time.Duration
		Cont bool
		VUs  int64
	***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Stages      []Stage
		Checkpoints []checkpoint
	***REMOVED******REMOVED***
		"none": ***REMOVED***
			[]Stage***REMOVED******REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, false, 0***REMOVED***,
				***REMOVED***10 * time.Second, false, 0***REMOVED***,
				***REMOVED***24 * time.Hour, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(10 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***10 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one/targeted": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(10 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 30***REMOVED***,
				***REMOVED***1 * time.Second, true, 40***REMOVED***,
				***REMOVED***1 * time.Second, true, 50***REMOVED***,
				***REMOVED***1 * time.Second, true, 60***REMOVED***,
				***REMOVED***1 * time.Second, true, 70***REMOVED***,
				***REMOVED***1 * time.Second, true, 80***REMOVED***,
				***REMOVED***1 * time.Second, true, 90***REMOVED***,
				***REMOVED***1 * time.Second, true, 100***REMOVED***,
				***REMOVED***1 * time.Second, false, 100***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***10 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/targeted": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 40***REMOVED***,
				***REMOVED***1 * time.Second, true, 60***REMOVED***,
				***REMOVED***1 * time.Second, true, 80***REMOVED***,
				***REMOVED***1 * time.Second, true, 100***REMOVED***,
				***REMOVED***1 * time.Second, true, 80***REMOVED***,
				***REMOVED***1 * time.Second, true, 60***REMOVED***,
				***REMOVED***1 * time.Second, true, 40***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***15 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three/targeted": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(50)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 30***REMOVED***,
				***REMOVED***1 * time.Second, true, 40***REMOVED***,
				***REMOVED***1 * time.Second, true, 50***REMOVED***,
				***REMOVED***1 * time.Second, true, 60***REMOVED***,
				***REMOVED***1 * time.Second, true, 70***REMOVED***,
				***REMOVED***1 * time.Second, true, 80***REMOVED***,
				***REMOVED***1 * time.Second, true, 90***REMOVED***,
				***REMOVED***1 * time.Second, true, 100***REMOVED***,
				***REMOVED***1 * time.Second, true, 80***REMOVED***,
				***REMOVED***1 * time.Second, true, 60***REMOVED***,
				***REMOVED***1 * time.Second, true, 40***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"mix": ***REMOVED***
			[]Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,

				***REMOVED***1 * time.Second, true, 4***REMOVED***,
				***REMOVED***1 * time.Second, true, 8***REMOVED***,
				***REMOVED***1 * time.Second, true, 12***REMOVED***,
				***REMOVED***1 * time.Second, true, 16***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,

				***REMOVED***1 * time.Second, true, 18***REMOVED***,
				***REMOVED***1 * time.Second, true, 16***REMOVED***,
				***REMOVED***1 * time.Second, true, 14***REMOVED***,
				***REMOVED***1 * time.Second, true, 12***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,

				***REMOVED***1 * time.Second, true, 10***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,

				***REMOVED***1 * time.Second, true, 12***REMOVED***,
				***REMOVED***1 * time.Second, true, 14***REMOVED***,
				***REMOVED***1 * time.Second, true, 16***REMOVED***,
				***REMOVED***1 * time.Second, true, 18***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,

				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,

				***REMOVED***1 * time.Second, true, 18***REMOVED***,
				***REMOVED***1 * time.Second, true, 16***REMOVED***,
				***REMOVED***1 * time.Second, true, 14***REMOVED***,
				***REMOVED***1 * time.Second, true, 12***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED***
				VUs:    null.IntFrom(0),
				VUsMax: null.IntFrom(100),
			***REMOVED***)
			assert.NoError(t, err)

			e.Stages = data.Stages
			for _, ckp := range data.Checkpoints ***REMOVED***
				t.Run((e.AtTime() + ckp.D).String(), func(t *testing.T) ***REMOVED***
					cont, err := e.processStages(ckp.D)
					assert.NoError(t, err)
					if ckp.Cont ***REMOVED***
						assert.True(t, cont, "test stopped")
					***REMOVED*** else ***REMOVED***
						assert.False(t, cont, "test not stopped")
					***REMOVED***
					assert.Equal(t, ckp.VUs, e.GetVUs())
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestEngineCollector(t *testing.T) ***REMOVED***
	testMetric := stats.New("test_metric", stats.Trend)
	c := &dummy.Collector***REMOVED******REMOVED***

	e, err, _ := newTestEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
		return []stats.Sample***REMOVED******REMOVED***Metric: testMetric***REMOVED******REMOVED***, nil
	***REMOVED***), Options***REMOVED***VUs: null.IntFrom(1), VUsMax: null.IntFrom(1)***REMOVED***)
	assert.NoError(t, err)
	e.Collector = c

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error)
	go func() ***REMOVED*** ch <- e.Run(ctx) ***REMOVED***()

	time.Sleep(100 * time.Millisecond)
	assert.True(t, e.IsRunning(), "engine not running")
	assert.True(t, c.IsRunning(), "collector not running")

	cancel()
	assert.NoError(t, <-ch)

	assert.False(t, e.IsRunning(), "engine still running")
	assert.False(t, c.IsRunning(), "collector still running")

	cSamples := []stats.Sample***REMOVED******REMOVED***
	for _, sample := range c.Samples ***REMOVED***
		if sample.Metric == testMetric ***REMOVED***
			cSamples = append(cSamples, sample)
		***REMOVED***
	***REMOVED***
	numCollectorSamples := len(cSamples)
	numEngineSamples := len(e.Metrics["test_metric"].Sink.(*stats.TrendSink).Values)
	assert.Equal(t, numEngineSamples, numCollectorSamples)
***REMOVED***

func TestEngine_processSamples(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	t.Run("metric", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
	***REMOVED***)
	t.Run("submetric", func(t *testing.T) ***REMOVED***
		ths, err := stats.NewThresholds([]string***REMOVED***`1+1==2`***REMOVED***)
		assert.NoError(t, err)

		e, err, _ := newTestEngine(nil, Options***REMOVED***
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

			e, err, _ := newTestEngine(nil, Options***REMOVED***Thresholds: thresholds***REMOVED***)
			assert.NoError(t, err)

			e.processSamples(
				stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
			)
			e.processThresholds()

			assert.Equal(t, data.pass, !e.IsTainted())
		***REMOVED***)
	***REMOVED***
***REMOVED***
