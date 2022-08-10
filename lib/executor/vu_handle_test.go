package executor

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/metrics"
)

func mockNextIterations() (uint64, uint64) ***REMOVED***
	return 12, 15
***REMOVED***

// this test is mostly interesting when -race is enabled
func TestVUHandleRace(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(testutils.NewTestOutput(t))
	// testLog.Level = logrus.DebugLevel
	logEntry := logrus.NewEntry(testLog)

	runner := &minirunner.MiniRunner***REMOVED******REMOVED***
	runner.Fn = func(ctx context.Context, _ *lib.State, out chan<- metrics.SampleContainer) error ***REMOVED***
		return nil
	***REMOVED***

	var getVUCount int64
	var returnVUCount int64
	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		return runner.NewVU(uint64(atomic.AddInt64(&getVUCount, 1)), 0, nil)
	***REMOVED***

	returnVU := func(_ lib.InitializedVU) ***REMOVED***
		atomic.AddInt64(&returnVUCount, 1)
		// do something
	***REMOVED***
	var interruptedIter int64
	var fullIterations int64

	runIter := func(ctx context.Context, vu lib.ActiveVU) bool ***REMOVED***
		_ = vu.RunOnce()
		select ***REMOVED***
		case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			atomic.AddInt64(&interruptedIter, 1)
			return false
		default:
			atomic.AddInt64(&fullIterations, 1)
			return true
		***REMOVED***
	***REMOVED***

	vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
	go vuHandle.runLoopsIfPossible(runIter)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() ***REMOVED***
		defer wg.Done()
		for i := 0; i < 10000; i++ ***REMOVED***
			err := vuHandle.start()
			require.NoError(t, err)
		***REMOVED***
	***REMOVED***()

	go func() ***REMOVED***
		defer wg.Done()
		for i := 0; i < 1000; i++ ***REMOVED***
			vuHandle.gracefulStop()
			time.Sleep(1 * time.Nanosecond)
		***REMOVED***
	***REMOVED***()

	go func() ***REMOVED***
		defer wg.Done()
		for i := 0; i < 100; i++ ***REMOVED***
			vuHandle.hardStop()
			time.Sleep(10 * time.Nanosecond)
		***REMOVED***
	***REMOVED***()
	wg.Wait()
	vuHandle.hardStop() // STOP it
	time.Sleep(time.Millisecond * 50)
	interruptedBefore := atomic.LoadInt64(&interruptedIter)
	fullBefore := atomic.LoadInt64(&fullIterations)
	_ = vuHandle.start()
	time.Sleep(time.Millisecond * 50) // just to be sure an iteration will squeeze in
	cancel()
	time.Sleep(time.Millisecond * 50)
	interruptedAfter := atomic.LoadInt64(&interruptedIter)
	fullAfter := atomic.LoadInt64(&fullIterations)
	assert.True(t, interruptedBefore >= interruptedAfter-1,
		"too big of a difference %d >= %d - 1", interruptedBefore, interruptedAfter)
	assert.True(t, fullBefore+1 <= fullAfter,
		"too small of a difference %d + 1 <= %d", fullBefore, fullAfter)
	require.Equal(t, atomic.LoadInt64(&getVUCount), atomic.LoadInt64(&returnVUCount))
***REMOVED***

// this test is mostly interesting when -race is enabled
func TestVUHandleStartStopRace(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(testutils.NewTestOutput(t))
	// testLog.Level = logrus.DebugLevel
	logEntry := logrus.NewEntry(testLog)

	runner := &minirunner.MiniRunner***REMOVED******REMOVED***
	runner.Fn = func(ctx context.Context, _ *lib.State, out chan<- metrics.SampleContainer) error ***REMOVED***
		return nil
	***REMOVED***

	var vuID uint64
	testIterations := 10000
	returned := make(chan struct***REMOVED******REMOVED***)

	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		returned = make(chan struct***REMOVED******REMOVED***)
		return runner.NewVU(atomic.AddUint64(&vuID, 1), 0, nil)
	***REMOVED***

	returnVU := func(v lib.InitializedVU) ***REMOVED***
		require.Equal(t, atomic.LoadUint64(&vuID), v.(*minirunner.VU).ID)
		close(returned)
	***REMOVED***
	var interruptedIter int64
	var fullIterations int64

	runIter := func(ctx context.Context, vu lib.ActiveVU) bool ***REMOVED***
		_ = vu.RunOnce()
		select ***REMOVED***
		case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			atomic.AddInt64(&interruptedIter, 1)
			return false
		default:
			atomic.AddInt64(&fullIterations, 1)
			return true
		***REMOVED***
	***REMOVED***

	vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
	go vuHandle.runLoopsIfPossible(runIter)
	for i := 0; i < testIterations; i++ ***REMOVED***
		err := vuHandle.start()
		vuHandle.gracefulStop()
		require.NoError(t, err)
		select ***REMOVED***
		case <-returned:
		case <-time.After(100 * time.Millisecond):
			go panic("returning took too long")
			time.Sleep(time.Second)
		***REMOVED***
	***REMOVED***

	vuHandle.hardStop() // STOP it
	time.Sleep(time.Millisecond * 5)
	interruptedBefore := atomic.LoadInt64(&interruptedIter)
	fullBefore := atomic.LoadInt64(&fullIterations)
	_ = vuHandle.start()
	time.Sleep(time.Millisecond * 50) // just to be sure an iteration will squeeze in
	cancel()
	time.Sleep(time.Millisecond * 5)
	interruptedAfter := atomic.LoadInt64(&interruptedIter)
	fullAfter := atomic.LoadInt64(&fullIterations)
	assert.True(t, interruptedBefore >= interruptedAfter-1,
		"too big of a difference %d >= %d - 1", interruptedBefore, interruptedAfter)
	assert.True(t, fullBefore+1 <= fullAfter,
		"too small of a difference %d + 1 <= %d", fullBefore, fullAfter)
***REMOVED***

type handleVUTest struct ***REMOVED***
	runner          *minirunner.MiniRunner
	getVUCount      uint32
	returnVUCount   uint32
	interruptedIter int64
	fullIterations  int64
***REMOVED***

func (h *handleVUTest) getVU() (lib.InitializedVU, error) ***REMOVED***
	return h.runner.NewVU(uint64(atomic.AddUint32(&h.getVUCount, 1)), 0, nil)
***REMOVED***

func (h *handleVUTest) returnVU(_ lib.InitializedVU) ***REMOVED***
	atomic.AddUint32(&h.returnVUCount, 1)
***REMOVED***

func (h *handleVUTest) runIter(ctx context.Context, _ lib.ActiveVU) bool ***REMOVED***
	select ***REMOVED***
	case <-time.After(time.Second):
	case <-ctx.Done():
	***REMOVED***

	select ***REMOVED***
	case <-ctx.Done():
		// Don't log errors or emit iterations metrics from cancelled iterations
		atomic.AddInt64(&h.interruptedIter, 1)
		return false
	default:
		atomic.AddInt64(&h.fullIterations, 1)
		return true
	***REMOVED***
***REMOVED***

func TestVUHandleSimple(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("start before gracefulStop finishes", func(t *testing.T) ***REMOVED***
		t.Parallel()
		logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
		testLog := logrus.New()
		testLog.AddHook(logHook)
		testLog.SetOutput(testutils.NewTestOutput(t))
		// testLog.Level = logrus.DebugLevel
		logEntry := logrus.NewEntry(testLog)
		test := &handleVUTest***REMOVED***runner: &minirunner.MiniRunner***REMOVED******REMOVED******REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		vuHandle := newStoppedVUHandle(ctx, test.getVU, test.returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(test.runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 50)
		vuHandle.gracefulStop()
		// time.Sleep(time.Millisecond * 5) // No sleep as we want to always not return the VU
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 0, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 0, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.fullIterations))
		cancel()
		wg.Wait()
		time.Sleep(time.Millisecond * 5)
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.fullIterations))
	***REMOVED***)

	t.Run("start after gracefulStop finishes", func(t *testing.T) ***REMOVED***
		t.Parallel()
		logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
		testLog := logrus.New()
		testLog.AddHook(logHook)
		testLog.SetOutput(testutils.NewTestOutput(t))
		// testLog.Level = logrus.DebugLevel
		logEntry := logrus.NewEntry(testLog)
		test := &handleVUTest***REMOVED***runner: &minirunner.MiniRunner***REMOVED******REMOVED******REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		vuHandle := newStoppedVUHandle(ctx, test.getVU, test.returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(test.runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 50)
		vuHandle.gracefulStop()
		time.Sleep(time.Millisecond * 1500)
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 0, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.fullIterations))
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		cancel()
		wg.Wait()

		time.Sleep(time.Millisecond * 50)
		assert.EqualValues(t, 2, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 2, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 2, atomic.LoadInt64(&test.fullIterations))
	***REMOVED***)

	t.Run("start after hardStop", func(t *testing.T) ***REMOVED***
		t.Parallel()
		logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
		testLog := logrus.New()
		testLog.AddHook(logHook)
		testLog.SetOutput(testutils.NewTestOutput(t))
		// testLog.Level = logrus.DebugLevel
		logEntry := logrus.NewEntry(testLog)
		test := &handleVUTest***REMOVED***runner: &minirunner.MiniRunner***REMOVED******REMOVED******REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		vuHandle := newStoppedVUHandle(ctx, test.getVU, test.returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(test.runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)
		vuHandle.hardStop()
		time.Sleep(time.Millisecond * 15)
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 0, atomic.LoadInt64(&test.fullIterations))
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		cancel()
		wg.Wait()

		time.Sleep(time.Millisecond * 5)
		assert.EqualValues(t, 2, atomic.LoadUint32(&test.getVUCount))
		assert.EqualValues(t, 2, atomic.LoadUint32(&test.returnVUCount))
		assert.EqualValues(t, 2, atomic.LoadInt64(&test.interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&test.fullIterations))
	***REMOVED***)
***REMOVED***

func BenchmarkVUHandleIterations(b *testing.B) ***REMOVED***
	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	// testLog.Level = logrus.DebugLevel
	logEntry := logrus.NewEntry(testLog)

	var (
		getVUCount      uint32
		returnVUCount   uint32
		interruptedIter int64
		fullIterations  int64
	)
	reset := func() ***REMOVED***
		getVUCount = 0
		returnVUCount = 0
		interruptedIter = 0
		fullIterations = 0
	***REMOVED***

	runner := &minirunner.MiniRunner***REMOVED******REMOVED***
	runner.Fn = func(ctx context.Context, _ *lib.State, out chan<- metrics.SampleContainer) error ***REMOVED***
		return nil
	***REMOVED***
	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		return runner.NewVU(uint64(atomic.AddUint32(&getVUCount, 1)), 0, nil)
	***REMOVED***

	returnVU := func(_ lib.InitializedVU) ***REMOVED***
		atomic.AddUint32(&returnVUCount, 1)
	***REMOVED***

	runIter := func(ctx context.Context, _ lib.ActiveVU) bool ***REMOVED***
		// Do nothing
		select ***REMOVED***
		case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			atomic.AddInt64(&interruptedIter, 1)
			return false
		default:
			atomic.AddInt64(&fullIterations, 1)
			return true
		***REMOVED***
	***REMOVED***

	reset()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, mockNextIterations, &BaseConfig***REMOVED******REMOVED***, logEntry)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		vuHandle.runLoopsIfPossible(runIter)
	***REMOVED***()
	start := time.Now()
	b.ResetTimer()
	err := vuHandle.start()
	require.NoError(b, err)
	time.Sleep(time.Second)
	cancel()
	wg.Wait()
	b.StopTimer()
	took := time.Since(start)
	b.ReportMetric(float64(atomic.LoadInt64(&fullIterations))/float64(took), "iterations/ns")
***REMOVED***
