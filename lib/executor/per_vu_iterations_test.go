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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

func getTestPerVUIterationsConfig() PerVUIterationsConfig ***REMOVED***
	return PerVUIterationsConfig***REMOVED***
		BaseConfig:  BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(1 * time.Second)***REMOVED***,
		VUs:         null.IntFrom(10),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(3 * time.Second),
	***REMOVED***
***REMOVED***

// Baseline test
func TestPerVUIterationsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var result sync.Map
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestPerVUIterationsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			state := lib.GetState(ctx)
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
		vuIters := value.(uint64)
		assert.Equal(t, uint64(100), vuIters)
		totalIters += vuIters
		return true
	***REMOVED***)
	assert.Equal(t, uint64(1000), totalIters)
***REMOVED***

// Test that when one VU "slows down", others will *not* pick up the workload.
// This is the reverse behavior of the SharedIterations executor.
func TestPerVUIterationsRunVariableVU(t *testing.T) ***REMOVED***
	t.Parallel()
	var (
		result   sync.Map
		slowVUID int64
	)
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestPerVUIterationsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			state := lib.GetState(ctx)
			// Pick one VU randomly and always slow it down.
			sid := atomic.LoadInt64(&slowVUID)
			if sid == int64(0) ***REMOVED***
				atomic.StoreInt64(&slowVUID, state.Vu)
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

	val, ok := result.Load(slowVUID)
	assert.True(t, ok)

	var totalIters uint64
	result.Range(func(key, value interface***REMOVED******REMOVED***) bool ***REMOVED***
		vuIters := value.(uint64)
		if key != slowVUID ***REMOVED***
			assert.Equal(t, uint64(100), vuIters)
		***REMOVED***
		totalIters += vuIters
		return true
	***REMOVED***)

	// The slow VU should complete 16 iterations given these timings,
	// while the rest should equally complete their assigned 100 iterations.
	assert.Equal(t, uint64(16), val)
	assert.Equal(t, uint64(916), totalIters)
***REMOVED***

func TestPerVuIterationsEmitDroppedIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)

	config := PerVUIterationsConfig***REMOVED***
		VUs:         null.IntFrom(5),
		Iterations:  null.IntFrom(20),
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