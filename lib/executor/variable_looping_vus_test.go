package executor

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
)

func getTestVariableLoopingVUsConfig() VariableLoopingVUsConfig ***REMOVED***
	return VariableLoopingVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(0)***REMOVED***,
		StartVUs:   null.IntFrom(5),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(5),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(0),
				Target:   null.IntFrom(3),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(3),
			***REMOVED***,
		***REMOVED***,
		GracefulRampDown: types.NullDurationFrom(0),
	***REMOVED***
***REMOVED***

func TestVariableLoopingVUsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var iterCount int64
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestVariableLoopingVUsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(200 * time.Millisecond)
			atomic.AddInt64(&iterCount, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()

	var (
		wg     sync.WaitGroup
		result []int64
	)

	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		result = append(result, es.GetCurrentlyActiveVUsCount())
		time.Sleep(1 * time.Second)
		result = append(result, es.GetCurrentlyActiveVUsCount())
		time.Sleep(1 * time.Second)
		result = append(result, es.GetCurrentlyActiveVUsCount())
	***REMOVED***()

	err := executor.Run(ctx, nil)

	wg.Wait()
	require.NoError(t, err)
	assert.Equal(t, []int64***REMOVED***5, 3, 0***REMOVED***, result)
	assert.Equal(t, int64(40), iterCount)
***REMOVED***
