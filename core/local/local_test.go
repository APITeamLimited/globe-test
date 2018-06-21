/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package local

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/netext"

	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestExecutorRun(t *testing.T) ***REMOVED***
	e := New(nil)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))

	ctx, cancel := context.WithCancel(context.Background())
	err := make(chan error, 1)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	cancel()
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetupTeardownRun(t *testing.T) ***REMOVED***
	t.Run("Normal", func(t *testing.T) ***REMOVED***
		setupC := make(chan struct***REMOVED******REMOVED***)
		teardownC := make(chan struct***REMOVED******REMOVED***)
		e := New(&lib.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
				close(setupC)
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				close(teardownC)
				return nil
			***REMOVED***,
		***REMOVED***)

		ctx, cancel := context.WithCancel(context.Background())
		err := make(chan error, 1)
		go func() ***REMOVED*** err <- e.Run(ctx, make(chan stats.SampleContainer, 100)) ***REMOVED***()
		cancel()
		<-setupC
		<-teardownC
		assert.NoError(t, <-err)
	***REMOVED***)
	t.Run("Setup Error", func(t *testing.T) ***REMOVED***
		e := New(&lib.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
				return nil, errors.New("setup error")
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***)
		assert.EqualError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 100)), "setup error")

		t.Run("Don't Run Setup", func(t *testing.T) ***REMOVED***
			e := New(&lib.MiniRunner***REMOVED***
				SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
					return nil, errors.New("setup error")
				***REMOVED***,
				TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					return errors.New("teardown error")
				***REMOVED***,
			***REMOVED***)
			e.SetRunSetup(false)
			e.SetEndIterations(null.IntFrom(1))
			assert.NoError(t, e.SetVUsMax(1))
			assert.NoError(t, e.SetVUs(1))
			assert.EqualError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 100)), "teardown error")
		***REMOVED***)
	***REMOVED***)
	t.Run("Teardown Error", func(t *testing.T) ***REMOVED***
		e := New(&lib.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***)
		e.SetEndIterations(null.IntFrom(1))
		assert.NoError(t, e.SetVUsMax(1))
		assert.NoError(t, e.SetVUs(1))
		assert.EqualError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 100)), "teardown error")

		t.Run("Don't Run Teardown", func(t *testing.T) ***REMOVED***
			e := New(&lib.MiniRunner***REMOVED***
				SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
					return nil, nil
				***REMOVED***,
				TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					return errors.New("teardown error")
				***REMOVED***,
			***REMOVED***)
			e.SetRunTeardown(false)
			e.SetEndIterations(null.IntFrom(1))
			assert.NoError(t, e.SetVUsMax(1))
			assert.NoError(t, e.SetVUs(1))
			assert.NoError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 100)))
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestExecutorSetLogger(t *testing.T) ***REMOVED***
	logger, _ := logtest.NewNullLogger()
	e := New(nil)
	e.SetLogger(logger)
	assert.Equal(t, logger, e.GetLogger())
***REMOVED***

func TestExecutorStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Duration time.Duration
		Stages   []lib.Stage
	***REMOVED******REMOVED***
		"one": ***REMOVED***
			1 * time.Second,
			[]lib.Stage***REMOVED******REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED******REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/targeted": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			e := New(&lib.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					time.Sleep(100 * time.Millisecond)
					return nil
				***REMOVED***,
				Options: lib.Options***REMOVED***
					MetricSamplesBufferSize: null.IntFrom(500),
				***REMOVED***,
			***REMOVED***)
			assert.NoError(t, e.SetVUsMax(10))
			e.SetStages(data.Stages)
			assert.NoError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 500)))
			assert.True(t, e.GetTime() >= data.Duration)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutorEndTime(t *testing.T) ***REMOVED***
	e := New(&lib.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(100 * time.Millisecond)
			return nil
		***REMOVED***,
		Options: lib.Options***REMOVED***MetricSamplesBufferSize: null.IntFrom(200)***REMOVED***,
	***REMOVED***)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))
	e.SetEndTime(types.NullDurationFrom(1 * time.Second))
	assert.Equal(t, types.NullDurationFrom(1*time.Second), e.GetEndTime())

	startTime := time.Now()
	assert.NoError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 200)))
	assert.True(t, time.Now().After(startTime.Add(1*time.Second)), "test did not take 1s")

	t.Run("Runtime Errors", func(t *testing.T) ***REMOVED***
		e := New(&lib.MiniRunner***REMOVED***
			Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				time.Sleep(10 * time.Millisecond)
				return errors.New("hi")
			***REMOVED***,
			Options: lib.Options***REMOVED***MetricSamplesBufferSize: null.IntFrom(200)***REMOVED***,
		***REMOVED***)
		assert.NoError(t, e.SetVUsMax(10))
		assert.NoError(t, e.SetVUs(10))
		e.SetEndTime(types.NullDurationFrom(100 * time.Millisecond))
		assert.Equal(t, types.NullDurationFrom(100*time.Millisecond), e.GetEndTime())

		l, hook := logtest.NewNullLogger()
		e.SetLogger(l)

		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 200)))
		assert.True(t, time.Now().After(startTime.Add(100*time.Millisecond)), "test did not take 100ms")

		assert.NotEmpty(t, hook.Entries)
		for _, e := range hook.Entries ***REMOVED***
			assert.Equal(t, "hi", e.Message)
		***REMOVED***
	***REMOVED***)

	t.Run("End Errors", func(t *testing.T) ***REMOVED***
		e := New(&lib.MiniRunner***REMOVED***
			Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				<-ctx.Done()
				return errors.New("hi")
			***REMOVED***,
			Options: lib.Options***REMOVED***MetricSamplesBufferSize: null.IntFrom(200)***REMOVED***,
		***REMOVED***)
		assert.NoError(t, e.SetVUsMax(10))
		assert.NoError(t, e.SetVUs(10))
		e.SetEndTime(types.NullDurationFrom(100 * time.Millisecond))
		assert.Equal(t, types.NullDurationFrom(100*time.Millisecond), e.GetEndTime())

		l, hook := logtest.NewNullLogger()
		e.SetLogger(l)

		startTime := time.Now()
		assert.NoError(t, e.Run(context.Background(), make(chan stats.SampleContainer, 200)))
		assert.True(t, time.Now().After(startTime.Add(100*time.Millisecond)), "test did not take 100ms")

		assert.Empty(t, hook.Entries)
	***REMOVED***)
***REMOVED***

func TestExecutorEndIterations(t *testing.T) ***REMOVED***
	metric := &stats.Metric***REMOVED***Name: "test_metric"***REMOVED***

	var i int64
	e := New(&lib.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		default:
			atomic.AddInt64(&i, 1)
		***REMOVED***
		out <- stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***
		return nil
	***REMOVED******REMOVED***)
	assert.NoError(t, e.SetVUsMax(1))
	assert.NoError(t, e.SetVUs(1))
	e.SetEndIterations(null.IntFrom(100))
	assert.Equal(t, null.IntFrom(100), e.GetEndIterations())

	samples := make(chan stats.SampleContainer, 201)
	assert.NoError(t, e.Run(context.Background(), samples))
	assert.Equal(t, int64(100), e.GetIterations())
	assert.Equal(t, int64(100), i)
	for i := 0; i < 100; i++ ***REMOVED***
		mySample, ok := <-samples
		require.True(t, ok)
		assert.Equal(t, stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***, mySample)
		sample, ok := <-samples
		require.True(t, ok)
		iterSample, ok := (sample).(stats.Sample)
		require.True(t, ok)
		assert.Equal(t, metrics.Iterations, iterSample.Metric)
		assert.Equal(t, float64(1), iterSample.Value)
	***REMOVED***
***REMOVED***

func TestExecutorIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e := New(nil)

	err := make(chan error)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	for !e.IsRunning() ***REMOVED***
	***REMOVED***
	cancel()
	for e.IsRunning() ***REMOVED***
	***REMOVED***
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUsMax(-1), "vu cap can't be negative")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(nil)

		assert.NoError(t, e.SetVUsMax(50))
		assert.Equal(t, int64(50), e.GetVUsMax())

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(50))
			assert.Equal(t, int64(50), e.GetVUsMax())
		***REMOVED***)
	***REMOVED***)

	t.Run("TooLow", func(t *testing.T) ***REMOVED***
		e := New(nil)
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())

		assert.EqualError(t, e.SetVUsMax(50), "can't lower vu cap (to 50) below vu count (100)")
	***REMOVED***)
***REMOVED***

func TestExecutorSetVUs(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(-1), "vu count can't be negative")
	***REMOVED***)

	t.Run("Too High", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(100), "can't raise vu count (to 100) above vu cap (0)")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(&lib.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			return nil
		***REMOVED******REMOVED***)
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				num++
				if assert.NotNil(t, handle.vu, "vu %d lacks impl", i) ***REMOVED***
					assert.Equal(t, int64(0), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
				assert.Nil(t, handle.ctx, "vu %d has ctx", i)
				assert.Nil(t, handle.cancel, "vu %d has cancel", i)
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(50))
		assert.Equal(t, int64(50), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				if i < 50 ***REMOVED***
					assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
					assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
					num++
				***REMOVED*** else ***REMOVED***
					assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					assert.Equal(t, int64(0), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
			***REMOVED***
			assert.Equal(t, 50, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
				assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
				num++
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(50))
			assert.Equal(t, int64(50), e.GetVUs())
			if assert.Len(t, e.vus, 100) ***REMOVED***
				num := 0
				for i, handle := range e.vus ***REMOVED***
					if i < 50 ***REMOVED***
						assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
						num++
					***REMOVED*** else ***REMOVED***
						assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					***REMOVED***
					assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
				***REMOVED***
				assert.Equal(t, 50, num)
			***REMOVED***

			t.Run("Raise", func(t *testing.T) ***REMOVED***
				assert.NoError(t, e.SetVUs(100))
				assert.Equal(t, int64(100), e.GetVUs())
				if assert.Len(t, e.vus, 100) ***REMOVED***
					for i, handle := range e.vus ***REMOVED***
						assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
						if i < 50 ***REMOVED***
							assert.Equal(t, int64(i+1), handle.vu.(*lib.MiniRunnerVU).ID)
						***REMOVED*** else ***REMOVED***
							assert.Equal(t, int64(50+i+1), handle.vu.(*lib.MiniRunnerVU).ID)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestRealTimeAndSetupTeardownMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
	import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";
	import ***REMOVED*** sleep ***REMOVED*** from "k6";

	var counter = new Counter("test_counter");

	export function setup() ***REMOVED***
		console.log("setup(), sleeping for 1 second");
		counter.add(1, ***REMOVED*** place: "setupBeforeSleep" ***REMOVED***);
		sleep(1);
		console.log("setup sleep is done");
		counter.add(2, ***REMOVED*** place: "setupAfterSleep" ***REMOVED***);
		return ***REMOVED*** "some": ["data"], "v": 1 ***REMOVED***;
	***REMOVED***

	export function teardown(data) ***REMOVED***
		console.log("teardown(" + JSON.stringify(data) + "), sleeping for 1 second");
		counter.add(3, ***REMOVED*** place: "teardownBeforeSleep" ***REMOVED***);
		sleep(1);
		if (!data || data.v != 1) ***REMOVED***
			throw new Error("incorrect data: " + JSON.stringify(data));
		***REMOVED***
		console.log("teardown sleep is done");
		counter.add(4, ***REMOVED*** place: "teardownAfterSleep" ***REMOVED***);
	***REMOVED***

	export default function (data) ***REMOVED***
		console.log("default(" + JSON.stringify(data) + ") with ENV=" + JSON.stringify(__ENV) + " for in ITER " + __ITER + " and VU " + __VU);
		counter.add(5, ***REMOVED*** place: "defaultBeforeSleep" ***REMOVED***);
		if (!data || data.v != 1) ***REMOVED***
			throw new Error("incorrect data: " + JSON.stringify(data));
		***REMOVED***
		sleep(1);
		console.log("default() for in ITER " + __ITER + " and VU " + __VU + " done!");
		counter.add(6, ***REMOVED*** place: "defaultAfterSleep" ***REMOVED***);
	***REMOVED***`)

	runner, err := js.New(
		&lib.SourceData***REMOVED***Filename: "/script.js", Data: script***REMOVED***,
		afero.NewMemMapFs(),
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	options := lib.Options***REMOVED***
		SystemTags:      lib.GetTagSet(lib.DefaultSystemTagList...),
		SetupTimeout:    types.NullDurationFrom(4 * time.Second),
		TeardownTimeout: types.NullDurationFrom(4 * time.Second),
	***REMOVED***
	runner.SetOptions(options)

	executor := New(runner)
	executor.SetEndIterations(null.IntFrom(2))
	executor.SetVUsMax(1)
	executor.SetVUs(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct***REMOVED******REMOVED***)
	sampleContainers := make(chan stats.SampleContainer)
	go func() ***REMOVED***
		assert.NoError(t, executor.Run(ctx, sampleContainers))
		close(done)
	***REMOVED***()

	expectIn := func(from, to time.Duration, expected stats.SampleContainer) ***REMOVED***
		start := time.Now()
		from = from * time.Millisecond
		to = to * time.Millisecond
		for ***REMOVED***
			select ***REMOVED***
			case sampleContainer := <-sampleContainers:
				now := time.Now()
				elapsed := now.Sub(start)
				if elapsed < from ***REMOVED***
					t.Errorf("Received sample earlier (%s) than expected (%s)", elapsed, from)
					return
				***REMOVED***
				assert.IsType(t, expected, sampleContainer)
				expSamples := expected.GetSamples()
				gotSamples := sampleContainer.GetSamples()
				if assert.Len(t, gotSamples, len(expSamples)) ***REMOVED***
					for i, s := range gotSamples ***REMOVED***
						expS := expSamples[i]
						if s.Metric != metrics.IterationDuration ***REMOVED***
							assert.Equal(t, expS.Value, s.Value)
						***REMOVED***
						assert.Equal(t, expS.Metric.Name, s.Metric.Name)
						assert.Equal(t, expS.Tags.CloneTags(), s.Tags.CloneTags())
						assert.InDelta(t, 0, now.Sub(s.Time), float64(50*time.Millisecond))
					***REMOVED***
				***REMOVED***
				return
			case <-time.After(to):
				t.Errorf("Did not receive sample in the maximum alotted time (%s)", to)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	getTags := func(args ...string) *stats.SampleTags ***REMOVED***
		tags := map[string]string***REMOVED******REMOVED***
		for i := 0; i < len(args)-1; i += 2 ***REMOVED***
			tags[args[i]] = args[i+1]
		***REMOVED***
		return stats.IntoSampleTags(&tags)
	***REMOVED***
	testCounter := stats.New("test_counter", stats.Counter)
	getSample := func(expValue float64, expMetric *stats.Metric, expTags ...string) stats.SampleContainer ***REMOVED***
		return stats.Sample***REMOVED***
			Metric: expMetric,
			Time:   time.Now(),
			Tags:   getTags(expTags...),
			Value:  expValue,
		***REMOVED***
	***REMOVED***
	getDummyTrail := func(group string) stats.SampleContainer ***REMOVED***
		return netext.NewDialer(net.Dialer***REMOVED******REMOVED***).GetTrail(time.Now(), time.Now(), getTags("group", group))
	***REMOVED***

	// Initially give a long time (5s) for the executor to start
	expectIn(0, 5000, getSample(1, testCounter, "group", "::setup", "place", "setupBeforeSleep"))
	expectIn(900, 1100, getSample(2, testCounter, "group", "::setup", "place", "setupAfterSleep"))
	expectIn(0, 100, getDummyTrail("::setup"))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep"))
	expectIn(0, 100, getDummyTrail(""))
	expectIn(0, 100, getSample(1, metrics.Iterations))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep"))
	expectIn(0, 100, getDummyTrail(""))
	expectIn(0, 100, getSample(1, metrics.Iterations))

	expectIn(0, 1000, getSample(3, testCounter, "group", "::teardown", "place", "teardownBeforeSleep"))
	expectIn(900, 1100, getSample(4, testCounter, "group", "::teardown", "place", "teardownAfterSleep"))
	expectIn(0, 100, getDummyTrail("::teardown"))

	for ***REMOVED***
		select ***REMOVED***
		case s := <-sampleContainers:
			t.Fatalf("Did not expect anything in the sample channel bug got %#v", s)
		case <-time.After(3 * time.Second):
			t.Fatalf("Local executor took way to long to finish")
		case <-done:
			return // Exit normally
		***REMOVED***
	***REMOVED***
***REMOVED***
