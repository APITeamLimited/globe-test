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
	"fmt"
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

func TestVariableLoopingVUsRun(t *testing.T) ***REMOVED***
	t.Parallel()

	config := VariableLoopingVUsConfig***REMOVED***
		BaseConfig:       BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(0)***REMOVED***,
		GracefulRampDown: types.NullDurationFrom(0),
		StartVUs:         null.IntFrom(5),
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
	***REMOVED***

	var iterCount int64
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, config, es,
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

// Ensure there's no wobble of VUs during graceful ramp-down, without segments.
// See https://github.com/loadimpact/k6/issues/1296
func TestVariableLoopingVUsRampDownNoWobble(t *testing.T) ***REMOVED***
	t.Parallel()

	config := VariableLoopingVUsConfig***REMOVED***
		BaseConfig:       BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(0)***REMOVED***,
		GracefulRampDown: types.NullDurationFrom(1 * time.Second),
		StartVUs:         null.IntFrom(0),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(3 * time.Second),
				Target:   null.IntFrom(10),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(2 * time.Second),
				Target:   null.IntFrom(0),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(1 * time.Second)
			return nil
		***REMOVED***),
	)
	defer cancel()

	var (
		wg     sync.WaitGroup
		result []int64
		m      sync.Mutex
	)

	sampleActiveVUs := func(delay time.Duration) ***REMOVED***
		time.Sleep(delay)
		m.Lock()
		result = append(result, es.GetCurrentlyActiveVUsCount())
		m.Unlock()
	***REMOVED***

	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		sampleActiveVUs(100 * time.Millisecond)
		sampleActiveVUs(3 * time.Second)
		time.AfterFunc(2*time.Second, func() ***REMOVED***
			sampleActiveVUs(0)
		***REMOVED***)
		time.Sleep(1 * time.Second)
		// Sample ramp-down at a higher frequency
		for i := 0; i < 15; i++ ***REMOVED***
			sampleActiveVUs(100 * time.Millisecond)
		***REMOVED***
	***REMOVED***()

	err := executor.Run(ctx, nil)

	wg.Wait()
	require.NoError(t, err)
	assert.Equal(t, int64(0), result[0])
	assert.Equal(t, int64(10), result[1])
	assert.Equal(t, int64(0), result[len(result)-1])

	var curr int64
	last := result[2]
	// Check all ramp-down samples
	for i := 3; i < len(result[2:]); i++ ***REMOVED***
		curr = result[i]
		// Detect ramp-ups, missteps (e.g. 7 -> 4), but ignore pauses
		if curr > last || (curr != last && curr != last-1) ***REMOVED***
			assert.FailNow(t,
				fmt.Sprintf("ramping down wobble bug - "+
					"current: %d, previous: %d\nVU samples: %v", curr, last, result))
		***REMOVED***
		last = curr
	***REMOVED***
***REMOVED***
