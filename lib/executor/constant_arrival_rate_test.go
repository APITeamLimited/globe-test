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

package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
)

func newExecutionSegmentFromString(str string) *lib.ExecutionSegment ***REMOVED***
	r, err := lib.NewExecutionSegmentFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r
***REMOVED***

func newExecutionSegmentSequenceFromString(str string) *lib.ExecutionSegmentSequence ***REMOVED***
	r, err := lib.NewExecutionSegmentSequenceFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return &r
***REMOVED***

func getTestConstantArrivalRateConfig() *ConstantArrivalRateConfig ***REMOVED***
	return &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(1 * time.Second)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(50),
		Duration:        types.NullDurationFrom(5 * time.Second),
		PreAllocatedVUs: null.IntFrom(10),
		MaxVUs:          null.IntFrom(20),
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunNotEnoughAllocatedVUsWarn(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	ctx, cancel, executor, logHook := setupExecutor(
		t, getTestConstantArrivalRateConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(time.Second)
			return nil
		***REMOVED***),
	)
	defer cancel()
	engineOut := make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	require.NoError(t, err)
	entries := logHook.Drain()
	require.NotEmpty(t, entries)
	for _, entry := range entries ***REMOVED***
		require.Equal(t,
			"Insufficient VUs, reached 20 active VUs and cannot initialize more",
			entry.Message)
		require.Equal(t, logrus.WarnLevel, entry.Level)
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunCorrectRate(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	ctx, cancel, executor, logHook := setupExecutor(
		t, getTestConstantArrivalRateConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			atomic.AddInt64(&count, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		// check that we got around the amount of VU iterations as we would expect
		var currentCount int64

		for i := 0; i < 5; i++ ***REMOVED***
			time.Sleep(time.Second)
			currentCount = atomic.SwapInt64(&count, 0)
			require.InDelta(t, 50, currentCount, 1)
		***REMOVED***
	***REMOVED***()
	engineOut := make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	wg.Wait()
	require.NoError(t, err)
	require.Empty(t, logHook.Drain())
***REMOVED***

//nolint:tparallel,paralleltest // this is flaky if ran with other tests
func TestConstantArrivalRateRunCorrectTiming(t *testing.T) ***REMOVED***
	// t.Parallel()
	tests := []struct ***REMOVED***
		segment  *lib.ExecutionSegment
		sequence *lib.ExecutionSegmentSequence
		start    time.Duration
		steps    []int64
	***REMOVED******REMOVED***
		***REMOVED***
			segment: newExecutionSegmentFromString("0:1/3"),
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 60, 60, 60, 60, 60, 60***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("1/3:2/3"),
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("2/3:1"),
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 60, 60, 60, 60, 60, 60***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("1/6:3/6"),
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 80, 40, 80, 40, 80, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("1/6:3/6"),
			sequence: newExecutionSegmentSequenceFromString("1/6,3/6"),
			start:    time.Millisecond * 20,
			steps:    []int64***REMOVED***40, 80, 40, 80, 40, 80, 40***REMOVED***,
		***REMOVED***,
		// sequences
		***REMOVED***
			segment:  newExecutionSegmentFromString("0:1/3"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 0,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("1/3:2/3"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 20,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("2/3:1"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 40,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 100***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		test := test

		t.Run(fmt.Sprintf("segment %s sequence %s", test.segment, test.sequence), func(t *testing.T) ***REMOVED***
			t.Parallel()
			et, err := lib.NewExecutionTuple(test.segment, test.sequence)
			require.NoError(t, err)
			es := lib.NewExecutionState(lib.Options***REMOVED***
				ExecutionSegment:         test.segment,
				ExecutionSegmentSequence: test.sequence,
			***REMOVED***, et, 10, 50)
			var count int64
			seconds := 2
			config := getTestConstantArrivalRateConfig()
			config.Duration.Duration = types.Duration(time.Second * time.Duration(seconds))
			newET, err := es.ExecutionTuple.GetNewExecutionTupleFromValue(config.MaxVUs.Int64)
			require.NoError(t, err)
			rateScaled := newET.ScaleInt64(config.Rate.Int64)
			startTime := time.Now()
			expectedTimeInt64 := int64(test.start)
			ctx, cancel, executor, logHook := setupExecutor(
				t, config, es,
				simpleRunner(func(ctx context.Context) error ***REMOVED***
					current := atomic.AddInt64(&count, 1)

					expectedTime := test.start
					if current != 1 ***REMOVED***
						expectedTime = time.Duration(atomic.AddInt64(&expectedTimeInt64,
							int64(time.Millisecond)*test.steps[(current-2)%int64(len(test.steps))]))
					***REMOVED***
					assert.WithinDuration(t,
						startTime.Add(expectedTime),
						time.Now(),
						time.Millisecond*12,
						"%d expectedTime %s", current, expectedTime,
					)

					return nil
				***REMOVED***),
			)

			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				// check that we got around the amount of VU iterations as we would expect
				var currentCount int64

				for i := 0; i < seconds; i++ ***REMOVED***
					time.Sleep(time.Second)
					currentCount = atomic.LoadInt64(&count)
					assert.InDelta(t, int64(i+1)*rateScaled, currentCount, 3)
				***REMOVED***
			***REMOVED***()
			startTime = time.Now()
			engineOut := make(chan stats.SampleContainer, 1000)
			err = executor.Run(ctx, engineOut)
			wg.Wait()
			require.NoError(t, err)
			require.Empty(t, logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArrivalRateCancel(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := map[string]lib.ExecutorConfig***REMOVED***
		"constant": getTestConstantArrivalRateConfig(),
		"ramping":  getTestRampingArrivalRateConfig(),
	***REMOVED***
	for name, config := range testCases ***REMOVED***
		config := config
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			ch := make(chan struct***REMOVED******REMOVED***)
			errCh := make(chan error, 1)
			weAreDoneCh := make(chan struct***REMOVED******REMOVED***)
			et, err := lib.NewExecutionTuple(nil, nil)
			require.NoError(t, err)
			es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
			ctx, cancel, executor, logHook := setupExecutor(
				t, config, es, simpleRunner(func(ctx context.Context) error ***REMOVED***
					select ***REMOVED***
					case <-ch:
						<-ch
					default:
					***REMOVED***
					return nil
				***REMOVED***))
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()

				engineOut := make(chan stats.SampleContainer, 1000)
				errCh <- executor.Run(ctx, engineOut)
				close(weAreDoneCh)
			***REMOVED***()

			time.Sleep(time.Second)
			ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			cancel()
			time.Sleep(time.Second)
			select ***REMOVED***
			case <-weAreDoneCh:
				t.Fatal("Run returned before all VU iterations were finished")
			default:
			***REMOVED***
			close(ch)
			<-weAreDoneCh
			wg.Wait()
			require.NoError(t, <-errCh)
			require.Empty(t, logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateDroppedIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)

	config := &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(10),
		Duration:        types.NullDurationFrom(950 * time.Millisecond),
		PreAllocatedVUs: null.IntFrom(5),
		MaxVUs:          null.IntFrom(5),
	***REMOVED***

	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	ctx, cancel, executor, logHook := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			atomic.AddInt64(&count, 1)
			<-ctx.Done()
			return nil
		***REMOVED***),
	)
	defer cancel()
	engineOut := make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	require.NoError(t, err)
	logs := logHook.Drain()
	require.Len(t, logs, 1)
	assert.Contains(t, logs[0].Message, "cannot initialize more")
	assert.Equal(t, int64(5), count)
	assert.Equal(t, float64(5), sumMetricValues(engineOut, metrics.DroppedIterations.Name))
***REMOVED***

func TestConstantArrivalRateGlobalIters(t *testing.T) ***REMOVED***
	t.Parallel()

	config := &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(100 * time.Millisecond)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(950 * time.Millisecond),
		Rate:            null.IntFrom(20),
		Duration:        types.NullDurationFrom(1 * time.Second),
		PreAllocatedVUs: null.IntFrom(5),
		MaxVUs:          null.IntFrom(5),
	***REMOVED***

	testCases := []struct ***REMOVED***
		seq, seg string
		expIters []uint64
	***REMOVED******REMOVED***
		***REMOVED***"0,1/4,3/4,1", "0:1/4", []uint64***REMOVED***1, 6, 11, 16, 21***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "1/4:3/4", []uint64***REMOVED***0, 2, 4, 5, 7, 9, 10, 12, 14, 15, 17, 19, 20***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "3/4:1", []uint64***REMOVED***3, 8, 13, 18***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%s_%s", tc.seq, tc.seg), func(t *testing.T) ***REMOVED***
			t.Parallel()
			ess, err := lib.NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(t, err)
			seg, err := lib.NewExecutionSegmentFromString(tc.seg)
			require.NoError(t, err)
			et, err := lib.NewExecutionTuple(seg, &ess)
			require.NoError(t, err)
			es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 5, 5)

			runner := &minirunner.MiniRunner***REMOVED******REMOVED***
			ctx, cancel, executor, _ := setupExecutor(t, config, es, runner)
			defer cancel()

			gotIters := []uint64***REMOVED******REMOVED***
			var mx sync.Mutex
			runner.Fn = func(ctx context.Context, _ chan<- stats.SampleContainer) error ***REMOVED***
				state := lib.GetState(ctx)
				mx.Lock()
				gotIters = append(gotIters, state.GetScenarioGlobalVUIter())
				mx.Unlock()
				return nil
			***REMOVED***

			engineOut := make(chan stats.SampleContainer, 100)
			err = executor.Run(ctx, engineOut)
			require.NoError(t, err)
			assert.Equal(t, tc.expIters, gotIters)
		***REMOVED***)
	***REMOVED***
***REMOVED***
