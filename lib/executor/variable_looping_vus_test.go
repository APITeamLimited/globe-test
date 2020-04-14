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
	"encoding/json"
	"fmt"
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
	assert.Equal(t, int64(29), iterCount)
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
		3400 * time.Millisecond,
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

func TestVariableLoopingVUsConfigExecutionPlanExample(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	conf := NewVariableLoopingVUsConfig("test")
	conf.StartVUs = null.IntFrom(4)
	conf.Stages = []Stage***REMOVED***
		***REMOVED***Target: null.IntFrom(6), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
	***REMOVED***

	expRawStepsNoZeroEnd := []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 17 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
	***REMOVED***
	rawStepsNoZeroEnd := conf.getRawExecutionSteps(et, false)
	assert.Equal(t, expRawStepsNoZeroEnd, rawStepsNoZeroEnd)
	endOffset, isFinal := lib.GetEndOffset(rawStepsNoZeroEnd)
	assert.Equal(t, 20*time.Second, endOffset)
	assert.Equal(t, false, isFinal)

	rawStepsZeroEnd := conf.getRawExecutionSteps(et, true)
	assert.Equal(t,
		append(expRawStepsNoZeroEnd, lib.ExecutionStep***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 0***REMOVED***),
		rawStepsZeroEnd,
	)
	endOffset, isFinal = lib.GetEndOffset(rawStepsZeroEnd)
	assert.Equal(t, 23*time.Second, endOffset)
	assert.Equal(t, true, isFinal)

	// GracefulStop and GracefulRampDown equal to the default 30 sec
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 41 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 53 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a longer GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(80 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 41 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 103 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a much shorter GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(3 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a zero GracefulStop
	conf.GracefulStop = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a zero GracefulStop and GracefulRampDown, i.e. raw steps with 0 end cap
	conf.GracefulRampDown = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, rawStepsZeroEnd, conf.GetExecutionRequirements(et))
***REMOVED***

func TestVariableLoopingVUsConfigExecutionPlanExampleOneThird(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(newExecutionSegmentFromString("0:1/3"), nil)
	require.NoError(t, err)
	conf := NewVariableLoopingVUsConfig("test")
	conf.StartVUs = null.IntFrom(4)
	conf.Stages = []Stage***REMOVED***
		***REMOVED***Target: null.IntFrom(6), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
	***REMOVED***

	expRawStepsNoZeroEnd := []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
		***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
		***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***
	rawStepsNoZeroEnd := conf.getRawExecutionSteps(et, false)
	assert.Equal(t, expRawStepsNoZeroEnd, rawStepsNoZeroEnd)
	endOffset, isFinal := lib.GetEndOffset(rawStepsNoZeroEnd)
	assert.Equal(t, 20*time.Second, endOffset)
	assert.Equal(t, true, isFinal)

	rawStepsZeroEnd := conf.getRawExecutionSteps(et, true)
	assert.Equal(t, expRawStepsNoZeroEnd, rawStepsZeroEnd)
	endOffset, isFinal = lib.GetEndOffset(rawStepsZeroEnd)
	assert.Equal(t, 20*time.Second, endOffset)
	assert.Equal(t, true, isFinal)

	// GracefulStop and GracefulRampDown equal to the default 30 sec
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 41 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a longer GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(80 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 41 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a much shorter GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(3 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a zero GracefulStop
	conf.GracefulStop = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a zero GracefulStop and GracefulRampDown, i.e. raw steps with 0 end cap
	conf.GracefulRampDown = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, rawStepsZeroEnd, conf.GetExecutionRequirements(et))
***REMOVED***

func TestVariableLoopingVUsExecutionTupleTests(t *testing.T) ***REMOVED***
	t.Parallel()

	conf := NewVariableLoopingVUsConfig("test")
	conf.StartVUs = null.IntFrom(4)
	conf.Stages = []Stage***REMOVED***
		***REMOVED***Target: null.IntFrom(6), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
	***REMOVED***
	/*

			Graph of the above:
			^
		8	|
		7	|
		6	| +
		5	|/ \       +           +--+
		4	+   \     / \     +-+  |  |       *
		3	|    \   /   \   /  |  |  |      /
		2	|     \ /     \ /   |  |  | +   /
		1	|      +       +    +--+  |/ \ /
		0	+-------------------------+---+------------------------------>
		    01234567890123456789012345678901234567890

	*/

	testCases := []struct ***REMOVED***
		expectedSteps []lib.ExecutionStep
		et            *lib.ExecutionTuple
	***REMOVED******REMOVED***
		***REMOVED***
			et: mustNewExecutionTuple(nil, nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 17 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 27 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 31 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 34 * time.Second, PlannedVUs: 4***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 27 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 31 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 34 * time.Second, PlannedVUs: 2***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 17 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		et := testCase.et
		expectedSteps := testCase.expectedSteps

		t.Run(et.String(), func(t *testing.T) ***REMOVED***
			rawStepsNoZeroEnd := conf.getRawExecutionSteps(et, false)
			assert.Equal(t, expectedSteps, rawStepsNoZeroEnd)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkVarriableArrivalRateGetRawExecutionSteps(b *testing.B) ***REMOVED***
	testCases := []struct ***REMOVED***
		seq string
		seg string
	***REMOVED******REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***seg: "0:1"***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.3"***REMOVED***,
		***REMOVED***seq: "0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.1"***REMOVED***,
		***REMOVED***seg: "2/5:4/5"***REMOVED***,
		***REMOVED***seg: "2235/5213:4/5"***REMOVED***, // just wanted it to be ugly ;D
	***REMOVED***

	stageCases := []struct ***REMOVED***
		name   string
		stages string
	***REMOVED******REMOVED***
		***REMOVED***
			name:   "normal",
			stages: `[***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":10000***REMOVED***,***REMOVED***"duration":"5m", "target":10000***REMOVED***]`,
		***REMOVED***, ***REMOVED***
			name: "rollercoaster",
			stages: `[***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***,
				***REMOVED***"duration":"5m", "target":5000***REMOVED***,***REMOVED***"duration":"5m", "target":0***REMOVED***]`,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		tc := tc
		b.Run(fmt.Sprintf("seq:%s;segment:%s", tc.seq, tc.seg), func(b *testing.B) ***REMOVED***
			ess, err := lib.NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(b, err)
			segment, err := lib.NewExecutionSegmentFromString(tc.seg)
			require.NoError(b, err)
			if tc.seg == "" ***REMOVED***
				segment = nil // specifically for the optimization
			***REMOVED***
			et, err := lib.NewExecutionTuple(segment, &ess)
			require.NoError(b, err)
			for _, stageCase := range stageCases ***REMOVED***
				var st []Stage
				require.NoError(b, json.Unmarshal([]byte(stageCase.stages), &st))
				vlvc := VariableLoopingVUsConfig***REMOVED***
					Stages: st,
				***REMOVED***
				b.Run(stageCase.name, func(b *testing.B) ***REMOVED***
					for i := 0; i < b.N; i++ ***REMOVED***
						_ = vlvc.getRawExecutionSteps(et, false)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
