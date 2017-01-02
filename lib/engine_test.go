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

// Helper for asserting the number of active/dead VUs.
func assertActiveVUs(t *testing.T, e *Engine, active, dead int) ***REMOVED***
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
	_, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(1), e.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
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
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
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
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, _ := context.WithTimeout(context.Background(), duration)
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("terminates subctx", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
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
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		d := 50 * time.Millisecond
		e.Stages = []Stage***REMOVED***Stage***REMOVED***Duration: d***REMOVED******REMOVED***
		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background()))
		assert.WithinDuration(t, startTime.Add(d), startTime.Add(e.AtTime()), 2*TickRate)
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
				e, err := NewEngine(RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
					return []stats.Sample***REMOVED***stats.Sample***REMOVED***Metric: testMetric, Value: 1.0***REMOVED******REMOVED***, reterr
				***REMOVED***), Options***REMOVED***VUsMax: null.IntFrom(1), VUs: null.IntFrom(1)***REMOVED***)
				assert.NoError(t, err)

				ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
				assert.NoError(t, e.Run(ctx))
				if !assert.True(t, e.numIterations > 0, "no iterations performed") ***REMOVED***
					return
				***REMOVED***

				sink := e.Metrics[testMetric].(*stats.TrendSink)
				assert.Equal(t, int(e.numIterations), len(sink.Values))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestEngineIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	go func() ***REMOVED*** assert.NoError(t, e.Run(ctx)) ***REMOVED***()
	runtime.Gosched()
	time.Sleep(1 * time.Millisecond)
	assert.True(t, e.IsRunning())

	cancel()
	runtime.Gosched()
	time.Sleep(1 * time.Millisecond)
	assert.False(t, e.IsRunning())
***REMOVED***

func TestEngineTotalTime(t *testing.T) ***REMOVED***
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		for _, d := range []time.Duration***REMOVED***0, 1 * time.Second, 10 * time.Second***REMOVED*** ***REMOVED***
			t.Run(d.String(), func(t *testing.T) ***REMOVED***
				e, err := NewEngine(nil, Options***REMOVED***Duration: null.StringFrom(d.String())***REMOVED***)
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
			"1,infinite": ***REMOVED***0, []Stage***REMOVED***Stage***REMOVED******REMOVED******REMOVED******REMOVED***,
			"2,infinite": ***REMOVED***0, []Stage***REMOVED***Stage***REMOVED***Duration: 10 * sec***REMOVED***, Stage***REMOVED******REMOVED******REMOVED******REMOVED***,
			"1,finite":   ***REMOVED***10 * sec, []Stage***REMOVED***Stage***REMOVED***Duration: 10 * sec***REMOVED******REMOVED******REMOVED***,
			"2,finite":   ***REMOVED***15 * sec, []Stage***REMOVED***Stage***REMOVED***Duration: 10 * sec***REMOVED***, Stage***REMOVED***Duration: 5 * sec***REMOVED******REMOVED******REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				e, err := NewEngine(nil, Options***REMOVED***Stages: data.Stages***REMOVED***)
				assert.NoError(t, err)
				assert.Equal(t, data.Duration, e.TotalTime())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestEngineAtTime(t *testing.T) ***REMOVED***
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	d := 50 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), d)
	startTime := time.Now()
	assert.NoError(t, e.Run(ctx))
	assert.WithinDuration(t, startTime.Add(d), startTime.Add(e.AtTime()), 2*TickRate)
***REMOVED***

func TestEngineSetPaused(t *testing.T) ***REMOVED***
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.False(t, e.IsPaused())

	e.SetPaused(true)
	assert.True(t, e.IsPaused())

	e.SetPaused(false)
	assert.False(t, e.IsPaused())
***REMOVED***

func TestEngineSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
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
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(-1), "vus-max can't be negative")
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set too low", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***
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
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Equal(t, int64(0), e.GetVUs())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(15)***REMOVED***)
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
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)

			e.numIterations = data.I
			e.numTaints = data.T
			assert.Equal(t, data.Expect, e.IsTainted())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestEngine_runVUOnceKeepsCounters(t *testing.T) ***REMOVED***
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), e.numIterations)
	assert.Equal(t, int64(0), e.numTaints)

	t.Run("success", func(t *testing.T) ***REMOVED***
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
		hook := logtest.NewGlobal()
		defer hook.Reset()

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

		err := hook.LastEntry().Data["error"].(error)
		assert.Equal(t, "this is an error", err.Error())
	***REMOVED***)
	t.Run("error/string", func(t *testing.T) ***REMOVED***
		hook := logtest.NewGlobal()
		defer hook.Reset()

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
	t.Run("taint", func(t *testing.T) ***REMOVED***
		hook := logtest.NewGlobal()
		defer hook.Reset()

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

		assert.Len(t, hook.Entries, 0)

	***REMOVED***)
***REMOVED***
