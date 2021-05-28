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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
)

func getTestSharedIterationsConfig() SharedIterationsConfig ***REMOVED***
	return SharedIterationsConfig***REMOVED***
		VUs:         null.IntFrom(10),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(5 * time.Second),
	***REMOVED***
***REMOVED***

// Baseline test
func TestSharedIterationsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var doneIters uint64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestSharedIterationsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			atomic.AddUint64(&doneIters, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	err = executor.Run(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), doneIters)
***REMOVED***

// Test that when one VU "slows down", others will pick up the workload.
// This is the reverse behavior of the PerVUIterations executor.
func TestSharedIterationsRunVariableVU(t *testing.T) ***REMOVED***
	t.Parallel()
	var (
		result   sync.Map
		slowVUID uint64
	)
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestSharedIterationsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(10 * time.Millisecond) // small wait to stabilize the test
			state := lib.GetState(ctx)
			// Pick one VU randomly and always slow it down.
			sid := atomic.LoadUint64(&slowVUID)
			if sid == uint64(0) ***REMOVED***
				atomic.StoreUint64(&slowVUID, state.Vu)
			***REMOVED***
			if sid == state.Vu ***REMOVED***
				time.Sleep(200 * time.Millisecond)
			***REMOVED***
			currIter, _ := result.LoadOrStore(state.Vu, uint64(0))
			result.Store(state.Vu, currIter.(uint64)+1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	err = executor.Run(ctx, nil)
	require.NoError(t, err)

	var totalIters uint64
	result.Range(func(key, value interface***REMOVED******REMOVED***) bool ***REMOVED***
		totalIters += value.(uint64)
		return true
	***REMOVED***)

	// The slow VU should complete 2 iterations given these timings,
	// while the rest should randomly complete the other 98 iterations.
	val, ok := result.Load(slowVUID)
	assert.True(t, ok)
	assert.Equal(t, uint64(2), val)
	assert.Equal(t, uint64(100), totalIters)
***REMOVED***

func TestSharedIterationsEmitDroppedIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)

	config := &SharedIterationsConfig***REMOVED***
		VUs:         null.IntFrom(5),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(1 * time.Second),
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
	assert.Empty(t, logHook.Drain())
	assert.Equal(t, int64(5), count)
	assert.Equal(t, float64(95), sumMetricValues(engineOut, metrics.DroppedIterations.Name))
***REMOVED***

func TestSharedIterationsGlobalIters(t *testing.T) ***REMOVED***
	t.Parallel()

	config := &SharedIterationsConfig***REMOVED***
		VUs:         null.IntFrom(5),
		Iterations:  null.IntFrom(50),
		MaxDuration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	testCases := []struct ***REMOVED***
		seq, seg string
		expIters []int64
	***REMOVED******REMOVED***
		***REMOVED***"0,1/4,3/4,1", "0:1/4", []int64***REMOVED***1, 6, 11, 16, 21, 26, 31, 36, 41, 46***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "1/4:3/4", []int64***REMOVED***0, 2, 4, 5, 7, 9, 10, 12, 14, 15, 17, 19, 20, 22, 24, 25, 27, 29, 30, 32, 34, 35, 37, 39, 40, 42, 44, 45, 47, 49***REMOVED******REMOVED***,
		// FIXME: The skewed values are because of the time.Sleep() in the
		// VU function below.
		// ***REMOVED***"0,1/4,3/4,1", "1/4:3/4", []int64***REMOVED***4, 5, 7, 9, 10, 12, 14, 15, 17, 19, 20, 22, 24, 25, 27, 29, 30, 32, 34, 35, 37, 39, 40, 42, 44, 45, 47, 49, 49, 49***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "3/4:1", []int64***REMOVED***3, 8, 13, 18, 23, 28, 33, 38, 43, 48***REMOVED******REMOVED***,
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

			gotIters := []int64***REMOVED******REMOVED***
			var mx sync.Mutex
			runner.Fn = func(ctx context.Context, _ chan<- stats.SampleContainer) error ***REMOVED***
				mx.Lock()
				gotIters = append(gotIters, executor.(*SharedIterations).getGlobalIter())
				// FIXME: This delay minimizes chances of flakiness, but
				// produces skewed values. :-/
				// time.Sleep(10 * time.Millisecond)
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
