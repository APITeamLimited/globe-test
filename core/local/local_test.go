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

package local

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestExecutorRun(t *testing.T) ***REMOVED***
	e := New(nil)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))

	ctx, cancel := context.WithCancel(context.Background())
	err := make(chan error, 1)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	cancel()
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetLogger(t *testing.T) ***REMOVED***
	logger, _ := logtest.NewNullLogger()
	e := New(nil)
	e.SetLogger(logger)
	assert.Equal(t, logger, e.GetLogger())
***REMOVED***

func TestExecutorStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Duration time.Duration
		Stages   []lib.Stage
	***REMOVED******REMOVED***
		"one": ***REMOVED***
			1 * time.Second,
			[]lib.Stage***REMOVED******REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED******REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/targeted": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			e := New(nil)
			assert.NoError(t, e.SetVUsMax(10))
			e.SetStages(data.Stages)
			assert.NoError(t, e.Run(context.Background(), nil))
			assert.True(t, e.GetTime() >= data.Duration)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutorEndTime(t *testing.T) ***REMOVED***
	e := New(nil)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))
	e.SetEndTime(types.NullDurationFrom(1 * time.Second))
	assert.Equal(t, types.NullDurationFrom(1*time.Second), e.GetEndTime())

	startTime := time.Now()
	assert.NoError(t, e.Run(context.Background(), nil))
	assert.True(t, time.Now().After(startTime.Add(1*time.Second)), "test did not take 1s")

	t.Run("Runtime Errors", func(t *testing.T) ***REMOVED***
		e := New(lib.MiniRunner***REMOVED***Fn: func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			return nil, errors.New("hi")
		***REMOVED******REMOVED***)
		assert.NoError(t, e.SetVUsMax(10))
		assert.NoError(t, e.SetVUs(10))
		e.SetEndTime(types.NullDurationFrom(100 * time.Millisecond))
		assert.Equal(t, types.NullDurationFrom(100*time.Millisecond), e.GetEndTime())

		l, hook := logtest.NewNullLogger()
		e.SetLogger(l)

		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background(), nil))
		assert.True(t, time.Now().After(startTime.Add(100*time.Millisecond)), "test did not take 100ms")

		assert.NotEmpty(t, hook.Entries)
		for _, e := range hook.Entries ***REMOVED***
			assert.Equal(t, "hi", e.Message)
		***REMOVED***
	***REMOVED***)

	t.Run("End Errors", func(t *testing.T) ***REMOVED***
		e := New(lib.MiniRunner***REMOVED***Fn: func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			<-ctx.Done()
			return nil, errors.New("hi")
		***REMOVED******REMOVED***)
		assert.NoError(t, e.SetVUsMax(10))
		assert.NoError(t, e.SetVUs(10))
		e.SetEndTime(types.NullDurationFrom(100 * time.Millisecond))
		assert.Equal(t, types.NullDurationFrom(100*time.Millisecond), e.GetEndTime())

		l, hook := logtest.NewNullLogger()
		e.SetLogger(l)

		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background(), nil))
		assert.True(t, time.Now().After(startTime.Add(100*time.Millisecond)), "test did not take 100ms")

		assert.Empty(t, hook.Entries)
	***REMOVED***)
***REMOVED***

func TestExecutorEndIterations(t *testing.T) ***REMOVED***
	metric := &stats.Metric***REMOVED***Name: "test_metric"***REMOVED***

	var i int64
	e := New(lib.MiniRunner***REMOVED***Fn: func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		default:
			atomic.AddInt64(&i, 1)
		***REMOVED***
		return []stats.Sample***REMOVED******REMOVED***Metric: metric, Value: 1.0***REMOVED******REMOVED***, nil
	***REMOVED******REMOVED***)
	assert.NoError(t, e.SetVUsMax(1))
	assert.NoError(t, e.SetVUs(1))
	e.SetEndIterations(null.IntFrom(100))
	assert.Equal(t, null.IntFrom(100), e.GetEndIterations())

	samples := make(chan []stats.Sample, 101)
	assert.NoError(t, e.Run(context.Background(), samples))
	assert.Equal(t, int64(100), e.GetIterations())
	assert.Equal(t, int64(100), i)

	for i := 0; i < 100; i++ ***REMOVED***
		samples := <-samples
		if assert.Len(t, samples, 2) ***REMOVED***
			assert.Equal(t, stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***, samples[0])
			assert.Equal(t, metrics.Iterations, samples[1].Metric)
			assert.Equal(t, float64(1), samples[1].Value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestExecutorIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e := New(nil)

	err := make(chan error)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	for !e.IsRunning() ***REMOVED***
	***REMOVED***
	cancel()
	for e.IsRunning() ***REMOVED***
	***REMOVED***
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUsMax(-1), "vu cap can't be negative")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(nil)

		assert.NoError(t, e.SetVUsMax(50))
		assert.Equal(t, int64(50), e.GetVUsMax())

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(50))
			assert.Equal(t, int64(50), e.GetVUsMax())
		***REMOVED***)
	***REMOVED***)

	t.Run("TooLow", func(t *testing.T) ***REMOVED***
		e := New(nil)
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())

		assert.EqualError(t, e.SetVUsMax(50), "can't lower vu cap (to 50) below vu count (100)")
	***REMOVED***)
***REMOVED***

func TestExecutorSetVUs(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(-1), "vu count can't be negative")
	***REMOVED***)

	t.Run("Too High", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(100), "can't raise vu count (to 100) above vu cap (0)")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(lib.MiniRunner***REMOVED***Fn: func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			return nil, nil
		***REMOVED******REMOVED***)
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				num++
				if assert.NotNil(t, handle.vu, "vu %d lacks impl", i) ***REMOVED***
					assert.Equal(t, int64(0), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
				assert.Nil(t, handle.ctx, "vu %d has ctx", i)
				assert.Nil(t, handle.cancel, "vu %d has cancel", i)
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(50))
		assert.Equal(t, int64(50), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				if i < 50 ***REMOVED***
					assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
					assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
					num++
				***REMOVED*** else ***REMOVED***
					assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					assert.Equal(t, int64(0), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
			***REMOVED***
			assert.Equal(t, 50, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
				assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
				num++
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(50))
			assert.Equal(t, int64(50), e.GetVUs())
			if assert.Len(t, e.vus, 100) ***REMOVED***
				num := 0
				for i, handle := range e.vus ***REMOVED***
					if i < 50 ***REMOVED***
						assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
						num++
					***REMOVED*** else ***REMOVED***
						assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					***REMOVED***
					assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
				assert.Equal(t, 50, num)
			***REMOVED***

			t.Run("Raise", func(t *testing.T) ***REMOVED***
				assert.NoError(t, e.SetVUs(100))
				assert.Equal(t, int64(100), e.GetVUs())
				if assert.Len(t, e.vus, 100) ***REMOVED***
					for i, handle := range e.vus ***REMOVED***
						assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
						if i < 50 ***REMOVED***
							assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
						***REMOVED*** else ***REMOVED***
							assert.Equal(t, int64(50+i+1), handle.vu.(*lib.MiniRunnerVU).ID)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
