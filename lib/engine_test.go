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
	"fmt"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"runtime"
	"testing"
	"time"
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
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		startTime := time.Now()
		duration := 100 * time.Millisecond
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, _ := context.WithTimeout(context.Background(), duration)
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
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		d := 50 * time.Millisecond
		e.Stages = []Stage***REMOVED******REMOVED***Duration: d***REMOVED******REMOVED***
		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background()))
		assert.WithinDuration(t, startTime.Add(d), startTime.Add(e.AtTime()), 100*TickRate)
	***REMOVED***)
	t.Run("exits with AbortOnTaint", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***AbortOnTaint: null.BoolFrom(true)***REMOVED***)
		assert.NoError(t, err)

		ch := make(chan error)
		go func() ***REMOVED*** ch <- e.Run(context.Background()) ***REMOVED***()
		time.Sleep(1 * time.Millisecond)
		assert.True(t, e.IsRunning())

		e.lock.Lock()
		e.numTaints++
		e.lock.Unlock()

		assert.EqualError(t, <-ch, "test is tainted")
		assert.False(t, e.IsRunning())
	***REMOVED***)
	t.Run("collects samples", func(t *testing.T) ***REMOVED***
		testMetric := stats.New("test_metric", stats.Trend)

		errors := map[string]error***REMOVED***
			"nil":   nil,
			"error": errors.New("error"),
			"taint": ErrVUWantsTaint,
		***REMOVED***
		for name, reterr := range errors ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err, _ := newTestEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return []stats.Sample***REMOVED******REMOVED***Metric: testMetric, Value: 1.0***REMOVED******REMOVED***, reterr
				***REMOVED***), Options***REMOVED***VUsMax: null.IntFrom(1), VUs: null.IntFrom(1)***REMOVED***)
				assert.NoError(t, err)

				ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
				assert.NoError(t, e.Run(ctx))

				e.lock.Lock()
				if !assert.True(t, e.numIterations > 0, "no iterations performed") ***REMOVED***
					e.lock.Unlock()
					return
				***REMOVED***
				sink := e.Metrics[testMetric].(*stats.TrendSink)
				assert.True(t, len(sink.Values) > int(float64(e.numIterations)*0.99), "more than 1%% of iterations missed")
				e.lock.Unlock()
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
				e, err, _ := newTestEngine(nil, Options***REMOVED***Duration: null.StringFrom(d.String())***REMOVED***)
				assert.NoError(t, err)

				assert.Len(t, e.Stages, 1)
				assert.Equal(t, Stage***REMOVED***Duration: d***REMOVED***, e.Stages[0])
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		// The lines get way too damn long if I have to write time.Second everywhere
		sec := time.Second

		testdata := map[string]struct ***REMOVED***
			Duration time.Duration
			Stages   []Stage
		***REMOVED******REMOVED***
			"nil":        ***REMOVED***0, nil***REMOVED***,
			"empty":      ***REMOVED***0, []Stage***REMOVED******REMOVED******REMOVED***,
			"1,infinite": ***REMOVED***0, []Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			"2,infinite": ***REMOVED***0, []Stage***REMOVED******REMOVED***Duration: 10 * sec***REMOVED***, ***REMOVED******REMOVED******REMOVED******REMOVED***,
			"1,finite":   ***REMOVED***10 * sec, []Stage***REMOVED******REMOVED***Duration: 10 * sec***REMOVED******REMOVED******REMOVED***,
			"2,finite":   ***REMOVED***15 * sec, []Stage***REMOVED******REMOVED***Duration: 10 * sec***REMOVED***, ***REMOVED***Duration: 5 * sec***REMOVED******REMOVED******REMOVED***,
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

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
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
		iterationSampleA1 := e.numIterations
		atTimeSampleA1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleA2 := e.numIterations
		atTimeSampleA2 := e.AtTime()
		assert.True(t, iterationSampleA2 > iterationSampleA1, "iteration counter did not increase")
		assert.True(t, atTimeSampleA2 > atTimeSampleA1, "timer did not increase")

		// ...stop increasing when you pause... (sleep to ensure outstanding VUs finish)
		e.SetPaused(true)
		assert.True(t, e.IsPaused(), "engine did not pause")
		time.Sleep(1 * time.Millisecond)
		iterationSampleB1 := e.numIterations
		atTimeSampleB1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleB2 := e.numIterations
		atTimeSampleB2 := e.AtTime()
		assert.Equal(t, iterationSampleB1, iterationSampleB2, "iteration counter changed while paused")
		assert.Equal(t, atTimeSampleB1, atTimeSampleB2, "timer changed while paused")

		// ...and resume when you unpause.
		e.SetPaused(false)
		assert.False(t, e.IsPaused(), "engine did not unpause")
		iterationSampleC1 := e.numIterations
		atTimeSampleC1 := e.AtTime()
		time.Sleep(100 * time.Millisecond)
		iterationSampleC2 := e.numIterations
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
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Equal(t, int64(0), e.GetVUs())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(15)***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUs(10))
		assert.Equal(t, int64(10), e.GetVUs())
		assertActiveVUs(t, e, 10, 5)

		t.Run("negative", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(-1), "vus can't be negative")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
		***REMOVED***)

		t.Run("too high", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(20), "more vus than allocated requested")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
		***REMOVED***)

		t.Run("lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(5))
			assert.Equal(t, int64(5), e.GetVUs())
			assertActiveVUs(t, e, 5, 10)
		***REMOVED***)

		t.Run("higher", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(15))
			assert.Equal(t, int64(15), e.GetVUs())
			assertActiveVUs(t, e, 15, 0)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngineIsTainted(t *testing.T) ***REMOVED***
	testdata := []struct ***REMOVED***
		I      int64
		T      int64
		Expect bool
	***REMOVED******REMOVED***
		***REMOVED***1, 0, false***REMOVED***,
		***REMOVED***1, 1, true***REMOVED***,
	***REMOVED***

	for _, data := range testdata ***REMOVED***
		t.Run(fmt.Sprintf("i=%d,t=%d", data.I, data.T), func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)

			e.numIterations = data.I
			e.numTaints = data.T
			assert.Equal(t, data.Expect, e.IsTainted())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestEngine_runVUOnceKeepsCounters(t *testing.T) ***REMOVED***
	e, err, hook := newTestEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), e.numIterations)
	assert.Equal(t, int64(0), e.numTaints)

	t.Run("success", func(t *testing.T) ***REMOVED***
		hook.Reset()
		e.numIterations = 0
		e.numTaints = 0
		e.runVUOnce(context.Background(), &vuEntry***REMOVED***
			VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
				return nil, nil
			***REMOVED***),
		***REMOVED***)
		assert.Equal(t, int64(1), e.numIterations)
		assert.Equal(t, int64(0), e.numTaints)
		assert.False(t, e.IsTainted(), "test is tainted")
	***REMOVED***)
	t.Run("error", func(t *testing.T) ***REMOVED***
		hook.Reset()
		e.numIterations = 0
		e.numTaints = 0
		e.runVUOnce(context.Background(), &vuEntry***REMOVED***
			VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
				return nil, errors.New("this is an error")
			***REMOVED***),
		***REMOVED***)
		assert.Equal(t, int64(1), e.numIterations)
		assert.Equal(t, int64(1), e.numTaints)
		assert.True(t, e.IsTainted(), "test is not tainted")
		assert.Equal(t, "this is an error", hook.LastEntry().Data["error"].(error).Error())

		t.Run("string", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numTaints = 0
			e.runVUOnce(context.Background(), &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, testErrorWithString("this is an error")
				***REMOVED***),
			***REMOVED***)
			assert.Equal(t, int64(1), e.numIterations)
			assert.Equal(t, int64(1), e.numTaints)
			assert.True(t, e.IsTainted(), "test is not tainted")

			entry := hook.LastEntry()
			assert.Equal(t, "this is an error", entry.Message)
			assert.Empty(t, entry.Data)
		***REMOVED***)
	***REMOVED***)
	t.Run("taint", func(t *testing.T) ***REMOVED***
		hook.Reset()
		e.numIterations = 0
		e.numTaints = 0
		e.runVUOnce(context.Background(), &vuEntry***REMOVED***
			VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
				return nil, ErrVUWantsTaint
			***REMOVED***),
		***REMOVED***)
		assert.Equal(t, int64(1), e.numIterations)
		assert.Equal(t, int64(1), e.numTaints)
		assert.True(t, e.IsTainted(), "test is not tainted")

		assert.Nil(t, hook.LastEntry())
	***REMOVED***)
	t.Run("cancelled", func(t *testing.T) ***REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		t.Run("success", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numTaints = 0
			e.runVUOnce(ctx, &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, nil
				***REMOVED***),
			***REMOVED***)
			assert.Equal(t, int64(0), e.numIterations)
			assert.Equal(t, int64(0), e.numTaints)
			assert.False(t, e.IsTainted(), "test is tainted")
		***REMOVED***)
		t.Run("error", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numTaints = 0
			e.runVUOnce(ctx, &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, errors.New("this is an error")
				***REMOVED***),
			***REMOVED***)
			assert.Equal(t, int64(0), e.numIterations)
			assert.Equal(t, int64(0), e.numTaints)
			assert.False(t, e.IsTainted(), "test is tainted")
			assert.Nil(t, hook.LastEntry())

			t.Run("string", func(t *testing.T) ***REMOVED***
				hook.Reset()
				e.numIterations = 0
				e.numTaints = 0
				e.runVUOnce(ctx, &vuEntry***REMOVED***
					VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
						return nil, testErrorWithString("this is an error")
					***REMOVED***),
				***REMOVED***)
				assert.Equal(t, int64(0), e.numIterations)
				assert.Equal(t, int64(0), e.numTaints)
				assert.False(t, e.IsTainted(), "test is tainted")

				assert.Nil(t, hook.LastEntry())
			***REMOVED***)
		***REMOVED***)
		t.Run("taint", func(t *testing.T) ***REMOVED***
			hook.Reset()
			e.numIterations = 0
			e.numTaints = 0
			e.runVUOnce(ctx, &vuEntry***REMOVED***
				VU: RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return nil, ErrVUWantsTaint
				***REMOVED***),
			***REMOVED***)
			assert.Equal(t, int64(0), e.numIterations)
			assert.Equal(t, int64(0), e.numTaints)
			assert.False(t, e.IsTainted(), "test is tainted")

			assert.Nil(t, hook.LastEntry())
		***REMOVED***)
	***REMOVED***)
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

	// Allow 10% of samples to get lost; NOT OPTIMAL, but I can't figure out why they get lost.
	numSamples := len(e.Metrics[testMetric].(*stats.TrendSink).Values)
	assert.True(t, numSamples > 0, "no samples")
	assert.True(t, numSamples > len(c.Samples)-(len(c.Samples)/10), "more than 10%% of samples omitted")
***REMOVED***
