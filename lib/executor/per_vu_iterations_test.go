package executor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
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
	doneIters := uint64(0)
	var ctx, cancel, executor, logHook = setupExecutor(
		t, getTestPerVUIterationsConfig(),
		func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			atomic.AddUint64(&doneIters, 1)
			return nil
		***REMOVED***,
		[]logrus.Level***REMOVED***logrus.InfoLevel***REMOVED***,
	)
	defer cancel()
	err := executor.Run(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), doneIters)

	entries := logHook.Drain()
	require.NotEmpty(t, entries)
	result := map[int64]uint64***REMOVED******REMOVED***
	for _, entry := range entries ***REMOVED***
		vuID := entry.Data["vu_id"].(int64)
		result[vuID]++
	***REMOVED***
	assert.Equal(t, 10, len(result))
	for _, vuIterCount := range result ***REMOVED***
		assert.Equal(t, uint64(100), vuIterCount)
	***REMOVED***
***REMOVED***
