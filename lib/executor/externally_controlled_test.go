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

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
)

func getTestExternallyControlledConfig() ExternallyControlledConfig ***REMOVED***
	return ExternallyControlledConfig***REMOVED***
		ExternallyControlledConfigParams: ExternallyControlledConfigParams***REMOVED***
			VUs:      null.IntFrom(2),
			MaxVUs:   null.IntFrom(10),
			Duration: types.NullDurationFrom(2 * time.Second),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func TestExternallyControlledRun(t *testing.T) ***REMOVED***
	t.Parallel()

	doneIters := new(uint64)
	runner := simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
		time.Sleep(200 * time.Millisecond)
		atomic.AddUint64(doneIters, 1)
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", lib.Options***REMOVED******REMOVED***, runner, getTestExternallyControlledConfig())
	defer test.cancel()

	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, 1)
		doneCh = make(chan struct***REMOVED******REMOVED***)
	)
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		test.state.MarkStarted()
		errCh <- test.executor.Run(test.ctx, nil, lib.GetTestWorkerInfo())
		test.state.MarkEnded()
		close(doneCh)
	***REMOVED***()

	updateConfig := func(vus, maxVUs int64, errMsg string) ***REMOVED***
		newConfig := ExternallyControlledConfigParams***REMOVED***
			VUs:      null.IntFrom(vus),
			MaxVUs:   null.IntFrom(maxVUs),
			Duration: types.NullDurationFrom(2 * time.Second),
		***REMOVED***
		err := test.executor.(*ExternallyControlled).UpdateConfig(test.ctx, newConfig) //nolint:forcetypeassert
		if errMsg != "" ***REMOVED***
			assert.EqualError(t, err, errMsg)
		***REMOVED*** else ***REMOVED***
			assert.NoError(t, err)
		***REMOVED***
	***REMOVED***

	var resultVUCount [][]int64
	snapshot := func() ***REMOVED***
		resultVUCount = append(resultVUCount,
			[]int64***REMOVED***test.state.GetCurrentlyActiveVUsCount(), test.state.GetInitializedVUsCount()***REMOVED***)
	***REMOVED***

	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		snapshotTicker := time.NewTicker(500 * time.Millisecond)
		ticks := 0
		for ***REMOVED***
			select ***REMOVED***
			case <-snapshotTicker.C:
				snapshot()
				switch ticks ***REMOVED***
				case 0, 2:
					updateConfig(4, 10, "")
				case 1:
					updateConfig(8, 20, "")
				case 3:
					updateConfig(15, 10,
						"invalid configuration supplied: the number of active VUs (15)"+
							" must be less than or equal to the number of maxVUs (10)")
					updateConfig(-1, 10,
						"invalid configuration supplied: the number of VUs can't be negative")
				***REMOVED***
				ticks++
			case <-doneCh:
				snapshotTicker.Stop()
				snapshot()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	wg.Wait()
	require.NoError(t, <-errCh)
	assert.InDelta(t, 48, int(atomic.LoadUint64(doneIters)), 2)
	assert.Equal(t, [][]int64***REMOVED******REMOVED***2, 10***REMOVED***, ***REMOVED***4, 10***REMOVED***, ***REMOVED***8, 20***REMOVED***, ***REMOVED***4, 10***REMOVED***, ***REMOVED***0, 10***REMOVED******REMOVED***, resultVUCount)
***REMOVED***
