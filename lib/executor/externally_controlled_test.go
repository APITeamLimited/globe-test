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

	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)

	doneIters := new(uint64)
	ctx, cancel, executor, _ := setupExecutor(
		t, getTestExternallyControlledConfig(), es,
		simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
			time.Sleep(200 * time.Millisecond)
			atomic.AddUint64(doneIters, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()

	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, 1)
		doneCh = make(chan struct***REMOVED******REMOVED***)
	)
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		es.MarkStarted()
		errCh <- executor.Run(ctx, nil)
		es.MarkEnded()
		close(doneCh)
	***REMOVED***()

	updateConfig := func(vus, maxVUs int64, errMsg string) ***REMOVED***
		newConfig := ExternallyControlledConfigParams***REMOVED***
			VUs:      null.IntFrom(vus),
			MaxVUs:   null.IntFrom(maxVUs),
			Duration: types.NullDurationFrom(2 * time.Second),
		***REMOVED***
		err := executor.(*ExternallyControlled).UpdateConfig(ctx, newConfig)
		if errMsg != "" ***REMOVED***
			assert.EqualError(t, err, errMsg)
		***REMOVED*** else ***REMOVED***
			assert.NoError(t, err)
		***REMOVED***
	***REMOVED***

	var resultVUCount [][]int64
	snapshot := func() ***REMOVED***
		resultVUCount = append(resultVUCount,
			[]int64***REMOVED***es.GetCurrentlyActiveVUsCount(), es.GetInitializedVUsCount()***REMOVED***)
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
