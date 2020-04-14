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
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

func getTestVariableArrivalRateConfig() VariableArrivalRateConfig ***REMOVED***
	return VariableArrivalRateConfig***REMOVED***
		BaseConfig: BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(1 * time.Second)***REMOVED***,
		TimeUnit:   types.NullDurationFrom(time.Second),
		StartRate:  null.IntFrom(10),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(time.Second * 1),
				Target:   null.IntFrom(10),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(time.Second * 1),
				Target:   null.IntFrom(50),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(time.Second * 1),
				Target:   null.IntFrom(50),
			***REMOVED***,
		***REMOVED***,
		PreAllocatedVUs: null.IntFrom(10),
		MaxVUs:          null.IntFrom(20),
	***REMOVED***
***REMOVED***

func TestVariableArrivalRateRunNotEnoughAllocatedVUsWarn(t *testing.T) ***REMOVED***
	t.Parallel()
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, logHook = setupExecutor(
		t, getTestVariableArrivalRateConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			time.Sleep(time.Second)
			return nil
		***REMOVED***),
	)
	defer cancel()
	var engineOut = make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	require.NoError(t, err)
	entries := logHook.Drain()
	require.NotEmpty(t, entries)
	for _, entry := range entries ***REMOVED***
		require.Equal(t,
			"Insufficient VUs, reached 20 active VUs and cannot allocate more",
			entry.Message)
		require.Equal(t, logrus.WarnLevel, entry.Level)
	***REMOVED***
***REMOVED***

func TestVariableArrivalRateRunCorrectRate(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var ctx, cancel, executor, logHook = setupExecutor(
		t, getTestVariableArrivalRateConfig(), es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			atomic.AddInt64(&count, 1)
			return nil
		***REMOVED***),
	)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		// check that we got around the amount of VU iterations as we would expect
		var currentCount int64

		time.Sleep(time.Second)
		currentCount = atomic.SwapInt64(&count, 0)
		assert.InDelta(t, 10, currentCount, 1)

		time.Sleep(time.Second)
		currentCount = atomic.SwapInt64(&count, 0)
		assert.InDelta(t, 30, currentCount, 2)

		time.Sleep(time.Second)
		currentCount = atomic.SwapInt64(&count, 0)
		assert.InDelta(t, 50, currentCount, 2)
	***REMOVED***()
	var engineOut = make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	wg.Wait()
	require.NoError(t, err)
	require.Empty(t, logHook.Drain())
***REMOVED***

func TestVariableArrivalRateRunCorrectRateWithSlowRate(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	var now = time.Now()
	et, err := lib.NewExecutionTuple(nil, nil)
	require.NoError(t, err)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, et, 10, 50)
	var expectedTimes = []time.Duration***REMOVED***
		time.Millisecond * 3464, time.Millisecond * 4898, time.Second * 6***REMOVED***
	var ctx, cancel, executor, logHook = setupExecutor(
		t, VariableArrivalRateConfig***REMOVED***
			TimeUnit: types.NullDurationFrom(time.Second),
			Stages: []Stage***REMOVED***
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 6),
					Target:   null.IntFrom(1),
				***REMOVED***,
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 0),
					Target:   null.IntFrom(0),
				***REMOVED***,
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 1),
					Target:   null.IntFrom(0),
				***REMOVED***,
			***REMOVED***,
			PreAllocatedVUs: null.IntFrom(10),
			MaxVUs:          null.IntFrom(20),
		***REMOVED***,
		es,
		simpleRunner(func(ctx context.Context) error ***REMOVED***
			current := atomic.AddInt64(&count, 1)
			if !assert.True(t, int(current) <= len(expectedTimes)) ***REMOVED***
				return nil
			***REMOVED***
			expectedTime := expectedTimes[current-1]
			assert.WithinDuration(t,
				now.Add(expectedTime),
				time.Now(),
				time.Millisecond*100,
				"%d expectedTime %s", current, expectedTime,
			)
			return nil
		***REMOVED***),
	)
	defer cancel()
	var engineOut = make(chan stats.SampleContainer, 1000)
	err = executor.Run(ctx, engineOut)
	require.NoError(t, err)
	require.Equal(t, int64(len(expectedTimes)), count)
	require.Empty(t, logHook.Drain())
***REMOVED***

func mustNewExecutionTuple(seg *lib.ExecutionSegment, seq *lib.ExecutionSegmentSequence) *lib.ExecutionTuple ***REMOVED***
	et, err := lib.NewExecutionTuple(seg, seq)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return et
***REMOVED***

func TestVariableArrivalRateCal(t *testing.T) ***REMOVED***
	t.Parallel()

	var (
		defaultTimeUnit = time.Second
		config          = VariableArrivalRateConfig***REMOVED***
			StartRate: null.IntFrom(0),
			Stages: []Stage***REMOVED*** // TODO make this even bigger and longer .. will need more time
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 5),
					Target:   null.IntFrom(1),
				***REMOVED***,
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 1),
					Target:   null.IntFrom(1),
				***REMOVED***,
				***REMOVED***
					Duration: types.NullDurationFrom(time.Second * 5),
					Target:   null.IntFrom(0),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
	)

	testCases := []struct ***REMOVED***
		expectedTimes []time.Duration
		et            *lib.ExecutionTuple
		timeUnit      time.Duration
	***REMOVED******REMOVED***
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 3162, time.Millisecond * 4472, time.Millisecond * 5500, time.Millisecond * 6527, time.Millisecond * 7837, time.Second * 11***REMOVED***,
			et:            mustNewExecutionTuple(nil, nil),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 4472, time.Millisecond * 7837***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), nil),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 4472, time.Millisecond * 7837***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,1")),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 4472, time.Millisecond * 7837***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), nil),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 4472, time.Millisecond * 7837***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), nil),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 3162, time.Millisecond * 6527***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("0:1/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 4472, time.Millisecond * 7837***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("1/3:2/3"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 5500, time.Millisecond * 11000***REMOVED***,
			et:            mustNewExecutionTuple(newExecutionSegmentFromString("2/3:1"), newExecutionSegmentSequenceFromString("0,1/3,2/3,1")),
		***REMOVED***,
		***REMOVED***
			expectedTimes: []time.Duration***REMOVED***time.Millisecond * 1825, time.Millisecond * 2581, time.Millisecond * 3162, time.Millisecond * 3651, time.Millisecond * 4082, time.Millisecond * 4472,
				time.Millisecond * 4830, time.Millisecond * 5166, time.Millisecond * 5499, time.Millisecond * 5833, time.Millisecond * 6169, time.Millisecond * 6527,
				time.Millisecond * 6917, time.Millisecond * 7348, time.Millisecond * 7837, time.Millisecond * 8418, time.Millisecond * 9174, time.Millisecond * 10999***REMOVED***,
			et:       mustNewExecutionTuple(nil, nil),
			timeUnit: time.Second / 3, // three  times as fast
		***REMOVED***,
		// TODO: extend more
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		et := testCase.et
		expectedTimes := testCase.expectedTimes
		config.TimeUnit = types.NewNullDuration(testCase.timeUnit, true)
		if testCase.timeUnit == 0 ***REMOVED***
			config.TimeUnit = types.NewNullDuration(defaultTimeUnit, true)
		***REMOVED***

		t.Run(fmt.Sprintf("%s timeunit %s", et, config.TimeUnit), func(t *testing.T) ***REMOVED***
			var ch = make(chan time.Duration)
			go config.cal(et, ch)
			var changes = make([]time.Duration, 0, len(expectedTimes))
			for c := range ch ***REMOVED***
				changes = append(changes, c)
			***REMOVED***
			assert.Equal(t, len(expectedTimes), len(changes))
			for i, expectedTime := range expectedTimes ***REMOVED***
				require.True(t, i < len(changes))
				change := changes[i]
				assert.InEpsilon(t, expectedTime, change, 0.001, "%s %s", expectedTime, change)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkCal(b *testing.B) ***REMOVED***
	for _, t := range []time.Duration***REMOVED***
		time.Second, time.Minute,
	***REMOVED*** ***REMOVED***
		t := t
		b.Run(t.String(), func(b *testing.B) ***REMOVED***
			var config = VariableArrivalRateConfig***REMOVED***
				TimeUnit:  types.NullDurationFrom(time.Second),
				StartRate: null.IntFrom(50),
				Stages: []Stage***REMOVED***
					***REMOVED***
						Duration: types.NullDurationFrom(t),
						Target:   null.IntFrom(49),
					***REMOVED***,
					***REMOVED***
						Duration: types.NullDurationFrom(t),
						Target:   null.IntFrom(50),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***
			et := mustNewExecutionTuple(nil, nil)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) ***REMOVED***
				for pb.Next() ***REMOVED***
					var ch = make(chan time.Duration, 20)
					go config.cal(et, ch)
					for c := range ch ***REMOVED***
						_ = c
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkCalRat(b *testing.B) ***REMOVED***
	for _, t := range []time.Duration***REMOVED***
		time.Second, time.Minute,
	***REMOVED*** ***REMOVED***
		t := t
		b.Run(t.String(), func(b *testing.B) ***REMOVED***
			var config = VariableArrivalRateConfig***REMOVED***
				TimeUnit:  types.NullDurationFrom(time.Second),
				StartRate: null.IntFrom(50),
				Stages: []Stage***REMOVED***
					***REMOVED***
						Duration: types.NullDurationFrom(t),
						Target:   null.IntFrom(49),
					***REMOVED***,
					***REMOVED***
						Duration: types.NullDurationFrom(t),
						Target:   null.IntFrom(50),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***
			et := mustNewExecutionTuple(nil, nil)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) ***REMOVED***
				for pb.Next() ***REMOVED***
					var ch = make(chan time.Duration, 20)
					go config.calRat(et, ch)
					for c := range ch ***REMOVED***
						_ = c
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestCompareCalImplementation(t *testing.T) ***REMOVED***
	t.Parallel()
	// This test checks that the cal and calRat implementation get roughly similar numbers
	// in my experiment the difference is 1(nanosecond) in 7 case for the whole test
	// the duration is 1 second for each stage as calRat takes way longer - a longer better test can
	// be done when/if it's performance is improved
	var config = VariableArrivalRateConfig***REMOVED***
		TimeUnit:  types.NullDurationFrom(time.Second),
		StartRate: null.IntFrom(0),
		Stages: []Stage***REMOVED***
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(200),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(200),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(2000),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(2000),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(300),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(300),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(1333),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(1334),
			***REMOVED***,
			***REMOVED***
				Duration: types.NullDurationFrom(1 * time.Second),
				Target:   null.IntFrom(1334),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	et := mustNewExecutionTuple(nil, nil)
	var chRat = make(chan time.Duration, 20)
	var ch = make(chan time.Duration, 20)
	go config.calRat(et, chRat)
	go config.cal(et, ch)
	count := 0
	var diff int
	for c := range ch ***REMOVED***
		count++
		cRat := <-chRat
		if !assert.InDelta(t, c, cRat, 1, "%d", count) ***REMOVED***
			diff++
		***REMOVED***
	***REMOVED***
	require.Equal(t, 0, diff)
***REMOVED***

// calRat code here is just to check how accurate the cal implemenattion is
// there are no other tests for it so it depends on the test of cal that it is actually accurate :D

//nolint:gochecknoglobals
var two = big.NewRat(2, 1)

// from https://groups.google.com/forum/#!topic/golang-nuts/aIcDf8T-Png
func sqrtRat(x *big.Rat) *big.Rat ***REMOVED***
	var z, a, b big.Rat
	var ns, ds big.Int
	ni, di := x.Num(), x.Denom()
	z.SetFrac(ns.Rsh(ni, uint(ni.BitLen())/2), ds.Rsh(di, uint(di.BitLen())/2))
	for i := 10; i > 0; i-- ***REMOVED*** //TODO: better termination
		a.Sub(a.Mul(&z, &z), x)
		f, _ := a.Float64()
		if f == 0 ***REMOVED***
			break
		***REMOVED***
		// fmt.Println(x, z, i)
		z.Sub(&z, b.Quo(&a, b.Mul(two, &z)))
	***REMOVED***
	return &z
***REMOVED***

// This implementation is just for reference and accuracy testing
func (varc VariableArrivalRateConfig) calRat(et *lib.ExecutionTuple, ch chan<- time.Duration) ***REMOVED***
	defer close(ch)

	start, offsets, _ := et.GetStripedOffsets(et.ES)
	li := -1
	next := func() int64 ***REMOVED***
		li++
		return offsets[li%len(offsets)]
	***REMOVED***
	iRat := big.NewRat(start+1, 1)

	var carry = big.NewRat(0, 1)
	var doneSoFar = big.NewRat(0, 1)
	var endCount = big.NewRat(0, 1)
	curr := varc.StartRate.ValueOrZero()
	var base time.Duration
	for _, stage := range varc.Stages ***REMOVED***
		target := stage.Target.ValueOrZero()
		if target != curr ***REMOVED***
			var (
				from = big.NewRat(curr, int64(time.Second))
				to   = big.NewRat(target, int64(time.Second))
				dur  = big.NewRat(time.Duration(stage.Duration.Duration).Nanoseconds(), 1)
			)
			// precalcualations :)
			toMinusFrom := new(big.Rat).Sub(to, from)
			fromSquare := new(big.Rat).Mul(from, from)
			durMulSquare := new(big.Rat).Mul(dur, fromSquare)
			fromMulDur := new(big.Rat).Mul(from, dur)
			oneOverToMinusFrom := new(big.Rat).Inv(toMinusFrom)

			endCount.Add(endCount,
				new(big.Rat).Mul(
					dur,
					new(big.Rat).Add(new(big.Rat).Mul(toMinusFrom, big.NewRat(1, 2)), from)))
			for ; endCount.Cmp(iRat) >= 0; iRat.Add(iRat, big.NewRat(next(), 1)) ***REMOVED***
				// even with all of this optimizations sqrtRat is taking so long this is still
				// extremely slow ... :(
				buf := new(big.Rat).Sub(iRat, doneSoFar)
				buf.Mul(buf, two)
				buf.Mul(buf, toMinusFrom)
				buf.Add(buf, durMulSquare)
				buf.Mul(buf, dur)
				buf.Sub(fromMulDur, sqrtRat(buf))
				buf.Mul(buf, oneOverToMinusFrom)

				r, _ := buf.Float64()
				ch <- base + time.Duration(-r) // the minus is because we don't deive by from-to but by to-from above
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			step := big.NewRat(int64(time.Second), target)
			first := big.NewRat(0, 1)
			first.Sub(first, carry)
			endCount.Add(endCount, new(big.Rat).Mul(big.NewRat(target, 1), big.NewRat(time.Duration(stage.Duration.Duration).Nanoseconds(), time.Duration(varc.TimeUnit.Duration).Nanoseconds())))

			for ; endCount.Cmp(iRat) >= 0; iRat.Add(iRat, big.NewRat(next(), 1)) ***REMOVED***
				res := new(big.Rat).Sub(iRat, doneSoFar) // this can get next added to it but will need to change the above for .. so
				r, _ := res.Mul(res, step).Float64()
				ch <- base + time.Duration(r)
				first.Add(first, step)
			***REMOVED***
		***REMOVED***
		doneSoFar.Set(endCount) // copy
		curr = target
		base += time.Duration(stage.Duration.Duration)
	***REMOVED***
***REMOVED***
