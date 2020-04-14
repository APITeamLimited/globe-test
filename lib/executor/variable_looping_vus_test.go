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
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			// Sleeping for a weird duration somewhat offset from the
			// executor ticks to hopefully keep race conditions out of
			// our control from failing the test.
			time.Sleep(300 * time.Millisecond)
			atomic.AddInt64(&iterCount, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()

	sampleTimes := []time.Duration***REMOVED***
		500 * time.Millisecond,
		1000 * time.Millisecond,
		700 * time.Millisecond,
	***REMOVED***

	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	var result = make([]int64, len(sampleTimes))
	for i, d := range sampleTimes ***REMOVED***
		time.Sleep(d)
		result[i] = es.GetCurrentlyActiveVUsCount()
	***REMOVED***

	require.NoError(t, <-errCh)

	assert.Equal(t, []int64***REMOVED***5, 3, 0***REMOVED***, result)
	assert.Equal(t, int64(29), atomic.LoadInt64(&iterCount))
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

	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, _ = setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(1 * time.Second)
			return nil
		***REMOVED***),
	)
	defer cancel()

	sampleTimes := []time.Duration***REMOVED***
		100 * time.Millisecond,
		3000 * time.Millisecond,
	***REMOVED***
	const rampDownSamples = 50

	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	var result = make([]int64, len(sampleTimes)+rampDownSamples)
	for i, d := range sampleTimes ***REMOVED***
		time.Sleep(d)
		result[i] = es.GetCurrentlyActiveVUsCount()
	***REMOVED***

	// Sample ramp-down at a higher rate
	for i := len(sampleTimes); i < rampDownSamples; i++ ***REMOVED***
		time.Sleep(50 * time.Millisecond)
		result[i] = es.GetCurrentlyActiveVUsCount()
	***REMOVED***

	require.NoError(t, <-errCh)

	// Some baseline checks
	assert.Equal(t, int64(0), result[0])
	assert.Equal(t, int64(10), result[1])
	assert.Equal(t, int64(0), result[len(result)-1])

	vuChanges := []int64***REMOVED***result[2]***REMOVED***
	// Check ramp-down consistency
	for i := 3; i < len(result[2:]); i++ ***REMOVED***
		if result[i] != result[i-1] ***REMOVED***
			vuChanges = append(vuChanges, result[i])
		***REMOVED***
	***REMOVED***
	assert.Equal(t, []int64***REMOVED***10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0***REMOVED***, vuChanges)
***REMOVED***
