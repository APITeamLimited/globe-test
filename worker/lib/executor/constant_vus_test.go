package executor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/types"
)

func getTestConstantVUsConfig() ConstantVUsConfig ***REMOVED***
	return ConstantVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(100 * time.Millisecond)***REMOVED***,
		VUs:        null.IntFrom(10),
		Duration:   types.NullDurationFrom(1 * time.Second),
	***REMOVED***
***REMOVED***

func TestConstantVUsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var result sync.Map

	runner := simpleRunner(func(ctx context.Context, state *lib.State) error ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		default:
		***REMOVED***
		currIter, _ := result.LoadOrStore(state.VUID, uint64(0))
		result.Store(state.VUID, currIter.(uint64)+1) //nolint:forcetypeassert
		time.Sleep(210 * time.Millisecond)
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", lib.Options***REMOVED******REMOVED***, runner, getTestConstantVUsConfig())
	defer test.cancel()

	require.NoError(t, test.executor.Run(test.ctx, nil, lib.GetTestWorkerInfo()))

	var totalIters uint64
	result.Range(func(key, value interface***REMOVED******REMOVED***) bool ***REMOVED***
		vuIters := value.(uint64)
		assert.Equal(t, uint64(5), vuIters)
		totalIters += vuIters
		return true
	***REMOVED***)
	assert.Equal(t, uint64(50), totalIters)
***REMOVED***
