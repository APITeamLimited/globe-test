package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

func newExecutionSegmentFromString(str string) *libWorker.ExecutionSegment ***REMOVED***
	r, err := libWorker.NewExecutionSegmentFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r
***REMOVED***

func newExecutionSegmentSequenceFromString(str string) *libWorker.ExecutionSegmentSequence ***REMOVED***
	r, err := libWorker.NewExecutionSegmentSequenceFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return &r
***REMOVED***

func getTestConstantArrivalRateConfig() *ConstantArrivalRateConfig ***REMOVED***
	return &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(1 * time.Second)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(50),
		Duration:        types.NullDurationFrom(5 * time.Second),
		PreAllocatedVUs: null.IntFrom(10),
		MaxVUs:          null.IntFrom(20),
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunNotEnoughAllocatedVUsWarn(t *testing.T) ***REMOVED***
	t.Parallel()

	runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
		time.Sleep(time.Second)
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, getTestConstantArrivalRateConfig())
	defer test.cancel()

	engineOut := make(chan workerMetrics.SampleContainer, 1000)
	require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
	entries := test.logHook.Drain()
	require.NotEmpty(t, entries)
	for _, entry := range entries ***REMOVED***
		require.Equal(t,
			"Insufficient VUs, reached 20 active VUs and cannot initialize more",
			entry.Message)
		require.Equal(t, logrus.WarnLevel, entry.Level)
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateRunCorrectRate(t *testing.T) ***REMOVED***
	t.Parallel()

	var count int64
	runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
		atomic.AddInt64(&count, 1)
		return nil
	***REMOVED***)

	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, getTestConstantArrivalRateConfig())
	defer test.cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		// check that we got around the amount of VU iterations as we would expect
		var currentCount int64

		for i := 0; i < 5; i++ ***REMOVED***
			time.Sleep(time.Second)
			currentCount = atomic.SwapInt64(&count, 0)
			require.InDelta(t, 50, currentCount, 1)
		***REMOVED***
	***REMOVED***()
	engineOut := make(chan workerMetrics.SampleContainer, 1000)
	require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
	wg.Wait()
	require.Empty(t, test.logHook.Drain())
***REMOVED***

//nolint:tparallel,paralleltest // this is flaky if ran with other tests
func TestConstantArrivalRateRunCorrectTiming(t *testing.T) ***REMOVED***
	// t.Parallel()
	tests := []struct ***REMOVED***
		segment  string
		sequence string
		start    time.Duration
		steps    []int64
	***REMOVED******REMOVED***
		***REMOVED***
			segment: "0:1/3",
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 60, 60, 60, 60, 60, 60***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: "1/3:2/3",
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: "2/3:1",
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 60, 60, 60, 60, 60, 60***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment: "1/6:3/6",
			start:   time.Millisecond * 20,
			steps:   []int64***REMOVED***40, 80, 40, 80, 40, 80, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  "1/6:3/6",
			sequence: "1/6,3/6",
			start:    time.Millisecond * 20,
			steps:    []int64***REMOVED***40, 80, 40, 80, 40, 80, 40***REMOVED***,
		***REMOVED***,
		// sequences
		***REMOVED***
			segment:  "0:1/3",
			sequence: "0,1/3,2/3,1",
			start:    time.Millisecond * 0,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  "1/3:2/3",
			sequence: "0,1/3,2/3,1",
			start:    time.Millisecond * 20,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 60, 40***REMOVED***,
		***REMOVED***,
		***REMOVED***
			segment:  "2/3:1",
			sequence: "0,1/3,2/3,1",
			start:    time.Millisecond * 40,
			steps:    []int64***REMOVED***60, 60, 60, 60, 60, 100***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		test := test

		t.Run(fmt.Sprintf("segment %s sequence %s", test.segment, test.sequence), func(t *testing.T) ***REMOVED***
			t.Parallel()

			var count int64
			startTime := time.Now()
			expectedTimeInt64 := int64(test.start)
			runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
				current := atomic.AddInt64(&count, 1)

				expectedTime := test.start
				if current != 1 ***REMOVED***
					expectedTime = time.Duration(atomic.AddInt64(&expectedTimeInt64,
						int64(time.Millisecond)*test.steps[(current-2)%int64(len(test.steps))]))
				***REMOVED***

				// FIXME: replace this check with a unit test asserting that the scheduling is correct,
				// without depending on the execution time itself
				assert.WithinDuration(t,
					startTime.Add(expectedTime),
					time.Now(),
					time.Millisecond*24,
					"%d expectedTime %s", current, expectedTime,
				)

				return nil
			***REMOVED***)

			config := getTestConstantArrivalRateConfig()
			seconds := 2
			config.Duration.Duration = types.Duration(time.Second * time.Duration(seconds))
			execTest := setupExecutorTest(
				t, test.segment, test.sequence, libWorker.Options***REMOVED******REMOVED***, runner, config,
			)
			defer execTest.cancel()

			newET, err := execTest.state.ExecutionTuple.GetNewExecutionTupleFromValue(config.MaxVUs.Int64)
			require.NoError(t, err)
			rateScaled := newET.ScaleInt64(config.Rate.Int64)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				// check that we got around the amount of VU iterations as we would expect
				var currentCount int64

				for i := 0; i < seconds; i++ ***REMOVED***
					time.Sleep(time.Second)
					currentCount = atomic.LoadInt64(&count)
					assert.InDelta(t, int64(i+1)*rateScaled, currentCount, 3)
				***REMOVED***
			***REMOVED***()
			startTime = time.Now()
			engineOut := make(chan workerMetrics.SampleContainer, 1000)
			err = execTest.executor.Run(execTest.ctx, engineOut, libWorker.GetTestWorkerInfo())
			wg.Wait()
			require.NoError(t, err)
			require.Empty(t, execTest.logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArrivalRateCancel(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := map[string]libWorker.ExecutorConfig***REMOVED***
		"constant": getTestConstantArrivalRateConfig(),
		"ramping":  getTestRampingArrivalRateConfig(),
	***REMOVED***
	for name, config := range testCases ***REMOVED***
		config := config
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			ch := make(chan struct***REMOVED******REMOVED***)
			errCh := make(chan error, 1)
			weAreDoneCh := make(chan struct***REMOVED******REMOVED***)

			runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
				select ***REMOVED***
				case <-ch:
					<-ch
				default:
				***REMOVED***
				return nil
			***REMOVED***)
			test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, config)
			defer test.cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()

				engineOut := make(chan workerMetrics.SampleContainer, 1000)
				errCh <- test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo())
				close(weAreDoneCh)
			***REMOVED***()

			time.Sleep(time.Second)
			ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			test.cancel()
			time.Sleep(time.Second)
			select ***REMOVED***
			case <-weAreDoneCh:
				t.Fatal("Run returned before all VU iterations were finished")
			default:
			***REMOVED***
			close(ch)
			<-weAreDoneCh
			wg.Wait()
			require.NoError(t, <-errCh)
			require.Empty(t, test.logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConstantArrivalRateDroppedIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64

	config := &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(10),
		Duration:        types.NullDurationFrom(950 * time.Millisecond),
		PreAllocatedVUs: null.IntFrom(5),
		MaxVUs:          null.IntFrom(5),
	***REMOVED***

	runner := simpleRunner(func(ctx context.Context, _ *libWorker.State) error ***REMOVED***
		atomic.AddInt64(&count, 1)
		<-ctx.Done()
		return nil
	***REMOVED***)
	test := setupExecutorTest(t, "", "", libWorker.Options***REMOVED******REMOVED***, runner, config)
	defer test.cancel()

	engineOut := make(chan workerMetrics.SampleContainer, 1000)
	require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
	logs := test.logHook.Drain()
	require.Len(t, logs, 1)
	assert.Contains(t, logs[0].Message, "cannot initialize more")
	assert.Equal(t, int64(5), count)
	assert.Equal(t, float64(5), sumMetricValues(engineOut, workerMetrics.DroppedIterationsName))
***REMOVED***

func TestConstantArrivalRateGlobalIters(t *testing.T) ***REMOVED***
	t.Parallel()

	config := &ConstantArrivalRateConfig***REMOVED***
		BaseConfig:      BaseConfig***REMOVED***GracefulStop: types.NullDurationFrom(100 * time.Millisecond)***REMOVED***,
		TimeUnit:        types.NullDurationFrom(950 * time.Millisecond),
		Rate:            null.IntFrom(20),
		Duration:        types.NullDurationFrom(1 * time.Second),
		PreAllocatedVUs: null.IntFrom(5),
		MaxVUs:          null.IntFrom(5),
	***REMOVED***

	testCases := []struct ***REMOVED***
		seq, seg string
		expIters []uint64
	***REMOVED******REMOVED***
		***REMOVED***"0,1/4,3/4,1", "0:1/4", []uint64***REMOVED***1, 6, 11, 16, 21***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "1/4:3/4", []uint64***REMOVED***0, 2, 4, 5, 7, 9, 10, 12, 14, 15, 17, 19, 20***REMOVED******REMOVED***,
		***REMOVED***"0,1/4,3/4,1", "3/4:1", []uint64***REMOVED***3, 8, 13, 18***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%s_%s", tc.seq, tc.seg), func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotIters := []uint64***REMOVED******REMOVED***
			var mx sync.Mutex
			runner := simpleRunner(func(ctx context.Context, state *libWorker.State) error ***REMOVED***
				mx.Lock()
				gotIters = append(gotIters, state.GetScenarioGlobalVUIter())
				mx.Unlock()
				return nil
			***REMOVED***)
			test := setupExecutorTest(t, tc.seg, tc.seq, libWorker.Options***REMOVED******REMOVED***, runner, config)
			defer test.cancel()

			engineOut := make(chan workerMetrics.SampleContainer, 100)
			require.NoError(t, test.executor.Run(test.ctx, engineOut, libWorker.GetTestWorkerInfo()))
			assert.Equal(t, tc.expIters, gotIters)
		***REMOVED***)
	***REMOVED***
***REMOVED***
