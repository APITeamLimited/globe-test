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
	"math/big"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
)

func TestRampingVUsConfigValidation(t *testing.T) ***REMOVED***
	t.Parallel()

	errs := NewRampingVUsConfig("default").Validate()
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "one stage has to be specified")

	c := NewRampingVUsConfig("stage")
	c.Stages = []Stage***REMOVED***
		***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(12 * time.Second)***REMOVED***,
	***REMOVED***
	errs = c.Validate()
	require.Empty(t, errs) // by default StartVUs is 1

	c.StartVUs = null.IntFrom(0)
	errs = c.Validate()
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "greater than 0")
***REMOVED***

func TestRampingVUsRun(t *testing.T) ***REMOVED***
	t.Parallel()

	config := RampingVUsConfig***REMOVED***
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
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)
	ctx, cancel, executor, _ := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
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
		900 * time.Millisecond,
	***REMOVED***

	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	result := make([]int64, len(sampleTimes))
	for i, d := range sampleTimes ***REMOVED***
		time.Sleep(d)
		result[i] = es.GetCurrentlyActiveVUsCount()
	***REMOVED***

	require.NoError(t, <-errCh)

	assert.Equal(t, []int64***REMOVED***5, 3, 0***REMOVED***, result)
	assert.Equal(t, int64(29), atomic.LoadInt64(&iterCount))
***REMOVED***

func TestRampingVUsGracefulStopWaits(t *testing.T) ***REMOVED***
	t.Parallel()

	config := RampingVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(time.Second)***REMOVED***,
		StartVUs:   null.IntFrom(1),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(1),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var (
		started = make(chan struct***REMOVED******REMOVED***) // the iteration started
		stopped = make(chan struct***REMOVED******REMOVED***) // the iteration stopped
		stop    = make(chan struct***REMOVED******REMOVED***) // the itearation should stop
	)

	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)
	ctx, cancel, executor, _ := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
			close(started)
			defer close(stopped)
			select ***REMOVED***
			case <-ctx.Done():
				t.Fatal("The iterations should've ended before the context")
			case <-stop:
			***REMOVED***
			return nil
		***REMOVED***),
	)
	defer cancel()
	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	<-started
	// 500 milliseconds more then the duration and 500 less then the gracefulStop
	time.Sleep(time.Millisecond * 1500)
	close(stop)
	<-stopped

	require.NoError(t, <-errCh)
***REMOVED***

func TestRampingVUsGracefulStopStops(t *testing.T) ***REMOVED***
	t.Parallel()

	config := RampingVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(time.Second)***REMOVED***,
		StartVUs:   null.IntFrom(1),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(1),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var (
		started = make(chan struct***REMOVED******REMOVED***) // the iteration started
		stopped = make(chan struct***REMOVED******REMOVED***) // the iteration stopped
		stop    = make(chan struct***REMOVED******REMOVED***) // the itearation should stop
	)

	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)
	ctx, cancel, executor, _ := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
			close(started)
			defer close(stopped)
			select ***REMOVED***
			case <-ctx.Done():
			case <-stop:
				t.Fatal("The iterations shouldn't have ended before the context")
			***REMOVED***
			return nil
		***REMOVED***),
	)
	defer cancel()
	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	<-started
	// 500 milliseconds more then the gracefulStop + duration
	time.Sleep(time.Millisecond * 2500)
	close(stop)
	<-stopped

	require.NoError(t, <-errCh)
***REMOVED***

func TestRampingVUsGracefulRampDown(t *testing.T) ***REMOVED***
	t.Parallel()

	config := RampingVUsConfig***REMOVED***
		BaseConfig:       BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(5 * time.Second)***REMOVED***,
		StartVUs:         null.IntFrom(2),
		GracefulRampDown: types.NullDurationFrom(5 * time.Second),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(2),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(0),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var (
		started = make(chan struct***REMOVED******REMOVED***) // the iteration started
		stopped = make(chan struct***REMOVED******REMOVED***) // the iteration stopped
		stop    = make(chan struct***REMOVED******REMOVED***) // the itearation should stop
	)

	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)
	ctx, cancel, executor, _ := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context, state *lib.State) error ***REMOVED***
			if state.VUID == 1 ***REMOVED*** // the first VU will wait here to do stuff
				close(started)
				defer close(stopped)
				select ***REMOVED***
				case <-ctx.Done():
					t.Fatal("The iterations can't have ended before the context")
				case <-stop:
				***REMOVED***
			***REMOVED*** else ***REMOVED*** // all other (1) VUs will just sleep long enough
				time.Sleep(2500 * time.Millisecond)
			***REMOVED***
			return nil
		***REMOVED***),
	)
	defer cancel()
	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	<-started
	// 500 milliseconds more then the gracefulRampDown + duration
	time.Sleep(2500 * time.Millisecond)
	close(stop)
	<-stopped

	select ***REMOVED***
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second): // way too much time
		t.Fatal("Execution should've ended already")
	***REMOVED***
***REMOVED***

// This test aims to check whether the ramping VU executor interrupts
// hanging/remaining VUs after the graceful rampdown period finishes.
//
//                   Rampdown         Graceful Rampdown
//                   Stage (40ms)     (+30ms)
//                   [               ][              ]
//         t 0---5---10---20---30---40---50---60---70
//     VU1       *..................................✔ (40+30=70ms)
//     VU2           *...................X            (20+30=50ms)
//
//     ✔=Finishes,X=Interrupted,.=Sleeps
func TestRampingVUsHandleRemainingVUs(t *testing.T) ***REMOVED***
	t.Parallel()

	const (
		maxVus                   = 2
		vuSleepDuration          = 65 * time.Millisecond // Each VU will sleep 65ms
		wantVuFinished    uint32 = 1                     // one VU should finish an iteration
		wantVuInterrupted uint32 = 1                     // one VU should be interrupted
	)

	cfg := RampingVUsConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***
			// Extend the total test duration 50ms more
			//
			// test duration = sum(stages) + GracefulStop
			//
			// This could have been 30ms but increased it to 50ms
			// to prevent the test to become flaky.
			GracefulStop: types.NullDurationFrom(50 * time.Millisecond),
		***REMOVED***,
		// Wait 30ms more for already started iterations
		// (Happens in the 2nd stage below: Graceful rampdown period)
		GracefulRampDown: types.NullDurationFrom(30 * time.Millisecond),
		// Total test duration is 50ms (excluding the GracefulRampdown period)
		Stages: []Stage***REMOVED***
			// Activate 2 VUs in 10ms
			***REMOVED***
				Duration: types.NullDurationFrom(10 * time.Millisecond),
				Target:   null.IntFrom(int64(maxVus)),
			***REMOVED***,
			// Rampdown to 0 VUs in 40ms
			***REMOVED***
				Duration: types.NullDurationFrom(40 * time.Millisecond),
				Target:   null.IntFrom(int64(0)),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var (
		gotVuInterrupted uint32
		gotVuFinished    uint32
	)
	iteration := func(ctx context.Context, _ *lib.State) error ***REMOVED***
		select ***REMOVED***
		case <-time.After(vuSleepDuration):
			atomic.AddUint32(&gotVuFinished, 1)
		case <-ctx.Done():
			atomic.AddUint32(&gotVuInterrupted, 1)
		***REMOVED***
		return nil
	***REMOVED***

	// run the executor: this should finish in ~70ms
	// sum(stages) + GracefulRampDown
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	ctx, cancel, executor, _ := setupExecutor(
		t, cfg,
		lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, maxVus, maxVus),
		simpleRunner(iteration),
	)
	defer cancel()
	require.NoError(t, executor.Run(ctx, nil))

	assert.Equal(t, wantVuInterrupted, atomic.LoadUint32(&gotVuInterrupted))
	assert.Equal(t, wantVuFinished, atomic.LoadUint32(&gotVuFinished))
***REMOVED***

// Ensure there's no wobble of VUs during graceful ramp-down, without segments.
// See https://github.com/k6io/k6/issues/1296
func TestRampingVUsRampDownNoWobble(t *testing.T) ***REMOVED***
	t.Parallel()

	config := RampingVUsConfig***REMOVED***
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
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, nil, 10, 50)
	ctx, cancel, executor, _ := setupExecutor(
		t, config, es,
		simpleRunner(func(ctx context.Context, _ *lib.State) error ***REMOVED***
			time.Sleep(500 * time.Millisecond)
			return nil
		***REMOVED***),
	)
	defer cancel()

	sampleTimes := []time.Duration***REMOVED***
		100 * time.Millisecond,
		3000 * time.Millisecond,
	***REMOVED***
	const rampDownSampleTime = 50 * time.Millisecond
	rampDownSamples := int((config.Stages[len(config.Stages)-1].Duration.TimeDuration() + config.GracefulRampDown.TimeDuration()) / rampDownSampleTime)

	errCh := make(chan error)
	go func() ***REMOVED*** errCh <- executor.Run(ctx, nil) ***REMOVED***()

	result := make([]int64, len(sampleTimes)+rampDownSamples)
	for i, d := range sampleTimes ***REMOVED***
		time.Sleep(d)
		result[i] = es.GetCurrentlyActiveVUsCount()
	***REMOVED***

	// Sample ramp-down at a higher rate
	for i := len(sampleTimes); i < rampDownSamples; i++ ***REMOVED***
		time.Sleep(rampDownSampleTime)
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

func TestRampingVUsConfigExecutionPlanExample(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	conf := NewRampingVUsConfig("test")
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
		***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 1***REMOVED***,
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
		***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 53 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a longer GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(80 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 4***REMOVED***,
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

func TestRampingVUsConfigExecutionPlanExampleOneThird(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(newExecutionSegmentFromString("0:1/3"), nil)
	require.NoError(t, err)
	conf := NewRampingVUsConfig("test")
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
		***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
		***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
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
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(et))

	// Try a longer GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(80 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 1***REMOVED***,
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

func TestRampingVUsExecutionTupleTests(t *testing.T) ***REMOVED***
	t.Parallel()

	conf := NewRampingVUsConfig("test")
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
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 17 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 27 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 30 * time.Second, PlannedVUs: 0***REMOVED***,
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
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), nil),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 27 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 30 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 31 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 34 * time.Second, PlannedVUs: 2***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 28 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 29 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 32 * time.Second, PlannedVUs: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			et: mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 0***REMOVED***,
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
			t.Parallel()
			rawStepsNoZeroEnd := conf.getRawExecutionSteps(et, false)
			assert.Equal(t, expectedSteps, rawStepsNoZeroEnd)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestRampingVUsGetRawExecutionStepsCornerCases(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name          string
		expectedSteps []lib.ExecutionStep
		et            *lib.ExecutionTuple
		stages        []Stage
		start         int64
	***REMOVED******REMOVED***
		***REMOVED***
			name: "going up then down straight away",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 3***REMOVED***,
			***REMOVED***,
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(3), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
			start: 2,
		***REMOVED***,
		***REMOVED***
			name: "jump up then go up again",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 5***REMOVED***,
			***REMOVED***,
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
			start: 3,
		***REMOVED***,
		***REMOVED***
			name: "up down up down",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "up down up down in half",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/2"), nil),
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "up down up down in the other half",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			et: mustNewExecutionTuple(newExecutionSegmentFromString("1/2:1"), nil),
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "up down up down in with nothing",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			et: mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "up down up down in with funky sequence", // panics if there are no localIndex == 0 guards
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,1/2,2/3,1")),
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(2), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "strange",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 6***REMOVED***,
				***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 7***REMOVED***,
				***REMOVED***TimeOffset: 35 * time.Second, PlannedVUs: 8***REMOVED***,
				***REMOVED***TimeOffset: 44 * time.Second, PlannedVUs: 9***REMOVED***,
			***REMOVED***,
			et: mustNewExecutionTuple(newExecutionSegmentFromString("0:0.3"), newExecutionSegmentSequenceFromString("0,0.3,0.6,0.9,1")),
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(20), Duration: types.NullDurationFrom(20 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(30), Duration: types.NullDurationFrom(30 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "more up and down",
			expectedSteps: []lib.ExecutionStep***REMOVED***
				***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 0***REMOVED***,
				***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 5***REMOVED***,
				***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 4***REMOVED***,
				***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 3***REMOVED***,
				***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
				***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 1***REMOVED***,
				***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 0***REMOVED***,
			***REMOVED***,
			stages: []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(0), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		conf := NewRampingVUsConfig("test")
		conf.StartVUs = null.IntFrom(testCase.start)
		conf.Stages = testCase.stages
		et := testCase.et
		if et == nil ***REMOVED***
			et = mustNewExecutionTuple(nil, nil)
		***REMOVED***
		expectedSteps := testCase.expectedSteps

		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			rawStepsNoZeroEnd := conf.getRawExecutionSteps(et, false)
			assert.Equal(t, expectedSteps, rawStepsNoZeroEnd)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkRampingVUsGetRawExecutionSteps(b *testing.B) ***REMOVED***
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
				vlvc := RampingVUsConfig***REMOVED***
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

// TODO: delete in favor of lib.generateRandomSequence() after
// https://github.com/k6io/k6/issues/1302 is done (can't import now due to
// import loops...)
func generateRandomSequence(t testing.TB, n, m int64, r *rand.Rand) lib.ExecutionSegmentSequence ***REMOVED***
	var err error
	ess := lib.ExecutionSegmentSequence(make([]*lib.ExecutionSegment, n))
	numerators := make([]int64, n)
	var denominator int64
	for i := int64(0); i < n; i++ ***REMOVED***
		numerators[i] = 1 + r.Int63n(m)
		denominator += numerators[i]
	***REMOVED***
	from := big.NewRat(0, 1)
	for i := int64(0); i < n; i++ ***REMOVED***
		to := new(big.Rat).Add(big.NewRat(numerators[i], denominator), from)
		ess[i], err = lib.NewExecutionSegment(from, to)
		require.NoError(t, err)
		from = to
	***REMOVED***

	return ess
***REMOVED***

func TestSumRandomSegmentSequenceMatchesNoSegment(t *testing.T) ***REMOVED***
	t.Parallel()

	const (
		numTests         = 10
		maxStages        = 10
		minStageDuration = 1 * time.Second
		maxStageDuration = 10 * time.Minute
		maxVUs           = 300
		segmentSeqMaxLen = 15
		maxNumerator     = 300
	)
	getTestConfig := func(r *rand.Rand, name string) RampingVUsConfig ***REMOVED***
		stagesCount := 1 + r.Int31n(maxStages)
		stages := make([]Stage, stagesCount)
		for s := int32(0); s < stagesCount; s++ ***REMOVED***
			dur := (minStageDuration + time.Duration(r.Int63n(int64(maxStageDuration-minStageDuration)))).Round(time.Second)
			stages[s] = Stage***REMOVED***Duration: types.NullDurationFrom(dur), Target: null.IntFrom(r.Int63n(maxVUs))***REMOVED***
		***REMOVED***

		c := NewRampingVUsConfig(name)
		c.GracefulRampDown = types.NullDurationFrom(0)
		c.GracefulStop = types.NullDurationFrom(0)
		c.StartVUs = null.IntFrom(r.Int63n(maxVUs))
		c.Stages = stages
		return c
	***REMOVED***

	subtractChildSteps := func(t *testing.T, parent, child []lib.ExecutionStep) ***REMOVED***
		t.Logf("subtractChildSteps()")
		for _, step := range child ***REMOVED***
			t.Logf("	child planned VUs for time offset %s: %d", step.TimeOffset, step.PlannedVUs)
		***REMOVED***
		sub := uint64(0)
		ci := 0
		for pi, p := range parent ***REMOVED***
			// We iterate over all parent steps and match them to child steps.
			// Once we have a match, we remove the child step's plannedVUs from
			// the parent steps until a new match, when we adjust the subtracted
			// amount again.
			if p.TimeOffset > child[ci].TimeOffset && ci != len(child)-1 ***REMOVED***
				t.Errorf("ERR Could not match child offset %s with any parent time offset", child[ci].TimeOffset)
			***REMOVED***
			if p.TimeOffset == child[ci].TimeOffset ***REMOVED***
				t.Logf("Setting sub to %d at t=%s", child[ci].PlannedVUs, child[ci].TimeOffset)
				sub = child[ci].PlannedVUs
				if ci != len(child)-1 ***REMOVED***
					ci++
				***REMOVED***
			***REMOVED***
			t.Logf("Subtracting %d VUs (out of %d) at t=%s", sub, p.PlannedVUs, p.TimeOffset)
			parent[pi].PlannedVUs -= sub
		***REMOVED***
	***REMOVED***

	for i := 0; i < numTests; i++ ***REMOVED***
		name := fmt.Sprintf("random%02d", i)
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			seed := time.Now().UnixNano()
			r := rand.New(rand.NewSource(seed)) //nolint:gosec
			t.Logf("Random source seeded with %d\n", seed)
			c := getTestConfig(r, name)
			ranSeqLen := 2 + r.Int63n(segmentSeqMaxLen-1)
			t.Logf("Config: %#v, ranSeqLen: %d", c, ranSeqLen)
			randomSequence := generateRandomSequence(t, ranSeqLen, maxNumerator, r)
			t.Logf("Random sequence: %s", randomSequence)
			fullSeg, err := lib.NewExecutionTuple(nil, nil)
			require.NoError(t, err)
			fullRawSteps := c.getRawExecutionSteps(fullSeg, false)

			for _, step := range fullRawSteps ***REMOVED***
				t.Logf("original planned VUs for time offset %s: %d", step.TimeOffset, step.PlannedVUs)
			***REMOVED***

			for s := 0; s < len(randomSequence); s++ ***REMOVED***
				et, err := lib.NewExecutionTuple(randomSequence[s], &randomSequence)
				require.NoError(t, err)
				segRawSteps := c.getRawExecutionSteps(et, false)
				subtractChildSteps(t, fullRawSteps, segRawSteps)
			***REMOVED***

			for _, step := range fullRawSteps ***REMOVED***
				if step.PlannedVUs != 0 ***REMOVED***
					t.Errorf("ERR Remaining planned VUs for time offset %s are not 0 but %d", step.TimeOffset, step.PlannedVUs)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
