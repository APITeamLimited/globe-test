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

	null "gopkg.in/guregu/null.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
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
	var doneIters uint64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestExternallyControlledConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(200 * time.Millisecond)
			atomic.AddUint64(&doneIters, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()

	var (
		wg            sync.WaitGroup
		errCh         = make(chan error, 1)
		doneCh        = make(chan struct***REMOVED******REMOVED***)
		resultVUCount []int64
	)
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		errCh <- executor.Run(ctx, nil)
		close(doneCh)
	***REMOVED***()

	updateConfig := func(vus int) ***REMOVED***
		newConfig := ExternallyControlledConfigParams***REMOVED***
			VUs:      null.IntFrom(int64(vus)),
			MaxVUs:   null.IntFrom(10),
			Duration: types.NullDurationFrom(2 * time.Second),
		***REMOVED***
		err := executor.(*ExternallyControlled).UpdateConfig(ctx, newConfig)
		assert.NoError(t, err)
	***REMOVED***

	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		es.MarkStarted()
		time.Sleep(150 * time.Millisecond) // wait for startup
		resultVUCount = append(resultVUCount, es.GetCurrentlyActiveVUsCount())
		time.Sleep(500 * time.Millisecond)
		updateConfig(4)
		time.Sleep(100 * time.Millisecond)
		resultVUCount = append(resultVUCount, es.GetCurrentlyActiveVUsCount())
		time.Sleep(500 * time.Millisecond)
		updateConfig(8)
		time.Sleep(100 * time.Millisecond)
		resultVUCount = append(resultVUCount, es.GetCurrentlyActiveVUsCount())
		time.Sleep(1 * time.Second)
		resultVUCount = append(resultVUCount, es.GetCurrentlyActiveVUsCount())
		es.MarkEnded()
	***REMOVED***()

	<-doneCh
	wg.Wait()
	require.NoError(t, <-errCh)
	assert.Equal(t, uint64(50), doneIters)
	assert.Equal(t, []int64***REMOVED***2, 4, 8, 0***REMOVED***, resultVUCount)
***REMOVED***
