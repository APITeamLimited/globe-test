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
	"errors"
	"net"
	"net/url"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/executor"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func newTestExecutionScheduler(
	t *testing.T, runner lib.Runner, logger *logrus.Logger, opts lib.Options, //nolint: golint
) (ctx context.Context, cancel func(), execScheduler *ExecutionScheduler, samples chan stats.SampleContainer) ***REMOVED***
	if runner == nil ***REMOVED***
		runner = &testutils.MiniRunner***REMOVED******REMOVED***
	***REMOVED***
	ctx, cancel = context.WithCancel(context.Background())
	newOpts, err := executor.DeriveExecutionFromShortcuts(lib.Options***REMOVED***
		MetricSamplesBufferSize: null.NewInt(200, false),
	***REMOVED***.Apply(runner.GetOptions()).Apply(opts))
	require.NoError(t, err)
	require.Empty(t, newOpts.Validate())

	require.NoError(t, runner.SetOptions(newOpts))

	if logger == nil ***REMOVED***
		logger = logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
	***REMOVED***

	execScheduler, err = NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	samples = make(chan stats.SampleContainer, newOpts.MetricSamplesBufferSize.Int64)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-samples:
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	require.NoError(t, execScheduler.Init(ctx, samples))

	return ctx, cancel, execScheduler, samples
***REMOVED***

func TestExecutionSchedulerRun(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, nil, nil, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	err := make(chan error, 1)
	go func() ***REMOVED*** err <- execScheduler.Run(ctx, samples) ***REMOVED***()
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutionSchedulerSetupTeardownRun(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Normal", func(t *testing.T) ***REMOVED***
		setupC := make(chan struct***REMOVED******REMOVED***)
		teardownC := make(chan struct***REMOVED******REMOVED***)
		runner := &testutils.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				close(setupC)
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				close(teardownC)
				return nil
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)

		err := make(chan error, 1)
		go func() ***REMOVED*** err <- execScheduler.Run(ctx, samples) ***REMOVED***()
		defer cancel()
		<-setupC
		<-teardownC
		assert.NoError(t, <-err)
	***REMOVED***)
	t.Run("Setup Error", func(t *testing.T) ***REMOVED***
		runner := &testutils.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, errors.New("setup error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)
		defer cancel()
		assert.EqualError(t, execScheduler.Run(ctx, samples), "setup error")
	***REMOVED***)
	t.Run("Don't Run Setup", func(t *testing.T) ***REMOVED***
		runner := &testutils.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, errors.New("setup error")
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			NoSetup:    null.BoolFrom(true),
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()
		assert.EqualError(t, execScheduler.Run(ctx, samples), "teardown error")
	***REMOVED***)

	t.Run("Teardown Error", func(t *testing.T) ***REMOVED***
		runner := &testutils.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()

		assert.EqualError(t, execScheduler.Run(ctx, samples), "teardown error")
	***REMOVED***)
	t.Run("Don't Run Teardown", func(t *testing.T) ***REMOVED***
		runner := &testutils.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			NoTeardown: null.BoolFrom(true),
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()
		assert.NoError(t, execScheduler.Run(ctx, samples))
	***REMOVED***)
***REMOVED***

func TestExecutionSchedulerStages(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		Duration time.Duration
		Stages   []lib.Stage
	***REMOVED******REMOVED***
		"one": ***REMOVED***
			1 * time.Second,
			[]lib.Stage***REMOVED******REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(1)***REMOVED******REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(1)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(2)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"four": ***REMOVED***
			4 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(3 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			runner := &testutils.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					time.Sleep(100 * time.Millisecond)
					return nil
				***REMOVED***,
			***REMOVED***
			ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
				VUs:    null.IntFrom(1),
				Stages: data.Stages,
			***REMOVED***)
			defer cancel()
			assert.NoError(t, execScheduler.Run(ctx, samples))
			assert.True(t, execScheduler.GetState().GetCurrentTestRunDuration() >= data.Duration)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerEndTime(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &testutils.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(100 * time.Millisecond)
			return nil
		***REMOVED***,
	***REMOVED***
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
		VUs:      null.IntFrom(10),
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 31*time.Second, endTime) // because of the default 30s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")
***REMOVED***

func TestExecutionSchedulerRuntimeErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &testutils.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(10 * time.Millisecond)
			return errors.New("hi")
		***REMOVED***,
		Options: lib.Options***REMOVED***
			VUs:      null.IntFrom(10),
			Duration: types.NullDurationFrom(1 * time.Second),
		***REMOVED***,
	***REMOVED***
	logger, hook := logtest.NewNullLogger()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 31*time.Second, endTime) // because of the default 30s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")

	assert.NotEmpty(t, hook.Entries)
	for _, e := range hook.Entries ***REMOVED***
		assert.Equal(t, "hi", e.Message)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerEndErrors(t *testing.T) ***REMOVED***
	t.Parallel()

	exec := executor.NewConstantLoopingVUsConfig("we_need_hard_stop")
	exec.VUs = null.IntFrom(10)
	exec.Duration = types.NullDurationFrom(1 * time.Second)
	exec.GracefulStop = types.NullDurationFrom(0 * time.Second)

	runner := &testutils.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			<-ctx.Done()
			return errors.New("hi")
		***REMOVED***,
		Options: lib.Options***REMOVED***
			Execution: lib.ExecutorConfigMap***REMOVED***exec.GetName(): exec***REMOVED***,
		***REMOVED***,
	***REMOVED***
	logger, hook := logtest.NewNullLogger()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 1*time.Second, endTime) // because of the 0s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")

	assert.Empty(t, hook.Entries)
***REMOVED***

func TestExecutionSchedulerEndIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	metric := &stats.Metric***REMOVED***Name: "test_metric"***REMOVED***

	options, err := executor.DeriveExecutionFromShortcuts(lib.Options***REMOVED***
		VUs:        null.IntFrom(1),
		Iterations: null.IntFrom(100),
	***REMOVED***)
	require.NoError(t, err)
	require.Empty(t, options.Validate())

	var i int64
	runner := &testutils.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
			default:
				atomic.AddInt64(&i, 1)
			***REMOVED***
			out <- stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***
			return nil
		***REMOVED***,
		Options: options,
	***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	samples := make(chan stats.SampleContainer, 300)
	require.NoError(t, execScheduler.Init(ctx, samples))
	require.NoError(t, execScheduler.Run(ctx, samples))

	assert.Equal(t, uint64(100), execScheduler.GetState().GetFullIterationCount())
	assert.Equal(t, uint64(0), execScheduler.GetState().GetPartialIterationCount())
	assert.Equal(t, int64(100), i)
	require.Equal(t, 100, len(samples)) //TODO: change to 200 https://github.com/loadimpact/k6/issues/1250
	for i := 0; i < 100; i++ ***REMOVED***
		mySample, ok := <-samples
		require.True(t, ok)
		assert.Equal(t, stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***, mySample)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerIsRunning(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &testutils.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			<-ctx.Done()
			return nil
		***REMOVED***,
	***REMOVED***
	ctx, cancel, execScheduler, _ := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)
	state := execScheduler.GetState()

	err := make(chan error)
	go func() ***REMOVED*** err <- execScheduler.Run(ctx, nil) ***REMOVED***()
	for !state.HasStarted() ***REMOVED***
		time.Sleep(10 * time.Microsecond)
	***REMOVED***
	cancel()
	for !state.HasEnded() ***REMOVED***
		time.Sleep(10 * time.Microsecond)
	***REMOVED***
	assert.NoError(t, <-err)
***REMOVED***

/*
//TODO: convert for the externally-controlled scheduler
func TestExecutionSchedulerSetVUs(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(-1), "vu count can't be negative")
	***REMOVED***)

	t.Run("Too High", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(100), "can't raise vu count (to 100) above vu cap (0)")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(&testutils.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
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
					assert.Equal(t, int64(0), handle.vu.(*testutils.MiniRunnerVU).ID)
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
					assert.Equal(t, int64(i+1), handle.vu.(*testutils.MiniRunnerVU).ID)
					num++
				***REMOVED*** else ***REMOVED***
					assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					assert.Equal(t, int64(0), handle.vu.(*testutils.MiniRunnerVU).ID)
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
				assert.Equal(t, int64(i+1), handle.vu.(*testutils.MiniRunnerVU).ID)
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
					assert.Equal(t, int64(i+1), handle.vu.(*testutils.MiniRunnerVU).ID)
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
							assert.Equal(t, int64(i+1), handle.vu.(*testutils.MiniRunnerVU).ID)
						***REMOVED*** else ***REMOVED***
							assert.Equal(t, int64(50+i+1), handle.vu.(*testutils.MiniRunnerVU).ID)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
*/

func TestRealTimeAndSetupTeardownMetrics(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip()
	***REMOVED***
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

	runner, err := js.New(&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***, nil, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	options, err := executor.DeriveExecutionFromShortcuts(runner.GetOptions().Apply(lib.Options***REMOVED***
		Iterations:      null.IntFrom(2),
		VUs:             null.IntFrom(1),
		SystemTags:      &stats.DefaultSystemTagSet,
		SetupTimeout:    types.NullDurationFrom(4 * time.Second),
		TeardownTimeout: types.NullDurationFrom(4 * time.Second),
	***REMOVED***))
	require.NoError(t, err)
	require.NoError(t, runner.SetOptions(options))

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct***REMOVED******REMOVED***)
	sampleContainers := make(chan stats.SampleContainer)
	go func() ***REMOVED***
		require.NoError(t, execScheduler.Init(ctx, sampleContainers))
		assert.NoError(t, execScheduler.Run(ctx, sampleContainers))
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
				t.Errorf("Did not receive sample in the maximum allotted time (%s)", to)
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
	getDummyTrail := func(group string, emitIterations bool) stats.SampleContainer ***REMOVED***
		return netext.NewDialer(net.Dialer***REMOVED******REMOVED***).GetTrail(time.Now(), time.Now(),
			true, emitIterations, getTags("group", group))
	***REMOVED***

	// Initially give a long time (5s) for the execScheduler to start
	expectIn(0, 5000, getSample(1, testCounter, "group", "::setup", "place", "setupBeforeSleep"))
	expectIn(900, 1100, getSample(2, testCounter, "group", "::setup", "place", "setupAfterSleep"))
	expectIn(0, 100, getDummyTrail("::setup", false))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep"))
	expectIn(0, 100, getDummyTrail("", true))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep"))
	expectIn(0, 100, getDummyTrail("", true))

	expectIn(0, 1000, getSample(3, testCounter, "group", "::teardown", "place", "teardownBeforeSleep"))
	expectIn(900, 1100, getSample(4, testCounter, "group", "::teardown", "place", "teardownAfterSleep"))
	expectIn(0, 100, getDummyTrail("::teardown", false))

	for ***REMOVED***
		select ***REMOVED***
		case s := <-sampleContainers:
			t.Fatalf("Did not expect anything in the sample channel bug got %#v", s)
		case <-time.After(3 * time.Second):
			t.Fatalf("Local execScheduler took way to long to finish")
		case <-done:
			return // Exit normally
		***REMOVED***
	***REMOVED***
***REMOVED***

// Just a lib.PausableExecutor implementation that can return an error
type pausableExecutor struct ***REMOVED***
	lib.Executor
	err error
***REMOVED***

func (p pausableExecutor) SetPaused(bool) error ***REMOVED***
	return p.err
***REMOVED***

func TestSetPaused(t *testing.T) ***REMOVED***
	t.Run("second pause is an error", func(t *testing.T) ***REMOVED***
		var runner = &testutils.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		var sched, err = NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***

		require.NoError(t, sched.SetPaused(true))
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution is already paused")
	***REMOVED***)

	t.Run("unpause at the start is an error", func(t *testing.T) ***REMOVED***
		var runner = &testutils.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		var sched, err = NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***
		err = sched.SetPaused(false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution wasn't paused")
	***REMOVED***)

	t.Run("second unpause is an error", func(t *testing.T) ***REMOVED***
		var runner = &testutils.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		var sched, err = NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***
		require.NoError(t, sched.SetPaused(true))
		require.NoError(t, sched.SetPaused(false))
		err = sched.SetPaused(false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution wasn't paused")
	***REMOVED***)

	t.Run("an error on pausing is propagated", func(t *testing.T) ***REMOVED***
		var runner = &testutils.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		var sched, err = NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		var expectedErr = errors.New("testing pausable executor error")
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: expectedErr***REMOVED******REMOVED***
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Equal(t, err, expectedErr)
	***REMOVED***)

	t.Run("can't pause unpausable executor", func(t *testing.T) ***REMOVED***
		var runner = &testutils.MiniRunner***REMOVED******REMOVED***
		options, err := executor.DeriveExecutionFromShortcuts(lib.Options***REMOVED***
			Iterations: null.IntFrom(2),
			VUs:        null.IntFrom(1),
		***REMOVED***.Apply(runner.GetOptions()))
		require.NoError(t, err)
		require.NoError(t, runner.SetOptions(options))

		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Contains(t, err.Error(), "doesn't support pause and resume operations after its start")
	***REMOVED***)
***REMOVED***
