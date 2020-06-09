package executor

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/testutils/minirunner"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		return &minirunner.VU***REMOVED***
			R: &minirunner.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					// TODO: do something
					return nil
				***REMOVED***,
			***REMOVED***,
		***REMOVED***, nil
	***REMOVED***

	returnVU := func(_ lib.InitializedVU) ***REMOVED***
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

	vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, &BaseConfig***REMOVED******REMOVED***, logEntry)
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
	time.Sleep(time.Millisecond)
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

func TestVUHandleSimple(t *testing.T) ***REMOVED***
	t.Parallel()

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.DebugLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(testutils.NewTestOutput(t))
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

	getVU := func() (lib.InitializedVU, error) ***REMOVED*** //nolint:unparam
		atomic.AddUint32(&getVUCount, 1)

		return &minirunner.VU***REMOVED***
			R: &minirunner.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					// TODO: do something
					return nil
				***REMOVED***,
			***REMOVED***,
		***REMOVED***, nil
	***REMOVED***

	returnVU := func(_ lib.InitializedVU) ***REMOVED***
		atomic.AddUint32(&returnVUCount, 1)
	***REMOVED***

	runIter := func(ctx context.Context, _ lib.ActiveVU) bool ***REMOVED***
		select ***REMOVED***
		case <-time.After(time.Second):
		case <-ctx.Done():
		***REMOVED***

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
	t.Run("start before gracefulStop finishes", func(t *testing.T) ***REMOVED***
		reset()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)
		vuHandle.gracefulStop()
		time.Sleep(time.Millisecond * 5)
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		assert.EqualValues(t, 1, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 0, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 0, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&fullIterations))
		cancel()
		wg.Wait()
		time.Sleep(time.Millisecond * 5)
		assert.EqualValues(t, 1, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&fullIterations))
	***REMOVED***)

	t.Run("start after gracefulStop finishes", func(t *testing.T) ***REMOVED***
		reset()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)
		vuHandle.gracefulStop()
		time.Sleep(time.Millisecond * 1500)
		assert.EqualValues(t, 1, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 0, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&fullIterations))
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		cancel()
		wg.Wait()

		time.Sleep(time.Millisecond * 5)
		assert.EqualValues(t, 2, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 2, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 2, atomic.LoadInt64(&fullIterations))
	***REMOVED***)

	t.Run("start after hardStop", func(t *testing.T) ***REMOVED***
		reset()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, &BaseConfig***REMOVED******REMOVED***, logEntry)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			vuHandle.runLoopsIfPossible(runIter)
		***REMOVED***()
		err := vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)
		vuHandle.hardStop()
		time.Sleep(time.Millisecond * 15)
		assert.EqualValues(t, 1, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 1, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 1, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 0, atomic.LoadInt64(&fullIterations))
		err = vuHandle.start()
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 1500)
		cancel()
		wg.Wait()

		time.Sleep(time.Millisecond * 5)
		assert.EqualValues(t, 2, atomic.LoadUint32(&getVUCount))
		assert.EqualValues(t, 2, atomic.LoadUint32(&returnVUCount))
		assert.EqualValues(t, 2, atomic.LoadInt64(&interruptedIter))
		assert.EqualValues(t, 1, atomic.LoadInt64(&fullIterations))
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

	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		atomic.AddUint32(&getVUCount, 1)

		return &minirunner.VU***REMOVED***
			R: &minirunner.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					// TODO: do something
					return nil
				***REMOVED***,
			***REMOVED***,
		***REMOVED***, nil
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
	vuHandle := newStoppedVUHandle(ctx, getVU, returnVU, &BaseConfig***REMOVED******REMOVED***, logEntry)
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
