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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
)

func getTestConstantLoopingVUsConfig() ConstantLoopingVUsConfig ***REMOVED***
	return ConstantLoopingVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(1 * time.Second)***REMOVED***,
		VUs:        null.IntFrom(10),
		Duration:   types.NullDurationFrom(1 * time.Second),
	***REMOVED***
***REMOVED***

func TestConstantLoopingVUsRun(t *testing.T) ***REMOVED***
	t.Parallel()
	var result sync.Map
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, getTestConstantLoopingVUsConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(200 * time.Millisecond)
			select ***REMOVED***
			case <-ctx.Done():
				return nil
			default:
			***REMOVED***
			state := lib.GetState(ctx)
			currIter, _ := result.LoadOrStore(state.Vu, uint64(0))
			result.Store(state.Vu, currIter.(uint64)+1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	err = executor.Run(ctx, nil)
	require.NoError(t, err)

	var totalIters uint64
	result.Range(func(key, value interface***REMOVED******REMOVED***) bool ***REMOVED***
		vuIters := value.(uint64)
		assert.Equal(t, uint64(5), vuIters)
		totalIters += vuIters
		return true
	***REMOVED***)
	assert.Equal(t, uint64(50), totalIters)
***REMOVED***
