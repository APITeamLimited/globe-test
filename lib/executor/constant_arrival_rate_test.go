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
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

func getTestConstantArrivalRateConfig() ConstantArrivalRateConfig ***REMOVED***
	return ConstantArrivalRateConfig***REMOVED***
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(50),
		Duration:        types.NullDurationFrom(5 * time.Second),
		PreAllocatedVUs: null.IntFrom(10),
		MaxVUs:          null.IntFrom(20),
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunNotEnoughAllocatedVUsWarn(t *testing.T) ***REMOVED***
	t.Parallel()
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, logHook = setupExecutor(
		t, getTestConstantArrivalRateConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(time.Second)
			return nil
		***REMOVED***),
	)
	defer cancel()
	var engineOut = make(chan stats.SampleContainer, 1000)
	err := executor.Run(ctx, engineOut)
	require.NoError(t, err)
	entries := logHook.Drain()
	require.NotEmpty(t, entries)
	for _, entry := range entries ***REMOVED***
		require.Equal(t,
			"Insufficient VUs, reached 20 active VUs and cannot allocate more",
			entry.Message)
		require.Equal(t, logrus.WarnLevel, entry.Level)
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunCorrectRate(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, logHook = setupExecutor(
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
	var engineOut = make(chan stats.SampleContainer, 1000)
	err := executor.Run(ctx, engineOut)
	wg.Wait()
	require.NoError(t, err)
	require.Empty(t, logHook.Drain())
***REMOVED***

func TestConstantArrivalRateRunCorrectTiming(t *testing.T) ***REMOVED***
	newExecutionSegmentFromString := func(str string) *lib.ExecutionSegment ***REMOVED***
		r, err := lib.NewExecutionSegmentFromString(str)
		require.NoError(t, err)
		return r
	***REMOVED***

	newExecutionSegmentSequenceFromString := func(str string) *lib.ExecutionSegmentSequence ***REMOVED***
		r, err := lib.NewExecutionSegmentSequenceFromString(str)
		require.NoError(t, err)
		return &r
	***REMOVED***

	var tests = []struct ***REMOVED***
		segment  *lib.ExecutionSegment
		sequence *lib.ExecutionSegmentSequence
		start    time.Duration
		step     time.Duration
	***REMOVED******REMOVED***
		***REMOVED***
			segment: newExecutionSegmentFromString("0:1/3"),
			start:   time.Millisecond * 00,
			step:    time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("1/3:2/3"),
			start:   time.Millisecond * 00,
			step:    time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("2/3:1"),
			start:   time.Millisecond * 00,
			step:    time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment: newExecutionSegmentFromString("1/6:3/6"),
			start:   time.Millisecond * 00,
			step:    time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("1/6:3/6"),
			sequence: &lib.ExecutionSegmentSequence***REMOVED******REMOVED***,
			start:    time.Millisecond * 00,
			step:     time.Millisecond * 60,
		***REMOVED***,
		// sequences
		***REMOVED***
			segment:  newExecutionSegmentFromString("0:1/3"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 00,
			step:     time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("1/3:2/3"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 20,
			step:     time.Millisecond * 60,
		***REMOVED***,
		***REMOVED***
			segment:  newExecutionSegmentFromString("2/3:1"),
			sequence: newExecutionSegmentSequenceFromString("0,1/3,2/3,1"),
			start:    time.Millisecond * 40,
			step:     time.Millisecond * 60,
		***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		test := test

		t.Run(fmt.Sprintf("segment %s sequence %s", test.segment, test.sequence), func(t *testing.T) ***REMOVED***
			t.Parallel()
			es := lib.NewExecutionState(lib.Options***REMOVED***
				ExecutionSegment: test.segment,
				ESS:              test.sequence,
			***REMOVED***, 10, 50)
			var count int64
			var startTime = time.Now()
			var ctx, cancel, executor, logHook = setupExecutor(
				t, getTestConstantArrivalRateConfig(), es,
				simpleRunner(func(ctx context.Context) error ***REMOVED***
					current := atomic.AddInt64(&count, 1)
					expectedTime := test.start + time.Duration(current-1)*test.step
					assert.WithinDuration(t,
						startTime.Add(expectedTime),
						time.Now(),
						time.Millisecond*10,
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

				for i := 0; i < 5; i++ ***REMOVED***
					time.Sleep(time.Second)
					currentCount = atomic.LoadInt64(&count)
					assert.InDelta(t, (i+1)*17, currentCount, 2)
				***REMOVED***
			***REMOVED***()
			startTime = time.Now()
			var engineOut = make(chan stats.SampleContainer, 1000)
			err := executor.Run(ctx, engineOut)
			wg.Wait()
			require.NoError(t, err)
			require.Empty(t, logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArrivalRateCancel(t *testing.T) ***REMOVED***
	t.Parallel()

	var testCases = map[string]lib.ExecutorConfig***REMOVED***
		"constant": getTestConstantArrivalRateConfig(),
		"variable": getTestVariableArrivalRateConfig(),
	***REMOVED***
	for name, config := range testCases ***REMOVED***
		config := config
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			var ch = make(chan struct***REMOVED******REMOVED***)
			var errCh = make(chan error, 1)
			var weAreDoneCh = make(chan struct***REMOVED******REMOVED***)
			es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
			var ctx, cancel, executor, logHook = setupExecutor(
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

				var engineOut = make(chan stats.SampleContainer, 1000)
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
