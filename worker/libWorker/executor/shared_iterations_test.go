package executor

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
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

	runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
		atomic.AddUint64(&doneIters, 1)
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, getTestSharedIterationsConfig())
	defer test.cancel()

	require.NoError(t, test.executor.Run(test.ctx, nil, libWorker.GetTestWorkerInfo()), libWorker.GetTestWorkerInfo())
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

	runner := simpleRunner(func(ctx context.Context, state *libWorker.State) error ***REMOVED***
		time.Sleep(10 * time.Millisecond) // small wait to stabilize the test
		// Pick one VU randomly and always slow it down.
		sid := atomic.LoadUint64(&slowVUID)
		if sid == uint64(0) ***REMOVED***
			atomic.StoreUint64(&slowVUID, state.VUID)
		***REMOVED***
		if sid == state.VUID ***REMOVED***
			time.Sleep(200 * time.Millisecond)
		***REMOVED***
		currIter, _ := result.LoadOrStore(state.VUID, uint64(0))
		result.Store(state.VUID, currIter.(uint64)+1) //nolint:forcetypeassert
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, getTestSharedIterationsConfig())
	defer test.cancel()

	require.NoError(t, test.executor.Run(test.ctx, nil, libWorker.GetTestWorkerInfo()))

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

	runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
		atomic.AddInt64(&count, 1)
		<-ctx.Done()
		return nil
	***REMOVED***)

	config := &SharedIterationsConfig***REMOVED***
		VUs:         null.IntFrom(5),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, config)
	defer test.cancel()

	engineOut := make(chan workerMetrics.SampleContainer, 1000)
	require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
	assert.Empty(t, test.logHook.Drain())
	assert.Equal(t, int64(5), count)
	assert.Equal(t, float64(95), sumMetricValues(engineOut, workerMetrics.DroppedIterationsName))
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
		expIters []uint64
	***REMOVED******REMOVED***
		***REMOVED***"0,1/4,3/4,1", "0:1/4", []uint64***REMOVED***1, 6, 11, 16, 21, 26, 31, 36, 41, 46***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "1/4:3/4", []uint64***REMOVED***0, 2, 4, 5, 7, 9, 10, 12, 14, 15, 17, 19, 20, 22, 24, 25, 27, 29, 30, 32, 34, 35, 37, 39, 40, 42, 44, 45, 47, 49***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "3/4:1", []uint64***REMOVED***3, 8, 13, 18, 23, 28, 33, 38, 43, 48***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%s_%s", tc.seq, tc.seg), func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotIters := []uint64***REMOVED******REMOVED***
			var mx sync.Mutex
			runner := simpleRunner(func(ctx context.Context, state *libWorker.State) error ***REMOVED***
				mx.Lock()
				gotIters = append(gotIters, state.GetScenarioGlobalVUIter())
				mx.Unlock()
				return nil
			***REMOVED***)

			test := setupExecutorTest(t, tc.seg, tc.seq, libWorker.Options***REMOVED******REMOVED***, runner, config)
			defer test.cancel()

			engineOut := make(chan workerMetrics.SampleContainer, 100)
			require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
			sort.Slice(gotIters, func(i, j int) bool ***REMOVED*** return gotIters[i] < gotIters[j] ***REMOVED***)
			assert.Equal(t, tc.expIters, gotIters)
		***REMOVED***)
	***REMOVED***
***REMOVED***
