package executor

import (
	"context"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func testConstantArrivalRateSetup(t *testing.T, vuFn func(context.Context, chan<- stats.SampleContainer) error) (
	context.Context, context.CancelFunc, lib.Executor, *testutils.SimpleLogrusHook) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	var config = ConstantArrivalRateConfig***REMOVED***
		TimeUnit:        types.NullDurationFrom(time.Second),
		Rate:            null.IntFrom(50),
		Duration:        types.NullDurationFrom(5 * time.Second),
		PreAllocatedVUs: null.IntFrom(10),
		MaxVUs:          null.IntFrom(20),
	***REMOVED***
	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)
	es := lib.NewExecutionState(lib.Options***REMOVED******REMOVED***, 10, 50)
	runner := lib.MiniRunner***REMOVED***
		Fn: vuFn,
	***REMOVED***

	es.SetInitVUFunc(func(_ context.Context, _ *logrus.Entry) (lib.VU, error) ***REMOVED***
		return &lib.MiniRunnerVU***REMOVED***R: runner***REMOVED***, nil
	***REMOVED***)

	initializeVUs(ctx, t, logEntry, es, 10)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)
	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
***REMOVED***

func TestConstantArrivalRateRunNotEnoughAllocatedVUsWarn(t *testing.T) ***REMOVED***
	t.Parallel()
	var ctx, cancel, executor, logHook = testConstantArrivalRateSetup(
		t, func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(time.Second)
			return nil
		***REMOVED***)
	defer cancel()
	var engineOut = make(chan stats.SampleContainer, 1000)
	err := executor.Run(ctx, engineOut)
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

func TestConstantArrivalRateRunCorrectRate(t *testing.T) ***REMOVED***
	t.Parallel()
	var count int64
	var ctx, cancel, executor, logHook = testConstantArrivalRateSetup(
		t, func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			atomic.AddInt64(&count, 1)
			return nil
		***REMOVED***)
	defer cancel()
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
	var engineOut = make(chan stats.SampleContainer, 1000)
	err := executor.Run(ctx, engineOut)
	wg.Wait()
	require.NoError(t, err)
	require.Empty(t, logHook.Drain())
***REMOVED***

func TestArrivalRateCancel(t *testing.T) ***REMOVED***
	t.Parallel()
	var mat = map[string]func(
		t *testing.T,
		vuFn func(context.Context, chan<- stats.SampleContainer) error,
	) (context.Context, context.CancelFunc, lib.Executor, *testutils.SimpleLogrusHook)***REMOVED***
		"constant": testConstantArrivalRateSetup,
		"variable": testVariableArrivalRateSetup,
	***REMOVED***
	for name, fn := range mat ***REMOVED***
		fn := fn
		t.Run(name, func(t *testing.T) ***REMOVED***
			var ch = make(chan struct***REMOVED******REMOVED***)
			var errCh = make(chan error, 1)
			var weAreDoneCh = make(chan struct***REMOVED******REMOVED***)
			var ctx, cancel, executor, logHook = fn(
				t, func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					select ***REMOVED***
					case <-ch:
						<-ch
					default:
					***REMOVED***
					return nil
				***REMOVED***)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()

				var engineOut = make(chan stats.SampleContainer, 1000)
				errCh <- executor.Run(ctx, engineOut)
				close(weAreDoneCh)
			***REMOVED***()

			time.Sleep(time.Second)
			ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			cancel()
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
			require.Empty(t, logHook.Drain())
		***REMOVED***)
	***REMOVED***
***REMOVED***
