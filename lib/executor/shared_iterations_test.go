package executor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
)

func getTestSharedIterationsConfig() SharedIterationsConfig ***REMOVED***
	return SharedIterationsConfig***REMOVED***
		VUs:         null.IntFrom(10),
		Iterations:  null.IntFrom(100),
		MaxDuration: types.NullDurationFrom(5 * time.Second),
	***REMOVED***
***REMOVED***

func TestSharedIterationsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var doneIters uint64
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestSharedIterationsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			atomic.AddUint64(&doneIters, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	err := executor.Run(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), doneIters)
***REMOVED***
