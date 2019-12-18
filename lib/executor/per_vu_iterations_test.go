package executor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
)

func getTestPerVUIterationsConfig() PerVUIterationsConfig ***REMOVED***
	return PerVUIterationsConfig***REMOVED***
		VUs:         null.IntFrom(10),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(5 * time.Second),
	***REMOVED***
***REMOVED***

func TestPerVUIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	var result sync.Map
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
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
	err := executor.Run(ctx, nil)
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
