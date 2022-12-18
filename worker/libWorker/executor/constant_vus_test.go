package executor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func getTestConstantVUsConfig() ConstantVUsConfig {
	return ConstantVUsConfig{
		BaseConfig: BaseConfig{GracefulStop: types.NullDurationFrom(100 * time.Millisecond)},
		VUs:        null.IntFrom(10),
		Duration:   types.NullDurationFrom(1 * time.Second),
	}
}

func TestConstantVUsRun(t *testing.T) {
	t.Parallel()
	var result sync.Map

	runner := simpleRunner(func(ctx context.Context, state *libWorker.State) error {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		currIter, _ := result.LoadOrStore(state.VUID, uint64(0))
		result.Store(state.VUID, currIter.(uint64)+1) //nolint:forcetypeassert
		time.Sleep(210 * time.Millisecond)
		return nil
	})

	test := setupExecutorTest(t, "", "", libWorker.Options{}, runner, getTestConstantVUsConfig())
	defer test.cancel()

	require.NoError(t, test.executor.Run(test.ctx, nil, libWorker.GetTestWorkerInfo()))

	var totalIters uint64
	result.Range(func(key, value interface{}) bool {
		vuIters := value.(uint64)
		assert.Equal(t, uint64(5), vuIters)
		totalIters += vuIters
		return true
	})
	assert.Equal(t, uint64(50), totalIters)
}
